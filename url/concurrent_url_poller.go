package main

import (
	"log"
	"net/http"
	"time"
)

type Resource struct {
	url string
	errCount int
}

type State struct {
	url string
	status string
}

const (
	numPollers = 2
	pollInterval = 60 * time.Second
	statusInterval = 10 * time.Second
	errTimeout = 10 * time.Second
)

var urls = []string {
	"http://www.google.com/1",
	"http://golang.org/",
	"http://blog.golang.org/",
}

func main() {
	pendingChan, completeChan := make(chan *Resource), make(chan *Resource)
	statusChan := StateMonitor(statusInterval)

	for i := 0; i < numPollers; i++ {
		go Poller(pendingChan, completeChan, statusChan)
	}

	go func() {
		for _, url := range urls {
			pendingChan <- &Resource{url: url}
		}
	}()

	for r := range completeChan{
		go r.Sleep(pendingChan)
	}
}

func StateMonitor(updateInterval time.Duration) chan <- State {
	// 状态更新通道
	updates := make(chan State)
	// status map
	urlStatus := make(map[string]string)

	ticker := time.NewTicker(updateInterval)

	go func() {
		for  {
			select {
			case <- ticker.C:
				logState(urlStatus)
			case s := <- updates:
				urlStatus[s.url] = s.status
			}
		}
	}()

	return updates
}

func logState(s map[string]string)  {
	log.Println("Current state:")
	for k, v := range s {
		log.Printf(" %s %s", k, v)
	}
}

func Poller(in <-chan *Resource, out chan<- *Resource, statusChan chan <- State)  {
	// 从通道中读取资源
	for r := range in {
		status := r.Poll()
		statusChan <- State{r.url, status}
		out <- r
	}
}

func (r *Resource) Poll() string {
	resp, err := http.Head(r.url)

	if err != nil {
		log.Println("Error", r.url, err)
		r.errCount++
		return err.Error()
	}

	r.errCount = 0
	return resp.Status
}

func (r *Resource) Sleep(doneChan chan<- *Resource)  {
	time.Sleep(pollInterval + errTimeout * time.Duration(r.errCount))

	doneChan <- r
}




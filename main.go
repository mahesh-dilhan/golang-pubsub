package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

func main() {
	srv1 := "srv1"
	srv2 := "srv2"

	srvchan1 := make(chan string)
	srvchan2 := make(chan string)

	srv3 := "srv3"
	srv4 := "srv4"

	srvchan3 := make(chan string)
	srvchan4 := make(chan string)

	m := map[string][]chan string{
		srv1: {srvchan1},
		srv2: {srvchan2},
		srv3: {srvchan3},
		srv4: {srvchan4},
	}

	var mux sync.RWMutex

	p := Pubsub{
		closed: false,
		subs:   m,
		mu:     mux,
	}

	for i := 1; i <= 4; i++ {
		go p.Publish("srv"+strconv.Itoa(i), "meg from service "+strconv.Itoa(i))
		fmt.Printf("receive msg '%v' \n ", <-p.Subscribe("srv"+strconv.Itoa(i)))
	}

	//s:=  p.Subscribe(srv1)
	//fmt.Printf("receive msg '%v' \n ",<-s)
	time.Sleep(5 * time.Second)

}

type Pubsub struct {
	mu     sync.RWMutex
	subs   map[string][]chan string
	closed bool
}

func (ps *Pubsub) Publish(topic string, msg string) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if ps.closed {
		return
	}
	fmt.Printf("pushing msg '%s' to topic '%s'\n", msg, topic)
	for _, ch := range ps.subs[topic] {
		go func(ch chan string) {
			ch <- msg
		}(ch)
	}
}

func (ps *Pubsub) Subscribe(topic string) <-chan string {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan string, 1)
	ps.subs[topic] = append(ps.subs[topic], ch)
	return ch
}

func (ps *Pubsub) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if !ps.closed {
		ps.closed = true
		for _, subs := range ps.subs {
			for _, ch := range subs {
				close(ch)
			}
		}
	}
}

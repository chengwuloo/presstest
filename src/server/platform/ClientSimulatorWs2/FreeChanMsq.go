package main

//
// Created by andy_ro@qq.com
// 			4/9/2019
//

import (
	"sync/atomic"
	"time"
)

//FreeChanMsq chan类型消息队列
type FreeChanMsq struct {
	Msgs chan interface{}
	n    int64
}

//
func newFreeChanMsq() MsgQueue {
	return &FreeChanMsq{Msgs: make(chan interface{}, 9000)}
}

//Push 入队列
func (s *FreeChanMsq) Push(msg interface{}) {
	s.Msgs <- msg
	atomic.AddInt64(&s.n, 1)
}

//Pop 出队列
func (s *FreeChanMsq) Pop() (msg interface{}, exit bool) {
	select {
	case q := <-s.Msgs:
		{
			if q == nil {
				close(s.Msgs)
				exit = true
				break
			} else {
				msg = q
			}
			atomic.AddInt64(&s.n, -1)
		}
	//case <-time.After(time.Nanosecond):
	//case <-time.After(time.Microsecond):
	case <-time.After(time.Millisecond):
		//log.Printf("--- *** ----------------------------- [%05d]Run time.After...\n", os.Getpid())
		//default:
	}
	return
}

//Pick 出队列
func (s *FreeChanMsq) Pick() (msgs []interface{}, exit bool) {
	msg, e := s.Pop()
	exit = e
	if msg != nil && !exit {
		msgs = append(msgs, msg)
	}
	return
}

//
func (s *FreeChanMsq) Count() int64 {
	return atomic.LoadInt64(&s.n)
}

//
func (s *FreeChanMsq) Signal() {
}

//
func (s *FreeChanMsq) EnableNonBlocking(bv bool) {
}

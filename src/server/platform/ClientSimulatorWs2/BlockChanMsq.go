package main

//
// Created by andy_ro@qq.com
// 			4/9/2019
//

import (
	"log"
	"server/platform/util"
	"sync/atomic"
	"time"
)

//
const (
	KBlocking int32 = iota
	KNonblocking
)

//BlockChanMsq chan类型消息队列
type BlockChanMsq struct {
	Msgs        chan interface{}
	signal      chan bool
	n           int64
	nonblocking int32
}

//
func newBlockChanMsq() MsgQueue {
	return &BlockChanMsq{Msgs: make(chan interface{}, 100), signal: make(chan bool, 1)}
}

//Push 入队列
func (s *BlockChanMsq) Push(msg interface{}) {
	s.Msgs <- msg
	atomic.AddInt64(&s.n, 1)
}

//
func (s *BlockChanMsq) blockPop() (msg interface{}, exit bool) {
	select {
	case q := <-s.Msgs:
		{
			if q == nil {
				close(s.Msgs)
				close(s.signal)
				exit = true
				break
			} else {
				msg = q
			}
			atomic.AddInt64(&s.n, -1)
		}
	case <-s.signal:
	}
	return
}

//
func (s *BlockChanMsq) nonblockPop() (msg interface{}, exit bool) {
	select {
	case q := <-s.Msgs:
		{
			if q == nil {
				close(s.Msgs)
				close(s.signal)
				exit = true
				break
			} else {
				msg = q
			}
			atomic.AddInt64(&s.n, -1)
		}
	case <-time.After(time.Microsecond):
	default:
	}
	return
}

//Pop 出队列
func (s *BlockChanMsq) Pop() (msg interface{}, exit bool) {
	if KBlocking == atomic.LoadInt32(&s.nonblocking) {
		msg, exit = s.blockPop()
	} else {
		msg, exit = s.nonblockPop()
	}
	return
}

//Pick 出队列
func (s *BlockChanMsq) Pick() (msgs []interface{}, exit bool) {
	msg, e := s.Pop()
	exit = e
	if msg != nil && !exit {
		msgs = append(msgs, msg)
	}
	return
}

//
func (s *BlockChanMsq) Count() int64 {
	return atomic.LoadInt64(&s.n)
}

//
func (s *BlockChanMsq) Signal() {
	if KBlocking == atomic.LoadInt32(&s.nonblocking) {
		s.signalx()
	}
}

//
func (s *BlockChanMsq) EnableNonBlocking(bv bool) {
	val := KBlocking
	if bv == true {
		val = KNonblocking
	}
	old := atomic.LoadInt32(&s.nonblocking)
	if old != val {
		atomic.StoreInt32(&s.nonblocking, val)
	}
	//阻塞变为非阻塞
	if old == KBlocking && val == KNonblocking {
		s.signalx()
	}
}

//
func (s *BlockChanMsq) signalx() {
	s.signal <- true
}

var msq = newBlockChanMsq()

func onInput1(str string) int {
	switch str {
	case "q":
		{
			msq.Push(nil)
			return -1
		}
	case "s":
		{
			msq.Signal()

		}
	case "b":
		{
			msq.EnableNonBlocking(false)
		}
	case "nb":
		{
			msq.EnableNonBlocking(true)
		}
	default:
		{

		}
	}
	return 0
}

func testChanMsq() {
	go func(msq MsgQueue) {
		for {
			msg, exit := msq.Pop()
			if exit == true {
				break
			}
			if msg != nil && !exit {
				log.Println("testChanMsq receive: ", msg)
			} else {
				log.Println("continue ...")
			}
		}
		log.Println("exit...")
	}(msq)
	for {
		util.ReadConsole(onInput1)
	}
}

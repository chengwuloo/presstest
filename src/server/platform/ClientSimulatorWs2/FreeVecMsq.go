package main

//
// Created by YangZhi
// 			4/9/2019
//

import (
	"sync"
	"sync/atomic"
)

//FreeVecMsq 非阻塞切片类型
type FreeVecMsq struct {
	Msgs []interface{}
	l    *sync.Mutex
	n    int64
}

//
func newFreeVecMsq() MsgQueue {
	s := &FreeVecMsq{l: &sync.Mutex{}}
	return s
}

//
func (s *FreeVecMsq) EnableNonBlocking(bv bool) {

}

//Push 入队列
func (s *FreeVecMsq) Push(msg interface{}) {
	{
		s.l.Lock()
		s.Msgs = append(s.Msgs, msg)
		s.l.Unlock()
		atomic.AddInt64(&s.n, 1)
	}
}

//Pop 出队列
func (s *FreeVecMsq) Pop() (msg interface{}, exit bool) {
	{
		s.l.Lock()
		if len(s.Msgs) > 0 {
			msg = s.Msgs[0]
			if msg == nil {
				exit = true
				s.reset()
			} else {
				s.Msgs = s.Msgs[1:]
				atomic.AddInt64(&s.n, -1)
			}
		}
		s.l.Unlock()
	}
	return
}

//Pick 出队列
func (s *FreeVecMsq) Pick() (msgs []interface{}, exit bool) {
	{
		s.l.Lock()
		for _, msg := range s.Msgs {
			if msg == nil {
				exit = true
				break
			} else {
				msgs = append(msgs, msg)
			}
		}
		s.reset()
		s.l.Unlock()
	}
	return
}

//
func (s *FreeVecMsq) Count() int64 {
	return atomic.LoadInt64(&s.n)
}

//
func (s *FreeVecMsq) Signal() {
}

//
func (s *FreeVecMsq) reset() {
	s.Msgs = s.Msgs[0:0]
	atomic.StoreInt64(&s.n, 0)
}

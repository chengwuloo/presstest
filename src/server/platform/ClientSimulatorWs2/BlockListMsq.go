package main

//
// Created by YangZhi
// 			4/9/2019
//

import (
	"container/list"
	"sync"
	"sync/atomic"
)

//BlockListMsq 链表类型
type BlockListMsq struct {
	Msgs        *list.List
	l           *sync.Mutex
	c           *sync.Cond
	n           int64
	nonblocking bool
}

//
func newBlockListMsq() MsgQueue {
	s := &BlockListMsq{Msgs: list.New(),
		l: &sync.Mutex{}}
	s.c = sync.NewCond(s.l)
	return s
}

//
func (s *BlockListMsq) EnableNonBlocking(bv bool) {
	if s.nonblocking != bv {
		s.nonblocking = bv
		// str := " FALSE"
		// if s.nonblocking {
		// 	str = " TRUE"
		// }
		// log.Println("NonBlocking: ", str)
	}
}

//Push 入队列
func (s *BlockListMsq) Push(msg interface{}) {
	{
		s.l.Lock()
		s.Msgs.PushBack(msg)
		s.l.Unlock()
	}
	atomic.AddInt64(&s.n, 1)
	s.c.Signal()
}

//Pop 出队列
func (s *BlockListMsq) Pop() (msg interface{}, exit bool) {
	{
		s.l.Lock()
		if !s.nonblocking && s.Msgs.Len() == 0 {
			s.c.Wait()
		}
		s.l.Unlock()
	}
	{
		s.l.Lock()
		if elem := s.Msgs.Front(); elem != nil {
			msg = elem.Value
			if msg == nil {
				exit = true
				s.reset()
			} else {
				s.Msgs.Remove(elem)
			}
			atomic.AddInt64(&s.n, -1)
		}
		s.l.Unlock()
	}
	return
}

//Pick 出队列
func (s *BlockListMsq) Pick() (msgs []interface{}, exit bool) {
	{
		s.l.Lock()
		if !s.nonblocking && s.Msgs.Len() == 0 {
			s.c.Wait()
		}
		s.l.Unlock()
	}
	{
		s.l.Lock()
		var next *list.Element
		for elem := s.Msgs.Front(); elem != nil; elem = next {
			next = elem.Next()
			msg := elem.Value
			s.Msgs.Remove(elem)
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
func (s *BlockListMsq) Count() int64 {
	return atomic.LoadInt64(&s.n)
}

//
func (s *BlockListMsq) Signal() {
	if !s.nonblocking {
		s.c.Signal()
	}
}

//
func (s *BlockListMsq) reset() {
	var next *list.Element
	for elem := s.Msgs.Front(); elem != nil; elem = next {
		next = elem.Next()
		s.Msgs.Remove(elem)
	}
	atomic.StoreInt64(&s.n, 0)
}

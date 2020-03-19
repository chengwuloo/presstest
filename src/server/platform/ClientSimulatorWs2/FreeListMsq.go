package main

import (
	"container/list"
	"sync"
	"sync/atomic"
)

//FreeListMsq 非阻塞链表类型
type FreeListMsq struct {
	Msgs *list.List
	l    *sync.Mutex
	n    int64
}

//
func newFreeListMsq() MsgQueue {
	s := &FreeListMsq{Msgs: list.New(),
		l: &sync.Mutex{}}
	return s
}

//
func (s *FreeListMsq) EnableNonBlocking(bv bool) {

}

//Push 入队列
func (s *FreeListMsq) Push(msg interface{}) {
	{
		s.l.Lock()
		s.Msgs.PushBack(msg)
		s.l.Unlock()
		atomic.AddInt64(&s.n, 1)
	}
}

//Pop 出队列
func (s *FreeListMsq) Pop() (msg interface{}, exit bool) {
	{
		s.l.Lock()
		if elem := s.Msgs.Front(); elem != nil {
			msg = elem.Value
			if msg == nil {
				exit = true
				s.reset()
			} else {
				atomic.AddInt64(&s.n, -1)
				s.Msgs.Remove(elem)
			}
		}
		s.l.Unlock()
	}
	return
}

//Pick 出队列
func (s *FreeListMsq) Pick() (msgs []interface{}, exit bool) {
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
func (s *FreeListMsq) Count() int64 {
	return atomic.LoadInt64(&s.n)
}

//
func (s *FreeListMsq) Signal() {
}

//
func (s *FreeListMsq) reset() {
	var next *list.Element
	for elem := s.Msgs.Front(); elem != nil; elem = next {
		next = elem.Next()
		s.Msgs.Remove(elem)
	}
	atomic.StoreInt64(&s.n, 0)
}

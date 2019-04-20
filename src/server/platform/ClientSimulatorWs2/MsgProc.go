package main

//
// Created by YangZhi
// 			4/9/2019
//

import (
	"sync"
	"sync/atomic"
)

//
const (
	Idle int32 = iota
	Running
)

//
const (
	Enter int = iota + 1
	Quit
)

//MsgProc 消息处理
type MsgProc struct {
	cell    MsgProcCell //处理单元
	l       *sync.Mutex
	c       *sync.Cond
	creator WorkerCreator
	sta     int32
}

func newMsgProc(creator WorkerCreator) *MsgProc {
	s := &MsgProc{creator: creator, l: &sync.Mutex{}, sta: Idle}
	s.c = sync.NewCond(s.l)
	return s
}

//
func (s *MsgProc) Sched(size int, cb func(int)) MsgProcCell {
	if atomic.CompareAndSwapInt32(&s.sta, Idle, Running) {
		go s.run(size, cb)
	}
	{
		s.l.Lock()
		for s.cell == nil {
			s.c.Wait()
		}
		s.l.Unlock()
	}
	return s.cell
}

//
func (s *MsgProc) run(size int, cb func(int)) {
	//time.Sleep(time.Second * 5)
	s.cell = newDefMsgProcCell(s.creator, size)
	//s.l.Lock()
	s.c.Signal()
	//s.l.Unlock()
	cb(Enter)
	s.cell.Run()
	atomic.StoreInt32(&s.sta, Idle)
	cb(Quit)
	s.cell = nil
}

package main

//
// Created by andy_ro@qq.com
// 			4/9/2019
//

import (
	"runtime"
	"sync/atomic"
)

//MsgProcPool 消息处理池
type MsgProcPool struct {
	cell  MsgProcCell
	cells []MsgProcCell
	x     int64
	n     int32
}

//
func NewMailbox() Mailbox {
	s := &MsgProcPool{}
	return s
}

//Start num cell数
//Start size 时间轮盘大小 size>=timeout>interval
func (s *MsgProcPool) Start(creator WorkerCreator, num, size int) MsgProcCell {
	proc := newMsgProc(creator)
	s.cell = proc.Sched(size, s.onCell)
	for i := 0; i < num; i++ {
		proc := newMsgProc(creator)
		cell := proc.Sched(size, s.onCell)
		s.cells = append(s.cells, cell)
	}
	return s.cell
}

//
func (s *MsgProcPool) onCell(b int) {
	if b == Enter {
		atomic.AddInt32(&s.n, 1)
	} else {
		atomic.AddInt32(&s.n, -1)
	}
}

//
func (s *MsgProcPool) GetNextCell() MsgProcCell {
	p := s.cell
	size := len(s.cells)
	if size > 0 {
		return s.cells[atomic.AddInt64(&s.x, 1)%int64(size)]
	}
	return p
}

//
func (s *MsgProcPool) Stop() {
	for _, cell := range s.cells {
		cell.GetMsq().Push(nil)
	}
	s.cell.GetMsq().Push(nil)
}

//
func (s *MsgProcPool) Wait() {
	for atomic.LoadInt32(&s.n) != 0 {
		runtime.Gosched()
	}
}

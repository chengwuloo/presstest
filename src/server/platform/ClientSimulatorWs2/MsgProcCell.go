package main

//
// Created by YangZhi
// 			4/9/2019
//

import (
	"log"
	"runtime"
	"runtime/debug"
	"server/platform/util"
	"sync"
	"sync/atomic"
	"time"
)

//MsgProcCell 消息处理单元
type MsgProcCell interface {
	AddTask(data *Event)
	AddReadTask(cmd uint32, msg interface{}, peer Session)
	AddReadTaskWith(handler ReadCallback, cmd uint32, msg interface{}, peer Session)
	AddCustomTask(cmd uint32, msg interface{}, peer Session)
	AddCustomTaskWith(handler CustomCallback, cmd uint32, msg interface{}, peer Session)
	//协程ID
	GetPID() uint32
	//消息队列
	GetMsq() MsgQueue
	//线程局部worker
	GetWorker() Worker
	//线程安全
	InCellThread() bool
	//空闲回调
	Exec(func())
	Append(func())
	//任务轮询
	Run()
	//添加会话
	Add(peer Session)
	//删除会话
	Remove(peer Session)
	//处理会话数
	Count() int64
}

//
type DefMsgProcCell struct {
	msq    MsgQueue       //任务队列
	msgLen int32          //任务数目
	pid    uint32         //协程ID
	worker Worker         //任务worker
	cbs    []func()       //空闲回调
	l      *sync.Mutex    //
	peers  map[int64]bool //注册会话
	c      int64          //注册会话数
	p      *sync.Mutex    //
}

//
func newDefMsgProcCell(creator WorkerCreator, size int) MsgProcCell {
	s := &DefMsgProcCell{
		pid: util.GoroutineID(),
		msq: newFreeVecMsq(),
		//msq:     newBlockVecMsq(),
		//msq:   newBlockChanMsq(),
		//msq:   newFreeChanMsq(),
		l:     &sync.Mutex{},
		peers: map[int64]bool{},
		p:     &sync.Mutex{},
	}
	s.worker = creator.CreateInstance(s)                //线程局部worker
	timeWheel := NewTimeWheel(s.pid, int32(size))       //指定时间轮大小
	timer := NewScopedTimer(s.pid, s.worker, timeWheel) //线程局部定时器
	timeWheel.SetTimer(timer)                           //绑定时间轮定时器
	s.worker.SetTimer(timer)                            //绑定worker定时器
	s.worker.SetCell(s)                                 //所在处理单元
	return s
}

//
func (s *DefMsgProcCell) Add(peer Session) {
	s.p.Lock()
	s.peers[peer.ID()] = true
	s.p.Unlock()
	atomic.AddInt64(&s.c, 1)
}

//
func (s *DefMsgProcCell) Remove(peer Session) {
	s.p.Lock()
	if _, ok := s.peers[peer.ID()]; ok {
		delete(s.peers, peer.ID())
		atomic.AddInt64(&s.c, -1)
	}
	s.p.Unlock()
}

//
func (s *DefMsgProcCell) Count() int64 {
	return atomic.LoadInt64(&s.c)
}

//
func (s *DefMsgProcCell) GetPID() uint32 {
	return s.pid
}

//GetMsq 消息队列
func (s *DefMsgProcCell) GetMsq() MsgQueue {
	return s.msq
}

//GetWorker 线程局部worker
func (s *DefMsgProcCell) GetWorker() Worker {
	return s.worker
}

//InCellThread 线程安全检查
func (s *DefMsgProcCell) InCellThread() bool {
	return util.GoroutineID() == s.pid
}

//Exec 执行空闲回调
func (s *DefMsgProcCell) Exec(cb func()) {
	if s.InCellThread() {
		cb()
	} else {
		s.Append(cb)
	}
}

//Append 添加空闲回调
func (s *DefMsgProcCell) Append(cb func()) {
	s.l.Lock()
	s.cbs = append(s.cbs, cb)
	s.l.Unlock()
	s.msq.Signal()
}

//Run 轮询
func (s *DefMsgProcCell) Run() {
	//捕获异常
	defer func() {
		if r := recover(); r != nil {
			log.Fatalln(debug.Stack())
		}
	}()
	worker := s.worker           //线程局部worker
	timer := s.worker.GetTimer() //线程局部定时器
	worker.OnInit()              //初始化worker
	i, t := 0, 200               //CPU分片调度
	flag := 0                    //单条或批量处理
	exit := false                //是否退出
	signal := make(chan bool, 1) //减轻cpu负担
	signal <- true
	s.msq.EnableNonBlocking(true)
EXIT:
	for {
		if i > t {
			i = 0
			runtime.Gosched()
		}
		i++
		//定时器轮询
		//log.Printf("--- *** ----------------------------- [%05d]Run Poll begin...\n", s.pid)
		bv := timer.Poll(s.pid, worker.OnTimer)
		//log.Printf("--- *** ----------------------------- [%05d]Run Poll end...\n", s.pid)
		if bv == true {
			//定时器已空
			//s.msq.EnableNonBlocking(false)
		} else {
			//定时器不空
			//s.msq.EnableNonBlocking(true)
		}
		//select {
		// case <-time.After(time.Microsecond):
		// 	time.Sleep(1e9)
		// 	signal <- true
		// 	//log.Printf("--- *** ----------------------------- [%05d]Run time.After...\n", s.pid)
		// case <-signal:
		// 	{
		// 		time.Sleep(1e9)
		// 		signal <- true
		// 		//log.Printf("--- *** ----------------------------- [%05d]Run signal...\n", s.pid)
		// 		break
		// 	}
		//default:
		{
			switch flag {
			case 0:
				{
					//单条消息处理
					msg, b := s.msq.Pop()
					exit = b
					if msg != nil && !exit {
						if _, ok := msg.(*Event); ok {
							//log.Printf("--- *** ----------------------------- [%05d]Run proc begin...\n", s.pid)
							s.proc(msg.(*Event), worker)
							//log.Printf("--- *** ----------------------------- [%05d]Run proc end...\n", s.pid)
						}
					}
					if nil == msg && !exit {
						//log.Printf("--- *** ----------------------------- [%05d]Run time.Sleep...\n", s.pid)
						time.Sleep(50 * time.Millisecond)
					}
					break
				}
			case 1:
				{
					//批量消息处理
					msgs, b := s.msq.Pick()
					exit = b
					for _, msg := range msgs {
						if _, ok := msg.(*Event); ok {
							//log.Printf("--- *** ----------------------------- [%05d]Run proc begin...\n", s.pid)
							s.proc(msg.(*Event), worker)
							//log.Printf("--- *** ----------------------------- [%05d]Run proc end...\n", s.pid)
						}
					}
					if 0 == len(msgs) && !exit {
						//log.Printf("--- *** ----------------------------- [%05d]Run time.Sleep...\n", s.pid)
						time.Sleep(50 * time.Millisecond)
					}
					break
				}
			}
			//处理空闲回调
			s.execFunc()
			if exit {
				timer.RemoveTimers()
				close(signal)
				break EXIT
			}
		}
		//}
	}
}

//execFunc 执行空闲回调
func (s *DefMsgProcCell) execFunc() {
	s.InCellThread()
	var cbs []func()
	{
		s.l.Lock()
		if len(s.cbs) > 0 {
			cbs = s.cbs[:]
			s.cbs = s.cbs[0:0]
		}
		s.l.Unlock()
	}
	for _, cb := range cbs {
		cb()
	}
}

//
func (s *DefMsgProcCell) AddTask(data *Event) {
	s.msq.Push(data)
}

//
func (s *DefMsgProcCell) AddReadTask(cmd uint32, msg interface{}, peer Session) {
	s.AddTask(createEvent(EVTRead, createReadEvent(cmd, msg, peer), nil))
}

//
func (s *DefMsgProcCell) AddReadTaskWith(handler ReadCallback, cmd uint32, msg interface{}, peer Session) {
	s.AddTask(createEvent(EVTRead, createReadEventWith(handler, cmd, msg, peer), nil))
}

//
func (s *DefMsgProcCell) AddCustomTask(cmd uint32, msg interface{}, peer Session) {
	s.AddTask(createEvent(EVTCustom, createCustomEvent(cmd, msg, peer), nil))
}

//
func (s *DefMsgProcCell) AddCustomTaskWith(handler CustomCallback, cmd uint32, msg interface{}, peer Session) {
	s.AddTask(createEvent(EVTCustom, createCustomEventWith(handler, cmd, msg, peer), nil))
}

//proc 处理任务队列
func (s *DefMsgProcCell) proc(data *Event, worker Worker) {
	worker.ResetDispatcher()
	switch data.ev {
	case EVTRead:
		ev := data.obj.(*readEvent)
		if ev.handler != nil {
			ev.handler(ev.cmd, ev.msg, ev.peer)
		} else {
			worker.OnRead(ev.cmd, ev.msg, ev.peer)
		}
	case EVTCustom:
		ev := data.obj.(*customEvent)
		if ev.handler != nil {
			ev.handler(ev.cmd, ev.msg, ev.peer)
		} else {
			worker.OnCustom(ev.cmd, ev.msg, ev.peer)
		}
	}
	if worker.GetDispatcher() != nil {
		worker.GetDispatcher().AddTask(data)
	}
}

package main

//
// Created by andy_ro@qq.com
// 			4/9/2019
//

// Worker 任务处理
type Worker interface {
	OnInit()
	OnConnected(peer Session, Type SesType)
	OnClosed(peer Session, Type SesType)
	RunAfter(delay int32, args interface{}) uint32
	RunAfterWith(delay int32, handler TimerCallback, args interface{}) uint32
	RunEvery(delay, interval int32, args interface{}) uint32
	RunEveryWith(delay, interval int32, handler TimerCallback, args interface{}) uint32
	RemoveTimer(timerID uint32)
	OnRead(cmd uint32, msg interface{}, peer Session)
	OnCustom(cmd uint32, msg interface{}, peer Session)
	OnTimer(timerID uint32, dt int32, args interface{}) bool
	SetTimer(t ScopedTimer)
	GetTimer() ScopedTimer
	GetTimeWheel() TimeWheel
	SetCell(c MsgProcCell)
	GetCell() MsgProcCell
	SetDispatcher(c MsgProcCell)
	GetDispatcher() MsgProcCell
	ResetDispatcher()
}

//
type DefWorker struct {
	t ScopedTimer //线程局部定时器
	c MsgProcCell //worker所在cell
	d MsgProcCell //分发到其它cell
}

//
func newWorker(c MsgProcCell) Worker {
	return &DefWorker{c: c}
}

//
func (s *DefWorker) SetTimer(t ScopedTimer) {
	s.t = t
}

//
func (s *DefWorker) GetTimer() ScopedTimer {
	return s.t
}

//
func (s *DefWorker) GetTimeWheel() TimeWheel {
	return s.t.GetTimeWheel()
}

//
func (s *DefWorker) SetCell(c MsgProcCell) {
	s.c = c
}

//
func (s *DefWorker) GetCell() MsgProcCell {
	return s.c
}

//
func (s *DefWorker) SetDispatcher(c MsgProcCell) {
	s.d = c
}

//
func (s *DefWorker) GetDispatcher() MsgProcCell {
	return s.d
}

//
func (s *DefWorker) ResetDispatcher() {
	s.d = nil
}

//
func (s *DefWorker) OnInit() {

}

//
func (s *DefWorker) RunAfter(delay int32, args interface{}) uint32 {
	return s.t.CreateTimer(delay, 0, args)
}

//
func (s *DefWorker) RunAfterWith(delay int32, handler TimerCallback, args interface{}) uint32 {
	return s.t.CreateTimerWithCB(delay, 0, handler, args)
}

//
func (s *DefWorker) RunEvery(delay, interval int32, args interface{}) uint32 {
	return s.t.CreateTimer(delay, interval, args)
}

//
func (s *DefWorker) RunEveryWith(delay, interval int32, handler TimerCallback, args interface{}) uint32 {
	return s.t.CreateTimerWithCB(delay, interval, handler, args)
}

//
func (s *DefWorker) RemoveTimer(timerID uint32) {
	s.t.RemoveTimer(timerID)
}

//
func (s *DefWorker) OnConnected(peer Session, stype SesType) {

}

//
func (s *DefWorker) OnClosed(peer Session, stype SesType) {

}

//
func (s *DefWorker) OnRead(cmd uint32, msg interface{}, session Session) {

}

//
func (s *DefWorker) OnCustom(cmd uint32, msg interface{}, session Session) {

}

//
func (s *DefWorker) OnTimer(timerID uint32, dt int32, args interface{}) bool {
	return false
}

//
type WorkerCreator interface {
	CreateInstance(c MsgProcCell) Worker
}

//
type workerCreator struct {
}

//
func NewWorkerCreator() WorkerCreator {
	return &workerCreator{}
}

//
func (s *workerCreator) CreateInstance(c MsgProcCell) Worker {
	return newWorker(c)
}

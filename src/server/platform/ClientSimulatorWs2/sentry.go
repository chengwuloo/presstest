package main

//
// Created by andy_ro@qq.com
// 			4/9/2019
//

//
type sentry struct {
	Worker
	access DefWorker
	main   *smain
	tickID uint32
}

//
func newsentry(c MsgProcCell) Worker {
	p := &sentry{}
	p.access.SetCell(c)
	p.main = newsmain(p)
	return p
}

//
func (s *sentry) OnInit() {
	s.main.initModuleHandlers()
}

//
func (s *sentry) RunAfter(delay int32, args interface{}) uint32 {
	return s.access.RunAfter(delay, args)
}

//
func (s *sentry) RunAfterWith(delay int32, handler TimerCallback, args interface{}) uint32 {
	return s.access.RunAfterWith(delay, handler, args)
}

//
func (s *sentry) RunEvery(delay, interval int32, args interface{}) uint32 {
	return s.access.RunEvery(delay, interval, args)
}

//
func (s *sentry) RunEveryWith(delay, interval int32, handler TimerCallback, args interface{}) uint32 {
	return s.access.RunEveryWith(delay, interval, handler, args)
}

//
func (s *sentry) RemoveTimer(timerID uint32) {
	s.access.RemoveTimer(timerID)
}

//
func (s *sentry) OnConnected(peer Session, stype SesType) {
	s.main.onConnected(peer, stype)
}

//
func (s *sentry) OnClosed(peer Session, stype SesType) {
	s.main.onClosed(peer, stype)
}

//
func (s *sentry) OnRead(cmd uint32, msg interface{}, peer Session) {
	s.main.onRead(cmd, msg, peer)
}

//
func (s *sentry) OnCustom(cmd uint32, msg interface{}, peer Session) {
	s.main.onCustom(cmd, msg, peer)
}

//
func (s *sentry) OnTimer(timerID uint32, dt int32, args interface{}) bool {
	return s.main.OnTimer(timerID, dt, args)
}

//
func (s *sentry) SetTimer(t ScopedTimer) {
	s.access.SetTimer(t)
}

//
func (s *sentry) GetTimer() ScopedTimer {
	return s.access.GetTimer()
}

//
func (s *sentry) GetTimeWheel() TimeWheel {
	return s.access.GetTimeWheel()
}

//
func (s *sentry) SetCell(c MsgProcCell) {
	s.access.SetCell(c)
}

//
func (s *sentry) GetCell() MsgProcCell {
	return s.access.GetCell()
}

//
func (s *sentry) SetDispatcher(c MsgProcCell) {
	s.access.SetDispatcher(c)
}

//
func (s *sentry) GetDispatcher() MsgProcCell {
	return s.access.GetDispatcher()
}

//
func (s *sentry) ResetDispatcher() {
	s.access.ResetDispatcher()
}

//
type SentryCreator struct {
	WorkerCreator
}

//
func NewSentryCreator() *SentryCreator {
	return &SentryCreator{}
}

//
func (s *SentryCreator) CreateInstance(c MsgProcCell) Worker {
	return newsentry(c)
}

package main

//
// Created by andy_ro@qq.com
// 			4/9/2019
//

//
type smain struct {
	entry    *sentry
	handlers CmdCallbacks
	player   *Player
}

//
func newsmain(s *sentry) *smain {
	return &smain{entry: s,
		player:   newPlayer(s),
		handlers: CmdCallbacks{}}
}

//
func (s *smain) onConnected(peer Session, stype SesType) {

}

//
func (s *smain) onClosed(peer Session, stype SesType) {

}

//
func (s *smain) onRead(cmd uint32, msg interface{}, peer Session) {
	if handler, ok := s.handlers[cmd]; ok {
		handler(msg, peer)
	} else {

	}
}

//
func (s *smain) onCustom(cmd uint32, msg interface{}, peer Session) {
	if handler, ok := s.handlers[cmd]; ok {
		handler(msg, peer)
	} else {

	}
}

//
func (s *smain) OnTimer(timerID uint32, dt int32, args interface{}) bool {
	return s.player.OnTimer(timerID, dt, args)
}

//
func (s *smain) initModuleHandlers() {
	s.entry.tickID = s.entry.RunEvery(500, 1000, nil)
	s.player.registerModuleHandler(s.handlers)
}

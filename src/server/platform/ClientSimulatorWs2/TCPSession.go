package main

//
// Created by andy_ro@qq.com
// 			4/9/2019
//

import (
	"log"
	"net"
	"sync"
	"sync/atomic"
)

//
type TCPSession struct {
	SesID         int64
	conn          net.Conn
	ctx           map[int]interface{}
	closing       int64
	Wg            sync.WaitGroup
	WMsq          MsgQueue    //写队列
	RMsq          MsgQueue    //读队列
	Channel       MsgTransmit //消息传输协议
	Cell          MsgProcCell //处理单元
	onConnected   OnConnected
	onClosed      OnClosed
	onMessage     OnMessage
	onWritten     OnWritten
	onError       OnError
	closeCallback CloseCallback
}

//
func newTCPSession(conn net.Conn) Session {
	sID := createSesID()
	peer := &TCPSession{
		SesID:   sID,
		conn:    conn,
		ctx:     map[int]interface{}{},
		WMsq:    newBlockVecMsq(),
		Channel: NewMyTCPTransmit()}
	return peer
}

//
func (s *TCPSession) ID() int64 {
	return s.SesID
}

//
func (s *TCPSession) GetCell() MsgProcCell {
	return s.Cell
}

//
func (s *TCPSession) Conn() interface{} {
	return s.conn
}

//事件回调
func (s *TCPSession) SetOnConnected(cb OnConnected) {
	s.onConnected = cb
}

//
func (s *TCPSession) SetOnClosed(cb OnClosed) {
	s.onClosed = cb
}

//
func (s *TCPSession) SetOnMessage(cb OnMessage) {
	s.onMessage = cb
}

//
func (s *TCPSession) SetOnError(cb OnError) {
	s.onError = cb
}

//
func (s *TCPSession) SetOnWritten(cb OnWritten) {
	s.onWritten = cb
}

//
func (s *TCPSession) SetCloseCallback(cb CloseCallback) {
	s.closeCallback = cb
}

//创建连接
func (s *TCPSession) OnEstablished() {
	s.Wg.Add(1)
	go s.readLoop()
	go s.writeLoop()
	if s.onConnected != nil {
		s.onConnected(s)
	}
}

//移除连接
func (s *TCPSession) OnDestroyed() {
	if s.onClosed != nil {
		s.onClosed(s)
	}
}

//
func (s *TCPSession) AddReadTask(cmd uint32, msg interface{}) {
	s.Cell.AddReadTask(cmd, msg, s)
}

//
func (s *TCPSession) AddReadTaskWith(handler ReadCallback, cmd uint32, msg interface{}) {
	s.Cell.AddReadTaskWith(handler, cmd, msg, s)
}

//
func (s *TCPSession) AddCustomTask(cmd uint32, msg interface{}, peer Session) {
	s.Cell.AddCustomTask(cmd, msg, s)
}

//
func (s *TCPSession) AddCustomTaskWith(handler CustomCallback, cmd uint32, msg interface{}) {
	s.Cell.AddCustomTaskWith(handler, cmd, msg, s)
}

//
func (s *TCPSession) SetCtx(key int, ctx interface{}) {
	if ctx != nil {
		//添加
		s.ctx[key] = ctx
	} else {
		//删除
		if _, ok := s.ctx[key]; ok {
			delete(s.ctx, key)
		}
	}
}

//
func (s *TCPSession) GetCtx(key int) interface{} {
	if p, ok := s.ctx[key]; ok {
		return p
	}
	return nil
}

//读协程
func (s *TCPSession) readLoop() {
	for {
		_, err := s.Channel.OnRecvMessage(s)
		if err != nil {
			if !IsEOFOrReadError(err) {
				if s.onError != nil {
					s.onError(s, err)
				}
			}
			break
		}
	}
	//对端关闭连接
	if 0 == atomic.LoadInt64(&s.closing) {
		//通知写退出
		s.WMsq.Push(nil)
	}
	//等待写退出
	s.Wg.Wait()
	s.conn = nil
	log.Println("exit readLoop...")
	if s.closeCallback != nil {
		s.closeCallback(s)
	}
}

//写协程
func (s *TCPSession) writeLoop() {
	for {
		msgs, exit := s.WMsq.Pick()
		for _, msg := range msgs {
			err := s.Channel.OnSendMessage(s, msg)
			if err != nil {
				if !IsEOFOrWriteError(err) {
					if s.onError != nil {
						s.onError(s, err)
					}
				}
				//break
			}
		}
		if exit {
			break
		}
	}
	//唤醒阻塞读
	if s.conn != nil {
		s.conn.Close()
	}
	log.Println("exit writeLoop...")
	s.Wg.Done()
}

//写
func (s *TCPSession) Write(msg interface{}) {
	s.WMsq.Push(msg)
}

//
func (s *TCPSession) Close() {
	//本端关闭连接
	if 0 == atomic.SwapInt64(&s.closing, 1) && s.conn != nil {
		//通知写退出
		s.WMsq.Push(nil)
	}
}

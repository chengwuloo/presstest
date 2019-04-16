package main

//
// Created by YangZhi
// 			4/9/2019
//

import (
	"log"
	"runtime/debug"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

//WSSession websocket会话
type WSSession struct {
	SesID         int64
	conn          *websocket.Conn
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
	reason        int32
	sta           int32
}

//
func newWSSession(conn *websocket.Conn) Session {
	sID := createSesID()
	peer := &WSSession{
		SesID:   sID,
		conn:    conn,
		ctx:     map[int]interface{}{},
		WMsq:    newBlockVecMsq(),
		Channel: NewMyWsTransmit()}
	return peer
}

//
func (s *WSSession) ID() int64 {
	return s.SesID
}

//
func (s *WSSession) GetCell() MsgProcCell {
	return s.Cell
}

//
func (s *WSSession) Conn() interface{} {
	return s.conn
}

//
func (s *WSSession) IsConnected() bool {
	return atomic.LoadInt32(&s.sta) == KConnected
}

//事件回调
func (s *WSSession) SetOnConnected(cb OnConnected) {
	s.onConnected = cb
}

//
func (s *WSSession) SetOnClosed(cb OnClosed) {
	s.onClosed = cb
}

//
func (s *WSSession) SetOnMessage(cb OnMessage) {
	s.onMessage = cb
}

//
func (s *WSSession) SetOnError(cb OnError) {
	s.onError = cb
}

//
func (s *WSSession) SetOnWritten(cb OnWritten) {
	s.onWritten = cb
}

//
func (s *WSSession) SetCloseCallback(cb CloseCallback) {
	s.closeCallback = cb
}

//创建连接
func (s *WSSession) OnEstablished() {
	s.Wg.Add(1)
	s.Cell = gMailbox.GetNextCell()
	//s.Cell.Add(s) //注册到Cell
	go s.readLoop()
	go s.writeLoop()
	if s.onConnected != nil {
		s.onConnected(s)
	}
}

//移除连接
func (s *WSSession) OnDestroyed() {
	if s.SesID == 0 {
		log.Panicln("SesID == 0")
	}
	//s.Cell.Remove(s) //从Cell中移除
	if s.onClosed != nil {
		s.onClosed(s)
	}
}

//
func (s *WSSession) AddReadTask(cmd uint32, msg interface{}) {
	s.Cell.AddReadTask(cmd, msg, s)
}

//
func (s *WSSession) AddReadTaskWith(handler ReadCallback, cmd uint32, msg interface{}) {
	s.Cell.AddReadTaskWith(handler, cmd, msg, s)
}

//
func (s *WSSession) AddCustomTask(cmd uint32, msg interface{}, peer Session) {
	s.Cell.AddCustomTask(cmd, msg, s)
}

//
func (s *WSSession) AddCustomTaskWith(handler CustomCallback, cmd uint32, msg interface{}) {
	s.Cell.AddCustomTaskWith(handler, cmd, msg, s)
}

//
func (s *WSSession) SetCtx(key int, ctx interface{}) {
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
func (s *WSSession) GetCtx(key int) interface{} {
	if p, ok := s.ctx[key]; ok {
		return p
	}
	return nil
}

//读协程
func (s *WSSession) readLoop() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalln(debug.Stack())
		}
	}()
	for {
		msg, err := s.Channel.OnRecvMessage(s)
		if err != nil {
			//log.Println("readLoop: ", err)
			// if !IsEOFOrReadError(err) {
			// 	if s.onError != nil {
			// 		s.onError(s, err)
			// 	}
			// }
			break
		}
		if msg == nil {
			log.Fatalln("readLoop: msg == nil")
		}
		if s.onMessage != nil {
			s.onMessage(msg, s)
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
	if s.closeCallback != nil {
		s.closeCallback(s)
	}
}

//写协程
func (s *WSSession) writeLoop() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalln(debug.Stack())
		}
	}()
	for {
		msgs, exit := s.WMsq.Pick()
		for _, msg := range msgs {
			err := s.Channel.OnSendMessage(s, msg)
			if err != nil {
				log.Println("writeLoop: ", err)
				if !IsEOFOrWriteError(err) {
					if s.onError != nil {
						s.onError(s, err)
					}
				}
				//break
			}
			if s.onWritten != nil {
				s.onWritten(msg, s)
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
	s.Wg.Done()
}

//写
func (s *WSSession) Write(msg interface{}) {
	s.WMsq.Push(msg)
}

//
func (s *WSSession) Close() {
	//本端关闭连接
	if 0 == atomic.SwapInt64(&s.closing, 1) && s.conn != nil {
		//通知写退出
		s.WMsq.Push(nil)
	}
}

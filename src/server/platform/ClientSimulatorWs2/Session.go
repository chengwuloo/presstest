package main

//
// Created by andy_ro@qq.com
// 			4/9/2019
//

import (
	"net"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

//关闭原因
const (
	ReasonPeerClosed int32 = iota + 1001 //对端关闭
	ReasonSelfClosed                     //本端关闭
	ReasonSelfExcept                     //本端异常
)

//状态
const (
	KConnected    int32 = iota + 1001 //已经连接
	KDisconnected                     //连接关闭
)

//
type Session interface {
	//会话ID
	ID() int64
	//连接
	Conn() interface{}
	//关闭
	Close()
	//写
	Write(msg interface{})
	//指定ctx
	SetCtx(key int, ctx interface{})
	//获取ctx
	GetCtx(key int) interface{}
	//获取处理单元
	GetCell() MsgProcCell
	//添加到读任务中处理
	AddReadTask(cmd uint32, msg interface{})
	AddReadTaskWith(handler ReadCallback, cmd uint32, msg interface{})
	//添加到自定义任务中处理
	AddCustomTask(cmd uint32, msg interface{}, peer Session)
	AddCustomTaskWith(handler CustomCallback, cmd uint32, msg interface{})
	//指定回调
	SetOnConnected(cb OnConnected)
	SetOnClosed(cb OnClosed)
	SetOnMessage(cb OnMessage)
	SetOnError(cb OnError)
	SetOnWritten(cb OnWritten)
	SetCloseCallback(cb CloseCallback)
	//连接建立
	OnEstablished()
	//连接移除
	OnDestroyed()
}

//
var gSesID int64

//
func createSesID() int64 {
	return atomic.AddInt64(&gSesID, 1)
}

//
func newSession(conn interface{}) Session {
	if c, ok := conn.(*websocket.Conn); ok {
		return newWSSession(c)
	}
	c, _ := conn.(net.Conn)
	return newTCPSession(c)
}

//
type SesType uint8

//
const (
	SesClient SesType = SesType(1)
	SesServer SesType = SesType(2)
)

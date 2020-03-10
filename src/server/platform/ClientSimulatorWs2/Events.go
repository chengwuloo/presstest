package main

//
// Created by andy_ro@qq.com
// 			4/9/2019
//

// 事件类型
const (
	EVTConnected int8 = iota + 10 // 连接
	EVTClosed                     // 关闭
	EVTRead                       // 读
	EVTSend                       // 写
	EVTCustom                     // 自定义
	EVTLogger                     // 日志
	EVTClose                      // 延迟关闭
)

// Event 事件数据
type Event struct {
	ev  int8
	obj interface{}
	ext interface{}
}

// createEvent 创建事件
func createEvent(ev int8, obj interface{}, ext interface{}) *Event {
	return &Event{ev, obj, ext}
}

//readEvent 读事件
type readEvent struct {
	cmd     uint32
	peer    Session
	msg     interface{}
	handler ReadCallback
}

//
func createReadEvent(cmd uint32, msg interface{}, peer Session) *readEvent {
	ev := &readEvent{cmd: cmd, msg: msg, peer: peer}
	//引用计数加1
	//ev.peer.AddRef("Event")
	return ev
}

//
func createReadEventWith(handler ReadCallback, cmd uint32, msg interface{}, peer Session) *readEvent {
	ev := &readEvent{handler: handler, cmd: cmd, msg: msg, peer: peer}
	//引用计数加1
	//ev.peer.AddRef("Event")
	return ev
}

//customEvent 自定义事件
type customEvent struct {
	cmd     uint32
	peer    Session
	msg     interface{}
	handler CustomCallback
}

//
func createCustomEvent(cmd uint32, msg interface{}, peer Session) *customEvent {
	ev := &customEvent{cmd: cmd, msg: msg, peer: peer}
	//引用计数加1
	//ev.peer.AddRef("Event")
	return ev
}

//
func createCustomEventWith(handler CustomCallback, cmd uint32, msg interface{}, peer Session) *customEvent {
	ev := &customEvent{handler: handler, cmd: cmd, msg: msg, peer: peer}
	//引用计数加1
	//ev.peer.AddRef("Event")
	return ev
}

package main

//
// Created by YangZhi
// 			4/9/2019
//

//
type CloseCallback func(peer Session)

//
type OnConnected func(peer Session)

//
type OnClosed func(peer Session)

//
type OnMessage func(msg interface{}, peer Session)

//
type OnWritten func(msg interface{}, peer Session)

//
type OnError func(peer Session, err error)

//
type ReadCallback func(cmd uint32, msg interface{}, peer Session)

//
type CustomCallback func(cmd uint32, msg interface{}, peer Session)

//
type CmdCallback func(msg interface{}, peer Session)

//
type CmdCallbacks map[uint32]CmdCallback

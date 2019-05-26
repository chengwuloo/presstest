package main

//
// Created by andy_ro@qq.com
// 			4/9/2019
//

//
type WSClient interface {
	//会话ID
	ID() int64
	//会话
	Session() Session
	//关闭
	Close()
	//写
	Write(msg interface{})
	//连接
	ConnectTCP(address string)
}

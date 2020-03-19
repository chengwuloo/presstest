package main

//
// Created by andy_ro@qq.com
// 			4/9/2019
//

//Mailbox 消息邮槽
type Mailbox interface {
	GetNextCell() MsgProcCell
	//Start creator worker工厂
	//Start num cell数
	//Start size 时间轮盘大小 size>=timeout>interval
	Start(creator WorkerCreator, num, size int) MsgProcCell
	Stop()
	Wait()
}

//
var gMailbox = NewMailbox()

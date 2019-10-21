# presstest

    服务器压测程序框架组件(golang版本)
		MsgProc.go 消息处理
		MsgProcCell.go 消息处理单元
		Mailbox.go(MsgProcPool.go) 邮槽(消息处理池)
		Worker.go 工作类(逻辑业务类继承入口)
		TimeWheel.go 时间轮用来处理超时连接
		ScopedTimer.go 线程局部定时器
		MsgTransmit.go 消息协议解析入口
		MsgQueue.go 消息队列
		Session.go 连接会话
		SessionMgr.go 会话管理类
		Timestamp.go 时间戳
		
    握手/交互协议：websocket&protobuf
    可以根据需要自定义修改底层通信协议
    LoadClientSimulatorWs 负责加载 ClientSimulatorWs2 子进程，windows/linux版
    配置在LoadClientSimulatorWs/conf.ini下

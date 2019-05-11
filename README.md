# presstest
棋牌游戏服压测客户端程序

    golang高并发压测客户端程序，协议：websocket&protobuf
    测试登陆大厅，大厅心跳，子游戏模块：龙虎斗/百人牛牛/二八杠/红黑大战
    可以根据需要自定义修改底层通信协议
    LoadClientSimulatorWs 负责加载 ClientSimulatorWs2 子进程，windows/linux版
    配置在LoadClientSimulatorWs/conf.ini下

# presstest

# golang写的h5压力测试客户端 websocket&protobuf
# 测试目前的棋牌游戏，有龙虎斗，百人牛牛，二八杠，红黑大战
# 可以根据需要自定义修改底层通信协议
# LoadClientSimulatorWs 负责加载 ClientSimulatorWs2 子进程，windows/linux版
# 配置在LoadClientSimulatorWs/conf.ini下
# 虽只是一个压测工具，但的确是一个能处理百万级并发的高性能框架，简单易用
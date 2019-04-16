echo on

LoadClientSimulatorWs.exe -children=1 -httpaddr=192.168.2.30:801 -wsaddr=192.168.2.75:10000 -mailboxs=10 -clients=10 -baseTest=123 -deltaClients=5 -deltaTime=3000 -interval=8000 -timeout=10000

@pause

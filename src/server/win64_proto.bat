
%GOBIN%/protoc64 --plugin=protoc-gen-gogofaster=%GOBIN%/protoc-gen-gogofaster.exe --gogofaster_out=%GOPATH%/src/server/pb/Game_Common --proto_path=Y:\Landy\TianXia\Program\proto\src Game.Common.proto
%GOBIN%/protoc64 --plugin=protoc-gen-gogofaster=%GOBIN%/protoc-gen-gogofaster.exe --gogofaster_out=%GOPATH%/src/server/pb/HallServer --proto_path=Y:\Landy\TianXia\Program\proto\src HallServer.Message.proto
%GOBIN%/protoc64 --plugin=protoc-gen-gogofaster=%GOBIN%/protoc-gen-gogofaster.exe --gogofaster_out=%GOPATH%/src/server/pb/GameServer --proto_path=Y:\Landy\TianXia\Program\proto\src GameServer.Message.proto

::红黑大战
%GOBIN%/protoc64 --plugin=protoc-gen-gogofaster=%GOBIN%/protoc-gen-gogofaster.exe --gogofaster_out=%GOPATH%/src/server/pb/HongHei --proto_path=Y:\Landy\TianXia\Program\proto\src HongHei.Message.proto
::百人牛牛
%GOBIN%/protoc64 --plugin=protoc-gen-gogofaster=%GOBIN%/protoc-gen-gogofaster.exe --gogofaster_out=%GOPATH%/src/server/pb/Brnn --proto_path=Y:\Landy\TianXia\Program\proto\src Brnn.Message.proto
::二八杠
%GOBIN%/protoc64 --plugin=protoc-gen-gogofaster=%GOBIN%/protoc-gen-gogofaster.exe --gogofaster_out=%GOPATH%/src/server/pb/ErBaGang --proto_path=Y:\Landy\TianXia\Program\proto\src ErBaGang.Message.proto
::龙虎斗
%GOBIN%/protoc64 --plugin=protoc-gen-gogofaster=%GOBIN%/protoc-gen-gogofaster.exe --gogofaster_out=%GOPATH%/src/server/pb/Longhu --proto_path=Y:\Landy\TianXia\Program\proto\src Longhu.Message.proto

@pause
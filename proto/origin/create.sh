#!/bin/sh

DST_DIR="$(go env GOPATH)/src/github.com/teamtalk-remix/ttr-client-incubator/test-client-go/proto"

#for golang, it's a little manual work, see https://github.com/golang/protobuf/issues/39
protoc -I=. --go_out=$DST_DIR/IM_BaseDefine $SRC_DIR/IM.BaseDefine.proto
protoc -I=. --go_out=$DST_DIR/IM_File $SRC_DIR/IM.File.proto
protoc -I=. --go_out=$DST_DIR/IM_Login $SRC_DIR/IM.Login.proto
protoc -I=. --go_out=$DST_DIR/IM_Other $SRC_DIR/IM.Other.proto
protoc -I=. --go_out=$DST_DIR/IM_SwitchService $SRC_DIR/IM.SwitchService.proto
protoc -I=. --go_out=$DST_DIR/IM_Buddy $SRC_DIR/IM.Buddy.proto
protoc -I=. --go_out=$DST_DIR/IM_Group $SRC_DIR/IM.Group.proto
protoc -I=. --go_out=$DST_DIR/IM_Message $SRC_DIR/IM.Message.proto
protoc -I=. --go_out=$DST_DIR/IM_Server $SRC_DIR/IM.Server.proto
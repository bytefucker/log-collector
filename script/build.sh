#!/bin/bash

ROOTPATH=$(pwd)

echo "构建根目录:$ROOTPATH"

echo "----------开始构建logAgent-------"
cd $ROOTPATH/logAgent
cp -R config bin/config
go build -x -o bin/logAgent logAgent
echo "----------构建logAgent成功-------"

echo "----------开始构建logManager-------"
cd $ROOTPATH/logManager
cp -R conf bin/conf
go build -x -o bin/logManager logManager/manager
echo "----------构建logManager成功-------"

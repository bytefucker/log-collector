#!/usr/bin/env bash
set -e

pwd
echo 开始编译镜像
docker build -f docker/Dockerfile -t ampregistry:5000/log-collector .
echo 开始推送镜像
docker push  ampregistry:5000/log-collector:latest
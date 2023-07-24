#!/usr/bin/bash

WHO=$(whoami | grep "root")
if [ -z ${WHO} ]; then
    # 切换到root用户
    sudo su
    source ~/.bash_profile
fi

WHO=$(whoami | grep "root")
if [ -z ${WHO} ]; then
    # 切换失败，退出
    echo "please change to root user"
    exit 1
fi

# 创建环境文件夹
mkdir -p /mnt/env
cd /mnt/env

yum install -y wget

if command -v go >/dev/null 2>&1 ; then
  echo "Golang is installed"
else
  echo "install golang..."
  wget https://go.dev/dl/go1.19.11.linux-amd64.tar.gz
  tar -zxvf go1.19.11.linux-amd64.tar.gz
  cp -r go1.19.11.linux-amd64/bin /usr/bin
fi

if command -v go >/dev/null 2>&1 ; then
  echo "node is installed"
else
  echo "install node..."
  wget https://registry.npmmirror.com/-/binary/node/v16.20.1/node-v16.20.1-linux-x64.tar.gz
  tar -zxvf node-v16.20.1-linux-x64.tar.gz
  cp -r node-v16.20.1-linux-x64.tar.gz/bin /usr/bin
fi

certbot certonly -d hk.cgtnew.com -n -m hkcgt@gmail.com --preferred-challenges http --standalone

npm install -g yarn

./build.sh

if command -v certbot >/dev/null 2>&1 ; then
  echo "certbot is installed"
else
  yum install epel-release -y
  yum install certbot -y
fi

read -p "Enter domain for apply SSL: " domain

certbot certonly -d ${domain} -m hkcgt@gmail.com --preferred-challenges http --standalone

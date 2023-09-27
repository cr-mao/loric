#!/bin/bash

# 无状态 发布脚本，仅供参考

if [ $# -lt 4 ];then
  echo "Usage: develop.sh [type_name] [branch_name] [port] [env]"
  exit 1
fi

if [[ "$1" != "gate" && "$1" != "auth" && "$1" != "limeng_chat" && "$1" != "limeng_chat1" ]];then
  echo "type_name error ,must one of (gate,auth,status_node,normal_node)"
  exit 1
fi


if [[ "$4" != "testing" && "$4" != "local" && "$4" != "production" ]];then
  echo "env error ,must one of (testing,local,production)"
  exit 1
fi

responsity_prefix="/opt/loricgameresponsity"
responsity_path="/opt/loricgameresponsity/loric"
responsity_ssh_url="git@github.com:cr-mao/loric.git"
main_path_prefix="/opt/loricgameresponsity/loric/example/"

supversior_path="/etc/supervisor/conf.d"

prefix="loric"

# 拼接 type_name+port  /opt/loricgame/gate9001, /opt/loricgame/auth10001
project_prefix_path="/opt/loricgame/"

git config --global --add safe.directory $responsity_path
if [ ! -d $responsity_path ];then
  if [ ! -d $responsity_prefix ];then
     mkdir -p $responsity_prefix
  fi
  cd $responsity_prefix
  git clone $responsity_ssh_url
  if [ $? != 0 ];then
    echo "代码啦不下来"
    exit 1
  fi
fi
if [ $# == 1 ];then
   branch="main"
else
   branch=$2
fi


cd $responsity_path
git fetch origin ${branch}
git fetch -p
git reset --hard origin/${branch}

# 编译go 程序

cd $main_path_prefix$1
if [ -e "main" ];then
  rm -f main
fi

GOROOT=/usr/local/go
GOBIN=$GOROOT/bin
# linux 下 gopath不管用。 都在用户下的go  /home/ubuntu/go
GOPATH=/opt/go1.18.6
GO111MODULE=on
GOPROXY=https://goproxy.cn,direct
PATH=$PATH:$GOBIN:$GOPATH/bin
VERSION=`git rev-parse --short HEAD`
/usr/local/go/bin/go mod tidy
BUILDTIME=`date +%FT%T`
LDFLAGS="-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILDTIME}"
CGO_ENABLED=0 GOOS=linux /usr/local/go/bin/go build -o  $1 -ldflags="-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILDTIME}"  main.go && go clean -cache

#if [ -d $project_prefix_path$1$3 ];then
#    echo "此端口项目已经存在"
#    exit 1
#fi

if [ ! -e $main_path_prefix$1/$4.config.yaml ];then
  echo "$main_path_prefix$1/$4.config.yaml 配置文件不存在"
  exit 1
fi

# 创建目录，cp 配置文件
if [ -d $project_prefix_path$1$3 ];then
  mkdir -p $project_prefix_path$1$3
fi





supver_conf_has=`supervisorctl status | grep $prefix_$1_$3 | wc -l`

if [ $supver_conf_has -eq 0 ];then

  cp $main_path_prefix$1/$1 $project_prefix_path$1$3/
  cp $main_path_prefix$1/$4.config.yaml $project_prefix_path$1$3/

  chmod +x $project_prefix_path$1$3/$1
  chown root:root  $project_prefix_path$1$3/$1

 # supervisor 一下。
 touch $supversior_path/${prefix}_$1_$3.conf
 cat > $supversior_path/${prefix}_$1_$3.conf <<EOF
[program:${prefix}_$1_$3]
directory=$project_prefix_path$1$3
command=$project_prefix_path$1$3/$1  --env=local
autostart=true
autorestart=true
startretries=10
redirect_stderr=true
stdout_logfile=$project_prefix_path$1$3/serve.log
EOF
 # 更新配置 启动
 supervisorctl update
else
  supervisorctl stop ${prefix}_$1_$3
  cp $main_path_prefix$1/$1 $project_prefix_path$1$3/
  cp $main_path_prefix$1/$4.config.yaml $project_prefix_path$1$3/
  chmod +x $project_prefix_path$1$3/$1
  chown root:root  $project_prefix_path$1$3/$1
  supervisorctl start ${prefix}_$1_$3
fi






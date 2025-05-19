#!/bin/bash

info() {
    echo -e "\033[32m[OpenResty Manager] $*\033[0m"
}

warning() {
    echo -e "\033[33m[OpenResty Manager] $*\033[0m"
}

abort() {
    echo -e "\033[31m[OpenResty Manager] $*\033[0m"
    exit 1
}

if [ -z "$BASH" ]; then
	abort "请使用 bash 执行此脚本，并参考最新的官方技术文档 https://www.uusec.com/"
fi

if [ "$EUID" -ne "0" ]; then
	abort "请以 root 权限运行"
fi

SCRIPT_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd  "$SCRIPT_PATH"

if [ ! $(command -v docker) ]; then
	warning "未检测到 Docker 引擎，我们将自动为您安装。过程很慢，请耐心等待 ..."
	sh install-docker.sh --mirror Aliyun
	if [ $? -ne "0" ]; then
		abort "Docker 引擎自动安装失败，请在执行此脚本之前手动安装它。"
	fi
	systemctl enable docker && systemctl daemon-reload && systemctl restart docker
fi

DC_CMD="docker compose"
$DC_CMD version > /dev/null 2>&1
if [ $? -ne "0" ]; then
	abort "你的 Docker 版本太低，缺少 'docker compose' 命令。请卸载并安装最新版本"
fi

if [ ! -f ".env" ];then
	echo "" > .env
	if [ $(command -v firewall-cmd) ]; then
		firewall-cmd --permanent --add-port={80,443,34567}/tcp > /dev/null 2>&1
		firewall-cmd --reload > /dev/null 2>&1
	elif [ $(command -v ufw) ]; then
		for port in 80 443 34567; do ufw allow $port/tcp > /dev/null 2>&1; done
		ufw reload > /dev/null 2>&1
	fi
fi

stop_om(){
	$DC_CMD down
}

uninstall_om(){
	stop_om
	docker rm -f openresty-manager > /dev/null 2>&1
	docker images|grep openresty-manager|awk '{print $3}'|xargs docker rmi -f > /dev/null 2>&1
	docker volume ls|grep _om_|awk '{print $2}'|xargs docker volume rm -f > /dev/null 2>&1
}

start_om(){
	if [ ! $(command -v netstat) ]; then
		$( command -v yum || command -v apt-get || command -v zypper ) -y install net-tools
	fi
	port_status=`netstat -nlt|grep -E ':(80|443|34567)\s'|wc -l`
	if [ $port_status -gt 0 ]; then
		abort "端口80、443、34567中的一个或多个被占用。请关闭相应的服务或修改其端口"
	fi
	$DC_CMD up -d --remove-orphans
}

upgrade_om(){
	$DC_CMD pull
	$DC_CMD up -d --remove-orphans
}

repair_om(){
	if [ $(command -v firewall-cmd) ]; then
		systemctl restart firewalld > /dev/null 2>&1
	elif [ $(command -v ufw) ]; then
		systemctl restart ufw > /dev/null 2>&1
	fi
	systemctl daemon-reload
	systemctl restart docker
}

restart_om(){
	stop_om
	start_om
}

start_menu(){
    clear
    echo "========================="
    echo "OpenResty Manager 管理"
    echo "========================="
    echo "1. 启动"
    echo "2. 停止"
    echo "3. 重启"
    echo "4. 更新"
    echo "5. 修复"
    echo "6. 卸载"
    echo "7. 退出"
    echo
    read -p "请输入数字序号: " num
    case "$num" in
	1)
	start_om
	info "启动已完成"
	;;
	2)
	stop_om
	info "停止已完成"
	;;
	3)
	restart_om
	info "重启已完成"
	;;
	4)
	upgrade_om
	info "升级已完成"
	;;
	5)
	repair_om
	info "修复已完成"
	;;
	6)
	uninstall_om
	info "卸载已完成"
	;;
	7)
	exit 1
	;;
	*)
	clear
	info "请输入正确的数字"
	;;
    esac
    sleep 3s
    start_menu
}

start_menu

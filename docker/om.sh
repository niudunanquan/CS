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
	abort "Please execute this script using bash and refer to the latest official technical documentation https://www.uusec.com/"
fi

if [ "$EUID" -ne "0" ]; then
	abort "Please run with root privileges"
fi

SCRIPT_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd  "$SCRIPT_PATH"

if [ ! $(command -v docker) ]; then
	warning "Docker Engine not detected, we will automatically install it for you. The process is slow, please be patient ..."
	sh install-docker.sh
	if [ $? -ne "0" ]; then
		abort "Automatic installation of Docker Engine failed. Please manually install it before executing this script"
	fi
	systemctl enable docker && systemctl daemon-reload && systemctl restart docker
fi

DC_CMD="docker compose"
$DC_CMD version > /dev/null 2>&1
if [ $? -ne "0" ]; then
	abort "Your Docker version is too low and lacks the 'docker compose' command. Please uninstall and install the latest version"
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
		abort "One or more of ports 80, 443, 34567 are occupied. Please shutdown the corresponding service or modify its port"
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
    echo "OpenResty Manager Management"
    echo "========================="
    echo "1. Start"
    echo "2. Stop"
    echo "3. Restart"
    echo "4. Upgrade"
    echo "5. Repair"
    echo "6. Uninstall"
    echo "7. Exit"
    echo
    read -p "Please enter the number: " num
    case "$num" in
	1)
	start_om
	info "Startup completed"
	;;
	2)
	stop_om
	info "Stop completed"
	;;
	3)
	restart_om
	info "Restart completed"
	;;
	4)
	upgrade_om
	info "Upgrade completed"
	;;
	5)
	repair_om
	info "Repair completed"
	;;
	6)
	uninstall_om
	info "Uninstall completed"
	;;
	7)
	exit 1
	;;
	*)
	clear
	info "Please enter the right number"
	;;
    esac
    sleep 3s
    start_menu
}

start_menu

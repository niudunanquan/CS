#!/bin/bash

# OpenResty Manager one click installation script
# Supported system: CentOS/RHEL 7+, Debian 11+, Ubuntu 18+, Fedora 32+, etc

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

if [[ $EUID -ne 0 ]]; then
    abort "此脚本必须以root权限运行"
fi

OS_ARCH=$(uname -m)
case "$OS_ARCH" in
    x86_64)
    ;;
    *)
    abort "不支持的 CPU 架构: $OS_ARCH"
    ;;
esac

if [ -f /etc/os-release ]; then
    source /etc/os-release
    OS_NAME=$ID
    OS_VERSION=$VERSION_ID
elif type lsb_release >/dev/null 2>&1; then
    OS_NAME=$(lsb_release -si | tr '[:upper:]' '[:lower:]')
    OS_VERSION=$(lsb_release -sr)
else
    abort "无法检测操作系统"
fi

check_ports() {
    if [ $(command -v ss) ]; then
        for port in 80 443 777 34567; do
            if ss -tln "( sport = :${port} )" | grep -q LISTEN; then
                abort "端口 ${port} 被占用, 请关闭该端口后重新安装"
            fi
        done
	fi
}

install_openresty_manager() {
    curl https://om.uusec.com/docker_cn.tgz -o /tmp/docker.tgz
    mkdir -p /opt && tar -zxf /tmp/docker.tgz -C /opt/
    if [ $? -ne "0" ]; then
        abort "OpenResty Manager安装失败"
    fi
}

allow_firewall_ports() {
    if [ ! -f "/opt/om/.fw" ];then
        echo "" > /opt/om/.fw
        if [ $(command -v firewall-cmd) ]; then
            firewall-cmd --permanent --add-port={80,443,34567}/tcp > /dev/null 2>&1
            firewall-cmd --reload > /dev/null 2>&1
        elif [ $(command -v ufw) ]; then
            for port in 80 443 34567; do ufw allow $port/tcp > /dev/null 2>&1; done
            ufw reload > /dev/null 2>&1
        fi
    fi
}

main() {
    info "检测到系统：${OS_NAME} ${OS_VERSION} ${OS_ARCH}"

    warning "检查端口冲突 ..."
    check_ports

    if [ ! -e "/opt/om" ]; then
        warning "安装OpenResty Manager..."
        install_openresty_manager
    else
        abort '目录 "/opt/om" 已存在, 请确认删除后再试'
    fi

    warning "添加防火墙端口例外..."
    allow_firewall_ports

    bash /opt/om/om.sh
}

main

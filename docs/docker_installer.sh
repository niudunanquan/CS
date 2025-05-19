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
    abort "This script must be run with root privileges"
fi

OS_ARCH=$(uname -m)
case "$OS_ARCH" in
    x86_64|arm*|aarch64)
    ;;
    *)
    abort "Unsupported CPU arch: $OS_ARCH"
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
    abort "Unable to detect operating system"
fi

check_ports() {
    if [ $(command -v ss) ]; then
        for port in 80 443 777 34567; do
            if ss -tln "( sport = :${port} )" | grep -q LISTEN; then
                abort "Port ${port} is occupied, please close it and try again"
            fi
        done
	fi
}

install_openresty_manager() {
    curl https://om.uusec.com/docker.tgz -o /tmp/docker.tgz
    mkdir -p /opt && tar -zxf /tmp/docker.tgz -C /opt/
    if [ $? -ne "0" ]; then
        abort "Installation of OpenResty Manager failed"
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
    info "Detected system: ${OS_NAME} ${OS_VERSION} ${OS_ARCH}"

    warning "Check for port conflicts ..."
    check_ports

    if [ ! -e "/opt/om" ]; then
        warning "Install OpenResty Manager ..."
        install_openresty_manager
    else
        abort 'The directory "/opt/om" already exists, please confirm to remove it and try again'
    fi

    warning "Add firewall ports exception ..."
    allow_firewall_ports

    bash /opt/om/om.sh
}

main

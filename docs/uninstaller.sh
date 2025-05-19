#!/bin/bash

# OpenResty Manager one click uninstallation script
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

main() {
    info "Detected system: ${OS_NAME} ${OS_VERSION} ${OS_ARCH}"

    warning "Uninstall OpenResty Manager ..."
    if [ -f "/opt/om/oms" ]; then
        /opt/om/oms -s stop > /dev/null 2>&1
        /opt/om/oms -s uninstall > /dev/null 2>&1
        rm -rf /opt/om
    elif [ -f "/opt/om/om.sh" ]; then
        SCRIPT_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
        cd  "$SCRIPT_PATH"
        docker compose down > /dev/null 2>&1
        docker rm -f openresty-manager > /dev/null 2>&1
	    docker images|grep openresty-manager|awk '{print $3}'|xargs docker rmi -f > /dev/null 2>&1
	    docker volume ls|grep _om_|awk '{print $2}'|xargs docker volume rm -f > /dev/null 2>&1
        rm -rf /opt/om
    else
        abort 'Not found OpenResty Manager in directory "/opt/om"'
    fi

    info "Congratulations on the successful uninstallation"
}

main

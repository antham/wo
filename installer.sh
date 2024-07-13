#!/bin/sh

set -u

DEFAULT_INSTALL_PATH="/usr/local/bin"

main() {
    _set_path=false
    while [ $# -gt 0 ]; do
        case "$1" in
        "-o" | "--output")
            shift
            _path=$1
            _set_path=true
            ;;
        "--platform")
            shift
            _user_platform=$1
            ;;
        "--help")
            printHelp
            ;;
        esac
        shift
    done

    if [[ ${_set_path} != true ]]; then
        _path=$DEFAULT_INSTALL_PATH
    fi
    get_architecture || return 1
    local _platform="${_user_platform:-$PLATFORM}"

    echo "Install to ${_path}"
    URL=$(curl -s https://api.github.com/repos/antham/wo/releases/latest|grep browser_download_url|grep -o "https://.*${_platform}.tar.gz")
    curl -s --output - -L "$URL"|tar -zx -C "${_path}" wo || err "A failure occurred when installing the app, ensure you have the access to the folder where you want to install the app"
    echo "Application installed"
}

get_architecture() {
    local _ostype _cputype
    _ostype="$(uname -s)"
    _cputype="$(uname -m)"

    if [ "$_ostype" = Darwin ] && [ "$_cputype" = i386 ]; then
        if sysctl hw.optional.x86_64 | grep -q ': 1'; then
            _cputype=x86_64
        fi
    fi

    case "$_cputype" in
    xscale | arm | armv6l | armv7l | armv8l | aarch64 | arm64)
        _cputype=arm64
        ;;
    x86_64 | x86-64 | x64 | amd64)
        _cputype=amd64
        ;;
    *)
        err "unknown CPU type: $_cputype"
        ;;
    esac

    case "$_ostype" in
    Linux | FreeBSD | NetBSD | DragonFly)
        _ostype=linux
        ;;
    Darwin)
        _ostype=darwin
        ;;
    *)
        err "unrecognized OS type: $_ostype"
        ;;
    esac

    PLATFORM="${_ostype}_${_cputype}"
}

err() {
    echo "$1" >&2
    exit 1
}

printHelp() {
    echo "Wo CLI shell installer

Usage:
    curl -sSf https://github.com/antham/wo/blob/main/installer.sh | sh
    curl -sSf https://github.com/antham/wo/blob/main/installer.sh | sh -s -- [options]

Options:
    -o, --output  Path to save the binary to
    --platform    Specify platform to assume, e.g. 'darwin_arm64', 'linux_amd64', ...
    --help        Print this help message
"
    exit 0
}

main "$@" || exit 1

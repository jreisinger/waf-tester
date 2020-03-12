#!/bin/bash
set -e

######## FUNCTION ########

# Helper function for printing formatted colorfule messages.
print_msg() {
    msg=$1
    # https://stackoverflow.com/questions/5947742/how-to-change-the-output-color-of-echo-in-linux
    COLOR='\033[0;34m'
    NC='\033[0m' # No Color
    printf "${COLOR}--> ${msg}${NC}\n"
}

# Setup temp dir, copy needed files into it and create dirs for logs.
prepare_dirs() {
    local dir=$1

    print_msg "Create a temporary directory"
    TMPDIR=$(mktemp -d)

    # Bail out if the temp directory wasn't created successfully.
    if [ ! -e $TMPDIR ]; then
        >&2 echo "Failed to create temp directory"
        exit 1
    fi

    print_msg "Copy all needed files to $TMPDIR"
    cp docker-compose.yaml $TMPDIR
    cp -r $dir/* $TMPDIR

    OLDDIR=$(pwd)
    cd $TMPDIR

    print_msg "Create directories for logs in /tmp"
    mkdir -p /tmp/var/log
    mkdir -p /tmp/var/log/nginx
}

# Remove docker containers and networks. Remove temp dir.
clean_up () {
    print_msg "Cleaning up"
    docker-compose down
    #rm -rf /tmp/var/log
    rm -rf "$TMPDIR"
}

# Build and run WAF and protected demo app.
build_waf() {
    local verbose=$1
    print_msg "Run docker images for WAF and protected web server"
    if [[ $verbose == 1 ]]; then
        docker-compose up -d --build
    else
        docker-compose up -d --build > /dev/null
    fi
}

# Test WAF has been started.
test_waf_is_up() {
    print_msg "Check WAF is proxying the requests (must be 200)"
    sleep 3
    curl --fail -sI localhost:80 | grep "^HTTP"
}

help_and_exit() {
    cat << EOF
${0##*/} [options] DIR

Build a WAF container and run it locally together with a test web server
container. DIR is directory containing the Dockerfile and possibly other
related files.

  options:
    -h, -?  help
    -v      be verbose
    -s      sleep (don't clean up containers) until Ctrl-C
EOF

    exit $1
}

######## MAIN ########

# A POSIX variable
OPTIND=1    # Reset in case getopts has been used previously in the shell.

# Initialize our own variables:
verbose=0
mysleep=0

# getops can't handle long options (--help) but can handle options clustering
# (-vf FILE)
while getopts "h?vs" opt; do
    case "$opt" in
    h|\?)
        help_and_exit 0
        ;;
    v)  verbose=1
        ;;
    s)  mysleep=1
        ;;
    esac
done

# Shift off the options and optional --
shift $((OPTIND-1))
[ "$1" = "--" ] && shift

dir="$@"
if [[ ! -d $dir ]]; then
    echo "'$dir' is not a valid directory"
    echo
    help_and_exit
fi

# Make sure we run cleanup actions even if the script exits abnormally.
trap "exit 1"   HUP INT PIPE QUIT TERM
trap clean_up   EXIT

prepare_dirs "$dir"
build_waf "$verbose"
test_waf_is_up
test_rules_basic
#test_rules_ftw

if [[ $mysleep == 1 ]]; then
    print_msg "Sleeping until Ctrl-C"
    sleep 999999
fi

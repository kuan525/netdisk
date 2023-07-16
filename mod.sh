#!/bin/bash

function go_mod_tidy_recursive() {
    local dir="$1"

    if [[ -f "$dir/go.mod" && $dir != $rootDir ]]; then
        echo "正在更新: $dir"
        cd "$dir"
        (go get github.com/kuan525/netdisk@$tag && go mod tidy)
    fi

    for subdir in "$dir"/*; do
        if [[ -d "$subdir" ]]; then
            go_mod_tidy_recursive "$subdir"
        fi
    done
}

if [[ -z "$1" ]]; then
    echo "请输入全局想要拉取netdisk的版本"
    exit 1
fi

tag=$1
rootDir=$(pwd)
go_mod_tidy_recursive $rootDir

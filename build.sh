#!/bin/bash

get_tag() {
    local tag=$(git log -1 --pretty="%ad-%h" --date=short $1)
    echo "$tag"
}

build_sub() {
    local tag=$(get_tag $1)
    local img="$LINKSAFE_REGISTRY/$2:$tag"
    local org_dir=$(pwd)
    local dockerfile="$3"
    docker manifest inspect $img > /dev/null
    if [ "$?" -ne "0" ]; then
        cd "$1"
        if [ ! -z "$dockerfile" ]; then
            docker build -f $dockerfile -t "$img" .
        else
            docker build -t "$img" .
        fi
        docker push "$img"
        cd "$org_dir"
        return 0
    fi
    echo "Skipping $img"
}

show_sub() {
    local tag=$(get_tag $1)
    local img="$LINKSAFE_REGISTRY/$2:$tag"
    echo "$img"
}

do_build() {
    build_sub . lancsnet/linksafe/partner-api
}

listing() {
    show_sub . lancsnet/linksafe/partner-api
}

if [ -z "$LINKSAFE_ENV" ]; then
    echo "LINKSAFE_ENV not setup"
    exit 1   
fi

case "$1" in 
    build)
       do_build
       ;;
    list)
       listing
       ;;
    *)
       echo "Usage: $0 {build|list}"
esac
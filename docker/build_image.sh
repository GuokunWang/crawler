#/bin/bash

function check {
    if [ ! $? =  0 ]; then
        echo 'check '$1' failed'
        exit
    fi
}

##### set env #####
GO_VERSION=$(go version | awk '{print $3}')
MAJOR_VERSION=$(echo $GO_VERSION | cut -d '.' -f1 | sed 's/go//')
MINOR_VERSION=$(echo $GO_VERSION | cut -d '.' -f2)

if [ $MAJOR_VERSION -lt 1 ] || [ $MAJOR_VERSION -eq 1 -a $MINOR_VERSION -lt 18 ]; then
    echo "Go version is 1.17 or lower, not supported, please udpate golang"
    exit
fi

export GOPROXY=https://goproxy.cn

##### current branch #####

##### build crawler #####
CURR_DIR=$(dirname $(realpath -s "$0"))
PROJ_DIR=$(dirname $CURR_DIR)
BRANCH=$(git rev-parse --abbrev-ref HEAD)

echo "项目目录： $PROJ_DIR"
echo "项目分支： $BRANCH"
cd $PROJ_DIR

echo 'build crawler'
go mod tidy && go build
check 'build crawler'

##### make container images #####
cd $CURR_DIR

##### make crawler container image #####
rm -fr ./crawler
cp $PROJ_DIR/crawler ./crawler
cp $PROJ_DIR/config.json ./config.json
check 'copy crawler'

docker build -t crawler:$BRANCH -f Dockerfile .
check 'build crawler container'
rm -fr ./crawler

#!/usr/bin/env bash

docker run -e GO111MODULE=on -v $(pwd):/go/src/github.com/nanananakam/twitterbot-post-tweet golang:1.12-stretch go build -o /go/src/github.com/nanananakam/twitterbot-post-tweet/main /go/src/github.com/nanananakam/twitterbot-post-tweet/main.go
zip Main.zip Main
rm go.sum
rm -f main
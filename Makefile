VERSION := $(shell git describe --tags)
BUILD := $(shell git rev-parse --short HEAD)
PROJECT := $(shell basename "$(PWD)")

GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
GOOS := "linux"
GOARCH := "amd64"

LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(Build)"

mirror:
	@go env -w GOPROXY=https://goproxy.cn,direct

install:
	@go get -u

build:
	@echo ">  Building binary"
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o wion *.go
upload:
	@ssh root@1.13.171.19 "systemctl stop wion"
	@scp -C wion app.toml root@1.13.171.19:/root/wion/server/
	@ssh root@1.13.171.19 "systemctl start wion"

upload-dev:
	@ssh -p 22022 aidapi@10.19.14.155 "sudo systemctl stop sfu"
	@scp -P 22022 -C wion aidapi@10.19.14.155:/home/aidapi/wion/
	@ssh -p 22022 aidapi@10.19.14.155 "sudo systemctl start sfu"
upload-test:
	@ssh root@jp.252798.xyz "systemctl stop wion"
	@scp -C wion app.toml config.toml root@jp.252798.xyz:/root/wion/
	@ssh root@jp.252798.xyz "systemctl start wion"
doc:
	@scp -P 22022 ~/Documents/insomnia/xr_scene.json aidapi@10.19.14.155:/aifs01/doc/xr-scene/insomnia.json


.PHONY: install build mirror upload upload-dev

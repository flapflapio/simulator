.PHONY: default build clean deploy docker docker-run docker-push dewit
default: build


SHELL				=	bash
STAGE				?=	dev
out					=	app
PORT				?=	8080
TAG 	?= latest
image 	= ghcr.io/flapflapio/simulator:$(TAG)

download:
	go get -d -v

build: download
	export GOOS=linux
	export GO111MODULE=on
	go build -o app

# Adds some flags for building the app statically linked (run: `file app` to
# ensure you are getting a static binary). This is needed for our multi-stage
# docker build, where we place only the executable and the config file inside a
# `FROM scratch` image
build-static: download
	export GOOS=linux
	export GO111MODULE=on
	export CGO_ENABLED=0
	go build \
		-ldflags="-extldflags=-static" \
		-tags osusergo,netgo \
		-o $(out)

clean:
	rm -rf ./$(out) ./vendor

docker:
	docker build -t "$(image)" -f docker/Dockerfile .
	docker image prune -f

docker-run:
	docker run \
		-p $(PORT):$(PORT) \
		--rm \
		--name simulator \
		$(image)

docker-push:
	docker push "$(image)"

dewit: docker docker-run

test:
	go clean --testcache && go test -v ./...

# TODO
install:
	cp $(out) /bin/asset-service
	mkdir --parents /etc/asset-service/
	cp config.yml /etc/asset-service/

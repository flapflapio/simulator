.PHONY: default build clean deploy docker docker-run dewit
default: build

SHELL	=	bash
STAGE	?=	dev
in		=	app/handlers
out		=	app

docker_image_tag = atomata/asset-service:latest

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
	docker build -t "$(docker_image_tag)" -f docker/Dockerfile .
	docker image prune -f

docker-run:
	docker run \
		-p 8080:8080 \
		-it \
		--rm \
		--env GIN_MODE=debug \
		--name asset-service \
		$(docker_image_tag)

dewit: docker docker-run

# NOTE: this build step runs tests but ignores packages within dir `test/`
# REASONING: I wanted to be able to provide a sub directory with mutliple
# test-related files (such as mocks and other tools) in each package.
test:
	go clean --testcache
	go test $$(go list ./... | grep -vP '(node_modules/|/test/|/test$$)')

# TODO
install:
	cp $(out) /bin/asset-service
	mkdir --parents /etc/asset-service/
	cp config.yml /etc/asset-service/

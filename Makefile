NAMESPACE		= ceptem
CATEGORY		= micro
REPOSITORY		= mailgun-forwarder
VERSION			= 1.0.0

SYSCONFDIR		= /etc/${NAMESPACE}/${CATEGORY}/${REPOSITORY}
IMAGE			= ${NAMESPACE}/${CATEGORY}_${REPOSITORY}:${VERSION}

DOCKER_BUILD		= docker build
DOCKER_BUILD_FLAGS	= -t "${IMAGE}"

all: build

build: Dockerfile
	${DOCKER_BUILD} ${DOCKER_BUILD_FLAGS} .

test-run: build
	docker run \
		--rm \
		-e 'PROFILEDIR=/conf/' \
		-e 'UID=$(shell id -u nobody)' \
		-e 'GID=$(shell id -g nobody)' \
		-e 'ADDRESS=0.0.0.0' \
		-e 'PORT=6543' \
		-v '$(shell pwd)/test/conf:/conf' \
		-p 6543:6543 \
	"${IMAGE}"

.PHONY: all build test-run

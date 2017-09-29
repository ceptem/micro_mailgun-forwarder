REPOSITORY		= ceptem
CATEGORY		= us
NAME			= mailgun-forward-email
VERSION			= 1.0.0

SYSCONFDIR		= /etc/${REPOSITORY}/${CATEGORY}/${NAME}
IMAGE			= ${REPOSITORY}/${CATEGORY}-${NAME}:${VERSION}

DOCKER_BUILD		= docker build
DOCKER_BUILD_FLAGS	= -t "${IMAGE}"

all: build

build: Dockerfile
	${DOCKER_BUILD} ${DOCKER_BUILD_FLAGS} .

test-run: build
	docker run --rm -v "`pwd`/test/conf:${SYSCONFDIR}" -p 6543:6543 "${IMAGE}"

.PHONY: all build test-run

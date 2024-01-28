ifeq '$(findstring ;,$(PATH))' ';'
    DETECTED_OS := windows
    DETECTED_ARCH := amd64
else
    DETECTED_OS := $(shell uname | tr '[:upper:]' '[:lower:]' 2> /dev/null || echo Unknown)
    DETECTED_OS := $(patsubst CYGWIN%,Cygwin,${DETECTED_OS})
    DETECTED_OS := $(patsubst MSYS%,MSYS,${DETECTED_OS})
    DETECTED_OS := $(patsubst MINGW%,MSYS,${DETECTED_OS})
    DETECTED_ARCH := $(shell dpkg --print-architecture 2>/dev/null || amd64)
endif

#colors:
B = \033[1;94m#   BLUE
G = \033[1;92m#   GREEN
Y = \033[1;93m#   YELLOW
R = \033[1;31m#   RED
M = \033[1;95m#   MAGENTA
K = \033[K#       ERASE END OF LINE
D = \033[0m#      DEFAULT
A = \007#         BEEP

APP :=ibot
# $(shell basename $(shell git remote get-url origin) | cut -d"." -f1)
# $(shell basename $(shell git remote get-url origin))
REGESTRY :=ghcr.io/k3ilona
BRANCH :=dev
VERSION=$(shell git describe --tags --abbrev=0 --always)-$(shell git rev-parse --short HEAD)-${BRANCH}
TARGETARCH := amd64 
TARGETOS=${DETECTED_OS}

format:
	gofmt -s -w ./

get:
	go get

lint:
	golint

test:
	go test -v

build: format get
	@printf "$GDetected OS/ARCH: $R${DETECTED_OS}/${DETECTED_ARCH}$D\n"
	@printf "$MDetected Version: $R${VERSION}$D\n"
	CGO_ENABLED=0 GOOS=${DETECTED_OS} GOARCH=${DETECTED_ARCH} go build -v -o ibot -ldflags "-X="github.com/k3ilona/publisher-slack-bot/cmd.appVersion=${VERSION}

linux: format get
	@printf "$GTarget OS/ARCH: $Rlinux/${DETECTED_ARCH}$D\n"
	CGO_ENABLED=0 GOOS=linux GOARCH=${DETECTED_ARCH} go build -v -o ibot -ldflags "-X="github.com/k3ilona/publisher-slack-bot/cmd.appVersion=${VERSION}
	docker build --build-arg name=linux -t ${REGESTRY}/${APP}:${VERSION}-linux-${DETECTED_ARCH} .

windows: format get
	@printf "$GTarget OS/ARCH: $Rwindows/${DETECTED_ARCH}$D\n"
	CGO_ENABLED=0 GOOS=windows GOARCH=${DETECTED_ARCH} go build -v -o ibot -ldflags "-X="github.com/k3ilona/publisher-slack-bot/cmd.appVersion=${VERSION}
	docker build --build-arg name=windows -t ${REGESTRY}/${APP}:${VERSION}-windows-${DETECTED_ARCH} .

darwin:format get
	@printf "$GTarget OS/ARCH: $Rdarwin/${DETECTED_ARCH}$D\n"
	CGO_ENABLED=0 GOOS=darwin GOARCH=${DETECTED_ARCH} go build -v -o ibot -ldflags "-X="github.com/k3ilona/publisher-slack-bot/cmd.appVersion=${VERSION}
	docker build --build-arg name=darwin -t ${REGESTRY}/${APP}:${VERSION}-darwin-${DETECTED_ARCH} .

arm: format get
	@printf "$GTarget OS/ARCH: $R${DETECTED_OS}/arm$D\n"
	CGO_ENABLED=0 GOOS=${DETECTED_OS} GOARCH=arm go build -v -o ibot -ldflags "-X="github.com/k3ilona/publisher-slack-bot/cmd.appVersion=${VERSION}
	docker build --build-arg name=arm -t ${REGESTRY}/${APP}:${VERSION}-${DETECTED_OS}-arm .

image:
	docker build . -t ${REGESTRY}/${APP}:${VERSION} --build-arg TARGETOS=${DETECTED_OS} --build-arg TARGETARCH=${TARGETARCH}

push:
	docker push ${REGESTRY}/${APP}:${VERSION}

dive: image
	IMG1=$$(docker images -q | head -n 1); \
	CI=true docker run -ti --rm -v /var/run/docker.sock:/var/run/docker.sock wagoodman/dive --ci --lowestEfficiency=0.99 $${IMG1}; \
	IMG2=$$(docker images -q | sed -n 2p); \
	docker rmi $${IMG1}; \
	docker rmi $${IMG2}

clean:
	@rm -rf ibot; \
	IMG1=$$(docker images -q | head -n 1); \
	if [ -n "$${IMG1}" ]; then  docker rmi -f $${IMG1}; else printf "$RImage not found$D\n"; fi

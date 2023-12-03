HOOK_DIR ?= ${USERPROFILE}\.githooks

build: hook.exe

hook.exe: main.go
	go build -ldflags "-s -w" -o hook.exe ./...

lint:
	golangci-lint run

test:
	go test ./...

${HOOK_DIR}:
	mkdir ${HOOK_DIR}

install: ${HOOK_DIR} build
	 copy hook.exe ${HOOK_DIR}\hook.exe

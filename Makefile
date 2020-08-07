# Copyright 2019 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
#

EXE_CLI = bin/kvgo-cli
APP_PATH = /usr/local/bin/kvgo-cli

all:
	go build -o ${EXE_CLI} cmd/cli/main.go
	strip -s ${EXE_CLI}
	upx -9 ${EXE_CLI}

install:
	install -m 755 ${EXE_CLI} ${APP_PATH}

clean:
	rm -f ${EXE_CLI}


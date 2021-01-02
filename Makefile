# Copyright 2019 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
#


BUILDCOLOR="\033[34;1m"
BINCOLOR="\033[37;1m"
ENDCOLOR="\033[0m"

ifndef V
	QUIET_BUILD = @printf '%b %b\n' $(BUILDCOLOR)BUILD$(ENDCOLOR) $(BINCOLOR)$@$(ENDCOLOR) 1>&2;
	QUIET_INSTALL = @printf '%b %b\n' $(BUILDCOLOR)INSTALL$(ENDCOLOR) $(BINCOLOR)$@$(ENDCOLOR) 1>&2;
endif

EXE_CLI = bin/kvgo-cli
APP_PATH = /usr/local/bin/kvgo-cli

BINDATA_CMD = httpsrv-bindata
BINDATA_ARGS_WEBUI = -src webui/ -dst bindata/webui/ -inc htm,js,css,svg

all: bin_build bindata_build
	@echo ""
	@echo "build complete"
	@echo ""

bin_build:
	$(QUIET_BUILD)go build -o ${EXE_CLI} cmd/kvgo-cli/main.go $(CCLINK)

bin_clean:
	rm -f ${EXE_CLI}

bindata_build:
	$(QUIET_BUILD)$(BINDATA_CMD) $(BINDATA_ARGS_WEBUI) $(CCLINK)

bindata_clean:
	rm -f bindata/webui/statik.go

install:
	install -m 755 ${EXE_CLI} ${APP_PATH}

clean: bin_clean bindata_clean
	@echo ""
	@echo "clean complete"
	@echo ""

// Copyright 2020 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package console

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hooto/hlog4g/hlog"
	"github.com/hooto/httpsrv"

	"github.com/lynkdb/kvgo-cli/config"
	"github.com/lynkdb/kvgo-cli/websrv/console"
)

func Start() {

	httpsrv.GlobalService.Config.UrlBasePath = "/kvgo"
	httpsrv.GlobalService.Config.HttpPort = uint16(config.Config.HttpPort)

	httpsrv.GlobalService.ModuleRegister("/", console.NewModule())

	go func() {
		hlog.Printf("info", "kvgo-cli start")
		if err := httpsrv.GlobalService.Start(); err != nil {
			hlog.Printf("info", "kvgo-cli start err %s", err.Error())
		}
	}()

	fmt.Printf("Setup Instances: %d\n",
		len(config.Config.Instances))

	fmt.Printf("Console URL: http://localhost:%d\n",
		config.Config.HttpPort)

	quit := make(chan os.Signal, 2)

	//
	signal.Notify(quit,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL)
	sg := <-quit

	time.Sleep(1e9)

	hlog.Printf("warn", "kvgo-cli signal quit %s", sg.String())
	hlog.Flush()
}

// Copyright 2019 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
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

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/gosuri/uitable"

	"github.com/lynkdb/kvgo-cli/cmd/accesskey"
	"github.com/lynkdb/kvgo-cli/cmd/instance"
	"github.com/lynkdb/kvgo-cli/cmd/table"
	"github.com/lynkdb/kvgo-cli/config"
	"github.com/lynkdb/kvgo-cli/data"
)

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func resetPrompt(l *readline.Instance) {
	instanceName := ""
	if data.DataInstance != "" {
		instanceName = "(" + data.DataInstance + ")"
	}
	l.SetPrompt("\033[31mkvgo-cli " + instanceName + ": \033[0m")
}

func main() {

	if err := config.Setup(); err != nil {
		log.Fatal(err)
	}

	if err := data.Setup(""); err != nil {
		fmt.Println(err)
	}

	l, err := readline.NewEx(&readline.Config{
		AutoComplete:        nil, // completer,
		HistoryFile:         fmt.Sprintf("%s/.%s.history", config.Homedir, config.AppName),
		InterruptPrompt:     "^C",
		EOFPrompt:           "exit",
		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	uitable.Separator = " | "

	for {
		resetPrompt(l)

		line, err := l.Readline()

		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		var (
			lineStr = strings.TrimSpace(line)
			lineArr = strings.Split(line, " ")
			out     string
		)

		switch {

		case lineStr == "instance list":
			out, err = instance.InstanceList()

		case strings.HasPrefix(lineStr, "instance use") && len(lineArr) == 3:
			out, err = instance.InstanceUse(lineArr[2])

		case lineStr == "instance new":
			out, err = instance.InstanceNew(l)

		case lineStr == "instance write-test-100":
			out, err = instance.InstanceWriteTest(100)

		case lineStr == "table list":
			out, err = table.TableList()

		case lineStr == "table set":
			out, err = table.TableSet(l)

		case lineStr == "access key list":
			out, err = accesskey.AccessKeyList()

		case lineStr == "access key get":
			out, err = accesskey.AccessKeyGet(l)

		case lineStr == "access key set":
			out, err = accesskey.AccessKeySet(l)

		case lineStr == "help", lineStr == "h":
			out, err = cmdHelp()

		case lineStr == "quit", lineStr == "exit":
			os.Exit(0)

		default:
			err = fmt.Errorf("unknown cmd %s\n", lineStr)
		}

		if err != nil {
			fmt.Println("Error:", err)
		} else if out != "" {
			fmt.Println(out)
		}
	}
}

func cmdHelp() (string, error) {
	return `kvgo-cli usage:
  instance list
  instance use <name>
  instance new
  table list
  table set
  access key list
  access key get
  access key set
  help
  quit
`, nil
}

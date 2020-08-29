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

package table

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/gosuri/uitable"
	"github.com/lynkdb/kvgo-cli/cmd/utils"
	"github.com/lynkdb/kvgo-cli/data"
	kv2 "github.com/lynkdb/kvspec/go/kvspec/v2"
)

func TableSet(l *readline.Instance) (string, error) {

	l.SetPrompt("table name: ")
	tableName, err := l.Readline()
	if err != nil {
		return "", err
	}

	l.SetPrompt("table desc: ")
	tableDesc, err := l.Readline()
	if err != nil {
		return "", err
	}

	req := kv2.NewSysCmdRequest("TableSet", &kv2.TableSetRequest{
		Name: tableName,
		Desc: tableDesc,
	})

	rs := data.Data.Connector().SysCmd(req)
	if !rs.OK() {
		return "", errors.New(rs.Message)
	}

	return fmt.Sprintf("OK %d", rs.Meta.IncrId), nil
}

func TableList() (string, error) {

	req := kv2.NewSysCmdRequest("TableList", &kv2.TableListRequest{})

	rs := data.Data.Connector().SysCmd(req)
	if !rs.OK() {
		return "", fmt.Errorf("error %s", rs.Message)
	}

	table := uitable.New()
	table.MaxColWidth = 40
	table.Wrap = true
	table.RightAlign(2)
	table.RightAlign(3)

	table.AddRow("ID", "Name", "Keys", "Size", "Log", "Incr", "Async",
		"Desc", "Created")

	sort.Slice(rs.Items, func(i, j int) bool {
		return rs.Items[i].Meta.IncrId < rs.Items[j].Meta.IncrId
	})

	for _, v := range rs.Items {

		var item kv2.TableItem
		if err := v.DataValue().Decode(&item, nil); err != nil {
			continue
		}

		if item.Status == nil {
			item.Status = &kv2.TableStatus{}
		}

		var (
			tc     = time.Unix(int64(v.Meta.Created)/1e3, 0)
			logid  = ""
			incrid = ""
			async  = ""
		)

		for k2, v2 := range item.Status.Options {

			if k2 == "log_id" {
				logid = fmt.Sprintf("%d", v2)
			} else if strings.HasPrefix(k2, "incr_id_") {
				if incrid != "" {
					incrid += "\n"
				}
				incrid += fmt.Sprintf("%s %d", strings.TrimPrefix(k2, "incr_id_"), v2)
			} else if strings.HasPrefix(k2, "async_") {
				if async != "" {
					async += "\n"
				}
				async += fmt.Sprintf("%s %d", strings.TrimPrefix(k2, "async_"), v2)
			}
		}

		table.AddRow(fmt.Sprintf("%d", v.Meta.IncrId),
			item.Name,
			fmt.Sprintf("%d", item.Status.KeyNum),
			utils.SizeTitle(int64(item.Status.DbSize)),
			logid,
			incrid,
			async,
			item.Desc,
			tc.Format("2006-01-02"),
		)
	}

	return table.String(), nil
}

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

package accesskey

import (
	"errors"
	"fmt"
	"strings"

	"github.com/chzyer/readline"
	"github.com/gosuri/uitable"
	"github.com/hooto/hauth/go/hauth/v1"
	"github.com/lynkdb/kvgo-cli/data"
	kv2 "github.com/lynkdb/kvspec/go/kvspec/v2"
)

func AccessKeyList() (string, error) {

	req := kv2.NewSysCmdRequest("AccessKeyList", nil)

	rs := data.Data.Connector().SysCmd(req)
	if !rs.OK() {
		return "", fmt.Errorf("error %s", rs.Message)
	}

	table := uitable.New()
	table.MaxColWidth = 32
	table.Wrap = true

	table.AddRow("Access Key ID", "Secret", "Roles", "Scope")

	for _, v := range rs.Items {

		var item hauth.AccessKey
		if err := v.DataValue().Decode(&item, nil); err != nil {
			continue
		}

		var (
			scopes = ""
		)
		for _, v2 := range item.Scopes {
			if scopes != "" {
				scopes += "\n"
			}
			scopes += fmt.Sprintf("%s = %s", v2.Name, v2.Value)
		}

		table.AddRow(
			item.Id,
			item.Secret,
			strings.Join(item.Roles, ","),
			scopes,
		)
	}

	return table.String(), nil
}

func AccessKeyGet(l *readline.Instance) (string, error) {

	l.SetPrompt("id: ")
	accessKeyId, err := l.Readline()
	if err != nil {
		return "", err
	}

	req := kv2.NewSysCmdRequest("AccessKeyList", nil)
	req.Body = []byte(accessKeyId)

	rs := data.Data.Connector().SysCmd(req)
	if !rs.OK() {
		return "", fmt.Errorf("error %s", rs.Message)
	}

	table := uitable.New()
	table.MaxColWidth = 50
	table.Wrap = true

	table.AddRow("Access Key ID", "Access Key Secret", "Roles", "Scopes")

	for _, v := range rs.Items {

		var item hauth.AccessKey
		if err := v.DataValue().Decode(&item, nil); err != nil {
			continue
		}

		var (
			scopes = ""
		)
		for _, v2 := range item.Scopes {
			if scopes != "" {
				scopes += "\n"
			}
			scopes += fmt.Sprintf("%s = %s", v2.Name, v2.Value)
		}

		table.AddRow(
			item.Id,
			item.Secret,
			strings.Join(item.Roles, ","),
			scopes,
		)
	}

	return table.String(), nil
}

func AccessKeySet(l *readline.Instance) (string, error) {

	l.SetPrompt("id: ")
	accessKeyId, err := l.Readline()
	if err != nil {
		return "", err
	}

	l.SetPrompt("roles: ")
	accessKeyRoles, err := l.Readline()
	if err != nil {
		return "", err
	}

	scopes := []*hauth.ScopeFilter{}
	for {

		l.SetPrompt(fmt.Sprintf("setup scope (ex: kvgo/table = table name) #%d : ", len(scopes)+1))

		accessKeyScope, err := l.Readline()
		if err != nil {
			return "", err
		}

		if accessKeyScope == "" {
			break
		}

		ar := strings.Split(accessKeyScope, "=")
		if len(ar) == 2 {
			scopes = append(scopes, &hauth.ScopeFilter{
				Name:  strings.TrimSpace(ar[0]),
				Value: strings.TrimSpace(ar[1]),
			})
		}
	}

	req := kv2.NewSysCmdRequest("AccessKeySet", &hauth.AccessKey{
		Id:     accessKeyId,
		Roles:  strings.Split(accessKeyRoles, ","),
		Scopes: scopes,
	})

	rs := data.Data.Connector().SysCmd(req)
	if !rs.OK() {
		return "", errors.New(rs.Message)
	}

	return fmt.Sprintf("OK %d", rs.Meta.Version), nil
}

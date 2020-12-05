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

package instance

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"time"

	"github.com/chzyer/readline"
	"github.com/gosuri/uitable"
	"github.com/hooto/hauth/go/hauth/v1"
	"github.com/lessos/lessgo/crypto/idhash"
	"github.com/lynkdb/kvgo"
	"github.com/lynkdb/kvgo-cli/config"
	"github.com/lynkdb/kvgo-cli/data"
)

var (
	InstanceNameRE = regexp.MustCompile("^[a-z]{1}[a-z0-9_]{0,31}$")
)

func InstanceNew(l *readline.Instance) (string, error) {

	validator := func(name, desc string, re *regexp.Regexp, fn func(v string) error) (string, error) {

		var (
			err   error
			value string
			tip   = fmt.Sprintf("input %s", name)
		)

		if desc != "" {
			tip += " (" + desc + ")"
		}

		tip += " : "

		for {

			l.SetPrompt(tip)
			value, err = l.Readline()

			if err != nil {
				fmt.Println(err)
				break
			}

			if re != nil && !re.MatchString(value) {
				fmt.Printf(" invalid %s, try again ...\n", name)
				continue
			}

			if fn != nil {
				if e := fn(value); e != nil {
					fmt.Printf("  invalid %s (%s), try again ...\n", name, e)
					continue
				}
			}

			break
		}

		return value, err
	}

	//
	name, err := validator("alias name of instance", "ex: prod, demo, ...",
		InstanceNameRE, nil)
	if err != nil {
		return "", err
	}

	//
	addr, err := validator("instace address", "ex. 127.0.0.1:9200",
		nil, func(v string) error {
			taddr, err := net.ResolveTCPAddr("tcp", v)
			if err != nil {
				return err
			}
			if taddr.IP == nil {
				return errors.New("invalid ip address")
			}
			if taddr.Port < 1 || taddr.Port > 65503 {
				return errors.New("invalid network port")
			}
			return nil
		})
	if err != nil {
		return "", err
	}

	//

	akId, err := validator("access key id", "",
		hauth.AccessKeyIdReg, nil)
	if err != nil {
		return "", err
	}

	//

	akSecret, err := validator("access key secret: ", "",
		hauth.AccessKeySecretRE, nil)
	if err != nil {
		return "", err
	}

	for _, v := range config.Config.Instances {
		if v.Name == name {
			return "", fmt.Errorf("instance name (%s) already exists", name)
		}
	}

	config.Config.Instances = append(config.Config.Instances, &config.ConfigInstance{
		Name: name,
		ClientConfig: &kvgo.ClientConfig{
			Addr: addr,
			AccessKey: &hauth.AccessKey{
				Id:     akId,
				Secret: akSecret,
			},
		},
	})

	err = config.Flush()
	if err == nil {
		return InstanceUse(name)
	}

	return "", err
}

func InstanceList() (string, error) {

	if len(config.Config.Instances) < 1 {
		return "", errors.New("no instance setup in" + config.ConfigFile)
	}

	table := uitable.New()
	table.MaxColWidth = 40
	table.Wrap = true

	table.AddRow("Name", "Address", "AccessKey ID", "Access Key Secret")

	for _, v := range config.Config.Instances {
		if v.AccessKey == nil {
			continue
		}
		var (
			akid  = ""
			aksec = ""
		)

		if len(v.AccessKey.Id) > 16 {
			akid = v.AccessKey.Id[:16] + "****"
		} else {
			akid = v.AccessKey.Id
		}
		if n := len(v.AccessKey.Secret); n > 8 {
			aksec = v.AccessKey.Secret[:4] + "****" + v.AccessKey.Secret[n-4:]
		} else {
			aksec = "********"
		}
		table.AddRow(v.Name, v.Addr, akid, aksec)
	}

	return table.String(), nil
}

func InstanceUse(name string) (string, error) {
	if err := data.Setup(name); err != nil {
		return "", err
	}
	return fmt.Sprintf("Use Instance %s", name), nil
}

func InstanceWriteTest(max int) (string, error) {
	var (
		tn = time.Now()
		n  = 0
		kp = "debug:test:"
		vl = 100
	)
	for i := 0; i < max; i++ {
		if rs := data.Data.NewWriter(
			[]byte(kp+idhash.RandHexString(16)),
			idhash.RandBase64String(vl)).Commit(); rs.OK() {
			n++
		} else {
			return "", errors.New(rs.Message)
		}
	}
	for i := 0; i < max; i++ {
		if rs := data.Data.NewWriter(
			[]byte(kp+"incr:"+idhash.RandHexString(16)),
			idhash.RandBase64String(vl)).
			IncrNamespaceSet("meta2").
			Commit(); rs.OK() {
			n++
		} else {
			return "", errors.New(rs.Message)
		}

	}
	return fmt.Sprintf("write %d records (key/value size %d/%d), time in %v\n",
		n, len(kp)+16, vl, time.Since(tn)), nil
}

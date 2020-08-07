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
	"time"

	"github.com/gosuri/uitable"
	"github.com/lessos/lessgo/crypto/idhash"
	"github.com/lynkdb/kvgo-cli/config"
	"github.com/lynkdb/kvgo-cli/data"
)

func InstanceList() (string, error) {

	if len(config.Config.Instances) < 1 {
		return "", errors.New("no instance setup in" + config.ConfigFile)
	}

	table := uitable.New()
	table.MaxColWidth = 30
	table.Wrap = true

	table.AddRow("Name", "Address", "Auth SecretKey")

	for _, v := range config.Config.Instances {
		sk := "********"
		if v.AuthKey == nil {
			continue
		}
		if len(v.AuthKey.SecretKey) > 32 {
			sk = v.AuthKey.SecretKey[:16] + " ..."
		}
		table.AddRow(v.Name, v.Addr, sk)
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

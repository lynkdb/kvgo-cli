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

package data

import (
	"errors"

	"github.com/hooto/hflag4g/hflag"

	"github.com/lynkdb/kvgo-cli/config"
	kv2 "github.com/lynkdb/kvspec/go/kvspec/v2"
)

var (
	Data         kv2.Client
	DataInstance = ""
	err          error
)

func Setup(instanceName string) error {

	if instanceName == "" {
		instanceName = hflag.Value("instance").String()
	}

	if instanceName == "" {
		if len(config.Config.Instances) < 1 {
			return errors.New("no instance config found in " + config.ConfigFile)
		}
		instanceName = config.Config.Instances[0].Name
	}

	var cfg *config.ConfigInstance
	for _, v := range config.Config.Instances {
		if v.Name == instanceName {
			cfg = v
			break
		}
	}

	if cfg == nil {
		return errors.New("no instance found")
	}

	if instanceName != "" {
		db, err := cfg.ClientConfig.NewClient()
		if err != nil {
			return err
		}

		if Data != nil {
			Data.Close()
		}

		Data = db
		DataInstance = instanceName
	}

	return nil
}

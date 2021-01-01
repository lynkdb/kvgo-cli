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

package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hooto/htoml4g/htoml"
	"github.com/lynkdb/kvgo"
)

var (
	AppName    = "kvgo-cli"
	Homedir    = ""
	Prefix     = ""
	err        error
	Config     ConfigCommon
	ConfigFile = ""
)

type ConfigCommon struct {
	HttpPort int `toml:"http_port,omitempty"`
	Client   struct {
		LastActiveInstance string `toml:"last_active_instance" json:"last_active_instance"`
	} `toml:"client" json:"client"`
	Instances []*ConfigInstance `toml:"instances" json:"instances"`
}

type ConfigInstance struct {
	*kvgo.ClientConfig
	Name string `toml:"name" json:"name"`
}

func Setup() error {

	Homedir, err = os.UserHomeDir()
	if err != nil {
		return err
	}

	ConfigFile = fmt.Sprintf("%s/.%s.conf", Homedir, AppName)

	if ConfigFile, err = filepath.Abs(ConfigFile); err != nil {
		return err
	}

	if Config.HttpPort == 0 {
		Config.HttpPort = 9201
	}

	if err = htoml.DecodeFromFile(&Config, ConfigFile); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

func Flush() error {
	return htoml.EncodeToFile(Config, ConfigFile, nil)
}

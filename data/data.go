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
	"fmt"
	"sync"

	"github.com/hooto/hflag4g/hflag"

	"github.com/lynkdb/kvgo-cli/config"
	kv2 "github.com/lynkdb/kvspec/go/kvspec/v2"
)

var (
	Data         kv2.Client
	DataInstance = ""
	err          error
	mu           sync.RWMutex
	dbSets       = map[string]kv2.Client{}
)

func Setup(instanceName string) error {

	if instanceName == "" {

		instanceName = hflag.Value("instance").String()

		if instanceName == "" {

			if config.Config.Client.LastActiveInstance != "" {
				instanceName = config.Config.Client.LastActiveInstance
			} else if len(config.Config.Instances) < 1 {
				return fmt.Errorf("no instance setup in %s, try to use 'instance new' to create new connection to kvgo-server",
					config.ConfigFile)
			} else {
				instanceName = config.Config.Instances[0].Name
			}
		}
	}

	var cfg *config.ConfigInstance
	for _, v := range config.Config.Instances {
		if v.Name == instanceName {
			cfg = v
			break
		}
	}

	if cfg == nil {
		return fmt.Errorf("no instance (%s) found, try to use 'instance new' to create new connection to kvgo-server",
			instanceName)
	}

	if Data == nil || instanceName != DataInstance {
		db, err := cfg.ClientConfig.NewClient()
		if err != nil {
			return err
		}

		if Data != nil {
			Data.Close()
		}

		Data = db
		DataInstance = instanceName

		if instanceName != config.Config.Client.LastActiveInstance {
			config.Config.Client.LastActiveInstance = instanceName
			config.Flush()
		}
	}

	return nil
}

func Connector(instanceName string) (kv2.ClientConnector, error) {

	if instanceName == "" {

		instanceName = hflag.Value("instance").String()

		if instanceName == "" {

			if config.Config.Client.LastActiveInstance != "" {
				instanceName = config.Config.Client.LastActiveInstance
			} else if len(config.Config.Instances) < 1 {
				return nil, fmt.Errorf("no instance setup in %s, try to use 'instance new' to create new connection to kvgo-server",
					config.ConfigFile)
			} else {
				instanceName = config.Config.Instances[0].Name
			}
		}
	}

	var cfg *config.ConfigInstance
	for _, v := range config.Config.Instances {
		if v.Name == instanceName {
			cfg = v
			break
		}
	}

	if cfg == nil {
		return nil, fmt.Errorf("no instance (%s) found, try to use 'instance new' to create new connection to kvgo-server",
			instanceName)
	}

	mu.Lock()
	defer mu.Unlock()
	db := dbSets[instanceName]

	if db == nil || instanceName != DataInstance {

		db, err = cfg.ClientConfig.NewClient()
		if err != nil {
			return nil, err
		}

		dbSets[instanceName] = db

		DataInstance = instanceName

		if instanceName != config.Config.Client.LastActiveInstance {
			config.Config.Client.LastActiveInstance = instanceName
			config.Flush()
		}
	}

	return db.Connector(), nil
}

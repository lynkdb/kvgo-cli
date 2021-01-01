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
	"github.com/hooto/httpsrv"
	"github.com/lessos/lessgo/types"
	"github.com/valuedig/apis/go/tsd/v1"

	kv2 "github.com/lynkdb/kvspec/go/kvspec/v2"

	"github.com/lynkdb/kvgo-cli/config"
	"github.com/lynkdb/kvgo-cli/data"
)

type Sys struct {
	*httpsrv.Controller
}

func (c Sys) InstanceListAction() {

	rsp := types.WebServiceResult{
		Kind: "InstanceList",
	}

	for _, v := range config.Config.Instances {
		rsp.Items = append(rsp.Items, &config.ConfigInstance{
			Name: v.Name,
		})
	}

	c.RenderJson(rsp)
}

func (c Sys) StatusAction() {

	var (
		rsp = types.WebServiceResult{
			Kind: "SysStatus",
		}
		instanceName = c.Params.Get("instance_name")
	)
	defer c.RenderJson(&rsp)

	if instanceName == "" {
		rsp.Kind, rsp.Message = "400", "no instance_name found"
		return
	}

	db, err := data.Connector(instanceName)
	if err != nil {
		rsp.Kind, rsp.Message = "400.1", err.Error()
		return
	}

	rs := db.SysCmd(kv2.NewSysCmdRequest("SysStatus", nil))
	if !rs.OK() {
		rsp.Kind, rsp.Message = "400.2", rs.Message
		return
	}

	var item kv2.SysStatus
	if err := rs.DataValue().Decode(&item, nil); err != nil {
		rsp.Kind, rsp.Message = "400", err.Error()
		return
	}

	rsp.Item = item
}

func (c Sys) MetricsAction() {

	var (
		rsp = types.WebServiceResult{
			Kind: "SysMetrics",
		}
		instanceName = c.Params.Get("instance_name")
	)
	defer c.RenderJson(&rsp)

	if instanceName == "" {
		rsp.Kind, rsp.Message = "400", "no instance_name found"
		return
	}

	db, err := data.Connector(instanceName)
	if err != nil {
		rsp.Kind, rsp.Message = "400.1", err.Error()
		return
	}

	opts := tsd.NewCycleExportOptionsFromHttp(c.Request.Request)

	rs := db.SysCmd(kv2.NewSysCmdRequest("SysMetrics", opts))
	if !rs.OK() {
		rsp.Kind, rsp.Message = "400.2", rs.Message
		return
	}

	var item tsd.CycleFeed
	if err := rs.DataValue().Decode(&item, nil); err != nil {
		rsp.Kind, rsp.Message = "400", err.Error()
		return
	}

	rsp.Item = item
}

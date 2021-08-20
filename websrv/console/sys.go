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
	tsd2 "github.com/valuedig/apis/go/tsd/v2"

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
		instanceName    = c.Params.Get("instance_name")
		lastTimeRange   = c.Params.Int64("last_time_range")
		alignmentPeriod = c.Params.Int64("alignment_period")
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

	req := tsd2.NewSampleQueryRequest()

	req.AlignmentPeriod = alignmentPeriod
	if req.AlignmentPeriod < 10 {
		req.AlignmentPeriod = 10
	}

	req.LastTimeRange = lastTimeRange
	if req.LastTimeRange < 10 {
		req.LastTimeRange = 600
	}

	req.MetricSelect("LogSyncCall").LabelSelect("*")
	req.MetricSelect("LogSyncLatency").LabelSelect("*")

	req.MetricSelect("ServiceCall").LabelSelect("*")
	req.MetricSelect("ServiceLatency").LabelSelect("*")

	req.MetricSelect("StorageCall").LabelSelect("*")
	req.MetricSelect("StorageLatency").LabelSelect("*")

	rs := db.SysCmd(kv2.NewSysCmdRequest("SysMetrics", req))
	if !rs.OK() {
		rsp.Kind, rsp.Message = "400.2", rs.Message
		return
	}

	var item tsd2.MetricSet
	if err := rs.DataValue().Decode(&item, nil); err != nil {
		rsp.Kind, rsp.Message = "400", err.Error()
		return
	}

	rsp.Item = item
}

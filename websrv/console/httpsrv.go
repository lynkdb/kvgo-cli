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

	"github.com/lynkdb/kvgo-cli/bindata"
	"github.com/lynkdb/kvgo-cli/config"
)

func NewModule() httpsrv.Module {

	module := httpsrv.NewModule("valueui_index")

	module.RouteSet(httpsrv.Route{
		Type:       httpsrv.RouteTypeStatic,
		Path:       "~",
		StaticPath: config.Prefix + "/webui",
		BinFs:      bindata.NewFs("webui"),
	})

	module.ControllerRegister(new(Index))
	module.ControllerRegister(new(Sys))

	return module
}

type Index struct {
	*httpsrv.Controller
}

func (c Index) IndexAction() {

	c.AutoRender = false
	c.Response.Out.Header().Set("Cache-Control", "no-cache")

	c.RenderString(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>kvgo console</title>
  <script src="/kvgo/~/valueui/main.js"></script>
  <script type="text/javascript">
    valueui.app_version = "` + config.AppVersion + `";
    valueui.basepath = "/kvgo/~/";
    window.onload = valueui.use("kvgo/main.js", function() {kvgo.Boot()});
  </script>
</head>
<body id="valueui-body">
<div class="incp-well" id="incp-well">
<div class="incp-well-box">
  <div class="incp-well-panel">
    <div class="body2c">
      <div class="body2c1">
      </div>
      <div class="body2c2">
        <div>lynkdb<br/>Enterprise Data Mining Engine</div>
      </div>
    </div>
    <div class="status status_dark" id="incp-well-status">loading</div>
  </div>
  <div class="footer">
  <span class="url-info">Powered by <a href="https://github.com/lynkdb" target="_blank">lynkdb</a></span>
  </div>
</div>
</div>
</body>
</html>
`)
}

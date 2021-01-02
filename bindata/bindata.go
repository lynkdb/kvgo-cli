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

package bindata

import (
	"net/http"
	"os"

	"github.com/hooto/hlog4g/hlog"
	"github.com/rakyll/statik/fs"

	_ "github.com/lynkdb/kvgo-cli/bindata/webui"
)

func NewFs(ns string) http.FileSystem {

	if os.Args[0] != "go" {
		binFs, err := fs.NewWithNamespace(ns)
		if err == nil && binFs != nil {
			hlog.Printf("info", "bindata load %s done", ns)
			return binFs
		}
	}

	return nil
}

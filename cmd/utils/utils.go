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

package utils

import (
	"fmt"

	kv2 "github.com/lynkdb/kvspec/go/kvspec/v2"
)

func SizeTitle(siz int64) string {

	sizeS := "0"
	if siz > kv2.TiB {
		sizeS = fmt.Sprintf("%d TB", siz/kv2.TiB)
	} else if siz > kv2.GiB {
		sizeS = fmt.Sprintf("%d GB", siz/kv2.GiB)
	} else if siz > kv2.MiB {
		sizeS = fmt.Sprintf("%d MB", siz/kv2.MiB)
	} else if siz > kv2.KiB {
		sizeS = fmt.Sprintf("%d KB", siz/kv2.KiB)
	}
	return sizeS
}

// Copyright 2016-2019 Yincheng Fang
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package extension

import (
	"github.com/apache/dubbo-go/common/gr_pool"
)

var (
	grPools = make(map[string]func() gr_pool.GrPool)
)

func SetGrPool(name string, v func() gr_pool.GrPool) {
	grPools[name] = v
}

func GetGrPool(name string) gr_pool.GrPool {
	if grPools[name] == nil {
		panic("thread_pool for " + name + " is not existing, make sure you have import the package.")
	}
	return grPools[name]()
}

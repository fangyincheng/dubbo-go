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

package limited

import (
	"strconv"
)

import (
	"github.com/dubbo/go-for-apache-dubbo/common"
	"github.com/dubbo/go-for-apache-dubbo/common/constant"
	"github.com/dubbo/go-for-apache-dubbo/common/extension"
	"github.com/dubbo/go-for-apache-dubbo/common/gr_pool"
	"github.com/dubbo/go-for-apache-dubbo/common/logger"
)

const LIMITED = "limited"

func init() {
	extension.SetGrPool(LIMITED, GetGrPool)
}

type FixPool struct {
	keeper *Keeper
}

func (fp *FixPool) CreatePool(url common.URL) {
	core, err := strconv.Atoi(url.GetParam(constant.COREGRS_KEY, constant.DEFAULT_COREGRS))
	if err != nil {
		logger.Errorf("[Execute] error: %v", err)
		panic(err)
	}
	grs, err := strconv.Atoi(url.GetParam(constant.GRS_KEY, constant.DEFAULT_GRS))
	if err != nil {
		logger.Errorf("[Execute] error: %v", err)
		panic(err)
	}
	queues, err := strconv.Atoi(url.GetParam(constant.QUEUES_KEY, constant.DEFAULT_QUEUES))
	if err != nil {
		logger.Errorf("[Execute] error: %v", err)
		panic(err)
	}
	fp.keeper = NewKeeper(core, grs, queues)
}

func (fp *FixPool) Close() {
	fp.keeper.Close()
}

func (fp *FixPool) Push(f func(index int)) error {
	return fp.keeper.PushTask(f, 0)
}

func GetGrPool() gr_pool.GrPool {
	return new(FixPool)
}

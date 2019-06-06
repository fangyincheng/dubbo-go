/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package limited

import (
	"strconv"
)

import (
	"github.com/apache/dubbo-go/common"
	"github.com/apache/dubbo-go/common/constant"
	"github.com/apache/dubbo-go/common/extension"
	"github.com/apache/dubbo-go/common/gr_pool"
	"github.com/apache/dubbo-go/common/logger"
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

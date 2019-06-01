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

package all

import (
	"github.com/dubbo/go-for-apache-dubbo/common"
	"github.com/dubbo/go-for-apache-dubbo/common/extension"
	"github.com/dubbo/go-for-apache-dubbo/protocol"
	"github.com/dubbo/go-for-apache-dubbo/protocol/dispatcher"
)

const ALL = "all"

func init() {
	extension.SetDispatcher(ALL, GetDispatcher)
}

var allDispatcher *AllDispatcher

type AllDispatcher struct {
}

func (ad *AllDispatcher) Dispatch(handler protocol.Handler, url common.URL) protocol.Handler {
	return NewAllHandler(handler, url)
}

func GetDispatcher() dispatcher.Dispatcher {
	if allDispatcher == nil {
		allDispatcher = &AllDispatcher{}
	}
	return allDispatcher
}

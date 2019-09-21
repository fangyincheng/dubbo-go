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

package getty

import (
	"bytes"
)

import (
	"github.com/dubbogo/getty"
	perrors "github.com/pkg/errors"
)
import (
	"github.com/apache/dubbo-go/common/logger"
)

////////////////////////////////////////////
// RpcClientPackageHandler
////////////////////////////////////////////

type RpcClientPackageHandler struct {
	client *Client
}

func NewRpcClientPackageHandler(client *Client) *RpcClientPackageHandler {
	return &RpcClientPackageHandler{client: client}
}

func (p *RpcClientPackageHandler) Read(ss getty.Session, data []byte) (interface{}, int, error) {
	p.client.pendingLock.RLock()
	defer p.client.pendingLock.RUnlock()
	pkg := &DubboPackage{}

	buf := bytes.NewBuffer(data)
	err := pkg.Unmarshal(buf, p.client)
	if err != nil {
		pkg.Err = perrors.WithStack(err) // client will get this err
		return pkg, len(data), nil
	}

	return pkg, len(data), nil
}

func (p *RpcClientPackageHandler) Write(ss getty.Session, pkg interface{}) error {
	req, ok := pkg.(*DubboPackage)
	if !ok {
		logger.Errorf("illegal pkg:%+v\n", pkg)
		return perrors.New("invalid rpc request")
	}

	buf, err := req.Marshal()
	if err != nil {
		logger.Warnf("binary.Write(req{%#v}) = err{%#v}", req, perrors.WithStack(err))
		return perrors.WithStack(err)
	}

	return perrors.WithStack(ss.WriteBytes(buf.Bytes()))
}

////////////////////////////////////////////
// RpcServerPackageHandler
////////////////////////////////////////////

type RpcServerPackageHandler struct {
}

func NewRpcServerPackageHandler() *RpcServerPackageHandler {
	return &RpcServerPackageHandler{}
}

func (p *RpcServerPackageHandler) Read(ss getty.Session, data []byte) (interface{}, int, error) {
	pkg := &DubboPackage{
		Body: make([]interface{}, 7),
	}
	buf := bytes.NewBuffer(data)
	err := pkg.UnmarshalHeader(buf)
	if err != nil {
		return nil, 0, perrors.WithStack(err)
	}

	return pkg, len(data), nil
}

func (p *RpcServerPackageHandler) Write(ss getty.Session, pkg interface{}) error {
	if b, ok := pkg.([]byte); ok {
		return perrors.WithStack(ss.WriteBytes(b))
	}
	return pkg.(error)
}

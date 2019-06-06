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

package dubbo

import (
	"bufio"
	"bytes"
	"fmt"
	"reflect"
	"time"
)

import (
	"github.com/dubbogo/hessian2"
	perrors "github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go/common"
	"github.com/apache/dubbo-go/common/constant"
	"github.com/apache/dubbo-go/common/logger"
)

// serial ID
type SerialID byte

const (
	S_Dubbo SerialID = 2
)

// call type
type CallType int32

const (
	CT_UNKOWN CallType = 0
	CT_OneWay CallType = 1
	CT_TwoWay CallType = 2
)

////////////////////////////////////////////
// dubbo package
////////////////////////////////////////////

type SequenceType int64

type DubboPackage struct {
	Codec   *hessian.HessianCodec
	Header  hessian.DubboHeader
	Service hessian.Service
	Body    interface{}
	Err     error
}

func (p DubboPackage) String() string {
	return fmt.Sprintf("DubboPackage: Header-%v, Path-%v, Body-%v", p.Header, p.Service, p.Body)
}

func (p *DubboPackage) Marshal() (*bytes.Buffer, error) {
	codec := hessian.NewHessianCodec(nil)

	pkg, err := codec.Write(p.Service, p.Header, p.Body)
	if err != nil {
		return nil, perrors.WithStack(err)
	}

	return bytes.NewBuffer(pkg), nil
}

func (p *DubboPackage) UnmarshalHeader(buf *bytes.Buffer) error {
	codec := hessian.NewHessianCodec(bufio.NewReader(buf))

	// read header
	err := codec.ReadHeader(&p.Header)
	if err != nil {
		return perrors.WithStack(err)
	}
	p.Codec = codec
	return nil
}

func (p *DubboPackage) UnmarshalBody(opts ...interface{}) error {
	if p.Codec == nil {
		return perrors.New("[UnmarshalBody] no codec!")
	}

	if len(opts) != 0 { // for client
		if client, ok := opts[0].(*Client); ok {

			r := client.pendingResponses[SequenceType(p.Header.ID)]
			if r == nil {
				return perrors.Errorf("pendingResponses[%v] = nil", p.Header.ID)
			}
			p.Body = client.pendingResponses[SequenceType(p.Header.ID)].reply
		} else {
			return perrors.Errorf("opts[0] is not *Client")
		}
	}

	if p.Header.Type&hessian.PackageHeartbeat != 0x00 {
		return nil
	}

	// read body
	err := p.Codec.ReadBody(p.Body)
	return perrors.WithStack(err)
}

func (p *DubboPackage) Unmarshal(buf *bytes.Buffer, opts ...interface{}) error {
	codec := hessian.NewHessianCodec(bufio.NewReader(buf))

	// read header
	err := codec.ReadHeader(&p.Header)
	if err != nil {
		return perrors.WithStack(err)
	}

	if len(opts) != 0 { // for client
		if client, ok := opts[0].(*Client); ok {

			r := client.pendingResponses[SequenceType(p.Header.ID)]
			if r == nil {
				return perrors.Errorf("pendingResponses[%v] = nil", p.Header.ID)
			}
			p.Body = client.pendingResponses[SequenceType(p.Header.ID)].reply
		} else {
			return perrors.Errorf("opts[0] is not *Client")
		}
	}

	if p.Header.Type&hessian.PackageHeartbeat != 0x00 {
		return nil
	}

	// read body
	err = codec.ReadBody(p.Body)
	return perrors.WithStack(err)
}

////////////////////////////////////////////
// PendingResponse
////////////////////////////////////////////

type PendingResponse struct {
	seq       uint64
	err       error
	start     time.Time
	readStart time.Time
	callback  AsyncCallback
	reply     interface{}
	opts      CallOptions
	done      chan struct{}
}

func NewPendingResponse() *PendingResponse {
	return &PendingResponse{
		start: time.Now(),
		done:  make(chan struct{}),
	}
}

func (r PendingResponse) GetCallResponse() CallResponse {
	return CallResponse{
		Opts:      r.opts,
		Cause:     r.err,
		Start:     r.start,
		ReadStart: r.readStart,
		Reply:     r.reply,
	}
}

////////////////////////////////////////////
// Codec
////////////////////////////////////////////

func DecodeBody(pkg interface{}) error {
	req, ok := pkg.(*DubboPackage)
	if !ok {
		logger.Errorf("illegal pkg:%+v\n, it is %+v", pkg, reflect.TypeOf(pkg))
		return perrors.New("invalid rpc request")
	}
	err := req.UnmarshalBody()
	if err != nil {
		return err
	}
	// convert params of request
	content := req.Body.([]interface{}) // length of body should be 7
	if len(content) > 0 {
		var dubboVersion, argsTypes string
		var args []interface{}
		var attachments map[interface{}]interface{}
		if content[0] != nil {
			dubboVersion = content[0].(string)
		}
		if content[1] != nil {
			req.Service.Path = content[1].(string)
		}
		if content[2] != nil {
			req.Service.Version = content[2].(string)
		}
		if content[3] != nil {
			req.Service.Method = content[3].(string)
		}
		if content[4] != nil {
			argsTypes = content[4].(string)
		}
		if content[5] != nil {
			args = content[5].([]interface{})
		}
		if content[6] != nil {
			attachments = content[6].(map[interface{}]interface{})
		}
		if interf, ok := attachments[constant.INTERFACE_KEY]; ok {
			req.Service.Interface = interf.(string)
		}
		req.Body = map[string]interface{}{
			"dubboVersion": dubboVersion,
			"argsTypes":    argsTypes,
			"args":         args,
			"service":      common.ServiceMap.GetService(DUBBO, req.Service.Interface),
			"attachments":  attachments,
		}
	}

	return nil
}

func Decode(data []byte) (interface{}, error) {
	pkg := &DubboPackage{
		Body: make([]interface{}, 7),
	}

	buf := bytes.NewBuffer(data)
	err := pkg.Unmarshal(buf)
	if err != nil {
		return nil, err
	}
	// convert params of request
	req := pkg.Body.([]interface{}) // length of body should be 7
	if len(req) > 0 {
		var dubboVersion, argsTypes string
		var args []interface{}
		var attachments map[interface{}]interface{}
		if req[0] != nil {
			dubboVersion = req[0].(string)
		}
		if req[1] != nil {
			pkg.Service.Path = req[1].(string)
		}
		if req[2] != nil {
			pkg.Service.Version = req[2].(string)
		}
		if req[3] != nil {
			pkg.Service.Method = req[3].(string)
		}
		if req[4] != nil {
			argsTypes = req[4].(string)
		}
		if req[5] != nil {
			args = req[5].([]interface{})
		}
		if req[6] != nil {
			attachments = req[6].(map[interface{}]interface{})
		}
		if interf, ok := attachments[constant.INTERFACE_KEY]; ok {
			pkg.Service.Interface = interf.(string)
		}
		pkg.Body = map[string]interface{}{
			"dubboVersion": dubboVersion,
			"argsTypes":    argsTypes,
			"args":         args,
			"service":      common.ServiceMap.GetService(DUBBO, pkg.Service.Interface),
			"attachments":  attachments,
		}
	}

	return pkg, nil
}

func Encode(pkg interface{}) interface{} {
	res, ok := pkg.(*DubboPackage)
	if !ok {
		logger.Errorf("illegal pkg:%+v\n, it is %+v", pkg, reflect.TypeOf(pkg))
		return perrors.New("invalid rpc response")
	}

	buf, err := res.Marshal()
	if err != nil {
		logger.Warnf("binary.Write(res{%#v}) = err{%#v}", res, perrors.WithStack(err))
		return perrors.WithStack(err)
	}

	return buf.Bytes()
}

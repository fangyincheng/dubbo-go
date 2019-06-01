// Copyright 2016-2019 Alex Stocks
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
	"fmt"
	"strconv"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	suite.Suite
	keeper *Keeper
}

func (suite *KeeperTestSuite) SetupSuite() {
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.keeper = NewKeeper(4, 0, 0)
}

func (suite *KeeperTestSuite) TearDownTest() {
	suite.keeper.Close()
	pendingNum := suite.keeper.PendingTaskNum()
	suite.T().Logf("pendingNum = %d", pendingNum)
}

func (suite *KeeperTestSuite) TearDownSuite() {
}

func (suite *KeeperTestSuite) TestKeeper() {
	f := func(d time.Duration, req string) Task {
		return func(id int) {
			time.Sleep(d)
			fmt.Printf("worker id:%d, req:%q\n", id, req)
			return
		}
	}

	for i := 0; i < 8; i++ {
		task := f(time.Duration(1e8*i), "f"+strconv.Itoa(i))
		err := suite.keeper.PushTask(task, 1e9)
		suite.Equalf(nil, err, "err != nil")
	}
	time.Sleep(8e9)
}

func (suite *KeeperTestSuite) TestPendingNum() {
	f := func(d time.Duration, req string) Task {
		return func(id int) {
			time.Sleep(d)
			fmt.Printf("worker id:%d, req:%q\n", id, req)
			return
		}
	}

	for i := 0; i < 12; i++ {
		task := f(time.Duration(1e8*i), "f"+strconv.Itoa(i))
		err := suite.keeper.PushTask(task, 1e9)
		suite.Equalf(nil, err, "err != nil")
	}
	time.Sleep(10e8)
	pendingNum := suite.keeper.PendingTaskNum()
	suite.T().Logf("pendingNum = %d", pendingNum)
	suite.Equalf(true, 1 <= pendingNum, "pendingNum = %d", pendingNum)
}

func (suite *KeeperTestSuite) TestPanic() {
	f := func(d time.Duration, req string) Task {
		return func(id int) {
			time.Sleep(d)
			fmt.Printf("worker id:%d, req:%q\n", id, req)
			if id == 1 {
				panic("PONG! PONG! PONG! I AM DIEING!")
			}
			return
		}
	}

	for i := 0; i < 8; i++ {
		task := f(time.Duration(1e9*i), "f"+strconv.Itoa(i))
		err := suite.keeper.PushTask(task, 1e9)
		suite.Equalf(nil, err, "err != nil")
	}
	time.Sleep(10e9)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

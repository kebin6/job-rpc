// Copyright 2023 The Ryan SU Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mqtask

import (
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"

	"github.com/suyuan32/simple-admin-job/internal/svc"
)

type MQTask struct {
	svcCtx *svc.ServiceContext
	mux    *asynq.ServeMux
}

func NewMQTask(svcCtx *svc.ServiceContext) *MQTask {
	return &MQTask{
		svcCtx: svcCtx,
	}
}

// Start starts the server.
func (m *MQTask) Start() {
	m.Register()
	if err := m.svcCtx.AsynqServer.Run(m.mux); err != nil {
		log.Fatal(fmt.Errorf("failed to start mqtask server, error: %v", err))
	}
}

// Stop stops the server.
func (m *MQTask) Stop() {
	time.Sleep(5 * time.Second)
	m.svcCtx.AsynqServer.Stop()
	m.svcCtx.AsynqServer.Shutdown()
}

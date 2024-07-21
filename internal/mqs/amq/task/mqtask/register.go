package mqtask

import (
	"github.com/hibiken/asynq"
	"github.com/suyuan32/simple-admin-job/internal/mqs/amq/handler/amq/game"
	"github.com/suyuan32/simple-admin-job/internal/mqs/amq/handler/amq/game/coingame"
	"github.com/suyuan32/simple-admin-job/internal/mqs/amq/handler/amq/game/tokengame"
	"github.com/suyuan32/simple-admin-job/internal/mqs/amq/types/pattern"
)

// Register adds task to cron. | 在此处定义任务处理逻辑，注册worker.
func (m *MQTask) Register() {
	mux := asynq.NewServeMux()

	// define the handler | 定义处理逻辑
	// mux.Handle(pattern.RecordHelloWorld, base.NewHelloWorldHandler(m.svcCtx))
	mux.Handle(pattern.ProcessCoinGame, coingame.NewProcessGameHandler(m.svcCtx))
	mux.Handle(pattern.ProcessTokenGame, tokengame.NewProcessGameHandler(m.svcCtx))
	mux.Handle(pattern.ProcessSyncGCICS, game.NewSyncGcicsHandler(m.svcCtx))

	m.mux = mux
}

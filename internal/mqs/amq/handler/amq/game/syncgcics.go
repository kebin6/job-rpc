package game

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/kebin6/wolflamp-rpc/types/wolflamp"
	"github.com/suyuan32/simple-admin-common/i18n"
	"github.com/suyuan32/simple-admin-job/ent/task"
	"github.com/suyuan32/simple-admin-job/internal/mqs/amq/types/pattern"
	"github.com/suyuan32/simple-admin-job/internal/mqs/amq/types/payload"
	"github.com/suyuan32/simple-admin-job/internal/svc"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

type SyncGcicsHandler struct {
	svcCtx  *svc.ServiceContext
	taskId  uint64
	payload payload.SyncGcicsPayload
}

func NewSyncGcicsHandler(svcCtx *svc.ServiceContext) *SyncGcicsHandler {
	taskInfo, err := svcCtx.DB.Task.Query().Where(task.PatternEQ(pattern.ProcessCoinGame)).First(context.Background())
	if err != nil || taskInfo == nil {
		return nil
	}

	var p payload.SyncGcicsPayload
	if err := json.Unmarshal([]byte(taskInfo.Payload), &p); err != nil {
		p = payload.SyncGcicsPayload{
			ChunkSize: 1,
		}
	}

	return &SyncGcicsHandler{
		svcCtx:  svcCtx,
		taskId:  taskInfo.ID,
		payload: p,
	}
}

// ProcessTask if return err != nil , asynq will retry | 如果返回错误不为空则会重试
func (l *SyncGcicsHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {

	if l == nil {
		return nil
	}
	if l.taskId == 0 {
		logx.Errorw("failed to load task info")
		return errorx.NewInternalError(i18n.DatabaseError)
	}

	sync, err := l.svcCtx.WolfLampRpc.SyncGcics(ctx, &wolflamp.SyncGcicsReq{ChunkSize: uint32(l.payload.ChunkSize)})
	if err != nil {
		fmt.Printf("ProcessSync Error: %v \n", err)
		return err
	}
	fmt.Printf("ProcessSync: Msg=%s \n", sync.Msg)
	return nil

}

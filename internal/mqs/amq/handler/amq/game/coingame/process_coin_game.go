package coingame

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/kebin6/wolflamp-rpc/common/enum/cachekey"
	"github.com/kebin6/wolflamp-rpc/common/enum/poolenum"
	"github.com/kebin6/wolflamp-rpc/common/enum/roundenum"
	"github.com/kebin6/wolflamp-rpc/types/wolflamp"
	"github.com/suyuan32/simple-admin-common/i18n"
	"github.com/suyuan32/simple-admin-job/ent/task"
	"github.com/suyuan32/simple-admin-job/internal/mqs/amq/types/pattern"
	"github.com/suyuan32/simple-admin-job/internal/mqs/amq/types/payload"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
	"math/rand"
	"time"

	"github.com/suyuan32/simple-admin-job/internal/svc"
)

type ProcessGameHandler struct {
	svcCtx  *svc.ServiceContext
	taskId  uint64
	payload payload.ProcessGamePayload
	mode    string
}

// InvestPrepare invest prepare | 投注实体
type InvestPrepare struct {
	RoundId    uint64
	PlayerId   uint64
	LambNum    uint32
	LambFoldNo uint32
}

// LambFoldAggregateInfo | 羊圈的聚合统计数据
type LambFoldAggregateInfo struct {
	// 羊圈编号
	LambFoldNo uint32
	// 盈亏数
	ProfitAndLoss float32
	// 胜利场数
	WinCount uint32
	// 平均胜率
	AvgWinRate float32
}

func NewProcessGameHandler(svcCtx *svc.ServiceContext) *ProcessGameHandler {
	taskInfo, err := svcCtx.DB.Task.Query().Where(task.PatternEQ(pattern.ProcessCoinGame)).First(context.Background())
	if err != nil || taskInfo == nil {
		return nil
	}

	var p payload.ProcessGamePayload
	if err := json.Unmarshal([]byte(taskInfo.Payload), &p); err != nil {
		p = payload.ProcessGamePayload{
			InvestTime: 60,
			OpenTime:   10,
		}
	}

	return &ProcessGameHandler{
		svcCtx:  svcCtx,
		taskId:  taskInfo.ID,
		payload: p,
		mode:    "coin",
	}
}

// ProcessTask if return err != nil , asynq will retry | 如果返回错误不为空则会重试
func (l *ProcessGameHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {

	if l == nil {
		return nil
	}
	if l.taskId == 0 {
		logx.Errorw("failed to load task info")
		return errorx.NewInternalError(i18n.DatabaseError)
	}

	currentRound, err := l.svcCtx.WolfLampRpc.FindRound(ctx, &wolflamp.FindRoundReq{Mode: l.mode})
	// 游戏数据不存在则创建新一轮游戏
	if err != nil {
		if status.Convert(err).Message() != "target does not exist" {
			return err
		}
		l.ClearCache(ctx)
		return l.ProcessNew(ctx, currentRound)
	}

	nowTime := time.Now().Unix()
	// 判断当前回合是否处于投注阶段
	if currentRound.Status == roundenum.Investing.Val() && currentRound.StartAt <= nowTime && currentRound.OpenAt > nowTime {
		return l.ProcessInvest(ctx, currentRound)
	}

	// 判断当前回合是否投注结束，只要过了投注时间且还没开奖，都优先开奖
	if currentRound.Status == roundenum.Investing.Val() && currentRound.OpenAt <= nowTime {
		return l.ProcessOpen(ctx, currentRound)
	}

	// 已开奖则可以创建新一轮游戏数据
	if currentRound.Status == roundenum.Opening.Val() {
		return l.ProcessNew(ctx, currentRound)
	}

	// 清除缓存，开始新的一轮游戏
	if currentRound.EndAt < nowTime {
		l.ClearCache(ctx)
		fmt.Printf("new round start[%s]\n", l.mode)
		fmt.Println("")
		return nil
	}
	return nil

}

// ClearCache 清除缓存
func (l *ProcessGameHandler) ClearCache(ctx context.Context) {
	l.svcCtx.Redis.Del(ctx, cachekey.CurrentGameRound.ModeVal(l.mode))
	l.svcCtx.Redis.Del(ctx, cachekey.CurrentGameLastRobotTime.ModeVal(l.mode))
	l.svcCtx.Redis.Del(ctx, cachekey.CurrentGameRobotNum.ModeVal(l.mode))
	l.svcCtx.Redis.Del(ctx, fmt.Sprintf("current_game:%s_opening_lock", l.mode))
	l.svcCtx.Redis.Del(ctx, cachekey.PreviousSelectedLambFoldNo.ModeVal(l.mode))
}

// ProcessInvest 执行投注逻辑
func (l *ProcessGameHandler) ProcessInvest(ctx context.Context, round *wolflamp.RoundInfo) error {

	mode := l.mode
	fmt.Printf("ProcessInvest[%s]:\n", mode)
	nowTime := time.Now().Unix()
	// 判断是否需要投入rob
	// 空闲X秒开始投放rob
	idleTime, err := l.svcCtx.WolfLampRpc.GetIdleTime(ctx, &wolflamp.Empty{})
	if err != nil || idleTime.IdleTime == 0 {
		fmt.Printf("ProcessInvest[%s]: empty idleTime, exit\n", mode)
		return nil
	}
	// 投放rob数量
	robNum, err := l.svcCtx.WolfLampRpc.GetRobotNum(ctx, &wolflamp.Empty{})
	if err != nil || robNum.Max == 0 {
		fmt.Printf("ProcessInvest[%s]: empty robNum, exit\n", mode)
		return nil
	}
	// 投放羊数量
	lampNum, err := l.svcCtx.WolfLampRpc.GetRobotLampNum(ctx, &wolflamp.Empty{})
	if err != nil || lampNum.Max == 0 {
		fmt.Printf("ProcessInvest[%s]: empty lampNum, exit\n", mode)
		return nil
	}
	// 判断投注人数有没有超过8人
	investNum, err := l.svcCtx.Redis.Get(ctx, cachekey.CurrentGameRobotNum.ModeVal(mode)).Uint64()
	if err != nil {
		investNum = 0
	}
	if investNum >= 8 {
		fmt.Printf("ProcessInvest[%s]: investNum >= 8, exit\n", mode)
		return nil
	}
	// 设置当上一轮结束，空闲【】秒，开始增加
	lastTime, err := l.svcCtx.Redis.Get(ctx, cachekey.CurrentGameLastRobotTime.ModeVal(mode)).Int64()
	if err != nil {
		lastTime = 0
	}
	if (nowTime < int64(idleTime.IdleTime)+round.StartAt) || (nowTime < lastTime+int64(idleTime.IdleTime)) {
		fmt.Printf("ProcessInvest[%s]: idleTime not reached, exit\n", mode)
		return nil
	}

	// 获取剩余机器人池数量
	robSumResp, err := l.svcCtx.WolfLampRpc.GetSum(ctx, &wolflamp.GetSumReq{Mode: mode, Status: 1, Type: uint32(poolenum.Robot)})
	if err != nil {
		return err
	}
	if robSumResp.Amount <= 0 {
		fmt.Printf("ProcessInvest[%s]: robot pool is not enough, exit\n", mode)
		return nil
	}

	// 生成rob数量在robNum.Min~robNum.Max之间
	rand.Seed(time.Now().UnixNano())
	robRand := rand.Intn(int(robNum.Max-robNum.Min+1)) + int(robNum.Min)
	prepare := make([]InvestPrepare, robRand)
	lampNum.Max = lampNum.Max / 100
	lampNum.Min = lampNum.Min / 100
	// 生成rob投羊数量
	for i := 0; i < robRand; i++ {
		// 生成羊数量在lampNum.Min~lampNum.Max之间
		// 羊数量要是100的整数倍
		rand.Seed(time.Now().UnixNano())
		lampRand := rand.Intn(int(lampNum.Max-lampNum.Min+1)) + int(lampNum.Min)
		prepareItem := InvestPrepare{
			LambNum:    uint32(lampRand) * 100,
			LambFoldNo: uint32(rand.Intn(8) + 1),
			PlayerId:   uint64(time.Now().UnixMilli()*10) + uint64(i),
			RoundId:    round.Id,
		}
		if float64(prepareItem.LambNum) > robSumResp.Amount {
			break
		}
		robSumResp.Amount -= float64(prepare[i].LambNum)
		prepare[i] = prepareItem
	}
	// 创建投注记录
	for _, v := range prepare {
		if v.PlayerId == 0 {
			continue
		}
		_, err := l.svcCtx.WolfLampRpc.Invest(ctx, &wolflamp.CreateInvestReq{
			PlayerId: v.PlayerId,
			RoundId:  v.RoundId,
			LambNum:  v.LambNum,
			FoldNo:   v.LambFoldNo,
			Mode:     mode,
		})
		if err != nil {
			return err
		}
		l.svcCtx.Redis.Incr(ctx, cachekey.CurrentGameRobotNum.ModeVal(mode))
		fmt.Printf("ProcessInvest[%s]: RoundId=%d, PlayerId=%d, LambNum=%d, FoldNo=%d \n", mode, v.RoundId, v.PlayerId, v.LambNum, v.LambFoldNo)
	}
	l.svcCtx.Redis.Set(ctx, cachekey.CurrentGameLastRobotTime.ModeVal(mode), time.Now().Unix(), 0)
	fmt.Println("")
	return nil

}

// ProcessOpen 执行开奖逻辑
func (l *ProcessGameHandler) ProcessOpen(ctx context.Context, round *wolflamp.RoundInfo) error {

	fmt.Printf("ProcessOpen[%s]\n", l.mode)

	openLock := fmt.Sprintf("current_game:%s_opening_lock", l.mode)
	result, err := l.svcCtx.Redis.SetNX(ctx, openLock, time.Now().Unix(), time.Minute*3).Result()
	if err != nil {
		return err
	}
	if !result {
		fmt.Printf("ProcessOpen[%s]: %s_opening_lock already exists, exit\n", l.mode, l.mode)
		return nil
	}

	// 被攻击的小羊按照其他羊圈用户投放的小羊数量占比进行分配
	_, err = l.svcCtx.WolfLampRpc.DealOpenGame(ctx, &wolflamp.DealOpenGameReq{Mode: l.mode})

	if err != nil {
		l.svcCtx.Redis.Del(ctx, openLock)
		return err
	}

	fmt.Println("")
	return nil

}

// ProcessNew 创建新一轮游戏
func (l *ProcessGameHandler) ProcessNew(ctx context.Context, round *wolflamp.RoundInfo) error {

	fmt.Printf("ProcessNew[%s]\n", l.mode)
	nowTime := time.Now().Unix()
	start := nowTime
	if round != nil && round.EndAt > nowTime {
		start = round.EndAt + 1
	}
	open := start + l.payload.InvestTime
	end := open + l.payload.OpenTime
	_, err := l.svcCtx.WolfLampRpc.CreateRound(ctx, &wolflamp.CreateRoundReq{
		StartAt: start, OpenAt: open, EndAt: end, Mode: l.mode,
	})
	if err != nil {
		return err
	}
	fmt.Println("")
	return nil

}

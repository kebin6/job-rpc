package game

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/kebin6/wolflamp-rpc/common/enum/cachekey"
	"github.com/kebin6/wolflamp-rpc/common/enum/roundenum"
	"github.com/kebin6/wolflamp-rpc/types/wolflamp"
	"github.com/suyuan32/simple-admin-common/i18n"
	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/suyuan32/simple-admin-job/ent/task"
	"github.com/suyuan32/simple-admin-job/internal/mqs/amq/types/pattern"
	"github.com/suyuan32/simple-admin-job/internal/mqs/amq/types/payload"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"math/rand"
	"time"

	"github.com/suyuan32/simple-admin-job/internal/svc"
)

type ProcessGameHandler struct {
	svcCtx  *svc.ServiceContext
	taskId  uint64
	payload payload.ProcessGamePayload
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
	taskInfo, err := svcCtx.DB.Task.Query().Where(task.PatternEQ(pattern.ProcessGame)).First(context.Background())
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

	currentRound, err := l.svcCtx.WolfLampRpc.FindRound(ctx, &wolflamp.FindRoundReq{})
	// 游戏数据不存在则创建新一轮游戏
	if err != nil {
		l.svcCtx.Redis.Del(ctx, cachekey.CurrentGameRound.Val())
		l.svcCtx.Redis.Del(ctx, cachekey.CurrentGameLastRobotTime.Val())
		l.svcCtx.Redis.Del(ctx, cachekey.CurrentGameRobotNum.Val())
		l.svcCtx.Redis.Del(ctx, "current_game:opening_lock")
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

	// 判断缓存中的当前回合是否变化
	if currentRound.EndAt < nowTime {
		l.svcCtx.Redis.Del(ctx, cachekey.CurrentGameRound.Val())
		l.svcCtx.Redis.Del(ctx, cachekey.CurrentGameLastRobotTime.Val())
		l.svcCtx.Redis.Del(ctx, cachekey.CurrentGameRobotNum.Val())
		l.svcCtx.Redis.Del(ctx, "current_game:opening_lock")
		fmt.Println("new round start")
		fmt.Println("")
		return nil
	}
	return nil

}

// ProcessInvest 执行投注逻辑
func (l *ProcessGameHandler) ProcessInvest(ctx context.Context, round *wolflamp.RoundInfo) error {

	fmt.Println("ProcessInvest")
	nowTime := time.Now().Unix()
	// 判断是否需要投入rob
	// 空闲X秒开始投放rob
	idleTime, err := l.svcCtx.WolfLampRpc.GetIdleTime(ctx, &wolflamp.Empty{})
	if err != nil || idleTime.IdleTime == 0 {
		fmt.Println("ProcessInvest: empty idleTime, exit")
		return nil
	}
	// 投放rob数量
	robNum, err := l.svcCtx.WolfLampRpc.GetRobotNum(ctx, &wolflamp.Empty{})
	if err != nil || robNum.Max == 0 {
		fmt.Println("ProcessInvest: empty robNum, exit")
		return nil
	}
	// 投放羊数量
	lampNum, err := l.svcCtx.WolfLampRpc.GetRobotLampNum(ctx, &wolflamp.Empty{})
	if err != nil || lampNum.Max == 0 {
		fmt.Println("ProcessInvest: empty lampNum, exit")
		return nil
	}
	// 判断投注人数有没有超过8人
	investNum, err := l.svcCtx.Redis.Get(ctx, cachekey.CurrentGameRobotNum.Val()).Uint64()
	if err != nil {
		investNum = 0
	}
	if investNum >= 8 {
		fmt.Println("ProcessInvest: investNum >= 8, exit")
		return nil
	}
	// 设置当上一轮结束，空闲【】秒，开始增加
	lastTime, err := l.svcCtx.Redis.Get(ctx, cachekey.CurrentGameLastRobotTime.Val()).Int64()
	if err != nil {
		lastTime = 0
	}
	if (nowTime < int64(idleTime.IdleTime)+round.StartAt) || (nowTime < lastTime+int64(idleTime.IdleTime)) {
		fmt.Println("ProcessInvest: idleTime not reached, exit")
		return nil
	}
	// 生成rob数量在robNum.Min~robNum.Max之间
	rand.Seed(time.Now().UnixNano())
	robRand := rand.Intn(int(robNum.Max-robNum.Min+1)) + int(robNum.Min)
	prepare := make([]InvestPrepare, robRand)
	// 生成rob投羊数量
	for i := 0; i < robRand; i++ {
		// 生成羊数量在lampNum.Min~lampNum.Max之间
		rand.Seed(time.Now().UnixNano())
		lampRand := rand.Intn(int(lampNum.Max-lampNum.Min+1)) + int(lampNum.Min)
		prepare[i].LambNum = uint32(lampRand)
		prepare[i].PlayerId = uint64(time.Now().UnixMilli()*10) + uint64(i)
		prepare[i].RoundId = round.Id
		prepare[i].LambFoldNo = uint32(rand.Intn(8) + 1)
	}
	// 创建投注记录
	for _, v := range prepare {
		_, err := l.svcCtx.WolfLampRpc.Invest(ctx, &wolflamp.CreateInvestReq{
			PlayerId: v.PlayerId,
			RoundId:  v.RoundId,
			LambNum:  v.LambNum,
			FoldNo:   v.LambFoldNo,
		})
		if err != nil {
			return err
		}
		l.svcCtx.Redis.Incr(ctx, cachekey.CurrentGameRobotNum.Val())
		fmt.Printf("ProcessInvest: RoundId=%d, PlayerId=%d, LambNum=%d, FoldNo=%d \n", v.RoundId, v.PlayerId, v.LambNum, v.LambFoldNo)
	}
	l.svcCtx.Redis.Set(ctx, cachekey.CurrentGameLastRobotTime.Val(), time.Now().Unix(), 0)
	fmt.Println("")
	return nil

}

// ProcessOpen 执行开奖逻辑
func (l *ProcessGameHandler) ProcessOpen(ctx context.Context, round *wolflamp.RoundInfo) error {

	fmt.Println("ProcessOpen")

	result, err := l.svcCtx.Redis.SetNX(ctx, "current_game:opening_lock", time.Now().Unix(), 0).Result()
	if err != nil {
		return err
	}
	if !result {
		fmt.Println("ProcessOpen: opening_lock already exists, exit")
		return nil
	}

	// 抽选羊圈
	foldChoice, err := l.ChooseLambFold(ctx, round)
	if err != nil {
		l.svcCtx.Redis.Del(ctx, "current_game:opening_lock")
		return nil
	}
	fmt.Printf("ProcessOpen: ChooseLambFold: %d", *foldChoice)

	// 被攻击的小羊按照其他羊圈用户投放的小羊数量占比进行分配
	_, err = l.svcCtx.WolfLampRpc.DealOpenGame(ctx, &wolflamp.DealOpenGameReq{LambFoldNo: *foldChoice})

	if err != nil {
		l.svcCtx.Redis.Del(ctx, "current_game:opening_lock")
		return err
	}

	fmt.Println("")
	return nil

}

// ProcessNew 创建新一轮游戏
func (l *ProcessGameHandler) ProcessNew(ctx context.Context, round *wolflamp.RoundInfo) error {

	fmt.Println("ProcessNew")
	nowTime := time.Now().Unix()
	start := nowTime
	if round != nil && round.EndAt > nowTime {
		start = round.EndAt + 1
	}
	open := start + l.payload.InvestTime
	end := open + l.payload.OpenTime
	_, err := l.svcCtx.WolfLampRpc.CreateRound(ctx, &wolflamp.CreateRoundReq{
		StartAt: start, OpenAt: open, EndAt: end,
	})
	if err != nil {
		return err
	}
	fmt.Println("")
	return nil

}

// ChooseLambFold 抽选羊圈
func (l *ProcessGameHandler) ChooseLambFold(ctx context.Context, round *wolflamp.RoundInfo) (*uint32, error) {

	fmt.Println("ChooseLambFold")

	rand.Seed(time.Now().UnixNano())
	// 随机选择10~30随机数，作为查询值
	randNum := uint64(rand.Intn(20) + 10)
	// 随机数作为历史回合数回溯各个羊圈的历史总盈亏
	// 先获取当前回合的累计回合数
	findRound, err := l.svcCtx.WolfLampRpc.FindRound(ctx, &wolflamp.FindRoundReq{})
	if err != nil {
		return nil, err
	}
	if findRound.TotalRoundCount <= 1 {
		// 如果当前是历史第一回合，则不需要聚合统计数据
		rand.Seed(time.Now().UnixNano())
		foldRand := rand.Intn(8)
		choice := pointy.GetPointer(uint32(foldRand + 1))
		return choice, nil
	} else if findRound.TotalRoundCount <= randNum {
		// 如果历史回合数小于随机数，则以历史回合数聚合统计
		randNum = findRound.TotalRoundCount - 1
	}

	// 根据查询值获取羊圈统计数据
	maxRoundCount := findRound.TotalRoundCount - 1
	minRoundCount := maxRoundCount - randNum + 1
	aggregateResult, err := l.svcCtx.WolfLampRpc.GetLambFoldAggregate(ctx, &wolflamp.GetLambFoldAggregateReq{
		TotalRoundCountMin: minRoundCount,
		TotalRoundCountMax: maxRoundCount,
	})
	if err != nil {
		return nil, err
	}
	// 最小的数
	excludeFirstOne := -1
	// 第二小的数
	excludeSecondOne := -1
	// 盈亏数最小的2个羊圈优先排除
	for i, v := range aggregateResult.Data {
		if excludeFirstOne == -1 {
			excludeFirstOne = i
			excludeSecondOne = i
			continue
		}
		if v.ProfitAndLossCount < aggregateResult.Data[excludeFirstOne].ProfitAndLossCount {
			excludeSecondOne = excludeFirstOne
			excludeFirstOne = i
		} else if v.ProfitAndLossCount < aggregateResult.Data[excludeSecondOne].ProfitAndLossCount {
			excludeSecondOne = i
		}
	}
	// 排除胜率最低的羊圈
	excludeThirdOne := -1
	for i, v := range aggregateResult.Data {
		if i == excludeFirstOne || i == excludeSecondOne {
			continue
		}
		if excludeThirdOne == -1 {
			excludeThirdOne = i
			continue
		}
		if v.AvgWinRate < aggregateResult.Data[excludeThirdOne].AvgWinRate {
			excludeThirdOne = i
		}
	}
	// 在剩下5个羊圈中随机抽选
	rand.Seed(time.Now().UnixNano())
	foldRand := rand.Intn(5)
	count := 0
	for i := 0; i < 8; i++ {
		if i == excludeFirstOne || i == excludeSecondOne || i == excludeThirdOne {
			continue
		}
		if count == foldRand {
			choice := pointy.GetPointer(uint32(i + 1))
			return choice, nil
		}
		count++
	}
	return nil, nil

}
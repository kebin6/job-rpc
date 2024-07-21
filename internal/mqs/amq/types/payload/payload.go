// Package payload defines all the payload structures used in tasks
package payload

type HelloWorldPayload struct {
	Name string `json:"name"`
}

type ProcessGamePayload struct {
	InvestTime int64 `json:"invest_time"` // 投注时长/秒
	OpenTime   int64 `json:"open_time"`   // 开奖-结束时长/秒
}

type SyncGcicsPayload struct {
	ChunkSize int64 `json:"chunk"` // 每次处理批次
}

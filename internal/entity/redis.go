package entity

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisWrapper はRedis.Clientをラップします。
type RedisWrapper interface {
	Publish(ctx context.Context, channel string, message interface{}) *IntCmd
	Subscribe(ctx context.Context, channels ...string) *PubSub
}

type RedisConn struct {
	Conn *redis.Client
}

func (r *RedisConn) Subscribe(ctx context.Context, channels ...string) *PubSub {
	ps := r.Conn.Subscribe(ctx, channels...)
	return &PubSub{ps: ps}
}

func (r *RedisConn) Publish(ctx context.Context, channel string, message interface{}) *IntCmd {
	cmd := r.Conn.Publish(ctx, channel, message)

	// go-redis の IntCmd から結果を取得
	val, err := cmd.Result()

	// 新しい entity.IntCmd インスタンスを作成して結果を設定
	return &IntCmd{
		baseCmd: baseCmd{
			ctx:  ctx,
			err:  err,
			args: []interface{}{channel, message},
		},
		val: val,
	}
}

type PubSub struct {
	ps *redis.PubSub
}

func (p *PubSub) ReceiveMessage(ctx context.Context) (*redis.Message, error) {
	return p.ps.ReceiveMessage(ctx)
}

func (p *PubSub) Close() error {
	return p.ps.Close()
}

type IntCmd struct {
	baseCmd

	val int64
}

func (cmd *IntCmd) Err() error {
	return cmd.baseCmd.err // baseCmd 構造体に保持されているエラーを返す
}

type baseCmd struct {
	ctx    context.Context
	args   []interface{}
	err    error
	keyPos int8

	_readTimeout *time.Duration
}

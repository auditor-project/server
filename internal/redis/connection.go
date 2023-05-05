package redis

import (
	"fmt"

	"auditor.z9fr.xyz/server/internal/lib"
	"github.com/hibiken/asynq"
)

type RedisConnection struct {
	asynq.RedisClientOpt
}

func NewRedisConnection(env *lib.Env) RedisConnection {
	redisOpt := asynq.RedisClientOpt{
		Addr: fmt.Sprintf("%s:%s", env.REDIS_HOST, env.REDIS_PORT),
	}

	return RedisConnection{
		redisOpt,
	}
}

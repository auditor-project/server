package worker

import (
	"go.uber.org/fx"
)

const (
	TaskStartAnalyser = "task:start-analyse"
)

var Module = fx.Options(
	fx.Provide(NewRedisTaskDistributor),
	fx.Provide(NewRedisTaskProcessor),
)

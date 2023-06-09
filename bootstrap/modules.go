package bootstrap

import (
	"auditor.z9fr.xyz/server/internal/db"
	"auditor.z9fr.xyz/server/internal/handler"
	"auditor.z9fr.xyz/server/internal/lib"
	"auditor.z9fr.xyz/server/internal/redis"
	"auditor.z9fr.xyz/server/internal/service"
	"auditor.z9fr.xyz/server/internal/worker"
	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	lib.Module,
	db.Module,
	service.Module,
	handler.Module,
	redis.Module,
	worker.Module,
)

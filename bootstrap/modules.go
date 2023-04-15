package bootstrap

import (
	"auditor.z9fr.xyz/server/internal/handler"
	"auditor.z9fr.xyz/server/internal/lib"
	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	lib.Module,
	handler.Module,
)

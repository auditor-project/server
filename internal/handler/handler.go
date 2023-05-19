package handler

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewErrorHandler),
	fx.Provide(NewParserHandlerImpl),
)

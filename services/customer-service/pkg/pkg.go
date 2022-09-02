package pkg

import "go.uber.org/fx"

// Module pkg
var Module = fx.Options(
	fx.Provide(NewEchoServer),
)

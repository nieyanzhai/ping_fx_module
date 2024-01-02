package ping_fx_module

import "go.uber.org/fx"

var module = fx.Module(
	"ping",
	fx.Provide(NewPing),
)

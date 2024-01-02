package ping_fx_module

import "go.uber.org/fx"

var Module = fx.Module(
	"ping",
	fx.Provide(NewPing),
)

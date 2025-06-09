package viper

import (
	"github.com/hdget/common/types"
	"go.uber.org/fx"
)

const (
	providerName = "param-viper"
)

var Capability = types.Capability{
	Category: types.ProviderCategoryConfig,
	Name:     providerName,
	Module: fx.Module(
		providerName,
		fx.Provide(New),
	),
}

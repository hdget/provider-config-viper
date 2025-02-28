package viper

import (
	"github.com/hdget/common/types"
	"github.com/hdget/provider-config-viper/pkg"
	"go.uber.org/fx"
)

const (
	providerName = "config-viper"
)

var Capability = &types.Capability{
	Category: types.ProviderCategoryConfig,
	Name:     providerName,
	Module: fx.Module(
		providerName,
		fx.Provide(pkg.New),
	),
}

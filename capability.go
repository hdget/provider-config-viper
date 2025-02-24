package viper

import (
	"github.com/hdget/common/types"
	"github.com/hdget/provider-config-viper/pkg"
	"go.uber.org/fx"
)

var Capability = &types.Capability{
	Category: types.ProviderCategoryConfig,
	Module: fx.Module(
		string(types.ProviderNameConfigViper),
		fx.Provide(pkg.New),
	),
}

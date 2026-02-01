package service

import "ride-sharing/services/trip-service/pkg/types"

func DefaultPricingConfig() *types.PricingConfig {
	return &types.PricingConfig{
		PricePerUnitOfDistance: 1.5,
		PricingPerMinute:       0.25,
	}
}

package fxboot

import (
	"go.uber.org/fx"
)

type ProvidingID int

const (
	// ProvidingAppID - app id
	ProvidingAppID ProvidingID = iota

	// ProvidingIDFXTimeouts - fx timeouts
	ProvidingIDFXTimeouts

	// ProvidingIDConfig - app config
	ProvidingIDConfig

	// ProvidingIDLogger - logger
	ProvidingIDLogger

	// ProvidingIDFXLogger - fx logger
	ProvidingIDFXLogger

	// ProvidingIDDBClients - db clients
	ProvidingIDDBClients

	// ProvidingIDRedisClient - redis
	ProvidingRedisClient

	// ProvidingIDBackoff - backoff
	ProvidingIDBackoff

	// ProvidingGRPCServer - grpc
	ProvidingGRPCServer

	// ProvidingIDBudgetModule - budget
	ProvidingIDBudgetModule
)

type OptionsMap struct {
	Providing map[ProvidingID]fx.Option
	Invokes   []fx.Option
}

func OptionsMapToSlice(optionsMap OptionsMap) []fx.Option {
	result := make([]fx.Option, 0)

	for _, option := range optionsMap.Providing {
		result = append(result, option)
	}

	result = append(result, optionsMap.Invokes...)

	return result
}

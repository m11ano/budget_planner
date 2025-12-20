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

	// ProvidingHTTPFiberServer - fiber
	ProvidingHTTPFiberServer

	// ProvidingIDDeliveryHTTP
	ProvidingIDDeliveryHTTP

	// ProvidingIDDeliveryGRPC
	ProvidingIDGRPCAuthClient

	// ProvidingIDDeliveryGRPC
	ProvidingIDGRPCLedgerClient
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

package models

import (
	"rates_project/internal/domain/types"
)

type CalculationParams struct {
	Method types.CalcMethod
	N      int
	M      int
}

package dcy

import (
	"github.com/lafaulx/decay/model"
)

type IDcy interface {
	Damage(d *model.Dcy) *model.Dcy
}

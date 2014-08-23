package dcy

import (
	"github.com/lafaulx/decay/misc"
	"github.com/lafaulx/decay/model"
)

type Text struct{}

func (k *Text) Damage(t *model.Dcy) *model.Dcy {
	victim := t.Content

	if misc.FortuneWheel(t.CallCount) {
		for i := 0; i < misc.GetDataCountToDamage(len(victim)); i++ {
			victim = misc.ReplaceCharInString(victim, misc.RandomCharFromString(victim), misc.RandIntFromInterval(0, len(victim)-1))
		}

		t.Content = victim
	}

	return t
}

package core

import (
	"github.com/isfir/wowsims-turtle/sim/core/proto"
	"github.com/isfir/wowsims-turtle/sim/core/stats"
)

type APLValueCurrentAttackPower struct {
	DefaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueCurrentAttackPower(_ *proto.APLValueCurrentAttackPower) APLValue {
	return &APLValueCurrentAttackPower{
		unit: rot.unit,
	}
}
func (value *APLValueCurrentAttackPower) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueCurrentAttackPower) GetFloat(_ *Simulation) float64 {
	return value.unit.GetStat(stats.AttackPower)
}
func (value *APLValueCurrentAttackPower) String() string {
	return "Current Attack Power"
}

type APLValueCurrentSpellHaste struct {
	DefaultAPLValueImpl
	unit *Unit
}

func (rot *APLRotation) newValueCurrentSpellHaste(_ *proto.APLValueCurrentSpellHaste) APLValue {
	return &APLValueCurrentSpellHaste{
		unit: rot.unit,
	}
}
func (value *APLValueCurrentSpellHaste) Type() proto.APLValueType {
	return proto.APLValueType_ValueTypeFloat
}
func (value *APLValueCurrentSpellHaste) GetFloat(_ *Simulation) float64 {
	return value.unit.PseudoStats.CastSpeedMultiplier
}
func (value *APLValueCurrentSpellHaste) String() string {
	return "Current Spell Haste"
}

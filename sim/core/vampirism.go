package core

import "github.com/isfir/wowsims-turtle/sim/core/stats"

const VampirismHealSpellID int32 = 45419

func (unit *Unit) registerVampirismSpell() {
	unit.vampirismSpell = unit.GetOrRegisterSpell(SpellConfig{
		ActionID: ActionID{SpellID: VampirismHealSpellID},
		ProcMask: ProcMaskEmpty,
		Flags:    SpellFlagPassiveSpell | SpellFlagNoLogs,
	})
}

func (unit *Unit) applyVampirism(sim *Simulation, damage float64, target *Unit) {
	vampirism := unit.GetStat(stats.Vampirism)
	if vampirism <= 0 || damage <= 0 || !unit.IsOpponent(target) {
		return
	}

	healAmount := damage * vampirism / 100.0

	result := unit.vampirismSpell.NewResult(unit)
	result.Damage = healAmount
	result.Outcome = OutcomeHit
	unit.vampirismSpell.DealHealing(sim, result)
}

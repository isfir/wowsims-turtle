package mage

import (
	"time"

	"github.com/isfir/wowsims-turtle/sim/core"
)

const ArcaneSurgeRanks = 4

var ArcaneSurgeSpellId = [ArcaneSurgeRanks + 1]int32{0, 51933, 51934, 51935, 51936}
var ArcaneSurgeBaseDamage = [ArcaneSurgeRanks + 1][]float64{{0}, {202, 245}, {290, 350}, {398, 475}, {517, 613}}
var ArcaneSurgeSpellCoeff = [ArcaneSurgeRanks + 1]float64{0, .65, .65, .65, .65}
var ArcaneSurgeManaCost = [ArcaneSurgeRanks + 1]float64{0, 85, 110, 140, 170}
var ArcaneSurgeLevel = [ArcaneSurgeRanks + 1]int{0, 32, 40, 48, 56}

func (mage *Mage) registerArcaneSurgeSpell() {
	mage.RegisterAura(core.Aura{
		Label:    "Arcane Surge Trigger",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if result.DidResist() {
				mage.ArcaneSurgeAura.Activate(sim)
			}
		},
	})

	mage.ArcaneSurgeAura = mage.RegisterAura(core.Aura{
		Label:    "Arcane Surge",
		ActionID: core.ActionID{SpellID: ArcaneSurgeSpellId[ArcaneSurgeRanks]},
		Duration: time.Second * 4,
	})

	mage.ArcaneSurge = make([]*core.Spell, ArcaneSurgeRanks+1)

	for rank := 1; rank <= ArcaneSurgeRanks; rank++ {
		config := mage.newArcaneSurgeSpellConfig(rank)

		if config.RequiredLevel <= int(mage.Level) {
			mage.ArcaneSurge[rank] = mage.GetOrRegisterSpell(config)
		}
	}
}

func (mage *Mage) newArcaneSurgeSpellConfig(rank int) core.SpellConfig {
	spellId := ArcaneSurgeSpellId[rank]
	baseDamageLow := ArcaneSurgeBaseDamage[rank][0]
	baseDamageHigh := ArcaneSurgeBaseDamage[rank][1]
	spellCoeff := ArcaneSurgeSpellCoeff[rank]
	manaCost := ArcaneSurgeManaCost[rank]
	level := ArcaneSurgeLevel[rank]

	spellConfig := core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellId},
		SpellSchool: core.SpellSchoolArcane,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagAPL | core.SpellFlagBinary,

		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * 8,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return mage.ArcaneSurgeAura.IsActive()
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicCrit)
			spell.DealDamage(sim, result)
			mage.ArcaneSurgeAura.Deactivate(sim)
		},
	}

	return spellConfig
}

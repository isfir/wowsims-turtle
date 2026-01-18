package mage

import (
	"fmt"
	"time"

	"github.com/isfir/wowsims-turtle/sim/core"
)

const ArcaneMissilesRanks = 8

var ArcaneMissilesChannelSpellId = [ArcaneMissilesRanks + 1]int32{0, 5143, 5144, 5145, 8416, 8417, 10211, 10212, 25345}
var ArcaneMissilesTickSpellId = [ArcaneMissilesRanks + 1]int32{0, 7268, 7269, 7270, 8419, 8418, 10273, 10274, 25346}
var ArcaneMissilesBaseTickDamage = [ArcaneMissilesRanks + 1]float64{0, 26, 38, 58, 86, 118, 154, 196, 230}
var ArcaneMissilesSpellCoeff = [ArcaneMissilesRanks + 1]float64{0, .180, .279, .328, .328, .328, .328, .328, .328}
var ArcaneMissilesCastTime = [ArcaneMissilesRanks + 1]int32{0, 3, 4, 5, 5, 5, 5, 5, 5}
var ArcaneMissilesManaCost = [ArcaneMissilesRanks + 1]float64{0, 85, 140, 235, 320, 410, 500, 595, 655}
var ArcaneMissilesLevel = [ArcaneMissilesRanks + 1]int{0, 8, 16, 24, 32, 40, 48, 56, 56}

func (mage *Mage) registerArcaneMissilesSpell() {
	mage.ArcaneMissiles = make([]*core.Spell, ArcaneMissilesRanks+1)
	mage.ArcaneMissilesTickSpell = make([]*core.Spell, ArcaneMissilesRanks+1)

	for rank := 1; rank <= ArcaneMissilesRanks; rank++ {
		config := mage.getArcaneMissilesSpellConfig(rank)

		if config.RequiredLevel <= int(mage.Level) {
			mage.ArcaneMissiles[rank] = mage.GetOrRegisterSpell(config)
		}
	}
}

func (mage *Mage) getArcaneMissilesSpellConfig(rank int) core.SpellConfig {
	spellId := ArcaneMissilesChannelSpellId[rank]
	baseTickDamage := ArcaneMissilesBaseTickDamage[rank]
	castTime := ArcaneMissilesCastTime[rank]
	manaCost := ArcaneMissilesManaCost[rank]
	level := ArcaneMissilesLevel[rank]

	numTicks := castTime
	tickLength := time.Second

	tickSpell := mage.getArcaneMissilesTickSpell(rank)
	mage.ArcaneMissilesTickSpell[rank] = tickSpell

	return core.SpellConfig{
		SpellCode:   SpellCode_MageArcaneMissiles,
		ActionID:    core.ActionID{SpellID: spellId},
		SpellSchool: core.SpellSchoolArcane,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagAPL | core.SpellFlagChanneled | core.SpellFlagNoMetrics,

		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: fmt.Sprintf("ArcaneMissiles-%d-%d", +rank, numTicks),
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					// TODO: This check is necessary to ensure the final tick occurs before
					// Arcane Blast stacks are dropped. To fix this, ticks need to reliably
					// occur before aura expirations.

					//TODO: Test interaction in classic code without aura
					dot := mage.ArcaneMissiles[rank].Dot(aura.Unit)
					if dot.TickCount < dot.NumberOfTicks {
						dot.TickCount++
						dot.TickOnce(sim)
					}
				},
			},
			NumberOfTicks:       numTicks,
			TickLength:          tickLength,
			AffectedByCastSpeed: mage.Talents.AcceleratedArcana,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				tickSpell.Cast(sim, target)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},
		ExpectedTickDamage: func(sim *core.Simulation, target *core.Unit, spell *core.Spell, _ bool) *core.SpellResult {
			// Should be updated to take into account Arcane Rupture buff
			return tickSpell.CalcDamage(sim, target, baseTickDamage, spell.OutcomeExpectedMagicHitAndCrit)
		},
	}
}

func (mage *Mage) getArcaneMissilesTickSpell(rank int) *core.Spell {
	spellId := ArcaneMissilesTickSpellId[rank]
	baseTickDamage := ArcaneMissilesBaseTickDamage[rank]
	spellCoeff := ArcaneMissilesSpellCoeff[rank]

	return mage.RegisterSpell(core.SpellConfig{
		SpellCode:    SpellCode_MageArcaneMissilesTick,
		ActionID:     core.ActionID{SpellID: spellId}.WithTag(1),
		SpellSchool:  core.SpellSchoolArcane,
		DefenseType:  core.DefenseTypeMagic,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        SpellFlagMage,
		MissileSpeed: 20,

		Rank: 1,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcDamage(sim, target, baseTickDamage, spell.OutcomeMagicHitAndCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})
}

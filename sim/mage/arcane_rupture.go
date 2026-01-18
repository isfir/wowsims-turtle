package mage

import (
	"time"

	"github.com/isfir/wowsims-turtle/sim/core"
)

const ArcaneRuptureRanks = 6

var ArcaneRuptureSpellId = [ArcaneRuptureRanks + 1]int32{0, 51949, 51950, 51951, 51952, 51953, 51954}
var ArcaneRuptureBaseDamage = [ArcaneRuptureRanks + 1][]float64{{0}, {101, 115}, {171, 191}, {302, 334}, {375, 434}, {528, 577}, {703, 766}}
var ArcaneRuptureSpellCoeff = [ArcaneRuptureRanks + 1]float64{0, .9, .9, .9, .9, .9, .9}
var ArcaneRuptureManaCost = [ArcaneRuptureRanks + 1]float64{0, 80, 145, 210, 270, 320, 390}
var ArcaneRuptureLevel = [ArcaneRuptureRanks + 1]int{0, 21, 28, 36, 44, 52, 60}

func (mage *Mage) registerArcaneRuptureSpell() {
	if !mage.Talents.ArcaneRupture {
		return
	}

	mage.ArcaneRuptureAura = mage.RegisterAura(core.Aura{
		Label:    "Arcane Rupture",
		ActionID: core.ActionID{SpellID: 52502},
		Duration: time.Second * 8,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range mage.ArcaneMissilesTickSpell {
				if spell != nil {
					spell.DamageMultiplierAdditive += 0.20
				}
			}
			for _, spell := range mage.ArcaneMissiles {
				if spell != nil {
					spell.Cost.Multiplier += 20
				}
			}
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			for _, spell := range mage.ArcaneMissilesTickSpell {
				if spell != nil {
					spell.DamageMultiplierAdditive -= 0.20
				}
			}
			for _, spell := range mage.ArcaneMissiles {
				if spell != nil {
					spell.Cost.Multiplier -= 20
				}
			}
		},
	})

	mage.ArcaneRupture = make([]*core.Spell, ArcaneRuptureRanks+1)

	for rank := 1; rank <= ArcaneRuptureRanks; rank++ {
		config := mage.newArcaneRuptureSpellConfig(rank)

		if config.RequiredLevel <= int(mage.Level) {
			mage.ArcaneRupture[rank] = mage.GetOrRegisterSpell(config)
		}
	}
}

func (mage *Mage) newArcaneRuptureSpellConfig(rank int) core.SpellConfig {
	spellId := ArcaneRuptureSpellId[rank]
	baseDamageLow := ArcaneRuptureBaseDamage[rank][0]
	baseDamageHigh := ArcaneRuptureBaseDamage[rank][1]
	spellCoeff := ArcaneRuptureSpellCoeff[rank]
	manaCost := ArcaneRuptureManaCost[rank]
	level := ArcaneRuptureLevel[rank]
	castTime := time.Millisecond * 2500

	spellConfig := core.SpellConfig{
		ActionID:    core.ActionID{SpellID: spellId},
		SpellCode:   SpellCode_MageArcaneRupture,
		SpellSchool: core.SpellSchoolArcane,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,
		Flags:       SpellFlagMage | core.SpellFlagAPL,

		RequiredLevel: level,
		Rank:          rank,

		ManaCost: core.ManaCostOptions{
			FlatCost: manaCost,
		},
		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD:      core.GCDDefault,
				CastTime: castTime,
			},
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * 15,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: spellCoeff,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			baseDamage := sim.Roll(baseDamageLow, baseDamageHigh)
			result := spell.CalcDamage(sim, target, baseDamage, spell.OutcomeMagicHitAndCrit)

			if result.Landed() {
				mage.ArcaneRuptureAura.Activate(sim)
			}

			spell.DealDamage(sim, result)
		},
	}

	return spellConfig
}

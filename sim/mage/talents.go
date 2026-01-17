package mage

import (
	"slices"
	"time"

	"github.com/wowsims/classic/sim/core"
)

func (mage *Mage) ApplyTalents() {
	mage.applyArcaneTalents()
	mage.applyFireTalents()
	mage.applyFrostTalents()
}

func (mage *Mage) applyArcaneTalents() {
	mage.applyMagicAbsorption()
	mage.applyArcaneConcentration()
	mage.applyTemporalConvergence()
	mage.registerPresenceOfMindCD()
	mage.registerArcanePowerCD()

	// Arcane Subtlety
	if mage.Talents.ArcaneSubtlety > 0 {
		// Target's resistance part does not seem to be implemented
		threatMultiplier := 1 - .20*float64(mage.Talents.ArcaneSubtlety)
		mage.OnSpellRegistered(func(spell *core.Spell) {
			if spell.SpellSchool.Matches(core.SpellSchoolArcane) && spell.Flags.Matches(SpellFlagMage) {
				spell.ThreatMultiplier *= threatMultiplier
			}
		})
	}

	// Arcane Focus
	if mage.Talents.ArcaneFocus > 0 {
		bonusHit := 2 * float64(mage.Talents.ArcaneFocus) * core.SpellHitRatingPerHitChance
		mage.OnSpellRegistered(func(spell *core.Spell) {
			if spell.SpellSchool.Matches(core.SpellSchoolArcane) && spell.Flags.Matches(SpellFlagMage) {
				spell.BonusHitRating += bonusHit
			}
		})
	}

	// Arcane Impact
	if mage.Talents.ArcaneImpact > 0 {
		bonusCrit := 2 * float64(mage.Talents.ArcaneImpact) * core.SpellCritRatingPerCritChance
		mage.OnSpellRegistered(func(spell *core.Spell) {
			if spell.SpellSchool.Matches(core.SpellSchoolArcane) && spell.Flags.Matches(SpellFlagMage) {
				spell.BonusCritRating += bonusCrit
			}
		})
	}

	// Arcane Meditation
	// TODO: Implement turtle version properly
	mage.PseudoStats.SpiritRegenRateCasting += 0.05 * float64(mage.Talents.ArcaneMeditation)

	// Arcane Potency
	if mage.Talents.ArcanePotency > 0 {
		critBonus := .50 * float64(mage.Talents.ArcanePotency)
		mage.OnSpellRegistered(func(spell *core.Spell) {
			if spell.SpellSchool.Matches(core.SpellSchoolArcane) && spell.Flags.Matches(SpellFlagMage) {
				spell.CritDamageBonus += critBonus
			}
		})
	}
}

func (mage *Mage) applyFireTalents() {
	mage.applyIgnite()
	mage.applyImprovedScorch()
	mage.applyMasterOfElements()

	mage.registerCombustionCD()

	// Burning Soul
	if mage.Talents.BurningSoul > 0 {
		threatMultiplier := 1 - .15*float64(mage.Talents.BurningSoul)
		mage.OnSpellRegistered(func(spell *core.Spell) {
			if spell.SpellSchool.Matches(core.SpellSchoolFire) && spell.Flags.Matches(SpellFlagMage) {
				spell.ThreatMultiplier *= threatMultiplier
			}
		})
	}

	// Critical Mass
	if mage.Talents.CriticalMass > 0 {
		bonusCrit := 2 * float64(mage.Talents.CriticalMass) * core.SpellCritRatingPerCritChance
		mage.OnSpellRegistered(func(spell *core.Spell) {
			if spell.SpellSchool.Matches(core.SpellSchoolFire) && spell.Flags.Matches(SpellFlagMage) {
				spell.BonusCritRating += bonusCrit
			}
		})
	}

	// Fire Power
	if mage.Talents.FirePower > 0 {
		bonusDamageMultiplierAdditive := 0.02 * float64(mage.Talents.FirePower)
		mage.OnSpellRegistered(func(spell *core.Spell) {
			// Fire Power buffs pretty much all mage fire spells EXCEPT ignite
			if spell.SpellSchool.Matches(core.SpellSchoolFire) && spell.Flags.Matches(SpellFlagMage) && spell.SpellCode != SpellCode_MageIgnite {
				spell.DamageMultiplierAdditive += bonusDamageMultiplierAdditive
			}
		})
	}
}

func (mage *Mage) applyFrostTalents() {
	mage.registerColdSnapCD()
	mage.registerIceBarrierSpell()
	mage.applyWintersChill()

	// Elemental Precision
	if mage.Talents.ElementalPrecision > 0 {
		bonusHit := 2 * float64(mage.Talents.ElementalPrecision) * core.SpellHitRatingPerHitChance

		mage.OnSpellRegistered(func(spell *core.Spell) {
			if spell.Flags.Matches(SpellFlagMage) && (spell.SpellSchool.Matches(core.SpellSchoolFire) || spell.SpellSchool.Matches(core.SpellSchoolFrost)) {
				spell.BonusHitRating += bonusHit
			}
		})
	}

	// Ice Shards
	if mage.Talents.IceShards > 0 {
		critBonus := .20 * float64(mage.Talents.IceShards)

		mage.OnSpellRegistered(func(spell *core.Spell) {
			if spell.SpellSchool.Matches(core.SpellSchoolFrost) && spell.Flags.Matches(SpellFlagMage) {
				spell.CritDamageBonus += critBonus
			}
		})
	}

	// Piercing Ice
	if mage.Talents.PiercingIce > 0 {
		bonusDamageMultiplierAdditive := 0.02 * float64(mage.Talents.PiercingIce)

		mage.OnSpellRegistered(func(spell *core.Spell) {
			if spell.SpellSchool.Matches(core.SpellSchoolFrost) && spell.Flags.Matches(SpellFlagMage) {
				spell.DamageMultiplierAdditive += bonusDamageMultiplierAdditive
			}
		})
	}

	// Frost Channeling
	if mage.Talents.FrostChanneling > 0 {
		manaCostMultiplier := 5 * mage.Talents.FrostChanneling
		threatMultiplier := 1 - .10*float64(mage.Talents.FrostChanneling)
		mage.OnSpellRegistered(func(spell *core.Spell) {
			if spell.SpellSchool.Matches(core.SpellSchoolFrost) && spell.Flags.Matches(SpellFlagMage) {
				spell.Cost.Multiplier -= manaCostMultiplier
				spell.ThreatMultiplier *= threatMultiplier
			}
		})
	}
}

func (mage *Mage) applyMagicAbsorption() {
	if mage.Talents.MagicAbsorption == 0 {
		return
	}

	spellID := []int32{29441, 29444, 29445}[mage.Talents.MagicAbsorption-1]
	magicAbsorptionBonus := []float64{4, 7, 10}[mage.Talents.MagicAbsorption-1]
	manaMetrics := mage.NewManaMetrics(core.ActionID{SpellID: 29442})

	mage.AddResistances(magicAbsorptionBonus)

	mage.RegisterAura(core.Aura{
		Label:    "Magic Absorption",
		ActionID: core.ActionID{SpellID: spellID},
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !spell.Flags.Matches(SpellFlagMage) {
				return
			}

			if result.DidResist() {
				mage.AddMana(sim, mage.MaxMana()*float64(mage.Talents.MagicAbsorption)/100, manaMetrics)
			}
		},
	})
}

func (mage *Mage) applyArcaneConcentration() {
	if mage.Talents.ArcaneConcentration == 0 {
		return
	}

	procChance := 0.02 * float64(mage.Talents.ArcaneConcentration)

	mage.ClearcastingAura = mage.RegisterAura(core.Aura{
		Label:    "Clearcasting",
		ActionID: core.ActionID{SpellID: 12536},
		Duration: time.Second * 15,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolCostMultiplier.AddToMagicSchools(-100)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			aura.Unit.PseudoStats.SchoolCostMultiplier.AddToMagicSchools(100)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// OnCastComplete is called after OnSpellHitDealt / etc, so don't deactivate if it was just activated.
			if aura.RemainingDuration(sim) == aura.Duration {
				return
			}
			if !spell.Flags.Matches(SpellFlagMage) {
				return
			}
			if spell.Cost != nil && spell.Cost.GetCurrentCost() == 0 {
				return
			}
			aura.Deactivate(sim)
		},
	})

	mage.RegisterAura(core.Aura{
		Label:    "Arcane Concentration",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.Flags.Matches(SpellFlagMage) {
				return
			}

			if sim.Proc(procChance, "Clearcasting") {
				mage.ClearcastingAura.Activate(sim)
			}
		},
	})
}

func (mage *Mage) registerPresenceOfMindCD() {
	if !mage.Talents.PresenceOfMind {
		return
	}

	actionID := core.ActionID{SpellID: 12043}
	cooldown := time.Second * 180

	affectedSpells := []*core.Spell{}
	pomAura := mage.RegisterAura(core.Aura{
		Label:    "Presence of Mind",
		ActionID: actionID,
		Duration: time.Second * 15,
		OnInit: func(aura *core.Aura, sim *core.Simulation) {
			for spellIdx := range mage.Spellbook {
				if spell := mage.Spellbook[spellIdx]; spell.DefaultCast.CastTime > 0 {
					affectedSpells = append(affectedSpells, spell)
				}
			}
		},
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(affectedSpells, func(spell *core.Spell) {
				spell.CastTimeMultiplier -= 1
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			core.Each(affectedSpells, func(spell *core.Spell) {
				spell.CastTimeMultiplier += 1
			})
			mage.PresenceOfMind.CD.Use(sim)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !slices.Contains(affectedSpells, spell) {
				return
			}

			aura.Deactivate(sim)
		},
	})

	mage.PresenceOfMind = mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: cooldown,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return mage.GCD.IsReady(sim)
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			pomAura.Activate(sim)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: mage.PresenceOfMind,
		Type:  core.CooldownTypeDPS,
	})
}

func (mage *Mage) registerArcanePowerCD() {
	if !mage.Talents.ArcanePower {
		return
	}

	var apPa *core.PendingAction
	actionID := core.ActionID{SpellID: 12042}
	manaMetricsPeriodic := mage.NewManaMetrics(core.ActionID{SpellID: 12042})
	manaMetricsDeath := mage.NewManaMetrics(core.ActionID{SpellID: 51941})

	mage.ArcanePowerAura = mage.RegisterAura(core.Aura{
		Label:    "Arcane Power",
		ActionID: actionID,
		Duration: time.Second * 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			mage.MultiplyCastSpeed(1.30)
			mage.PseudoStats.ManaGainMultiplier *= 0.5
			apPa = core.NewPeriodicAction(sim, core.PeriodicActionOptions{
				Period: time.Second,
				OnAction: func(s *core.Simulation) {
					mage.SpendMana(sim, mage.MaxMana()/100, manaMetricsPeriodic)
					if mage.CurrentManaPercent() < 0.1 {
						// Simulate death
						mage.RemoveHealth(sim, 999999)
						mage.SpendMana(sim, 999999, manaMetricsDeath)
						mage.PseudoStats.ManaGainMultiplier = 0
					}
				},
			})
			sim.AddPendingAction(apPa)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			mage.MultiplyCastSpeed(1 / 1.30)
			mage.PseudoStats.ManaGainMultiplier /= 0.5
			apPa.Cancel(sim)
		},
	})
	core.RegisterPercentDamageModifierEffect(mage.ArcanePowerAura, 1.3)

	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Second * 180,
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			mage.ArcanePowerAura.Activate(sim)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}

func (mage *Mage) applyImprovedScorch() {
	if mage.Talents.FireVulnerability == 0 {
		return
	}

	mage.ImprovedScorchAuras = mage.NewEnemyAuraArray(func(unit *core.Unit) *core.Aura {
		return core.ImprovedScorchAura(unit)
	})
}

func (mage *Mage) applyMasterOfElements() {
	if mage.Talents.MasterOfElements == 0 {
		return
	}

	refundCoeff := 0.1 * float64(mage.Talents.MasterOfElements)
	manaMetrics := mage.NewManaMetrics(core.ActionID{SpellID: 29076})

	mage.RegisterAura(core.Aura{
		Label:    "Master of Elements",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if spell.ProcMask.Matches(core.ProcMaskMeleeOrRanged) {
				return
			}
			if spell.CurCast.Cost == 0 {
				return
			}
			if result.DidCrit() {
				mage.AddMana(sim, spell.Cost.BaseCost*refundCoeff, manaMetrics)
			}
		},
	})
}

func (mage *Mage) registerCombustionCD() {
	if !mage.Talents.Combustion {
		return
	}

	actionID := core.ActionID{SpellID: 11129}
	cd := core.Cooldown{
		Timer:    mage.NewTimer(),
		Duration: time.Minute * 3,
	}

	var fireSpells []*core.Spell
	mage.OnSpellRegistered(func(spell *core.Spell) {
		if spell.SpellSchool.Matches(core.SpellSchoolFire) && spell.Flags.Matches(SpellFlagMage) {
			fireSpells = append(fireSpells, spell)
		}
	})

	numCrits := 0
	critPerStack := 10.0 * core.SpellCritRatingPerCritChance

	mage.CombustionAura = mage.RegisterAura(core.Aura{
		Label:     "Combustion",
		ActionID:  actionID,
		Duration:  core.NeverExpires,
		MaxStacks: 20,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			numCrits = 0
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			cd.Use(sim)
			mage.UpdateMajorCooldowns()
		},
		OnStacksChange: func(aura *core.Aura, sim *core.Simulation, oldStacks int32, newStacks int32) {
			bonusCrit := critPerStack * float64(newStacks-oldStacks)
			for _, spell := range fireSpells {
				spell.BonusCritRating += bonusCrit
			}
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || numCrits >= 3 || !spell.SpellSchool.Matches(core.SpellSchoolFire) || !spell.Flags.Matches(SpellFlagMage) {
				return
			}

			// Ignite, Living Bomb explosions, and Fire Blast with Overheart don't consume crit stacks
			// To Do: Classic - I don't believe ignite can crit so can probably remove this check?
			if spell.SpellCode == SpellCode_MageIgnite {
				return
			}

			// TODO: This wont work properly with flamestrike
			aura.AddStack(sim)

			if result.DidCrit() {
				numCrits++
				if numCrits == 3 {
					aura.Deactivate(sim)
				}
			}
		},
	})

	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID: actionID,
		Flags:    core.SpellFlagNoOnCastComplete,
		Cast: core.CastConfig{
			CD: cd,
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !mage.CombustionAura.IsActive()
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			mage.CombustionAura.Activate(sim)
			mage.CombustionAura.AddStack(sim)
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}

func (mage *Mage) registerColdSnapCD() {
	if !mage.Talents.ColdSnap {
		return
	}

	// Grab all frost spells with a CD > 0
	var affectedSpells = []*core.Spell{}
	mage.OnSpellRegistered(func(spell *core.Spell) {
		if spell.SpellSchool.Matches(core.SpellSchoolFrost) && spell.CD.Duration > 0 {
			affectedSpells = append(affectedSpells, spell)
		}
	})

	spell := mage.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 12472},
		Flags:    core.SpellFlagNoOnCastComplete,

		Cast: core.CastConfig{
			CD: core.Cooldown{
				Timer:    mage.NewTimer(),
				Duration: time.Duration(time.Minute * 10),
			},
		},
		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			for _, spell := range affectedSpells {
				spell.CD.Reset()
			}
		},
	})

	mage.AddMajorCooldown(core.MajorCooldown{
		Spell: spell,
		Type:  core.CooldownTypeDPS,
	})
}

func (mage *Mage) applyWintersChill() {
	if mage.Talents.WintersChill == 0 {
		return
	}

	procChance := float64(mage.Talents.WintersChill) * 0.2

	wcAuras := mage.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return core.WintersChillAura(target)
	})
	mage.Env.RegisterPreFinalizeEffect(func() {
		for _, spell := range mage.GetSpellsMatchingSchool(core.SpellSchoolFrost) {
			spell.RelatedAuras = append(spell.RelatedAuras, wcAuras)
		}
	})

	mage.RegisterAura(core.Aura{
		Label:    "Winters Chill Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || !spell.SpellSchool.Matches(core.SpellSchoolFrost) {
				return
			}

			if sim.Proc(procChance, "Winters Chill") {
				aura := wcAuras.Get(result.Target)
				aura.Activate(sim)
				if aura.IsActive() {
					aura.AddStack(sim)
				}
			}
		},
	})
}

func (mage *Mage) applyTemporalConvergence() {
	if mage.Talents.TemporalConvergence == 0 {
		return
	}

	procChance := 0.05 * float64(mage.Talents.TemporalConvergence)
	manaMetrics := mage.NewManaMetrics(core.ActionID{SpellID: 51962})
	icdDuration := time.Second * 15

	icd := core.Cooldown{
		Timer:    mage.NewTimer(),
		Duration: icdDuration,
	}

	// Should maybe only be local
	mage.TemporalConvergenceAura = mage.RegisterAura(core.Aura{
		Label:    "Temporal Convergence",
		ActionID: core.ActionID{SpellID: 51961},
		Duration: time.Second * 15,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			// OnCastComplete is called after OnSpellHitDealt / etc, so don't deactivate if it was just activated.
			if aura.RemainingDuration(sim) == aura.Duration {
				return
			}

			if spell.SpellCode != SpellCode_MageArcaneRupture {
				return
			}

			mage.AddMana(sim, spell.Cost.BaseCost, manaMetrics)
			aura.Deactivate(sim)
		},
	})

	mage.RegisterAura(core.Aura{
		Label:    "Temporal Convergence Talent",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnSpellHitDealt: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
			if !result.Landed() || spell.SpellCode != SpellCode_MageArcaneMissilesTick {
				return
			}

			if !icd.IsReady(sim) {
				return
			}

			if sim.Proc(procChance, "Temporal Convergence") {
				icd.Use(sim)

				for _, spell := range mage.ArcaneRupture {
					if spell != nil {
						spell.CD.Reset()
					}
				}
				mage.TemporalConvergenceAura.Activate(sim)
			}
		},
	})
}

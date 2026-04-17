package item_effects

import (
	"time"

	"github.com/isfir/wowsims-turtle/sim/common/guardians"
	"github.com/isfir/wowsims-turtle/sim/core"
	"github.com/isfir/wowsims-turtle/sim/core/stats"
)

const (
	SpellpowerGogglesXtremePlusPlus = 33095
	RingOfBurningTalons             = 33154
	DropletOfNordrassil             = 33294
	SpellwovenNobilityDrape         = 55056
	JewelOfWildMagics               = 55087
	RemainsOfOverwhelmingPower      = 55093
	BindingsOfContainedMagic        = 55106
	SphereOfTheEndlessGulch         = 55501
	TheScytheOfElune                = 55505
	TrueBandOfSulfuras              = 58088
	SigilOfAncientAccord            = 58244
)

func init() {
	core.AddEffectsToTest = false

	///////////////////////////////////////////////////////////////////////////
	//                                 Trinkets
	///////////////////////////////////////////////////////////////////////////

	// https://database.turtlecraft.gg/?item=33294
	// Increases damage done by magical spells and effects by up to 80 and increases spell hit chance by 3% for 10 sec whenever your spells are fully or partially resisted.
	core.NewItemEffect(DropletOfNordrassil, func(agent core.Agent) {
		character := agent.GetCharacter()

		buffAura := character.NewTemporaryStatsAura(
			"Nordrassil's Reprieve",
			core.ActionID{SpellID: 58227},
			stats.Stats{
				stats.SpellDamage: 80,
				stats.SpellHit:    3,
			},
			time.Second*10,
		)

		core.MakeItemProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Nordrassil's Reprieve Passive",
			Callback: core.CallbackOnSpellHitDealt,
			Outcome:  core.OutcomeResisted,
			ProcMask: core.ProcMaskSpellDamage,
			ICD:      time.Second * 4,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				buffAura.Activate(sim)
			},
		})
	})

	// https://database.turtlecraft.gg/?item=55087
	// Unleash a blast of a random element, dealing 490 - 540 damage of that element within 10 yards, additionally causing an effect depending on the element.
	core.NewItemEffect(JewelOfWildMagics, func(agent core.Agent) {
		character := agent.GetCharacter()

		frostDebuffAuras := character.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
			aura := target.GetOrRegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 51005},
				Label:    "Emanating Frost (Jewel)",
				Duration: time.Second * 10,
				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.AddMoveSpeedModifier(&aura.ActionID, 0.50)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.RemoveMoveSpeedModifier(&aura.ActionID)
				},
			})
			core.AtkSpeedReductionEffect(aura, 1.15)
			return aura
		})

		frostExplosion := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 51004},
			SpellSchool: core.SpellSchoolFrost,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskSpellDamage,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					result := spell.CalcAndDealDamage(sim, aoeTarget, sim.Roll(491, 541), spell.OutcomeMagicHitAndCrit)
					if result.Landed() {
						frostDebuffAuras.Get(aoeTarget).Activate(sim)
					}
				}
			},
		})

		fireExplosion := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 51006},
			SpellSchool: core.SpellSchoolFire,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskSpellDamage,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			Dot: core.DotConfig{
				Aura: core.Aura{
					ActionID: core.ActionID{SpellID: 51007},
					Label:    "Scorched Earth (Jewel)",
				},
				NumberOfTicks: 3,
				TickLength:    time.Second * 2,
				OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
					dot.Snapshot(target, 110, isRollover)
				},
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					result := spell.CalcAndDealDamage(sim, aoeTarget, sim.Roll(491, 541), spell.OutcomeMagicHitAndCrit)
					if result.Landed() {
						spell.Dot(aoeTarget).Apply(sim)
					}
				}
			},
		})

		arcaneSurgeAura := character.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 51009},
			Label:    "Arcane Surge (Jewel)",
			Duration: time.Second * 12,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatDynamic(sim, stats.SpellDamage, 50)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatDynamic(sim, stats.SpellDamage, -50)
			},
		}).AttachMultiplyCastSpeed(&character.Unit, 1.03)

		arcaneExplosion := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 51008},
			SpellSchool: core.SpellSchoolArcane,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskSpellDamage,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					spell.CalcAndDealDamage(sim, aoeTarget, sim.Roll(491, 541), spell.OutcomeMagicHitAndCrit)
				}
				arcaneSurgeAura.Activate(sim)
			},
		})

		renewingBlast := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 51011},
			SpellSchool: core.SpellSchoolHoly,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagHelpful,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for _, partyAgent := range character.Party.PlayersAndPets {
					spell.CalcAndDealHealing(sim, &partyAgent.GetCharacter().Unit, sim.Roll(400, 450), spell.OutcomeHealingCrit)
				}
			},
		})

		holyExplosion := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 51010},
			SpellSchool: core.SpellSchoolHoly,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskSpellDamage,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					spell.CalcAndDealDamage(sim, aoeTarget, sim.Roll(491, 541), spell.OutcomeMagicHitAndCrit)
				}
				renewingBlast.Cast(sim, &character.Unit)
			},
		})

		effects := []*core.Spell{frostExplosion, fireExplosion, arcaneExplosion, holyExplosion}

		cdSpell := character.RegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{ItemID: JewelOfWildMagics},
			Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 5,
				},
				SharedCD: core.Cooldown{
					Timer:    character.GetOffensiveTrinketCD(),
					Duration: time.Second * 20,
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				effectIdx := int(sim.RandomFloat("Jewel of Wild Magics") * float64(len(effects)))
				effects[effectIdx].Cast(sim, target)
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Spell: cdSpell,
			Type:  core.CooldownTypeDPS,
		})
	})

	// https://database.turtlecraft.gg/?item=55093
	// Use: Summons a Minor Arcane Elemental that will protect you for 60 sec. While this elemental is alive it increases your damage and healing done by magical spells and effects by up to 55.
	core.NewItemEffect(RemainsOfOverwhelmingPower, func(agent core.Agent) {
		character := agent.GetCharacter()
		actionID := core.ActionID{SpellID: 51143}

		spell := character.RegisterSpell(core.SpellConfig{
			ActionID: actionID,
			Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 5,
				},
				SharedCD: core.Cooldown{
					Timer:    character.GetOffensiveTrinketCD(),
					Duration: time.Second * 20,
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for _, petAgent := range character.PetAgents {
					if elemental, ok := petAgent.(*guardians.MinorArcaneElemental); ok {
						elemental.EnableWithTimeout(sim, elemental, time.Minute)
						break
					}
				}
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeDPS,
			Spell: spell,
		})
	})

	// https://database.turtlecraft.gg/?item=58244
	// Your direct damaging spells have a 8% chance to detonate the surrounding magic, dealing 100 Arcane damage to the target and all targets within 0. The main target takes 300 additional Arcane damage.
	core.NewItemEffect(SigilOfAncientAccord, func(agent core.Agent) {
		character := agent.GetCharacter()

		procSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 44095},
			SpellSchool: core.SpellSchoolArcane,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.BonusCoefficient = 0.15
				spell.CalcAndDealDamage(sim, target, 300, spell.OutcomeMagicHitAndCrit)

				spell.BonusCoefficient = 0.05
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					spell.CalcAndDealDamage(sim, aoeTarget, 100, spell.OutcomeMagicHitAndCrit)
				}
			},
		})

		core.MakeItemProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Ancient Accord Passive",
			Callback:   core.CallbackOnSpellHitDealt,
			Outcome:    core.OutcomeLanded,
			ProcMask:   core.ProcMaskSpellDamage,
			ProcChance: 0.08,
			ICD:        time.Second * 2,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procSpell.Cast(sim, result.Target)
			},
		})
	})

	// https://database.turtlecraft.gg/?item=55501
	// Gives a 20% chance when your harmful spells land to grant you Wisdom of the Mak'aru. Upon gaining Wisdom of the Mak'aru 20 times, you enter an enlightened state, increasing casting speed by 20% for 12 sec.
	// (Proc chance: 20%, ICD: 3s)
	core.NewItemEffect(SphereOfTheEndlessGulch, func(agent core.Agent) {
		character := agent.GetCharacter()

		enlightenedAura := character.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 51270},
			Label:    "Enlightened State",
			Duration: time.Second * 12,
		}).AttachMultiplyCastSpeed(&character.Unit, 1.20)

		wisdomAura := character.RegisterAura(core.Aura{
			ActionID:  core.ActionID{SpellID: 51271},
			Label:     "Wisdom of the Mak'aru",
			Duration:  time.Second * 30,
			MaxStacks: 8,
		})

		core.MakeItemProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Wisdom of the Mak'aru Passive",
			Callback:   core.CallbackOnSpellHitDealt,
			Outcome:    core.OutcomeLanded,
			ProcMask:   core.ProcMaskSpellDamage,
			ProcChance: 0.20,
			ICD:        time.Second * 3,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				wisdomAura.Activate(sim)
				wisdomAura.AddStack(sim)

				if wisdomAura.GetStacks() >= wisdomAura.MaxStacks {
					wisdomAura.Deactivate(sim)
					enlightenedAura.Activate(sim)
				}
			},
		})
	})

	// https://database.turtlecraft.gg/?item=55505
	// Equip: Elune infuses your magic with astral power, giving a 5% chance when your harmful spells are cast to unleash astral energy, dealing 375 to 501 Arcane damage and imbuing your next harmful spell with moonlight. Spells imbued by moonlight increase the damage the target takes from its spell school by 8% and trigger an additional effect for 10 sec depending on the spell school of the harmful spell.
	// Use: Pierce the veil between Azeroth and Vorgendor, summoning a Scytheclaw Pureborn to protect you for 60 sec.
	core.NewItemEffect(TheScytheOfElune, func(agent core.Agent) {
		character := agent.GetCharacter()

		eluneInfusionICD := core.Cooldown{
			Timer:    character.NewTimer(),
			Duration: time.Second * 15,
		}
		eluneInfusion := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 57658},
			SpellSchool: core.SpellSchoolArcane,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BonusCoefficient: 0.15,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, sim.Roll(375, 501), spell.OutcomeMagicHitAndCrit)
			},
		})

		elunesMoonlightAura := character.GetOrRegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 57663},
			Label:    "Elune's Wrath (Moonlight)",
			Duration: time.Second * 10,
		})
		core.ApplyProcTriggerCallback(&character.Unit, elunesMoonlightAura, core.ProcTrigger{
			Callback:          core.CallbackOnSpellHitDealt,
			Outcome:           core.OutcomeLanded,
			ProcMask:          core.ProcMaskSpellDamage,
			SpellFlagsExclude: core.SpellFlagHelpful | core.SpellFlagPassiveSpell,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				eluneInfusionICD.Use(sim)
				elunesMoonlightAura.Deactivate(sim)
				eluneInfusion.Cast(sim, result.Target)
			},
		})

		applyElunesMoonlight := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 57663},
			SpellSchool: core.SpellSchoolArcane,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				elunesMoonlightAura.Activate(sim)
			},
		})

		elunesWrathAuras := character.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
			return target.GetOrRegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 52408},
				Label:    "Elune's Wrath",
				Duration: time.Second * 10,

				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] *= 1.08
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexArcane] /= 1.08
				},
			})
		})

		applyElunesWrath := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 52408},
			SpellSchool: core.SpellSchoolArcane,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagBinary | core.SpellFlagNoOnCastComplete | core.SpellFlagNoOnDamageDealt | core.SpellFlagPassiveSpell,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
				if result.Landed() {
					elunesWrathAuras.Get(target).Activate(sim)
				}

				applyElunesMoonlight.Cast(sim, &character.Unit)
			},
		})

		elunesRageTriggerDamage := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 52375},
			SpellSchool: core.SpellSchoolFire,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 325, spell.OutcomeMagicHitAndCrit)
			},
		})

		magicCritTakenSchools := []stats.SchoolIndex{
			stats.SchoolIndexArcane,
			stats.SchoolIndexFire,
			stats.SchoolIndexFrost,
			stats.SchoolIndexHoly,
			stats.SchoolIndexNature,
			stats.SchoolIndexShadow,
		}
		elunesRageAuras := character.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
			icd := core.Cooldown{
				Timer:    target.NewTimer(),
				Duration: time.Second * 2,
			}

			return target.GetOrRegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 52404},
				Label:    "Elune's Rage",
				Duration: time.Second * 10,

				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] *= 1.08
					for _, school := range magicCritTakenSchools {
						aura.Unit.PseudoStats.SchoolCritTakenChance[school] += 0.02
					}
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFire] /= 1.08
					for _, school := range magicCritTakenSchools {
						aura.Unit.PseudoStats.SchoolCritTakenChance[school] -= 0.02
					}
				},
				OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !result.Landed() || !result.DidCrit() {
						return
					}
					if spell.Flags.Matches(core.SpellFlagHelpful) {
						return
					}
					if spell.DefenseType != core.DefenseTypeMagic {
						return
					}
					if !icd.IsReady(sim) {
						return
					}

					icd.Use(sim)
					elunesRageTriggerDamage.Cast(sim, aura.Unit)
				},
			})
		})

		applyElunesRage := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 52404},
			SpellSchool: core.SpellSchoolFire,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagBinary | core.SpellFlagNoOnCastComplete | core.SpellFlagNoOnDamageDealt | core.SpellFlagPassiveSpell,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
				if result.Landed() {
					elunesRageAuras.Get(target).Activate(sim)
				}
			},
		})

		elunesIreTriggerDamage := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 52411},
			SpellSchool: core.SpellSchoolFrost,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagNoOnDamageDealt,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 275, spell.OutcomeMagicHitAndCrit)
			},
		})

		elunesIreAuras := character.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
			icd := core.Cooldown{
				Timer:    target.NewTimer(),
				Duration: time.Second * 2,
			}

			aura := target.GetOrRegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 57662},
				Label:    "Elune's Ire",
				Duration: time.Second * 10,

				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] *= 1.08
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexFrost] /= 1.08
				},
				OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !result.Landed() {
						return
					}
					if spell.Flags.Matches(core.SpellFlagHelpful) {
						return
					}
					if spell.DefenseType != core.DefenseTypeMagic {
						return
					}
					if !icd.IsReady(sim) {
						return
					}

					icd.Use(sim)
					elunesIreTriggerDamage.Cast(sim, aura.Unit)
				},
			})
			core.AtkSpeedReductionEffect(aura, 1.25)
			return aura
		})

		applyElunesIre := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 52406},
			SpellSchool: core.SpellSchoolFrost,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagBinary | core.SpellFlagNoOnCastComplete | core.SpellFlagNoOnDamageDealt | core.SpellFlagPassiveSpell,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
				if result.Landed() {
					elunesIreAuras.Get(target).Activate(sim)
				}
			},
		})

		elunesRadianceAttackerHeal := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 52409},
			SpellSchool: core.SpellSchoolHoly,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagHelpful,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealHealing(sim, target, 130, spell.OutcomeHealingCrit)
			},
		})

		elunesRadianceHitPenalty := stats.Stats{
			stats.MeleeHit: -3 * core.MeleeHitRatingPerHitChance,
			stats.SpellHit: -3 * core.SpellHitRatingPerHitChance,
		}
		elunesRadianceAuras := character.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
			return target.GetOrRegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 57666},
				Label:    "Elune's Radiance",
				Duration: time.Second * 10,

				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexHoly] *= 1.08
					aura.Unit.AddStatsDynamic(sim, elunesRadianceHitPenalty)
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexHoly] /= 1.08
					aura.Unit.AddStatsDynamic(sim, elunesRadianceHitPenalty.Invert())
				},
				OnSpellHitTaken: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
					if !result.Landed() {
						return
					}
					if spell.Flags.Matches(core.SpellFlagHelpful) {
						return
					}

					elunesRadianceAttackerHeal.Cast(sim, spell.Unit)
				},
			})
		})

		applyElunesRadiance := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 57661},
			SpellSchool: core.SpellSchoolHoly,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagBinary | core.SpellFlagNoOnCastComplete | core.SpellFlagNoOnDamageDealt | core.SpellFlagPassiveSpell,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
				if result.Landed() {
					elunesRadianceAuras.Get(target).Activate(sim)
				}
			},
		})

		elunesGraceTriggerDamage := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 52410},
			SpellSchool: core.SpellSchoolNature,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell | core.SpellFlagNoOnDamageDealt,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.CalcAndDealDamage(sim, target, 200, spell.OutcomeMagicHitAndCrit)
			},
		})

		elunesGraceAuras := character.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
			icd := core.Cooldown{
				Timer:    target.NewTimer(),
				Duration: time.Second * 2,
			}

			return target.GetOrRegisterAura(core.Aura{
				ActionID: core.ActionID{SpellID: 57664},
				Label:    "Elune's Grace",
				Duration: time.Second * 10,

				OnGain: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] *= 1.08
					aura.Unit.PseudoStats.DamageDealtMultiplier *= 0.92
				},
				OnExpire: func(aura *core.Aura, sim *core.Simulation) {
					aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexNature] /= 1.08
					aura.Unit.PseudoStats.DamageDealtMultiplier /= 0.92
				},
				OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
					if spell.Flags.Matches(core.SpellFlagPassiveSpell) {
						return
					}
					if !icd.IsReady(sim) {
						return
					}

					icd.Use(sim)
					elunesGraceTriggerDamage.Cast(sim, aura.Unit)
				},
			})
		})

		applyElunesGrace := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 52405},
			SpellSchool: core.SpellSchoolNature,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagBinary | core.SpellFlagNoOnCastComplete | core.SpellFlagNoOnDamageDealt | core.SpellFlagPassiveSpell,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
				if result.Landed() {
					elunesGraceAuras.Get(target).Activate(sim)
				}
			},
		})

		elunesTwilightManaDrainMetrics := map[int32]*core.ResourceMetrics{}

		elunesTwilightDot := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 57665},
			SpellSchool: core.SpellSchoolShadow,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagPureDot | core.SpellFlagNoOnCastComplete | core.SpellFlagPassiveSpell,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			Dot: core.DotConfig{
				Aura: core.Aura{
					ActionID: core.ActionID{SpellID: 57665},
					Label:    "Elune's Twilight",

					OnGain: func(aura *core.Aura, sim *core.Simulation) {
						aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] *= 1.08
						aura.Unit.PseudoStats.HealingTakenMultiplier *= 0.75
					},
					OnExpire: func(aura *core.Aura, sim *core.Simulation) {
						aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexShadow] /= 1.08
						aura.Unit.PseudoStats.HealingTakenMultiplier /= 0.75
					},
				},
				NumberOfTicks: 10,
				TickLength:    time.Second,

				OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
					dot.Snapshot(target, 120, isRollover)
				},
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)

					if target.HasManaBar() {
						metrics := elunesTwilightManaDrainMetrics[target.UnitIndex]
						if metrics == nil {
							metrics = target.NewManaMetrics(core.ActionID{SpellID: 57665})
							elunesTwilightManaDrainMetrics[target.UnitIndex] = metrics
						}
						target.AddMana(sim, -100, metrics)
					}
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				spell.Dot(target).Apply(sim)
			},
		})

		applyElunesTwilight := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 52407},
			SpellSchool: core.SpellSchoolShadow,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagBinary | core.SpellFlagNoOnCastComplete | core.SpellFlagNoOnDamageDealt | core.SpellFlagPassiveSpell,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
				if result.Landed() {
					elunesTwilightDot.Dot(target).Apply(sim)
				}
			},
		})

		core.MakeItemProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:                 "Elune Infusion Passive",
			Callback:             core.CallbackOnSpellHitDealt,
			Outcome:              core.OutcomeLanded,
			ProcMask:             core.ProcMaskSpellDamage,
			ProcChance:           0.05,
			LandedExecutionIndex: 1,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				if !eluneInfusionICD.IsReady(sim) {
					return
				}

				switch {
				case spell.SpellSchool.Matches(core.SpellSchoolArcane):
					applyElunesWrath.Cast(sim, result.Target)
				case spell.SpellSchool.Matches(core.SpellSchoolFire):
					applyElunesRage.Cast(sim, result.Target)
				case spell.SpellSchool.Matches(core.SpellSchoolFrost):
					applyElunesIre.Cast(sim, result.Target)
				case spell.SpellSchool.Matches(core.SpellSchoolHoly):
					applyElunesRadiance.Cast(sim, result.Target)
				case spell.SpellSchool.Matches(core.SpellSchoolNature):
					applyElunesGrace.Cast(sim, result.Target)
				case spell.SpellSchool.Matches(core.SpellSchoolShadow):
					applyElunesTwilight.Cast(sim, result.Target)
				}

				eluneInfusion.Cast(sim, result.Target)
			},
		})

		useSpell := character.RegisterSpell(core.SpellConfig{
			ActionID: core.ActionID{ItemID: TheScytheOfElune},
			Flags:    core.SpellFlagNoOnCastComplete | core.SpellFlagOffensiveEquipment,

			Cast: core.CastConfig{
				CD: core.Cooldown{
					Timer:    character.NewTimer(),
					Duration: time.Minute * 10,
				},
				SharedCD: core.Cooldown{
					Timer:    character.GetOffensiveTrinketCD(),
					Duration: time.Minute,
				},
			},

			ApplyEffects: func(sim *core.Simulation, _ *core.Unit, spell *core.Spell) {
				for _, petAgent := range character.PetAgents {
					if pureborn, ok := petAgent.(*guardians.ScytheclawPureborn); ok {
						pureborn.EnableWithTimeout(sim, pureborn, time.Minute)
						break
					}
				}
			},
		})

		character.AddMajorCooldown(core.MajorCooldown{
			Type:  core.CooldownTypeDPS,
			Spell: useSpell,
		})
	})

	///////////////////////////////////////////////////////////////////////////
	//                                 Other
	///////////////////////////////////////////////////////////////////////////

	// https://database.turtlecraft.gg/?item=55106
	// Gives a chance when your harmful spells land to increase the damage of your spells and effects by up to 100 for 6 sec.
	// (Proc chance: 10%, ICD: 18s)
	core.NewItemEffect(BindingsOfContainedMagic, func(agent core.Agent) {
		character := agent.GetCharacter()

		buffAura := character.NewTemporaryStatsAura(
			"Uncontained Magic",
			core.ActionID{SpellID: 51060},
			stats.Stats{stats.SpellDamage: 100},
			time.Second*6,
		)

		core.MakeItemProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Uncontained Magic Passive",
			Callback:   core.CallbackOnSpellHitDealt,
			Outcome:    core.OutcomeLanded,
			ProcMask:   core.ProcMaskSpellDamage,
			ProcChance: 0.10,
			ICD:        time.Second * 18,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				buffAura.Activate(sim)
			},
		})
	})

	// https://database.turtlecraft.gg/?item=33154
	// Your harmful spells have a 10% chance to Conflagrate the target, dealing 100 Fire damage every second for 5 sec.
	// While affected, the target periodically scorches nearby enemies for 40 Fire damage.
	core.NewItemEffect(RingOfBurningTalons, func(agent core.Agent) {
		character := agent.GetCharacter()

		splashSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 52984},
			SpellSchool: core.SpellSchoolFire,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagPassiveSpell | core.SpellFlagNoOnCastComplete,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,
			BonusCoefficient: 0.02,

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				for _, aoeTarget := range sim.Encounter.TargetUnits {
					if aoeTarget == target {
						continue
					}

					spell.CalcAndDealDamage(sim, aoeTarget, 40, spell.OutcomeMagicHit)
				}
			},
		})

		conflagrationSpell := character.RegisterSpell(core.SpellConfig{
			ActionID:    core.ActionID{SpellID: 52983},
			SpellSchool: core.SpellSchoolFire,
			DefenseType: core.DefenseTypeMagic,
			ProcMask:    core.ProcMaskEmpty,
			Flags:       core.SpellFlagPureDot | core.SpellFlagPassiveSpell,

			DamageMultiplier: 1,
			ThreatMultiplier: 1,

			Dot: core.DotConfig{
				Aura: core.Aura{
					ActionID: core.ActionID{SpellID: 52983},
					Label:    "Conflagration",
				},
				NumberOfTicks:    5,
				TickLength:       time.Second,
				BonusCoefficient: 0.04,
				OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
					dot.Snapshot(target, 80, isRollover)
				},
				OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
					splashSpell.Cast(sim, target)
					dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
				},
			},

			ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
				result := spell.CalcOutcome(sim, target, spell.OutcomeMagicHit)
				if result.Landed() {
					spell.Dot(target).Apply(sim)
				}
			},
		})

		core.MakeItemProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:                 "Conflagration Passive",
			Callback:             core.CallbackOnSpellHitDealt,
			Outcome:              core.OutcomeLanded,
			ProcMask:             core.ProcMaskSpellDamage,
			ProcChance:           0.10,
			LandedExecutionIndex: 1,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				conflagrationSpell.Cast(sim, result.Target)
			},
		})
	})

	// https://database.turtlecraft.gg/?item=33095
	// Your direct harmful spells have a 8% chance to increase damage done by magical spells and effects by up to 200, but reduce casting speed by 10%. Lasts 6
	// sec.
	// (Proc chance: 8%)
	core.NewItemEffect(SpellpowerGogglesXtremePlusPlus, func(agent core.Agent) {
		character := agent.GetCharacter()

		buffAura := character.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 52907},
			Label:    "Overpowered",
			Duration: time.Second * 6,
			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatDynamic(sim, stats.SpellDamage, 200)
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				character.AddStatDynamic(sim, stats.SpellDamage, -200)
			},
		}).AttachMultiplyCastSpeed(&character.Unit, 0.9)

		core.MakeItemProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:                 "Overpowered Passive",
			Callback:             core.CallbackOnSpellHitDealt,
			Outcome:              core.OutcomeLanded,
			ProcMask:             core.ProcMaskSpellDamage,
			ProcChance:           0.08,
			LandedExecutionIndex: 1,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				buffAura.Activate(sim)
			},
		})
	})

	// https://database.turtlecraft.gg/?item=55056
	// Your harmful spells have a chance to grant you 'Highborne Insight', increasing your Intellect by 150 for 6 sec.
	// (Proc chance: 10%)
	core.NewItemEffect(SpellwovenNobilityDrape, func(agent core.Agent) {
		character := agent.GetCharacter()

		procAura := character.NewTemporaryStatsAura(
			"Highborne Insight",
			core.ActionID{SpellID: 56546},
			stats.Stats{stats.Intellect: 150},
			time.Second*6,
		)

		core.MakeItemProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:       "Highborne Insight Passive",
			Callback:   core.CallbackOnSpellHitDealt,
			Outcome:    core.OutcomeLanded,
			ProcMask:   core.ProcMaskSpellDamage,
			ProcChance: 0.10,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procAura.Activate(sim)
			},
		})
	})

	// https://database.turtlecraft.gg/?item=58088
	// Your landing damaging spells have a 8% chance to increase your casting speed by 5% for 6 sec. Fire spells have a 50% increased chance to trigger this effect.
	core.NewItemEffect(TrueBandOfSulfuras, func(agent core.Agent) {
		character := agent.GetCharacter()
		buffAura := character.RegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 42027},
			Label:    "Sulfuron Blaze",
			Duration: time.Second * 6,
		}).AttachMultiplyCastSpeed(&character.Unit, 1.05)

		core.MakeItemProcTriggerAura(&character.Unit, core.ProcTrigger{
			Name:     "Sulfuron Blaze Passive",
			Callback: core.CallbackOnSpellHitDealt,
			Outcome:  core.OutcomeLanded,
			ProcMask: core.ProcMaskSpellDamage,
			Handler: func(sim *core.Simulation, spell *core.Spell, result *core.SpellResult) {
				procChance := 0.08
				if spell.SpellSchool.Matches(core.SpellSchoolFire) {
					procChance *= 1.5
				}

				if sim.Proc(character.Unit.ApplyFortuneToProcChance(procChance), "Sulfuron Blaze Passive") {
					buffAura.Activate(sim)
				}
			},
		})
	})

	core.AddEffectsToTest = true
}

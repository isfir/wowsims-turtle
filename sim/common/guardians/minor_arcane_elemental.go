package guardians

import (
	"time"

	"github.com/isfir/wowsims-turtle/sim/core"
	"github.com/isfir/wowsims-turtle/sim/core/stats"
)

const RemainsOfOverwhelmingPower = int32(55093)

// https://database.turtlecraft.gg/?item=55093
// https://database.turtlecraft.gg/?spell=51143
// https://database.turtlecraft.gg/?npc=59992
//
// Minor Arcane Elemental:
// - Summoned for 60 sec
// - While active, owner gains Arcane Link (+55 spell power)
// - Casts Arcane Missiles (Rank 4) every 15 sec

type MinorArcaneElemental struct {
	core.Pet

	arcaneLink         *core.Aura
	arcaneMissiles     *core.Spell
	arcaneMissilesTick *core.Spell

	firstMissilesAt time.Duration
}

func NewMinorArcaneElemental(character *core.Character) *MinorArcaneElemental {
	elementalBaseStats := stats.Stats{
		// From leak; not sure if it got changed since then.
		stats.SpellDamage: 45,
		stats.SpellCrit:   5 * core.CritRatingPerCritChance,
	}

	elemental := &MinorArcaneElemental{
		Pet: core.NewPet("Minor Arcane Elemental", character, elementalBaseStats, minorArcaneElementalStatInheritance(), false, true),
	}

	elemental.Level = 60

	// Not sure about these values; tried to replicate what I saw in logs.
	elemental.EnableAutoAttacks(elemental, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin: 40,
			BaseDamageMax: 80,
			SwingSpeed:    2.0,
			SpellSchool:   core.SpellSchoolPhysical,
		},
		AutoSwingMelee: true,
	})

	elemental.arcaneLink = character.NewTemporaryStatsAura(
		"Arcane Link",
		core.ActionID{SpellID: 51142},
		stats.Stats{stats.SpellPower: 55},
		core.NeverExpires,
	)

	elemental.ApplyOnPetEnable(func(sim *core.Simulation) {
		elemental.firstMissilesAt = sim.CurrentTime + 3*time.Second
		elemental.arcaneLink.Activate(sim)
	})
	elemental.ApplyOnPetDisable(func(sim *core.Simulation) {
		elemental.arcaneLink.Deactivate(sim)
	})

	return elemental
}

func minorArcaneElementalStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{}
	}
}

func (elemental *MinorArcaneElemental) Initialize() {
	elemental.registerArcaneMissilesSpell()
}

func (elemental *MinorArcaneElemental) ExecuteCustomRotation(sim *core.Simulation) {
	if !elemental.IsEnabled() {
		return
	}

	target := elemental.CurrentTarget
	if target == nil {
		target = elemental.Owner.CurrentTarget
	}
	if target == nil {
		elemental.WaitUntil(sim, sim.CurrentTime+time.Second)
		return
	}

	if elemental.IsCasting(sim) || elemental.IsChanneling(sim) {
		elemental.WaitUntil(sim, elemental.Hardcast.Expires)
		return
	}

	if sim.CurrentTime < elemental.firstMissilesAt {
		elemental.WaitUntil(sim, elemental.firstMissilesAt)
		return
	}

	if elemental.arcaneMissiles.IsReady(sim) {
		if elemental.arcaneMissiles.Cast(sim, target) {
			return
		}
	}

	nextActionAt := sim.CurrentTime + elemental.arcaneMissiles.TimeToReady(sim)
	if nextActionAt <= sim.CurrentTime {
		nextActionAt = sim.CurrentTime + time.Millisecond
	}
	elemental.WaitUntil(sim, nextActionAt)
}

func (elemental *MinorArcaneElemental) Reset(sim *core.Simulation) {
	elemental.Disable(sim)
}

func (elemental *MinorArcaneElemental) GetPet() *core.Pet {
	return &elemental.Pet
}

func (elemental *MinorArcaneElemental) registerArcaneMissilesSpell() {
	elemental.arcaneMissilesTick = elemental.RegisterSpell(core.SpellConfig{
		ActionID:     core.ActionID{SpellID: 8419},
		SpellSchool:  core.SpellSchoolArcane,
		DefenseType:  core.DefenseTypeMagic,
		ProcMask:     core.ProcMaskSpellDamage,
		Flags:        core.SpellFlagNoOnCastComplete | core.SpellFlagNoOnDamageDealt,
		MissileSpeed: 20,
		Rank:         4,

		DamageMultiplier: 1,
		ThreatMultiplier: 1,
		BonusCoefficient: 0.328, // Not sure if it's using bonus coeff or fixed damage

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			if !elemental.IsEnabled() {
				return
			}

			// It seems to never miss
			result := spell.CalcDamage(sim, target, 83, spell.OutcomeMagicCrit)

			spell.WaitTravelTime(sim, func(sim *core.Simulation) {
				spell.DealDamage(sim, result)
			})
		},
	})

	elemental.arcaneMissiles = elemental.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 8416},
		SpellSchool: core.SpellSchoolArcane,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskEmpty,
		Flags:       core.SpellFlagChanneled | core.SpellFlagNoOnCastComplete | core.SpellFlagNoOnDamageDealt,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    elemental.NewTimer(),
				Duration: time.Second * 15,
			},
		},

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Arcane Missiles (Minor Arcane Elemental)",
			},
			NumberOfTicks:       5,
			TickLength:          time.Second,
			AffectedByCastSpeed: false,
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				elemental.arcaneMissilesTick.Cast(sim, target)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			spell.Dot(target).Apply(sim)
		},
	})
}

func constructMinorArcaneElemental(character *core.Character) {
	if character.HasTrinketEquipped(RemainsOfOverwhelmingPower) {
		character.AddPet(NewMinorArcaneElemental(character))
	}
}

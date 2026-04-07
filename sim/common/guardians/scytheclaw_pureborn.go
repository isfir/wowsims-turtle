package guardians

import (
	"time"

	"github.com/isfir/wowsims-turtle/sim/core"
	"github.com/isfir/wowsims-turtle/sim/core/stats"
)

const TheScytheOfElune = int32(55505)

// https://database.turtlecraft.gg/?item=55505
// https://database.turtlecraft.gg/?spell=57667
// https://database.turtlecraft.gg/?npc=40054
//
// Scytheclaw Pureborn:
// - Summoned for 60 sec
// - Crooked Claw roughly every 9 sec
// - Pure Moonfire roughly every 12 sec
// - Elune's Courage roughly every 30 sec

type ScytheclawPureborn struct {
	core.Pet

	crookedClaw   *core.Spell
	pureMoonfire  *core.Spell
	elunesCourage *core.Spell

	courageBuffs []*core.Aura

	firstCourageAt  time.Duration
	firstCrookedAt  time.Duration
	firstMoonfireAt time.Duration
}

func NewScytheclawPureborn(character *core.Character) *ScytheclawPureborn {
	baseStats := stats.Stats{
		stats.Health: 7581,
		stats.Mana:   23864,
		stats.Armor:  4293,
	}

	pureborn := &ScytheclawPureborn{
		Pet: core.NewPet("Scytheclaw Pureborn", character, baseStats, scytheclawPurebornStatInheritance(), false, true),
	}

	pureborn.Level = 61
	pureborn.EnableManaBar()

	pureborn.EnableAutoAttacks(pureborn, core.AutoAttackOptions{
		MainHand: core.Weapon{
			BaseDamageMin: 40.0,
			BaseDamageMax: 80.0,
			SwingSpeed:    2.0,
			SpellSchool:   core.SpellSchoolPhysical,
		},
		AutoSwingMelee: true,
	})

	pureborn.ApplyOnPetEnable(func(sim *core.Simulation) {
		pureborn.firstCourageAt = sim.CurrentTime + 2*time.Second
		pureborn.firstCrookedAt = sim.CurrentTime + 9*time.Second
		pureborn.firstMoonfireAt = sim.CurrentTime + 12*time.Second
	})

	return pureborn
}

func scytheclawPurebornStatInheritance() core.PetStatInheritance {
	return func(ownerStats stats.Stats) stats.Stats {
		return stats.Stats{}
	}
}

func (pureborn *ScytheclawPureborn) Initialize() {
	pureborn.registerElunesCourageSpell()
	pureborn.registerCrookedClawSpell()
	pureborn.registerPureMoonfireSpell()
}

func (pureborn *ScytheclawPureborn) ExecuteCustomRotation(sim *core.Simulation) {
	if !pureborn.IsEnabled() {
		return
	}

	target := pureborn.Owner.CurrentTarget
	if target == nil {
		target = pureborn.CurrentTarget
	}
	if target == nil {
		pureborn.WaitUntil(sim, sim.CurrentTime+time.Second)
		return
	}

	if pureborn.IsCasting(sim) || pureborn.IsChanneling(sim) {
		pureborn.WaitUntil(sim, pureborn.Hardcast.Expires)
		return
	}

	if sim.CurrentTime >= pureborn.firstCourageAt && pureborn.elunesCourage.IsReady(sim) {
		if pureborn.elunesCourage.Cast(sim, target) {
			return
		}
	}
	if sim.CurrentTime >= pureborn.firstCrookedAt && pureborn.crookedClaw.IsReady(sim) {
		if pureborn.crookedClaw.Cast(sim, target) {
			return
		}
	}
	if sim.CurrentTime >= pureborn.firstMoonfireAt && pureborn.pureMoonfire.IsReady(sim) {
		if pureborn.pureMoonfire.Cast(sim, target) {
			return
		}
	}

	nextActionAt := sim.CurrentTime + time.Second

	if sim.CurrentTime < pureborn.firstCourageAt {
		nextActionAt = pureborn.firstCourageAt
	} else if t := sim.CurrentTime + pureborn.elunesCourage.TimeToReady(sim); t < nextActionAt {
		nextActionAt = t
	}

	if sim.CurrentTime < pureborn.firstCrookedAt {
		if pureborn.firstCrookedAt < nextActionAt {
			nextActionAt = pureborn.firstCrookedAt
		}
	} else if t := sim.CurrentTime + pureborn.crookedClaw.TimeToReady(sim); t < nextActionAt {
		nextActionAt = t
	}

	if sim.CurrentTime < pureborn.firstMoonfireAt {
		if pureborn.firstMoonfireAt < nextActionAt {
			nextActionAt = pureborn.firstMoonfireAt
		}
	} else if t := sim.CurrentTime + pureborn.pureMoonfire.TimeToReady(sim); t < nextActionAt {
		nextActionAt = t
	}

	if nextActionAt <= sim.CurrentTime {
		nextActionAt = sim.CurrentTime + time.Millisecond
	}
	pureborn.WaitUntil(sim, nextActionAt)
}

func (pureborn *ScytheclawPureborn) Reset(sim *core.Simulation) {
	pureborn.Disable(sim)
}

func (pureborn *ScytheclawPureborn) GetPet() *core.Pet {
	return &pureborn.Pet
}

func (pureborn *ScytheclawPureborn) registerElunesCourageSpell() {
	pureborn.courageBuffs = make([]*core.Aura, 0, len(pureborn.Owner.Party.PlayersAndPets)+1)

	includedSelf := false
	for _, raidAgent := range pureborn.Owner.Party.PlayersAndPets {
		raidChar := raidAgent.GetCharacter()
		unit := &raidChar.Unit

		if unit == &pureborn.Unit {
			includedSelf = true
		}

		buff := raidChar.GetOrRegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 57692},
			Label:    "Elune's Courage",
			Duration: time.Second * 8,
		}).AttachMultiplyAttackSpeed(unit, 1.10).AttachMultiplyCastSpeed(unit, 1.10)

		pureborn.courageBuffs = append(pureborn.courageBuffs, buff)
	}

	if !includedSelf {
		buff := pureborn.GetOrRegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 57692},
			Label:    "Elune's Courage",
			Duration: time.Second * 8,
		}).AttachMultiplyAttackSpeed(&pureborn.Unit, 1.10).AttachMultiplyCastSpeed(&pureborn.Unit, 1.10)

		pureborn.courageBuffs = append(pureborn.courageBuffs, buff)
	}

	pureborn.elunesCourage = pureborn.RegisterSpell(core.SpellConfig{
		ActionID: core.ActionID{SpellID: 57692},
		Flags:    core.SpellFlagHelpful,

		ManaCost: core.ManaCostOptions{
			FlatCost: 150,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    pureborn.NewTimer(),
				Duration: time.Second * 30,
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			for _, aura := range pureborn.courageBuffs {
				aura.Activate(sim)
			}
		},
	})
}

func (pureborn *ScytheclawPureborn) registerCrookedClawSpell() {
	physicalVulnAuras := pureborn.Owner.NewEnemyAuraArray(func(target *core.Unit) *core.Aura {
		return target.GetOrRegisterAura(core.Aura{
			ActionID: core.ActionID{SpellID: 57669},
			Label:    "Crooked Claw",
			Duration: time.Second * 5,

			OnGain: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical] *= 1.03
			},
			OnExpire: func(aura *core.Aura, sim *core.Simulation) {
				aura.Unit.PseudoStats.SchoolDamageTakenMultiplier[stats.SchoolIndexPhysical] /= 1.03
			},
		})
	})

	pureborn.crookedClaw = pureborn.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 57669},
		SpellSchool: core.SpellSchoolPhysical,
		DefenseType: core.DefenseTypeMelee,
		ProcMask:    core.ProcMaskEmpty,

		ManaCost: core.ManaCostOptions{
			FlatCost: 250,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    pureborn.NewTimer(),
				Duration: time.Second * 9,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealDamage(sim, target, sim.Roll(300, 351), spell.OutcomeMeleeSpecialHitAndCrit)
			if result.Landed() {
				physicalVulnAuras.Get(target).Activate(sim)
			}
		},
	})
}

func (pureborn *ScytheclawPureborn) registerPureMoonfireSpell() {
	pureborn.pureMoonfire = pureborn.RegisterSpell(core.SpellConfig{
		ActionID:    core.ActionID{SpellID: 57668},
		SpellSchool: core.SpellSchoolArcane,
		DefenseType: core.DefenseTypeMagic,
		ProcMask:    core.ProcMaskSpellDamage,

		ManaCost: core.ManaCostOptions{
			FlatCost: 325,
		},

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
			CD: core.Cooldown{
				Timer:    pureborn.NewTimer(),
				Duration: time.Second * 12,
			},
		},

		DamageMultiplier: 1,
		ThreatMultiplier: 1,

		Dot: core.DotConfig{
			Aura: core.Aura{
				Label: "Pure Moonfire",
			},
			NumberOfTicks: 2,
			TickLength:    time.Second * 3,

			OnSnapshot: func(sim *core.Simulation, target *core.Unit, dot *core.Dot, isRollover bool) {
				dot.Snapshot(target, 75, isRollover)
			},
			OnTick: func(sim *core.Simulation, target *core.Unit, dot *core.Dot) {
				dot.CalcAndDealPeriodicSnapshotDamage(sim, target, dot.OutcomeTick)
			},
		},

		ApplyEffects: func(sim *core.Simulation, target *core.Unit, spell *core.Spell) {
			result := spell.CalcAndDealDamage(sim, target, sim.Roll(300, 351), spell.OutcomeMagicHitAndCrit)
			if result.Landed() {
				spell.Dot(target).Apply(sim)
			}
		},
	})
}

func constructScytheclawPureborn(character *core.Character) {
	if character.HasTrinketEquipped(TheScytheOfElune) {
		character.AddPet(NewScytheclawPureborn(character))
	}
}

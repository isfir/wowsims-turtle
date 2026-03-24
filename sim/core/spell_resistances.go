package core

import (
	"math"

	"github.com/isfir/wowsims-turtle/sim/core/stats"
)

func (result *SpellResult) applyResistances(sim *Simulation, spell *Spell, isPeriodic bool, attackTable *AttackTable) {
	resistanceMultiplier, outcome := spell.ResistanceMultiplier(sim, isPeriodic, attackTable)

	result.Damage *= resistanceMultiplier
	result.Outcome |= outcome

	result.ResistanceMultiplier = resistanceMultiplier
	result.PreOutcomeDamage = result.Damage
}

// Modifies damage based on Armor or Magic resistances, depending on the damage type.
func (spell *Spell) ResistanceMultiplier(sim *Simulation, isPeriodic bool, attackTable *AttackTable) (float64, HitOutcome) {
	if spell.Flags.Matches(SpellFlagIgnoreResists) {
		return 1, OutcomeEmpty
	}

	if spell.SpellSchool.Matches(SpellSchoolPhysical) {
		if spell.SchoolIndex == stats.SchoolIndexPhysical || MultiSchoolShouldUseArmor(spell, attackTable.Defender) {
			// All physical dots (Bleeds) ignore armor.
			if isPeriodic {
				return 1, OutcomeEmpty
			}

			// Physical resistance (armor).
			return attackTable.GetArmorDamageModifier(), OutcomeEmpty
		}
	}

	// Magical resistance.
	if spell.Flags.Matches(SpellFlagBinary) {
		return 1, OutcomeEmpty
	}

	resistanceRoll := sim.RandomFloat("Partial Resist")

	threshold00, threshold25, threshold50 := attackTable.GetPartialResistThresholds(spell, isPeriodic)
	//if sim.Log != nil {
	//	sim.Log("Resist thresholds: %0.04f, %0.04f, %0.04f", threshold00, threshold25, threshold50)
	//}

	if resistanceRoll > threshold00 {
		// No partial resist.
		return 1, OutcomeEmpty
	} else if resistanceRoll > threshold25 {
		return 0.75, OutcomePartial1_4
	} else if resistanceRoll > threshold50 {
		return 0.5, OutcomePartial2_4
	} else {
		return 0.25, OutcomePartial3_4
	}
}

// Decide whether to use armor for physical multi school spells.
//
// TODO: This is most likely not accurate for the case: armor near resistance but not 0
//
// A short test showed that the game uses armor if it's far enough below resistance,
// but not simply if it's lower.
// 49 (and above) armor vs 57 res => used resistance
// 7 (and below) armor vs 57 res => used armor/no partials anymore
// If level based resist is used in this decission process is also not known as it was tested PvP.
//
// For most purposes this should work fine for now, but should be properly tested and fixed if
// spells using it become important and boss armor can actually go below (level based) resistance values.
func MultiSchoolShouldUseArmor(spell *Spell, target *Unit) bool {
	resistance := 100000.0
	lowestIsArmor := true
	for _, baseSchoolIndex := range spell.SchoolBaseIndices {
		resiVal := target.GetResistanceForSchool(baseSchoolIndex)
		if resiVal < resistance {
			resistance = resiVal
			lowestIsArmor = baseSchoolIndex == stats.SchoolIndexPhysical
		}
	}
	return lowestIsArmor
}

func (at *AttackTable) GetArmorDamageModifier() float64 {
	armorPenRating := at.Attacker.stats[stats.ArmorPenetration]
	defenderArmor := max(at.Defender.Armor()-armorPenRating, 0.0)
	return 1 - defenderArmor/(defenderArmor+400+85*float64(at.Attacker.Level))
}

func (at *AttackTable) GetPartialResistThresholds(spell *Spell, dot bool) (float64, float64, float64) {
	return at.Defender.partialResistRollThresholds(spell, at.Attacker, dot)
}

func (at *AttackTable) GetBinaryHitChance(spell *Spell) float64 {
	return at.Defender.binaryHitChance(spell, at.Attacker)
}

// Only for base schools!
func (unit *Unit) GetResistanceForSchool(schoolIndex stats.SchoolIndex) float64 {
	switch schoolIndex {
	case stats.SchoolIndexNone:
		return 0
	case stats.SchoolIndexPhysical:
		return unit.GetStat(stats.Armor)
	case stats.SchoolIndexArcane:
		return unit.GetStat(stats.ArcaneResistance)
	case stats.SchoolIndexFire:
		return unit.GetStat(stats.FireResistance)
	case stats.SchoolIndexFrost:
		return unit.GetStat(stats.FrostResistance)
	case stats.SchoolIndexHoly:
		return 0 // Holy resistance doesn't exist.
	case stats.SchoolIndexNature:
		return unit.GetStat(stats.NatureResistance)
	case stats.SchoolIndexShadow:
		return unit.GetStat(stats.ShadowResistance)
	default:
		return 0
	}
}

func (unit *Unit) resistCoeff(spell *Spell, attacker *Unit, binary bool) float64 {
	if spell.SchoolIndex <= stats.SchoolIndexPhysical {
		return 0
	}

	var resistance float64

	if spell.SchoolIndex.IsMultiSchool() {
		// Multi school: Choose the lowest resistance available.
		resistance = 1000.0
		for _, baseSchoolIndex := range spell.SchoolBaseIndices {
			resiVal := unit.GetResistanceForSchool(baseSchoolIndex)
			if resiVal < resistance {
				resistance = resiVal
			}
		}
	} else {
		resistance = unit.GetResistanceForSchool(spell.SchoolIndex)
	}

	resistance = max(0, resistance-attacker.stats[stats.SpellPenetration])

	resistanceCap := float64(attacker.Level * 5)
	resistanceCoef := resistance / resistanceCap

	if !binary && unit.Type == EnemyUnit && unit.Level > attacker.Level {
		avgMitigationAdded := AverageMagicPartialResistPerLevelMultiplier * float64(unit.Level-attacker.Level)
		// coeff is scaled 0..1, while avg mitigation is 0..0.75
		resistanceCoef += avgMitigationAdded / 0.75
	}

	return min(1, resistanceCoef)
}

func (unit *Unit) binaryHitChance(spell *Spell, attacker *Unit) float64 {
	resistCoeff := unit.resistCoeff(spell, attacker, true)
	return 1 - 0.75*resistCoeff
}

type resistanceValues struct {
	resist100    float64
	resist75     float64
	resist50     float64
	resist25     float64
	resist0      float64
	chanceResist float64
}

var partialResistTable = []resistanceValues{
	{0, 0, 0, 0, 100, 0},    // 0
	{0, 0, 2, 6, 92, 3},     // 10
	{0, 1, 4, 12, 84, 5},    // 20
	{0, 1, 5, 18, 76, 8},    // 30
	{0, 1, 7, 23, 69, 10},   // 40
	{0, 2, 9, 28, 61, 13},   // 50
	{0, 2, 11, 33, 54, 15},  // 60
	{0, 2, 13, 37, 37, 18},  // 70
	{0, 3, 15, 41, 41, 20},  // 80
	{1, 3, 17, 46, 36, 23},  // 90
	{1, 4, 19, 47, 29, 25},  // 100
	{1, 5, 21, 48, 24, 28},  // 110
	{1, 6, 24, 49, 20, 30},  // 120
	{1, 8, 28, 47, 17, 33},  // 130
	{1, 9, 33, 43, 14, 35},  // 140
	{1, 11, 37, 39, 12, 38}, // 150
	{1, 13, 41, 35, 10, 40}, // 160
	{1, 16, 45, 30, 8, 43},  // 170
	{1, 18, 48, 26, 7, 45},  // 180
	{2, 20, 48, 24, 6, 48},  // 190
	{4, 23, 48, 21, 4, 50},  // 200
	{5, 25, 47, 19, 3, 53},  // 210
	{7, 28, 45, 17, 2, 55},  // 220
	{9, 31, 43, 16, 2, 58},  // 230
	{11, 34, 40, 14, 1, 60}, // 240
	{13, 37, 37, 12, 1, 62}, // 250
	{15, 41, 33, 10, 1, 65}, // 260
	{18, 44, 29, 8, 1, 68},  // 270
	{20, 48, 25, 7, 1, 70},  // 280
	{23, 51, 20, 5, 1, 73},  // 290
	{25, 55, 16, 3, 1, 75},  // 300
}

// Returns raw VMaNGOS-style bucket chances in percent.
// resist75Pct already includes resist100, because server folds 100% into 75%.
func interpolatePartialResistBucketsPct(resistanceChancePct float64) (resist25Pct float64, resist50Pct float64, resist75Pct float64) {
	resistanceChancePct = max(0, min(75, resistanceChancePct))
	if resistanceChancePct <= 0 {
		return 0, 0, 0
	}

	prev := partialResistTable[0]
	next := partialResistTable[len(partialResistTable)-1]

	for i := 1; i < len(partialResistTable); i++ {
		if partialResistTable[i].chanceResist >= resistanceChancePct {
			prev = partialResistTable[i-1]
			next = partialResistTable[i]
			break
		}
	}

	coeff := 0.0
	if next.chanceResist > prev.chanceResist {
		coeff = (resistanceChancePct - prev.chanceResist) / (next.chanceResist - prev.chanceResist)
	}

	resist25 := prev.resist25 + (next.resist25-prev.resist25)*coeff
	resist50 := prev.resist50 + (next.resist50-prev.resist50)*coeff
	resist75 := prev.resist75 + (next.resist75-prev.resist75)*coeff
	resist100 := prev.resist100 + (next.resist100-prev.resist100)*coeff

	return resist25, resist50, resist75 + resist100
}

// VMaNGOS uses urand(0, 99) and compares "ran < cumulativePct".
// Number of winning integer buckets is ceil(cumulativePct), clamped to [0,100].
func discreteRollThresholdFromPct(cumulativePct float64) float64 {
	if cumulativePct <= 0 {
		return 0
	}
	if cumulativePct >= 100 {
		return 1
	}
	return math.Ceil(cumulativePct) / 100.0
}

// Roll threshold for each type of partial resist.
func (unit *Unit) partialResistRollThresholds(spell *Spell, attacker *Unit, dot bool) (float64, float64, float64) {
	// Convert coeff (0..1) to average resist chance (0..0.75).
	resistanceChance := 0.75 * unit.resistCoeff(spell, attacker, false)

	// Applies DoT reduction to the final resistance chance.
	if dot {
		resistanceChance *= 0.1
	}

	resist25Pct, resist50Pct, resist75Pct := interpolatePartialResistBucketsPct(resistanceChance * 100.0)

	// Build cumulative chances in percent, then quantize to 1% buckets.
	threshold50 := discreteRollThresholdFromPct(resist75Pct)
	threshold25 := discreteRollThresholdFromPct(resist75Pct + resist50Pct)
	threshold00 := discreteRollThresholdFromPct(resist75Pct + resist50Pct + resist25Pct)

	return threshold00, threshold25, threshold50
}

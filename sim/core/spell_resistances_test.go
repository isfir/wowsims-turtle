package core

import (
	"math"
	"testing"

	"github.com/isfir/wowsims-turtle/sim/core/proto"
	"github.com/isfir/wowsims-turtle/sim/core/simsignals"
	"github.com/isfir/wowsims-turtle/sim/core/stats"
)

func expectedPartialResistThresholds(resistanceChance float64) (float64, float64, float64) {
	resistanceChance = max(0, min(0.75, resistanceChance))

	resist25Pct, resist50Pct, resist75Pct := interpolatePartialResistBucketsPct(resistanceChance * 100.0)

	threshold50 := discreteRollThresholdFromPct(resist75Pct)
	threshold25 := discreteRollThresholdFromPct(resist75Pct + resist50Pct)
	threshold00 := discreteRollThresholdFromPct(resist75Pct + resist50Pct + resist25Pct)

	return threshold00, threshold25, threshold50
}

func Test_PartialResistsVsPlayer(t *testing.T) {
	attacker := &Unit{
		Type:  EnemyUnit,
		Level: 63,
		stats: stats.Stats{},
	}
	defender := &Unit{
		Type:  PlayerUnit,
		Level: 60,
		stats: stats.Stats{},
	}

	attackTable := NewAttackTable(attacker, defender, nil)

	sim := NewSim(&proto.RaidSimRequest{
		SimOptions: &proto.SimOptions{},
		Encounter:  &proto.Encounter{},
		Raid:       &proto.Raid{},
	}, simsignals.CreateSignals())

	schoolMask := SpellSchoolFire
	spell := &Spell{
		SpellSchool:       schoolMask,
		SchoolIndex:       schoolMask.GetSchoolIndex(),
		SchoolBaseIndices: schoolMask.GetBaseIndices(),
	}

	for resist := 0; resist < 5_000; resist += 1 {
		defender.stats[stats.FireResistance] = float64(resist)

		threshold00, threshold25, threshold50 := attackTable.Defender.partialResistRollThresholds(spell, attackTable.Attacker, false)

		resistanceChance := 0.75 * attackTable.Defender.resistCoeff(spell, attackTable.Attacker, false)
		expectedT00, expectedT25, expectedT50 := expectedPartialResistThresholds(resistanceChance)

		if math.Abs(threshold00-expectedT00) > 1e-9 ||
			math.Abs(threshold25-expectedT25) > 1e-9 ||
			math.Abs(threshold50-expectedT50) > 1e-9 {
			t.Errorf("resist = %d, thresholds = (%.2f, %.2f, %.2f), expected = (%.2f, %.2f, %.2f)",
				resist, threshold00, threshold25, threshold50, expectedT00, expectedT25, expectedT50)
			return
		}

		expectedAr, _, _, _, _ := GetChancesAndMitFromThresholds(expectedT00, expectedT25, expectedT50)

		const n = 10_000

		outcomes := make(map[HitOutcome]int, n)
		var totalDamage float64
		for iter := 0; iter < n; iter++ {
			result := SpellResult{
				Outcome: OutcomeHit,
				Damage:  1000,
			}

			result.applyResistances(sim, spell, false, attackTable)

			outcomes[result.Outcome]++
			totalDamage += result.Damage
		}

		if math.Abs(expectedAr-(1-totalDamage/float64(1000*n))) > 0.01 {
			t.Logf("after %d iterations, resist = %d, ar = %.2f%% vs. damage lost = %.2f%%, outcomes = %v\n",
				n, resist, expectedAr*100, 100-100*totalDamage/float64(1000*n), outcomes)
		}
	}
}

func GetChancesAndMitFromThresholds(t0 float64, t25 float64, t50 float64) (float64, float64, float64, float64, float64) {
	chance0 := 1 - t0
	chance25 := t0 - t25
	chance50 := t25 - t50
	chance75 := t50
	avgResist := chance25*0.25 + chance50*0.50 + chance75*0.75
	return avgResist, chance0, chance25, chance50, chance75
}

func CloseEnough(f1 float64, f2 float64, eps float64) bool {
	return math.Abs(f1-f2) < eps
}

func ResistanceCheck(t *testing.T, isDoT bool) {
	attacker := &Unit{
		Type:  PlayerUnit,
		Level: 60,
		stats: stats.Stats{},
	}
	defender := &Unit{
		Type:  EnemyUnit,
		Level: attacker.Level + 3,
		stats: stats.Stats{},
	}

	attackTable := NewAttackTable(attacker, defender, nil)

	schoolMask := SpellSchoolFromIndex(stats.SchoolIndexNature)
	spell := &Spell{
		SpellSchool:       schoolMask,
		SchoolIndex:       schoolMask.GetSchoolIndex(),
		SchoolBaseIndices: schoolMask.GetBaseIndices(),
	}

	maxResist := float64(attacker.Level) * 5.0

	// Resist coeff itself is unchanged by DoT handling now.
	defender.stats[stats.NatureResistance] = 0
	coef := defender.resistCoeff(spell, attacker, false)
	if coef != 0.08 {
		t.Errorf("Resist coef is %.3f at 0 resistance, but should be 0.08!", coef)
		return
	}

	// VMaNGOS-style known values at 200 resistance vs +3 target.
	defender.stats[stats.NatureResistance] = 200
	expectedMitigation := 0.5425
	expectedChances := []float64{0.02, 0.17, 0.44, 0.37}
	if isDoT {
		expectedMitigation = 0.065
		expectedChances = []float64{0.81, 0.13, 0.05, 0.01}
	}

	threshold00, threshold25, threshold50 := attackTable.GetPartialResistThresholds(spell, isDoT)
	avgResist, chance0, chance25, chance50, chance75 := GetChancesAndMitFromThresholds(threshold00, threshold25, threshold50)

	if !CloseEnough(avgResist, expectedMitigation, 0.01) {
		t.Errorf("Avg mitigation %.3f at 200 resistance, but should be %.3f!", avgResist, expectedMitigation)
		return
	}
	if !CloseEnough(chance0, expectedChances[0], 0.01) ||
		!CloseEnough(chance25, expectedChances[1], 0.01) ||
		!CloseEnough(chance50, expectedChances[2], 0.01) ||
		!CloseEnough(chance75, expectedChances[3], 0.01) {
		t.Errorf("Bucket chances do not match known values at 200 resistance. Known %v, returned %v!",
			expectedChances, []float64{chance0, chance25, chance50, chance75})
		return
	}

	// Check various resistance values
	for resist := 0.0; resist < maxResist; resist += 1.0 {
		defender.stats[stats.NatureResistance] = resist

		resistanceCap := float64(attacker.Level * 5)
		levelBased := float64(max(defender.Level-attacker.Level, 0)) * 0.02
		expectedCoef := min(1, resist/resistanceCap+levelBased*1/0.75)

		resistCoef := defender.resistCoeff(spell, attacker, false)
		if math.Abs(resistCoef-expectedCoef) > 0.001 {
			t.Errorf("Resist coef is %.3f but expected %.3f at resistance %f", resistCoef, expectedCoef, resist)
			return
		}

		expectedResistanceChance := 0.75 * expectedCoef
		if isDoT {
			expectedResistanceChance *= 0.1
		}

		expectedT00, expectedT25, expectedT50 := expectedPartialResistThresholds(expectedResistanceChance)
		expectedAvgMitigation, _, _, _, _ := GetChancesAndMitFromThresholds(expectedT00, expectedT25, expectedT50)

		threshold00, threshold25, threshold50 := attackTable.GetPartialResistThresholds(spell, isDoT)

		if math.Abs(threshold00-expectedT00) > 1e-9 ||
			math.Abs(threshold25-expectedT25) > 1e-9 ||
			math.Abs(threshold50-expectedT50) > 1e-9 {
			t.Errorf("resist = %.2f, thresholds = (%.2f, %.2f, %.2f), expected = (%.2f, %.2f, %.2f)",
				resist, threshold00, threshold25, threshold50, expectedT00, expectedT25, expectedT50)
			return
		}

		avgResist, _, _, _, _ := GetChancesAndMitFromThresholds(threshold00, threshold25, threshold50)
		if math.Abs(avgResist-expectedAvgMitigation) > 0.005 {
			t.Errorf("resist = %.2f, thresholds = %f, resultingAr = %.2f%%, expectedAr = %.2f%%",
				resist, threshold00, avgResist, expectedAvgMitigation)
			return
		}
	}
}

func Test_ResistsVsBoss(t *testing.T) {
	t.Run("Direct", func(t *testing.T) { ResistanceCheck(t, false) })
	t.Run("DoT", func(t *testing.T) { ResistanceCheck(t, true) })
}

func Test_ResistBinary(t *testing.T) {
	attacker := &Unit{
		Type:  PlayerUnit,
		Level: 60,
		stats: stats.Stats{},
	}
	defender := &Unit{
		Type:  EnemyUnit,
		Level: attacker.Level + 3,
		stats: stats.Stats{},
	}

	attackTable := NewAttackTable(attacker, defender, nil)

	schoolMask := SpellSchoolFromIndex(stats.SchoolIndexNature)
	spell := &Spell{
		Flags:             SpellFlagBinary,
		SpellSchool:       schoolMask,
		SchoolIndex:       schoolMask.GetSchoolIndex(),
		SchoolBaseIndices: schoolMask.GetBaseIndices(),
	}

	// Check if coef is 0.0 at 0 resistance, binary spells do not get level based resistance!
	defender.stats[stats.NatureResistance] = 0
	coef := defender.resistCoeff(spell, attacker, true)
	if coef != 0.0 {
		t.Errorf("Resist coef is %.3f at 0 resistance for binary spell, but should be 0.0!", coef)
		return
	}

	// Should not partial resist
	dmgMult, outcome := spell.ResistanceMultiplier(nil, false, attackTable)
	if dmgMult != 1 || outcome != OutcomeEmpty {
		t.Errorf("ResistanceMultiplier for binary spell did not return mult=1 and empty outcome, got %.3f and outcome %d!", dmgMult, outcome)
		return
	}

	// Hit chance
	tests := [][]float64{
		{0.0, 1},
		{100.0, 0.75},
		{200.0, 0.5},
		{300.0, 0.25},
	}
	for _, test := range tests {
		resistance := test[0]
		defender.stats[stats.NatureResistance] = resistance
		expectedResult := test[1]
		result := attackTable.GetBinaryHitChance(spell)
		if !CloseEnough(result, expectedResult, 0.000001) {
			t.Errorf("Binary hit chance result at %.0f resistance was %.3f, expected %.3f!", resistance, result, expectedResult)
			return
		}
	}
}

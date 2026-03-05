package dpswarrior

import (
	"testing"

	_ "github.com/isfir/wowsims-turtle/sim/common" // imported to get item effects included.
	"github.com/isfir/wowsims-turtle/sim/core"
	"github.com/isfir/wowsims-turtle/sim/core/proto"
)

func init() {
	RegisterDpsWarrior()
}

func TestP1DPSWarrior(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassWarrior,
			Phase:      1,
			Race:       proto.Race_RaceOrc,
			OtherRaces: []proto.Race{proto.Race_RaceHuman},

			Talents:  P1Talents,
			GearSet:  core.GetGearSet("../../../ui/warrior/gear_sets", "phase_1"),
			Rotation: core.GetAplRotation("../../../ui/warrior/apls", "dps_reck"),
			OtherRotations: []core.RotationCombo{
				core.GetAplRotation("../../../ui/warrior/apls", "dps_no_reck"),
			},
			Buffs:       core.FullBuffs,
			Consumes:    P1Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "DPS", SpecOptions: PlayerOptionsFury},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
	}))
}

var P1Talents = "30305001302-05050005525010051"

var P1Consumes = core.ConsumesCombo{
	Label: "P1-Consumes",
	Consumes: &proto.Consumes{
		AgilityElixir:   proto.AgilityElixir_ElixirOfTheMongoose,
		AttackPowerBuff: proto.AttackPowerBuff_JujuMight,
		DefaultPotion:   proto.Potions_MightyRagePotion,
		Food:            proto.Food_FoodSmokedDesertDumpling,
		MainHandImbue:   proto.WeaponImbue_Windfury,
		OffHandImbue:    proto.WeaponImbue_ElementalSharpeningStone,
		StrengthBuff:    proto.StrengthBuff_JujuPower,
	},
}

var PlayerOptionsArms = &proto.Player_Warrior{
	Warrior: &proto.Warrior{
		Options: warriorOptions,
	},
}

var PlayerOptionsFury = &proto.Player_Warrior{
	Warrior: &proto.Warrior{
		Options: warriorOptions,
	},
}

var warriorOptions = &proto.Warrior_Options{
	StartingRage: 50,
	Shout:        proto.WarriorShout_WarriorShoutBattle,
}

var ItemFilters = core.ItemFilter{
	ArmorType: proto.ArmorType_ArmorTypePlate,

	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeAxe,
		proto.WeaponType_WeaponTypeSword,
		proto.WeaponType_WeaponTypeMace,
		proto.WeaponType_WeaponTypeDagger,
		proto.WeaponType_WeaponTypeFist,
	},
}

var Stats = []proto.Stat{
	proto.Stat_StatStrength,
	proto.Stat_StatAgility,
	proto.Stat_StatAttackPower,
	proto.Stat_StatMeleeCrit,
	proto.Stat_StatMeleeHit,
}

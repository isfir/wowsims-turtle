package warden

import (
	"testing"

	"github.com/isfir/wowsims-turtle/sim/core"
	"github.com/isfir/wowsims-turtle/sim/core/proto"
)

func init() {
	RegisterWardenShaman()
}

func TestWardenShaman(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassShaman,
			Phase:      1,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceOrc},

			Talents:     DefaultTalents,
			GearSet:     core.GetGearSet("../../../ui/warden_shaman/gear_sets", "blank"),
			Rotation:    core.GetAplRotation("../../../ui/warden_shaman/apls", "default"),
			Buffs:       core.FullBuffs,
			Consumes:    Phase1Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Default", SpecOptions: PlayerOptionsBasic},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
	}))
}

var PlayerOptionsBasic = &proto.Player_WardenShaman{
	WardenShaman: &proto.WardenShaman{
		Options: &proto.WardenShaman_Options{},
	},
}

var DefaultTalents = "5203015-0505000145503151"

var Phase1Consumes = core.ConsumesCombo{
	Label: "P1-Consumes",
	Consumes: &proto.Consumes{
		AttackPowerBuff: proto.AttackPowerBuff_JujuMight,
		AgilityElixir:   proto.AgilityElixir_ElixirOfTheMongoose,
		DefaultConjured: proto.Conjured_ConjuredDemonicRune,
		DefaultPotion:   proto.Potions_MajorManaPotion,
		FirePowerBuff:   proto.FirePowerBuff_ElixirOfGreaterFirepower,
		Flask:           proto.Flask_FlaskOfTheTitans,
		Food:            proto.Food_FoodBlessSunfruit,
		MainHandImbue:   proto.WeaponImbue_RockbiterWeapon,
		SpellPowerBuff:  proto.SpellPowerBuff_GreaterArcaneElixir,
		StrengthBuff:    proto.StrengthBuff_JujuPower,
	},
}

var ItemFilters = core.ItemFilter{
	ArmorType: proto.ArmorType_ArmorTypeMail,

	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeAxe,
		proto.WeaponType_WeaponTypeSword,
		proto.WeaponType_WeaponTypeMace,
		proto.WeaponType_WeaponTypeDagger,
		proto.WeaponType_WeaponTypeFist,
		proto.WeaponType_WeaponTypeShield,
	},
}

var Stats = []proto.Stat{
	proto.Stat_StatStrength,
	proto.Stat_StatAttackPower,
	proto.Stat_StatSpellPower,
	proto.Stat_StatArmor,
	proto.Stat_StatDodge,
	proto.Stat_StatParry,
	proto.Stat_StatBlockValue,
	proto.Stat_StatDefense,
}

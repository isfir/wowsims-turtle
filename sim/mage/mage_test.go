package mage

import (
	"testing"

	_ "github.com/isfir/wowsims-turtle/sim/common"
	"github.com/isfir/wowsims-turtle/sim/core"
	"github.com/isfir/wowsims-turtle/sim/core/proto"
)

func init() {
	RegisterMage()
}

func TestMage(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:        proto.Class_ClassMage,
			Phase:        1,
			Race:         proto.Race_RaceGnome,
			OtherRaces:   []proto.Race{proto.Race_RaceTroll},
			InfiniteMana: true,

			Talents:     DefaultTalents,
			GearSet:     core.GetGearSet("../../ui/mage/gear_sets", "arcane_st_bis"),
			Rotation:    core.GetAplRotation("../../ui/mage/apls", "arcane_st"),
			Buffs:       core.FullBuffs,
			Consumes:    DefaultConsumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Arcane", SpecOptions: PlayerOptions},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatSpellPower,
			StatsToWeigh:    Stats,
		},
	}))
}

var DefaultTalents = "2350550310033311251-50003"

var PlayerOptions = &proto.Player_Mage{
	Mage: &proto.Mage{
		Options: &proto.Mage_Options{
			Armor: proto.Mage_Options_MageArmor,
		},
	},
}

var DefaultConsumes = core.ConsumesCombo{
	Label: "Default-Consumes",
	Consumes: &proto.Consumes{
		Alcohol:                    proto.Alcohol_AlcoholMedivhMerlotBlue,
		BlastedLandsBuff:           proto.BlastedLandsBuff_CerebralCortexCompound,
		DefaultConjured:            proto.Conjured_ConjuredDemonicRune,
		DefaultPotion:              proto.Potions_QuicknessPotion,
		DreamshardElixir:           true,
		Dreamtonic:                 true,
		ElixirOfGreaterArcanePower: true,
		FirePowerBuff:              proto.FirePowerBuff_ElixirOfGreaterFirepower,
		Flask:                      proto.Flask_FlaskOfSupremePower,
		Food:                       proto.Food_FoodDanonzosTelAbimMedley,
		FrostPowerBuff:             proto.FrostPowerBuff_ElixirOfGreaterFrostPower,
		HealthElixir:               proto.HealthElixir_ElixirOfFortitude,
		MagebloodPotion:            true,
		MainHandImbue:              proto.WeaponImbue_BrilliantWizardOil,
		MiscConsumes:               &proto.MiscConsumes{JujuFlurry: true},
		NordanaarHerbalTea:         true,
		SapperExplosive:            proto.SapperExplosive_SapperGoblinSapper,
		SpellPowerBuff:             proto.SpellPowerBuff_GreaterArcaneElixir,
		ZanzaBuff:                  proto.ZanzaBuff_SpiritOfZanza,
	},
}

var ItemFilters = core.ItemFilter{
	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeDagger,
		proto.WeaponType_WeaponTypeSword,
		proto.WeaponType_WeaponTypeOffHand,
		proto.WeaponType_WeaponTypeStaff,
	},
	ArmorType: proto.ArmorType_ArmorTypeCloth,
	RangedWeaponTypes: []proto.RangedWeaponType{
		proto.RangedWeaponType_RangedWeaponTypeWand,
	},
}

var Stats = []proto.Stat{
	proto.Stat_StatIntellect,
	proto.Stat_StatSpirit,
	proto.Stat_StatSpellPower,
	proto.Stat_StatSpellDamage,
	proto.Stat_StatArcanePower,
	proto.Stat_StatFirePower,
	proto.Stat_StatFrostPower,
	proto.Stat_StatSpellHit,
	proto.Stat_StatSpellCrit,
	proto.Stat_StatSpellHaste,
	proto.Stat_StatMP5,
	proto.Stat_StatFortune,
}

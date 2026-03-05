package enhancement

import (
	"testing"

	_ "github.com/isfir/wowsims-turtle/sim/common" // imported to get item effects included.
	"github.com/isfir/wowsims-turtle/sim/core"
	"github.com/isfir/wowsims-turtle/sim/core/proto"
)

func init() {
	RegisterEnhancementShaman()
}

func TestEnhancement(t *testing.T) {
	core.RunTestSuite(t, t.Name(), core.FullCharacterTestSuiteGenerator([]core.CharacterSuiteConfig{
		{
			Class:      proto.Class_ClassShaman,
			Phase:      1,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceOrc},

			Talents:     DefaultTalents,
			GearSet:     core.GetGearSet("../../../ui/enhancement_shaman/gear_sets", "phase_1"),
			Rotation:    core.GetAplRotation("../../../ui/enhancement_shaman/apls", "default"),
			Buffs:       core.FullBuffs,
			Consumes:    Phase1Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Sync Auto", SpecOptions: PlayerOptionsSyncAuto},
			OtherSpecOptions: []core.SpecOptionsCombo{
				{Label: "Sync Delay OH", SpecOptions: PlayerOptionsSyncDelayOH},
			},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassShaman,
			Phase:      2,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceOrc},

			Talents:     DefaultTalents,
			GearSet:     core.GetGearSet("../../../ui/enhancement_shaman/gear_sets", "phase_2"),
			Rotation:    core.GetAplRotation("../../../ui/enhancement_shaman/apls", "default"),
			Buffs:       core.FullBuffs,
			Consumes:    Phase1Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Sync Auto", SpecOptions: PlayerOptionsSyncAuto},
			OtherSpecOptions: []core.SpecOptionsCombo{
				{Label: "Sync Delay OH", SpecOptions: PlayerOptionsSyncDelayOH},
			},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassShaman,
			Phase:      3,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceOrc},

			Talents:     DefaultTalents,
			GearSet:     core.GetGearSet("../../../ui/enhancement_shaman/gear_sets", "phase_3"),
			Rotation:    core.GetAplRotation("../../../ui/enhancement_shaman/apls", "default"),
			Buffs:       core.FullBuffs,
			Consumes:    Phase1Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Sync Auto", SpecOptions: PlayerOptionsSyncAuto},
			OtherSpecOptions: []core.SpecOptionsCombo{
				{Label: "Sync Delay OH", SpecOptions: PlayerOptionsSyncDelayOH},
			},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
		{
			Class:      proto.Class_ClassShaman,
			Phase:      5,
			Race:       proto.Race_RaceTroll,
			OtherRaces: []proto.Race{proto.Race_RaceOrc},

			Talents:     DefaultTalents,
			GearSet:     core.GetGearSet("../../../ui/enhancement_shaman/gear_sets", "phase_5"),
			Rotation:    core.GetAplRotation("../../../ui/enhancement_shaman/apls", "default"),
			Buffs:       core.FullBuffs,
			Consumes:    Phase1Consumes,
			SpecOptions: core.SpecOptionsCombo{Label: "Sync Auto", SpecOptions: PlayerOptionsSyncAuto},
			OtherSpecOptions: []core.SpecOptionsCombo{
				{Label: "Sync Delay OH", SpecOptions: PlayerOptionsSyncDelayOH},
			},

			ItemFilter:      ItemFilters,
			EPReferenceStat: proto.Stat_StatAttackPower,
			StatsToWeigh:    Stats,
		},
	}))
}

var DefaultTalents = "05-5025002105023051-05105301"

var PlayerOptionsSyncDelayOH = &proto.Player_EnhancementShaman{
	EnhancementShaman: &proto.EnhancementShaman{
		Options: optionsSyncDelayOffhand,
	},
}

var PlayerOptionsSyncAuto = &proto.Player_EnhancementShaman{
	EnhancementShaman: &proto.EnhancementShaman{
		Options: optionsSyncAuto,
	},
}

var optionsSyncDelayOffhand = &proto.EnhancementShaman_Options{
	SyncType: proto.ShamanSyncType_DelayOffhandSwings,
}

var optionsSyncAuto = &proto.EnhancementShaman_Options{
	SyncType: proto.ShamanSyncType_Auto,
}

var Phase1Consumes = core.ConsumesCombo{
	Label: "P1-Consumes",
	Consumes: &proto.Consumes{
		AttackPowerBuff: proto.AttackPowerBuff_JujuMight,
		AgilityElixir:   proto.AgilityElixir_ElixirOfTheMongoose,
		DefaultConjured: proto.Conjured_ConjuredDemonicRune,
		DefaultPotion:   proto.Potions_MajorManaPotion,
		FirePowerBuff:   proto.FirePowerBuff_ElixirOfGreaterFirepower,
		Flask:           proto.Flask_FlaskOfSupremePower,
		Food:            proto.Food_FoodBlessSunfruit,
		MainHandImbue:   proto.WeaponImbue_WindfuryWeapon,
		SpellPowerBuff:  proto.SpellPowerBuff_GreaterArcaneElixir,
		StrengthBuff:    proto.StrengthBuff_JujuPower,
	},
}

var ItemFilters = core.ItemFilter{
	WeaponTypes: []proto.WeaponType{
		proto.WeaponType_WeaponTypeAxe,
		proto.WeaponType_WeaponTypeDagger,
		proto.WeaponType_WeaponTypeFist,
		proto.WeaponType_WeaponTypeMace,
		proto.WeaponType_WeaponTypeOffHand,
		proto.WeaponType_WeaponTypeShield,
		proto.WeaponType_WeaponTypeStaff,
	},
	ArmorType: proto.ArmorType_ArmorTypeMail,
	RangedWeaponTypes: []proto.RangedWeaponType{
		proto.RangedWeaponType_RangedWeaponTypeTotem,
	},
}

var Stats = []proto.Stat{
	proto.Stat_StatStrength,
	proto.Stat_StatAgility,
	proto.Stat_StatAttackPower,
	proto.Stat_StatMeleeHit,
	proto.Stat_StatMeleeCrit,
	proto.Stat_StatSpellPower,
}

package database

import (
	"encoding/json"
	"strings"

	"github.com/wowsims/classic/sim/core/proto"
)

// TurtleItemJSON based on TypeScript interface (from items.json)
type TurtleItemJSON struct {
	ID               int32              `json:"id"`
	Name             string             `json:"name"`
	Level            int32              `json:"level"`
	RequiredLevel    int32              `json:"requiredLevel"`
	ParentURL        string             `json:"parentUrl"`
	URL              string             `json:"url"`
	Icon             string             `json:"icon"`
	Quality          int32              `json:"quality"`
	ArmorType        int32              `json:"armorType"`  // ArmorSubclass enum
	WeaponType       int32              `json:"weaponType"` // WeaponSubclass enum
	Slot             string             `json:"slot"`
	BOE              bool               `json:"boe"`
	Unique           bool               `json:"unique"`
	Armor            int32              `json:"armor"`
	Durability       int32              `json:"durability"`
	Strength         int32              `json:"strength"`
	Agility          int32              `json:"agility"`
	Stamina          int32              `json:"stamina"`
	Intellect        int32              `json:"intellect"`
	Spirit           int32              `json:"spirit"`
	FireRes          int32              `json:"fireRes"`
	NatureRes        int32              `json:"natureRes"`
	FrostRes         int32              `json:"frostRes"`
	ShadowRes        int32              `json:"shadowRes"`
	ArcaneRes        int32              `json:"arcaneRes"`
	Effects          TurtleEffects      `json:"effects"`
	SetName          string             `json:"setName"`
	SetBonuses       json.RawMessage    `json:"setBonuses"` // Record<number, Effects>
	SetItems         []string           `json:"setItems"`
	Location         string             `json:"location"`
	ObtainType       string             `json:"obtainType"`
	ObtainFrom       string             `json:"obtainFrom"`
	SeeAlsoNum       int32              `json:"seeAlsoNum"`
	PvP              bool               `json:"pvp"`
	DropChance       *float64           `json:"dropChance,omitempty"`
	AllowableClasses int32              `json:"allowableClasses"`
	WeaponStats      *TurtleWeaponStats `json:"weaponStats,omitempty"`
	Score            *float64           `json:"score,omitempty"`
	OnUseScore       *float64           `json:"onUseScore,omitempty"`
	SetBonusScore    *float64           `json:"setBonusScore,omitempty"`
}

type TurtleEffects struct {
	Mana                   *float64 `json:"mana,omitempty"`
	HP                     *float64 `json:"hp,omitempty"`
	Intellect              *float64 `json:"intellect,omitempty"`
	Stamina                *float64 `json:"stamina,omitempty"`
	Spirit                 *float64 `json:"spirit,omitempty"`
	Agility                *float64 `json:"agility,omitempty"`
	Strength               *float64 `json:"strength,omitempty"`
	Armor                  *float64 `json:"armor,omitempty"`
	Dodge                  *float64 `json:"dodge,omitempty"`
	ThreatReductionPercent *float64 `json:"threatReductionPercent,omitempty"`
	MovementSpeedPercent   *float64 `json:"movementSpeedPercent,omitempty"`
	ManaRegenWhileCasting  *float64 `json:"manaRegenWhileCasting,omitempty"`
	MP5                    float64  `json:"mp5"`
	HP5                    float64  `json:"hp5"`
	HealPower              float64  `json:"healPower"`
	SpellCrit              float64  `json:"spellCrit"`
	SpellHit               float64  `json:"spellHit"`
	SpellPen               float64  `json:"spellPen"`
	SpellHaste             float64  `json:"spellHaste"`
	SpellPower             float64  `json:"spellPower"`
	SpellVamp              float64  `json:"spellVamp"`
	FrostSpellPower        float64  `json:"frostSpellPower"`
	FireSpellPower         float64  `json:"fireSpellPower"`
	ArcaneSpellPower       float64  `json:"arcaneSpellPower"`
	NatureSpellPower       float64  `json:"natureSpellPower"`
	ShadowSpellPower       float64  `json:"shadowSpellPower"`
	HolySpellPower         float64  `json:"holySpellPower"`
	TargetTypes            int      `json:"targetTypes"`
	OnUse                  bool     `json:"onUse"`
	Custom                 []string `json:"custom"`
}

type TurtleWeaponStats struct {
	MinDmg float64 `json:"minDmg"`
	MaxDmg float64 `json:"maxDmg"`
	Speed  float64 `json:"speed"`
	DPS    float64 `json:"dps"`
}

// ParseTurtleGearplannerDB reads raw Turtle WoW gear planner JSON and returns a WowDatabase.
func ParseTurtleGearplannerDB(jsonData string) *WowDatabase {
	db := NewWowDatabase()

	var items []TurtleItemJSON
	if err := json.Unmarshal([]byte(jsonData), &items); err != nil {
		panic("Failed to unmarshal Turtle gearplanner JSON: " + err.Error())
	}

	for i := range items {
		item := convertTurtleItem(&items[i])
		if item != nil {
			db.MergeItem(item)
		}
	}

	return db
}

func convertTurtleItem(turtle *TurtleItemJSON) *proto.UIItem {
	icon := turtle.Icon
	// Strip .jpg extension if present (UI will add it back)
	if strings.HasSuffix(icon, ".jpg") {
		icon = strings.TrimSuffix(icon, ".jpg")
	}

	item := &proto.UIItem{
		Id:      turtle.ID,
		Name:    turtle.Name,
		Icon:    icon,
		Ilvl:    turtle.Level, // TODO: verify mapping (level vs ilvl)
		Quality: proto.ItemQuality(turtle.Quality),
		Unique:  turtle.Unique,
	}

	// Map slot to ItemType
	item.Type = mapSlotToItemType(turtle.Slot)
	// Map armorType to ArmorType (only for armor items)
	item.ArmorType = mapArmorSubclassToArmorType(turtle.ArmorType)
	// Map weaponType to WeaponType and RangedWeaponType
	item.WeaponType = mapWeaponSubclassToWeaponType(turtle.WeaponType)
	item.RangedWeaponType = mapWeaponSubclassToRangedWeaponType(turtle.WeaponType)
	// Map hand type based on slot
	item.HandType = mapSlotToHandType(turtle.Slot, turtle.WeaponType)

	// Stats array (size = number of Stat enum values)
	stats := make([]float64, 44) // Stat enum size
	stats[proto.Stat_StatStrength] = float64(turtle.Strength)
	stats[proto.Stat_StatAgility] = float64(turtle.Agility)
	stats[proto.Stat_StatStamina] = float64(turtle.Stamina)
	stats[proto.Stat_StatIntellect] = float64(turtle.Intellect)
	stats[proto.Stat_StatSpirit] = float64(turtle.Spirit)
	stats[proto.Stat_StatFireResistance] = float64(turtle.FireRes)
	stats[proto.Stat_StatNatureResistance] = float64(turtle.NatureRes)
	stats[proto.Stat_StatFrostResistance] = float64(turtle.FrostRes)
	stats[proto.Stat_StatShadowResistance] = float64(turtle.ShadowRes)
	stats[proto.Stat_StatArcaneResistance] = float64(turtle.ArcaneRes)
	stats[proto.Stat_StatArmor] = float64(turtle.Armor)

	// Effects stats
	stats[proto.Stat_StatMP5] = turtle.Effects.MP5
	stats[proto.Stat_StatSpellHit] = turtle.Effects.SpellHit
	stats[proto.Stat_StatSpellCrit] = turtle.Effects.SpellCrit
	stats[proto.Stat_StatSpellHaste] = turtle.Effects.SpellHaste
	stats[proto.Stat_StatSpellPenetration] = turtle.Effects.SpellPen
	stats[proto.Stat_StatSpellPower] = turtle.Effects.SpellPower
	stats[proto.Stat_StatArcanePower] = turtle.Effects.ArcaneSpellPower
	stats[proto.Stat_StatFirePower] = turtle.Effects.FireSpellPower
	stats[proto.Stat_StatFrostPower] = turtle.Effects.FrostSpellPower
	stats[proto.Stat_StatHolyPower] = turtle.Effects.HolySpellPower
	stats[proto.Stat_StatNaturePower] = turtle.Effects.NatureSpellPower
	stats[proto.Stat_StatShadowPower] = turtle.Effects.ShadowSpellPower
	stats[proto.Stat_StatHealingPower] = turtle.Effects.HealPower
	if turtle.Effects.Dodge != nil {
		stats[proto.Stat_StatDodge] = *turtle.Effects.Dodge
	}
	// Additional stat pointers
	if turtle.Effects.Mana != nil {
		stats[proto.Stat_StatMana] = *turtle.Effects.Mana
	}
	if turtle.Effects.HP != nil {
		stats[proto.Stat_StatHealth] = *turtle.Effects.HP
	}
	if turtle.Effects.Armor != nil {
		// Assume bonus armor
		stats[proto.Stat_StatBonusArmor] = *turtle.Effects.Armor
	}
	// TODO: threat reduction, movement speed, mana regen while casting -> pseudo stats

	item.Stats = stats

	// Weapon stats
	if turtle.WeaponStats != nil {
		item.WeaponDamageMin = turtle.WeaponStats.MinDmg
		item.WeaponDamageMax = turtle.WeaponStats.MaxDmg
		item.WeaponSpeed = turtle.WeaponStats.Speed
		// Compute weapon skills based on weapon type
		// item.WeaponSkills = getWeaponSkills(turtle.WeaponType) // TODO: map to skill bonuses
	}

	// Class restrictions - use allowableClasses bitfield directly
	// Bitmask: 2=Warrior, 4=Paladin, 8=Hunter, 16=Rogue, 32=Priest, 64=Shaman, 128=Mage, 256=Warlock, 512=Druid
	// 1023 = all classes (bits 1-9 set: 2+4+8+16+32+64+128+256+512 = 1022? Actually 1023 includes bit 0)
	if turtle.AllowableClasses != 1023 { // 1023 = all classes
		item.ClassAllowlist = mapAllowableClasses(turtle.AllowableClasses)
	}

	// TODO: Parse sources from location/obtainType/obtainFrom
	// TODO: Set name and set bonuses
	// TODO: Faction restriction (maybe from location)
	// TODO: Required profession

	return item
}

func mapSlotToItemType(slot string) proto.ItemType {
	switch strings.ToLower(slot) {
	case "head":
		return proto.ItemType_ItemTypeHead
	case "neck":
		return proto.ItemType_ItemTypeNeck
	case "shoulder":
		return proto.ItemType_ItemTypeShoulder
	case "back":
		return proto.ItemType_ItemTypeBack
	case "chest":
		return proto.ItemType_ItemTypeChest
	case "wrist":
		return proto.ItemType_ItemTypeWrist
	case "hands":
		return proto.ItemType_ItemTypeHands
	case "waist":
		return proto.ItemType_ItemTypeWaist
	case "legs":
		return proto.ItemType_ItemTypeLegs
	case "feet":
		return proto.ItemType_ItemTypeFeet
	case "finger":
		return proto.ItemType_ItemTypeFinger
	case "trinket":
		return proto.ItemType_ItemTypeTrinket
	case "mainhand", "offhand", "twohand":
		return proto.ItemType_ItemTypeWeapon
	case "ranged":
		return proto.ItemType_ItemTypeRanged
	default:
		return proto.ItemType_ItemTypeUnknown
	}
}

func mapArmorSubclassToArmorType(armorSubclass int32) proto.ArmorType {
	switch armorSubclass {
	case 1: // Cloth
		return proto.ArmorType_ArmorTypeCloth
	case 2: // Leather
		return proto.ArmorType_ArmorTypeLeather
	case 3: // Mail
		return proto.ArmorType_ArmorTypeMail
	case 4: // Plate
		return proto.ArmorType_ArmorTypePlate
	case 6: // Shield (treated as weapon type)
		return proto.ArmorType_ArmorTypeUnknown
	default:
		return proto.ArmorType_ArmorTypeUnknown
	}
}

func mapWeaponSubclassToWeaponType(weaponSubclass int32) proto.WeaponType {
	switch weaponSubclass {
	case -100: // Empty
		return proto.WeaponType_WeaponTypeUnknown
	case 0: // OneHandedAxe
		return proto.WeaponType_WeaponTypeAxe
	case 1: // TwoHandedAxe
		return proto.WeaponType_WeaponTypeAxe
	case 4: // OneHandedMace
		return proto.WeaponType_WeaponTypeMace
	case 5: // TwoHandedMace
		return proto.WeaponType_WeaponTypeMace
	case 6: // Polearm
		return proto.WeaponType_WeaponTypePolearm
	case 7: // OneHandedSword
		return proto.WeaponType_WeaponTypeSword
	case 8: // TwoHandedSword
		return proto.WeaponType_WeaponTypeSword
	case 10: // Staff
		return proto.WeaponType_WeaponTypeStaff
	case 13: // FistWeapon
		return proto.WeaponType_WeaponTypeFist
	case 15: // Dagger
		return proto.WeaponType_WeaponTypeDagger
	case 11: // OneHandedExotic (maybe fist?)
		return proto.WeaponType_WeaponTypeFist
	case 12: // TwoHandedExotic (maybe polearm?)
		return proto.WeaponType_WeaponTypePolearm
	case 14: // Miscellaneous
		return proto.WeaponType_WeaponTypeUnknown
	case 16: // Thrown (ranged)
		return proto.WeaponType_WeaponTypeUnknown
	case 17: // Spear (polearm?)
		return proto.WeaponType_WeaponTypePolearm
	case 18: // Crossbow (ranged)
		return proto.WeaponType_WeaponTypeUnknown
	case 19: // Wand (ranged)
		return proto.WeaponType_WeaponTypeUnknown
	case 20: // FishingPole
		return proto.WeaponType_WeaponTypeUnknown
	case 2: // Bow (ranged)
		return proto.WeaponType_WeaponTypeUnknown
	case 3: // Gun (ranged)
		return proto.WeaponType_WeaponTypeUnknown
	default:
		return proto.WeaponType_WeaponTypeUnknown
	}
}

func mapWeaponSubclassToRangedWeaponType(weaponSubclass int32) proto.RangedWeaponType {
	switch weaponSubclass {
	case 2: // Bow
		return proto.RangedWeaponType_RangedWeaponTypeBow
	case 3: // Gun
		return proto.RangedWeaponType_RangedWeaponTypeGun
	case 18: // Crossbow
		return proto.RangedWeaponType_RangedWeaponTypeCrossbow
	case 16: // Thrown
		return proto.RangedWeaponType_RangedWeaponTypeThrown
	case 19: // Wand
		return proto.RangedWeaponType_RangedWeaponTypeWand
	case 8: // Libram (armor type 7)
		return proto.RangedWeaponType_RangedWeaponTypeLibram
	case 9: // Idol (armor type 8)
		return proto.RangedWeaponType_RangedWeaponTypeIdol
	case 10: // Totem (armor type 9)
		return proto.RangedWeaponType_RangedWeaponTypeTotem
	default:
		return proto.RangedWeaponType_RangedWeaponTypeUnknown
	}
}

func mapSlotToHandType(slot string, weaponSubclass int32) proto.HandType {
	switch strings.ToLower(slot) {
	case "mainhand":
		return proto.HandType_HandTypeMainHand
	case "offhand":
		// Check if it's a shield (armorType 6)
		if weaponSubclass == -100 {
			// Might be shield (armorType 6) or held in offhand
			return proto.HandType_HandTypeOffHand
		}
		return proto.HandType_HandTypeOffHand
	case "twohand":
		return proto.HandType_HandTypeTwoHand
	default:
		return proto.HandType_HandTypeUnknown
	}
}

func getWeaponSkills(weaponSubclass int32) []proto.WeaponSkill {
	var skills []proto.WeaponSkill
	switch weaponSubclass {
	case 0: // OneHandedAxe
		skills = append(skills, proto.WeaponSkill_WeaponSkillAxes)
	case 1: // TwoHandedAxe
		skills = append(skills, proto.WeaponSkill_WeaponSkillTwoHandedAxes)
	case 4: // OneHandedMace
		skills = append(skills, proto.WeaponSkill_WeaponSkillMaces)
	case 5: // TwoHandedMace
		skills = append(skills, proto.WeaponSkill_WeaponSkillTwoHandedMaces)
	case 6: // Polearm
		skills = append(skills, proto.WeaponSkill_WeaponSkillPolearms)
	case 7: // OneHandedSword
		skills = append(skills, proto.WeaponSkill_WeaponSkillSwords)
	case 8: // TwoHandedSword
		skills = append(skills, proto.WeaponSkill_WeaponSkillTwoHandedSwords)
	case 10: // Staff
		skills = append(skills, proto.WeaponSkill_WeaponSkillStaves)
	case 13: // FistWeapon
		skills = append(skills, proto.WeaponSkill_WeaponSkillUnarmed)
	case 15: // Dagger
		skills = append(skills, proto.WeaponSkill_WeaponSkillDaggers)
	case 2: // Bow
		skills = append(skills, proto.WeaponSkill_WeaponSkillBows)
	case 3: // Gun
		skills = append(skills, proto.WeaponSkill_WeaponSkillGuns)
	case 18: // Crossbow
		skills = append(skills, proto.WeaponSkill_WeaponSkillCrossbows)
	case 16: // Thrown
		skills = append(skills, proto.WeaponSkill_WeaponSkillThrown)
	default:
		// No weapon skills
	}
	return skills
}

func mapAllowableClasses(allowableClasses int32) []proto.Class {
	// Bitmask matches TypeScript parser: 2=Warrior, 4=Paladin, 8=Hunter, 16=Rogue, 32=Priest, 64=Shaman, 128=Mage, 256=Warlock, 512=Druid
	var classes []proto.Class
	if allowableClasses&2 != 0 { // Bit 1: Warrior
		classes = append(classes, proto.Class_ClassWarrior)
	}
	if allowableClasses&4 != 0 { // Bit 2: Paladin
		classes = append(classes, proto.Class_ClassPaladin)
	}
	if allowableClasses&8 != 0 { // Bit 3: Hunter
		classes = append(classes, proto.Class_ClassHunter)
	}
	if allowableClasses&16 != 0 { // Bit 4: Rogue
		classes = append(classes, proto.Class_ClassRogue)
	}
	if allowableClasses&32 != 0 { // Bit 5: Priest
		classes = append(classes, proto.Class_ClassPriest)
	}
	if allowableClasses&64 != 0 { // Bit 6: Shaman
		classes = append(classes, proto.Class_ClassShaman)
	}
	if allowableClasses&128 != 0 { // Bit 7: Mage
		classes = append(classes, proto.Class_ClassMage)
	}
	if allowableClasses&256 != 0 { // Bit 8: Warlock
		classes = append(classes, proto.Class_ClassWarlock)
	}
	if allowableClasses&512 != 0 { // Bit 9: Druid
		classes = append(classes, proto.Class_ClassDruid)
	}
	return classes
}

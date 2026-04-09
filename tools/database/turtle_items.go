package database

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/isfir/wowsims-turtle/sim/core/proto"
)

type itemSetInfo struct {
	SetID   int32
	SetName string
}

type SpellEffect struct {
	Stat                proto.Stat
	WeaponSkill         proto.WeaponSkill
	BonusPhysicalDamage float64
	Value               float64
}

type SpellAnalysis struct {
	Effects            []SpellEffect
	HasUnhandledEffect bool
}

const (
	statsLen       = int(proto.Stat_StatVampirism) + 1
	weaponSkillLen = int(proto.WeaponSkill_WeaponSkillFeralCombat) + 1
)

// Inventory types
const (
	invTypeNonEquip    = 0
	invTypeHead        = 1
	invTypeNeck        = 2
	invTypeShoulders   = 3
	invTypeBody        = 4
	invTypeChest       = 5
	invTypeWaist       = 6
	invTypeLegs        = 7
	invTypeFeet        = 8
	invTypeWrists      = 9
	invTypeHands       = 10
	invTypeFinger      = 11
	invTypeTrinket     = 12
	invTypeWeapon      = 13
	invTypeShield      = 14
	invTypeRanged      = 15
	invTypeCloak       = 16
	invType2HWeapon    = 17
	invTypeBag         = 18
	invTypeTabard      = 19
	invTypeRobe        = 20
	invTypeWeaponMain  = 21
	invTypeWeaponOff   = 22
	invTypeHoldable    = 23
	invTypeAmmo        = 24
	invTypeThrown      = 25
	invTypeRangedRight = 26
	invTypeQuiver      = 27
	invTypeRelic       = 28
)

const (
	spellEffectApplyAura = 6
)

// Aura types (only ones we care about)
const (
	auraModDamageDone                 = 13
	auraModStat                       = 29
	auraModResistance                 = 22
	auraModBaseResistance             = 83
	auraModSpellHitChance             = 55
	auraModSpellCritChance            = 57
	auraModSpellCritChanceSchool      = 71
	auraModMeleeHaste                 = 138
	auraModRangedHaste                = 140
	auraModCastingSpeedNotStack       = 65
	auraModAttackPower                = 99
	auraModRangedAttackPower          = 124
	auraModRangedAttackPowerVersus    = 131
	auraModIncreaseHealth             = 34
	auraModIncreaseEnergy             = 35
	auraModPowerRegen                 = 85
	auraModManaRegenInterrupt         = 134
	auraModCritPercent                = 52
	auraModHitChance                  = 54
	auraModHealingDone                = 135
	auraModHealing                    = 115
	auraModSkill                      = 30
	auraModSkillTalent                = 98
	auraModParrySkill                 = 46
	auraModDodgeSkill                 = 48
	auraModBlockSkill                 = 50
	auraModParryPercent               = 47
	auraModDodgePercent               = 49
	auraModBlockPercent               = 51
	auraModShieldBlockValue           = 158
	auraModAttackerSpellHitChance     = 186
	auraModAttackerMeleeHitChance     = 184
	auraModAttackerRangedHitChance    = 185
	auraModTargetResistance           = 123
	auraModFortune                    = 223
	auraModPeriodicDamagePercentTaken = 197
	auraModCritDamageTaken            = 198
)

const (
	powerTypeMana   = 0
	powerTypeRage   = 1
	powerTypeFocus  = 2
	powerTypeEnergy = 3
)

const (
	spellSchoolMaskNormal = 1
	spellSchoolMaskSpell  = 124
	spellSchoolMaskMagic  = 126
	spellSchoolMaskAll    = 127
)

var classMaskMap = []struct {
	ClassID int32
	Class   proto.Class
}{
	{1, proto.Class_ClassWarrior},
	{2, proto.Class_ClassPaladin},
	{3, proto.Class_ClassHunter},
	{4, proto.Class_ClassRogue},
	{5, proto.Class_ClassPriest},
	{7, proto.Class_ClassShaman},
	{8, proto.Class_ClassMage},
	{9, proto.Class_ClassWarlock},
	{11, proto.Class_ClassDruid},
}

var allClassesMask = func() uint32 {
	var mask uint32
	for _, c := range classMaskMap {
		mask |= 1 << (c.ClassID - 1)
	}
	return mask
}()

var allianceRaceMask = func() uint32 {
	var mask uint32
	mask |= 1 << (1 - 1) // Human
	mask |= 1 << (3 - 1) // Dwarf
	mask |= 1 << (4 - 1) // Night Elf
	mask |= 1 << (7 - 1) // Gnome
	return mask
}()

var hordeRaceMask = func() uint32 {
	var mask uint32
	mask |= 1 << (2 - 1) // Orc
	mask |= 1 << (5 - 1) // Undead
	mask |= 1 << (6 - 1) // Tauren
	mask |= 1 << (8 - 1) // Troll
	return mask
}()

func isAllZeros(slice []float64) bool {
	if slice == nil {
		return true
	}
	for _, v := range slice {
		if v != 0 {
			return false
		}
	}
	return true
}

func ParseTurtleItemsDB(itemsCSV, itemDisplayInfoCSV, spellCSV, factionCSV, itemSetCSV string) *WowDatabase {
	db := NewWowDatabase()

	displayInfo := parseItemDisplayInfoCSV(itemDisplayInfoCSV)
	spellEffects := parseSpellEffectsCSV(spellCSV)
	factionMap := parseFactionMapping(factionCSV)
	itemSets, itemToSet := parseItemSetCSV(itemSetCSV)
	for _, itemSet := range itemSets {
		db.MergeItemSet(itemSet)
	}

	items := parseItemsCSV(itemsCSV, displayInfo, spellEffects, factionMap, itemToSet)
	for _, item := range items {
		db.MergeItem(item)
	}

	// TODO: random suffixes (ItemRandomProperties + SpellItemEnchantment)

	return db
}

func parseItemDisplayInfoCSV(csvData string) map[int32]string {
	r := csv.NewReader(strings.NewReader(csvData))
	if _, err := r.Read(); err != nil {
		log.Fatalf("Cannot read item display info csv header: %v", err)
	}

	iconMap := make(map[int32]string)
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Cannot read item display info csv row: %v", err)
		}
		if len(row) < 6 {
			continue
		}
		id, err := strconv.Atoi(row[0])
		if err != nil {
			continue
		}
		icon := row[5]
		if icon != "" {
			iconMap[int32(id)] = strings.ToLower(icon)
		}
	}
	return iconMap
}

func parseItemsCSV(csvData string, displayInfo map[int32]string, spellEffects map[int32]SpellAnalysis, factionMap map[int32]proto.UIItem_FactionRestriction, itemToSet map[int32]itemSetInfo) []*proto.UIItem {
	r := csv.NewReader(strings.NewReader(csvData))
	headers, err := r.Read()
	if err != nil {
		log.Fatalf("Cannot read items csv header: %v", err)
	}

	colIdx := make(map[string]int)
	for i, name := range headers {
		colIdx[name] = i
	}

	requiredCols := []string{
		"itemID", "itemDisplayID", "qualityID", "itemLevel",
		"inventorySlotID", "itemClassID", "itemSubClassID", "name1",
	}
	for _, col := range requiredCols {
		if _, ok := colIdx[col]; !ok {
			log.Fatalf("Missing required column %s in items csv", col)
		}
	}

	var items []*proto.UIItem
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Cannot read items csv row: %v", err)
		}

		item := parseItemRow(row, colIdx, displayInfo, spellEffects, factionMap, itemToSet)
		if item != nil {
			items = append(items, item)
		}
	}
	return items
}

func parseItemRow(row []string, colIdx map[string]int, displayInfo map[int32]string, spellEffects map[int32]SpellAnalysis, factionMap map[int32]proto.UIItem_FactionRestriction, itemToSet map[int32]itemSetInfo) *proto.UIItem {
	id, ok := getInt(row, colIdx, "itemID")
	if !ok {
		return nil
	}

	displayID, _ := getInt(row, colIdx, "itemDisplayID")
	icon := displayInfo[displayID]

	name := getString(row, colIdx, "name1")
	if name == "" {
		return nil
	}

	quality, _ := getInt(row, colIdx, "qualityID")
	ilvl, _ := getInt(row, colIdx, "itemLevel")
	invType, _ := getInt(row, colIdx, "inventorySlotID")
	itemClass, _ := getInt(row, colIdx, "itemClassID")
	itemSubClass, _ := getInt(row, colIdx, "itemSubClassID")

	itemType := mapInventorySlotToItemType(invType)

	item := &proto.UIItem{
		Id:      id,
		Name:    name,
		Icon:    icon,
		Ilvl:    ilvl,
		Quality: proto.ItemQuality(quality),
		Type:    itemType,
	}

	item.ArmorType = mapArmorType(itemClass, itemSubClass)
	item.WeaponType, item.RangedWeaponType = mapWeaponTypes(itemClass, itemSubClass, invType)
	item.HandType = mapHandType(invType)

	stats := parseItemStats(row, colIdx)
	spellBonuses, weaponSkills, bonusPhysicalDmg := parseItemSpellEffects(row, colIdx, spellEffects)

	for i := range stats {
		stats[i] += spellBonuses[i]
	}
	if !isAllZeros(stats) {
		item.Stats = stats
	}
	if !isAllZeros(weaponSkills) {
		item.WeaponSkills = weaponSkills
	}
	item.BonusPhysicalDamage = bonusPhysicalDmg

	if hasWeaponDamage(row, colIdx) {
		minDmg, maxDmg, speed := parseWeaponDamage(row, colIdx)
		item.WeaponDamageMin = minDmg
		item.WeaponDamageMax = maxDmg
		item.WeaponSpeed = speed
	}

	item.Effects = parseItemEffects(row, colIdx, spellEffects)

	if classMask, ok := getUint(row, colIdx, "classBinFlag"); ok {
		item.ClassAllowlist = parseClassMask(uint32(classMask))
	}

	if raceMask, ok := getUint(row, colIdx, "raceBinFlag"); ok {
		item.FactionRestriction = factionRestrictionFromRaceMask(uint32(raceMask))
	}

	if item.FactionRestriction == proto.UIItem_FACTION_RESTRICTION_UNSPECIFIED {
		if factionID, ok := getInt(row, colIdx, "requiredFactionID"); ok && factionID != 0 {
			if restriction, ok := factionMap[factionID]; ok {
				item.FactionRestriction = restriction
			}
		}
	}

	if setID, ok := getInt(row, colIdx, "itemSetID"); ok && setID != 0 {
		item.SetId = setID
	}

	if unique, ok := getInt(row, colIdx, "stackUnique"); ok {
		item.Unique = unique > 0
	}

	if prof := mapRequiredProfession(row, colIdx); prof != proto.Profession_ProfessionUnknown {
		item.RequiredProfession = prof
	}

	if bonding, ok := getInt(row, colIdx, "bondID"); ok {
		item.BindType = mapItemBindType(bonding)
	}

	if requiredLevel, ok := getInt(row, colIdx, "requiredLevel"); ok && requiredLevel > 0 {
		item.RequiredLevel = requiredLevel
	}

	if maxDurability, ok := getInt(row, colIdx, "durabilityValue"); ok && maxDurability > 0 {
		item.MaxDurability = maxDurability
	}

	if setID, ok := getInt(row, colIdx, "itemSetID"); ok && setID != 0 {
		item.SetId = setID
	}

	if info, ok := itemToSet[item.Id]; ok {
		if item.SetId == 0 {
			item.SetId = info.SetID
		}
		if item.SetName == "" {
			item.SetName = info.SetName
		}
	}

	return item
}

func parseItemStats(row []string, colIdx map[string]int) []float64 {
	stats := make([]float64, statsLen)

	for i := 1; i <= 10; i++ {
		statIDCol := fmt.Sprintf("stat%dID", i)
		statValueCol := fmt.Sprintf("stat%dValue", i)

		statIDStr := getString(row, colIdx, statIDCol)
		statValueStr := getString(row, colIdx, statValueCol)

		if statIDStr == "" || statValueStr == "" || statIDStr == "0" {
			continue
		}

		statID, err1 := strconv.Atoi(statIDStr)
		statValue, err2 := strconv.ParseFloat(statValueStr, 64)
		if err1 != nil || err2 != nil {
			continue
		}

		protoStat := mapItemModTypeToStat(int32(statID))
		if protoStat != proto.Stat(-1) {
			stats[protoStat] = statValue
		}
	}

	resistCols := map[string]proto.Stat{
		"resistPhysical": proto.Stat_StatArmor,
		"resistFire":     proto.Stat_StatFireResistance,
		"resistNature":   proto.Stat_StatNatureResistance,
		"resistFrost":    proto.Stat_StatFrostResistance,
		"resistShadow":   proto.Stat_StatShadowResistance,
		"resistArcane":   proto.Stat_StatArcaneResistance,
	}

	for col, stat := range resistCols {
		if val, ok := getFloat(row, colIdx, col); ok {
			stats[stat] = val
		}
	}

	if blockValue, ok := getFloat(row, colIdx, "blockValue"); ok {
		stats[proto.Stat_StatBlockValue] = blockValue
	}

	return stats
}

func parseItemSpellEffects(row []string, colIdx map[string]int, spellEffects map[int32]SpellAnalysis) ([]float64, []float64, float64) {
	bonuses := make([]float64, statsLen)
	weaponSkills := make([]float64, weaponSkillLen)

	var bonusPhysicalDamage float64

	for i := 1; i <= 5; i++ {
		spellIDCol := fmt.Sprintf("spell%dID", i)
		triggerIDCol := fmt.Sprintf("spell%dTriggerID", i)

		spellID, okID := getInt(row, colIdx, spellIDCol)
		triggerID, okTrig := getInt(row, colIdx, triggerIDCol)
		if !okID || !okTrig || spellID == 0 {
			continue
		}

		if triggerID != 1 { // only on-equip contributes to passive item stats
			continue
		}

		analysis, ok := spellEffects[spellID]
		if !ok {
			continue
		}

		for _, effect := range analysis.Effects {
			switch {
			case effect.WeaponSkill != proto.WeaponSkill_WeaponSkillUnknown:
				weaponSkills[effect.WeaponSkill] += effect.Value

			case effect.BonusPhysicalDamage != 0:
				bonusPhysicalDamage += effect.BonusPhysicalDamage

			case effect.Stat != proto.Stat(-1):
				bonuses[effect.Stat] += effect.Value
			}
		}
	}

	return bonuses, weaponSkills, bonusPhysicalDamage
}

func parseItemEffects(row []string, colIdx map[string]int, spellEffects map[int32]SpellAnalysis) []*proto.UIItemEffect {
	var effects []*proto.UIItemEffect

	for i := 1; i <= 5; i++ {
		spellIDCol := fmt.Sprintf("spell%dID", i)
		triggerIDCol := fmt.Sprintf("spell%dTriggerID", i)

		spellID, okID := getInt(row, colIdx, spellIDCol)
		triggerID, okTrig := getInt(row, colIdx, triggerIDCol)
		if !okID || !okTrig || spellID == 0 {
			continue
		}

		triggerType := mapItemEffectTriggerType(triggerID)
		if triggerType == proto.ItemEffectTriggerType_ItemEffectTriggerTypeUnknown {
			continue
		}

		analysis, hasAnalysis := spellEffects[spellID]

		// Pure passive on-equip stat auras are already folded into item.stats.
		// Skip them here to avoid duplicate tooltip lines.
		if triggerType == proto.ItemEffectTriggerType_ItemEffectTriggerTypeOnEquip &&
			hasAnalysis &&
			!analysis.HasUnhandledEffect &&
			len(analysis.Effects) > 0 {
			continue
		}

		effects = append(effects, &proto.UIItemEffect{
			TriggerType: triggerType,
			SpellId:     spellID,
		})
	}

	return effects
}

func specialCaseSpellAnalysis(spellID int32) (SpellAnalysis, bool) {
	switch spellID {
	case 45420, 45421, 45422, 45423, 45424:
		return SpellAnalysis{
			Effects: []SpellEffect{
				{
					Stat:  proto.Stat_StatVampirism,
					Value: float64(spellID - 45419),
				},
			},
		}, true
	default:
		return SpellAnalysis{}, false
	}
}

func hasWeaponDamage(row []string, colIdx map[string]int) bool {
	minDmgStr := getString(row, colIdx, "damage1Min")
	return minDmgStr != "" && minDmgStr != "0.0"
}

func parseWeaponDamage(row []string, colIdx map[string]int) (minDmg, maxDmg, speed float64) {
	minDmgStr := getString(row, colIdx, "damage1Min")
	maxDmgStr := getString(row, colIdx, "damage1Max")
	speedStr := getString(row, colIdx, "weaponDelay")

	if minDmgStr != "" && minDmgStr != "0.0" {
		minDmg, _ = strconv.ParseFloat(minDmgStr, 64)
	}
	if maxDmgStr != "" && maxDmgStr != "0.0" {
		maxDmg, _ = strconv.ParseFloat(maxDmgStr, 64)
	}
	if speedStr != "" && speedStr != "0" {
		speedVal, _ := strconv.ParseFloat(speedStr, 64)
		speed = speedVal / 1000.0
	}

	return minDmg, maxDmg, speed
}

func parseClassMask(classMask uint32) []proto.Class {
	if classMask == 0 || (classMask&allClassesMask) == allClassesMask {
		return nil
	}

	var classes []proto.Class
	for _, c := range classMaskMap {
		if classMask&(1<<(c.ClassID-1)) != 0 {
			classes = append(classes, c.Class)
		}
	}
	return classes
}

func factionRestrictionFromRaceMask(mask uint32) proto.UIItem_FactionRestriction {
	if mask == 0 || mask == 0xFFFFFFFF {
		return proto.UIItem_FACTION_RESTRICTION_UNSPECIFIED
	}

	hasAlliance := (mask & allianceRaceMask) != 0
	hasHorde := (mask & hordeRaceMask) != 0

	switch {
	case hasAlliance && !hasHorde:
		return proto.UIItem_FACTION_RESTRICTION_ALLIANCE_ONLY
	case hasHorde && !hasAlliance:
		return proto.UIItem_FACTION_RESTRICTION_HORDE_ONLY
	default:
		return proto.UIItem_FACTION_RESTRICTION_UNSPECIFIED
	}
}

func targetResistanceMaskToEffects(mask int32, value float64) ([]SpellEffect, bool) {
	value = math.Abs(value)

	switch mask {
	case spellSchoolMaskNormal:
		return []SpellEffect{
			{Stat: proto.Stat_StatArmorPenetration, Value: value},
		}, true

	case spellSchoolMaskSpell, spellSchoolMaskMagic:
		return []SpellEffect{
			{Stat: proto.Stat_StatSpellPenetration, Value: value},
		}, true

	case spellSchoolMaskAll:
		return []SpellEffect{
			{Stat: proto.Stat_StatArmorPenetration, Value: value},
			{Stat: proto.Stat_StatSpellPenetration, Value: value},
		}, true

	default:
		return nil, false
	}
}

func mapInventorySlotToItemType(slot int32) proto.ItemType {
	switch slot {
	case invTypeHead:
		return proto.ItemType_ItemTypeHead
	case invTypeNeck:
		return proto.ItemType_ItemTypeNeck
	case invTypeShoulders:
		return proto.ItemType_ItemTypeShoulder
	case invTypeChest, invTypeRobe:
		return proto.ItemType_ItemTypeChest
	case invTypeWaist:
		return proto.ItemType_ItemTypeWaist
	case invTypeLegs:
		return proto.ItemType_ItemTypeLegs
	case invTypeFeet:
		return proto.ItemType_ItemTypeFeet
	case invTypeWrists:
		return proto.ItemType_ItemTypeWrist
	case invTypeHands:
		return proto.ItemType_ItemTypeHands
	case invTypeFinger:
		return proto.ItemType_ItemTypeFinger
	case invTypeTrinket:
		return proto.ItemType_ItemTypeTrinket
	case invTypeCloak:
		return proto.ItemType_ItemTypeBack
	case invTypeWeapon, invTypeWeaponMain, invTypeWeaponOff, invType2HWeapon, invTypeShield, invTypeHoldable:
		return proto.ItemType_ItemTypeWeapon
	case invTypeRanged, invTypeRangedRight, invTypeThrown, invTypeRelic:
		return proto.ItemType_ItemTypeRanged
	default:
		return proto.ItemType_ItemTypeUnknown
	}
}

func mapArmorType(itemClass, itemSubClass int32) proto.ArmorType {
	if itemClass != 4 {
		return proto.ArmorType_ArmorTypeUnknown
	}
	switch itemSubClass {
	case 1:
		return proto.ArmorType_ArmorTypeCloth
	case 2:
		return proto.ArmorType_ArmorTypeLeather
	case 3:
		return proto.ArmorType_ArmorTypeMail
	case 4:
		return proto.ArmorType_ArmorTypePlate
	default:
		return proto.ArmorType_ArmorTypeUnknown
	}
}

func mapWeaponTypes(itemClass, itemSubClass, invType int32) (proto.WeaponType, proto.RangedWeaponType) {
	var weaponType proto.WeaponType
	var rangedType proto.RangedWeaponType

	if itemClass == 2 {
		switch itemSubClass {
		case 0, 1:
			weaponType = proto.WeaponType_WeaponTypeAxe
		case 4, 5:
			weaponType = proto.WeaponType_WeaponTypeMace
		case 7, 8:
			weaponType = proto.WeaponType_WeaponTypeSword
		case 6:
			weaponType = proto.WeaponType_WeaponTypePolearm
		case 10:
			weaponType = proto.WeaponType_WeaponTypeStaff
		case 13:
			weaponType = proto.WeaponType_WeaponTypeFist
		case 15:
			weaponType = proto.WeaponType_WeaponTypeDagger
		case 17:
			weaponType = proto.WeaponType_WeaponTypePolearm
		case 20:
			weaponType = proto.WeaponType_WeaponTypeStaff // fishing pole
		}

		switch itemSubClass {
		case 2:
			rangedType = proto.RangedWeaponType_RangedWeaponTypeBow
		case 3:
			rangedType = proto.RangedWeaponType_RangedWeaponTypeGun
		case 18:
			rangedType = proto.RangedWeaponType_RangedWeaponTypeCrossbow
		case 19:
			rangedType = proto.RangedWeaponType_RangedWeaponTypeWand
		case 16:
			rangedType = proto.RangedWeaponType_RangedWeaponTypeThrown
		}
	}

	if itemClass == 4 {
		switch itemSubClass {
		case 6:
			weaponType = proto.WeaponType_WeaponTypeShield
		case 7:
			rangedType = proto.RangedWeaponType_RangedWeaponTypeLibram
		case 8:
			rangedType = proto.RangedWeaponType_RangedWeaponTypeIdol
		case 9:
			rangedType = proto.RangedWeaponType_RangedWeaponTypeTotem
		}
	}

	if invType == invTypeHoldable {
		weaponType = proto.WeaponType_WeaponTypeOffHand
	}

	return weaponType, rangedType
}

func mapHandType(invType int32) proto.HandType {
	switch invType {
	case invTypeWeapon:
		return proto.HandType_HandTypeOneHand
	case invTypeWeaponMain:
		return proto.HandType_HandTypeMainHand
	case invTypeWeaponOff, invTypeShield, invTypeHoldable:
		return proto.HandType_HandTypeOffHand
	case invType2HWeapon:
		return proto.HandType_HandTypeTwoHand
	default:
		return proto.HandType_HandTypeUnknown
	}
}

func mapProfessionBySkillID(skillID int32) proto.Profession {
	switch skillID {
	case 171:
		return proto.Profession_Alchemy
	case 164:
		return proto.Profession_Blacksmithing
	case 333:
		return proto.Profession_Enchanting
	case 202:
		return proto.Profession_Engineering
	case 182:
		return proto.Profession_Herbalism
	case 165:
		return proto.Profession_Leatherworking
	case 186:
		return proto.Profession_Mining
	case 393:
		return proto.Profession_Skinning
	case 197:
		return proto.Profession_Tailoring
	default:
		return proto.Profession_ProfessionUnknown
	}
}

func mapRequiredProfession(row []string, colIdx map[string]int) proto.Profession {
	skillID, ok := getInt(row, colIdx, "requiredSkillID")
	if !ok || skillID == 0 {
		return proto.Profession_ProfessionUnknown
	}

	return mapProfessionBySkillID(skillID)
}

func mapItemBindType(bonding int32) proto.ItemBindType {
	switch bonding {
	case 1:
		return proto.ItemBindType_ItemBindTypeBindOnPickup
	case 2:
		return proto.ItemBindType_ItemBindTypeBindOnEquip
	case 3:
		return proto.ItemBindType_ItemBindTypeBindOnUse
	case 4:
		return proto.ItemBindType_ItemBindTypeQuestItem
	default:
		return proto.ItemBindType_ItemBindTypeUnknown
	}
}

func mapItemEffectTriggerType(triggerID int32) proto.ItemEffectTriggerType {
	switch triggerID {
	case 0:
		return proto.ItemEffectTriggerType_ItemEffectTriggerTypeOnUse
	case 1:
		return proto.ItemEffectTriggerType_ItemEffectTriggerTypeOnEquip
	case 2:
		return proto.ItemEffectTriggerType_ItemEffectTriggerTypeChanceOnHit
	default:
		return proto.ItemEffectTriggerType_ItemEffectTriggerTypeUnknown
	}
}

func parseSpellEffectsCSV(csvData string) map[int32]SpellAnalysis {
	r := csv.NewReader(strings.NewReader(csvData))
	headers, err := r.Read()
	if err != nil {
		log.Fatalf("Cannot read spell csv header: %v", err)
	}

	colIdx := make(map[string]int)
	for i, name := range headers {
		colIdx[name] = i
	}

	requiredCols := []string{
		"id", "effect_1", "effect_2", "effect_3",
		"effectApplyAura_1", "effectApplyAura_2", "effectApplyAura_3",
		"effectBasePoints_1", "effectBasePoints_2", "effectBasePoints_3",
		"effectBaseDice_1", "effectBaseDice_2", "effectBaseDice_3",
		"effectMiscValue_1", "effectMiscValue_2", "effectMiscValue_3",
	}
	for _, col := range requiredCols {
		if _, ok := colIdx[col]; !ok {
			log.Fatalf("Missing required column %s in spell csv", col)
		}
	}

	analyses := make(map[int32]SpellAnalysis)

	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Cannot read spell csv row: %v", err)
		}

		spellID, err := strconv.Atoi(row[colIdx["id"]])
		if err != nil {
			continue
		}

		if analysis, ok := specialCaseSpellAnalysis(int32(spellID)); ok {
			analyses[int32(spellID)] = analysis
			continue
		}

		analysis := SpellAnalysis{}
		var resilienceDotValue float64
		var resilienceCritValue float64

		for i := 1; i <= 3; i++ {
			effectCol := fmt.Sprintf("effect_%d", i)
			auraCol := fmt.Sprintf("effectApplyAura_%d", i)
			basePointsCol := fmt.Sprintf("effectBasePoints_%d", i)
			baseDiceCol := fmt.Sprintf("effectBaseDice_%d", i)
			miscValueCol := fmt.Sprintf("effectMiscValue_%d", i)

			effectStr := row[colIdx[effectCol]]
			if effectStr == "" || effectStr == "0" {
				continue
			}

			effect, err := strconv.Atoi(effectStr)
			if err != nil {
				continue
			}

			// Non-aura effect => not reducible to raw item stats, likely a real item effect / proc / use effect.
			if effect != spellEffectApplyAura {
				analysis.HasUnhandledEffect = true
				continue
			}

			auraStr := row[colIdx[auraCol]]
			if auraStr == "" || auraStr == "0" {
				analysis.HasUnhandledEffect = true
				continue
			}
			aura, err := strconv.Atoi(auraStr)
			if err != nil {
				analysis.HasUnhandledEffect = true
				continue
			}

			basePoints, _ := strconv.ParseFloat(row[colIdx[basePointsCol]], 64)
			baseDice, _ := strconv.ParseFloat(row[colIdx[baseDiceCol]], 64)
			miscValue, _ := strconv.Atoi(row[colIdx[miscValueCol]])

			value := basePoints + baseDice
			if value == 0 {
				continue
			}
			switch int32(aura) {
			case auraModPeriodicDamagePercentTaken:
				if int32(miscValue) == 127 {
					resilienceDotValue += math.Abs(value)
					continue
				}
			case auraModCritDamageTaken:
				resilienceCritValue += math.Abs(value)
				continue
			}

			converted, handled := auraTypeToEffects(int32(aura), int32(miscValue), value)
			if !handled {
				analysis.HasUnhandledEffect = true
				continue
			}
			analysis.Effects = append(analysis.Effects, converted...)
		}

		normalizeSpellAnalysis(&analysis, resilienceDotValue, resilienceCritValue)

		if len(analysis.Effects) > 0 || analysis.HasUnhandledEffect {
			analyses[int32(spellID)] = analysis
		}
	}

	return analyses
}

func normalizeSpellAnalysis(
	analysis *SpellAnalysis,
	resilienceDotValue float64,
	resilienceCritValue float64,
) {
	statTotals := make([]float64, statsLen)
	weaponSkillTotals := make([]float64, weaponSkillLen)
	var bonusPhysicalDamage float64

	for _, effect := range analysis.Effects {
		if effect.Stat != proto.Stat(-1) {
			statIdx := int(effect.Stat)
			if statIdx >= 0 && statIdx < len(statTotals) {
				statTotals[statIdx] += effect.Value
			}
		}

		if effect.WeaponSkill != proto.WeaponSkill_WeaponSkillUnknown {
			wsIdx := int(effect.WeaponSkill)
			if wsIdx >= 0 && wsIdx < len(weaponSkillTotals) {
				weaponSkillTotals[wsIdx] += effect.Value
			}
		}

		if effect.BonusPhysicalDamage != 0 {
			bonusPhysicalDamage += effect.BonusPhysicalDamage
		}
	}

	if resilienceDotValue > 0 || resilienceCritValue > 0 {
		if resilienceDotValue == resilienceCritValue {
			statTotals[int(proto.Stat_StatResilience)] += resilienceDotValue
		} else {
			analysis.HasUnhandledEffect = true
		}
	}

	if statTotals[int(proto.Stat_StatSpellDamage)] > 0 && statTotals[proto.Stat_StatSpellDamage] == statTotals[proto.Stat_StatHealingPower] {
		statTotals[int(proto.Stat_StatSpellPower)] = statTotals[int(proto.Stat_StatSpellDamage)]
		statTotals[int(proto.Stat_StatSpellDamage)] = 0
		statTotals[int(proto.Stat_StatHealingPower)] = 0
	}

	analysis.Effects = analysis.Effects[:0]

	for statIdx, value := range statTotals {
		if value != 0 {
			analysis.Effects = append(analysis.Effects, SpellEffect{
				Stat:  proto.Stat(statIdx),
				Value: value,
			})
		}
	}

	for wsIdx, value := range weaponSkillTotals {
		if value != 0 {
			analysis.Effects = append(analysis.Effects, SpellEffect{
				WeaponSkill: proto.WeaponSkill(wsIdx),
				Value:       value,
			})
		}
	}

	if bonusPhysicalDamage != 0 {
		analysis.Effects = append(analysis.Effects, SpellEffect{
			BonusPhysicalDamage: bonusPhysicalDamage,
		})
	}
}

func auraTypeToEffects(auraType, miscValue int32, value float64) ([]SpellEffect, bool) {
	var effects []SpellEffect

	switch auraType {
	case auraModDamageDone:
		stats, hasPhysical := spellSchoolMaskToStats(miscValue)
		for _, stat := range stats {
			effects = append(effects, SpellEffect{Stat: stat, Value: value})
		}
		if hasPhysical {
			effects = append(effects, SpellEffect{BonusPhysicalDamage: value})
		}
		return effects, true

	case auraModStat:
		stat := mapItemModTypeToStat(miscValue)
		if stat != proto.Stat(-1) {
			effects = append(effects, SpellEffect{Stat: stat, Value: value})
			return effects, true
		}
		return nil, false

	case auraModResistance, auraModBaseResistance:
		stat := spellSchoolToResistanceStat(miscValue)
		if stat != proto.Stat(-1) {
			effects = append(effects, SpellEffect{Stat: stat, Value: value})
			return effects, true
		}
		return nil, false

	case auraModSpellHitChance, auraModAttackerSpellHitChance:
		return []SpellEffect{{Stat: proto.Stat_StatSpellHit, Value: value}}, true

	case auraModSpellCritChance, auraModSpellCritChanceSchool:
		return []SpellEffect{{Stat: proto.Stat_StatSpellCrit, Value: value}}, true

	case auraModMeleeHaste:
		return []SpellEffect{{Stat: proto.Stat_StatMeleeHaste, Value: value}}, true

	case auraModRangedHaste:
		// Intentionally ignored.
		// The sim has no ranged haste stat, and many item effects include both melee+ranged haste.
		// Treating this as "handled but no-op" avoids doubling haste and avoids exposing it as an extra effect.
		return nil, true

	case auraModCastingSpeedNotStack:
		return []SpellEffect{{Stat: proto.Stat_StatSpellHaste, Value: value}}, true

	case auraModAttackPower:
		return []SpellEffect{{Stat: proto.Stat_StatAttackPower, Value: value}}, true

	case auraModRangedAttackPower, auraModRangedAttackPowerVersus:
		return []SpellEffect{{Stat: proto.Stat_StatRangedAttackPower, Value: value}}, true

	case auraModIncreaseHealth:
		return []SpellEffect{{Stat: proto.Stat_StatHealth, Value: value}}, true

	case auraModIncreaseEnergy:
		switch miscValue {
		case powerTypeMana:
			return []SpellEffect{{Stat: proto.Stat_StatMana, Value: value}}, true
		case powerTypeRage:
			return []SpellEffect{{Stat: proto.Stat_StatRage, Value: value}}, true
		case powerTypeEnergy:
			return []SpellEffect{{Stat: proto.Stat_StatEnergy, Value: value}}, true
		default:
			return []SpellEffect{{Stat: proto.Stat_StatMana, Value: value}}, true
		}

	case auraModPowerRegen, auraModManaRegenInterrupt:
		return []SpellEffect{{Stat: proto.Stat_StatMP5, Value: value}}, true

	case auraModCritPercent:
		return []SpellEffect{{Stat: proto.Stat_StatMeleeCrit, Value: value}}, true

	case auraModHitChance, auraModAttackerMeleeHitChance, auraModAttackerRangedHitChance:
		return []SpellEffect{{Stat: proto.Stat_StatMeleeHit, Value: value}}, true

	case auraModHealingDone, auraModHealing:
		return []SpellEffect{{Stat: proto.Stat_StatHealingPower, Value: value}}, true

	case auraModSkill, auraModSkillTalent:
		if stat := skillIDToStat(miscValue); stat != proto.Stat(-1) {
			return []SpellEffect{{Stat: stat, Value: value}}, true
		} else if ws := skillIDToWeaponSkill(miscValue); ws != proto.WeaponSkill_WeaponSkillUnknown {
			return []SpellEffect{{WeaponSkill: ws, Value: value}}, true
		}
		return nil, false

	case auraModParrySkill, auraModParryPercent:
		return []SpellEffect{{Stat: proto.Stat_StatParry, Value: value}}, true

	case auraModDodgeSkill, auraModDodgePercent:
		return []SpellEffect{{Stat: proto.Stat_StatDodge, Value: value}}, true

	case auraModBlockSkill, auraModBlockPercent:
		return []SpellEffect{{Stat: proto.Stat_StatBlock, Value: value}}, true

	case auraModShieldBlockValue:
		return []SpellEffect{{Stat: proto.Stat_StatBlockValue, Value: value}}, true

	case auraModTargetResistance:
		return targetResistanceMaskToEffects(miscValue, value)

	case auraModFortune:
		return []SpellEffect{{Stat: proto.Stat_StatFortune, Value: value}}, true

	default:
		return nil, false
	}
}

func spellSchoolMaskToStats(mask int32) ([]proto.Stat, bool) {
	hasPhysical := mask&1 != 0

	if mask == 124 || mask == 126 || mask == 127 {
		return []proto.Stat{proto.Stat_StatSpellDamage}, hasPhysical
	}

	var stats []proto.Stat
	if mask&64 != 0 {
		stats = append(stats, proto.Stat_StatArcanePower)
	}
	if mask&4 != 0 {
		stats = append(stats, proto.Stat_StatFirePower)
	}
	if mask&16 != 0 {
		stats = append(stats, proto.Stat_StatFrostPower)
	}
	if mask&2 != 0 {
		stats = append(stats, proto.Stat_StatHolyPower)
	}
	if mask&8 != 0 {
		stats = append(stats, proto.Stat_StatNaturePower)
	}
	if mask&32 != 0 {
		stats = append(stats, proto.Stat_StatShadowPower)
	}

	if len(stats) == 0 && mask != 0 && !hasPhysical {
		stats = append(stats, proto.Stat_StatSpellDamage)
	}

	return stats, hasPhysical
}

func mapItemModTypeToStat(modType int32) proto.Stat {
	switch modType {
	case 0:
		return proto.Stat_StatMana
	case 1:
		return proto.Stat_StatHealth
	case 3:
		return proto.Stat_StatAgility
	case 4:
		return proto.Stat_StatStrength
	case 5:
		return proto.Stat_StatIntellect
	case 6:
		return proto.Stat_StatSpirit
	case 7:
		return proto.Stat_StatStamina
	default:
		return proto.Stat(-1)
	}
}

func spellSchoolToResistanceStat(school int32) proto.Stat {
	switch school {
	case 0:
		return proto.Stat_StatArmor
	case 2:
		return proto.Stat_StatFireResistance
	case 3:
		return proto.Stat_StatNatureResistance
	case 4:
		return proto.Stat_StatFrostResistance
	case 5:
		return proto.Stat_StatShadowResistance
	case 6:
		return proto.Stat_StatArcaneResistance
	default:
		return proto.Stat(-1)
	}
}

func skillIDToStat(skillID int32) proto.Stat {
	switch skillID {
	case 95:
		return proto.Stat_StatDefense
	default:
		return proto.Stat(-1)
	}
}

func skillIDToWeaponSkill(skillID int32) proto.WeaponSkill {
	switch skillID {
	case 43:
		return proto.WeaponSkill_WeaponSkillSwords
	case 44:
		return proto.WeaponSkill_WeaponSkillAxes
	case 45:
		return proto.WeaponSkill_WeaponSkillBows
	case 46:
		return proto.WeaponSkill_WeaponSkillGuns
	case 54:
		return proto.WeaponSkill_WeaponSkillMaces
	case 55:
		return proto.WeaponSkill_WeaponSkillTwoHandedSwords
	case 136:
		return proto.WeaponSkill_WeaponSkillStaves
	case 160:
		return proto.WeaponSkill_WeaponSkillTwoHandedMaces
	case 162:
		return proto.WeaponSkill_WeaponSkillUnarmed
	case 172:
		return proto.WeaponSkill_WeaponSkillTwoHandedAxes
	case 173:
		return proto.WeaponSkill_WeaponSkillDaggers
	case 176:
		return proto.WeaponSkill_WeaponSkillThrown
	case 226:
		return proto.WeaponSkill_WeaponSkillCrossbows
	case 229:
		return proto.WeaponSkill_WeaponSkillPolearms
	case 473:
		return proto.WeaponSkill_WeaponSkillUnarmed
	case 134:
		return proto.WeaponSkill_WeaponSkillFeralCombat
	default:
		return proto.WeaponSkill_WeaponSkillUnknown
	}
}

func parseFactionMapping(csvData string) map[int32]proto.UIItem_FactionRestriction {
	factionMap := make(map[int32]proto.UIItem_FactionRestriction)
	if csvData == "" {
		return factionMap
	}

	r := csv.NewReader(strings.NewReader(csvData))
	headers, err := r.Read()
	if err != nil {
		return factionMap
	}

	colIdx := make(map[string]int)
	for i, name := range headers {
		colIdx[name] = i
	}

	requiredCols := []string{"id", "repRaceMask_1", "repRaceMask_2", "repRaceMask_3", "repRaceMask_4"}
	for _, col := range requiredCols {
		if _, ok := colIdx[col]; !ok {
			return factionMap
		}
	}

	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		id, err := strconv.Atoi(row[colIdx["id"]])
		if err != nil {
			continue
		}

		var allianceCount, hordeCount int
		for i := 1; i <= 4; i++ {
			maskCol := fmt.Sprintf("repRaceMask_%d", i)
			maskStr := row[colIdx[maskCol]]
			if maskStr == "" || maskStr == "0" {
				continue
			}
			mask, err := strconv.ParseUint(maskStr, 10, 32)
			if err != nil {
				continue
			}

			if uint32(mask)&allianceRaceMask != 0 {
				allianceCount++
			}
			if uint32(mask)&hordeRaceMask != 0 {
				hordeCount++
			}
		}

		if allianceCount > 0 && hordeCount == 0 {
			factionMap[int32(id)] = proto.UIItem_FACTION_RESTRICTION_ALLIANCE_ONLY
		} else if hordeCount > 0 && allianceCount == 0 {
			factionMap[int32(id)] = proto.UIItem_FACTION_RESTRICTION_HORDE_ONLY
		}
	}

	return factionMap
}

func parseItemSetCSV(csvData string) ([]*proto.UIItemSet, map[int32]itemSetInfo) {
	var sets []*proto.UIItemSet
	itemToSet := make(map[int32]itemSetInfo)

	if csvData == "" {
		return sets, itemToSet
	}

	r := csv.NewReader(strings.NewReader(csvData))
	headers, err := r.Read()
	if err != nil {
		log.Printf("Cannot read item set csv header: %v", err)
		return sets, itemToSet
	}

	colIdx := make(map[string]int)
	for i, name := range headers {
		colIdx[name] = i
	}

	requiredCols := []string{"id", "name_enUS"}
	for _, col := range requiredCols {
		if _, ok := colIdx[col]; !ok {
			log.Printf("Missing required column %s in item set csv", col)
			return sets, itemToSet
		}
	}

	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		setID, err := strconv.Atoi(row[colIdx["id"]])
		if err != nil || setID == 0 {
			continue
		}

		setName := row[colIdx["name_enUS"]]
		if setName == "" {
			continue
		}

		itemSet := &proto.UIItemSet{
			Id:   int32(setID),
			Name: setName,
		}

		// Items
		for i := 1; i <= 17; i++ {
			col := fmt.Sprintf("itemId_%d", i)
			idx, ok := colIdx[col]
			if !ok || idx >= len(row) {
				continue
			}

			itemID, err := strconv.Atoi(row[idx])
			if err != nil || itemID == 0 {
				continue
			}

			itemSet.ItemIds = append(itemSet.ItemIds, int32(itemID))
			itemToSet[int32(itemID)] = itemSetInfo{
				SetID:   int32(setID),
				SetName: setName,
			}
		}

		// Bonuses
		for i := 1; i <= 8; i++ {
			spellCol := fmt.Sprintf("setSpellId_%d", i)
			thresholdCol := fmt.Sprintf("setThreshold_%d", i)

			spellID, okSpell := getInt(row, colIdx, spellCol)
			threshold, okThreshold := getInt(row, colIdx, thresholdCol)

			if !okSpell || !okThreshold || spellID == 0 || threshold == 0 {
				continue
			}

			itemSet.Bonuses = append(itemSet.Bonuses, &proto.UIItemSetBonus{
				PiecesRequired: threshold,
				SpellId:        spellID,
			})
		}

		// Sort bonuses by threshold, stable so duplicate thresholds keep CSV order.
		sort.SliceStable(itemSet.Bonuses, func(i, j int) bool {
			return itemSet.Bonuses[i].PiecesRequired < itemSet.Bonuses[j].PiecesRequired
		})

		// Optional profession requirement
		if requiredSkillID, ok := getInt(row, colIdx, "requiredSkillId"); ok && requiredSkillID != 0 {
			itemSet.RequiredProfession = mapProfessionBySkillID(requiredSkillID)
		}
		if requiredSkillRank, ok := getInt(row, colIdx, "requiredSkillRank"); ok && requiredSkillRank > 0 {
			itemSet.RequiredSkillRank = requiredSkillRank
		}

		sets = append(sets, itemSet)
	}

	return sets, itemToSet
}

// Helpers
func getString(row []string, colIdx map[string]int, col string) string {
	idx, ok := colIdx[col]
	if !ok || idx >= len(row) {
		return ""
	}
	return row[idx]
}

func getInt(row []string, colIdx map[string]int, col string) (int32, bool) {
	s := getString(row, colIdx, col)
	if s == "" {
		return 0, false
	}
	val, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, false
	}
	return int32(val), true
}

func getUint(row []string, colIdx map[string]int, col string) (uint32, bool) {
	s := getString(row, colIdx, col)
	if s == "" {
		return 0, false
	}
	val, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, false
	}
	return uint32(val), true
}

func getFloat(row []string, colIdx map[string]int, col string) (float64, bool) {
	s := getString(row, colIdx, col)
	if s == "" {
		return 0, false
	}
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}
	return val, true
}

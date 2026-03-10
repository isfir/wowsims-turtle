package database

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/isfir/wowsims-turtle/sim/core/proto"
)

const (
	spellAttrExChanneled1 = 0x00000004
	spellAttrExChanneled2 = 0x00000040
)

type rawSpellRow struct {
	ID            int32
	Name          string
	Icon          string
	Rank          int32
	RequiredLevel int32
	ManaCost      int32
	Range         int32
	CastTime      int32
	IsChannel     bool
	Duration      int32
	Cooldown      int32
	ProcChance    int32
	StackAmount   int32
	Description   string
	Subtext       string

	EffectBasePoints  [3]int32
	EffectDieSides    [3]int32
	EffectBaseDice    [3]int32
	EffectAmplitude   [3]int32
	EffectRadius      [3]int32
	EffectChainTarget [3]int32
}

var scaledSpellTokenRE = regexp.MustCompile(`\$\/(\d+);(?:(\d+))?([A-Za-z])(\d+)?`)
var spellTokenRE = regexp.MustCompile(`\$(?:(\d+))?([A-Za-z])(\d+)?`)
var inlineWhitespaceRE = regexp.MustCompile(`[ \t]{2,}`)

func maxInt32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func ParseTurtleSpellDB(
	spellCSV string,
	spellIconCSV string,
	spellRangeCSV string,
	spellDurationCSV string,
	spellCastTimesCSV string,
	spellRadiusCSV string,
) []*proto.UISpell {
	iconMap := parseSpellIconCSV(spellIconCSV)
	rangeMap := parseSpellRangeCSV(spellRangeCSV)
	durationMap := parseSpellDurationCSV(spellDurationCSV)
	castTimeMap := parseSpellCastTimesCSV(spellCastTimesCSV)
	radiusMap := parseSpellRadiusCSV(spellRadiusCSV)

	rawSpells := parseRawSpellCSV(spellCSV, iconMap, rangeMap, durationMap, castTimeMap, radiusMap)

	var spells []*proto.UISpell
	for _, raw := range rawSpells {
		resolvedDescription := resolveSpellDescription(raw.Description, raw.ID, rawSpells)

		spells = append(spells, &proto.UISpell{
			Id:            raw.ID,
			Name:          raw.Name,
			Icon:          raw.Icon,
			Rank:          raw.Rank,
			RequiredLevel: raw.RequiredLevel,
			ManaCost:      raw.ManaCost,
			CastTime:      raw.CastTime,
			IsChannel:     raw.IsChannel,
			Duration:      raw.Duration,
			Cooldown:      raw.Cooldown,
			Description:   resolvedDescription,
			Subtext:       raw.Subtext,
			Range:         raw.Range,
		})
	}

	return spells
}

func parseSpellIconCSV(csvData string) map[int32]string {
	r := csv.NewReader(strings.NewReader(csvData))
	if _, err := r.Read(); err != nil {
		log.Fatalf("Cannot read spell icon csv header: %v", err)
	}

	iconMap := make(map[int32]string)
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Cannot read spell icon csv row: %v", err)
		}
		if len(row) < 2 {
			continue
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			continue
		}

		texture := row[1]
		parts := strings.Split(texture, "\\")
		if len(parts) > 0 {
			texture = parts[len(parts)-1]
		}
		iconMap[int32(id)] = strings.ToLower(texture)
	}
	return iconMap
}

func parseSpellRangeCSV(csvData string) map[int32]int32 {
	r := csv.NewReader(strings.NewReader(csvData))
	headers, err := r.Read()
	if err != nil {
		log.Fatalf("Cannot read SpellRange csv header: %v", err)
	}

	colIdx := map[string]int{}
	for i, name := range headers {
		colIdx[name] = i
	}

	result := make(map[int32]int32)
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Cannot read SpellRange csv row: %v", err)
		}

		id, ok := getInt(row, colIdx, "id")
		if !ok {
			continue
		}

		rangeMax, ok := getFloat(row, colIdx, "rangeMax")
		if !ok {
			continue
		}

		result[id] = int32(rangeMax)
	}
	return result
}

func parseSpellDurationCSV(csvData string) map[int32]int32 {
	r := csv.NewReader(strings.NewReader(csvData))
	headers, err := r.Read()
	if err != nil {
		log.Fatalf("Cannot read SpellDuration csv header: %v", err)
	}

	colIdx := map[string]int{}
	for i, name := range headers {
		colIdx[name] = i
	}

	result := make(map[int32]int32)
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Cannot read SpellDuration csv row: %v", err)
		}

		id, ok := getInt(row, colIdx, "id")
		if !ok {
			continue
		}

		maxVal, ok := getInt(row, colIdx, "max")
		if !ok {
			continue
		}

		result[id] = maxVal
	}
	return result
}

func parseSpellCastTimesCSV(csvData string) map[int32]int32 {
	r := csv.NewReader(strings.NewReader(csvData))
	headers, err := r.Read()
	if err != nil {
		log.Fatalf("Cannot read SpellCastTimes csv header: %v", err)
	}

	colIdx := map[string]int{}
	for i, name := range headers {
		colIdx[name] = i
	}

	result := make(map[int32]int32)
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Cannot read SpellCastTimes csv row: %v", err)
		}

		id, ok := getInt(row, colIdx, "id")
		if !ok {
			continue
		}

		minimum, ok := getInt(row, colIdx, "minimum")
		if !ok {
			continue
		}

		result[id] = minimum
	}
	return result
}

func parseSpellRadiusCSV(csvData string) map[int32]int32 {
	r := csv.NewReader(strings.NewReader(csvData))
	headers, err := r.Read()
	if err != nil {
		log.Fatalf("Cannot read SpellRadius csv header: %v", err)
	}

	colIdx := map[string]int{}
	for i, name := range headers {
		colIdx[name] = i
	}

	result := make(map[int32]int32)
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Cannot read SpellRadius csv row: %v", err)
		}

		id, ok := getInt(row, colIdx, "id")
		if !ok {
			continue
		}

		radius, ok := getFloat(row, colIdx, "radius")
		if !ok {
			continue
		}

		result[id] = int32(math.Round(radius))
	}
	return result
}

func parseRawSpellCSV(
	csvData string,
	iconMap map[int32]string,
	rangeMap map[int32]int32,
	durationMap map[int32]int32,
	castTimeMap map[int32]int32,
	radiusMap map[int32]int32,
) map[int32]*rawSpellRow {
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
		"id",
		"spellIconId",
		"activeIconId",
		"name_enUS",
		"subtext_enUS",
		"description_enUS",
		"spellLevel",
		"manaCost",
		"castingTimeIndex",
		"durationIndex",
		"recoveryTime",
		"categoryRecoveryTime",
		"rangeIndex",
		"attribute_2",
		"procChance",
		"stackAmount",
	}
	for _, col := range requiredCols {
		if _, ok := colIdx[col]; !ok {
			log.Fatalf("Missing required column %s in spell csv", col)
		}
	}

	spells := make(map[int32]*rawSpellRow)

	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Cannot read spell csv row: %v", err)
		}

		id, ok := getInt(row, colIdx, "id")
		if !ok || id == 0 {
			continue
		}

		name := getString(row, colIdx, "name_enUS")
		if name == "" {
			continue
		}

		spellIconID, _ := getInt(row, colIdx, "spellIconId")
		activeIconID, _ := getInt(row, colIdx, "activeIconId")
		icon := iconMap[spellIconID]
		if icon == "" {
			icon = iconMap[activeIconID]
		}

		subtext := strings.TrimSpace(getString(row, colIdx, "subtext_enUS"))
		description := strings.TrimSpace(getString(row, colIdx, "description_enUS"))

		rank := int32(0)
		if rankStr, ok := strings.CutPrefix(subtext, "Rank "); ok {
			if rankNum, err := strconv.Atoi(rankStr); err == nil {
				rank = int32(rankNum)
			}
		}

		requiredLevel, _ := getInt(row, colIdx, "spellLevel")
		manaCost, _ := getInt(row, colIdx, "manaCost")
		castingTimeIndex, _ := getInt(row, colIdx, "castingTimeIndex")
		durationIndex, _ := getInt(row, colIdx, "durationIndex")
		rangeIndex, _ := getInt(row, colIdx, "rangeIndex")
		recoveryTime, _ := getInt(row, colIdx, "recoveryTime")
		categoryRecoveryTime, _ := getInt(row, colIdx, "categoryRecoveryTime")
		procChance, _ := getInt(row, colIdx, "procChance")
		stackAmount, _ := getInt(row, colIdx, "stackAmount")
		attribute2, _ := getInt(row, colIdx, "attribute_2")

		raw := &rawSpellRow{
			ID:            id,
			Name:          name,
			Icon:          icon,
			Rank:          rank,
			RequiredLevel: requiredLevel,
			ManaCost:      manaCost,
			Range:         rangeMap[rangeIndex],
			CastTime:      castTimeMap[castingTimeIndex],
			IsChannel:     (attribute2&spellAttrExChanneled1) != 0 || (attribute2&spellAttrExChanneled2) != 0,
			Duration:      durationMap[durationIndex],
			Cooldown:      maxInt32(recoveryTime, categoryRecoveryTime),
			ProcChance:    procChance,
			StackAmount:   stackAmount,
			Description:   description,
			Subtext:       subtext,
		}

		for i := 1; i <= 3; i++ {
			raw.EffectBasePoints[i-1], _ = getInt(row, colIdx, fmt.Sprintf("effectBasePoints_%d", i))
			raw.EffectDieSides[i-1], _ = getInt(row, colIdx, fmt.Sprintf("effectDieSides_%d", i))
			raw.EffectBaseDice[i-1], _ = getInt(row, colIdx, fmt.Sprintf("effectBaseDice_%d", i))
			raw.EffectAmplitude[i-1], _ = getInt(row, colIdx, fmt.Sprintf("effectAmplitude_%d", i))

			radiusIndex, _ := getInt(row, colIdx, fmt.Sprintf("effectRadiusIndex_%d", i))
			raw.EffectRadius[i-1] = radiusMap[radiusIndex]

			raw.EffectChainTarget[i-1], _ = getInt(row, colIdx, fmt.Sprintf("effectChainTarget_%d", i))
		}

		spells[id] = raw
	}

	return spells
}

func resolveSpellDescription(input string, currentSpellID int32, spells map[int32]*rawSpellRow) string {
	if input == "" {
		return ""
	}

	input = strings.ReplaceAll(input, "$B", "\n")
	input = strings.ReplaceAll(input, "$b", "\n")

	// First resolve scaled tokens like $/1000;s1
	resolved := scaledSpellTokenRE.ReplaceAllStringFunc(input, func(token string) string {
		matches := scaledSpellTokenRE.FindStringSubmatch(token)
		if len(matches) != 5 {
			return token
		}

		divisor, err := strconv.ParseFloat(matches[1], 64)
		if err != nil || divisor == 0 {
			return token
		}

		refSpell := resolveReferencedSpell(matches[2], currentSpellID, spells)
		if refSpell == nil {
			return token
		}

		index := parseSpellTokenIndex(matches[4])

		if text, ok := formatScaledSpellToken(refSpell, matches[3], index, divisor); ok {
			return text
		}

		return token
	})

	// Then resolve normal tokens like $s1, $22959u, $h, $a1
	resolved = spellTokenRE.ReplaceAllStringFunc(resolved, func(token string) string {
		matches := spellTokenRE.FindStringSubmatch(token)
		if len(matches) != 4 {
			return token
		}

		refSpell := resolveReferencedSpell(matches[1], currentSpellID, spells)
		if refSpell == nil {
			return token
		}

		index := parseSpellTokenIndex(matches[3])

		if text, ok := formatSpellToken(refSpell, matches[2], index); ok {
			return text
		}

		return token
	})

	// Clean repeated spaces per line, while keeping line breaks.
	lines := strings.Split(resolved, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(inlineWhitespaceRE.ReplaceAllString(line, " "))
	}

	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func resolveReferencedSpell(rawSpellID string, currentSpellID int32, spells map[int32]*rawSpellRow) *rawSpellRow {
	refSpellID := currentSpellID

	if rawSpellID != "" {
		if parsed, err := strconv.Atoi(rawSpellID); err == nil {
			refSpellID = int32(parsed)
		}
	}

	return spells[refSpellID]
}

func parseSpellTokenIndex(rawIndex string) int {
	if rawIndex == "" {
		return 0
	}

	parsed, err := strconv.Atoi(rawIndex)
	if err != nil || parsed <= 0 {
		return 0
	}

	return parsed - 1
}

func formatSpellToken(spell *rawSpellRow, op string, index int) (string, bool) {
	switch op {
	case "s", "S":
		return formatSpellEffectValue(spell, index), true
	case "o", "O":
		return formatSpellEffectOverTime(spell, index), true
	case "d", "D":
		return formatSpellDurationText(spell.Duration), true
	case "t", "T":
		return formatSpellEffectAmplitudeText(spell, index), true
	case "a", "A":
		return formatSpellRadiusValue(spell, index), true
	case "h", "H":
		return formatSpellProcChance(spell), true
	case "u", "U":
		return formatSpellStackAmount(spell), true
	case "m":
		return formatSpellMinEffectValue(spell, index), true
	case "M":
		return formatSpellMaxEffectValue(spell, index), true
	case "n", "N":
		return formatSpellName(spell), true
	case "r", "R":
		return formatSpellRangeValue(spell), true
	case "x", "X":
		return formatSpellChainTargetValue(spell, index), true
	default:
		return "", false
	}
}

func formatScaledSpellToken(spell *rawSpellRow, op string, index int, divisor float64) (string, bool) {
	switch op {
	case "s", "S":
		return formatScaledSpellEffectValue(spell, index, divisor), true
	case "o", "O":
		return formatScaledSpellEffectOverTime(spell, index, divisor), true
	case "d", "D":
		return formatScaledNumber(float64(spell.Duration) / divisor), true
	case "t", "T":
		if spell == nil || index < 0 || index >= 3 {
			return "", false
		}
		return formatScaledNumber(float64(spell.EffectAmplitude[index]) / divisor), true
	case "a", "A":
		if spell == nil || index < 0 || index >= 3 {
			return "", false
		}
		return formatScaledNumber(float64(spell.EffectRadius[index]) / divisor), true
	case "h", "H":
		return formatScaledNumber(float64(spell.ProcChance) / divisor), true
	case "u", "U":
		return formatScaledNumber(float64(spell.StackAmount) / divisor), true
	case "m":
		low, _, ok := getSpellEffectValueBounds(spell, index)
		if !ok {
			return "", false
		}
		return formatScaledNumber(float64(low) / divisor), true
	case "M":
		_, high, ok := getSpellEffectValueBounds(spell, index)
		if !ok {
			return "", false
		}
		return formatScaledNumber(float64(high) / divisor), true
	case "r", "R":
		if spell == nil {
			return "", false
		}
		return formatScaledNumber(float64(spell.Range) / divisor), true
	case "x", "X":
		if spell == nil || index < 0 || index >= 3 {
			return "", false
		}
		return formatScaledNumber(float64(spell.EffectChainTarget[index]) / divisor), true
	default:
		return "", false
	}
}

func getSpellEffectValueBounds(spell *rawSpellRow, index int) (int32, int32, bool) {
	if spell == nil || index < 0 || index >= 3 {
		return 0, 0, false
	}

	base := spell.EffectBasePoints[index] + spell.EffectBaseDice[index]
	dieSides := spell.EffectDieSides[index]

	// Fixed value: use the resolved base directly.
	if dieSides <= 1 {
		value := int32(math.Abs(float64(base)))
		return value, value, true
	}

	other := base + dieSides

	low := int32(math.Abs(float64(base)))
	high := int32(math.Abs(float64(other)))
	if low > high {
		low, high = high, low
	}

	return low, high, true
}

func getSpellEffectOverTimeBounds(spell *rawSpellRow, index int) (int32, int32, bool) {
	if spell == nil || index < 0 || index >= 3 {
		return 0, 0, false
	}

	amplitude := spell.EffectAmplitude[index]
	if amplitude <= 0 || spell.Duration <= 0 {
		return getSpellEffectValueBounds(spell, index)
	}

	ticks := spell.Duration / amplitude
	if ticks <= 0 {
		return getSpellEffectValueBounds(spell, index)
	}

	basePerTick := spell.EffectBasePoints[index] + spell.EffectBaseDice[index]
	dieSides := spell.EffectDieSides[index]

	// Fixed value over time: scale the resolved base directly.
	if dieSides <= 1 {
		total := int32(math.Abs(float64(basePerTick * ticks)))
		return total, total, true
	}

	otherPerTick := basePerTick + dieSides

	low := int32(math.Abs(float64(basePerTick * ticks)))
	high := int32(math.Abs(float64(otherPerTick * ticks)))
	if low > high {
		low, high = high, low
	}

	return low, high, true
}

func formatSpellMinEffectValue(spell *rawSpellRow, index int) string {
	low, _, ok := getSpellEffectValueBounds(spell, index)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%d", low)
}

func formatSpellMaxEffectValue(spell *rawSpellRow, index int) string {
	_, high, ok := getSpellEffectValueBounds(spell, index)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%d", high)
}

func formatSpellEffectValue(spell *rawSpellRow, index int) string {
	low, high, ok := getSpellEffectValueBounds(spell, index)
	if !ok {
		return ""
	}

	if high > low {
		return fmt.Sprintf("%d to %d", low, high)
	}

	return fmt.Sprintf("%d", low)
}

func formatSpellEffectOverTime(spell *rawSpellRow, index int) string {
	low, high, ok := getSpellEffectOverTimeBounds(spell, index)
	if !ok {
		return ""
	}

	if high > low {
		return fmt.Sprintf("%d to %d", low, high)
	}

	return fmt.Sprintf("%d", low)
}

func formatSpellDurationText(ms int32) string {
	if ms <= 0 {
		return "0 sec"
	}

	if ms%60000 == 0 {
		minutes := ms / 60000
		if minutes == 1 {
			return "1 min"
		}
		return fmt.Sprintf("%d min", minutes)
	}

	if ms%1000 == 0 {
		seconds := ms / 1000
		if seconds == 1 {
			return "1 sec"
		}
		return fmt.Sprintf("%d sec", seconds)
	}

	return fmt.Sprintf("%.1f sec", float64(ms)/1000.0)
}

func formatSpellRadiusValue(spell *rawSpellRow, index int) string {
	if spell == nil || index < 0 || index >= 3 {
		return ""
	}

	radius := spell.EffectRadius[index]
	if radius <= 0 {
		return "0"
	}

	return fmt.Sprintf("%d", radius)
}

func formatSpellProcChance(spell *rawSpellRow) string {
	if spell == nil || spell.ProcChance <= 0 {
		return "0"
	}

	return fmt.Sprintf("%d", spell.ProcChance)
}

func formatSpellStackAmount(spell *rawSpellRow) string {
	if spell == nil || spell.StackAmount <= 0 {
		return "0"
	}

	return fmt.Sprintf("%d", spell.StackAmount)
}

func formatSpellEffectAmplitudeText(spell *rawSpellRow, index int) string {
	if spell == nil || index < 0 || index >= 3 {
		return "0 sec"
	}

	return formatSpellDurationText(spell.EffectAmplitude[index])
}

func formatSpellName(spell *rawSpellRow) string {
	if spell == nil {
		return ""
	}
	return spell.Name
}

func formatSpellRangeValue(spell *rawSpellRow) string {
	if spell == nil || spell.Range <= 0 {
		return "0"
	}
	return fmt.Sprintf("%d", spell.Range)
}

func formatSpellChainTargetValue(spell *rawSpellRow, index int) string {
	if spell == nil || index < 0 || index >= 3 {
		return "0"
	}

	value := spell.EffectChainTarget[index]
	if value <= 0 {
		return "0"
	}

	return fmt.Sprintf("%d", value)
}

func formatScaledSpellEffectValue(spell *rawSpellRow, index int, divisor float64) string {
	low, high, ok := getSpellEffectValueBounds(spell, index)
	if !ok || divisor == 0 {
		return ""
	}

	if high > low {
		return fmt.Sprintf("%s to %s", formatScaledNumber(float64(low)/divisor), formatScaledNumber(float64(high)/divisor))
	}

	return formatScaledNumber(float64(low) / divisor)
}

func formatScaledSpellEffectOverTime(spell *rawSpellRow, index int, divisor float64) string {
	low, high, ok := getSpellEffectOverTimeBounds(spell, index)
	if !ok || divisor == 0 {
		return ""
	}

	if high > low {
		return fmt.Sprintf("%s to %s", formatScaledNumber(float64(low)/divisor), formatScaledNumber(float64(high)/divisor))
	}

	return formatScaledNumber(float64(low) / divisor)
}

func formatScaledNumber(value float64) string {
	if math.Abs(value-math.Round(value)) < 0.000001 {
		return fmt.Sprintf("%d", int64(math.Round(value)))
	}

	text := strconv.FormatFloat(value, 'f', 3, 64)
	text = strings.TrimRight(text, "0")
	text = strings.TrimRight(text, ".")
	return text
}

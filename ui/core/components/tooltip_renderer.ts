import {
	Class,
	HandType,
	ItemQuality,
	ItemType,
	Profession,
	Stat,
	WeaponSkill,
} from '../proto/common';
import { UIItem,UISpell } from '../proto/ui';
import { Database } from '../proto_utils/database';
import { getEnchantDescription } from '../proto_utils/enchants';
import {
	armorTypeNames,
	classNames,
	getClassStatName,
	handTypeNames,
	itemBindTypeNames,
	itemEffectTriggerNames,
	itemTypeTooltipNames,
	professionNames,
	rangedWeaponTypeNames,
	statOrder,
	weaponSkillNames,
	weaponSkillOrder,
	weaponTypeNames,
} from '../proto_utils/names';
import { ItemTooltipContext, SpellTooltipContext } from '../tooltip_context';

type ResolvedItemEffectLine = {
	text: string;
	isImplemented: boolean;
};

type ResolvedSetBonus = {
	piecesRequired: number;
	spellId: number;
	description: string;
	isImplemented: boolean;
};

type TooltipRenderableItem = UIItem & {
	tooltipResolvedEnchantText?: string;
	tooltipResolvedSetName?: string;
	tooltipSetItems?: UIItem[];
	tooltipResolvedEffectLines?: ResolvedItemEffectLine[];
	tooltipResolvedSetBonuses?: ResolvedSetBonus[];
};

type ItemTooltipRenderContext = Partial<ItemTooltipContext> & {
	playerClass?: Class;
	professions?: Profession[];
};

const ITEM_TYPE_SORT_ORDER = new Map<number, number>([
	[ItemType.ItemTypeHead, 0],
	[ItemType.ItemTypeNeck, 1],
	[ItemType.ItemTypeShoulder, 2],
	[ItemType.ItemTypeBack, 3],
	[ItemType.ItemTypeChest, 4],
	[ItemType.ItemTypeWrist, 5],
	[ItemType.ItemTypeHands, 6],
	[ItemType.ItemTypeWaist, 7],
	[ItemType.ItemTypeLegs, 8],
	[ItemType.ItemTypeFeet, 9],
	[ItemType.ItemTypeFinger, 10],
	[ItemType.ItemTypeTrinket, 11],
	[ItemType.ItemTypeWeapon, 12],
	[ItemType.ItemTypeRanged, 13],
]);

function escapeHtml(text: string): string {
	const div = document.createElement('div');
	div.textContent = text;
	return div.innerHTML;
}

function classList(...classes: Array<string | false | null | undefined>): string {
	return classes.filter(Boolean).join(' ');
}

function textLine(className: string, text: string): string {
	return `<div class="${className}">${escapeHtml(text)}</div>`;
}

function multilineTextLine(className: string, text: string): string {
	return `<div class="${className}">${escapeHtml(text).replace(/\n/g, '<br>')}</div>`;
}

function pairedTextLine(className: string, left: string, right?: string | null): string {
	return `
		<div class="${className} tooltip-line-pair">
			<span>${escapeHtml(left)}</span>
			<span>${escapeHtml(right ?? '')}</span>
		</div>
	`;
}

function requirementLine(text: string, unmet = false): string {
	return `<div class="${classList('tooltip-req', unmet && 'is-unmet')}">${escapeHtml(text)}</div>`;
}

function isDefined<T>(value: T | null | undefined): value is T {
	return value !== null && value !== undefined;
}

function formatNumber(value: number, maxFractionDigits = 2): string {
	return value.toLocaleString(undefined, {
		minimumFractionDigits: 0,
		maximumFractionDigits: maxFractionDigits,
	});
}

function formatWholeNumber(value: number): string {
	return formatNumber(value, 0);
}

function formatSignedWholeNumber(value: number): string {
	const sign = value >= 0 ? '+' : '-';
	return `${sign}${formatWholeNumber(Math.abs(value))}`;
}

function formatFixed(value: number, digits: number): string {
	return value.toLocaleString(undefined, {
		minimumFractionDigits: digits,
		maximumFractionDigits: digits,
	});
}

function formatShortDuration(ms: number): string {
	if (ms <= 0) {
		return '0 sec';
	}

	if (ms % 60000 === 0) {
		const minutes = ms / 60000;
		return `${minutes} min`;
	}

	if (ms % 1000 === 0) {
		const seconds = ms / 1000;
		return `${seconds} sec`;
	}

	return `${(ms / 1000).toFixed(1)} sec`;
}

function ensureSentence(text: string): string {
	const trimmed = text.trim();
	if (!trimmed) {
		return trimmed;
	}
	return /[.!?]$/.test(trimmed) ? trimmed : `${trimmed}.`;
}

function getQualityClassName(quality: ItemQuality): string {
	switch (quality) {
		case ItemQuality.ItemQualityJunk:
			return 'is-junk';
		case ItemQuality.ItemQualityCommon:
			return 'is-common';
		case ItemQuality.ItemQualityUncommon:
			return 'is-uncommon';
		case ItemQuality.ItemQualityRare:
			return 'is-rare';
		case ItemQuality.ItemQualityEpic:
			return 'is-epic';
		case ItemQuality.ItemQualityLegendary:
			return 'is-legendary';
		case ItemQuality.ItemQualityArtifact:
			return 'is-artifact';
		case ItemQuality.ItemQualityHeirloom:
			return 'is-heirloom';
		default:
			return 'is-common';
	}
}

function getStatValue(item: UIItem, stat: Stat): number {
	return item.stats?.[stat] ?? 0;
}

function getWeaponSkillValue(item: UIItem, skill: WeaponSkill): number {
	return item.weaponSkills?.[skill] ?? 0;
}

function getDisplayedItemLevel(item: UIItem, params?: ItemTooltipRenderContext): number {
	return params?.itemLevel && params.itemLevel > 0 ? params.itemLevel : item.ilvl;
}

function getLeftTypeLabel(item: UIItem): string | null {
	if (item.type === ItemType.ItemTypeRanged) {
		return itemTypeTooltipNames.get(item.type) ?? null;
	}

	if (item.type === ItemType.ItemTypeWeapon) {
		const handType = handTypeNames.get(item.handType);
		if (handType && item.handType !== HandType.HandTypeUnknown) {
			return handType;
		}
	}

	return itemTypeTooltipNames.get(item.type) ?? null;
}

function getRightTypeLabel(item: UIItem): string | null {
	if (item.weaponType) {
		return weaponTypeNames.get(item.weaponType) ?? null;
	}

	if (item.rangedWeaponType) {
		return rangedWeaponTypeNames.get(item.rangedWeaponType) ?? null;
	}

	if (item.armorType) {
		return armorTypeNames.get(item.armorType) ?? null;
	}

	return null;
}

function getArmorLine(item: UIItem): string | null {
	const armor = getStatValue(item, Stat.StatArmor);
	if (armor <= 0) {
		return null;
	}

	return `${formatWholeNumber(armor)} Armor`;
}

function getWeaponSectionLines(item: UIItem): string[] {
	const lines: string[] = [];

	const minDamage = item.weaponDamageMin ?? 0;
	const maxDamage = item.weaponDamageMax ?? 0;
	const weaponSpeed = item.weaponSpeed ?? 0;

	if ((minDamage > 0 || maxDamage > 0) && weaponSpeed > 0) {
		lines.push(
			pairedTextLine(
				'tooltip-stat',
				`${formatWholeNumber(minDamage)} - ${formatWholeNumber(maxDamage)} Damage`,
				`Speed ${formatFixed(weaponSpeed, 2)}`,
			),
		);

		const dps = ((minDamage + maxDamage) / 2) / weaponSpeed;
		lines.push(textLine('tooltip-note', `(${formatFixed(dps, 1)} damage per second)`));
	} else if (minDamage > 0 || maxDamage > 0) {
		lines.push(textLine('tooltip-stat', `${formatWholeNumber(minDamage)} - ${formatWholeNumber(maxDamage)} Damage`));
	}

	if ((item.bonusPhysicalDamage ?? 0) > 0) {
		lines.push(
			multilineTextLine(
				'tooltip-effect',
				`Equip: Increases weapon damage by ${formatWholeNumber(item.bonusPhysicalDamage)}.`,
			),
		);
	}

	return lines;
}

function formatDirectStatLine(stat: Stat, value: number, playerClass?: Class): string | null {
	switch (stat) {
		case Stat.StatHealth:
			return `${formatSignedWholeNumber(value)} Health`;
		case Stat.StatMana:
			return `${formatSignedWholeNumber(value)} Mana`;

		case Stat.StatStamina:
		case Stat.StatStrength:
		case Stat.StatAgility:
		case Stat.StatIntellect:
		case Stat.StatSpirit:
			return `${formatSignedWholeNumber(value)} ${getClassStatName(stat, playerClass)}`;

		case Stat.StatAttackPower:
			return `${formatSignedWholeNumber(value)} Attack Power`;
		case Stat.StatRangedAttackPower:
			return `${formatSignedWholeNumber(value)} Ranged Attack Power`;
		case Stat.StatDefense:
			return `${formatSignedWholeNumber(value)} Defense`;
		case Stat.StatBlockValue:
			return `${formatSignedWholeNumber(value)} Block Value`;
		case Stat.StatBonusArmor:
			return `${formatSignedWholeNumber(value)} Bonus Armor`;

		case Stat.StatArcaneResistance:
		case Stat.StatFireResistance:
		case Stat.StatFrostResistance:
		case Stat.StatNatureResistance:
		case Stat.StatShadowResistance:
			return `${formatSignedWholeNumber(value)} ${getClassStatName(stat, playerClass)}`;

		default:
			return null;
	}
}

function formatEffectStatLine(stat: Stat, value: number): string | null {
	if (value <= 0) {
		return null;
	}

	switch (stat) {
		case Stat.StatSpellPower:
			return `Equip: Increases damage and healing done by magical spells and effects by up to ${formatWholeNumber(value)}.`;
		case Stat.StatSpellDamage:
			return `Equip: Increases damage done by magical spells and effects by up to ${formatWholeNumber(value)}.`;
		case Stat.StatHealingPower:
			return `Equip: Increases healing done by spells and effects by up to ${formatWholeNumber(value)}.`;

		case Stat.StatArcanePower:
			return `Equip: Increases damage done by Arcane spells and effects by up to ${formatWholeNumber(value)}.`;
		case Stat.StatFirePower:
			return `Equip: Increases damage done by Fire spells and effects by up to ${formatWholeNumber(value)}.`;
		case Stat.StatFrostPower:
			return `Equip: Increases damage done by Frost spells and effects by up to ${formatWholeNumber(value)}.`;
		case Stat.StatHolyPower:
			return `Equip: Increases damage done by Holy spells and effects by up to ${formatWholeNumber(value)}.`;
		case Stat.StatNaturePower:
			return `Equip: Increases damage done by Nature spells and effects by up to ${formatWholeNumber(value)}.`;
		case Stat.StatShadowPower:
			return `Equip: Increases damage done by Shadow spells and effects by up to ${formatWholeNumber(value)}.`;

		case Stat.StatSpellHit:
			return `Equip: Improves your chance to hit with spells by ${formatNumber(value)}%.`;
		case Stat.StatSpellCrit:
			return `Equip: Improves your chance to get a critical strike with spells by ${formatNumber(value)}%.`;
		case Stat.StatSpellPenetration:
			return `Equip: Decreases the magical resistances of your spell targets by ${formatWholeNumber(value)}.`;
		case Stat.StatMP5:
			return `Equip: Restores ${formatWholeNumber(value)} mana per 5 sec.`;

		case Stat.StatMeleeHit:
			return `Equip: Improves your chance to hit by ${formatNumber(value)}%.`;
		case Stat.StatMeleeCrit:
			return `Equip: Improves your chance to get a critical strike by ${formatNumber(value)}%.`;
		case Stat.StatArmorPenetration:
			return `Equip: Your attacks ignore ${formatWholeNumber(value)} of your target's armor.`;
		case Stat.StatExpertise:
			return `Equip: Increases your expertise by ${formatNumber(value)}.`;

		case Stat.StatFeralAttackPower:
			return `Equip: Increases attack power in Cat, Bear, Dire Bear, and Moonkin forms only by ${formatWholeNumber(value)}.`;

		case Stat.StatDodge:
			return `Equip: Improves your chance to dodge by ${formatNumber(value)}%.`;
		case Stat.StatParry:
			return `Equip: Improves your chance to parry by ${formatNumber(value)}%.`;
		case Stat.StatBlock:
			return `Equip: Improves your chance to block by ${formatNumber(value)}%.`;

		case Stat.StatFortune:
			return `Equip: Increases your chance to trigger effects from equipped items by ${formatNumber(value)}%.`;
		case Stat.StatResilience:
			return `Equip: Reduces damage taken from critical hits and damage over time effects by ${formatNumber(value)}%.`;

		case Stat.StatEnergy:
			return `Equip: Increases your maximum Energy by ${formatWholeNumber(value)}.`;
		case Stat.StatRage:
			return `Equip: Increases your maximum Rage by ${formatWholeNumber(value)}.`;

		default:
			return null;
	}
}

function getDirectStatLines(item: UIItem, playerClass?: Class): string[] {
	const lines: string[] = [];

	for (const stat of statOrder) {
		if (stat === Stat.StatArmor) {
			continue;
		}

		const value = getStatValue(item, stat);
		if (!value) {
			continue;
		}

		const line = formatDirectStatLine(stat, value, playerClass);
		if (line) {
			lines.push(line);
		}
	}

	return lines;
}

function getWeaponSkillEffectLines(item: UIItem): string[] {
	const lines: string[] = [];

	for (const skill of weaponSkillOrder) {
		const value = getWeaponSkillValue(item, skill);
		if (!value) {
			continue;
		}

		const skillName = weaponSkillNames.get(skill);
		if (!skillName || skill === WeaponSkill.WeaponSkillUnknown) {
			continue;
		}

		lines.push(`Equip: Increases your skill with ${skillName} by ${formatWholeNumber(value)}.`);
	}

	return lines;
}

function getEffectStatLines(item: UIItem): string[] {
	const lines: string[] = [];

	const spellHaste = getStatValue(item, Stat.StatSpellHaste);
	const meleeHaste = getStatValue(item, Stat.StatMeleeHaste);

	for (const stat of statOrder) {
		if (stat === Stat.StatArmor) {
			continue;
		}

		if (stat === Stat.StatSpellHaste) {
			if (spellHaste > 0 && meleeHaste > 0 && spellHaste === meleeHaste) {
				lines.push(`Equip: Increases your attack and casting speed by ${formatNumber(spellHaste)}%.`);
			} else if (spellHaste > 0) {
				lines.push(`Equip: Increases your spell casting speed by ${formatNumber(spellHaste)}%.`);
			}
			continue;
		}

		if (stat === Stat.StatMeleeHaste) {
			if (spellHaste > 0 && meleeHaste > 0 && spellHaste === meleeHaste) {
				continue;
			}

			if (meleeHaste > 0) {
				lines.push(`Equip: Increases your attack speed by ${formatNumber(meleeHaste)}%.`);
			}
			continue;
		}

		const value = getStatValue(item, stat);
		if (!value) {
			continue;
		}

		const line = formatEffectStatLine(stat, value);
		if (line) {
			lines.push(line);
		}
	}

	lines.push(...getWeaponSkillEffectLines(item));

	return lines;
}

function getClassRestrictionLine(item: UIItem, params?: ItemTooltipRenderContext): string | null {
	if (!item.classAllowlist?.length) {
		return null;
	}

	const classListText = item.classAllowlist
		.map(classId => classNames.get(classId))
		.filter(isDefined)
		.join(', ');

	if (!classListText) {
		return null;
	}

	const unmet = params?.playerClass !== undefined && !item.classAllowlist.includes(params.playerClass);
	return requirementLine(`Classes: ${classListText}`, unmet);
}

function getProfessionRestrictionLine(item: UIItem, params?: ItemTooltipRenderContext): string | null {
	if (!item.requiredProfession) {
		return null;
	}

	const professionName = professionNames.get(item.requiredProfession);
	if (!professionName) {
		return null;
	}

	const professions = params?.professions ?? [];
	const unmet = professions.length > 0 && !professions.includes(item.requiredProfession);

	return requirementLine(`Requires ${professionName}`, unmet);
}

function getFactionRestrictionLine(item: UIItem): string | null {
	switch (item.factionRestriction) {
		case 1:
			return requirementLine('Alliance Only');
		case 2:
			return requirementLine('Horde Only');
		default:
			return null;
	}
}

function getItemRequiredLevelLine(item: UIItem, params?: ItemTooltipRenderContext): string | null {
	if (!item.requiredLevel || item.requiredLevel <= 0) {
		return null;
	}

	const unmet = params?.level !== undefined && params.level < item.requiredLevel;
	return requirementLine(`Requires Level ${item.requiredLevel}`, unmet);
}

function getSpellCostText(spell: UISpell): string | null {
	if (!spell.manaCost || spell.manaCost <= 0) {
		return null;
	}
	return `${spell.manaCost} Mana`;
}

function getSpellRangeText(spell: UISpell): string | null {
	if (!spell.range || spell.range <= 0) {
		return null;
	}
	return `${spell.range} yd range`;
}

function getSpellCastText(spell: UISpell): string {
	if (spell.isChannel) {
		if (spell.duration > 0) {
			return `Channeled (${formatShortDuration(spell.duration)} cast)`;
		}
		return 'Channeled';
	}

	if (spell.castTime > 0) {
		return `${formatShortDuration(spell.castTime)} cast`;
	}

	return 'Instant';
}

function getSpellCooldownText(spell: UISpell): string | null {
	if (!spell.cooldown || spell.cooldown <= 0) {
		return null;
	}
	return `${formatShortDuration(spell.cooldown)} cooldown`;
}

function getItemEffectPrefix(triggerType: number): string | null {
	const prefix = itemEffectTriggerNames.get(triggerType as never);
	return prefix || null;
}

function resolveItemEffectLines(item: UIItem, db: Database): ResolvedItemEffectLine[] {
	const lines: ResolvedItemEffectLine[] = [];

	for (const effect of item.effects ?? []) {
		const prefix = getItemEffectPrefix(effect.triggerType);
		if (!prefix) {
			continue;
		}

		let body = '';
		const spell = db.getSpellById(effect.spellId);

		if (spell?.description) {
			body = ensureSentence(spell.description);
		} else if (spell?.name) {
			body = spell.name;
		} else {
			body = `Spell ${effect.spellId}`;
		}

		lines.push({
			text: `${prefix}: ${body}`,
			isImplemented: item.hasImplementedEffects,
		});
	}

	return lines;
}

function resolveSetBonusDescription(spell: UISpell | null | undefined, spellId: number): string {
	if (spell?.description?.trim()) {
		return spell.description.trim();
	}

	if (spell?.name?.trim()) {
		return spell.name.trim();
	}

	return `Spell ${spellId}`;
}

function getSetItems(item: TooltipRenderableItem): UIItem[] {
	const setItems = item.tooltipSetItems ?? [];

	return [...setItems].sort((left, right) => {
		const leftOrder = ITEM_TYPE_SORT_ORDER.get(left.type) ?? 999;
		const rightOrder = ITEM_TYPE_SORT_ORDER.get(right.type) ?? 999;

		if (leftOrder !== rightOrder) {
			return leftOrder - rightOrder;
		}

		return left.name.localeCompare(right.name);
	});
}

function resolveItemSetData(
	item: UIItem,
	db: Database,
): Pick<TooltipRenderableItem, 'tooltipResolvedSetName' | 'tooltipSetItems' | 'tooltipResolvedSetBonuses'> {
	if (!item.setId) {
		return {
			tooltipResolvedSetName: item.setName,
			tooltipSetItems: [],
			tooltipResolvedSetBonuses: [],
		};
	}

	const itemSet = db.getItemSetById(item.setId);
	if (!itemSet) {
		return {
			tooltipResolvedSetName: item.setName,
			tooltipSetItems: [],
			tooltipResolvedSetBonuses: [],
		};
	}

	const tooltipSetItems = (itemSet.itemIds ?? [])
		.map(itemId => db.getItemById(itemId))
		.filter((setItem): setItem is UIItem => Boolean(setItem));

	const tooltipResolvedSetBonuses: ResolvedSetBonus[] = (itemSet.bonuses ?? []).map(bonus => {
		const spell = db.getSpellById(bonus.spellId);

		return {
			piecesRequired: bonus.piecesRequired,
			spellId: bonus.spellId,
			description: resolveSetBonusDescription(spell, bonus.spellId),
			isImplemented: bonus.isImplemented,
		};
	});

	return {
		tooltipResolvedSetName: itemSet.name || item.setName,
		tooltipSetItems,
		tooltipResolvedSetBonuses,
	};
}

function renderItemHeader(item: UIItem): string {
	const qualityClass = getQualityClassName(item.quality);
	const name = item.name || `Item ${item.id}`;

	return `<div class="tooltip-name ${qualityClass}">${escapeHtml(name)}</div>`;
}

function renderSetSection(item: TooltipRenderableItem, params?: ItemTooltipRenderContext): string {
	const setName = item.tooltipResolvedSetName ?? item.setName;
	if (!setName) {
		return '';
	}

	const lines: string[] = [];
	const setItems = getSetItems(item);
	const setBonuses = item.tooltipResolvedSetBonuses ?? [];
	const equippedSetPieceIds = new Set(params?.setPieceIds ?? []);

	const equippedCount = setItems.filter(setItem => equippedSetPieceIds.has(setItem.id)).length;
	const totalCount = setItems.length;

	const header = totalCount > 0
		? `${setName} (${equippedCount}/${totalCount})`
		: setName;

	lines.push(textLine('tooltip-set-name', header));

	for (const setItem of setItems) {
		const itemClass = equippedSetPieceIds.has(setItem.id)
			? 'tooltip-set-item is-equipped'
			: 'tooltip-set-item';

		lines.push(textLine(itemClass, setItem.name));
	}

	for (let i = 0; i < setBonuses.length; i++) {
		const bonus = setBonuses[i];
		const text = bonus.isImplemented
			? `(${bonus.piecesRequired}) Set: ${bonus.description}`
			: `(${bonus.piecesRequired}) Set: ${bonus.description} (unimplemented)`;

		const className = classList(
			'tooltip-set-bonus',
			i === 0 && 'tooltip-set-bonus-start',
			equippedCount >= bonus.piecesRequired && 'is-active',
			!bonus.isImplemented && 'is-unimplemented',
		);

		lines.push(multilineTextLine(className, text));
	}

	return lines.join('');
}

export function renderItemTooltip(item: UIItem, tooltipContext?: Partial<ItemTooltipContext>): string {
	const tooltipItem = item as TooltipRenderableItem;
	const renderParams = tooltipContext as ItemTooltipRenderContext | undefined;
	const lines: string[] = [];

	lines.push(renderItemHeader(tooltipItem));

	const displayedIlvl = getDisplayedItemLevel(tooltipItem, renderParams);
	if (displayedIlvl > 0) {
		lines.push(textLine('tooltip-item-level', `Item Level ${displayedIlvl}`));
	}

	const bindLabel = itemBindTypeNames.get(tooltipItem.bindType);
	if (bindLabel) {
		lines.push(textLine('tooltip-meta', bindLabel));
	}

	if (tooltipItem.unique) {
		lines.push(textLine('tooltip-meta', 'Unique'));
	}

	const leftType = getLeftTypeLabel(tooltipItem);
	const rightType = getRightTypeLabel(tooltipItem);
	if (leftType) {
		lines.push(pairedTextLine('tooltip-meta tooltip-type-line', leftType, rightType));
	}

	const armorLine = getArmorLine(tooltipItem);
	if (armorLine) {
		lines.push(textLine('tooltip-stat', armorLine));
	}

	lines.push(...getWeaponSectionLines(tooltipItem));

	for (const statLine of getDirectStatLines(tooltipItem, renderParams?.playerClass)) {
		lines.push(textLine('tooltip-stat', statLine));
	}

	if (tooltipItem.tooltipResolvedEnchantText) {
		lines.push(textLine('tooltip-effect', tooltipItem.tooltipResolvedEnchantText));
	}

	const classRestrictionLine = getClassRestrictionLine(tooltipItem, renderParams);
	if (classRestrictionLine) {
		lines.push(classRestrictionLine);
	}

	const professionRestrictionLine = getProfessionRestrictionLine(tooltipItem, renderParams);
	if (professionRestrictionLine) {
		lines.push(professionRestrictionLine);
	}

	const factionRestrictionLine = getFactionRestrictionLine(tooltipItem);
	if (factionRestrictionLine) {
		lines.push(factionRestrictionLine);
	}

	const requiredLevelLine = getItemRequiredLevelLine(tooltipItem, renderParams);
	if (requiredLevelLine) {
		lines.push(requiredLevelLine);
	}

	if (item.maxDurability) {
		lines.push(textLine('tooltip-meta', `Durability ${item.maxDurability}/${item.maxDurability}`));
	}

	for (const effectLine of getEffectStatLines(tooltipItem)) {
		lines.push(textLine('tooltip-effect', effectLine));
	}

	for (const effectLine of tooltipItem.tooltipResolvedEffectLines ?? []) {
		const text = effectLine.isImplemented
			? effectLine.text
			: `${effectLine.text} (unimplemented)`;

		lines.push(
			textLine(
				classList('tooltip-effect', !effectLine.isImplemented && 'is-unimplemented'),
				text,
			),
		);
	}

	const hasSetSection =
		Boolean(tooltipItem.tooltipResolvedSetName) ||
		Boolean(tooltipItem.setId) ||
		(tooltipItem.tooltipSetItems?.length ?? 0) > 0 ||
		(tooltipItem.tooltipResolvedSetBonuses?.length ?? 0) > 0;

	if (hasSetSection) {
		const setSection = renderSetSection(tooltipItem, renderParams);
		if (setSection) {
			lines.push(setSection);
		}
	}

	return `<div class="wow-tooltip">${lines.join('')}</div>`;
}

export function renderSpellTooltip(spell: UISpell, _tooltipContext?: Partial<SpellTooltipContext>): string {
	const lines: string[] = [];

	lines.push(`
        <div class="tooltip-header">
            <div class="tooltip-name is-common">${escapeHtml(spell.name || `Spell ${spell.id}`)}</div>
            ${spell.subtext ? `<div class="tooltip-subtext">${escapeHtml(spell.subtext)}</div>` : ''}
        </div>
    `);

	const costText = getSpellCostText(spell);
	const rangeText = getSpellRangeText(spell);
	if (costText || rangeText) {
		lines.push(pairedTextLine('tooltip-meta', costText ?? '', rangeText ?? ''));
	}

	const castText = getSpellCastText(spell);
	const cooldownText = getSpellCooldownText(spell);
	if (castText || cooldownText) {
		lines.push(pairedTextLine('tooltip-meta', castText ?? '', cooldownText ?? ''));
	}

	if (spell.requiredLevel > 0) {
		lines.push(textLine('tooltip-req', `Requires level ${spell.requiredLevel}`));
	}

	if (spell.description) {
		lines.push(multilineTextLine('tooltip-description', spell.description));
	}

	return `<div class="wow-tooltip">${lines.join('')}</div>`;
}

export function renderLoadingTooltip(): string {
	return `
		<div class="wow-tooltip">
			<div class="tooltip-loading">Loading tooltip…</div>
		</div>
	`;
}

export function renderErrorTooltip(message = 'Failed to load tooltip.'): string {
	return `
		<div class="wow-tooltip">
			<div class="tooltip-error">${escapeHtml(message)}</div>
		</div>
	`;
}

export function renderSimpleTooltip(title: string, description?: string): string {
	return `
		<div class="wow-tooltip">
			<div class="tooltip-name is-common">${escapeHtml(title)}</div>
			${description ? `<div class="tooltip-description">${escapeHtml(description)}</div>` : ''}
		</div>
	`;
}

export async function getItemTooltipData(itemId: number, tooltipContext?: Partial<ItemTooltipContext>): Promise<UIItem | null> {
	try {
		const db = await Database.get();
		const item = db.getItemById(itemId);
		if (!item) {
			return null;
		}

		let tooltipResolvedEnchantText: string | undefined;
		if (tooltipContext?.enchantId && tooltipContext.enchantId > 0) {
			const enchant = db.getEnchantByEffectId(tooltipContext.enchantId);
			if (enchant) {
				try {
					tooltipResolvedEnchantText = await getEnchantDescription(enchant);
				} catch (error) {
					console.warn(`Failed to load enchant description for enchant ${enchant.effectId}:`, error);
					tooltipResolvedEnchantText = enchant.name;
				}
			}
		}

		const tooltipResolvedEffectLines = resolveItemEffectLines(item, db);
		const { tooltipResolvedSetName, tooltipSetItems, tooltipResolvedSetBonuses } = resolveItemSetData(item, db);

		return {
			...item,
			tooltipResolvedEnchantText,
			tooltipResolvedEffectLines,
			tooltipResolvedSetName,
			tooltipSetItems,
			tooltipResolvedSetBonuses,
		} as TooltipRenderableItem;
	} catch (error) {
		console.error(`Failed to load item tooltip data for item ${itemId}:`, error);
		return null;
	}
}

export async function getSpellTooltipData(spellId: number): Promise<UISpell | null> {
	try {
		const db = await Database.get();
		return db.getSpellById(spellId) ?? null;
	} catch {
		return null;
	}
}

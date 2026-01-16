import { getWowheadLanguagePrefix } from '../constants/lang';
import { MAX_CHARACTER_LEVEL } from '../constants/mechanics';
import { ResourceType } from '../proto/api';
import { ActionID as ActionIdProto, ItemRandomSuffix, OtherAction } from '../proto/common';
import { IconData, UIItem as Item } from '../proto/ui';
import { buildWowheadTooltipDataset, WowheadTooltipItemParams, WowheadTooltipSpellParams } from '../wowhead';
import { Database } from './database';

// Used to filter action IDs by level
export interface ActionIdConfig {
	id: number;
	minLevel?: number;
	maxLevel?: number;
}

// Uniquely identifies a specific item / spell / thing in WoW. This object is immutable.
export class ActionId {
	readonly itemId: number;
	readonly randomSuffixId: number;
	readonly spellId: number;
	readonly otherId: OtherAction;
	readonly tag: number;
	readonly rank: number;

	readonly baseName: string; // The name without any tag additions.
	readonly name: string;
	readonly iconUrl: string;
	readonly spellIdTooltipOverride: number | null;

	private constructor(
		itemId: number,
		spellId: number,
		otherId: OtherAction,
		tag: number,
		baseName: string,
		name: string,
		iconUrl: string,
		rank: number,
		randomSuffixId?: number,
	) {
		this.itemId = itemId;
		this.randomSuffixId = randomSuffixId || 0;
		this.spellId = spellId;
		this.otherId = otherId;
		(this.rank = rank), (this.tag = tag);

		switch (otherId) {
			case OtherAction.OtherActionNone:
				break;
			case OtherAction.OtherActionWait:
				baseName = 'Wait';
				iconUrl = '/classic/assets/icons/inv_misc_pocketwatch_01.jpg';
				break;
			case OtherAction.OtherActionManaRegen:
				name = 'Mana Tick';
				iconUrl = resourceTypeToIcon[ResourceType.ResourceTypeMana];
				if (tag === 1) {
					name += ' (Casting)';
				} else if (tag === 2) {
					name += ' (Not Casting)';
				}
				break;
			case OtherAction.OtherActionEnergyRegen:
				baseName = 'Energy Tick';
				iconUrl = resourceTypeToIcon[ResourceType.ResourceTypeEnergy];
				break;
			case OtherAction.OtherActionComboPoints:
				baseName = 'Combo Point Gain';
				iconUrl = resourceTypeToIcon[ResourceType.ResourceTypeComboPoints];
				break;
			case OtherAction.OtherActionFocusRegen:
				baseName = 'Focus Tick';
				iconUrl = resourceTypeToIcon[ResourceType.ResourceTypeFocus];
				break;
			case OtherAction.OtherActionManaGain:
				baseName = 'Mana Gain';
				iconUrl = resourceTypeToIcon[ResourceType.ResourceTypeMana];
				break;
			case OtherAction.OtherActionRageGain:
				baseName = 'Rage Gain';
				iconUrl = resourceTypeToIcon[ResourceType.ResourceTypeRage];
				break;
			case OtherAction.OtherActionAttack:
				name = 'Melee';
				iconUrl = '/classic/assets/icons/inv_sword_04.jpg';
				if (tag === 1) {
					name += ' (Main-Hand)';
				} else if (tag === 2) {
					name += ' (Off-Hand)';
				} else if (tag === 3) {
					name += ' (Extra Attack)';
				}
				break;
			case OtherAction.OtherActionShoot:
				name = 'Shoot';
				iconUrl = '/classic/assets/icons/ability_marksmanship.jpg';
				if (tag === 3) {
					name += ' (Extra Attack)';
				}
				break;
			case OtherAction.OtherActionMove:
				name = 'Move';
				iconUrl = '/classic/assets/icons/inv_boots_02.jpg';
				break;
			case OtherAction.OtherActionPet:
				break;
			case OtherAction.OtherActionRefund:
				baseName = 'Refund';
				iconUrl = '/classic/assets/icons/inv_misc_coin_01.jpg';
				break;
			case OtherAction.OtherActionDamageTaken:
				baseName = 'Damage Taken';
				iconUrl = '/classic/assets/icons/inv_sword_04.jpg';
				break;
			case OtherAction.OtherActionHealingModel:
				baseName = 'Incoming HPS';
				iconUrl = '/classic/assets/icons/spell_holy_renew.jpg';
				break;
			case OtherAction.OtherActionPotion:
				baseName = 'Potion';
				iconUrl = '/classic/assets/icons/trade_alchemy.jpg';
				break;
			case OtherAction.OtherActionExplosives:
				baseName = 'Explosive';
				iconUrl = '/classic/assets/icons/inv_misc_bomb_06.jpg';
				break;
			case OtherAction.OtherActionOffensiveEquip:
				baseName = 'Offensive Equipment';
				iconUrl = '/classic/assets/icons/inv_trinket_naxxramas03.jpg';
				break;
			case OtherAction.OtherActionDefensiveEquip:
				baseName = 'Defensive Equipment';
				iconUrl = '/classic/assets/icons/inv_trinket_naxxramas05.jpg';
				break;
		}
		this.baseName = baseName;
		this.name = name || baseName;
		this.iconUrl = iconUrl;
		this.spellIdTooltipOverride = this.spellTooltipOverride?.spellId || null;
		if (this.name) this.name += rank ? ` (Rank ${rank})` : '';
	}

	anyId(): number {
		return this.itemId || this.spellId || this.otherId;
	}

	equals(other: ActionId): boolean {
		return this.equalsIgnoringTag(other) && this.tag === other.tag;
	}

	equalsIgnoringTag(other: ActionId): boolean {
		return this.itemId === other.itemId && this.randomSuffixId === other.randomSuffixId && this.spellId === other.spellId && this.otherId === other.otherId;
	}

	setBackground(elem: HTMLElement) {
		if (this.iconUrl) {
			elem.style.backgroundImage = `url('${this.iconUrl}')`;
		}
	}

	static makeItemUrl(id: number, randomSuffixId?: number): string {
		return `https://database.turtlecraft.gg/?item=${id}`;
	}
	static makeSpellUrl(id: number): string {
		return `https://database.turtlecraft.gg/?spell=${id}`;
	}
	static async makeItemTooltipData(id: number, params?: Omit<WowheadTooltipItemParams, 'itemId'>) {
		return buildWowheadTooltipDataset({ itemId: id, ...params });
	}
	static async makeSpellTooltipData(id: number, params?: Omit<WowheadTooltipSpellParams, 'spellId'>) {
		return buildWowheadTooltipDataset({ spellId: id, ...params });
	}
	static makeQuestUrl(id: number): string {
		return `https://database.turtlecraft.gg/?quest=${id}`;
	}
	static makeNpcUrl(id: number): string {
		return `https://database.turtlecraft.gg/?npc=${id}`;
	}
	static makeZoneUrl(id: number): string {
		const langPrefix = getWowheadLanguagePrefix();
		return `https://wowhead.com/classic/${langPrefix}zone=${id}`;
	}

	setWowheadHref(elem: HTMLAnchorElement) {
		if (this.itemId) {
			elem.href = ActionId.makeItemUrl(this.itemId, this.randomSuffixId);
		} else if (this.spellId) {
			elem.href = ActionId.makeSpellUrl(this.spellIdTooltipOverride || this.spellId);
		}
	}

	async setWowheadDataset(elem: HTMLElement, params?: Omit<WowheadTooltipItemParams, 'itemId'> | Omit<WowheadTooltipSpellParams, 'spellId'>) {
		(this.itemId
			? ActionId.makeItemTooltipData(this.itemId, params)
			: ActionId.makeSpellTooltipData(this.spellIdTooltipOverride || this.spellId, params)
		).then(url => {
			if (elem) elem.dataset.wowhead = url;
		});
	}

	setBackgroundAndHref(elem: HTMLAnchorElement) {
		this.setBackground(elem);
		this.setWowheadHref(elem);
	}

	async fillAndSet(elem: HTMLAnchorElement, setHref: boolean, setBackground: boolean): Promise<ActionId> {
		const filled = await this.fill();
		if (setHref) {
			filled.setWowheadHref(elem);
		}
		if (setBackground) {
			filled.setBackground(elem);
		}
		return filled;
	}

	// Returns an ActionId with the name and iconUrl fields filled.
	// playerIndex is the optional index of the player to whom this ID corresponds.
	async fill(playerIndex?: number): Promise<ActionId> {
		if (this.name || this.iconUrl) {
			return this;
		}

		if (this.otherId) {
			return this;
		}

		const tooltipData = await ActionId.getTooltipData(this);

		const baseName = tooltipData['name'];
		let name = baseName;
		switch (baseName) {
			case 'Master Demonologist':
				switch (this.tag) {
					case 1:
						name = `${name} (Imp)`;
						break;
					case 2:
						name = `${name} (Voidwalker)`;
						break;
					case 3:
						name = `${name} (Succubus)`;
						break;
					case 4:
						name = `${name} (Felhunter)`;
						break;
				}
				break;
			case 'Berserking':
				if (this.tag !== 0) name = `${name} (${this.tag * 5}%)`;
				break;
			// Burn Spells
			case 'Fireball':
			case 'Pyroblast':
			case 'Flame Shock':
				if (this.tag === 1) name = `${name} (DoT)`;
				break;
			// Channeled Tick Spells
			case 'Evocation':
			case 'Mind Flay':
				if (this.tag > 0) name = `${name} (${this.tag} Tick)`;
				break;
			// Combo Point Spenders
			case 'Eviscerate':
			case 'Expose Armor':
			case 'Rupture':
			case 'Slice and Dice':
				if (this.tag) name += ` (${this.tag} CP)`;
				break;
			case 'Deadly Poison':
			case 'Deadly Poison II':
			case 'Deadly Poison III':
			case 'Deadly Poison IV':
			case 'Deadly Poison V':
			case 'Instant Poison':
			case 'Instant Poison II':
			case 'Instant Poison III':
			case 'Instant Poison IV':
			case 'Instant Poison V':
			case 'Instant Poison VI':
			case 'Wound Poison':
				if (this.tag === 1) {
					name += ' (Shiv)';
				} else if (this.tag === 2) {
					name += ' (Deadly Brew)';
				} else if (this.tag === 100) {
					name += ' (Tick)';
				}
				break;
			// Dual-hit MH/OH spells and weapon imbues
			case 'Holy Strength': // Weapon - Crusader Enchant
				if (this.tag === 1) {
					name = `${name} (Main-Hand)`;
				} else if (this.tag === 2) {
					name = `${name} (Off-Hand)`;
				}
				break;
			case 'Holy Shield':
				if (this.tag === 1) {
					name += ' (Proc)';
				}
				break;
			// For targetted buffs, tag is the source player's raid index or -1 if none.
			case 'Innervate':
			case 'Mana Tide Totem':
			case 'Power Infusion':
				if (this.tag !== -1) {
					if (this.tag === playerIndex || playerIndex === undefined) {
						name += ` (self)`;
					} else {
						name += ` (from #${this.tag + 1})`;
					}
				} else {
					name += ' (raid)';
				}
				break;
			case 'Battle Shout':
				if (this.tag === 1) {
					name += ' (Snapshot)';
				}
				break;
			case 'Heroic Strike':
			case 'Cleave':
			case 'Maul':
				if (this.tag === 1) {
					name += ' (Queue)';
				}
				break;
			case 'Raptor Strike':
				if (this.tag === 1) name = `${name} (Hit)`;
				else if (this.tag === 3) name = `${name} (Queue)`;
				break;
			case 'Thunderfury':
				if (this.tag === 1) {
					name += ' (Main)';
				} else if (this.tag === 2) {
					name += ' (Bounce)';
				}
				break;
			case 'Power of the Guardian':
				switch (this.spellId) {
					case 28142:
						name = `${name} (Mage)`;
						break;
					case 28143:
						name = `${name} (Warlock)`;
						break;
					case 28144:
						name = `${name} (Priest)`;
						break;
					case 28145:
						name = `${name} (Druid)`;
						break;
				}
				break;
			default:
				if (this.tag) {
					name += ' (??)';
				}
				break;
		}

		const iconUrl = ActionId.makeIconUrl(tooltipData['icon']);

		return new ActionId(this.itemId, this.spellId, this.otherId, this.tag, baseName, name, iconUrl, this.rank || tooltipData.rank, this.randomSuffixId);
	}

	toString(): string {
		return this.toStringIgnoringTag() + (this.tag ? '-' + this.tag : '');
	}

	toStringIgnoringTag(): string {
		if (this.itemId) {
			return 'item-' + this.itemId;
		} else if (this.spellId) {
			return 'spell-' + this.spellId;
		} else if (this.otherId) {
			return 'other-' + this.otherId;
		} else {
			throw new Error('Empty action id!');
		}
	}

	toProto(): ActionIdProto {
		const protoId = ActionIdProto.create({
			tag: this.tag,
		});

		if (this.itemId) {
			protoId.rawId = {
				oneofKind: 'itemId',
				itemId: this.itemId,
			};
		} else if (this.spellId) {
			protoId.rawId = {
				oneofKind: 'spellId',
				spellId: this.spellId,
			};
			protoId.rank = this.rank;
		} else if (this.otherId) {
			protoId.rawId = {
				oneofKind: 'otherId',
				otherId: this.otherId,
			};
		}

		return protoId;
	}

	toProtoString(): string {
		return ActionIdProto.toJsonString(this.toProto());
	}

	withoutTag(): ActionId {
		return new ActionId(this.itemId, this.spellId, this.otherId, 0, this.baseName, this.baseName, this.iconUrl, this.rank, this.randomSuffixId);
	}

	static fromEmpty(): ActionId {
		return new ActionId(0, 0, OtherAction.OtherActionNone, 0, '', '', '', 0);
	}

	static fromItemId(itemId: number, tag?: number, randomSuffixId?: number): ActionId {
		return new ActionId(itemId, 0, OtherAction.OtherActionNone, tag || 0, '', '', '', 0, randomSuffixId || 0);
	}

	static fromSpellId(spellId: number, rank = 0, tag?: number): ActionId {
		return new ActionId(0, spellId, OtherAction.OtherActionNone, tag || 0, '', '', '', rank);
	}

	static fromOtherId(otherId: OtherAction, tag?: number): ActionId {
		return new ActionId(0, 0, otherId, tag || 0, '', '', '', 0);
	}

	static fromPetName(petName: string): ActionId {
		return petNameToActionId[petName] || new ActionId(0, 0, OtherAction.OtherActionPet, 0, petName, petName, petNameToIcon[petName] || '', 0);
	}

	static fromItem(item: Item): ActionId {
		return ActionId.fromItemId(item.id);
	}

	static fromRandomSuffix(item: Item, randomSuffix: ItemRandomSuffix): ActionId {
		return ActionId.fromItemId(item.id, 0, randomSuffix.id);
	}

	static fromProto(protoId: ActionIdProto): ActionId {
		if (protoId.rawId.oneofKind === 'spellId') {
			return ActionId.fromSpellId(protoId.rawId.spellId, protoId.rank, protoId.tag);
		} else if (protoId.rawId.oneofKind === 'itemId') {
			return ActionId.fromItemId(protoId.rawId.itemId, protoId.tag);
		} else if (protoId.rawId.oneofKind === 'otherId') {
			return ActionId.fromOtherId(protoId.rawId.otherId, protoId.tag);
		} else {
			return ActionId.fromEmpty();
		}
	}

	private static readonly logRegex = /{((SpellID)|(ItemID)|(OtherID)): (\d+)(, Tag: (-?\d+))?}/;
	private static readonly logRegexGlobal = new RegExp(ActionId.logRegex, 'g');
	private static fromMatch(match: RegExpMatchArray): ActionId {
		const idType = match[1];
		const id = parseInt(match[5]);
		return new ActionId(
			idType === 'ItemID' ? id : 0,
			idType === 'SpellID' ? id : 0,
			idType === 'OtherID' ? id : 0,
			match[7] ? parseInt(match[7]) : 0,
			'',
			'',
			'',
			0,
		);
	}
	static fromLogString(str: string): ActionId {
		const match = str.match(ActionId.logRegex);
		if (match) {
			return ActionId.fromMatch(match);
		} else {
			console.warn('Failed to parse action id from log: ' + str);
			return ActionId.fromEmpty();
		}
	}

	static async replaceAllInString(str: string): Promise<string> {
		const matches = [...str.matchAll(ActionId.logRegexGlobal)];

		const replaceData = await Promise.all(
			matches.map(async match => {
				const actionId = ActionId.fromMatch(match);
				const filledId = await actionId.fill();
				return {
					firstIndex: match.index || 0,
					len: match[0].length,
					actionId: filledId,
				};
			}),
		);

		// Loop in reverse order so we can greedily apply the string replacements.
		for (let i = replaceData.length - 1; i >= 0; i--) {
			const data = replaceData[i];
			str = str.substring(0, data.firstIndex) + data.actionId.name + str.substring(data.firstIndex + data.len);
		}

		return str;
	}

	private static makeIconUrl(iconLabel: string): string {
		return `/classic/assets/icons/${iconLabel}.jpg`;
	}

	static async getTooltipData(actionId: ActionId): Promise<IconData> {
		if (actionId.itemId) {
			return await Database.getItemIconData(actionId.itemId);
		} else {
			return await Database.getSpellIconData(actionId.spellId);
		}
	}
	get spellIconOverride(): ActionId | null {
		const override = spellIdIconOverrides.get(JSON.stringify({ spellId: this.spellId }));
		if (!override) return null;
		return override.itemId ? ActionId.fromItemId(override.itemId) : ActionId.fromItemId(override.spellId!);
	}

	get spellTooltipOverride(): ActionId | null {
		const override = spellIdTooltipOverrides.get(JSON.stringify({ spellId: this.spellId, tag: this.tag }));
		if (!override) return null;
		return override.itemId ? ActionId.fromItemId(override.itemId) : ActionId.fromSpellId(override.spellId!);
	}
}

type ActionIdOverride = { itemId?: number; spellId?: number };

// Some items/spells have weird icons, so use this to show a different icon instead.
const spellIdIconOverrides: Map<string, ActionIdOverride> = new Map([
	[JSON.stringify({ spellId: 449288 }), { itemId: 221309 }], // Darkmoon Card: Sandstorm
	[JSON.stringify({ spellId: 455864 }), { spellId: 9907 }], // Tier 1 Balance Druid "Improved Faerie Fire"
	[JSON.stringify({ spellId: 457544 }), { spellId: 10408 }], // Tier 1 Shaman Tank "Improved Stoneskin / Windwall Totem"
]);

const spellIdTooltipOverrides: Map<string, ActionIdOverride> = new Map([]);

const spellIDsToShowBuffs = new Set([
	702, // https://database.turtlecraft.gg/?spell=702
	704, // https://database.turtlecraft.gg/?spell=704
	770, // https://database.turtlecraft.gg/?spell=770
	778, // https://database.turtlecraft.gg/?spell=778
	1108, // https://database.turtlecraft.gg/?spell=1108
	1490, // https://database.turtlecraft.gg/?spell=1490
	6205, // https://database.turtlecraft.gg/?spell=6205
	7646, // https://database.turtlecraft.gg/?spell=7646
	7658, // https://database.turtlecraft.gg/?spell=7658
	7659, // https://database.turtlecraft.gg/?spell=7659
	9749, // https://database.turtlecraft.gg/?spell=9749
	9907, // https://database.turtlecraft.gg/?spell=9907
	11707, // https://database.turtlecraft.gg/?spell=11707
	11708, // https://database.turtlecraft.gg/?spell=11708
	11717, // https://database.turtlecraft.gg/?spell=11717
	11721, // https://database.turtlecraft.gg/?spell=11721
	11722, // https://database.turtlecraft.gg/?spell=11722
	14201, // https://database.turtlecraft.gg/?spell=14201
	16257, // https://database.turtlecraft.gg/?spell=16257
	16277, // https://database.turtlecraft.gg/?spell=16277
	16278, // https://database.turtlecraft.gg/?spell=16278
	16279, // https://database.turtlecraft.gg/?spell=16279
	16280, // https://database.turtlecraft.gg/?spell=16280
	17862, // https://database.turtlecraft.gg/?spell=17862
	17937, // https://database.turtlecraft.gg/?spell=17937
	18789, // https://database.turtlecraft.gg/?spell=18789
	18790, // https://database.turtlecraft.gg/?spell=18790
	18791, // https://database.turtlecraft.gg/?spell=18791
	18792, // https://database.turtlecraft.gg/?spell=18792
	20186, // https://database.turtlecraft.gg/?spell=20186
	20300, // https://database.turtlecraft.gg/?spell=20300
	20355, // https://database.turtlecraft.gg/?spell=20355
	20301, // https://database.turtlecraft.gg/?spell=20301
	20302, // https://database.turtlecraft.gg/?spell=20302
	20303, // https://database.turtlecraft.gg/?spell=20303
	23060, // https://database.turtlecraft.gg/?spell=23060
	23736, // https://database.turtlecraft.gg/?spell=23736
	23737, // https://database.turtlecraft.gg/?spell=23737
	23738, // https://database.turtlecraft.gg/?spell=23738
	23766, // https://database.turtlecraft.gg/?spell=23766
	23768, // https://database.turtlecraft.gg/?spell=23768
	24907, // https://database.turtlecraft.gg/?spell=24907
	24932, // https://database.turtlecraft.gg/?spell=24932
	402808, // https://database.turtlecraft.gg/?spell=402808
	425415, // https://database.turtlecraft.gg/?spell=425415
	461252, // https://database.turtlecraft.gg/?spell=461252
	461270, // https://database.turtlecraft.gg/?spell=461270
	1214279, // https://database.turtlecraft.gg/?spell=1214279
]);

export const defaultTargetIcon = '/classic/assets/icons/spell_shadow_metamorphosis.jpg';

const petNameToActionId: Record<string, ActionId> = {
	'Eye of the Void': ActionId.fromSpellId(402789),
	'Frozen Orb 1': ActionId.fromSpellId(440802),
	'Frozen Orb 2': ActionId.fromSpellId(440802),
	Homunculi: ActionId.fromSpellId(402799),
	Shadowfiend: ActionId.fromSpellId(401977),
};

// https://wowhead.com/classic/hunter-pets
const petNameToIcon: Record<string, string> = {
	Bat: '/classic/assets/icons/ability_hunter_pet_bat.jpg',
	Bear: '/classic/assets/icons/ability_hunter_pet_bear.jpg',
	'Bird of Prey': '/classic/assets/icons/ability_hunter_pet_owl.jpg',
	Boar: '/classic/assets/icons/ability_hunter_pet_boar.jpg',
	'Carrion Bird': '/classic/assets/icons/ability_hunter_pet_vulture.jpg',
	Cat: '/classic/assets/icons/ability_hunter_pet_cat.jpg',
	Chimaera: '/classic/assets/icons/ability_hunter_pet_chimera.jpg',
	'Core Hound': '/classic/assets/icons/ability_hunter_pet_corehound.jpg',
	Crab: '/classic/assets/icons/ability_hunter_pet_crab.jpg',
	Crocolisk: '/classic/assets/icons/ability_hunter_pet_crocolisk.jpg',
	Devilsaur: '/classic/assets/icons/ability_hunter_pet_devilsaur.jpg',
	Dragonhawk: '/classic/assets/icons/ability_hunter_pet_dragonhawk.jpg',
	'Emerald Dragon Whelp': '/classic/assets/icons/inv_misc_head_dragon_green.jpg',
	Eskhandar: '/classic/assets/icons/inv_misc_head_tiger_01.jpg',
	Felguard: '/classic/assets/icons/spell_shadow_summonfelguard.jpg',
	Felhunter: '/classic/assets/icons/spell_shadow_summonfelhunter.jpg',
	'Spirit Wolves': '/classic/assets/icons/spell_shaman_feralspirit.jpg',
	Infernal: '/classic/assets/icons/spell_shadow_summoninfernal.jpg',
	Gorilla: '/classic/assets/icons/ability_hunter_pet_gorilla.jpg',
	Hyena: '/classic/assets/icons/ability_hunter_pet_hyena.jpg',
	Imp: '/classic/assets/icons/spell_shadow_summonimp.jpg',
	'Mirror Image': '/classic/assets/icons/spell_magic_lesserinvisibilty.jpg',
	Moth: '/classic/assets/icons/ability_hunter_pet_moth.jpg',
	'Nether Ray': '/classic/assets/icons/ability_hunter_pet_netherray.jpg',
	Owl: '/classic/assets/icons/ability_hunter_pet_owl.jpg',
	Raptor: '/classic/assets/icons/ability_hunter_pet_raptor.jpg',
	Ravager: '/classic/assets/icons/ability_hunter_pet_ravager.jpg',
	Rhino: '/classic/assets/icons/ability_hunter_pet_rhino.jpg',
	Scorpid: '/classic/assets/icons/ability_hunter_pet_scorpid.jpg',
	Serpent: '/classic/assets/icons/spell_nature_guardianward.jpg',
	Silithid: '/classic/assets/icons/ability_hunter_pet_silithid.jpg',
	Spider: '/classic/assets/icons/ability_hunter_pet_spider.jpg',
	'Spirit Beast': '/classic/assets/icons/ability_druid_primalprecision.jpg',
	'Spore Bat': '/classic/assets/icons/ability_hunter_pet_sporebat.jpg',
	Succubus: '/classic/assets/icons/spell_shadow_summonsuccubus.jpg',
	Tallstrider: '/classic/assets/icons/ability_hunter_pet_tallstrider.jpg',
	Treants: '/classic/assets/icons/ability_druid_forceofnature.jpg',
	Turtle: '/classic/assets/icons/ability_hunter_pet_turtle.jpg',
	Voidwalker: '/classic/assets/icons/spell_shadow_summonvoidwalker.jpg',
	'Warp Stalker': '/classic/assets/icons/ability_hunter_pet_warpstalker.jpg',
	Wasp: '/classic/assets/icons/ability_hunter_pet_wasp.jpg',
	'Wind Serpent': '/classic/assets/icons/ability_hunter_pet_windserpent.jpg',
	Wolf: '/classic/assets/icons/ability_hunter_pet_wolf.jpg',
	Worm: '/classic/assets/icons/ability_hunter_pet_worm.jpg',
};

export function getPetIconFromName(name: string): string | ActionId | undefined {
	return petNameToActionId[name] || petNameToIcon[name];
}

export const resourceTypeToIcon: Record<ResourceType, string> = {
	[ResourceType.ResourceTypeNone]: '',
	[ResourceType.ResourceTypeHealth]: '/classic/assets/icons/inv_potion_01.jpg',
	[ResourceType.ResourceTypeMana]: '/classic/assets/icons/inv_potion_02.jpg',
	[ResourceType.ResourceTypeEnergy]: '/classic/assets/icons/spell_shadow_shadowworddominate.jpg',
	[ResourceType.ResourceTypeRage]: '/classic/assets/icons/inv_potion_03.jpg',
	[ResourceType.ResourceTypeComboPoints]: '/classic/assets/icons/inv_sword_04.jpg',
	[ResourceType.ResourceTypeFocus]: '/classic/assets/icons/ability_hunter_aimedshot.jpg',
};

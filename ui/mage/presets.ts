import { Player } from '../core/player.js';
import * as PresetUtils from '../core/preset_utils';
import {
	Alcohol,
	BlastedLandsBuff,
	Conjured,
	Consumes,
	Debuffs,
	FirePowerBuff,
	Flask,
	Food,
	FrostPowerBuff,
	HealthElixir,
	IndividualBuffs,
	MiscConsumes,
	PartyBuffs,
	Potions,
	Profession,
	Race,
	RaidBuffs,
	SapperExplosive,
	SaygesFortune,
	SpellPowerBuff,
	Stat,
	TristateEffect,
	WeaponImbue,
	ZanzaBuff,
} from '../core/proto/common';
import { Mage_Options as MageOptions, Mage_Options_ArmorType as ArmorType } from '../core/proto/mage';
import { SavedTalents } from '../core/proto/ui';
import { Stats } from '../core/proto_utils/stats.js';
import { TypedEvent } from '../core/typed_event.js';
import ArcaneAoEAPL from './apls/arcane_aoe.apl.json';
import ArcaneSTAPL from './apls/arcane_st.apl.json';
import FireST from './apls/fire_st.apl.json';
import ArcaneAoEGearBiS from './gear_sets/arcane_aoe_bis.gear.json';
import ArcaneSTGearBiS from './gear_sets/arcane_st_bis.gear.json';
import FireSTBiS from './gear_sets/fire_st_bis.gear.json';

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////
export const DefaultBossEncounter = 'Default/Raid Boss (1 Lvl 63 Demon)';
export const DefaultTrashEncounter = 'Default/Trash (5 lvl 60 Undead)';

export const GearArcaneSTBiS = PresetUtils.makePresetGear('Arcane ST BiS', ArcaneSTGearBiS);
export const GearArcaneAoEBiS = PresetUtils.makePresetGear('Arcane AoE BiS', ArcaneAoEGearBiS);
export const GearFireSTBiS = PresetUtils.makePresetGear('Fire ST BiS', FireSTBiS);

export const GearPresets = [GearArcaneSTBiS, GearArcaneAoEBiS, GearFireSTBiS];

export const DefaultGear = GearArcaneSTBiS;

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const ROTATION_PRESET_ARCANE_ST = PresetUtils.makePresetAPLRotation('Arcane ST', ArcaneSTAPL, {});
export const ROTATION_PRESET_ARCANE_AOE = PresetUtils.makePresetAPLRotation('Arcane AoE', ArcaneAoEAPL, {});
export const ROTATION_PRESET_FIRE_ST = PresetUtils.makePresetAPLRotation('Fire ST', FireST, {});

export const APLPresets = [ROTATION_PRESET_ARCANE_ST, ROTATION_PRESET_ARCANE_AOE, ROTATION_PRESET_FIRE_ST];

export const DefaultAPL = ROTATION_PRESET_ARCANE_ST;

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the talent string format (numeric), make the talents on
// a talent calculator and copy the numbers in the url.

export const TalentsArcane = PresetUtils.makePresetTalents('Arcane', SavedTalents.create({ talentsString: '2350550310033311251-50003' }));
export const TalentsFire = PresetUtils.makePresetTalents('Fire', SavedTalents.create({ talentsString: '2300020000000000000-50523231230313251-003' }));

export const TalentPresets = [TalentsArcane, TalentsFire];

export const DefaultTalents = TalentsArcane;

///////////////////////////////////////////////////////////////////////////
//                                 EP Presets
///////////////////////////////////////////////////////////////////////////

export const EPArcaneST = PresetUtils.makePresetEpWeights('Arcane ST', Stats.fromMap(
	{
		[Stat.StatIntellect]: 0.32,
		[Stat.StatSpellPower]: 1,
		[Stat.StatSpellDamage]: 1,
		[Stat.StatArcanePower]: 1,
		[Stat.StatFirePower]: 0,
		[Stat.StatFrostPower]: 0,
		[Stat.StatSpellHit]: 0,
		[Stat.StatSpellCrit]: 14.97,
		[Stat.StatSpellHaste]: 21.81,
		[Stat.StatFortune]: 0.92,
	},
	{},
))

export const EPArcaneAoE = PresetUtils.makePresetEpWeights('Arcane AoE', Stats.fromMap(
	{
		[Stat.StatIntellect]: 0.48,
		[Stat.StatSpellPower]: 1,
		[Stat.StatSpellDamage]: 1,
		[Stat.StatArcanePower]: 1,
		[Stat.StatFirePower]: 0,
		[Stat.StatFrostPower]: 0,
		[Stat.StatSpellHit]: 0,
		[Stat.StatSpellCrit]: 22.62,
		[Stat.StatSpellHaste]: 0,
		[Stat.StatFortune]: 3.11,
	},
	{},
))

export const EPFireST = PresetUtils.makePresetEpWeights('Fire ST', Stats.fromMap(
	{
		[Stat.StatIntellect]: 0.59,
		[Stat.StatSpellPower]: 1,
		[Stat.StatSpellDamage]: 1,
		[Stat.StatArcanePower]: 0,
		[Stat.StatFirePower]: 1,
		[Stat.StatFrostPower]: 0,
		[Stat.StatSpellHit]: 25.39,
		[Stat.StatSpellCrit]: 27.05,
		[Stat.StatSpellHaste]: 13.26,
		[Stat.StatFortune]: 0.51,
	},
	{},
))

export const DefaultEP = EPArcaneST;

///////////////////////////////////////////////////////////////////////////
//                                Build Presets
///////////////////////////////////////////////////////////////////////////

export const PresetBuildArcaneST = PresetUtils.makePresetBuild('Arcane ST', {
	gear: GearArcaneSTBiS,
	talents: TalentsArcane,
	rotation: ROTATION_PRESET_ARCANE_ST,
	epWeights: EPArcaneST,
	encounter: PresetUtils.makePresetEncounter('Default/Raid Boss (1 Lvl 63 Demon)'),
});
export const PresetBuildArcaneAoE = PresetUtils.makePresetBuild('Arcane AoE', {
	gear: GearArcaneAoEBiS,
	talents: TalentsArcane,
	rotation: ROTATION_PRESET_ARCANE_AOE,
	epWeights: EPArcaneAoE,
	encounter: PresetUtils.makePresetEncounter('Default/Trash (5 lvl 60 Undead)'),
});
export const PresetBuildFireST = PresetUtils.makePresetBuild('Fire ST', {
	gear: GearFireSTBiS,
	talents: TalentsFire,
	rotation: ROTATION_PRESET_FIRE_ST,
	epWeights: EPFireST,
	encounter: PresetUtils.makePresetEncounter('Default/Raid Boss (1 Lvl 63 Demon)'),
});

export const BuildPresets = [PresetBuildArcaneST, PresetBuildArcaneAoE, PresetBuildFireST]


///////////////////////////////////////////////////////////////////////////
//                                 Options
///////////////////////////////////////////////////////////////////////////

export const DefaultOptions = MageOptions.create({
	armor: ArmorType.MageArmor,
});

export const DefaultConsumes = Consumes.create({
	defaultPotion: Potions.QuicknessPotion,
	defaultConjured: Conjured.ConjuredDemonicRune,
	nordanaarHerbalTea: true,
	flask: Flask.FlaskOfSupremePower,
	mainHandImbue: WeaponImbue.BrilliantWizardOil,
	food: Food.FoodDanonzosTelAbimMedley,
	alcohol: Alcohol.AlcoholMedivhMerlotBlue,
	healthElixir: HealthElixir.ElixirOfFortitude,
	dreamshardElixir: true,
	dreamtonic: true,
	spellPowerBuff: SpellPowerBuff.GreaterArcaneElixir,
	firePowerBuff: FirePowerBuff.ElixirOfGreaterFirepower,
	frostPowerBuff: FrostPowerBuff.ElixirOfGreaterFrostPower,
	elixirOfGreaterArcanePower: true,
	magebloodPotion: true,
	zanzaBuff: ZanzaBuff.SpiritOfZanza,
	blastedLandsBuff: BlastedLandsBuff.CerebralCortexCompound,
	miscConsumes: MiscConsumes.create({ jujuFlurry: true }),
	sapperExplosive: SapperExplosive.SapperGoblinSapper,
});

export const DefaultRaidBuffs = RaidBuffs.create({
	emeraldBlessing: true,
	arcaneBrilliance: true,
	divineSpirit: true,
	giftOfTheWild: TristateEffect.TristateEffectImproved,
	powerWordFortitude: TristateEffect.TristateEffectImproved,
	manaSpringTotem: TristateEffect.TristateEffectMissing,
});

export const DefaultPartyBuffs = PartyBuffs.create({
	atieshMage: true,
	atieshDruid: true,
	atieshWarlock: true,
	moonkinAura: true,
});

export const DefaultIndividualBuffs = IndividualBuffs.create({
	blessingOfKings: true,
	blessingOfWisdom: TristateEffect.TristateEffectImproved,
	rallyingCryOfTheDragonslayer: false,
	saygesFortune: SaygesFortune.SaygesUnknown,
	slipkiksSavvy: false,
	songflowerSerenade: false,
	spiritOfZandalar: false,
	warchiefsBlessing: false,
});

export const DefaultDebuffs = Debuffs.create({
	fireVulnerability: true,
	judgementOfWisdom: true,
	wintersChill: false,
	curseOfShadow: true,
	curseOfElements: true,
});

export const OtherDefaults = {
	distanceFromTarget: 20,
	profession1: Profession.Enchanting,
	profession2: Profession.Engineering,
	race: Race.RaceGnome,
	channelClipDelay: 50,
	infiniteMana: true,
};

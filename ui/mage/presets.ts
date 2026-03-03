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
	TristateEffect,
	WeaponImbue,
	ZanzaBuff,
} from '../core/proto/common';
import { Mage_Options as MageOptions, Mage_Options_ArmorType as ArmorType } from '../core/proto/mage';
import { SavedTalents } from '../core/proto/ui';
import DEFAULTAPL from './apls/default.apl.json';
import BISGear from './gear_sets/bis.gear.json';

///////////////////////////////////////////////////////////////////////////
//                                 Gear Presets
///////////////////////////////////////////////////////////////////////////

export const GearBIS = PresetUtils.makePresetGear('BiS', BISGear);

export const GearPresets = [GearBIS];

export const DefaultGear = GearBIS;

///////////////////////////////////////////////////////////////////////////
//                                 APL Presets
///////////////////////////////////////////////////////////////////////////

export const APLPDEFAULT = PresetUtils.makePresetAPLRotation('Arcane', DEFAULTAPL);

export const APLPresets = [APLPDEFAULT];

export const DefaultAPL = APLPDEFAULT;

///////////////////////////////////////////////////////////////////////////
//                                 Talent Presets
///////////////////////////////////////////////////////////////////////////

// Default talents. Uses the wowhead calculator format, make the talents on
// https://wowhead.com/talent-calc and copy the numbers in the url.

export const TalentsArcane = PresetUtils.makePresetTalents('Arcane', SavedTalents.create({ talentsString: '2350550310033311251-50003' }));

export const TalentPresets = [TalentsArcane];

export const DefaultTalents = TalentsArcane;

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
	improvedScorch: false,
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
};

import { HunterTalents } from '../proto/hunter.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import HunterTalentJson from './trees/hunter.json';

export const hunterTalentsConfig: TalentsConfig<HunterTalents> = newTalentsConfig(HunterTalentJson);

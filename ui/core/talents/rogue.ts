import { RogueTalents } from '../proto/rogue.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import RogueTalentJson from './trees/rogue.json';

export const rogueTalentsConfig: TalentsConfig<RogueTalents> = newTalentsConfig(RogueTalentJson);

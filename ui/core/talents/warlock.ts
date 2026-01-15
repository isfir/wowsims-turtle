import { WarlockTalents } from '../proto/warlock.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import WarlockTalentJson from './trees/warlock.json';

export const warlockTalentsConfig: TalentsConfig<WarlockTalents> = newTalentsConfig(WarlockTalentJson);

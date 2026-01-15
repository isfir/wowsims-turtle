import { PriestTalents } from '../proto/priest.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import PriestTalentJson from './trees/priest.json';

export const priestTalentsConfig: TalentsConfig<PriestTalents> = newTalentsConfig(PriestTalentJson);

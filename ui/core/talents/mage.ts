import { MageTalents } from '../proto/mage.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import MageTalentJson from './trees/mage.json';

export const mageTalentsConfig: TalentsConfig<MageTalents> = newTalentsConfig(MageTalentJson);

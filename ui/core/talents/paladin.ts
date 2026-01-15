import { PaladinTalents } from '../proto/paladin.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import PaladinTalentJson from './trees/paladin.json';

export const paladinTalentsConfig: TalentsConfig<PaladinTalents> = newTalentsConfig(PaladinTalentJson);

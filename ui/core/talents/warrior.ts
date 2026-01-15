import { WarriorTalents } from '../proto/warrior.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import WarriorTalentJson from './trees/warrior.json';

export const warriorTalentsConfig: TalentsConfig<WarriorTalents> = newTalentsConfig(WarriorTalentJson);

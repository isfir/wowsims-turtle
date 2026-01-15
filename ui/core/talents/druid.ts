import { DruidTalents } from '../proto/druid.js';
import { newTalentsConfig, TalentsConfig } from './talents_picker.js';
import DruidTalentsJson from './trees/druid.json';

export const druidTalentsConfig: TalentsConfig<DruidTalents> = newTalentsConfig(DruidTalentsJson);

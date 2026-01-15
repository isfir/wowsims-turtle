import { RaidSimSettings } from '../../../core/proto/ui';
import { RaidExporter } from '../raid_exporter';
import { RaidSimUI } from '../raid_sim_ui';

export class RaidJsonExporter extends RaidExporter {
	constructor(parent: HTMLElement, simUI: RaidSimUI) {
		super(parent, simUI, { title: 'JSON Export', allowDownload: true });
	}

	getData(): string {
		return JSON.stringify(RaidSimSettings.toJson(this.simUI.toProto()), null, 2);
	}
}

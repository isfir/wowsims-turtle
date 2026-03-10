import { ref } from 'tsx-vanilla';

import { setItemQualityCssClass } from '../../css_utils';
import { Player } from '../../player';
import { ItemSlot } from '../../proto/common';
import { ActionId } from '../../proto_utils/action_id';
import { getEnchantDescription } from '../../proto_utils/enchants';
import { EquippedItem } from '../../proto_utils/equipped_item';
import { slotNames } from '../../proto_utils/names.js';
import { itemTypeToSlotsMap } from '../../proto_utils/utils.js';
import { SimUI } from '../../sim_ui';
import { EventID } from '../../typed_event';
import { Component } from '../component';
import { TooltipManager } from '../tooltip_manager';
import { GearData } from './item_list';
import SelectorModal, { SelectorModalTabs } from './selector_modal';
import { getEmptySlotIconUrl } from './utils';

export default class GearPicker extends Component {
	// ItemSlot is used as the index
	readonly itemPickers: Array<ItemPicker>;
	readonly selectorModal: SelectorModal;

	constructor(parent: HTMLElement, simUI: SimUI, player: Player<any>) {
		super(parent, 'gear-picker-root');

		const leftSide = document.createElement('div');
		leftSide.classList.add('gear-picker-left', 'tab-panel-col');
		this.rootElem.appendChild(leftSide);

		const rightSide = document.createElement('div');
		rightSide.classList.add('gear-picker-right', 'tab-panel-col');
		this.rootElem.appendChild(rightSide);

		const leftItemPickers = [
			ItemSlot.ItemSlotHead,
			ItemSlot.ItemSlotNeck,
			ItemSlot.ItemSlotShoulder,
			ItemSlot.ItemSlotBack,
			ItemSlot.ItemSlotChest,
			ItemSlot.ItemSlotWrist,
			ItemSlot.ItemSlotMainHand,
			ItemSlot.ItemSlotOffHand,
			ItemSlot.ItemSlotRanged,
		].map(slot => new ItemPicker(leftSide, this, simUI, player, slot));

		const rightItemPickers = [
			ItemSlot.ItemSlotHands,
			ItemSlot.ItemSlotWaist,
			ItemSlot.ItemSlotLegs,
			ItemSlot.ItemSlotFeet,
			ItemSlot.ItemSlotFinger1,
			ItemSlot.ItemSlotFinger2,
			ItemSlot.ItemSlotTrinket1,
			ItemSlot.ItemSlotTrinket2,
		].map(slot => new ItemPicker(rightSide, this, simUI, player, slot));

		this.itemPickers = leftItemPickers.concat(rightItemPickers).sort((a, b) => a.slot - b.slot);

		this.selectorModal = new SelectorModal(simUI.rootElem, simUI, player, this);
	}
}

export class ItemRenderer extends Component {
	private readonly player: Player<any>;

	readonly iconElem: HTMLAnchorElement;
	readonly nameElem: HTMLAnchorElement;
	readonly ilvlElem: HTMLSpanElement;
	readonly enchantElem: HTMLAnchorElement;

	constructor(parent: HTMLElement, root: HTMLElement, player: Player<any>) {
		super(parent, 'item-picker-root', root);
		this.player = player;

		const iconElem = ref<HTMLAnchorElement>();
		const nameElem = ref<HTMLAnchorElement>();
		const ilvlElem = ref<HTMLSpanElement>();
		const enchantElem = ref<HTMLAnchorElement>();
		this.rootElem.appendChild(
			<>
				<div className="item-picker-icon-wrapper">
					<span className="item-picker-ilvl" ref={ilvlElem} />
					<a ref={iconElem} className="item-picker-icon" href="javascript:void(0)" attributes={{ role: 'button' }}></a>
				</div>
				<div className="item-picker-labels-container">
					<a ref={nameElem} className="item-picker-name" href="javascript:void(0)" attributes={{ role: 'button' }}></a>
					<a ref={enchantElem} className="item-picker-enchant hide" href="javascript:void(0)" attributes={{ role: 'button' }}></a>
				</div>
			</>,
		);

		this.iconElem = iconElem.value!;
		this.nameElem = nameElem.value!;
		this.ilvlElem = ilvlElem.value!;
		this.enchantElem = enchantElem.value!;
	}

	clear() {
		TooltipManager.detachTooltip(this.iconElem);
		this.iconElem.removeAttribute('href');
		TooltipManager.detachTooltip(this.nameElem);
		this.nameElem.removeAttribute('href');
		TooltipManager.detachTooltip(this.enchantElem);
		this.enchantElem.removeAttribute('href');
		this.enchantElem.classList.add('hide');

		this.iconElem.style.backgroundImage = '';

		this.nameElem.innerText = '';
		this.ilvlElem.innerText = '';
		this.enchantElem.innerText = '';
	}

	update(newItem: EquippedItem) {
		this.nameElem.textContent = newItem.item.name;
		this.ilvlElem.textContent = newItem.item.ilvl.toString();

		if (newItem.randomSuffix) {
			this.nameElem.textContent += ' ' + newItem.randomSuffix.name;
		}

		setItemQualityCssClass(this.nameElem, newItem.item.quality);

		this.player.attachTooltip(newItem, this.iconElem);
		this.player.attachTooltip(newItem, this.nameElem);
		newItem
			.asActionId()
			.fill()
			.then(filledId => {
				filledId.setBackground(this.iconElem);
				filledId.setDatabaseHref(this.iconElem);
				filledId.setDatabaseHref(this.nameElem);
			});

		if (newItem.enchant) {
			getEnchantDescription(newItem.enchant).then(description => {
				this.enchantElem.textContent = description;
			});
			// Make enchant text hover have a tooltip.
			if (newItem.enchant.spellId) {
				this.enchantElem.href = ActionId.makeSpellUrl(newItem.enchant.spellId);
				const actionId = ActionId.fromSpellId(newItem.enchant.spellId);
				TooltipManager.attachTooltip(this.enchantElem, actionId);
			} else {
				this.enchantElem.href = ActionId.makeItemUrl(newItem.enchant.itemId);
				const actionId = ActionId.fromItemId(newItem.enchant.itemId);
				TooltipManager.attachTooltip(this.enchantElem, actionId);
			}
			this.enchantElem.classList.remove('hide');
		} else {
			this.enchantElem.classList.add('hide');
		}
	}
}

export class ItemPicker extends Component {
	readonly slot: ItemSlot;

	private readonly simUI: SimUI;
	private readonly player: Player<any>;

	private readonly onUpdateCallbacks: (() => void)[] = [];

	private readonly itemElem: ItemRenderer;
	private readonly gearPicker: GearPicker;

	// All items and enchants that are eligible for this slot
	private _equippedItem: EquippedItem | null = null;

	constructor(parent: HTMLElement, gearPicker: GearPicker, simUI: SimUI, player: Player<any>, slot: ItemSlot) {
		super(parent, 'item-picker-root');

		this.gearPicker = gearPicker;
		this.simUI = simUI;
		this.player = player;
		this.slot = slot;
		this.itemElem = new ItemRenderer(parent, this.rootElem, player);

		this.item = player.getEquippedItem(slot);

		player.sim.waitForInit().then(() => {
			const openGearSelector = (event: Event) => {
				event.preventDefault();
				this.openSelectorModal(SelectorModalTabs.Items);
			};
			const openEnchantSelector = (event: Event) => {
				event.preventDefault();
				this.openSelectorModal(SelectorModalTabs.Enchants);
			};

			this.itemElem.iconElem.addEventListener('click', openGearSelector);
			this.itemElem.nameElem.addEventListener('click', openGearSelector);
			this.itemElem.enchantElem.addEventListener('click', openEnchantSelector);
		});

		player.gearChangeEmitter.on(() => {
			this.item = this.player.getEquippedItem(this.slot);
		});

		player.professionChangeEmitter.on(() => {
			if (!!this._equippedItem) {
				this.player.attachTooltip(this._equippedItem, this.itemElem.iconElem);
			}
		});
	}

	createGearData(): GearData {
		return {
			equipItem: (eventID: EventID, equippedItem: EquippedItem | null) => {
				this.player.equipItem(eventID, this.slot, equippedItem);
			},
			getEquippedItem: () => this.player.getEquippedItem(this.slot),
			changeEvent: this.player.gearChangeEmitter,
		};
	}

	get item(): EquippedItem | null {
		return this._equippedItem;
	}

	set item(newItem: EquippedItem | null) {
		// Clear everything first
		this.itemElem.clear();
		this.itemElem.nameElem.textContent = slotNames.get(this.slot) ?? '';
		setItemQualityCssClass(this.itemElem.nameElem, null);

		if (!!newItem) {
			this.itemElem.update(newItem);
		} else {
			this.itemElem.iconElem.style.backgroundImage = `url('${getEmptySlotIconUrl(this.slot)}')`;
		}

		this._equippedItem = newItem;
		this.onUpdateCallbacks.forEach(callback => callback());
	}

	onUpdate(callback: () => void) {
		this.onUpdateCallbacks.push(callback);
	}

	openSelectorModal(selectedTab: SelectorModalTabs) {
		this.gearPicker.selectorModal.openTab(this.slot, selectedTab, this.createGearData());
	}
}

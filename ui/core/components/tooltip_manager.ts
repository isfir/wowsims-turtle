import tippy, { Instance } from 'tippy.js';

import { UIItem, UISpell } from '../proto/ui';
import { ActionId } from '../proto_utils/action_id';
import { ItemTooltipContext, SpellTooltipContext } from '../tooltip_context';
import {
	getItemTooltipData,
	getSpellTooltipData,
	renderErrorTooltip,
	renderItemTooltip,
	renderLoadingTooltip,
	renderSimpleTooltip,
	renderSpellTooltip,
} from './tooltip_renderer';

const TOOLTIP_THEME = 'wow-tooltip';

type TooltipBinding =
	| {
		kind: 'action';
		actionId: ActionId;
		tooltipContext?: Partial<ItemTooltipContext> | Partial<SpellTooltipContext>;
	}
	| {
		kind: 'item';
		item: UIItem;
		tooltipContext?: Partial<ItemTooltipContext>;
	}
	| {
		kind: 'spell';
		spell: UISpell;
		tooltipContext?: Partial<SpellTooltipContext>;
	};

type TooltipListeners = {
	mouseenter: EventListener;
	focusin: EventListener;
	mouseleave: EventListener;
	focusout: EventListener;
};

const bindingMap = new WeakMap<HTMLElement, TooltipBinding>();
const listenerMap = new WeakMap<HTMLElement, TooltipListeners>();
const instanceMap = new WeakMap<HTMLElement, Instance>();
const requestVersionMap = new WeakMap<Instance, number>();

let activeTooltipElem: HTMLElement | null = null;
let globalGuardsInstalled = false;

function nextRequestVersion(instance: Instance): number {
	const nextVersion = (requestVersionMap.get(instance) ?? 0) + 1;
	requestVersionMap.set(instance, nextVersion);
	return nextVersion;
}

function isRequestStale(instance: Instance, requestVersion: number): boolean {
	const reference = instance.reference;
	return (
		instance.state.isDestroyed ||
		!(reference instanceof HTMLElement) ||
		!reference.isConnected ||
		requestVersionMap.get(instance) !== requestVersion
	);
}

function createBaseTooltipOptions(content: string) {
	return {
		theme: TOOLTIP_THEME,
		allowHTML: true,
		arrow: false,
		animation: false,
		duration: 0,
		hideOnClick: false,
		interactive: false,
		placement: 'auto-start' as const,
		delay: [150, 0] as [number, number],
		offset: [8, 6] as [number, number],
		maxWidth: 320,
		appendTo: () => document.body,
		trigger: 'manual' as const,
		content,
	};
}

function createTooltipInstance(elem: HTMLElement): Instance | null {
	if (!elem.isConnected) {
		return null;
	}

	const instance = tippy(elem, createBaseTooltipOptions(renderLoadingTooltip()));
	instanceMap.set(elem, instance);
	return instance;
}

function installGlobalGuards() {
	if (globalGuardsInstalled) {
		return;
	}

	globalGuardsInstalled = true;

	const hideActiveTooltip = () => {
		if (activeTooltipElem) {
			TooltipManager.destroyTooltipInstance(activeTooltipElem);
		}
	};

	window.addEventListener('scroll', hideActiveTooltip, true);
	window.addEventListener('resize', hideActiveTooltip, true);
	document.addEventListener('visibilitychange', () => {
		if (document.hidden) {
			hideActiveTooltip();
		}
	});
}

function bindTooltip(elem: HTMLElement, binding: TooltipBinding) {
	TooltipManager.detachTooltip(elem);
	installGlobalGuards();

	bindingMap.set(elem, binding);

	const show = () => TooltipManager.showTooltip(elem);
	const hide = () => TooltipManager.destroyTooltipInstance(elem);

	const listeners: TooltipListeners = {
		mouseenter: show,
		focusin: show,
		mouseleave: hide,
		focusout: hide,
	};

	elem.addEventListener('mouseenter', listeners.mouseenter);
	elem.addEventListener('focusin', listeners.focusin);
	elem.addEventListener('mouseleave', listeners.mouseleave);
	elem.addEventListener('focusout', listeners.focusout);

	listenerMap.set(elem, listeners);
}

export class TooltipManager {
	static attachTooltip(
		elem: HTMLElement,
		actionId: ActionId,
		tooltipContext?: Partial<ItemTooltipContext> | Partial<SpellTooltipContext>,
	) {
		bindTooltip(elem, {
			kind: 'action',
			actionId,
			tooltipContext: tooltipContext,
		});
	}

	static attachItemTooltip(
		elem: HTMLElement,
		item: UIItem,
		tooltipContext?: Partial<ItemTooltipContext>,
	) {
		bindTooltip(elem, {
			kind: 'item',
			item,
			tooltipContext: tooltipContext,
		});
	}

	static attachSpellTooltip(
		elem: HTMLElement,
		spell: UISpell,
		tooltipContext?: Partial<SpellTooltipContext>,
	) {
		bindTooltip(elem, {
			kind: 'spell',
			spell,
			tooltipContext: tooltipContext,
		});
	}

	static showTooltip(elem: HTMLElement) {
		if (!elem.isConnected) {
			return;
		}

		const binding = bindingMap.get(elem);
		if (!binding) {
			return;
		}

		if (activeTooltipElem && activeTooltipElem !== elem) {
			this.destroyTooltipInstance(activeTooltipElem);
		}

		// Always recreate fresh so we never keep old orphaned poppers around.
		this.destroyTooltipInstance(elem);

		const instance = createTooltipInstance(elem);
		if (!instance) {
			return;
		}

		activeTooltipElem = elem;

		switch (binding.kind) {
			case 'action':
				instance.setContent(renderLoadingTooltip());
				instance.show();
				void updateActionTooltipContent(instance, binding.actionId, binding.tooltipContext);
				break;

			case 'item':
				instance.setContent(renderItemTooltip(binding.item, binding.tooltipContext));
				instance.show();

				if (binding.item.id > 0) {
					void updateItemTooltipContent(instance, binding.item.id, binding.item, binding.tooltipContext);
				}
				break;

			case 'spell':
				instance.setContent(renderSpellTooltip(binding.spell, binding.tooltipContext));
				instance.show();
				break;
		}
	}

	static destroyTooltipInstance(elem: HTMLElement) {
		const instance = instanceMap.get(elem);
		if (!instance) {
			return;
		}

		instance.destroy();
		instanceMap.delete(elem);
		requestVersionMap.delete(instance);

		if (activeTooltipElem === elem) {
			activeTooltipElem = null;
		}
	}

	static detachTooltip(elem: HTMLElement) {
		this.destroyTooltipInstance(elem);

		const listeners = listenerMap.get(elem);
		if (listeners) {
			elem.removeEventListener('mouseenter', listeners.mouseenter);
			elem.removeEventListener('focusin', listeners.focusin);
			elem.removeEventListener('mouseleave', listeners.mouseleave);
			elem.removeEventListener('focusout', listeners.focusout);
			listenerMap.delete(elem);
		}

		bindingMap.delete(elem);
	}
}

async function updateActionTooltipContent(
	instance: Instance,
	actionId: ActionId,
	tooltipContext?: Partial<ItemTooltipContext> | Partial<SpellTooltipContext>,
) {
	const requestVersion = nextRequestVersion(instance);
	instance.setContent(renderLoadingTooltip());

	try {
		const filledId = await actionId.fill();

		let content = renderSimpleTooltip(filledId.name || 'Unknown');

		if (filledId.itemId > 0) {
			const item = await getItemTooltipData(filledId.itemId, tooltipContext as Partial<ItemTooltipContext>);
			content = item
				? renderItemTooltip(item, tooltipContext as Partial<ItemTooltipContext>)
				: renderSimpleTooltip(filledId.name || `Item ${filledId.itemId}`, 'No item data available.');
		} else if (filledId.spellIdTooltipOverride || filledId.spellId > 0) {
			const spell = await getSpellTooltipData(filledId.spellIdTooltipOverride || filledId.spellId);
			content = spell
				? renderSpellTooltip(spell, tooltipContext as Partial<SpellTooltipContext>)
				: renderSimpleTooltip(filledId.name || `Spell ${filledId.spellId}`, 'No spell data available.');
		}

		if (isRequestStale(instance, requestVersion)) {
			return;
		}

		instance.setContent(content);
	} catch (error) {
		console.error('Failed to load tooltip content:', error);

		if (isRequestStale(instance, requestVersion)) {
			return;
		}

		instance.setContent(renderErrorTooltip('Failed to load tooltip.'));
	}
}

async function updateItemTooltipContent(
	instance: Instance,
	itemId: number,
	fallbackItem: UIItem,
	tooltipContext?: Partial<ItemTooltipContext>,
) {
	const requestVersion = nextRequestVersion(instance);

	try {
		const item = await getItemTooltipData(itemId, tooltipContext);
		if (!item) {
			return;
		}

		if (isRequestStale(instance, requestVersion)) {
			return;
		}

		instance.setContent(renderItemTooltip(item, tooltipContext));
	} catch (error) {
		console.error(`Failed to enrich item tooltip for item ${itemId}:`, error);

		if (isRequestStale(instance, requestVersion)) {
			return;
		}

		instance.setContent(renderItemTooltip(fallbackItem, tooltipContext));
	}
}

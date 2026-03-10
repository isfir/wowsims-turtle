export type ItemTooltipContext = {
    /**Item id for the tooltip being resolved.*/
    itemId: number;

    /**Optional override for the displayed item level.*/
    itemLevel?: number;

    /**Viewer/player level, used for requirement coloring or level-dependent display.*/
    level?: number;

    /**Applied enchant effect id to render alongside the item.*/
    enchantId?: number;

    /**Socketed gem item ids. Reserved for future local tooltip rendering.*/
    gemIds?: number[];

    /**Whether the item has an extra socket. Reserved for future local tooltip rendering.*/
    hasExtraSocket?: boolean;

    /**Equipped set piece item ids, used to highlight active set bonuses.*/
    setPieceIds?: number[];

    /**Random enchant/suffix id. Reserved for future local tooltip rendering.*/
    randomEnchantmentId?: number;

    /**Reforge id. Reserved for future local tooltip rendering.*/
    reforgeId?: number;

    /**Transmog source item id. Reserved for future local tooltip rendering.*/
    transmogId?: number;
};

export type SpellTooltipContext = {
    /**Spell id for the tooltip being resolved.*/
    spellId: number;

    /**Viewer/player level, used for level-dependent display if needed.*/
    level?: number;

    /**Optional difficulty context.*/
    difficultyId?: 14 | 15 | 16;
};

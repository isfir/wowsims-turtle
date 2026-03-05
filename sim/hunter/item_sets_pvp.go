package hunter

import (
	"github.com/isfir/wowsims-turtle/sim/core"
	"github.com/isfir/wowsims-turtle/sim/core/stats"
)

///////////////////////////////////////////////////////////////////////////
//                            Classic Phase 2
///////////////////////////////////////////////////////////////////////////

// https://www.wowhead.com/classic/item-set=362/lieutenant-commanders-pursuit
// var ItemSetLieutenantCommandersPursuit = core.NewItemSet(core.ItemSet{
// 	Name: "Lieutenant Commander's Pursuit",
// 	Bonuses: map[int32]core.ApplyEffect{
// 		// Increases your chance to parry an attack by 1%.
// 		2: func(agent core.Agent) {
// 			c := agent.GetCharacter()
// 			c.AddStat(stats.Parry, 1*core.ParryRatingPerParryChance)
// 		},
// 		// Reduces the cooldown of your Concussive Shot by 1 sec.
// 		4: func(agent core.Agent) {
// 			// Nothing to do
// 		},
// 		// +15 Stamina.
// 		6: func(agent core.Agent) {
// 			c := agent.GetCharacter()
// 			c.AddStat(stats.Stamina, 15)
// 		},
// 	},
// })

// https://www.wowhead.com/classic/item-set=361/champions-pursuit
// var ItemSetChampionsPursuit = core.NewItemSet(core.ItemSet{
// 	Name: "Champion's Pursuit",
// 	Bonuses: map[int32]core.ApplyEffect{
// 		// Increases your chance to parry an attack by 1%.
// 		2: func(agent core.Agent) {
// 			c := agent.GetCharacter()
// 			c.AddStat(stats.Parry, 1*core.ParryRatingPerParryChance)
// 		},
// 		// Reduces the cooldown of your Concussive Shot by 1 sec.
// 		4: func(agent core.Agent) {
// 			// Nothing to do
// 		},
// 		// +15 Stamina.
// 		6: func(agent core.Agent) {
// 			c := agent.GetCharacter()
// 			c.AddStat(stats.Stamina, 15)
// 		},
// 	},
// })

// https://www.wowhead.com/classic/item-set=543/champions-pursuance
var ItemSetChampionsPursuance = core.NewItemSet(core.ItemSet{
	Name: "Champion's Pursuance",
	Bonuses: map[int32]core.ApplyEffect{
		// +20 Agility.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Agility, 20)
		},
		// Reduces the cooldown of your Concussive Shot by 1 sec.
		4: func(agent core.Agent) {
			// Nothing to do
		},
		// +20 Stamina.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 20)
		},
	},
})

// https://www.wowhead.com/classic/item-set=550/lieutenant-commanders-pursuance
var ItemSetLieutenantCommandersPursuance = core.NewItemSet(core.ItemSet{
	Name: "Lieutenant Commander's Pursuance",
	Bonuses: map[int32]core.ApplyEffect{
		// +20 Agility.
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Agility, 20)
		},
		// Reduces the cooldown of your Concussive Shot by 1 sec.
		4: func(agent core.Agent) {
			// Nothing to do
		},
		// +20 Stamina.
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 20)
		},
	},
})

///////////////////////////////////////////////////////////////////////////
//                            Classic Phase 3
///////////////////////////////////////////////////////////////////////////

var ItemSetWarlordsPursuit = core.NewItemSet(core.ItemSet{
	Name: "Warlord's Pursuit",
	Bonuses: map[int32]core.ApplyEffect{
		// 20 Stamina
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 20)
		},
		// Reduces the cooldown of your Concussive Shot by 1 sec.
		4: func(agent core.Agent) {
			// Nothing to do
		},
		// +20 Agi
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Agility, 20)
		},
	},
})

var ItemSetFieldMarshalsPursuit = core.NewItemSet(core.ItemSet{
	Name: "Field Marshal's Pursuit",
	Bonuses: map[int32]core.ApplyEffect{
		// 20 stamina
		2: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Stamina, 20)
		},
		// Reduces the cooldown of your Concussive Shot by 1 sec.
		4: func(agent core.Agent) {
			// Nothing to do
		},
		// +20 Agi
		6: func(agent core.Agent) {
			c := agent.GetCharacter()
			c.AddStat(stats.Agility, 20)
		},
	},
})

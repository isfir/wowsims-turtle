package encounters

import (
	"github.com/isfir/wowsims-turtle/sim/core"
)

func init() {
	// TODO: Classic encounters?
	// naxxramas.Register()
	addDefaultRaidBoss()
	addDefaultTrash()
	addVaelastraszTheCorrupt("Classic")
}

func AddSingleTargetBossEncounter(presetTarget *core.PresetTarget) {
	core.AddPresetTarget(presetTarget)
	core.AddPresetEncounter(presetTarget.Config.Name, []string{
		presetTarget.Path(),
	})
}

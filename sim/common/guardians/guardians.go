package guardians

import "github.com/isfir/wowsims-turtle/sim/core"

func ConstructGuardians(character *core.Character) {
	constructEmeralDragonWhelps(character)
	constructEskhandar(character)
	constructCoreHound(character)
	constructMinorArcaneElemental(character)
	constructScytheclawPureborn(character)
}

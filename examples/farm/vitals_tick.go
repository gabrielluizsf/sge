package farm

import "github.com/gabrielluizsf/sge/vitals"

// restSeconds is how long a night's sleep feeds on the player, and
// restHeal how much health it gives back.
const (
	restSeconds = 270
	restHeal    = 40
)

// vitalsTick passes dt seconds for the player's health and food, and
// lets the red flash of a blow die down. Health run out, they die.
func (f *farm) vitalsTick(dt float64) {
	f.flash = max(0, f.flash-dt)
	f.bloodObj.Hidden, f.stainObj.Hidden = f.room != nil, f.room != nil
	if f.question != nil {
		return
	}
	f.vit.Tick(dt)
	if f.vit.Dead() {
		f.fall()
	}
}

// rest is the player sleeping the night: hungrier for it, but mended.
func (f *farm) rest() {
	f.vit.Tick(restSeconds)
	f.vit.Heal(restHeal)
}

// fall is the player dying of a blow or of hunger: their heir takes over,
// healthy and fed, or the line ends.
func (f *farm) fall() {
	f.throwOff()
	f.passOn(f.clock.Day, false)
	f.vit = vitals.Full()
	f.flash = 0
}

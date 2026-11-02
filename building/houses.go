package building

import "image/color"

// The houses' paints, in sRGB.
var (
	colonialBlue = color.NRGBA{70, 104, 150, 255}
	bareWood     = color.NRGBA{}
	white        = color.NRGBA{236, 234, 226, 255}
	bottleGreen  = color.NRGBA{58, 96, 66, 255}
)

// farmhouseStyle is a farm's house: one long storey, whitewashed, under a
// hipped roof of clay tiles, a porch along its front, blue shutters.
var farmhouseStyle = House{Width: 8, Depth: 6, Storeys: 1, Storey: 3, Wall: Plaster, Tint: color.NRGBA{250, 246, 236, 255},
	Roof: Roof{Cover: Tiles, Hipped: true, Pitch: 26, Overhang: 0.6}, Chimney: true, Porch: true,
	Trim: colonialBlue, Door: colonialBlue, Shutters: &colonialBlue}

// Farmhouse is the house of farmhouseStyle.
var Farmhouse = farmhouseStyle.Building("farmhouse")

// cottageStyle is a small house of stone under a steep hipped thatch.
var cottageStyle = House{Width: 6.5, Depth: 5, Storeys: 1, Storey: 2.7, Wall: Stone,
	Roof: Roof{Cover: Thatch, Hipped: true, Pitch: 48, Overhang: 0.45}, Chimney: true, Shutters: &bareWood}

// Cottage is the house of cottageStyle.
var Cottage = cottageStyle.Building("cottage")

// cabinStyle is a farmhand's house of boards under a gabled tile roof.
var cabinStyle = House{Width: 6.5, Depth: 5, Storeys: 1, Storey: 2.7, Wall: Planks,
	Roof: Roof{Cover: Tiles, Pitch: 30, Overhang: 0.45}, Chimney: true}

// Cabin is the house of cabinStyle.
var Cabin = cabinStyle.Building("cabin")

// townHouseStyle is a half-timbered house of two storeys, its plaster cream
// between dark timbers, under a steep gable of tiles.
var townHouseStyle = House{Width: 6.5, Depth: 5.5, Storeys: 2, Wall: Plaster, Tint: color.NRGBA{240, 226, 190, 255}, Timbered: true,
	Roof: Roof{Cover: Tiles, Pitch: 42, Overhang: 0.35}, Chimney: true, Door: color.NRGBA{150, 70, 45, 255}}

// TownHouse is the house of townHouseStyle.
var TownHouse = townHouseStyle.Building("townhouse")

// brickHouseStyle is a town house of brick, two storeys under slates, its
// frames painted white and its door green.
var brickHouseStyle = House{Width: 7, Depth: 6, Storeys: 2, Wall: Brick,
	Roof: Roof{Cover: Slates, Pitch: 38, Overhang: 0.3}, Chimney: true, Trim: white, Door: bottleGreen}

// BrickHouse is the house of brickHouseStyle.
var BrickHouse = brickHouseStyle.Building("brickhouse")

// sobradoStyle is a colonial town house of two storeys, its plaster ochre,
// under a low hipped roof of tiles, its frames white, its shutters green.
var sobradoStyle = House{Width: 8, Depth: 6.5, Storeys: 2, Wall: Plaster, Tint: color.NRGBA{236, 196, 120, 255},
	Roof: Roof{Cover: Tiles, Hipped: true, Pitch: 24, Overhang: 0.55}, Trim: white, Door: bottleGreen, Shutters: &bottleGreen}

// Sobrado is the house of sobradoStyle.
var Sobrado = sobradoStyle.Building("sobrado")

// Houses are every style of house.
var Houses = []Building{Farmhouse, Cottage, Cabin, TownHouse, BrickHouse, Sobrado}

// houseStyles are the styles of Houses, in their order.
var houseStyles = []House{farmhouseStyle, cottageStyle, cabinStyle, townHouseStyle, brickHouseStyle, sobradoStyle}

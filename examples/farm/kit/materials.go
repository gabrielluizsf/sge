package kit

import "github.com/gabrielluizsf/sge/building"

// Clad is a texture the farm's buildings are clad in: the material of
// package building's it is, and whether its boards run across its image.
type Clad struct {
	Material building.Material
	Across   bool
}

// Materials are the textures, by their assets' ids, the farm's buildings
// are clad in.
var Materials = map[string]Clad{
	"painted_plaster_wall":       {Material: building.Plaster},
	"brick_wall_006":             {Material: building.Brick},
	"stone_wall":                 {Material: building.Stone},
	"brown_planks_07":            {Material: building.Planks, Across: true},
	"rough_wood":                 {Material: building.Timber},
	"clay_roof_tiles_02":         {Material: building.Tiles},
	"reed_roof_04":               {Material: building.Thatch},
	"roof_slates_02":             {Material: building.Slates},
	"brown_planks_09":            {Material: building.Floorboards, Across: true},
	"pine_bark":                  {Material: building.Bark},
	"hessian_380":                {Material: building.Hessian},
	"fabric_pattern_07":          {Material: building.Check},
	"curly_teddy_checkered":      {Material: building.Plaid},
	"quatrefoil_jacquard_fabric": {Material: building.Jacquard},
	"rough_linen":                {Material: building.Linen},
}

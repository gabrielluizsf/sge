package farm

// snapshot is the game's state that is needed to resume it, as plain data:
// no pointers, things refer to one another by index in a list or by id. It
// holds the clock, the player (kit, hotbar, place, who they are), the
// family tree, the plots, the trees being cut and the stumps, the smith's
// errand, the towns woken and their people. It does not hold what is
// rebuilt from the world's seed (the map, the trees standing, the houses),
// the menus, a conversation or a room being open, the villagers' errands
// of the moment and the notices on screen.
type snapshot struct {
	Day     int     `json:"day"`
	Hour    float64 `json:"hour"`
	Seen    int     `json:"seen"`
	Fainted int     `json:"fainted"`
	Ids     int     `json:"ids"`
	Luck    uint32  `json:"luck"`
	Fate    uint32  `json:"fate"`
	// Level and XP are the family's; an old save has neither, and is level 1.
	Level int `json:"level"`
	XP    int `json:"xp"`
	// Vitals are the player's health and food; an old save has none, and
	// its player is at their best.
	Vitals *vitalsState `json:"vitals"`

	Player kitState   `json:"player"`
	Hotbar []string   `json:"hotbar"`
	Active int        `json:"active"`
	At     point      `json:"at"`
	Self   selfState  `json:"self"`
	Family []kinState `json:"family"`

	Plots      []plotState  `json:"plots"`
	Stumps     []stumpState `json:"stumps"`
	Hits       []hitState   `json:"hits"`
	Errand     orderState   `json:"errand"`
	ErrandDone bool         `json:"errandDone"`

	// Home is the furniture of the player's house and Yard the chests out of
	// doors; an old save has neither, and its house has only the bed.
	Home []pieceState `json:"home"`
	Yard []pieceState `json:"yard"`

	// Herd is the farm's animals; an old save has none, and its pen has
	// the first cow.
	Herd []creatureState `json:"herd,omitempty"`
	// Pens are the farm's, each animal saying which by its index.
	Pens []penState `json:"pens,omitempty"`

	Towns []cellState `json:"towns"`
	// Away is when each town that was asleep was last seen, to settle what
	// its people did meanwhile; an old save has none, and its towns
	// have nothing to settle.
	Away   []awayState     `json:"away,omitempty"`
	People []villagerState `json:"people"`
}

// awayState is when the town of a region was last seen, in hours since the
// first day began.
type awayState struct {
	Col int     `json:"col"`
	Row int     `json:"row"`
	At  float64 `json:"at"`
}

// creatureState is an animal of the farm: its kind and coat, where it is,
// and one more than the day it last gave milk, 0 for never.
type creatureState struct {
	Kind   string `json:"kind"`
	Coat   int    `json:"coat"`
	At     point  `json:"at"`
	Milked int    `json:"milked"`
	Pen    int    `json:"pen"`
}

// penState is a pen of the farm: the kind it keeps and its corner cell.
type penState struct {
	Kind string `json:"kind"`
	Col  int    `json:"col"`
	Row  int    `json:"row"`
}

// point is a place in the world.
type point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// cellState is a cell of the world, or a town's region.
type cellState struct {
	Col int `json:"col"`
	Row int `json:"row"`
}

// stackState is n of an item.
type stackState struct {
	Item string `json:"item"`
	N    int    `json:"n"`
}

// kitState is a kit: money, bag, the slot and seeds in hand, the axe,
// what was sold and how much of the can is poured.
type kitState struct {
	Money int          `json:"money"`
	Bag   []stackState `json:"bag"`
	Slot  int          `json:"slot"`
	Seed  int          `json:"seed"`
	Axe   int          `json:"axe"`
	Sold  int          `json:"sold"`
	Used  int          `json:"used"`
	Exact bool         `json:"exact"`
}

// orderState is what someone asks for and pays.
type orderState struct {
	Want   []stackState `json:"want"`
	Reward int          `json:"reward"`
}

// lifeState is a life: the day born and the span in years.
type lifeState struct {
	Born int     `json:"born"`
	Span float64 `json:"span"`
}

// expectState is a child being expected: the parents' ids and the day due.
type expectState struct {
	Parents [2]string `json:"parents"`
	Due     int       `json:"due"`
}

// selfState is who the player is. Spouse is an index of People, or -1.
type selfState struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Woman     bool         `json:"woman"`
	Genes     [5]int       `json:"genes"`
	Life      lifeState    `json:"life"`
	Spouse    int          `json:"spouse"`
	Expecting *expectState `json:"expecting"`
}

// kinState is a person of the family tree.
type kinState struct {
	ID      string    `json:"id"`
	Parents [2]string `json:"parents"`
	Born    int       `json:"born"`
	Dead    bool      `json:"dead"`
}

// plotState is a plot of soil. Town is the region of the town whose
// garden it is, nil for the player's field.
type plotState struct {
	Col     int        `json:"col"`
	Row     int        `json:"row"`
	Town    *cellState `json:"town"`
	Crop    int        `json:"crop"`
	Stage   int        `json:"stage"`
	Days    int        `json:"days"`
	Watered bool       `json:"watered"`
	Dead    bool       `json:"dead"`
	// Owner is 1 plus the index in People of the villager whose crop it is
	// on the farm, 0 for the player's (and for a save with no owners).
	Owner int `json:"owner"`
}

// stumpState is a tree felled, growing back.
type stumpState struct {
	Col  int     `json:"col"`
	Row  int     `json:"row"`
	Tree int     `json:"tree"`
	Left float64 `json:"left"`
}

// hitState is the strokes a tree being cut has taken.
type hitState struct {
	Col int `json:"col"`
	Row int `json:"row"`
	N   int `json:"n"`
}

// villagerState is a villager. Town is the region of their town, nil for
// none; Home an index of its houses or -1; Partner an index of People or -1.
type villagerState struct {
	Name      string       `json:"name"`
	Woman     bool         `json:"woman"`
	Role      string       `json:"role"`
	Genes     [5]int       `json:"genes"`
	Kit       kitState     `json:"kit"`
	Town      *cellState   `json:"town"`
	Home      int          `json:"home"`
	Order     orderState   `json:"order"`
	Wait      float64      `json:"wait"`
	At        point        `json:"at"`
	Inside    bool         `json:"inside"`
	HomeFor   float64      `json:"homeFor"`
	Fuel      float64      `json:"fuel"`
	ID        string       `json:"id"`
	Life      lifeState    `json:"life"`
	Bond      bondState    `json:"bond"`
	Partner   int          `json:"partner"`
	Expecting *expectState `json:"expecting"`
	Ours      bool         `json:"ours"`
	Luck      uint32       `json:"luck"`
}

// bondState is the player's bond with a villager.
type bondState struct {
	Affection int `json:"affection"`
	Courted   int `json:"courted"`
	Gifted    int `json:"gifted"`
	Talked    int `json:"talked"`
	Wedding   int `json:"wedding"`
}

// vitalsState is the player's health and food, each from 0 to 100.
type vitalsState struct {
	Health float64 `json:"health"`
	Food   float64 `json:"food"`
}

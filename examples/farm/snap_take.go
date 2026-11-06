package farm

// takeSnapshot is f's state as plain data.
func takeSnapshot(f *farm) snapshot {
	s := snapshot{
		Day: f.clock.Day, Hour: f.clock.Hour, Seen: f.seen, Fainted: f.fainted, Ids: f.ids,
		Luck: uint32(*f.luck), Fate: uint32(*f.fate),
		Player: takeKit(f.kit), Hotbar: append([]string(nil), f.pack.Hotbar...), Active: f.pack.Active,
		At: takePlace(f), Self: takeSelf(f), Family: takeFamily(f),
		Errand: takeOrder(f.errand), ErrandDone: f.errandDone,
		Towns: takeTowns(f), People: takePeople(f),
	}
	s.Level, s.XP = f.prog.Level, f.prog.XP
	s.Vitals = &vitalsState{Health: f.vit.Health, Food: f.vit.Food}
	s.Home, s.Yard = takeHouse(f), takeYard(f)
	s.Plots, s.Stumps, s.Hits = takePlots(f), takeStumps(f), takeHits(f)
	return s
}

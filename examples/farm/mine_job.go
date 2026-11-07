package farm

// mineJob is a job breaking one of the rocks about the villager's town,
// with a pickaxe they take up if they have none; never the player's own
// land. The rock grows back as the player's do.
func (m *mind) mineJob(f *farm) *job {
	v := m.v
	if v.kit.bag.Count(itemPick) == 0 {
		v.kit.bag.Add(itemPick, 1)
	}
	return m.strike(f, f.mineQuarry(v))
}

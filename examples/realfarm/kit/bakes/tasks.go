package bakes

// Tasks are everything the game is made of: the grounds, the shores of
// water, the crops, the buildings, the models of trees and things, and
// the villagers, grown-up first and the other stages after. Those with
// Core set come before the rest.
func (b *Baker) Tasks() []Task {
	var all []Task
	all = append(all, b.groundTasks()...)
	all = append(all, b.waterTasks()...)
	all = append(all, b.cropTasks()...)
	all = append(all, b.iconTasks()...)
	all = append(all, b.portraitTasks()...)
	all = append(all, b.buildingTasks()...)
	all = append(all, b.modelTasks()...)
	all = append(all, b.treeTasks()...)
	all = append(all, b.importTasks()...)
	all = append(all, b.peopleTasks()...)
	var core, rest []Task
	for _, t := range all {
		if t.Core {
			core = append(core, t)
		} else {
			rest = append(rest, t)
		}
	}
	return append(core, rest...)
}

package farm

import "fmt"

// checkLife is nil if l is a life the game could draw.
func checkLife(what string, l lifeState) error {
	return firstError(inRange(what+" born", l.Born, -maxDay, maxDay), finite(what+" span", l.Span, 0, maxYears))
}

// checkGenes is nil if g are five numbers a person is drawn from.
func checkGenes(what string, g [5]int) error {
	for _, n := range g {
		if err := inRange(what+" genes", n, 0, maxGene); err != nil {
			return err
		}
	}
	return nil
}

// checkExpecting is nil if e is nil or a child expected.
func checkExpecting(what string, e *expectState) error {
	if e == nil {
		return nil
	}
	for _, id := range e.Parents {
		if id != "" {
			if err := plain(what+" parent", id); err != nil {
				return err
			}
		}
	}
	return inRange(what+" due", e.Due, 0, 2*maxDay)
}

// checkSelf is nil if s is who a player could be.
func checkSelf(s *snapshot) error {
	p := s.Self
	return firstError(
		plain("player id", p.ID), plain("player name", p.Name),
		checkGenes("player", p.Genes), checkLife("player", p.Life),
		inRange("spouse", p.Spouse, -1, len(s.People)-1),
		checkExpecting("player's child", p.Expecting),
		checkSpouse(s),
	)
}

// checkSpouse is nil if the player's spouse, if any, lives in their house.
func checkSpouse(s *snapshot) error {
	if i := s.Self.Spouse; i >= 0 && i < len(s.People) && !s.People[i].Ours {
		return fmt.Errorf("save: the spouse does not live in the player's house")
	}
	return nil
}

// checkClock is nil if the time is one the game could be at.
func checkClock(s *snapshot) error {
	return firstError(
		inRange("day", s.Day, 0, maxDay), finite("hour", s.Hour, 0, 23.999999),
		inRange("seen day", s.Seen, 0, s.Day), inRange("fainted day", s.Fainted, -1, s.Day),
		inRange("ids", s.Ids, 0, maxFamily*4),
	)
}

package i18n

// Add gives the Portuguese words of a game, by their English keys (a key
// ending in "|f" or "|m" being a woman's or a man's word), to those T
// and TF translate with; a key given again takes its new word. A game
// adds its words as it starts, in an init, before T is first called.
func Add(words map[string]string) {
	for k, v := range words {
		translations[k] = v
	}
}

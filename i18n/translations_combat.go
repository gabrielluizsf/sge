package i18n

// combatWords are the Portuguese words of health, food and fighting,
// joined to translations.
var combatWords = map[string]string{
	"F: attack (the axe hits harder)":  "F: atacar (o machado fere mais)",
	"HP %d":                            "Vida %d",
	"Food %d":                          "Fome %d",
	"+%d food":                         "+%d comida",
	"Full":                             "Satisfeito",
	"%s was killed.":                   "%s foi morto.",
	"%s died. %s takes over the farm.": "%s morreu. %s assume a fazenda.",
	"Eat: +%d food":                    "Comer: +%d comida",
	"Eats the crop in hand":            "Come a colheita na mão",
}

func init() {
	for k, v := range combatWords {
		translations[k] = v
	}
}

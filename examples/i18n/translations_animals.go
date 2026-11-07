package i18n

// animalWords are the Portuguese words of the animals and what they give,
// joined to translations.
var animalWords = map[string]string{
	"Milk bottle":                         "Garrafa de leite",
	"Wool":                                "Lã",
	"From your cow, once a day":           "Da sua vaca, uma vez por dia",
	"Beds are made of it":                 "Com ela se fazem camas",
	"E: milk the cow":                     "E: tirar leite da vaca",
	"Already milked today":                "Já tirou leite hoje",
	"Sell goods $%d":                      "Vender produtos $%d",
	"E: buy animals":                      "E: comprar animais",
	"Stable":                              "Estábulo",
	"%s $%d":                              "%s $%d",
	"%s was killed by a wolf.":            "%s foi morto por um lobo.",
	"Wolf killed":                         "Lobo morto",
	"Bread":                               "Pão",
	"Beef":                                "Carne de boi",
	"Mutton":                              "Carne de carneiro",
	"Cow":                                 "Vaca",
	"Bakery":                              "Padaria",
	"Butcher":                             "Açougue",
	"baker|f":                             "padeira",
	"baker|m":                             "padeiro",
	"butcher|f":                           "açougueira",
	"butcher|m":                           "açougueiro",
	"E: enter the bakery":                 "E: entrar na padaria",
	"E: buy meat":                         "E: comprar carne",
	"E: bake bread":                       "E: assar pão",
	"Baked from wheat at a forge":         "Assado de trigo numa forja",
	"From the butcher, or your own cow":   "Do açougue ou da sua vaca",
	"From the butcher, or your own sheep": "Do açougue ou da sua ovelha",
	"Buy %s $%d":                          "Comprar %s $%d",
	"Sell wheat $%d each (you have %d)":   "Vender trigo a $%d cada (você tem %d)",
	"Sell your %s $%d":                    "Vender seu animal (%s) $%d",
	"You have none":                       "Você não tem nenhum",
	"+%d %s":                              "+%d %s",
	"Drink: +%d food":                     "Beber: +%d comida",
	"Sheep":                               "Ovelha",
	"Horse":                               "Cavalo",
	"Pen":                                 "Cercado",
	"Pen moved":                           "Cercado movido",
	"Enter: move this pen":                "Enter: mover este cercado",
	"livestock dealer":                    "pecuarista",
	"E: shear the sheep":                  "E: tosquiar a ovelha",
	"Wool grows back in %d days":          "A lã volta a crescer em %d dias",
	"There is no room for another pen":    "Não há lugar para outro cercado",
	"Your %s is in a pen by the farm: press B there to move it": "Seu animal (%s) está num cercado perto da fazenda: aperte B lá para movê-lo",
}

func init() {
	for k, v := range animalWords {
		translations[k] = v
	}
}

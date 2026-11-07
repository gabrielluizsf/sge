package i18n

// houseWords are the Portuguese words of the craft tab, the arranging of
// the house, the chests and the chairs, joined to translations.
var houseWords = map[string]string{
	"Craft":                       "Criar",
	"Made: %s":                    "Criado: %s",
	"Chest":                       "Baú",
	"Stone":                       "Pedra",
	"Pickaxe":                     "Picareta",
	"Breaks rocks into stone":     "Quebra pedras em pedaços",
	"%d hits per rock":            "%d golpes por rocha",
	"+%d stone":                   "+%d pedra",
	"E: mine":                     "E: minerar",
	"Bed":                         "Cama",
	"Table":                       "Mesa",
	"Chair":                       "Cadeira",
	"Shelf":                       "Estante",
	"Stove":                       "Fogão",
	"Plant pot":                   "Vaso de planta",
	"Furniture is made of it":     "Móveis são feitos com ela",
	"Sleep in it to end the day":  "Durma nela para terminar o dia",
	"A place to set things down":  "Um lugar para apoiar coisas",
	"Sit on it with E":            "Sente-se nela com E",
	"Tall and tidy":               "Alta e arrumada",
	"Warms the house":             "Aquece a casa",
	"Greens up a corner":          "Deixa um canto verde",
	"Keeps 10 stacks of items":    "Guarda 10 pilhas de itens",
	"B: arrange the house":        "B: arrumar a casa",
	"Pick it back up to carry it": "Pegue de volta para carregar",
	"It does not fit there":       "Não cabe aí",
	"Nothing to place: craft furniture in the backpack": "Nada para colocar: crie móveis na mochila",
	"Empty the chest first":                             "Esvazie o baú primeiro",
	"You need a bed to sleep in":                        "Você precisa de uma cama para dormir",
	"E: open chest":                                     "E: abrir o baú",
	"E: sit":                                            "E: sentar",
	"Move to stand up":                                  "Ande para levantar",
	"Arranging: nothing chosen":                         "Arrumando: nada escolhido",
	"Holding: %s":                                       "Na mão: %s",
	"Place: %s (Q/E: change)":                           "Colocar: %s (Q/E: trocar)",
	"Arrows or mouse: move   Enter or click: lift/place   R or right click: turn   X: store   Esc/B: done": "Setas ou mouse: mover   Enter ou clique: pegar/colocar   R ou botão direito: girar   X: guardar   Esc/B: pronto",
	"Enter: craft   Q/E: tabs   Esc: close":                         "Enter: criar   Q/E: abas   Esc: fechar",
	"Enter: pick up   Shift+Enter: send across   Esc: close":        "Enter: pegar   Shift+Enter: passar para o outro lado   Esc: fechar",
	"Enter: pick up   1-6: shortcut   Q/E: tabs   I/Tab/Esc: close": "Enter: pegar   1-6: atalho   Q/E: abas   I/Tab/Esc: fechar",
	"E: interact (talk, shop, doors, sit, chests)":                  "E: interagir (conversar, loja, portas, sentar, baús)",
	"B: arrange furniture (in your house; chests outside)":          "B: arrumar móveis (na sua casa; baús lá fora)",
	"Backpack: Q/E change tab, Craft makes furniture":               "Mochila: Q/E troca de aba, Criar faz móveis",
}

func init() {
	for k, v := range houseWords {
		translations[k] = v
	}
}

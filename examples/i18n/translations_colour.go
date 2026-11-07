package i18n

// colourWords are the Portuguese words of the colour code field and the
// eyedropper in the character creator, joined to translations.
var colourWords = map[string]string{
	"Item code":   "Código da peça",
	"Bottom code": "Código da parte de baixo",
	"Custom":      "Personalizada",
	"Pick":        "Pegar",
	"Clothes cannot be see-through: transparency is ignored": "Roupa não é transparente: a transparência é ignorada",
	"Not a colour code. Try #RRGGBB or rgba(r, g, b, a)":     "Código inválido. Use #RRGGBB ou rgba(r, g, b, a)",
	"The eyedropper is not available on this system":         "O conta-gotas não está disponível neste sistema",
	"The eyedropper could not read the screen":               "O conta-gotas não conseguiu ler a tela",
	"Click any pixel of the screen. Esc: cancel":             "Clique em qualquer pixel da tela. Esc: cancelar",
}

func init() {
	for k, v := range colourWords {
		translations[k] = v
	}
}

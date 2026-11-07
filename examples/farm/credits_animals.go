package main

import "github.com/gabrielluizsf/sge/farm"

// animalCredits are the animals' models, by others, under CC-BY 4.0, which
// asks for their authors to be named and the changes made to be said: see
// animals/SOURCE.txt.
func animalCredits() []farm.Credit {
	by := func() string { return ccBy4License }
	return []farm.Credit{
		{Name: "Farm Cow Animated Dairy Cattle (changed: walk and grazing added, rendered to sprites)", Author: "Rukh3D", URL: "https://sketchfab.com/3d-models/farm-cow-animated-dairy-cattle-0e780cb5ab6d457198ded7d4e191a03e", License: by},
		{Name: "Animated Sheep (changed: skeleton and motions added, rendered to sprites)", Author: "igor-lir", URL: "https://sketchfab.com/3d-models/animated-sheep-b99698502dea4905b916fce0bcf2dfc0", License: by},
		{Name: "Animated Wolf (changed: skeleton and motions added, rendered to sprites)", Author: "igor-lir", URL: "https://sketchfab.com/3d-models/animated-wolf-23a3a563d48049d5b9844ddc1b5b6e79", License: by},
		{Name: "Draft Horse Walk (changed: split in two horses, harness dyed as leather, rendered to sprites)", Author: "honved444", URL: "https://sketchfab.com/3d-models/draft-horse-walk-0541f4f11cb44383bf32b906aac5288e", License: by},
	}
}

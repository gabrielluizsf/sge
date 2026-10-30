package human

// macroTargets are the names of MakeHuman's macro targets, in its data's
// targets folder: by ancestry, by build, by height and by proportions.
func macroTargets() []string {
	var names []string
	genders := []string{"female", "male"}
	ages := []string{"baby", "child", "young", "old"}
	levels := func(s string) []string { return []string{"min" + s, "average" + s, "max" + s} }
	for _, race := range []string{"african", "asian", "caucasian"} {
		for _, g := range genders {
			for _, a := range ages {
				names = append(names, "macrodetails/"+race+"-"+g+"-"+a)
			}
		}
	}
	for _, g := range genders {
		for _, a := range ages {
			for _, m := range levels("muscle") {
				for _, w := range levels("weight") {
					build := g + "-" + a + "-" + m + "-" + w
					names = append(names, "macrodetails/universal-"+build)
					for _, h := range []string{"minheight", "maxheight"} {
						names = append(names, "macrodetails/height/"+build+"-"+h)
					}
					if a != "baby" {
						for _, p := range []string{"idealproportions", "uncommonproportions"} {
							names = append(names, "macrodetails/proportions/"+build+"-"+p)
						}
					}
				}
			}
		}
	}
	return names
}

// stem is a macro target's name without its folders, the words that
// weigh it.
func stem(name string) string {
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '/' {
			return name[i+1:]
		}
	}
	return name
}

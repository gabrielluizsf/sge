package human

import "image"

// eyesTexture is the colour texture of the eyes' material in file name.
func (k *Kit) eyesTexture(name string) (*image.NRGBA, error) {
	k.mu.Lock()
	defer k.mu.Unlock()
	data, err := k.lim.read(k.fsys, name)
	if err != nil {
		return nil, err
	}
	file, err := beside(name, parseMaterial(data).diffuse)
	if err != nil {
		return nil, err
	}
	return k.texture(file)
}

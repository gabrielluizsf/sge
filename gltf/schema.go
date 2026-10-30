package gltf

// The glTF JSON, as much of it as this package reads.
type document struct {
	Scene       *int           `json:"scene"`
	Scenes      []docScene     `json:"scenes"`
	Nodes       []docNode      `json:"nodes"`
	Meshes      []docMesh      `json:"meshes"`
	Materials   []docMaterial  `json:"materials"`
	Textures    []docTexture   `json:"textures"`
	Images      []docImage     `json:"images"`
	Accessors   []docAccessor  `json:"accessors"`
	BufferViews []docView      `json:"bufferViews"`
	Buffers     []docBuffer    `json:"buffers"`
	Skins       []docSkin      `json:"skins"`
	Animations  []docAnimation `json:"animations"`
}

type docScene struct {
	Nodes []int `json:"nodes"`
}

type docNode struct {
	Name        string    `json:"name"`
	Children    []int     `json:"children"`
	Mesh        *int      `json:"mesh"`
	Skin        *int      `json:"skin"`
	Matrix      []float64 `json:"matrix"`
	Translation []float64 `json:"translation"`
	Rotation    []float64 `json:"rotation"`
	Scale       []float64 `json:"scale"`
}

type docMesh struct {
	Name       string         `json:"name"`
	Primitives []docPrimitive `json:"primitives"`
}

type docPrimitive struct {
	Attributes map[string]int `json:"attributes"`
	Indices    *int           `json:"indices"`
	Material   *int           `json:"material"`
	Mode       *int           `json:"mode"`
}

type docMaterial struct {
	Name string `json:"name"`
	PBR  struct {
		BaseColorFactor  []float64 `json:"baseColorFactor"`
		BaseColorTexture *struct {
			Index int `json:"index"`
		} `json:"baseColorTexture"`
		MetallicRoughnessTexture *struct {
			Index int `json:"index"`
		} `json:"metallicRoughnessTexture"`
		MetallicFactor  *float64 `json:"metallicFactor"`
		RoughnessFactor *float64 `json:"roughnessFactor"`
	} `json:"pbrMetallicRoughness"`
	NormalTexture *struct {
		Index int      `json:"index"`
		Scale *float64 `json:"scale"`
	} `json:"normalTexture"`
	OcclusionTexture *struct {
		Index    int      `json:"index"`
		Strength *float64 `json:"strength"`
	} `json:"occlusionTexture"`
	AlphaMode   string   `json:"alphaMode"`
	AlphaCutoff *float64 `json:"alphaCutoff"`
	DoubleSided bool     `json:"doubleSided"`
}

type docTexture struct {
	Source *int `json:"source"`
}

type docImage struct {
	URI        string `json:"uri"`
	BufferView *int   `json:"bufferView"`
	MimeType   string `json:"mimeType"`
}

type docAccessor struct {
	BufferView    *int      `json:"bufferView"`
	ByteOffset    int64     `json:"byteOffset"`
	ComponentType int       `json:"componentType"`
	Normalized    bool      `json:"normalized"`
	Count         int64     `json:"count"`
	Type          string    `json:"type"`
	Sparse        *struct{} `json:"sparse"`
}

type docView struct {
	Buffer     int   `json:"buffer"`
	ByteOffset int64 `json:"byteOffset"`
	ByteLength int64 `json:"byteLength"`
	ByteStride int64 `json:"byteStride"`
}

type docBuffer struct {
	URI        string `json:"uri"`
	ByteLength int64  `json:"byteLength"`
}

type docSkin struct {
	Joints              []int `json:"joints"`
	InverseBindMatrices *int  `json:"inverseBindMatrices"`
}

type docAnimation struct {
	Name     string `json:"name"`
	Channels []struct {
		Sampler int `json:"sampler"`
		Target  struct {
			Node *int   `json:"node"`
			Path string `json:"path"`
		} `json:"target"`
	} `json:"channels"`
	Samplers []struct {
		Input         int    `json:"input"`
		Output        int    `json:"output"`
		Interpolation string `json:"interpolation"`
	} `json:"samplers"`
}

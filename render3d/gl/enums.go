//go:build linux && amd64

package gl

// GL's enums, those the device uses.
const (
	glTriangles          = 0x0004
	glLess               = 0x0201
	glBack               = 0x0405
	glCullFace           = 0x0B44
	glDepthTest          = 0x0B71
	glPackAlignment      = 0x0D05
	glUnpackAlignment    = 0x0CF5
	glTexture2D          = 0x0DE1
	glUnsignedByte       = 0x1401
	glUnsignedInt        = 0x1405
	glFloat              = 0x1406
	glDepthComponent     = 0x1902
	glRGBA               = 0x1908
	glVersion            = 0x1F02
	glNearest            = 0x2600
	glLinear             = 0x2601
	glLinearMipmapLinear = 0x2703
	glTexMagFilter       = 0x2800
	glTexMinFilter       = 0x2801
	glTexWrapS           = 0x2802
	glTexWrapT           = 0x2803
	glRepeat             = 0x2901
	glClampToEdge        = 0x812F
	glNone               = 0
	glCCW                = 0x0901
	glDepthBufferBit     = 0x0100
	glColorBufferBit     = 0x4000
	glRGBA8              = 0x8058
	glSRGB8Alpha8        = 0x8C43
	glDepthComponent24   = 0x81A6
	glDepthComponent32F  = 0x8CAC
	glTexture0           = 0x84C0
	glArrayBuffer        = 0x8892
	glElementArrayBuffer = 0x8893
	glUniformBuffer      = 0x8A11
	glStaticDraw         = 0x88E4
	glDynamicDraw        = 0x88E8
	glFragmentShader     = 0x8B30
	glVertexShader       = 0x8B31
	glCompileStatus      = 0x8B81
	glLinkStatus         = 0x8B82
	glInfoLogLength      = 0x8B84
	glFramebuffer        = 0x8D40
	glRenderbuffer       = 0x8D41
	glColorAttachment0   = 0x8CE0
	glDepthAttachment    = 0x8D00
	glFramebufferDone    = 0x8CD5 // GL_FRAMEBUFFER_COMPLETE
	glFramebufferSRGB    = 0x8DB9
)

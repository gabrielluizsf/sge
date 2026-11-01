//go:build linux && amd64

package gl

import (
	"errors"

	"github.com/gabrielluizsf/sge/native"
)

// procs are the GL functions the device calls, looked up once its
// context is current.
type procs struct {
	getError, getString, enable, disable, viewport          native.Proc
	clearColor, clearDepthf, clear, depthFunc, frontFace    native.Proc
	cullFace, pixelStorei, finish, readPixels               native.Proc
	genBuffers, deleteBuffers, bindBuffer, bufferData       native.Proc
	bufferSubData, bindBufferRange                          native.Proc
	genTextures, deleteTextures, bindTexture, activeTexture native.Proc
	texImage2D, texParameteri, generateMipmap               native.Proc
	genFramebuffers, deleteFramebuffers, bindFramebuffer    native.Proc
	framebufferTexture2D, framebufferRenderbuffer           native.Proc
	genRenderbuffers, deleteRenderbuffers, bindRenderbuffer native.Proc
	renderbufferStorage, checkFramebufferStatus             native.Proc
	drawBuffer, readBuffer                                  native.Proc
	createShader, shaderSource, compileShader, getShaderiv  native.Proc
	getShaderInfoLog, deleteShader, createProgram           native.Proc
	attachShader, linkProgram, getProgramiv                 native.Proc
	getProgramInfoLog, useProgram, deleteProgram            native.Proc
	genVertexArrays, bindVertexArray                        native.Proc
	enableVertexAttribArray, disableVertexAttribArray       native.Proc
	vertexAttribPointer, drawElementsInstanced              native.Proc
}

// named is a function's name and where its address goes.
type named struct {
	name string
	p    *native.Proc
}

// list is every function of g by its name.
func (g *procs) list() []named {
	return []named{
		{"glGetError", &g.getError}, {"glGetString", &g.getString}, {"glEnable", &g.enable},
		{"glDisable", &g.disable}, {"glViewport", &g.viewport}, {"glClearColor", &g.clearColor},
		{"glClearDepthf", &g.clearDepthf}, {"glClear", &g.clear}, {"glDepthFunc", &g.depthFunc},
		{"glFrontFace", &g.frontFace}, {"glCullFace", &g.cullFace}, {"glPixelStorei", &g.pixelStorei},
		{"glFinish", &g.finish}, {"glReadPixels", &g.readPixels},
		{"glGenBuffers", &g.genBuffers}, {"glDeleteBuffers", &g.deleteBuffers}, {"glBindBuffer", &g.bindBuffer},
		{"glBufferData", &g.bufferData}, {"glBufferSubData", &g.bufferSubData}, {"glBindBufferRange", &g.bindBufferRange},
		{"glGenTextures", &g.genTextures}, {"glDeleteTextures", &g.deleteTextures}, {"glBindTexture", &g.bindTexture},
		{"glActiveTexture", &g.activeTexture}, {"glTexImage2D", &g.texImage2D}, {"glTexParameteri", &g.texParameteri},
		{"glGenerateMipmap", &g.generateMipmap},
		{"glGenFramebuffers", &g.genFramebuffers}, {"glDeleteFramebuffers", &g.deleteFramebuffers},
		{"glBindFramebuffer", &g.bindFramebuffer}, {"glFramebufferTexture2D", &g.framebufferTexture2D},
		{"glFramebufferRenderbuffer", &g.framebufferRenderbuffer}, {"glGenRenderbuffers", &g.genRenderbuffers},
		{"glDeleteRenderbuffers", &g.deleteRenderbuffers}, {"glBindRenderbuffer", &g.bindRenderbuffer},
		{"glRenderbufferStorage", &g.renderbufferStorage}, {"glCheckFramebufferStatus", &g.checkFramebufferStatus},
		{"glDrawBuffer", &g.drawBuffer}, {"glReadBuffer", &g.readBuffer},
		{"glCreateShader", &g.createShader}, {"glShaderSource", &g.shaderSource}, {"glCompileShader", &g.compileShader},
		{"glGetShaderiv", &g.getShaderiv}, {"glGetShaderInfoLog", &g.getShaderInfoLog}, {"glDeleteShader", &g.deleteShader},
		{"glCreateProgram", &g.createProgram}, {"glAttachShader", &g.attachShader}, {"glLinkProgram", &g.linkProgram},
		{"glGetProgramiv", &g.getProgramiv}, {"glGetProgramInfoLog", &g.getProgramInfoLog},
		{"glUseProgram", &g.useProgram}, {"glDeleteProgram", &g.deleteProgram},
		{"glGenVertexArrays", &g.genVertexArrays}, {"glBindVertexArray", &g.bindVertexArray},
		{"glEnableVertexAttribArray", &g.enableVertexAttribArray},
		{"glDisableVertexAttribArray", &g.disableVertexAttribArray},
		{"glVertexAttribPointer", &g.vertexAttribPointer}, {"glDrawElementsInstanced", &g.drawElementsInstanced},
	}
}

// load looks every function of g up by getProc, eglGetProcAddress.
func (g *procs) load(getProc native.Proc) error {
	for _, f := range g.list() {
		*f.p = native.Proc(call(getProc, ptr(cstring(f.name))))
		if *f.p == 0 {
			return errors.New("gl: no " + f.name)
		}
	}
	return nil
}

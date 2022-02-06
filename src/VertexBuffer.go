package src

import (
	"unsafe"

	"github.com/go-gl/gl/v4.3-core/gl"
)

type VertexBuffer struct {
	rendererID uint32
}

func InitVertexBuffer(data unsafe.Pointer, size int) *VertexBuffer {
	buffer := VertexBuffer{}
	//Memory reference of buffer in VRAM
	gl.GenBuffers(1, &buffer.rendererID)

	//Selects the buffer
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer.rendererID)

	//Put data into buffer and specify size
	gl.BufferData(gl.ARRAY_BUFFER, size, data, gl.STATIC_DRAW)
	return &buffer
}

func (buffer VertexBuffer) Close() {
	gl.DeleteBuffers(1, &buffer.rendererID)
}

func (buffer VertexBuffer) Bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, buffer.rendererID)
}

func (buffer VertexBuffer) UnBind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

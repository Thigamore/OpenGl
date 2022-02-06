package src

import (
	"github.com/go-gl/gl/v4.3-core/gl"
)

type IndexBuffer struct {
	rendererID uint32
	count      int
}

func InitIndexBuffer(data []int, count int) *IndexBuffer {
	buffer := IndexBuffer{count: count}
	//Memory reference of buffer in VRAM
	gl.GenBuffers(1, &buffer.rendererID)

	//Selects the buffer
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, buffer.rendererID)

	//Put data into buffer and specify size
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, count*4, gl.Ptr(data), gl.STATIC_DRAW)
	return &buffer
}

func (buffer IndexBuffer) Close() {
	gl.DeleteBuffers(1, &buffer.rendererID)
}

func (buffer IndexBuffer) GetCount() int {
	return buffer.count
}

func (buffer IndexBuffer) Bind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, buffer.rendererID)
}

func (buffer IndexBuffer) UnBind() {
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}

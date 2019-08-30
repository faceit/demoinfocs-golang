package bitread

import "io"

type Reader interface {
	ReadString() string
	ReadFloat() float32
	ReadVarInt32() uint32
	ReadSignedVarInt32() int32
	ReadUBitInt() uint
	Pool()
	LazyPosition() int
	ActualPosition() int
	Open(underlying io.Reader, bufferSize int)
	OpenWithBuffer(underlying io.Reader, buffer []byte)
	Close()
	ReadBit() bool
	ReadBits(n int) []byte
	ReadSingleByte() byte
	ReadBitsToByte(n int) byte
	ReadInt(n int) uint
	ReadBytes(n int) []byte
	ReadBytesInto(out *[]byte, n int)
	ReadCString(n int) string
	ReadSignedInt(n int) int
	BeginChunk(n int)
	EndChunk()
	ChunkFinished() bool
	Skip(n int)
}

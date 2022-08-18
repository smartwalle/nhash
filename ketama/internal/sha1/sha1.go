package sha1

import (
	"crypto/sha1"
	"hash"
)

type digest struct {
	hash.Hash
}

func New() hash.Hash32 {
	var h = &digest{}
	h.Hash = sha1.New()
	h.Reset()
	return h
}

func (this *digest) Sum32() uint32 {
	var vBytes = this.Sum(nil)
	return uint32(vBytes[19]) | uint32(vBytes[18])<<8 | uint32(vBytes[17])<<16 | uint32(vBytes[16])<<24
}

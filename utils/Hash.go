package utils

import (
	"math/rand"
)

const regMask = uint64(1274391726635)

type Hash uint64

type HashGen struct {
	uniq uint64
}

func NewGen() *HashGen {
	val := uint64(0)
	amountOfBits := CountBits(val)
	for amountOfBits < 24 || amountOfBits > 40 {
		// Ensure that rand is seeded at initialization
		val ^= rand.Uint64()
		amountOfBits = CountBits(val)
	}
	return &HashGen{
		uniq: val,
	}
}

func (gen HashGen) String(s string) Hash {
	hashVal := RotateRight(gen.uniq, 16)
	hashReg := RotateLeft(gen.uniq^regMask, 16)
	hash := make([]uint8, 8)
	i := 0
	for _, v := range s {
		hash[i&7] = (uint8(hashVal) & uint8(^v)) ^ (uint8(hashReg) & uint8(v))
		i++
		hashVal = RotateRight(hashVal, 8)
		hashReg = RotateLeft(hashReg, 8)
	}
	retVal := uint64(0)
	for _, v := range hash {
		retVal |= uint64(v)
		retVal <<= 8
	}
	return Hash(retVal)
}

package hash

import (
	"SimpleRTMPServer/utils"
	"math/rand"
)

const regMask = uint64(1274391726635)

type Type uint64

type Gen struct {
	uniq uint64
}

func InitGen() Gen {
	val := uint64(0)
	amountOfBits := utils.CountBits(val)
	for amountOfBits < 24 || amountOfBits > 40 {
		// Ensure that rand is seeded at initialization
		val ^= rand.Uint64()
		amountOfBits = utils.CountBits(val)
	}
	return Gen{
		uniq: val,
	}
}

func (gen Gen) String(s string) Type {
	hashVal := utils.RotateBitsRight(gen.uniq, 16)
	hashReg := utils.RotateBitsLeft(gen.uniq^regMask, 16)
	hash := make([]uint8, 8)
	i := 0
	for _, v := range s {
		hash[i&7] = (uint8(hashVal) & uint8(^v)) ^ (uint8(hashReg) & uint8(v))
		i++
		hashVal = utils.RotateBitsRight(hashVal, 8)
		hashReg = utils.RotateBitsLeft(hashReg, 8)
	}
	retVal := uint64(0)
	for _, v := range hash {
		retVal |= uint64(v)
		retVal <<= 8
	}
	return Type(retVal)
}

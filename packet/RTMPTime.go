package packet

type RTMPTime uint32

func (t RTMPTime) Uint32() uint32 {
	return uint32(t)
}

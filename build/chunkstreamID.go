package build

import "errors"

func chunkStreamID(fmt, csID int) ([]byte, error) {
	fmt = fmt & 3
	fmt = fmt << 6
	if csID < 2 {
		return nil, errors.New("Unsupported csID")
	} else if csID < 64 {
		return []byte{byte(fmt | csID)}, nil
	} else if csID <= 319 {
		csID = csID - 64
		return []byte{byte(fmt | 1), byte(csID - 64)}, nil
	} else {
		csID = csID - 64
		return []byte{byte(fmt | 2), byte(csID % 256), byte(csID / 256)}, nil
	}
}

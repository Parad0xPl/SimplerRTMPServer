package build

import "errors"

func chunkstreamID(fmt, csid int) ([]byte, error) {
	fmt = fmt & 3
	fmt = fmt << 6
	if csid < 2 {
		return nil, errors.New("Unsupported csid")
	} else if csid < 64 {
		return []byte{byte(fmt | csid)}, nil
	} else if csid <= 319 {
		csid = csid - 64
		return []byte{byte(fmt | 1), byte(csid - 64)}, nil
	} else {
		csid = csid - 64
		return []byte{byte(fmt | 2), byte(csid % 256), byte(csid / 256)}, nil
	}
}

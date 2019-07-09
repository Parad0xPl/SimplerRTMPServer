package utils

// Concat slices with prealloc
// https://stackoverflow.com/questions/37884361/concat-multiple-slices-in-golang/40678026
func Concat(slcs ...[]byte) []byte {
	var finlen int
	for _, i := range slcs {
		finlen += len(i)
	}

	slc := make([]byte, finlen)
	var i int
	for _, s := range slcs {
		i += copy(slc[i:], s)
	}
	return slc
}

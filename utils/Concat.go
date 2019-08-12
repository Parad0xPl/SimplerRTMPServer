package utils

// Concat slices with prealloc
// https://stackoverflow.com/questions/37884361/concat-multiple-slices-in-golang/40678026
func Concat(slices ...[]byte) []byte {
	var finalLen int
	for _, i := range slices {
		finalLen += len(i)
	}

	slc := make([]byte, finalLen)
	var i int
	for _, s := range slices {
		i += copy(slc[i:], s)
	}
	return slc
}

package packet

type headerCacheMap map[int]*Header

type HeadersCache struct {
	cache headerCacheMap
}

func NewHeadersCache() *HeadersCache {
	return &HeadersCache{
		cache: make(headerCacheMap, 100),
	}
}

// TODO Optimize memory. Possible dos when client will send many different chunkID
// Insert Header into cache
func (h HeadersCache) Insert(chunkID int, val *Header) {
	h.cache[chunkID] = val
}

// Get matched Header from cache
func (h HeadersCache) Get(chunkID int) *Header {
	return h.cache[chunkID]
}

package types

import "strconv"

// id truncates and deduplicates identifiers to ensure that the identifier
// stays within limits of databases and do not conflict with each other.
//
// When called with the same sqlTypeConverter, the returned identifier will
// always be the same if the provided identifier is identical.
func (s sqlTypeConverter) id(original string) (safe string) {
	// Look in cache.
	if mapped, ok := s.idMappedCache[original]; ok {
		return mapped
	}
	defer func() {
		s.idMappedCache[original] = safe
	}()
	// Truncate to 50 characters.
	truncated := original
	if len(original) > 50 {
		truncated = original[:50]
	}
	// Append a number if it will cause a collision otherwise.
	count := s.idUsedCache[truncated]
	s.idUsedCache[truncated]++
	if count != 0 {
		return truncated + "_" + strconv.Itoa(count)
	}
	return truncated
}

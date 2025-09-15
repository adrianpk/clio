package am

// ToPtrSlice converts a slice of values to a slice of pointers.
func ToPtrSlice[T any](s []T) []*T {
	p := make([]*T, len(s))
	for i := range s {
		p[i] = &s[i]
	}
	return p
}

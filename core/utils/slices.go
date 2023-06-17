package utils

func Filter[T any](a []T, predicate func(b T) bool) *T {
	for _, v := range a {
		v := v
		if predicate(v) {
			return &v
		}
	}
	return nil
}

func AllMatches[T any](a []T, predicate func(b T) bool) []T {
	var l = []T{}
	for _, v := range a {
		v := v
		if predicate(v) {
			l = append(l, v)
		}
	}
	return l
}

func AnyMatch[T any](a []T, predicate func(b T) bool) bool {
	match := Filter(a, predicate)
	return match != nil
}

func AnyMatchString(a []string, find string) bool {
	match := Filter(a, func(b string) bool {
		return b == find
	})
	return match != nil
}

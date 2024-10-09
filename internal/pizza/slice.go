package pizza

func Map[I, R interface{}](objs []I, fn func(obj I) R) []R {
	s := []R{}
	for _, obj := range objs {
		s = append(s, fn(obj))
	}
	return s
}

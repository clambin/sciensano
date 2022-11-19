package set

type Set map[string]struct{}

func (s *Set) IsNew(name string) bool {
	if *s == nil {
		*s = make(map[string]struct{})
	}
	_, exists := (*s)[name]
	if exists {
		return false
	}
	(*s)[name] = struct{}{}
	return true
}

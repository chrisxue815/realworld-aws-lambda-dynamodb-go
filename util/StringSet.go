package util

type StringSet map[string]bool

func NewStringSetFromSlice(slice []string) StringSet {
	s := make(StringSet, len(slice))
	for _, value := range slice {
		s[value] = true
	}
	return s
}

func (s StringSet) Difference(other StringSet) StringSet {
	result := make(StringSet)
	for value := range s {
		if !other[value] {
			result[value] = true
		}
	}
	return result
}

func (s StringSet) ToSlice() []string {
	result := make([]string, 0, len(s))
	for value := range s {
		result = append(result, value)
	}
	return result
}

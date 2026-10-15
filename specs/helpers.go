package specs

import "maps"

func (s Specs) Has(name string) bool {
	_, ok := s[name]
	return ok
}

func (s Specs) Get(name string) (string, bool) {
	value, ok := s[name]
	return value, ok
}

func (s Specs) GetOr(name, fallback string) string {
	if value, ok := s[name]; ok {
		return value
	}

	s[name] = fallback

	return fallback
}

func (s Specs) Set(name, value string) {
	s[name] = value
}

func (s Specs) Delete(name string) Specs {
	delete(s, name)
	return s
}

func (s Specs) Clone() Specs {
	result := make(Specs, len(s))

	maps.Copy(result, s)

	return result
}

func (s Specs) Merge(other Specs) Specs {
	maps.Copy(s, other)
	return s
}

func (s Specs) IsEmpty() bool {
	return len(s) == 0
}

func (s Specs) Len() int {
	return len(s)
}

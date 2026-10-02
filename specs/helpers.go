package specs

func (s Specs) Has(name string) bool {
	_, ok := s[name]
	return ok
}

func (s Specs) Get(name string) string {
	return s[name]
}

func (s Specs) GetOr(name, fallback string) string {
	if value, ok := s[name]; ok {
		return value
	}

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

	for key, value := range s {
		result[key] = value
	}

	return result
}

func (s Specs) Merge(other Specs) Specs {
	for key, value := range other {
		s[key] = value
	}

	return s
}

func (s Specs) IsEmpty() bool {
	return len(s) == 0
}

func (s Specs) Len() int {
	return len(s)
}

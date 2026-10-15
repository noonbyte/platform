package specs

import (
	"errors"
	"fmt"
	"strings"
)

type Sort struct {
	Field string
	Order string
}

func (s Specs) Sort() Sort {
	field, _ := s.GetString("sort_by")
	order, _ := s.GetString("sort_order")

	if order == "" {
		order = "asc"
	}

	return Sort{
		Field: field,
		Order: strings.ToLower(order),
	}
}

var (
	ErrUnsupportedSortOrder = errors.New("unsupported sort order")
)

func (s Sort) Resolve() (string, error) {
	if s.Field == "" {
		return "", nil
	}

	order := strings.ToLower(s.Order)
	if order == "" {
		order = "asc"
	}

	if order != "asc" && order != "desc" {
		return "", fmt.Errorf("%w: %s", ErrUnsupportedSortOrder, s.Field)
	}

	return s.Field + " " + order, nil
}

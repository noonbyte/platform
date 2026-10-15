package specs

type Pagination struct {
	Page  int64 `json:"page"`
	Limit int64 `json:"limit"`
}

type PaginationConfig struct {
	DefaultLimit int64
	MaxLimit     int64
}

func DefaultPaginationConfig() PaginationConfig {
	return PaginationConfig{
		DefaultLimit: 20,
		MaxLimit:     100,
	}
}

func (s Specs) Pagination() (Pagination, error) {
	return s.PaginationWithConfig(DefaultPaginationConfig())
}

func (s Specs) PaginationWithConfig(
	config PaginationConfig,
) (Pagination, error) {
	if config.DefaultLimit < 1 {
		config.DefaultLimit = 20
	}

	if config.MaxLimit < config.DefaultLimit {
		config.MaxLimit = 100
	}

	p := Pagination{
		Page:  1,
		Limit: config.DefaultLimit,
	}

	page, exists, err := s.GetInt64("page")
	if err != nil {
		return p, err
	}
	if exists {
		p.Page = page
	}

	limit, exists, err := s.GetInt64("limit")
	if err != nil {
		return p, err
	}
	if exists {
		p.Limit = limit
	}

	if p.Page < 1 {
		p.Page = 1
	}

	if p.Limit < 1 {
		p.Limit = config.DefaultLimit
	}

	if p.Limit > config.MaxLimit {
		p.Limit = config.MaxLimit
	}

	return p, nil
}

func (p Pagination) Offset() int64 {
	return (p.Page - 1) * p.Limit
}

func (p Pagination) HasMore(total int64) bool {
	return p.Offset()+p.Limit < total
}

func (p Pagination) Meta(total int64) PageMeta {
	return PageMeta{
		Total:   total,
		Page:    p.Page,
		Limit:   p.Limit,
		HasMore: p.HasMore(total),
	}
}

type PageMeta struct {
	Total   int64 `json:"total"`
	Page    int64 `json:"page"`
	Limit   int64 `json:"limit"`
	HasMore bool  `json:"has_more"`
}

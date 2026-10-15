package specs

type Pagination struct {
	Specs  Specs
	Config PaginationConfig
}

type PaginationConfig struct {
	DefaultLimit int
	MaxLimit     int
}

func DefaultPaginationConfig() PaginationConfig {
	return PaginationConfig{
		DefaultLimit: 20,
		MaxLimit:     100,
	}
}

func (s Specs) Pagination(config PaginationConfig) Pagination {
	return Pagination{
		Specs:  s,
		Config: config,
	}
}

func (p Pagination) Set(page, limit *int) {
	config := p.Config

	if config.DefaultLimit <= 0 {
		config.DefaultLimit = 20
	}

	if config.MaxLimit <= 0 {
		config.MaxLimit = 100
	}

	pageValue := 1
	if page != nil && *page > 0 {
		pageValue = *page
	}

	limitValue := config.DefaultLimit
	if limit != nil && *limit > 0 {
		limitValue = *limit
	}

	if limitValue > config.MaxLimit {
		limitValue = config.MaxLimit
	}

	p.Specs.WithInt("page", &pageValue)
	p.Specs.WithInt("limit", &limitValue)
}

func (p Pagination) Page() int64 {
	value, _, _ := p.Specs.GetInt("page")
	return int64(value)
}

func (p Pagination) Limit() int64 {
	value, _, _ := p.Specs.GetInt("limit")
	return int64(value)
}

func (p Pagination) Offset() int64 {
	return (p.Page() - 1) * p.Limit()
}

func (p Pagination) HasMore(total int64) bool {
	return p.Offset()+p.Limit() < total
}

func (p Pagination) Meta(total int64) PageMeta {
	return PageMeta{
		Total:   total,
		Page:    p.Page(),
		Limit:   p.Limit(),
		HasMore: p.HasMore(total),
	}
}

type PageMeta struct {
	Total   int64 `json:"total"`
	Page    int64 `json:"page"`
	Limit   int64 `json:"limit"`
	HasMore bool  `json:"has_more"`
}

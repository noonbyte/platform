package specs

const (
	DefaultPage  = 1
	DefaultLimit = 20
	MaxLimit     = 100
)

type Pagination struct {
	specs Specs
}

type PageMeta struct {
	Total   int  `json:"total"`
	Page    int  `json:"page"`
	Limit   int  `json:"limit"`
	HasMore bool `json:"has_more"`
}

func (s Specs) Pagination() Pagination {
	return Pagination{
		specs: s,
	}
}

func (p Pagination) Set(page, limit *int) {
	pageValue := DefaultPage
	if page != nil && *page > 0 {
		pageValue = *page
	}

	limitValue := DefaultLimit
	if limit != nil && *limit > 0 {
		limitValue = *limit
	}

	if limitValue > MaxLimit {
		limitValue = MaxLimit
	}

	p.specs.WithInt("page", &pageValue)
	p.specs.WithInt("limit", &limitValue)
}

func (p Pagination) Page() int {
	value, _, _ := p.specs.GetInt("page")
	return value
}

func (p Pagination) Limit() int {
	value, _, _ := p.specs.GetInt("limit")
	return value
}

func (p Pagination) Offset() int {
	return (p.Page() - 1) * p.Limit()
}

func (p Pagination) HasMore(total int) bool {
	return p.Offset()+p.Limit() < total
}

func (p Pagination) Meta(total int) PageMeta {
	return PageMeta{
		Total:   total,
		Page:    p.Page(),
		Limit:   p.Limit(),
		HasMore: p.HasMore(total),
	}
}

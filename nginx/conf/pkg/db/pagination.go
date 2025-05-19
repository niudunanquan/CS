package db

type Pagination struct {
	Size  int         `json:"size"`
	Page  int         `json:"page"`
	Sort  string      `json:"sort,omitempty"`
	Total int64       `json:"total"`
	Data  interface{} `json:"data"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetSize()
}

func (p *Pagination) GetSize() int {
	if p.Size == 0 {
		p.Size = 10
	}
	return p.Size
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "id desc"
	}
	return p.Sort
}

func (p *Pagination) Pages() int {
	if p.Total == 0 || p.Size == 0 {
		return 0
	}
	totalPage := int(p.Total) / p.Size
	if int(p.Total)%p.Size > 0 {
		totalPage = totalPage + 1
	}
	return totalPage
}

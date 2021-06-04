package pagination

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	defaultPageSize = 100
	maxPageSize     = 1000
	pageVar         = "page"
	pageSizeVar     = "per_page"
)

type Pages struct {
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	PageCount  int         `json:"page_count"`
	TotalCount int         `json:"total_count"`
	Items      interface{} `json:"item"`
}

func New(page, perPage, total int) *Pages {
	if perPage <= 0 {
		perPage = defaultPageSize
	} else if perPage > maxPageSize {
		perPage = maxPageSize
	}

	pageCount := -1
	if total >= 0 {
		pageCount = (total + perPage - 1) / perPage
		if page > pageCount {
			page = pageCount
		}
	}
	if page < 1 {
		page = 1
	}

	return &Pages{
		Page:       page,
		PerPage:    perPage,
		PageCount:  pageCount,
		TotalCount: total,
	}
}

func parseInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	if r, err := strconv.Atoi(value); err == nil {
		return r
	}
	return defaultValue
}

func NewFromRequest(req *http.Request, count int) *Pages {
	page := parseInt(req.URL.Query().Get(pageVar), 1)
	perPage := parseInt(req.URL.Query().Get(pageSizeVar), defaultPageSize)
	return New(page, perPage, count)
}

func (p *Pages) Offset() int {
	return (p.Page - 1) * p.PerPage
}

// Limit returns the limitation value can be used in a SQL statement
func (p *Pages) Limit() int {
	return p.PerPage
}

const pageLinkPerPageCount = 4

func (p *Pages) BuildLinks(baseURL string, defaultPerPage int) [pageLinkPerPageCount]string {
	var links [pageLinkPerPageCount]string
	pageCount := p.PageCount
	page := p.Page

	buildLink := func(v int) string {
		return fmt.Sprintf("%v%v=%v", baseURL, pageVar, v)
	}

	if pageCount >= 0 && page > pageCount {
		page = pageCount
	}
	if strings.Contains(baseURL, "?") {
		baseURL += "&"
	} else {
		baseURL += "?"
	}

	if page > 1 {
		links[0] = buildLink(1)
		links[1] = buildLink(page - 1)
	}

	if pageCount >= 0 && page < pageCount {
		links[2] = buildLink(page + 1)
		links[3] = buildLink(pageCount)
	} else if pageCount < 0 {
		links[2] = buildLink(page + 1)
	}
	if perPage := p.PerPage; perPage != defaultPerPage {
		for i := 0; i < pageLinkPerPageCount; i++ {
			if links[i] != "" {
				links[i] += fmt.Sprintf("&%v=%v", pageSizeVar, perPage)
			}
		}
	}
	return links
}

func (p *Pages) BuildLinkHeader(baseURL string, defaultPerPage int) string {
	links := p.BuildLinks(baseURL, defaultPerPage)
	header := ""
	if links[0] != "" {
		header += fmt.Sprintf(`<%v>; rel="first", `, links[0])
		header += fmt.Sprintf(`<%v>; rel="prev"`, links[1])
	}
	if links[2] != "" {
		if header != "" {
			header += ", "
		}
		header += fmt.Sprintf(`<%v>; rel="next"`, links[2])
		if links[3] != "" {
			header += fmt.Sprintf(`, <%v>; rel="last"`, links[3])
		}
	}

	return header
}

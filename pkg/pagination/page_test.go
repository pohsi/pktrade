package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	data := []struct {
		tag                                                                    string
		page, perPage, total                                                   int
		expectedPage, expectedPerPage, expectedTotal, pageCount, offset, limit int
	}{
		// varying page
		{"t1", 1, 20, 50, 1, 20, 50, 3, 0, 20},
		{"t2", 2, 20, 50, 2, 20, 50, 3, 20, 20},
		{"t3", 3, 20, 50, 3, 20, 50, 3, 40, 20},
		{"t4", 4, 20, 50, 3, 20, 50, 3, 40, 20},
		{"t5", 0, 20, 50, 1, 20, 50, 3, 0, 20},

		// varying perPage
		{"t6", 1, 0, 50, 1, 100, 50, 1, 0, 100},
		{"t7", 1, -1, 50, 1, 100, 50, 1, 0, 100},
		{"t8", 1, 100, 50, 1, 100, 50, 1, 0, 100},
		{"t9", 1, 1001, 50, 1, 1000, 50, 1, 0, 1000},

		// varying total
		{"t10", 1, 20, 0, 1, 20, 0, 0, 0, 20},
		{"t10", 1, 20, -1, 1, 20, -1, -1, 0, 20},
	}

	for _, d := range data {
		p := New(d.page, d.perPage, d.total)
		assert.Equal(t, d.expectedPage, p.Page, d.tag)
		assert.Equal(t, d.expectedPerPage, p.PerPage, d.tag)
		assert.Equal(t, d.expectedTotal, p.TotalCount, d.tag)
		assert.Equal(t, d.pageCount, p.PageCount, d.tag)
		assert.Equal(t, d.offset, p.Offset(), d.tag)
		assert.Equal(t, d.limit, p.Limit(), d.tag)
	}
}

func TestBuildLinkHeader(t *testing.T) {
	baseURL := "/tokens"
	defaultPerPage := 10
	data := []struct {
		tag                  string
		page, perPage, total int
		header               string
	}{
		{"t1", 1, 20, 50, `</tokens?page=2&per_page=20>; rel="next", </tokens?page=3&per_page=20>; rel="last"`},
		{"t2", 2, 20, 50, `</tokens?page=1&per_page=20>; rel="first", </tokens?page=1&per_page=20>; rel="prev", </tokens?page=3&per_page=20>; rel="next", </tokens?page=3&per_page=20>; rel="last"`},
		{"t3", 3, 20, 50, `</tokens?page=1&per_page=20>; rel="first", </tokens?page=2&per_page=20>; rel="prev"`},
		{"t4", 0, 20, 50, `</tokens?page=2&per_page=20>; rel="next", </tokens?page=3&per_page=20>; rel="last"`},
		{"t5", 4, 20, 50, `</tokens?page=1&per_page=20>; rel="first", </tokens?page=2&per_page=20>; rel="prev"`},
		{"t6", 1, 20, 0, ""},
		{"t7", 4, 20, -1, `</tokens?page=1&per_page=20>; rel="first", </tokens?page=3&per_page=20>; rel="prev", </tokens?page=5&per_page=20>; rel="next"`},
	}

	for _, v := range data {
		p := New(v.page, v.perPage, v.total)
		assert.Equal(t, v.header, p.BuildLinkHeader(baseURL, defaultPerPage), v.tag)
	}

	baseURL = `/tokens?from=10`
	p := New(1, 20, 50)
	assert.Equal(t, `</tokens?from=10&page=2&per_page=20>; rel="next", </tokens?from=10&page=3&per_page=20>; rel="last"`,
		p.BuildLinkHeader(baseURL, defaultPerPage))
}

func Test_parseInt(t *testing.T) {
	type args struct {
		value        string
		defaultValue int
	}
	data := []struct {
		name string
		args
		want int
	}{
		{"t1", args{"123", 100}, 123},
		{"t2", args{"", 100}, 100},
		{"t3", args{"a", 100}, 100},
	}

	for _, v := range data {
		t.Run(v.name, func(t *testing.T) {
			if r := parseInt(v.value, v.defaultValue); r != v.want {
				t.Errorf("parseInt() = %v, want %v", r, v.want)
			}
		})
	}
}

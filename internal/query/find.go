package query

import (
	"bytes"
	"html"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/kirillgrachoff/optparf-check/internal/types"
)

var (
	rowSep    = []byte(`<div class="row table-body">`)
	articleRe = regexp.MustCompile(`<div class="col-auto cols-width">\s*([^<\s]+)\s*</div>`)
	nameRe    = regexp.MustCompile(`(?s)<a href="([^"]*)">(.*?)</a>`)
	priceRe   = regexp.MustCompile(`svc-price="(\d+)">([^<]*)<`)
	tagRe     = regexp.MustCompile(`<[^>]*>`)
)

type Parser struct {
	base *url.URL
}

// NewParser creates a parser; base is used to resolve relative item links.
func NewParser(base string) *Parser {
	u, err := url.Parse(base)
	if err != nil {
		u = &url.URL{}
	}
	return &Parser{base: u}
}

func (p Parser) Parse(in []byte) []types.Item {
	var items []types.Item
	rows := bytes.Split(in, rowSep)
	for _, row := range rows[1:] {
		price := priceRe.FindSubmatch(row)
		if price == nil {
			continue
		}
		item := types.Item{
			Id:    string(price[1]),
			Price: parsePrice(string(price[2])),
		}
		if m := articleRe.FindSubmatch(row); m != nil {
			item.Article = string(m[1])
		}
		if m := nameRe.FindSubmatch(row); m != nil {
			item.URL = p.resolve(html.UnescapeString(string(m[1])))
			item.Name = strings.Join(strings.Fields(html.UnescapeString(tagRe.ReplaceAllString(string(m[2]), ""))), " ")
		}
		items = append(items, item)
	}
	return items
}

func (p Parser) resolve(href string) string {
	u, err := url.Parse(href)
	if err != nil {
		return href
	}
	return p.base.ResolveReference(u).String()
}

func parsePrice(s string) float64 {
	s = strings.ReplaceAll(s, "&nbsp;", "")
	s = strings.Join(strings.Fields(s), "")
	s = strings.ReplaceAll(s, ",", ".")
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

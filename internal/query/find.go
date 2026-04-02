package query

import (
	"bytes"
	"regexp"
)

type Filter struct {
	catch  *regexp.Regexp
	delete string
}

func NewFilter(catch, delete string) *Filter {
	return &Filter{
		catch: regexp.MustCompile(catch),
		delete: delete,
	}
}

func (f Filter) Filter(in []byte) []byte {
	splitted := bytes.Split(in, []byte{'\n'})
	var result [][]byte
	for _, b := range splitted {
		if f.catch.Match(b) {
			b = bytes.ReplaceAll(b, []byte("&nbsp;"), nil)
			result = append(result, bytes.Trim(b, " \t"))
		}
	}
	return bytes.Join(result, []byte{'\n'})
}

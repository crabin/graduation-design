package gobusterfuzz

import (
	"github.com/OJ/gobuster/v3/libgobuster"
)

// OptionsFuzz is the struct to hold all options for this plugin
type OptionsFuzz struct {
	libgobuster.HTTPOptions
	ExcludedStatusCodes       string
	ExcludedStatusCodesParsed libgobuster.Set[int]
	ExcludeLength             []int
	RequestBody               string
}

// NewOptionsFuzz returns a new initialized OptionsFuzz
func NewOptionsFuzz() *OptionsFuzz {
	return &OptionsFuzz{
		ExcludedStatusCodesParsed: libgobuster.NewSet[int](),
	}
}

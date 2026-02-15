package ref

import (
	"fmt"
	"strings"
)

// RefFormatter formats a slice of resolved references into a string representation.
type RefFormatter interface {
	Format(resolved []*Resolved) (string, error)
}

// GetFormatter returns a formatter for the given style name.
func GetFormatter(style string) (RefFormatter, error) {
	switch style {
	case "canonical":
		return &canonicalFormatter{}, nil
	case "abbr":
		return &abbrFormatter{}, nil
	case "2letter":
		return &nLetterFormatter{n: 2, withPeriod: false}, nil
	case "3letter":
		return &nLetterFormatter{n: 3, withPeriod: false}, nil
	case "2letter.":
		return &nLetterFormatter{n: 2, withPeriod: true}, nil
	case "3letter.":
		return &nLetterFormatter{n: 3, withPeriod: true}, nil
	default:
		return nil, fmt.Errorf("unknown style %q", style)
	}
}

// GetAvailableStyles returns a list of all available style names.
func GetAvailableStyles() []string {
	return []string{
		"canonical",
		"abbr",
		"2letter",
		"3letter",
		"2letter.",
		"3letter.",
	}
}

// canonicalFormatter formats references with full book names.
type canonicalFormatter struct{}

func (f *canonicalFormatter) Format(resolved []*Resolved) (string, error) {
	refs := make([]string, 0, len(resolved))
	for _, r := range resolved {
		ref, err := r.CompactRef()
		if err != nil {
			return "", err
		}
		refs = append(refs, ref)
	}
	return strings.Join(refs, "; "), nil
}

// abbrFormatter formats references with preferred abbreviations.
type abbrFormatter struct{}

func (f *abbrFormatter) Format(resolved []*Resolved) (string, error) {
	refs := make([]string, 0, len(resolved))
	for _, r := range resolved {
		ref, err := r.AbbreviatedRef(WithAbbreviations(Abbreviations))
		if err != nil {
			return "", err
		}
		refs = append(refs, ref)
	}
	return strings.Join(refs, "; "), nil
}

// nLetterFormatter formats references with N-letter abbreviations.
type nLetterFormatter struct {
	n          int
	withPeriod bool
}

func (f *nLetterFormatter) Format(resolved []*Resolved) (string, error) {
	refs := make([]string, 0, len(resolved))
	for _, r := range resolved {
		abbr, err := Abbreviations.NLetterAbbreviation(r.Book.Name, f.n, f.withPeriod)
		if err != nil {
			return "", err
		}
		ref, err := formatResolvedWithName(r, abbr)
		if err != nil {
			return "", err
		}
		refs = append(refs, ref)
	}
	return strings.Join(refs, "; "), nil
}

// formatResolvedWithName formats a resolved reference with a custom book name.
// It delegates to Resolved.compactRef to avoid duplicating formatting logic.
func formatResolvedWithName(r *Resolved, name string) (string, error) {
	return r.compactRef(name)
}

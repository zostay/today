package flag

import (
	"strings"

	"github.com/zostay/today/cmd/output"
)

// OutputFormat is a flag type that represents one of the supported output
type OutputFormat struct {
	Value output.Format
}

// String returns the string representation of the output format.
func (o *OutputFormat) String() string {
	if o.Value.Name == "" {
		return output.DefaultFormat().Name
	}
	return o.Value.Name
}

// IsSet returns true if the output format has been set.
func (o *OutputFormat) IsSet() bool {
	return o.Value.Name != ""
}

// Set parses the input string and sets the output format value.
func (o *OutputFormat) Set(v string) error {
	var err error
	o.Value, err = output.LookupFormat(v)
	if err != nil {
		return err
	}
	return nil
}

// Type returns the type of the flag.
func (o *OutputFormat) Type() string {
	return "format"
}

func ListOutputFormats() string {
	return strings.Join(output.ListFormats(), ", ")
}

package flag

import (
	"time"

	"github.com/markusmobius/go-dateparser"
	"github.com/markusmobius/go-dateparser/date"
)

// Date is a flag type that represents a date string in various formats.
type Date struct {
	Value date.Date
}

// String returns the string representation of the date.
func (d *Date) String() string {
	return d.Value.Time.Format(time.RFC3339)
}

// Set parses the input string and sets the date value.
func (d *Date) Set(v string) (err error) {
	d.Value, err = dateparser.Parse(nil, v)
	if err != nil {
		return err
	}
	return nil
}

// Type returns the type of the flag.
func (d *Date) Type() string {
	return "datestring"
}

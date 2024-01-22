package flag

import (
	"time"

	"github.com/markusmobius/go-dateparser"
	"github.com/markusmobius/go-dateparser/date"
)

type Date struct {
	Value date.Date
}

func (d *Date) String() string {
	return d.Value.Time.Format(time.RFC3339)
}

func (d *Date) Set(v string) (err error) {
	d.Value, err = dateparser.Parse(nil, v)
	if err != nil {
		return err
	}
	return nil
}

func (d *Date) Type() string {
	return "datestring"
}

package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Proposals struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Created int64  `json:"created"`
	State   string `json:"state"`
	Space   Space  `json:"space"`
}

type Space struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// NullableTime represents a chat member.
// @Description nil
type NullableTime struct {
	Time  *time.Time `json:"time,omitempty"`
	Valid bool       `json:"valid"`
}

func (nt *NullableTime) Scan(value interface{}) error {
	if value == nil {
		nt.Time = nil
		nt.Valid = false
		return nil
	}

	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("NullableTime: cannot scan type %T into NullableTime", value)
	}

	nt.Time = &t
	nt.Valid = true
	return nil
}

func (nt NullableTime) Value() (driver.Value, error) {
	if !nt.Valid || nt.Time == nil {
		return nil, nil
	}
	return *nt.Time, nil
}

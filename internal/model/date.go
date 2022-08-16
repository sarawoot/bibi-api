package model

import (
	"fmt"
	"strings"
	"time"
)

type Date time.Time

func (d Date) MarshalJSON() ([]byte, error) {
	dateStr := ""
	if !time.Time(d).IsZero() {
		dateStr = time.Time(d).Local().Format("2006-01-02")
	}

	return []byte(fmt.Sprintf(`"%s"`, dateStr)), nil
}

func (d *Date) UnmarshalJSON(b []byte) error {
	timeString := strings.Trim(string(b), `"`)

	t, err := time.Parse("2006-01-02", timeString)
	if err == nil {
		*d = Date(t.Local())
		return nil
	}

	return fmt.Errorf("invalid date format: %s", timeString)
}

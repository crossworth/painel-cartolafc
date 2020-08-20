package updater

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Durations []time.Duration

func (ds Durations) Value() (driver.Value, error) {
	stringVals := []string{}
	for _, d := range ds {
		stringVals = append(stringVals, d.String())
	}

	return "{" + strings.Join(stringVals, ",") + "}", nil
}

func (ds *Durations) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	source, ok := src.([]byte)
	if !ok {
		return errors.New("source failed type assertion to []byte")
	}
	var durs []time.Duration

	source = source[1 : len(source)-1]
	if len(source) == 0 {
		return nil
	}

	parts := bytes.Split(source, []byte(","))
	for _, p := range parts {
		d, err := time.ParseDuration(string(p))
		if err != nil {
			return fmt.Errorf("could not scan Durations value, %v", err)
		}
		durs = append(durs, d)
	}
	newDS := Durations(durs)
	*ds = newDS
	return nil
}

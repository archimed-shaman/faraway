package config

import (
	"fmt"
	"time"
)

type sDuration time.Duration

func (sd sDuration) String() string { return time.Duration(sd).String() }

func (sd *sDuration) UnmarshalYAML(unmarshal func(any) error) error {
	var data string

	if err := unmarshal(&data); err != nil {
		return err
	}

	if data != "" && data[0] == '"' {
		data = data[1:]
	}

	if data != "" && data[len(data)-1] == '"' {
		data = data[:len(data)-1]
	}

	t, err := time.ParseDuration(data)
	if err != nil {
		return fmt.Errorf("failed to parse '%s' to time.Duration: %w", data, err)
	}

	*sd = sDuration(t)

	return nil
}

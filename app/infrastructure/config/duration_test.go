package config

import (
	"faraway/wow/pkg/test"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestSDuration_String(t *testing.T) {
	t.Parallel()

	duration := sDuration(time.Hour)
	expected := "1h0m0s"

	test.Check(t, "String method", expected, duration.String())
}

func TestSDuration_UnmarshalYAML_ValidDurationWithoutQuotes(t *testing.T) {
	t.Parallel()

	var duration sDuration

	yamlInput := "1h30m"
	expected := 1*time.Hour + 30*time.Minute

	err := yaml.Unmarshal([]byte(yamlInput), &duration)
	test.Nil(t, "UnmarshalYAML error", err)
	test.Check(t, "Unmarshaled duration", sDuration(expected), duration)
}

func TestSDuration_UnmarshalYAML_ValidDurationWithQuotes(t *testing.T) {
	t.Parallel()

	var duration sDuration

	yamlInput := `"1h30m"`
	expected := 1*time.Hour + 30*time.Minute

	err := yaml.Unmarshal([]byte(yamlInput), &duration)
	test.Nil(t, "UnmarshalYAML error", err)
	test.Check(t, "Unmarshaled duration", sDuration(expected), duration)
}

func TestSDuration_UnmarshalYAML_InvalidDuration(t *testing.T) {
	t.Parallel()

	var duration sDuration

	yamlInput := "invalid"

	err := yaml.Unmarshal([]byte(yamlInput), &duration)
	test.NotNil(t, "UnmarshalYAML error", err)
}

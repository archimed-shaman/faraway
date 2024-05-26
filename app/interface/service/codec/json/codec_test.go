package json

import (
	"encoding/json"
	"faraway/wow/pkg/test"
	"testing"
	"time"
)

func TestCodec_Marshal_Unmarshal(t *testing.T) {
	t.Parallel()

	// Test struct for encoding/decoding
	type InnerStruct struct {
		InnerField string `json:"innerField"`
	}

	type TestStruct struct {
		StringField    string         `json:"stringField"`
		IntField       int            `json:"intField"`
		FloatField     float64        `json:"floatField"`
		BoolField      bool           `json:"boolField"`
		TimeField      time.Time      `json:"timeField"`
		StructField    InnerStruct    `json:"structField"`
		SliceField     []string       `json:"sliceField"`
		MapField       map[string]int `json:"mapField"`
		InterfaceField interface{}    `json:"interfaceField"`
	}

	original := &TestStruct{
		StringField:    "test",
		IntField:       42,
		FloatField:     3.14,
		BoolField:      true,
		TimeField:      time.Now(),
		StructField:    InnerStruct{InnerField: "inner"},
		SliceField:     []string{"slice1", "slice2"},
		MapField:       map[string]int{"key1": 1, "key2": 2},
		InterfaceField: map[string]any{"nestedKey": "nestedValue"},
	}
	codec := NewCodec()

	// Marshal using goccy/go-json
	goccyData, err := codec.Marshal(original)
	test.Nil(t, "Marshal error with goccy/go-json", err)

	// Unmarshal using goccy/go-json
	var goccyUnmarshaled TestStruct
	err = codec.Unmarshal(goccyData, &goccyUnmarshaled)
	test.Nil(t, "Unmarshal error with goccy/go-json", err)

	// Ensure goccy/go-json marshaled and unmarshaled data is correct
	test.Check(t, "Unmarshaled struct with goccy/go-json", *original, goccyUnmarshaled)

	// Marshal using standard library encoding/json
	stdData, err := json.Marshal(original)
	test.Nil(t, "Marshal error with encoding/json", err)

	// Unmarshal using standard library encoding/json
	var stdUnmarshaled TestStruct
	err = json.Unmarshal(stdData, &stdUnmarshaled)
	test.Nil(t, "Unmarshal error with encoding/json", err)

	// Ensure encoding/json marshaled and unmarshaled data is correct
	test.Check(t, "Unmarshaled struct with encoding/json", *original, stdUnmarshaled)

	// Compare goccy/go-json data with encoding/json data
	test.Check(t, "goccy/go-json vs encoding/json data", goccyData, stdData)
}

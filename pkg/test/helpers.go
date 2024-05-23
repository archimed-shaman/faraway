package test

import (
	"errors"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-cmp/cmp"
)

func Nil(t *testing.T, what string, got any) {
	t.Helper()

	if !cmp.Equal(nil, got) {
		t.Errorf("%s expected: [nil], got: [%v]", what, spew.Sdump(got))
	}
}

func Err(t *testing.T, what string, expected, got error) {
	t.Helper()

	if !errors.Is(got, expected) {
		t.Errorf("%s expected: [%v], got: [%v]", what, expected, got)
	}
}

func Check(t *testing.T, what string, expected, got any) {
	t.Helper()

	if !cmp.Equal(expected, got) {
		t.Errorf("%s expected: [%v], got: [%v]", what, spew.Sdump(expected), spew.Sdump(got))
	}
}

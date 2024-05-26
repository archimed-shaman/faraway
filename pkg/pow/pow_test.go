package pow

import (
	"context"
	"faraway/wow/pkg/test"
	"fmt"
	"testing"
)

func TestGenChallenge_Success(t *testing.T) {
	t.Parallel()

	challengeLen := 32
	zeroLowerBits := 16

	challenge, err := GenChallenge(challengeLen, zeroLowerBits)
	test.Nil(t, "GenChallenge error", err)

	if len(challenge) != challengeLen {
		t.Fatalf("GenChallenge returned challenge of length %d, expected %d", len(challenge), challengeLen)
	}

	hash := Hash(challenge)
	expected := false

	got, err := CheckLowerBitsZero(hash, zeroLowerBits)
	test.Nil(t, "CheckLowerBitsZero error", err)

	test.Check(t, fmt.Sprintf("Generated challenge hasn't %d lower bits as zero", zeroLowerBits), expected, got)
}

func TestGenChallenge_InvalidChallengeLen(t *testing.T) {
	t.Parallel()

	_, err := GenChallenge(0, 16)
	test.Err(t, "Invalid challenge length", ErrUnableGenerate, err)
}

func TestGenChallenge_InvalidZeroLowerBits(t *testing.T) {
	t.Parallel()

	_, err := GenChallenge(32, 0)
	test.Err(t, "Invalid zero lower bits", ErrUnableGenerate, err)
}

func TestGenChallenge_ZeroLowerBitsExceedsChallengeLen(t *testing.T) {
	t.Parallel()

	_, err := GenChallenge(1, 16)
	test.Err(t, "Zero lower bits exceeds challenge length", ErrUnableGenerate, err)
}

//nolint:paralleltest //Test must be ran without parallel as it replaces global random generator
func TestGenChallenge_MaxAttemptsExceeded(t *testing.T) {
	originalRandBytes := randBytes
	defer func() { randBytes = originalRandBytes }()

	randBytes = func(b []byte) (int, error) {
		// Generating a challenge that will always fail CheckLowerBitsZero
		// 368 or [1 112] gives the necessary hash
		b[0] = 1
		b[1] = 112

		return len(b), nil
	}

	_, err := GenChallenge(2, 4)
	test.Err(t, "Max attempts exceeded", ErrUnableGenerate, err)
}

func TestCheckLowerBitsZero_AllBitsZero(t *testing.T) {
	t.Parallel()

	data := []byte{0x00}
	n := 8
	expected := true

	got, err := CheckLowerBitsZero(data, n)
	test.Nil(t, "CheckLowerBitsZero error", err)

	test.Check(t, "All bits are zero", expected, got)
}

func TestCheckLowerBitsZero_OneBitSet(t *testing.T) {
	t.Parallel()

	data := []byte{0x01}
	n := 8
	expected := false

	got, err := CheckLowerBitsZero(data, n)
	test.Nil(t, "CheckLowerBitsZero error", err)

	test.Check(t, "One bit is set", expected, got)
}

func TestCheckLowerBitsZero_TwoBytesZero(t *testing.T) {
	t.Parallel()

	data := []byte{0x00, 0x00}
	n := 16
	expected := true

	got, err := CheckLowerBitsZero(data, n)
	test.Nil(t, "CheckLowerBitsZero error", err)

	test.Check(t, "Two bytes are zero", expected, got)
}

func TestCheckLowerBitsZero_LeastSignificantByteNotZero(t *testing.T) {
	t.Parallel()

	data := []byte{0x00, 0x01}
	n := 16
	expected := false

	got, err := CheckLowerBitsZero(data, n)
	test.Nil(t, "CheckLowerBitsZero error", err)

	test.Check(t, "Least significant byte is not zero", expected, got)
}

func TestCheckLowerBitsZero_ThreeBytesZero(t *testing.T) {
	t.Parallel()

	data := []byte{0x00, 0x00, 0x00}
	n := 24
	expected := true

	got, err := CheckLowerBitsZero(data, n)
	test.Nil(t, "CheckLowerBitsZero error", err)

	test.Check(t, "Three bytes are zero", expected, got)
}

func TestCheckLowerBitsZero_LeastSignificantByteNotZero_ThreeBytes(t *testing.T) {
	t.Parallel()

	data := []byte{0x00, 0x00, 0x01}
	n := 24
	expected := false

	got, err := CheckLowerBitsZero(data, n)
	test.Nil(t, "CheckLowerBitsZero error", err)

	test.Check(t, "Three bytes, least significant byte is not zero", expected, got)
}

func TestCheckLowerBitsZero_LeastSignificant4BitsZero(t *testing.T) {
	t.Parallel()

	data := []byte{0x00, 0x10}
	n := 4
	expected := true

	got, err := CheckLowerBitsZero(data, n)
	test.Nil(t, "CheckLowerBitsZero error", err)

	test.Check(t, "Least significant 4 bits are zero", expected, got)
}

func TestCheckLowerBitsZero_LeastSignificantBitSet(t *testing.T) {
	t.Parallel()

	data := []byte{0x00, 0x01}
	n := 4
	expected := false

	got, err := CheckLowerBitsZero(data, n)
	test.Nil(t, "CheckLowerBitsZero error", err)

	test.Check(t, "Least significant bit is set", expected, got)
}

func TestCheckLowerBitsZero_Last4BitsZero(t *testing.T) {
	t.Parallel()

	data := []byte{0xFF, 0x00}
	n := 4
	expected := true

	got, err := CheckLowerBitsZero(data, n)
	test.Nil(t, "CheckLowerBitsZero error", err)

	test.Check(t, "Last 4 bits are zero", expected, got)
}

func TestCheckLowerBitsZero_Last4BitsNotZero(t *testing.T) {
	t.Parallel()

	data := []byte{0xFF, 0x01}
	n := 4
	expected := false

	got, err := CheckLowerBitsZero(data, n)
	test.Nil(t, "CheckLowerBitsZero error", err)

	test.Check(t, "Last 4 bits are not zero", expected, got)
}

func TestCheckLowerBitsZero_NExceedsSizeOfData(t *testing.T) {
	t.Parallel()

	data := []byte{0x00}
	n := 16
	expected := false

	got, err := CheckLowerBitsZero(data, n)
	test.Err(t, "CheckLowerBitsZero error", ErrInvalidBitCheck, err)

	test.Check(t, "N exceeds size of data", expected, got)
}

func TestCheckLowerBitsZero_NExceedsSizeOfData_TwoBytes(t *testing.T) {
	t.Parallel()

	data := []byte{0x00, 0x00}
	n := 24
	expected := false

	got, err := CheckLowerBitsZero(data, n)
	test.Err(t, "CheckLowerBitsZero error", ErrInvalidBitCheck, err)

	test.Check(t, "N exceeds size of data (two bytes)", expected, got)
}

func TestCheckLowerBitsZero_LongDataAllZero(t *testing.T) {
	t.Parallel()

	data := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	n := 160
	expected := true

	got, err := CheckLowerBitsZero(data, n)
	test.Nil(t, "CheckLowerBitsZero error", err)

	test.Check(t, "Long data all zero", expected, got)
}

func TestCheckLowerBitsZero_LongDataLastByteNotZero(t *testing.T) {
	t.Parallel()

	data := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
	}
	n := 160
	expected := false

	got, err := CheckLowerBitsZero(data, n)
	test.Nil(t, "CheckLowerBitsZero error", err)

	test.Check(t, "Long data, last byte not zero", expected, got)
}

func TestCheckLowerBitsZero_LongDataLastByteZero(t *testing.T) {
	t.Parallel()

	data := []byte{
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00,
	}
	n := 8
	expected := true

	got, err := CheckLowerBitsZero(data, n)
	test.Nil(t, "CheckLowerBitsZero error", err)

	test.Check(t, "Long data, last byte zero", expected, got)
}

func TestCheckLowerBitsZero_LongDataFewZeros(t *testing.T) {
	t.Parallel()

	data := []byte{
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xE0, 0x00, 0x00, 0x00, 0x00,
	}
	n := 37
	expected := true

	got, err := CheckLowerBitsZero(data, n)
	test.Nil(t, "CheckLowerBitsZero error", err)

	test.Check(t, "Long data, few zeros", expected, got)
}

func TestResolve_Success(t *testing.T) {
	t.Parallel()

	challenge := []byte{0x05, 0x04, 0x03, 0x02, 0x01, 0x00}

	expected := []byte{0x01, 0xf9} // Precalculated solution

	got, err := Resolve(context.Background(), challenge, 8)
	test.Nil(t, "Resolve error", err)

	test.Check(t, "Solution for 0x050403020100", expected, got)
}

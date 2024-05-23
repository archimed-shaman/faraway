package pow

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"math/big"

	pkgerr "github.com/pkg/errors"
)

const byteSize = 8

var (
	ErrUnableGenerate  = errors.New("unable generate appropriate challenge")
	ErrInvalidBitCheck = errors.New("n is greater than data contains")
)

var randBytes = rand.Read

// GenChallenge generates a random challenge of given length and ensures
// that its hash does not end with a specified number of zero lower bits.
func GenChallenge(challengeLen, zeroLowerBits int) ([]byte, error) {
	const maxAttempts = 100

	if challengeLen < 1 || zeroLowerBits <= 0 || challengeLen*byteSize < zeroLowerBits {
		return nil, ErrUnableGenerate
	}

	b := make([]byte, challengeLen)

	for i := 0; i < maxAttempts; i++ {
		_, err := randBytes(b)
		if err != nil {
			return nil, pkgerr.Wrap(err, "failed to read random bytes")
		}

		// Check the hash we've generated does not contain the necessary lower zero bits
		hash := Hash(b)

		ok, err := CheckLowerBitsZero(hash, zeroLowerBits)
		if err != nil {
			return nil, err
		}

		if !ok {
			return b, nil
		}
	}

	return nil, ErrUnableGenerate
}

// CheckSolution verifies if the solution to the challenge is correct by checking the lower bits.
func CheckSolution(challenge, solution []byte, zeroLowerBits int) (bool, error) {
	result := make([]byte, len(challenge)+len(solution))

	copy(result, challenge)
	copy(result[len(challenge):], solution)

	hash := Hash(result)

	return CheckLowerBitsZero(hash, zeroLowerBits)
}

// Resolve tries to find a solution for the given challenge by brute-forcing until the lower bits are zero.
func Resolve(challenge []byte, zeroLowerBits int) ([]byte, error) {
	const maxInt64 = int64(^uint64(0) >> 1)

	for i := int64(0); i < maxInt64; i++ {
		solution := big.NewInt(i).Bytes()

		ok, err := CheckSolution(challenge, solution, zeroLowerBits)
		if err != nil {
			return nil, err
		}

		if ok {
			return solution, nil
		}
	}

	return nil, ErrUnableGenerate
}

// Hash computes the SHA-256 hash of the given data.
func Hash(data []byte) []byte {
	h := sha256.New()
	h.Write(data)

	return h.Sum(nil)
}

// CheckLowerBitsZero checks if the last n lower bits of the given bytes are zero.
func CheckLowerBitsZero(data []byte, n int) (bool, error) {
	if n > len(data)*byteSize {
		return false, ErrInvalidBitCheck
	}

	bytes := n / byteSize
	bits := n % byteSize

	for i := 1; i <= bytes; i++ {
		if data[len(data)-i] != 0 {
			return false, nil
		}
	}

	if bits > 0 {
		lastByte := data[len(data)-bytes-1]
		mask := byte((1 << bits) - 1)

		if lastByte&mask != 0 {
			return false, nil
		}
	}

	return true, nil
}

// package pow

// import (
// 	"crypto/rand"
// 	"crypto/sha256"
// 	"errors"
// 	"math/big"

// 	pkgerr "github.com/pkg/errors"
// )

// const byteSize = 8

// var ErrUnableGenerate = errors.New("unable generate appropriate challenge")

// var randBytes = rand.Read

// // GenChallenge generates a random challenge of given length and ensures
// // that its hash does not end with a specified number of zero lower bits.
// func GenChallenge(challengeLen, zeroLowerBits int) ([]byte, error) {
// 	const maxAttempts = 100

// 	if challengeLen < 1 {
// 		return nil, ErrUnableGenerate
// 	}

// 	if zeroLowerBits <= 0 {
// 		return nil, ErrUnableGenerate
// 	}

// 	if challengeLen*byteSize < zeroLowerBits {
// 		return nil, ErrUnableGenerate
// 	}

// 	b := make([]byte, challengeLen)

// 	for i := 0; i < maxAttempts; i++ {
// 		_, err := randBytes(b)
// 		if err != nil {
// 			return nil, pkgerr.Wrap(err, "failed to read random bytes")
// 		}

// 		// Check the hash we've generated does not contain the necessary lower zero bits
// 		hash := Hash(b)

// 		if !CheckLowerBitsZero(hash, zeroLowerBits) {
// 			return b, nil
// 		}
// 	}

// 	return nil, ErrUnableGenerate
// }

// func CheckSolution(challenge, solution []byte, zeroLowerBits int) bool {
// 	result := make([]byte, len(challenge)+len(solution))

// 	copy(result, challenge)
// 	copy(result[len(challenge):], solution)

// 	hash := Hash(result)

// 	return CheckLowerBitsZero(hash, zeroLowerBits)
// }

// func Resolve(challenge []byte, zeroLowerBits int) []byte {
// 	const maxInt64 = int64(^uint64(0) >> 1)

// 	for i := int64(0); i < maxInt64; i++ {
// 		solution := big.NewInt(i).Bytes()

// 		if CheckSolution(challenge, solution, zeroLowerBits) {
// 			return solution
// 		}
// 	}

// 	return nil
// }

// // Hash computes the SHA-256 hash of the given data.
// func Hash(data []byte) []byte {
// 	h := sha256.New()
// 	h.Write(data)

// 	return h.Sum(nil)
// }

// // CheckLowerBitsZero checks if the last n lower bits of the given bytes are zero.
// func CheckLowerBitsZero(data []byte, n int) bool {
// 	if n > len(data)*byteSize {
// 		panic("n is greater than data contains")
// 	}

// 	bytes := n / byteSize
// 	bits := n % byteSize

// 	for i := 1; i <= bytes; i++ {
// 		if data[len(data)-i] != 0 {
// 			return false
// 		}
// 	}

// 	if bits > 0 {
// 		lastByte := data[len(data)-bytes-1]
// 		mask := byte((1 << bits) - 1)

// 		if lastByte&mask != 0 {
// 			return false
// 		}
// 	}

// 	return true
// }

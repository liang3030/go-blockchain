package main

import (
	"fmt"
	"math/big"
)

const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// Base58 alphabet, excluding visually similar characters.
var base = big.NewInt(int64(len(alphabet)))

var alphabetMap map[byte]int

func init() {
	alphabetMap = make(map[byte]int)
	for i, char := range []byte(alphabet) {
		alphabetMap[char] = i
	}
}

func Base58Encode(input []byte) []byte {
	if len(input) == 0 {
		return nil
	}

	x := new(big.Int).SetBytes(input)
	mod := new(big.Int)
	var result []byte

	// Convert to base58
	for x.Sign() > 0 {
		x.DivMod(x, base, mod)
		result = append(result, alphabet[mod.Int64()])
	}

	// Add leading zero bytes
	for _, b := range input {
		if b == 0x00 {
			result = append(result, alphabet[0])
		} else {
			break
		}
	}

	// Since we construct the encoded string in reverse order, we need to reverse it before returning
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result
}

func Base58Decode(input []byte) ([]byte, error) {
	x := big.NewInt(0)
	for _, b := range input {
		value, exists := alphabetMap[b]
		if !exists {
			return nil, fmt.Errorf("invalid character '%c' in input", b)
		}
		x.Mul(x, base)
		x.Add(x, big.NewInt(int64(value)))
	}
	result := x.Bytes()
	zeroes := 0
	for zeroes < len(input) && input[zeroes] == alphabet[0] {
		zeroes++
	}
	return append(make([]byte, zeroes), result...), nil
}

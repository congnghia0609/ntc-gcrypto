/**
 *
 * @author nghiatc
 * @since Jan 2, 2020
 */
package sss

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"strings"
)

// The largest PRIME 256-bit big.Int
// https://primes.utm.edu/lists/2small/200bit.html
// PRIME = 2^n - k = 2^256 - 189
const (
	DefaultPrimeStr = "115792089237316195423570985008687907853269984665640564039457584007913129639747"
)

// var PRIME *big.Int
var PRIME, _ = big.NewInt(0).SetString(DefaultPrimeStr, 10)

// Returns a random number from the range (0, PRIME-1) inclusive
func random() *big.Int {
	result := big.NewInt(0).Set(PRIME)
	result = result.Sub(result, big.NewInt(1))
	result, _ = rand.Int(rand.Reader, result)
	return result
}

// Converts a byte array into an a 256-bit big.Int, array based upon size of
// the input byte; all values are right-padded to length 256, even if the most
// significant bit is zero.
func splitByteToInt(secret []byte) []*big.Int {
	hex_data := hex.EncodeToString(secret)
	count := int(math.Ceil(float64(len(hex_data)) / 64.0))

	result := make([]*big.Int, count)
	for i := 0; i < count; i++ {
		if (i+1)*64 < len(hex_data) {
			result[i], _ = big.NewInt(0).SetString(hex_data[i*64:(i+1)*64], 16)
		} else {
			data := strings.Join([]string{hex_data[i*64:], strings.Repeat("0", 64-(len(hex_data)-i*64))}, "")
			result[i], _ = big.NewInt(0).SetString(data, 16)
		}
	}

	return result
}

// Converts an array of big.Ints to the original byte array, removing any
// least significant nulls
func mergeIntToByte(secret []*big.Int) []byte {
	var hex_data = ""
	for i := range secret {
		tmp := fmt.Sprintf("%x", secret[i])
		hex_data += strings.Join([]string{strings.Repeat("0", (64 - len(tmp))), tmp}, "")
	}
	result, _ := hex.DecodeString(hex_data)
	result = bytes.TrimRight(result, "\x00")

	return result
}

// Evaluates a polynomial with coefficients specified in reverse order:
// evaluatePolynomial([a, b, c, d], x):
//
//	return a + bx + cx^2 + dx^3
//
// Horner's method: ((dx + c)x + b)x + a
func evaluatePolynomial(polynomial []*big.Int, value *big.Int) *big.Int {
	last := len(polynomial) - 1
	var result *big.Int = big.NewInt(0).Set(polynomial[last])

	for s := last - 1; s >= 0; s-- {
		result = result.Mul(result, value)
		result = result.Add(result, polynomial[s])
		result = result.Mod(result, PRIME)
	}

	return result
}

// inNumbers(array, value) returns boolean whether or not value is in array
func inNumbers(numbers []*big.Int, value *big.Int) bool {
	for n := range numbers {
		if numbers[n].Cmp(value) == 0 {
			return true
		}
	}

	return false
}

// Returns the big.Int number base10 in base64 representation; note: this is
// not a string representation; the base64 output is exactly 256 bits long
func toBase64(number *big.Int) string {
	hexdata := fmt.Sprintf("%x", number)
	for i := 0; len(hexdata) < 64; i++ {
		hexdata = "0" + hexdata
	}
	bytedata, err := hex.DecodeString(hexdata)
	if err != nil {
		fmt.Println(err)
	}
	return base64.URLEncoding.EncodeToString(bytedata)
}

// Returns the big.Int number base10 in Hex representation; note: this is
// not a string representation; the Hex output is exactly 256 bits long
func toHex(number *big.Int) string {
	hexdata := fmt.Sprintf("%x", number)
	for i := 0; len(hexdata) < 64; i++ {
		hexdata = "0" + hexdata
	}
	return hexdata
}

// Returns the number base64 in base 10 big.Int representation; note: this is
// not coming from a string representation; the base64 input is exactly 256
// bits long, and the output is an arbitrary size base 10 integer.
//
// Returns -1 on failure
func fromBase64(number string) *big.Int {
	bytedata, err := base64.URLEncoding.DecodeString(number)
	if err != nil {
		return big.NewInt(-1)
	}
	hexdata := hex.EncodeToString(bytedata)
	result, ok := big.NewInt(0).SetString(hexdata, 16)
	if !ok {
		return big.NewInt(-1)
	}

	return result
}

// Returns the number Hex in base 10 big.Int representation; note: this is
// not coming from a string representation; the Hex input is exactly 256
// bits long, and the output is an arbitrary size base 10 integer.
//
// Returns -1 on failure
func fromHex(number string) *big.Int {
	result, ok := big.NewInt(0).SetString(number, 16)
	if !ok {
		return big.NewInt(-1)
	}

	return result
}

// Computes the multiplicative inverse of the number on the field PRIME; more
// specifically, number * inverse == 1; Note: number should never be zero
func modInverse(number *big.Int) *big.Int {
	copy := big.NewInt(0).Set(number)
	copy = copy.Mod(copy, PRIME)
	pcopy := big.NewInt(0).Set(PRIME)
	x := big.NewInt(0)
	y := big.NewInt(0)

	copy.GCD(x, y, pcopy, copy)

	result := big.NewInt(0).Set(PRIME)
	result = result.Add(result, y)
	result = result.Mod(result, PRIME)

	return result
}

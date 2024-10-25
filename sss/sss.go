/**
 *
 * @author nghiatc
 * @since Jan 2, 2020
 */
package sss

import (
	"errors"
	"math/big"
)

var (
	errCannotRequireMoreShares = errors.New("cannot require more shares then existing")
	errOneOfTheSharesIsInvalid = errors.New("one of the shares is invalid")
)

// Returns a new array of secret shares (encoding x,y pairs as Base64 or Hex strings)
// created by Shamir's Secret Sharing Algorithm requiring a minimum number of
// share to recreate, of length shares, from the input secret raw as a string
func Create(minimum int, shares int, secret string, isBase64 bool) ([]string, error) {
	// Verify minimum isn't greater than shares; there is no way to recreate
	// the original polynomial in our current setup, therefore it doesn't make
	// sense to generate fewer shares than are needed to reconstruct the secrets.
	if minimum <= 0 || shares <= 0 {
		return []string{""}, errors.New("minimum or shares is invalid")
	}
	if minimum > shares {
		return []string{""}, errCannotRequireMoreShares
	}
	if len(secret) == 0 {
		return []string{""}, errors.New("secret is empty")
	}

	// Convert the secrets to its respective 256-bit big.Int representation
	var secrets []*big.Int = splitByteToInt([]byte(secret))

	// List of currently used numbers in the polynomial
	var numbers []*big.Int = make([]*big.Int, 0)
	numbers = append(numbers, big.NewInt(0))

	// Create the polynomial of degree (minimum - 1); that is, the highest
	// order term is (minimum-1), though as there is a constant term with
	// order 0, there are (minimum) number of coefficients.
	//
	// However, the polynomial object is a 2d array, because we are constructing
	// a different polynomial for each part of the secrets.
	//
	// polynomial[parts][minimum]
	var polynomial [][]*big.Int = make([][]*big.Int, len(secrets))
	for i := range polynomial {
		polynomial[i] = make([]*big.Int, minimum)
		polynomial[i][0] = secrets[i]

		for j := range polynomial[i][1:] {
			//fmt.Println("j =", j)
			// Each coefficient should be unique
			number := random()
			for inNumbers(numbers, number) {
				number = random()
			}
			numbers = append(numbers, number)

			polynomial[i][j+1] = number
		}
	}

	// Create the points object; this holds the (x, y) points of each share.
	// Again, because secrets is an array, each share could have multiple parts
	// over which we are computing Shamir's Algorithm. The last dimension is
	// always two, as it is storing an x, y pair of points.
	//
	// Note: this array is technically unnecessary due to creating result
	// in the inner loop. Can disappear later if desired.
	var result []string = make([]string, shares)

	// For every share...
	len := len(secrets)
	for i := 0; i < shares; i++ {
		// and every part of the secrets...
		for j := 0; j < len; j++ {
			// generate a new x-coordinate.
			x := random()
			for inNumbers(numbers, x) {
				x = random()
			}
			numbers = append(numbers, x)

			// and evaluate the polynomial at that point.
			y := evaluatePolynomial(polynomial[j], x)

			// add it to results.
			if isBase64 {
				result[i] += toBase64(x)
				result[i] += toBase64(y)
			} else {
				result[i] += toHex(x)
				result[i] += toHex(y)
			}
		}
	}

	return result, nil
}

// Takes a string array of shares encoded in Base64 or Hex created via Shamir's Algorithm;
//
// Note: the polynomial will converge if the specified minimum number of shares
//
//	or more are passed to this function. Passing thus does not affect it
//	Passing fewer however, simply means that the returned secret is wrong.
func Combine(shares []string, isBase64 bool) (string, error) {
	if len(shares) == 0 {
		return "", errors.New("shares is empty")
	}

	// Recreate the original object of x, y points, based upon number of shares
	// and size of each share (number of parts in the secrets).
	//
	// points[shares][parts][2]
	var points [][][]*big.Int
	if isBase64 {
		points, _ = decodeShareBase64(shares)
	} else {
		points, _ = decodeShareHex(shares)
	}

	// Use Lagrange Polynomial Interpolation (LPI) to reconstruct the secrets.
	// For each part of the secrets (clearest to iterate over)...
	var secrets []*big.Int = make([]*big.Int, len(points[0]))
	for j := range secrets {
		secrets[j] = big.NewInt(0)
		// and every share...
		for i := range points { // LPI sum loop
			// remember the current x and y values.
			ax := points[i][j][0]        // ax
			ay := points[i][j][1]        // ay
			numerator := big.NewInt(1)   // LPI numerator
			denominator := big.NewInt(1) // LPI denominator
			// and for every other point...
			for k := range points { // LPI product loop
				if k != i {
					// combine them via half products.
					// x=0 ==> [(0-bx)/(ax-bx)] * ...
					bx := points[k][j][0] // bx
					negbx := big.NewInt(0)
					negbx = negbx.Mul(bx, big.NewInt(-1)) // (0-bx)
					axbx := big.NewInt(0)
					axbx = axbx.Sub(ax, bx) // (ax-bx)

					numerator = numerator.Mul(numerator, negbx) // (0-bx)*...
					numerator = numerator.Mod(numerator, PRIME)

					denominator = denominator.Mul(denominator, axbx) // (ax-bx)*...
					denominator = denominator.Mod(denominator, PRIME)
				}
			}

			// LPI product: x=0, y = ay * [(x-bx)/(ax-bx)] * ...
			// multiply together the points (ay)(numerator)(denominator)^-1 ...
			fx := big.NewInt(0).Set(ay)
			fx = fx.Mul(fx, numerator)
			fx = fx.Mul(fx, modInverse(denominator))

			// LPI sum: s = fx + fx + ...
			secrets[j] = secrets[j].Add(secrets[j], fx)
			secrets[j] = secrets[j].Mod(secrets[j], PRIME)
		}
	}

	// recover secret string.
	return string(mergeIntToByte(secrets)), nil
}

// Takes a string array of shares encoded in Base64 created via Shamir's
// Algorithm; each string must be of equal length of a multiple of 88 characters
// as a single 88 character share is a pair of 256-bit numbers (x, y).
func decodeShareBase64(shares []string) ([][][]*big.Int, error) {
	// Recreate the original object of x, y points, based upon number of shares
	// and size of each share (number of parts in the secret).
	var secrets [][][]*big.Int = make([][][]*big.Int, len(shares))

	// For each share...
	for i := range shares {
		// ensure that it is valid.
		if !isValidShareBase64(shares[i]) {
			return nil, errOneOfTheSharesIsInvalid
		}

		// find the number of parts it represents.
		share := shares[i]
		count := len(share) / 88
		secrets[i] = make([][]*big.Int, count)

		// and for each part, find the x,y pair...
		for j := range secrets[i] {
			cshare := share[j*88 : (j+1)*88]
			secrets[i][j] = make([]*big.Int, 2)
			// decoding from Base64.
			secrets[i][j][0] = fromBase64(cshare[0:44])
			secrets[i][j][1] = fromBase64(cshare[44:])
		}
	}

	return secrets, nil
}

// Takes a string array of shares encoded in Hex created via Shamir's
// Algorithm; each string must be of equal length of a multiple of 128 characters
// as a single 128 character share is a pair of 256-bit numbers (x, y).
func decodeShareHex(shares []string) ([][][]*big.Int, error) {
	// Recreate the original object of x, y points, based upon number of shares
	// and size of each share (number of parts in the secret).
	var secrets [][][]*big.Int = make([][][]*big.Int, len(shares))

	// For each share...
	for i := range shares {
		// ensure that it is valid.
		if !isValidShareHex(shares[i]) {
			return nil, errOneOfTheSharesIsInvalid
		}

		// find the number of parts it represents.
		share := shares[i]
		count := len(share) / 128
		secrets[i] = make([][]*big.Int, count)

		// and for each part, find the x,y pair...
		for j := range secrets[i] {
			cshare := share[j*128 : (j+1)*128]
			secrets[i][j] = make([]*big.Int, 2)
			// decoding from Hex.
			secrets[i][j][0] = fromHex(cshare[0:64])
			secrets[i][j][1] = fromHex(cshare[64:])
		}
	}

	return secrets, nil
}

// Takes in a given string to check if it is a valid secret
//
// Requirements:
//
//	Length multiple of 88
//	Can decode each 44 character block as Base64
//
// Returns only success/failure (bool)
func isValidShareBase64(candidate string) bool {
	if len(candidate) == 0 || len(candidate)%88 != 0 {
		return false
	}
	count := len(candidate) / 44
	for j := 0; j < count; j++ {
		part := candidate[j*44 : (j+1)*44]
		decode := fromBase64(part)
		if decode.Cmp(big.NewInt(0)) <= 0 || decode.Cmp(PRIME) >= 0 {
			return false
		}
	}

	return true
}

// Takes in a given string to check if it is a valid secret
//
// Requirements:
//
//	Length multiple of 128
//	Can decode each 64 character block as Hex
//
// Returns only success/failure (bool)
func isValidShareHex(candidate string) bool {
	if len(candidate) == 0 || len(candidate)%128 != 0 {
		return false
	}
	count := len(candidate) / 64
	for j := 0; j < count; j++ {
		part := candidate[j*64 : (j+1)*64]
		decode := fromHex(part)
		if decode.Cmp(big.NewInt(0)) <= 0 || decode.Cmp(PRIME) >= 0 {
			return false
		}
	}

	return true
}

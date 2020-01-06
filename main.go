/**
 *
 * @author nghiatc
 * @since Jan 2, 2020
 */
package main

import (
	"fmt"
	"ntc-gcrypto/sss"
)

/* https://github.com/SSSaaS/sssa-golang */
func main() {
	// creates a set of shares
	s := "nghiatcxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	arr, err := sss.Create2(3, 6, s)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(arr)
	//json_arr, _ := json.Marshal(arr)
	//fmt.Println(string(json_arr))
	fmt.Println("\nshares.size:", len(arr))
	for i:=0; i<len(arr); i++ {
		fmt.Printf("shares[%d]: %s\n", i, arr[i])
	}

	fmt.Println("\nsecret:", s)
	fmt.Println("secret.length:", len(s))
	// combines shares into secret
	s1, err := sss.Combine2(arr[:3])
	fmt.Println("combines shares 1 length =", len(arr[:3]))
	fmt.Println("secret:", s1)
	fmt.Println("secret.length:", len(s1))

	s2, err := sss.Combine2(arr[3:])
	fmt.Println("combines shares 2 length =", len(arr[3:]))
	fmt.Println("secret:", s2)
	fmt.Println("secret.length:", len(s2))

	s3, err := sss.Combine2(arr[1:5])
	fmt.Println("combines shares 3 length =", len(arr[1:5]))
	fmt.Println("secret:", s3)
	fmt.Println("secret.length:", len(s3))
}

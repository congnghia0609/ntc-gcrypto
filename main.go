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
	fmt.Println("secret:", s)
	arr, err := sss.Create(3, 6, s)
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

	// combines shares into secret
	fmt.Println("\ncombines shares 1")
	s1, err := sss.Combine(arr[:3])
	fmt.Println(s1)

	fmt.Println("combines shares 2")
	s2, err := sss.Combine(arr[3:])
	fmt.Println(s2)
}

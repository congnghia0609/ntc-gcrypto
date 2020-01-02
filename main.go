package main

import (
	"encoding/json"
	"fmt"
	"ntc-gcrypto/sss"
)

func main() {
	// creates a set of shares
	s := "nghiatcxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	fmt.Println("secret:", s)
	arr, err := sss.Create(3, 6, s)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(arr)
	json_arr, _ := json.Marshal(arr)
	fmt.Println(string(json_arr))

	// combines shares into secret
	s1, err := sss.Combine(arr[:3])
	fmt.Println(s1)

	s2, err := sss.Combine(arr[3:])
	fmt.Println(s2)
}

# ntc-gcrypto
ntc-gcrypto is module golang cryptography  

## 1. An implementation of Shamir's Secret Sharing Algorithm 256-bits in golang

### Usage
**Use encode/decode Base64**  
```go
// creates a set of shares
s := "nghiatcxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
arr, err := sss.Create(3, 6, s, true)
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
s1, err := sss.Combine(arr[:3], true)
fmt.Println("combines shares 1 length =", len(arr[:3]))
fmt.Println("secret:", s1)
fmt.Println("secret.length:", len(s1))

s2, err := sss.Combine(arr[3:], true)
fmt.Println("combines shares 2 length =", len(arr[3:]))
fmt.Println("secret:", s2)
fmt.Println("secret.length:", len(s2))

s3, err := sss.Combine(arr[1:5], true)
fmt.Println("combines shares 3 length =", len(arr[1:5]))
fmt.Println("secret:", s3)
fmt.Println("secret.length:", len(s3))
```

**Use encode/decode Hex**  
```go
// creates a set of shares
s := "nghiatcxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
arr, err := sss.Create(3, 6, s, false)
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
s1, err := sss.Combine(arr[:3], false)
fmt.Println("combines shares 1 length =", len(arr[:3]))
fmt.Println("secret:", s1)
fmt.Println("secret.length:", len(s1))

s2, err := sss.Combine(arr[3:], false)
fmt.Println("combines shares 2 length =", len(arr[3:]))
fmt.Println("secret:", s2)
fmt.Println("secret.length:", len(s2))

s3, err := sss.Combine(arr[1:5], false)
fmt.Println("combines shares 3 length =", len(arr[1:5]))
fmt.Println("secret:", s3)
fmt.Println("secret.length:", len(s3))
```

## License
This code is under the [Apache Licence v2](https://www.apache.org/licenses/LICENSE-2.0).  

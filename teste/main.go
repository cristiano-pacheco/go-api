package main

import (
	"fmt"
	"time"

	"github.com/gbrlsnchs/jwt/v3"
)

type CustomPayload struct {
	jwt.Payload
	Foo string `json:"foo,omitempty"`
	Bar int    `json:"bar,omitempty"`
}

var hs = jwt.NewHS256([]byte("secret"))

func main() {
	now := time.Now()
	pl := CustomPayload{
		Payload: jwt.Payload{
			Issuer:         "gbrlsnchs",
			Subject:        "someone",
			Audience:       jwt.Audience{"https://golang.org", "https://jwt.io"},
			ExpirationTime: jwt.NumericDate(now.Add(24 * 30 * 12 * time.Hour)),
			NotBefore:      jwt.NumericDate(now.Add(30 * time.Minute)),
			IssuedAt:       jwt.NumericDate(now),
			JWTID:          "foobar",
		},
		Foo: "foo",
		Bar: 1337,
	}

	token, err := jwt.Sign(pl, hs)
	if err != nil {
		fmt.Printf("deu erro")
	}

	tokenString := string(token)
	token2 := []byte(tokenString)

	var pl2 = &CustomPayload{}

	hd, err := jwt.Verify(token2, hs, &pl2)
	if err != nil {
		fmt.Printf("deu erro validação")
		return
	}

	fmt.Println(hd)

	// ...
}

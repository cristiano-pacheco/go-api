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
			ExpirationTime: jwt.NumericDate(now.Add(2 * time.Second)),
			IssuedAt:       jwt.NumericDate(now),
		},
	}

	token, err := jwt.Sign(pl, hs)
	if err != nil {
		fmt.Printf("deu erro")
	}

	time.Sleep(3 * time.Second)

	now = time.Now()
	iatValidator := jwt.IssuedAtValidator(now)
	expValidator := jwt.ExpirationTimeValidator(now)

	var pl2 = &CustomPayload{}

	validatePayload := jwt.ValidatePayload(&pl2.Payload, iatValidator, expValidator)

	hd, err := jwt.Verify(token, hs, &pl2, validatePayload)
	if err != nil {
		fmt.Printf("deu erro validação")
		return
	}

	fmt.Println(hd)
}

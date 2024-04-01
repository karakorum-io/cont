package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {

	start := time.Now()
	data := [6]int{2, 3, 5, 7, 11, 13}

	ctx := context.WithValue(context.Background(), "some_key", data)
	userID := 10
	val, err := fetUserData(ctx, userID)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("result", val)
	fmt.Println("took", time.Since(start))
}

type Response struct {
	value int
	err   error
}

func fetUserData(ctx context.Context, userId int) (int, error) {

	val := ctx.Value("some_key")
	fmt.Println(val)

	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	respch := make(chan Response)

	go func() {
		val, err := fetchThirdPartyStuffWhichCanBeSlow()
		respch <- Response{
			value: val,
			err:   err,
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return 0, fmt.Errorf("Fetching data from third party took too long ...")
		case resp := <-respch:
			return resp.value, resp.err
		}
	}
}

func fetchThirdPartyStuffWhichCanBeSlow() (int, error) {
	time.Sleep(time.Millisecond * 150)

	return 666, nil
}

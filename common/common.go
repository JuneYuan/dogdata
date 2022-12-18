package common

import "fmt"

func CheckErr(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%v: %v", msg, err))
	}
}

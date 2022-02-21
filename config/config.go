package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	InitialBalanceAmount int `json:"initialBalanceAmount"`
	MinimumBalanceAmount int `json:"minimumBalanceAmount"`
}

var C = &Config{}

func init() {
	file, err := os.Open(".config/" + env + ".json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	read, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(read, C)
	if err != nil {
		panic(err)
	}

	fmt.Println(C.InitialBalanceAmount)
}

package main

import "fmt"

func main() {
	data := make(map[string]string)

	data["aa"] = "111"

	if data["bb"] == "" {
		fmt.Println(22222)
	}
	fmt.Println(11111)
}

package main

import (
	"fmt"

	"github.com/haikeloksama/gobank/util"
)

func main() {
	name, err := util.RandomOwner()
	
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Println("teh name is", name)
}
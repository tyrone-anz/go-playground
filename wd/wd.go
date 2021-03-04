package wd

import (
	"fmt"
	"os"
)

func Print() {
	dir, err := os.Getwd()
	fmt.Println("Dir: ", dir)
	fmt.Println(err)
}

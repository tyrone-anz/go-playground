package awesome

import "fmt"

func init() {
	fmt.Println("init was run when imported")
}

func SaySomething() {
	fmt.Println("something")
}

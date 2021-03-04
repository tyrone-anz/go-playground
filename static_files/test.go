package static_files

import "fmt"

func Dummy() {
	fmt.Println("dummy things")
}

func init() {
	fmt.Println("init was run when imported")
}

func SaySomething() {
	fmt.Println("something")
}

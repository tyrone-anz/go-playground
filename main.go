package main

import (
	"encoding/base64"
	"fmt"

	"github.com/tyrone-anz/go-playground/awesome"
)

func main() {
	awesome.SaySomething()

	var b []byte
	fmt.Println(string(b))

	// plan, _ := ioutil.ReadFile("testusers.json")
	// var data map[string]struct {
	// 	Username string `json:"username"`
	// 	Password string `json:"password"`
	// }
	//
	// if err := json.Unmarshal(plan, &data); err != nil {
	// 	panic(err)
	// }
	//
	// cnt := 1
	// for _, u := range data {
	// 	// fmt.Println(u.Username, u.Password, basicAuth(u.Username, u.Password))
	//
	// 	fmt.Println(fmt.Sprintf("#%d ", cnt))
	// 	fmt.Println("username: ", u.Username)
	// 	fmt.Println("password: ", u.Password)
	// 	fmt.Println("auth:     ", basicAuth(u.Username, u.Password))
	// 	fmt.Println()
	//
	// 	cnt++
	// }
}

func basicAuth(username, password string) string {
	fmt.Println(fmt.Sprintf("#%s ", username))
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

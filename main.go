package main

import (
    "fmt"
    "net/http"
    "reflect"
    "unsafe"

    "github.com/grpc-ecosystem/grpc-gateway/runtime"
    "github.com/grpc-ecosystem/grpc-gateway/utilities"
    "go.opentelemetry.io/otel"
)

func main() {
    // awesome.SaySomething()
    //
    // var b []byte
    // fmt.Println(string(b))
    //
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

    // var topics map[string]map[int64]int64
    //
    // // topics["abc"] = append(topics["abc"], append(topics["abc"][1], 1))
    // m, ok := topics["abc"]
    // if !ok {
    // 	topics["abc"] = s
    // }

    // Test channel close
//     ch1 := make(chan *string, 1)
//     close(ch1)
//     close(ch1)
//
//     // Test channel push, and receive on close
//     ch := make(chan *string, 1)
//     m1 := "hello1"
//     m2 := "hello2"
//     ch <- &m1
//
//     select {
//     case ch <- &m2:
//         fmt.Println("able to send")
//     default:
//         fmt.Println("failed to send")
//     }
//
//     close(ch)
//
// floop:
//     for {
//         select {
//         case m := <-ch:
//             fmt.Println("received ", m)
//             if m == nil {
//                 break floop
//             }
//         default:
//             fmt.Println("failed to receive")
//             break floop
//         }
//     }
//
//     var lock sync.Mutex
//
//     lock.Lock()
//     lock.Unlock()
//
//     pkgMap := map[string]int{}
//     pkgMap["hello"] = 1
//     if pkgMap["unknown"] != 2 {
//         fmt.Println("no panic here")
//     }

    // var testSlice []string
    //
    // newSlice := append(testSlice, "hello")
    // fmt.Println(newSlice)
    //
    // testSlice = nil
    // newSlice2 := append(testSlice, "hello")
    // fmt.Println(newSlice2)

    // sample := unexported.GetSample()
    //
    // v := reflect.ValueOf(*sample)
    // field := v.FieldByName("field")
    // fmt.Printf("%v\n", field)
    // fmt.Println(v.CanSet())
    //
    // // elem := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
    // // c, ok := elem.Interface().([]copyUnexportedStruct)
    // // elem := field.Elem()
    // if field.Kind() == reflect.Slice {
    //     for i := 0; i < field.Len(); i++ {
    //         item := field.Index(i)
    //         if item.Kind() == reflect.Struct {
    //             itemVal := reflect.Indirect(item)
    //             fmt.Printf("%v\n", itemVal)
    //             fmt.Printf("%v\n", itemVal.FieldByName("pattern"))
    //             fmt.Printf("%v\n", itemVal.FieldByName("fn"))
    //
    //             fmt.Println(itemVal.CanSet())
    //             fmt.Println(field.CanSet())
    //             fmt.Println(item.CanSet())
    //             fmt.Println(itemVal.CanSet())
    //             fmt.Println(itemVal.FieldByName("fn").CanSet())
    //
    //             // fn := reflect.ValueOf(printWorld)
    //
    //             // for j := 0; j < itemVal.NumField(); j++ {
    //             //     fmt.Println(itemVal.Type().Field(j).Name, itemVal.Field(j).Interface())
    //             // }
    //         }
    //     }
    // }
    //
    // sample.Run()
    // c, ok := elem.Interface().([]copyUnexportedStruct)
    // fmt.Println(ok)
    // fmt.Printf("%v", c)

    pat := runtime.Pattern{

    }
    fmt.Println(pat.String())


    mux := runtime.NewServeMux()
    mux.Handle("GET", mustPattern("help"), func(_ http.ResponseWriter, _ *http.Request, _ map[string]string) {})
    mux.Handle("GET", mustPattern("health"), func(_ http.ResponseWriter, _ *http.Request, _ map[string]string) {})
    mux.Handle("GET", mustPattern("accounts"), func(_ http.ResponseWriter, _ *http.Request, _ map[string]string) {})
    mux.Handle("POST", mustPattern("accounts"), func(_ http.ResponseWriter, _ *http.Request, _ map[string]string) {})

    muxV := reflect.ValueOf(*mux)
    hanldersMapV := muxV.FieldByName("handlers")
    if hanldersMapV.Kind() == reflect.Map {
        mapRange := hanldersMapV.MapRange()
        for  mapRange.Next() {
            method := mapRange.Key().String()
            fmt.Println(method)

            handlers := mapRange.Value()
            if handlers.Kind() == reflect.Slice {
                for i := 0; i < handlers.Len(); i++ {
                    h := handlers.Index(i)
                    if h.Kind() == reflect.Struct {
                        // patternV := h.FieldByName("pat")
                        // fmt.Println(patternV.Type())
                        // fmt.Println(patternV.Addr().Type())

                        // patternV := h.FieldByName("pat").Addr()
                        // patternElem := reflect.NewAt(patternV.Type(), unsafe.Pointer(patternV.UnsafeAddr())).Elem()
                        // fmt.Println(patternElem)

                        patternV := h.FieldByName("pat")
                        patternElem := reflect.NewAt(patternV.Type(), unsafe.Pointer(patternV.UnsafeAddr()))
                        pattern, ok := patternElem.Interface().(*runtime.Pattern)
                        if ok {
                            fmt.Println(pattern.String())
                        }

                        // fmt.Println(patternElem)
                        // fmt.Println(patternV.FieldByName("pool").)
                        // reflect.NewAt(patternV)
                        // fmt.Println(patternV.FieldByName("pool").)
                        // fmt.Println(patternV.Addr().CanInterface())
                        // method := patternV.MethodByName("String")
                        // fmt.Println(method.Call(nil))
                        // pat,  ok := patternV.Addr().Interface().(*runtime.Pattern)
                        // if ok {
                        //     fmt.Println(pat.String())
                        // }
                    }
                }
            }
        }
    }

    // dir, err := os.Getwd()
    // fmt.Println("Dir: ", dir)
    // fmt.Println(err)
    //
    // wd.Print()

    // fmt.Println(runtime.Caller(0))

    otel.Meter("").NewInt64Counter("http.client.duration")
    otel.Meter("").NewInt64Counter("http.server.duration")
}

type unexportedStruct struct {
    pattern string
    fn func()
}

func printWorld() {
    fmt.Println("world")
}

func mustPattern(op string) runtime.Pattern {
    p, _ := runtime.NewPattern(
        1,
        []int{int(utilities.OpLitPush), 0},
        []string{op},
        "")
    return p
}

// func basicAuth(username, password string) string {
// 	fmt.Println(fmt.Sprintf("#%s ", username))
// 	auth := username + ":" + password
// 	return base64.StdEncoding.EncodeToString([]byte(auth))
// }

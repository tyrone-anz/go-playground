package unexported

import "fmt"

type Sample struct {
	field []unexportedStruct
}

func (s *Sample) Run() {
	for _, i := range s.field {
		i.fn()
	}
}

type unexportedStruct struct {
	pattern string
	fn func()
}

func GetSample() *Sample{
	return &Sample{
		field: []unexportedStruct{
			{"hello", func(){fmt.Println("print hello")}},
			{"hi", func(){fmt.Println("print hi")}},
		},
	}
}

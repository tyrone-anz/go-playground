package runtimepattern

import (
	"fmt"
	"strings"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
)

const (
	validVersion = 1
	anything     = 0
)

func TestNewPattern(t *testing.T) {
	for _, spec := range []struct {
		ops  []int
		pool []string
		verb string

		stackSizeWant, tailLenWant int
	}{
		{},
		{
			ops:           []int{int(utilities.OpNop), anything},
			stackSizeWant: 0,
			tailLenWant:   0,
		},
		{
			ops:           []int{int(utilities.OpPush), anything},
			stackSizeWant: 1,
			tailLenWant:   0,
		},
		{
			ops:           []int{int(utilities.OpLitPush), 0},
			pool:          []string{"abc"},
			stackSizeWant: 1,
			tailLenWant:   0,
		},
		{
			ops:           []int{int(utilities.OpPushM), anything},
			stackSizeWant: 1,
			tailLenWant:   0,
		},
		{
			ops: []int{
				int(utilities.OpPush), anything,
				int(utilities.OpConcatN), 1,
			},
			stackSizeWant: 1,
			tailLenWant:   0,
		},
		{
			ops: []int{
				int(utilities.OpPush), anything,
				int(utilities.OpConcatN), 1,
				int(utilities.OpCapture), 0,
			},
			pool:          []string{"abc"},
			stackSizeWant: 1,
			tailLenWant:   0,
		},
		{
			ops: []int{
				int(utilities.OpPush), anything,
				int(utilities.OpLitPush), 0,
				int(utilities.OpLitPush), 1,
				int(utilities.OpPushM), anything,
				int(utilities.OpConcatN), 2,
				int(utilities.OpCapture), 2,
			},
			pool:          []string{"lit1", "lit2", "var1"},
			stackSizeWant: 4,
			tailLenWant:   0,
		},
		{
			ops: []int{
				int(utilities.OpPushM), anything,
				int(utilities.OpConcatN), 1,
				int(utilities.OpCapture), 2,
				int(utilities.OpLitPush), 0,
				int(utilities.OpLitPush), 1,
			},
			pool:          []string{"lit1", "lit2", "var1"},
			stackSizeWant: 2,
			tailLenWant:   2,
		},
		{
			ops: []int{
				int(utilities.OpLitPush), 0,
				int(utilities.OpLitPush), 1,
				int(utilities.OpPushM), anything,
				int(utilities.OpLitPush), 2,
				int(utilities.OpConcatN), 3,
				int(utilities.OpLitPush), 3,
				int(utilities.OpCapture), 4,
			},
			pool:          []string{"lit1", "lit2", "lit3", "lit4", "var1"},
			stackSizeWant: 4,
			tailLenWant:   2,
		},
		{
			ops:           []int{int(utilities.OpLitPush), 0},
			pool:          []string{"abc"},
			verb:          "LOCK",
			stackSizeWant: 1,
			tailLenWant:   0,
		},
	} {
		pat, err := runtime.NewPattern(validVersion, spec.ops, spec.pool, spec.verb)
		if err != nil {
			t.Errorf("NewPattern(%d, %v, %q, %q) failed with %v; want success", validVersion, spec.ops, spec.pool, spec.verb, err)
			continue
		}
		fmt.Println(pat.String())
	}
}

func TestMatch(t *testing.T) {
	for _, spec := range []struct {
		ops  []int
		pool []string
		verb string

		match    []string
		notMatch []string
	}{
		{
			match:    []string{""},
			notMatch: []string{"example"},
		},
		{
			ops:      []int{int(utilities.OpNop), anything},
			match:    []string{""},
			notMatch: []string{"example", "path/to/example"},
		},
		{
			ops:      []int{int(utilities.OpPush), anything},
			match:    []string{"abc", "def"},
			notMatch: []string{"", "abc/def"},
		},
		{
			ops:      []int{int(utilities.OpLitPush), 0},
			pool:     []string{"v1"},
			match:    []string{"v1"},
			notMatch: []string{"", "v2"},
		},
		{
			ops:   []int{int(utilities.OpPushM), anything},
			match: []string{"", "abc", "abc/def", "abc/def/ghi"},
		},
		{
			ops: []int{
				int(utilities.OpPushM), anything,
				int(utilities.OpLitPush), 0,
			},
			pool:  []string{"tail"},
			match: []string{"tail", "abc/tail", "abc/def/tail"},
			notMatch: []string{
				"", "abc", "abc/def",
				"tail/extra", "abc/tail/extra", "abc/def/tail/extra",
			},
		},
		{
			ops: []int{
				int(utilities.OpLitPush), 0,
				int(utilities.OpLitPush), 1,
				int(utilities.OpPush), anything,
				int(utilities.OpConcatN), 1,
				int(utilities.OpCapture), 2,
			},
			pool:  []string{"v1", "bucket", "name"},
			match: []string{"v1/bucket/my-bucket", "v1/bucket/our-bucket"},
			notMatch: []string{
				"",
				"v1",
				"v1/bucket",
				"v2/bucket/my-bucket",
				"v1/pubsub/my-topic",
			},
		},
		{
			ops: []int{
				int(utilities.OpLitPush), 0,
				int(utilities.OpLitPush), 1,
				int(utilities.OpPushM), anything,
				int(utilities.OpConcatN), 2,
				int(utilities.OpCapture), 2,
			},
			pool: []string{"v1", "o", "name"},
			match: []string{
				"v1/o",
				"v1/o/my-bucket",
				"v1/o/our-bucket",
				"v1/o/my-bucket/dir",
				"v1/o/my-bucket/dir/dir2",
				"v1/o/my-bucket/dir/dir2/obj",
			},
			notMatch: []string{
				"",
				"v1",
				"v2/o/my-bucket",
				"v1/b/my-bucket",
			},
		},
		{
			ops: []int{
				int(utilities.OpLitPush), 0,
				int(utilities.OpLitPush), 1,
				int(utilities.OpPush), anything,
				int(utilities.OpConcatN), 2,
				int(utilities.OpCapture), 2,
				int(utilities.OpLitPush), 3,
				int(utilities.OpPush), anything,
				int(utilities.OpConcatN), 1,
				int(utilities.OpCapture), 4,
			},
			pool: []string{"v2", "b", "name", "o", "oname"},
			match: []string{
				"v2/b/my-bucket/o/obj",
				"v2/b/our-bucket/o/obj",
				"v2/b/my-bucket/o/dir",
			},
			notMatch: []string{
				"",
				"v2",
				"v2/b",
				"v2/b/my-bucket",
				"v2/b/my-bucket/o",
			},
		},
		{
			ops:      []int{int(utilities.OpLitPush), 0},
			pool:     []string{"v1"},
			verb:     "LOCK",
			match:    []string{"v1:LOCK"},
			notMatch: []string{"v1", "LOCK"},
		},
	} {
		pat, err := runtime.NewPattern(validVersion, spec.ops, spec.pool, spec.verb)
		if err != nil {
			t.Errorf("NewPattern(%d, %v, %q, %q) failed with %v; want success", validVersion, spec.ops, spec.pool, spec.verb, err)
			continue
		}

		for _, path := range spec.match {
			_, err = pat.Match(segments(path))
			if err != nil {
				t.Errorf("pat.Match(%q) failed with %v; want success; pattern = (%v, %q)", path, err, spec.ops, spec.pool)
			}
		}
	}
}

func segments(path string) (components []string, verb string) {
	if path == "" {
		return nil, ""
	}
	components = strings.Split(path, "/")
	l := len(components)
	c := components[l-1]
	if idx := strings.LastIndex(c, ":"); idx >= 0 {
		components[l-1], verb = c[:idx], c[idx+1:]
	}
	return components, verb
}

func TestPatternString(t *testing.T) {
	for _, spec := range []struct {
		ops  []int
		pool []string

		want string
	}{
		{
			want: "/",
		},
		{
			ops:  []int{int(utilities.OpNop), anything},
			want: "/",
		},
		{
			ops:  []int{int(utilities.OpPush), anything},
			want: "/*",
		},
		{
			ops:  []int{int(utilities.OpLitPush), 0},
			pool: []string{"endpoint"},
			want: "/endpoint",
		},
		{
			ops:  []int{int(utilities.OpPushM), anything},
			want: "/**",
		},
		{
			ops: []int{
				int(utilities.OpPush), anything,
				int(utilities.OpConcatN), 1,
			},
			want: "/*",
		},
		{
			ops: []int{
				int(utilities.OpPush), anything,
				int(utilities.OpConcatN), 1,
				int(utilities.OpCapture), 0,
			},
			pool: []string{"name"},
			want: "/{name=*}",
		},
		{
			ops: []int{
				int(utilities.OpLitPush), 0,
				int(utilities.OpLitPush), 1,
				int(utilities.OpPush), anything,
				int(utilities.OpConcatN), 2,
				int(utilities.OpCapture), 2,
				int(utilities.OpLitPush), 3,
				int(utilities.OpPushM), anything,
				int(utilities.OpLitPush), 4,
				int(utilities.OpConcatN), 3,
				int(utilities.OpCapture), 6,
				int(utilities.OpLitPush), 5,
			},
			pool: []string{"v1", "buckets", "bucket_name", "objects", ".ext", "tail", "name"},
			want: "/v1/{bucket_name=buckets/*}/{name=objects/**/.ext}/tail",
		},
	} {
		p, err := runtime.NewPattern(validVersion, spec.ops, spec.pool, "")
		if err != nil {
			t.Errorf("NewPattern(%d, %v, %q, %q) failed with %v; want success", validVersion, spec.ops, spec.pool, "", err)
			continue
		}
		if got, want := p.String(), spec.want; got != want {
			t.Errorf("%#v.String() = %q; want %q", p, got, want)
		}

		verb := "LOCK"
		p, err = runtime.NewPattern(validVersion, spec.ops, spec.pool, verb)
		if err != nil {
			t.Errorf("NewPattern(%d, %v, %q, %q) failed with %v; want success", validVersion, spec.ops, spec.pool, verb, err)
			continue
		}
		if got, want := p.String(), fmt.Sprintf("%s:%s", spec.want, verb); got != want {
			t.Errorf("%#v.String() = %q; want %q", p, got, want)
		}
	}
}

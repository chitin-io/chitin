// +build gofuzz

package schema

import "bytes"

/*
Fuzz testing support files

https://github.com/dvyukov/go-fuzz

Usage:

    ./admin/fuzz-schema

See schema/testdata/fuzz/crashers for results.
*/

func Fuzz(data []byte) int {
	r := bytes.NewReader(data)
	_, err := Parse(r)
	if err != nil {
		return 0
	}
	return 1
}

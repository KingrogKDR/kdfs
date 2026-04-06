package custom_error

import "io"

func Check(e error) {
	if e != nil && e != io.EOF {
		panic(e)
	}
}

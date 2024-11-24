package main

import "io"

func unsafeClose(a ...io.Closer) {
	for _, f := range a {
		if f != nil {
			_ = f.Close()
		}
	}
}

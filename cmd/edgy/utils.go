package main

import "io"

func closeAll(a ...io.Closer) {
	for _, f := range a {
		if f != nil {
			_ = f.Close()
		}
	}
}

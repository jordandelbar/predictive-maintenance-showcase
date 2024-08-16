package main

import (
	"io"
	"log"
)

func closeOrLog(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Fatal(err.Error())
	}
}

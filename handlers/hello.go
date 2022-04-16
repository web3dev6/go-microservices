package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// log Request
	h.l.Println("Request received @hello handler")

	// Read from Request
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Error!", http.StatusBadRequest)
		h.l.Println("err: ", err)
		return
	}
	// log.Printf("Data : %s\n", d)
	h.l.Printf("Data : %s\n", d)

	// Write to response
	// fmt.Fprintf(rw, "Hello %s", d)
	rw.Write([]byte(fmt.Sprintf("Hello, %s!", d)))
}

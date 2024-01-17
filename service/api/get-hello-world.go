package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"fmt"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) getHelloWorld(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("Hello World!"))

	// testing error message
	err := fmt.Errorf("\nsomething went wrong")

	// Set the HTTP status code to indicate an error
	w.WriteHeader(http.StatusInternalServerError)

	// Send the error message as the response body
	fmt.Fprint(w, err.Error())
}

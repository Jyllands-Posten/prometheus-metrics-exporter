package server

import (
	"fmt"
	"log"
	"net/http"
	. "simpleTestServer/server/middleware"
	. "simpleTestServer/server/writers"
)

func Server(port *string) {

	jsonHandler := http.HandlerFunc(JsonWriter)
	htmlHandler := http.HandlerFunc(HtmlWriter)
	http.Handle("/jsonNoBasicAuth", MethodValidatorMiddleware(jsonHandler))
	http.Handle("/jsonBasicAuth", MethodValidatorMiddleware(BasicAuthMiddleware(jsonHandler)))
	http.Handle("/htmlNoBasicAuth", MethodValidatorMiddleware(htmlHandler))
	http.Handle("/htmlBasicAuth", MethodValidatorMiddleware(BasicAuthMiddleware(htmlHandler)))

	fPort := fmt.Sprintf(":%s", *port)

	err := http.ListenAndServe(fPort, nil)

	if err != nil {
		log.Fatal(err.Error())
	}
}

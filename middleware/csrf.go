package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/nosurf"
	"net/http"
)

func CSRF(ver string) {
	//TODO:
	csrf := nosurf.New(gin.Default())
	csrf.SetFailureHandler(http.HandlerFunc(csrfFailHandler))

	http.ListenAndServe(":8000", csrf)
}

func csrfFailHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s\n", nosurf.Reason(r))
}

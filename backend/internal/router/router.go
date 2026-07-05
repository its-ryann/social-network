package router

import (
	"net/http"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	
	return middleware.Logger(mux)
}
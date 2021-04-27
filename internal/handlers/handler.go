package handler

import "net/http"

type Handler interface {
	Register(router *http.ServeMux)
}

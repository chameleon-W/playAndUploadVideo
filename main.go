package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type middleWareHandler struct {
	r *httprouter.Router
	l *ConnLimiter
}

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()
	router.GET("/videos/:vid-id", streamHandler)

	router.POST("/upload/:vid-id", uploadHandler)

	// router.GET("/testpage", testPageHandler)

	return router
}

func NewMiddleWareHandler(r *httprouter.Router, num int) http.Handler {
	m := middleWareHandler{}
	m.r = r
	m.l = NewConnLimiter(num)
	return m
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//进行校验
	if !m.l.GetConn() {
		sendErrorResponse(w, http.StatusTooManyRequests, "Too many requests")
		return
	}

	m.r.ServeHTTP(w, r)
	defer m.l.ReleaseConn()
}

// 1.注册router
// 2.注册middleware
// 3.开启server
func main() {
	var num int
	r := RegisterHandlers()
	mh := NewMiddleWareHandler(r, num)
	http.ListenAndServe(":9000", mh)
}

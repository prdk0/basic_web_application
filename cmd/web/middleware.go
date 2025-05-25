package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/justinas/nosurf"
)

type MiddleWare func(http.Handler) http.Handler

func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hit the page")
		next.ServeHTTP(w, r)
	})
}

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

func MiddleWaretest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeNow := time.Now()
		log.Println("MiddleWareTest log: ", timeNow)
		next.ServeHTTP(w, r)
	})
}

func MiddleWareRecoverTest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %s\n%s", err, debug.Stack())
				http.Error(w, "Something went wrong", http.StatusInternalServerError)
			}
			log.Println("RecoverMiddleWare: done")
		}()
		next.ServeHTTP(w, r)
	})
}

func MiddleWareChainTest(h http.Handler, middleWare ...MiddleWare) http.Handler {
	for _, m := range middleWare {
		h = m(h)
	}
	return h
}

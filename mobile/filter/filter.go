package filter

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/khitrov-aleksandr/proxyguard/faker"
	"github.com/khitrov-aleksandr/proxyguard/logger"
	"github.com/khitrov-aleksandr/proxyguard/repository"
	"github.com/rs/zerolog/log"
)

type Filter struct {
	rp repository.Repository
	lg *logger.HandlerLogger
}

func New(rp repository.Repository, lg *logger.HandlerLogger) *Filter {
	return &Filter{rp: rp, lg: lg}
}

func (f *Filter) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if f.blockIpByRegister(r, f.rp) {
			JSONResponse(w, faker.GetTokenResponse(), http.StatusOK)
			return
		}

		if f.denyLogin(r, f.rp) {
			fmt.Println("denyLogin")
			JSONResponse(w, faker.GetLoginResponse(), http.StatusOK)
			return
			//fmt.Println("denyLogin1")
		}

		next.ServeHTTP(w, r)
	})

	/*
		return func(c echo.Context) error {
			if h.blockIpByRegister(c, h.rp) {
				return c.JSONPretty(http.StatusOK, faker.GetTokenResponse(), "")
			}

			if h.denyLogin(c, h.rp) {
				return c.JSONPretty(http.StatusOK, faker.GetLoginResponse(), "")
			}

			return next(c)
		}
	*/
}

func JSONResponse(w http.ResponseWriter, d any, code int) {
	data, err := json.Marshal(d)

	if err != nil {
		log.Print(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

package filter

import (
	"encoding/json"
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
			JSONResponse(w, faker.GetLoginResponse(), http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
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

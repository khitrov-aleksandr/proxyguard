package proxy

import (
	"net/http"

	"github.com/khitrov-aleksandr/proxyguard/faker"
	"github.com/khitrov-aleksandr/proxyguard/logger"
	"github.com/khitrov-aleksandr/proxyguard/repository"
)

type MobileService struct {
	rp repository.Repository
	lg *logger.HandlerLogger
}

func New(rp repository.Repository, lg *logger.HandlerLogger) *MobileService {
	return &MobileService{rp: rp, lg: lg}
}

func (f *MobileService) Handle(next http.Handler) http.Handler {
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

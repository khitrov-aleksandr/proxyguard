package blocker

import (
	"fmt"

	"github.com/khitrov-aleksandr/proxyguard/repository"
)

const (
	BlockTime = 86400
)

type RegisterBlocker struct {
	repository repository.Repository
}

func NewRegisterBlocker(repository repository.Repository) *RegisterBlocker {
	return &RegisterBlocker{
		repository: repository,
	}
}

func (r *RegisterBlocker) Block(ip string) {
	r.repository.Save(getKey(ip), ip, BlockTime)
}

func (r *RegisterBlocker) IsBlocked(ip string) bool {
	return r.repository.Get(getKey(ip)) == ip
}

func getKey(ip string) string {
	return fmt.Sprintf("reg_block:%s", ip)
}

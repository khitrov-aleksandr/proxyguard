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
	key := fmt.Sprintf("reg_block:%s", ip)
	fmt.Println(key)
	r.repository.Save(key, ip, BlockTime)
}

func (r *RegisterBlocker) IsBlocked(ip string) bool {
	key := fmt.Sprintf("reg_block:%s", ip)
	return r.repository.Get(key) == ip
}

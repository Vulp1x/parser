package domain

import "github.com/inst-api/parser/internal/pb/instaproxy"

func (b Bots) Proto() []*instaproxy.Bot {
	protos := make([]*instaproxy.Bot, len(b))
	for i := range b {
		protos[i] = (*instaproxy.Bot)(b[i])
	}

	return protos
}

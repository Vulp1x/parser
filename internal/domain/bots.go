package domain

//
// import (
// 	"github.com/inst-api/parser/internal/dbmodel"
// 	"github.com/inst-api/parser/pkg/api"
// )
//
// type Bots []*api.Bot
//
// func (b Bots) ToSaveParams() dbmodel.SaveBotsParams {
// 	var usernames, sessionIDs, proxies []string
//
// 	for _, bot := range b {
// 		if bot == nil {
// 			continue
// 		}
//
// 		usernames = append(usernames, bot.Username)
// 		sessionIDs = append(sessionIDs, bot.SessionId)
// 		proxies = append(proxies, dbmodel.Proxy{
// 			Host:  bot.GetProxy().Host,
// 			Port:  bot.GetProxy().Port,
// 			Login: bot.GetProxy().Login,
// 			Pass:  bot.GetProxy().Pass,
// 		}.String())
//
// 	}
//
// 	return dbmodel.SaveBotsParams{
// 		Usernames:  usernames,
// 		SessionIds: sessionIDs,
// 		Proxies:    proxies,
// 	}
// }

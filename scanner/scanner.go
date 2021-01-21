package scanner

import (
	"github.com/AkameMoe/BilibiliFilter/database"
	"github.com/AkameMoe/BilibiliFilter/proxy"
	"github.com/AkameMoe/BilibiliFilter/utils"
	"github.com/panjf2000/ants/v2"
	"strconv"
	"time"
)

var (
	pool, _ = ants.NewPool(50)
	faileduid []int
)

func StartScannerModule()  {
	uid := 1
	for true {
		if proxy.Proxyupdatedone {
			if pool.Free() != 0 {
				if len(faileduid) != 0 {
					for indexuid := range faileduid{
						scanUid(faileduid[indexuid])
					}
					faileduid = faileduid[:0]
				} else {
					scanUid(uid)
					uid++
				}
			} else {
				utils.Logger.Info().Msg(strconv.Itoa(pool.Running()))
				time.Sleep(time.Second*10)
				continue
			}
		}
	}
}

func scanUid(uid int) {
	utils.Logger.Info().Msg("Scanning " + strconv.Itoa(uid))
	pool.Submit(func() {
		user, exist, fail:= RequestUserData(uid)
		if fail {
			utils.Logger.Error().Msg(strconv.Itoa(uid) + " failed")
			faileduid = append(faileduid, uid)
		} else if exist {
			database.SaveUser(user)
			utils.Logger.Info().Msg(strconv.Itoa(user.Uid) + " / " + user.Name)
		}
	})
}
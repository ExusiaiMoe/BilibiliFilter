package scanner

import (
	"github.com/AkameMoe/BilibiliFilter/define"
	customproxy "github.com/AkameMoe/BilibiliFilter/proxy"
	"github.com/AkameMoe/BilibiliFilter/utils"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"strconv"
)

func RequestUserData(uid int) (user *define.User, exist bool, fail bool) {
	url := "https://api.bilibili.com/x/space/acc/info?mid=" + strconv.Itoa(uid)
	proxy, random := customproxy.GetProxy()
	client := &http.Client{Transport: &http.Transport{Dial: proxy.Dial}}

	response, err := client.Get(url)
	if err != nil {
		customproxy.RemoveProxy(random)
		utils.Logger.Error().Msg(err.Error())
		return nil, true, true
	}
	defer response.Body.Close()

	//statuscode := response.Header.Get("bili-status-code")
	//if statuscode != "" && statuscode == "-404" {
	//	return nil, false, false
	//}
	//
	//if statuscode != "" && statuscode == "-412" {
	//	return nil, true, true
	//}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		utils.Logger.Error().Msg(err.Error())
	}

	json, err := simplejson.NewJson(body)
	if err != nil {
		utils.Logger.Error().Msg(err.Error())
	}
	code, err := json.Get("code").Int()
	if err != nil {
		utils.Logger.Error().Msg(err.Error())
	}
	if code == -412 {
		return nil, true, true
	} else if code == -404 {
		return nil, false, false
	}


	level, err := json.Get("data").Get("level").Int()
	if err != nil {
		utils.Logger.Info().Msg(string(body))
		utils.Logger.Error().Msg(err.Error())
	}
	silenceInt, err := json.Get("data").Get("silence").Int()
	if err != nil {
		utils.Logger.Error().Msg(err.Error())
	}
	var silence bool
	if silenceInt == 0 {
		silence = false
	} else {
		silence = true
	}
	name, err := json.Get("data").Get("name").String()
	if err != nil {
		utils.Logger.Error().Msg(err.Error())
	}
	userdata := &define.User{Uid: uid, Level: level, Silence: silence, Name: name}
	return userdata, true, false
}

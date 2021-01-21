package proxy

import (
	"bufio"
	"github.com/AkameMoe/BilibiliFilter/utils"
	"io/ioutil"
	"crypto/rand"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"golang.org/x/net/proxy"
	"github.com/robfig/cron/v3"
)

var (
	proxyList       []string
	locker sync.Mutex
	Proxyupdatedone bool
)

func GetProxy() (dialer proxy.Dialer, indexnum int) {
	//rand.Seed(time.Now().UnixNano())
	random64, err := rand.Int(rand.Reader, big.NewInt(int64(len(proxyList))))
	random, err := strconv.Atoi(random64.String())
	address := proxyList[random]
	dia,err := proxy.SOCKS5("tcp", address, nil, proxy.Direct)
	if err != nil {
		utils.Logger.Info().Msg("can't connect to the proxy:" + err.Error())
	}
	return dia, random
}

func RemoveProxy(indexnum int)  {
	utils.Logger.Info().Msg("Removing " + proxyList[indexnum])
	locker.Lock()
	proxyList = append(proxyList[:indexnum], proxyList[indexnum+1:]...)
	locker.Unlock()
}

func StartProxyModule() {
	Proxyupdatedone = false
	getAddress()
	proxylistcron := cron.New()
	proxylistcron.AddFunc("@every 10m", getAddress)
	proxylistcron.Start()
	utils.Logger.Info().Msg("Proxy Module Started Successfully")
	select {}
}

func getAddress() {
	Proxyupdatedone = false
	response, err := http.Get("https://www.proxy-list.download/api/v1/get?type=socks5&anon=elite&country=CN")
	if err != nil {
		utils.Logger.Warn().Msg(err.Error())
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		utils.Logger.Warn().Msg(err.Error())
	}
	scanner := bufio.NewScanner(strings.NewReader(string(body)))
	proxyList = proxyList[:0]
	for scanner.Scan() {
		utils.Logger.Info().Msg("Checking proxy " + scanner.Text())
		if checkProxy(scanner.Text()) {
			locker.Lock()
			proxyList = append(proxyList, scanner.Text())
			locker.Unlock()
		}
	}
	utils.Logger.Info().Msg("Get Proxy List Success")
}

func checkProxy(address string) bool {
	dialer, err := proxy.SOCKS5("tcp", address, nil, proxy.Direct)
	if err != nil {
		return false
	}
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport, Timeout: 5*time.Second}
	httpTransport.Dial = dialer.Dial
	resp, err := httpClient.Get("https://api.bilibili.com/x/space/acc/info?mid=1")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	utils.Logger.Info().Msg("Proxy " + address + " avaliable")
	Proxyupdatedone = true
	return true
}

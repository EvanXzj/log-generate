package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// user-agent list
var uaList = []string{
	"Mozilla/5.0 (Windows; U; Windows NT 5.2) Gecko/2008070208 Firefox/3.0.1",
	"Mozilla/5.0 (Windows; U; Windows NT 5.1) Gecko/20070309 Firefox/2.0.0.3",
	"Mozilla/5.0 (Windows; U; Windows NT 5.1) Gecko/20070803 Firefox/1.5.0.12",
	"Mozilla/5.0 (Windows; U; Windows NT 5.2) AppleWebKit/525.13 (KHTML, like Gecko) Version/3.1 Safari/525.13",
	"Mozilla/5.0 (iPhone; U; CPU like Mac OS X) AppleWebKit/420.1 (KHTML, like Gecko) Version/3.0 Mobile/4A93 Safari/419.3",
	"Mozilla/5.0 (Windows; U; Windows NT 5.2) AppleWebKit/525.13 (KHTML, like Gecko) Chrome/0.2.149.27 Safari/525.13",
	"Mozilla/5.0 (Windows; U; Windows NT 5.1; en-US; rv:1.8.1.12) Gecko/20080219 Firefox/2.0.0.12 Navigator/9.0.0.6",
	"Mozilla/5.0 (Linux; U; Android 4.0.3; zh-cn; M032 Build/IML74K) AppleWebKit/533.1 (KHTML, like Gecko)Version/4.0 MQQBrowser/4.1 Mobile Safari/533.1",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 5_1_1 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9B206 Safari/7534.48.3",
}

type resource struct {
	url,
	target string
	start,
	end int
}

func main() {
	// parse command-line params
	totalNum := flag.Int("totalNum", 100, "total number of log rows need to generate")
	logFilePath := flag.String("logFilePath", "/tmp/logs/dig.log", "log file path")
	flag.Parse()

	// get resource
	res := genResourceList()
	// build url
	uList := genURL(res)

	logStr := ""
	uLength := len(uList)
	uaLength := len(uaList)
	// generate log
	for i := 0; i < *totalNum; i++ {
		u := uList[randomInt(0, uLength)]
		refer := uList[randomInt(0, uLength)]
		ua := uaList[randomInt(0, uaLength)]

		logStr = logStr + genLog(u, refer, ua) + "\n"
	}

	// write log to file
	file, err := os.OpenFile(*logFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Printf("Open File: filePath=%s Error=%s\n", *logFilePath, err.Error())

		return
	}

	file.Write([]byte(logStr))
	defer file.Close()

	// prompt hint
	fmt.Println("Done", time.Now().Format("2006-01-02 15:04:05"))
}

// get ran random int number
func randomInt(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	if min > max {
		return max
	}

	return r.Intn(max-min) + min
}

func genResourceList() []resource {
	var res []resource

	// 首页
	homepage := resource{
		url:    "http://localhost:7888/",
		target: "",
		start:  0,
		end:    0,
	}

	// 列表页
	category := resource{
		url:    "http://localhost:7888/list/${id}.html",
		target: "${id}",
		start:  1,
		end:    21,
	}

	// 详情页
	detail := resource{
		url:    "http://localhost:7888/move/${id}.html",
		target: "${id}",
		start:  1,
		end:    12924,
	}

	res = append(res, homepage, category, detail)

	return res
}

func genURL(res []resource) []string {
	var list []string

	for _, item := range res {
		if len(item.target) == 0 {
			list = append(list, item.url)
		} else {
			for i := item.start; i <= item.end; i++ {
				urlStr := strings.Replace(item.url, item.target, strconv.Itoa(i), -1)
				list = append(list, urlStr)
			}
		}
	}

	return list
}

func genLog(currentURL, refer, ua string) string {
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	// timeStr := time.Now().UTC().String()

	u := url.Values{}

	u.Set("time", timeStr)
	u.Set("url", currentURL)
	u.Set("refer", refer)
	u.Set("ua", ua)
	paramsStr := u.Encode()

	logTemplate := "127.0.0.1 - - [${timeStr}] \"OPTIONS /dig?${paramsStr} HTTP/1.1\" 200 43 \"-\" \"${ua}\" \"-\""

	log := strings.Replace(logTemplate, "${paramsStr}", paramsStr, -1)
	log = strings.Replace(log, "${timeStr}", timeStr, -1)
	log = strings.Replace(log, "${ua}", ua, -1)

	return log
}

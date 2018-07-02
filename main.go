package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	totalNum := flag.Int("totalNum", 10, "total number of log rows need to generate")
	logFilePath := flag.String("logFilePath", "/tmp/logs/runtime.log", "log file path")
	flag.Parse()

	fmt.Println(*totalNum, *logFilePath)

	res := genResourceList()
	list := genURL(res)

	var logStr string
	// generate log
	for i := 0; i < *totalNum; i++ {
		index := randomNumber()
		logStr += genLog(list[index], list[index], "asdfa") + "\n"
	}

	file, err := os.OpenFile(*logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Printf("Open File: filePath=%s Error=%s\n", *logFilePath, err.Error())

		return
	}
	file.Write([]byte(logStr))
	defer file.Close()

	fmt.Printf("URL list length=%d\n", len(list))
	fmt.Println(list[0])
	fmt.Println("Done")
}

func randomNumber() int {

	return 0
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

func genLog(url, refer, ua string) string {

	return "success"
}

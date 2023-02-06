package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/SteamHostSync/fileIO"
	"github.com/levigross/grequests"
)

var cf *fileIO.Config
var L int

func init() {
	cf, L = fileIO.ReadConfig()
	cstSh, _ := time.LoadLocation("Asia/Shanghai")
	time.Local = cstSh
}

var result string
var resultTmp string
var Tempplat string

func main() {
	for i := 0; i < L; i++ {
		for i, v := range (*cf).Platforms[i] {
			if i == 0 {
				resultTmp += fmt.Sprintf("#%s Start\n", v)
				fmt.Printf("####################%s Start####################\n", v)
				Tempplat = v
				continue
			}
			ip, flag := getip(v)
			log.Println(flag)
			if flag == false {
				resultTmp += fmt.Sprintf("%s\t\t\t%s\n", ip, v)
				fmt.Printf("%s\t\t\t%s\n", ip, v)
			} else {
				resultTmp += fmt.Sprintf("####%s\t\t\t%s\n", ip, v)
				fmt.Printf("####%s\t\t\t%s\n", ip, v)
			}

		}
		resultTmp += fmt.Sprintf("#%s End\n", Tempplat)
		resultTmp += fmt.Sprintf("# Last Update Time : %s \n\n", time.Now().Format("2006-01-02 15:04:05"))
		fmt.Printf("####################%s End####################\n", Tempplat)
		fmt.Printf("# Last Update Time :  %s \n\n", time.Now().Format("2006-01-02 15:04:05"))
		fileIO.WriteHost(resultTmp, "Hosts"+"_"+Tempplat)
		result += resultTmp
		resultTmp = ""
	}
	result += "#Github: https://github.com/Clov614/SteamHostSync\n"
	fileIO.WriteHost(result, "Hosts")
	_, err := os.Stat("./README_TEMP.md")
	if err != nil {
		return
	}
	content := fileIO.Read_tmp()
	content = strings.Replace(content, "HOST_TARGET", result, 1)
	fileIO.WriteFile(content, "README.md")
}

func getip(url string) (res string, flag bool) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("获取ip错误：", err)
			flag = true
		}
	}()
	//RO := grequests.RequestOptions{
	//	//Headers: map[string]string{
	//	//	"User-Agent": "[{\"key\":\"User-Agent\",\"value\":\"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36\",\"description\":\"\",\"type\":\"text\",\"enabled\":true}]",
	//	//	"Host":       "www.ipaddress.com",
	//	//},
	//	InsecureSkipVerify: true,
	//	//UserAgent:          (*cf).UA,
	//	//Host:               "www.ipaddress.com",
	//}
	resp, err := grequests.Get("https://www.ipaddress.com/site/"+strings.TrimSpace(url), nil)
	if err != nil {
		log.Fatalln("Unable to make request:", err)
	}
	re := regexp.MustCompile(`/ipv4/(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})">`)
	result := re.FindStringSubmatch(resp.String())
	res = result[1]
	return
}

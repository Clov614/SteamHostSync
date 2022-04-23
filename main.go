package main

import (
	"fmt"
	"github.com/SteamHostSync/fileIO"
	"github.com/levigross/grequests"
	"log"
	"regexp"
	"strings"
	"time"
)

var cf *fileIO.Config
var L int

func init() {
	cf, L = fileIO.ReadConfig()
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
			ip := getip(v)
			resultTmp += fmt.Sprintf("%s\t\t\t%s\n", ip, v)
			fmt.Printf("%s\t\t\t%s\n", ip, v)
		}
		resultTmp += fmt.Sprintf("#%s End\n", Tempplat)
		resultTmp += fmt.Sprintf("# Last Update Time : %s \n\n", time.Now().Format("2006-01-02 15:04:05"))
		fmt.Printf("####################%s End####################\n", Tempplat)
		fmt.Printf("# Last Update Time :  %s \n\n", time.Now().Format("2006-01-02 15:04:05"))
		fileIO.WriteHost(resultTmp, "Hosts"+"_"+Tempplat)
		result += resultTmp
		resultTmp = ""
	}
	result += "Github: https://github.com/Clov614/SteamHostSync\n"
	fileIO.WriteHost(result, "Hosts")
	content := fileIO.Read_tmp()
	content = strings.Replace(content, "HOST_TARGET", result, 1)
	fileIO.WriteFile(content, "README.md")
	HTML := fileIO.ReadHtml()
	HTML = strings.Replace(HTML, "#TARGET#", result, 1)
	fileIO.WriteFile(HTML, "index.html")
}

func getip(url string) string {
	RO := grequests.RequestOptions{
		UserAgent: (*cf).UA,
	}
	resp, err := grequests.Get("https://ipaddress.com/website/"+url, &RO)
	if err != nil {
		log.Fatalln("Unable to make request:", err)
	}

	re := regexp.MustCompile(`<li>(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})</li>`)
	result := re.FindStringSubmatch(resp.String())

	return result[1]
}

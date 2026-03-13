package fileIO

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// PathExists
func init() {
	_, err := os.Stat("./Source.yaml")
	if err != nil {
		CreateSource()
		fmt.Println("生成默认Source.yaml")
		return
	}
}

func CreateSource() {
	cf := Config{
		UA:        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36",
		Platforms: urls,
	}

	bytes, err := yaml.Marshal(&cf)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ioutil.WriteFile("./Source.yaml", bytes, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}

var urls [][]string = [][]string{
	{"github",
		"alive.github.com",
		"live.github.com",
		"github.githubassets.com",
		"central.github.com",
		"desktop.githubusercontent.com",
		"assets-cdn.github.com",
		"camo.githubusercontent.com",
		"github.map.fastly.net",
		"github.global.ssl.fastly.net",
		"gist.github.com",
		"github.io",
		"github.com",
		"github.blog",
		"api.github.com",
		"raw.githubusercontent.com",
		"user-images.githubusercontent.com",
		"favicons.githubusercontent.com",
		"avatars5.githubusercontent.com",
		"avatars4.githubusercontent.com",
		"avatars3.githubusercontent.com",
		"avatars2.githubusercontent.com",
		"avatars1.githubusercontent.com",
		"avatars0.githubusercontent.com",
		"avatars.githubusercontent.com",
		"codeload.github.com",
		"github-cloud.s3.amazonaws.com",
		"github-com.s3.amazonaws.com",
		"github-production-release-asset-2e65be.s3.amazonaws.com",
		"github-production-user-asset-6210df.s3.amazonaws.com",
		"github-production-repository-file-5c1aeb.s3.amazonaws.com",
		"githubstatus.com",
		"github.community",
		"github.dev",
		"media.githubusercontent.com"},
	{
		"steam",
		"steamcommunity.com",
		"www.steamcommunity.com",
		"store.steampowered.com",
		"api.steampowered.com",
		"help.steampowered.com",
		"store.akamai.steamstatic.com",
		"steamcdn-a.akamaihd.net",
		"cdn.akamai.steamstatic.com",
		"steam-chat.com",
		"community.akamai.steamstatic.com",
	},
	{
		"Ubisoft_download",
		"static3.cdn.Ubi.com",
		"static2.cdn.Ubi.com",
		"static1.cdn.Ubi.com",
	},
}

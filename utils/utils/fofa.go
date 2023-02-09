package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"scan/utils/requests"
)

type FofaResp struct {
	Err     bool       `json:"error"`
	Mode    string     `json:"mode"`
	Page    int        `json:"page"`
	Query   string     `json:"query"`
	Results [][]string `json:"results"`
	Size    int        `json:"size"`
}

var _headers = requests.Header{

	"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
	"Accept-Language": "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7",
	"Cache-Control":   "max-age=0",
	"Connection":      "close",
	"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36",
}

var cloudOrg = []string{"Alibaba", "Huawei Cloud", "Tencent"}

func (fofa *FofaResp) Crawler(filename string, rule string, size int, pageNum int, isCloud bool) {

	var fofaResp FofaResp

	req := requests.Requests()

	fofaEmail, fofaApiKey := fofaReadFile(filename)
	if fofaEmail == "" ||  fofaApiKey == ""{
		fmt.Println("请检查 fofa token 文件")
		return;
	}


	ruleBase64 := StandBase64([]byte(rule))

	params := requests.Params{
		"email":   fofaEmail,
		"key":     fofaApiKey,
		"qbase64": ruleBase64,
		"size":    fmt.Sprintf("%d", size),
		"page":    fmt.Sprintf("%d", pageNum),
		"fields":  "host,title,ip,domain,port,country,as_organization,server,protocol",
	}

	resp, err := req.Get("https://fofa.info/api/v1/search/all", _headers, params)

	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(resp.Content(), &fofaResp)

	if err != nil {
		fmt.Printf("err=%v", err)
	}

	fmt.Printf("get %d host\n", len(fofaResp.Results))
	fmt.Println(fofaResp.Size)
	for _, value := range fofaResp.Results {
		host := value[0]
		region := value[5]
		org := value[6]
		if isCloud {
			if !IsContain(cloudOrg, org) {
				fmt.Printf("ip: %s  reg: %s  org: %s \n", host, region, org)

			}

		} else {
			fmt.Printf("ip: %s  reg: %s  org: %s \n", host, region, org)
		}

	}

}

func fofaReadFile(filename string) (string, string) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", ""
	}

	str := string(content)

	emailCompileRegex := regexp.MustCompile("email\\s+=\\s+(.*)\\s+")
	tokenCompileRegex := regexp.MustCompile("token\\s+=\\s+(.*)\\s+")
	if len(emailCompileRegex.FindStringSubmatch(str)) == 0 ||  len(tokenCompileRegex.FindStringSubmatch(str)) == 0{
		return "", ""
	}
	email := emailCompileRegex.FindStringSubmatch(str)[1]
	token := tokenCompileRegex.FindStringSubmatch(str)[1]


	return email, token

}

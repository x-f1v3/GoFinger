package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"scan/json_new"
	"scan/utils/requests"
	"scan/utils/utils"
	"strings"
	"time"

	"github.com/spaolacci/murmur3"
)

var _headers = requests.Header{

	"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
	"Accept-Language": "en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7",
	"Cache-Control":   "max-age=0",
	"Cookie":          "rememberMe=ABCD;rememberme=ABCD;Rememberme=ABCD;RememberMe=ABCD;",
	"Connection":      "close",
	"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36",
}

// 获取证书内容，参考byro07/fwhatweb
func getCerts(resp *http.Response) []byte {
	var certs []byte
	if resp.TLS != nil {
		cert := resp.TLS.PeerCertificates[0]
		var str string
		if js, err := json.Marshal(cert); err == nil {
			certs = js
		}
		str = string(certs) + cert.Issuer.String() + cert.Subject.String()
		certs = []byte(str)
	}
	return certs
}

func HeaderToString(header http.Header) string {
	res := ""
	for name, values := range header {
		for _, value := range values {
			res = res + fmt.Sprintf("%s: %s", name, value) + ","
		}
	}
	return res
}

func GetOtherReq(Uri string, args ...interface{}) string {

	url := ""
	uri := "/" + Uri
	req := requests.Requests()
	for _, arg := range args {
		switch args := arg.(type) { //.(type)是判断interface{}的具体类型，只能在switch case语句中使用。
		case *requests.Request:
			req = args
		case string:
			url = args
		}
	}
	req.SkipVerify()
	resp, err := req.Get(url+uri, _headers)
	if err != nil {
		return ""
	}
	if resp.R.StatusCode != 200 {
		return ""

	}

	
	return string(resp.Content())

}

func GetIcoHashOnce(url string) string {

	req := requests.Requests()
	req.SkipVerify()
	return GetIcoHash(req,url,[]byte(GetOtherReq("",url)))


}

func xegexpjs(reg string, resp string) (reslut1 [][]string) {
	reg1 := regexp.MustCompile(reg)
	if reg1 == nil {
		return nil
	}
	result1 := reg1.FindAllStringSubmatch(resp, -1)
	return result1
}

func GetIcoHash(req *requests.Request, url string, Content []byte) string {
	url = strings.Trim(url, "/")
	faviconpaths := xegexpjs(`href="*(.*?favicon....)"*`, string(Content))
	var faviconpath string
	if len(faviconpaths) > 0 {
		fav := faviconpaths[0][1]
		if fav[:2] == "//" {
			faviconpath = "http:" + fav
		} else {
			if fav[:4] == "http" {
				faviconpath = fav
			} else if strings.HasPrefix(fav, "/"){
				faviconpath = url + fav
				
				
			}else{
				faviconpath = url + "/" + fav
			}
		}
	} else {
		faviconpath = url + "/favicon.ico"
	}

	resp, err := req.Get(faviconpath, _headers)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if resp.R.StatusCode == 200 {
		mHash := murmur3.New32()
		_, errorr := mHash.Write([]byte(utils.StandBase64(resp.Content())))
		if errorr != nil {
			//fmt.Println(err) 只用于调试
		}
		hashNum := mHash.Sum32()
		
		

		return fmt.Sprintf("%d", int32(hashNum))
	} else {
		return ""
	}
}

func Reqdata(url string, proxy string) (*json_new.FetchResult, *requests.Request, string, error) {
	var req_data = json_new.FetchResult{}
	var certs []byte = nil
	retried := false
	req := requests.Requests()
	req.SetTimeout(time.Duration(15))
	req.SkipVerify()
	req.SetRedirect(5)
	if proxy != "" {
		req.Proxy(proxy)
	}
	if !strings.HasPrefix(url, "http") {

		url = "https://" + url
	}

retry:
	resp, err := req.Get(url, _headers)
	//print(err.Error())

	if err != nil {
		if !retried {
			//fmt.Println(err) 只用于调试
			url = strings.Replace(url, "https://", "http://", 1)
			retried = true
			goto retry
		}

		return nil, nil, "", errors.New("Error")

	}

	if strings.HasPrefix(url, "http") {
		certs = getCerts(resp.R)
	}
	defer resp.R.Body.Close()

	var resqRedirect *requests.Response = htmlRedirect(string(resp.Content()), url, req, certs)

	if resqRedirect != nil {
		resp = resqRedirect
	}
	req_data = json_new.FetchResult{
		Url:          url,
		Content:      resp.Content(),
		Headers:      resp.R.Header,
		Status:       resp.R.Status,
		HeaderString: HeaderToString(resp.R.Header),
		FaviconHash:  GetIcoHash(req,url,resp.Content()), // GetOtherReq("favicon.ico", req, url)
		OtherString:  GetOtherReq("robots.txt", req, url),
		Certs:        certs,
	}

	return &req_data, req, url, nil

}

func htmlRedirect(respContent string, url string, req *requests.Request, certs []byte) *requests.Response {
	var result bool = false
	var redirectUrl string = ""

	needRedirect := []string{
		"[<html>]*\\s*<script\\s*[type\\s*=\\s*\"text/javascript\"]*\\s*>\\s*window.location\\s*=\\s*[\"'](.*?)[\"'][;]*\\s*</script>\\s*[</html>]*\\s*$",
		"(?i)<meta\\s*http-equiv=\"refresh\"\\s*content=\"\\d;\\s*URL=(.*?)\"", //忽略大小写
		"[<html>]*.*?<script\\s*[type\\s*=\\s*\"text/javascript\"]*\\s*>\\s*location.href\\s*=\\s*[\"'](.*?)[\"'][;]*\\s*</script>\\s*.*?\\s*[<body>]*\\s*[</body>]*\\s*[</html>]*\\s*$",
		// "(window|top)\\.location\\.href\\s*=\\s*[\"'](.*?)[\"']",
		// "redirectUrl = \"(.*?)\"",
		// "<meta.*?http-equiv=.*?refresh.*?url=(.*?)>",
	}

	for _, rule := range needRedirect {
		compileRegex := regexp.MustCompile(rule)
		matchArr := compileRegex.FindStringSubmatch(respContent)

		if len(matchArr) > 0 {
			result = true
			redirectUrl = matchArr[len(matchArr)-1]
			break
		}

	}

	if strings.HasPrefix(redirectUrl, "/") || strings.HasPrefix(redirectUrl, "\\") {
		url = url + redirectUrl[1:]
	} else {
		url = url + redirectUrl
	}

	if result && redirectUrl != "" {

		resp, err := req.Get(url, _headers)
		if err != nil {
			//fmt.Println(err)  只用于调试
			return nil

		}

		defer resp.R.Body.Close()

		return resp

	}

	return nil

}

func ReqWafdata(req *requests.Request, url string, proxy string) (*json_new.FetchResult, error) {
	var req_data = json_new.FetchResult{}
	var payload string
	//var payload = requests.Params{
	//	"a": "cat%20/etc/passwd",
	//}
	if proxy != "" {
		req.Proxy(proxy)
	}

	if strings.HasSuffix(url,"/"){
		 payload  = "?a=cat%20/etc/passwd"
	}else{
		payload  = "?a=cat%20/etc/passwd"
	}

	req.SetRedirect(-1)
	resp, err := req.Get(url+payload, _headers)
	if err != nil {
		return nil, errors.New("Error")

	}
	req_data = json_new.FetchResult{
		Url:          url,
		Content:      resp.Content(),
		Headers:      resp.R.Header,
		HeaderString: HeaderToString(resp.R.Header),
		FaviconHash:  "",
		OtherString:  "",
		Certs:        nil,
	}

	return &req_data, nil

}

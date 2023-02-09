package json_new

import (
	"fmt"
	"regexp"
	"scan/utils/utils"
	"strings"
)

func Detect(resp *FetchResult, finger []Finger, silent bool) (res string) {
	if resp.IsEmpty() {
		return
	}

	products := make([]string, 0)
	//获取网页返回数据并赋值
	web_Content := string(resp.Content)
	//web_Headers := resp.Headers
	certString := string(resp.Certs)
	web_Certs := resp.Certs
	web_HeaderString := strings.ToLower(resp.HeaderString)
	web_faviconHash := resp.FaviconHash
	headerServerString := strings.ToLower(fmt.Sprintf("Server : %v\n", resp.Headers["Server"]))
	web_Title := webTitle(web_Content)
	web_Content = strings.ToLower(web_Content)
	Other_String := strings.ToLower(resp.OtherString)

	for _, fp := range finger {
		//fofa指纹中的最后一项
		rules := fp.Rules
		matchFlag := false
		//其中的match
		//matchFlag := false
		//对每个json的最后一项进行迭代
		for _, onerule := range rules {

			//控制继续器
			ruleMatchContinueFlag := true

			onerule_len := len(onerule)
			for k, rule := range onerule {

				if !ruleMatchContinueFlag {
					break
				}

				lowerRuleContent := strings.ToLower(rule.Content)
				// fmt.Println(lowerRuleContent)
				// fmt.Println(strings.Split(rule.Match, "_")[0])

				switch strings.Split(rule.Match, "_")[0] {
				case "banner":
					reBanner := regexp.MustCompile(`(?im)<\s*banner.*>(.*?)<\s*/\s*banner>`)
					matchResults := reBanner.FindAllString(web_Content, -1)
					if len(matchResults) == 0 {
						ruleMatchContinueFlag = false
						break
					}

					for _, matchResult := range matchResults {
						if !strings.Contains(strings.ToLower(matchResult), lowerRuleContent) {
							ruleMatchContinueFlag = false
							break
						}

					}
				case "title":
					if web_Title == "" {
						ruleMatchContinueFlag = false
						break
					}
					if !strings.Contains(strings.ToLower(web_Title), lowerRuleContent) {
						ruleMatchContinueFlag = false
					}

				case "body":
					if !strings.Contains(web_Content, lowerRuleContent) {
						ruleMatchContinueFlag = false
					}

				case "favicon":
					if !strings.Contains(web_faviconHash, lowerRuleContent) {
						ruleMatchContinueFlag = false
					}

				case "header":
					if !strings.Contains(web_HeaderString, rule.Content) {
						ruleMatchContinueFlag = false
					}
				case "server":
					if !strings.Contains(headerServerString, rule.Content) {
						ruleMatchContinueFlag = false
					}

				case "robots":
					if !strings.Contains(Other_String, rule.Content) {
						ruleMatchContinueFlag = false
					}
				case "cert":
					if (web_Certs == nil) || (web_Certs != nil && !strings.Contains(certString, rule.Content)) {
						ruleMatchContinueFlag = false
					}
				default:
					ruleMatchContinueFlag = false

				}

				// 单个rule之间是AND关系，匹配成功
				if ruleMatchContinueFlag && k == onerule_len-1 {
					matchFlag = true
					break
				}

			}

		}

		// 多个rule之间是OR关系，匹配成功
		if matchFlag {
			products = append(products, fp.Product)
		}

	}

	utils.Success(resp.Url, web_Title, strings.Join(products, ", "), silent)

	return strings.Join(products, ", ")

}

func webTitle(content string) string {
	title := regexp.MustCompile(`(?im)<\s*title.*>(.*?)<\s*/\s*title>`)
	matchResults := title.FindStringSubmatch(content)
	if len(matchResults) == 0 {
		return ""
	}
	return string(utils.ByteToUTF8([]byte(matchResults[1])))

}

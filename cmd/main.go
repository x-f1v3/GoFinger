package main

import (
	"flag"
	"fmt"
	"os"
	"scan/json_new"
	"scan/network"
	"strings"

	"scan/utils/requests"
	"scan/utils/utils"

	"github.com/remeh/sizedwaitgroup"
)

var (
	file       string
	fav        bool
	fofa       bool
	isCloud    bool
	isWaf      bool
	isOnlyWaf  bool
	silent     bool
	nuclei     bool
	url        string
	fingerfile string
	waffile    string
	threads    = 2
	proxy      string
	fofaToken  string
	fofaRule   string
	fofaFile   string
	fofaPage   int
	fofaSize   int
	filterStatus string
	swg        sizedwaitgroup.SizedWaitGroup
)

func usage() {
	fmt.Fprintf(os.Stderr, `a simple tool for pentest

Options:
`)
	flag.PrintDefaults()
}

func init() {
	flag.StringVar(&file, "f", "", "文件名称")
	flag.StringVar(&url, "u", "", "单个URL")
	flag.StringVar(&proxy, "p", "", "代理,eg http://127.0.0.1:8080")
	flag.StringVar(&filterStatus,"fst","502,503","过滤状态码,默认过滤 502,503")
	flag.IntVar(&threads, "t", 10, "线程数")
	flag.BoolVar(&fav, "fav", false, "开关，计算 favicon 的 hash ,和 -u 联动")
	flag.BoolVar(&isWaf, "isWaf", false, "同时识别waf和CMS")
	flag.BoolVar(&isOnlyWaf, "isOnlyWaf", false, "只识别waf")
	flag.IntVar(&fofaPage, "fp", 1, "fofa 搜索页数")
	flag.IntVar(&fofaSize, "fs", 20, "fofa 搜索一页数量")
	flag.StringVar(&fofaRule, "fr", "", "fofa 规则")
	flag.StringVar(&fofaFile, "ft", "token.txt", "fofa token 文件")
	flag.BoolVar(&isCloud, "fc", false, "开关，结果是否过滤云服务器")
	flag.BoolVar(&nuclei, "exp", false, "开关,结果是否进行nuclei扫描")
	flag.BoolVar(&fofa, "fofa", false, "开关,fofa 搜索 ,和 -fp -fs -fr -ff -fc 联动")
	flag.BoolVar(&silent, "silent", false, "不花里胡哨的输出")
	flag.StringVar(&fingerfile, "ffile", "finger.json", "cms指纹文件")
	flag.StringVar(&waffile, "wfile", "waf.json", "waf指纹文件")

	// 重写 usage
	flag.Usage = usage
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func processFinger(url string, proxy string, finger []json_new.Finger, nuclei bool,filterStatus string) (*requests.Request, string) {
	res, req, url, error := network.Reqdata(url, proxy)
	fsl := strings.Split(filterStatus,",")
	product := ""
	if error == nil {
			for _, status := range fsl {
				if strings.Contains(res.Status, status){
					return nil, ""
				}
			}
		//fmt.Println(res.Content)
		product = json_new.Detect(res, finger, silent)
		if nuclei {
			utils.NucleiRun(url, product)

		}
		return req, url
	}

	return nil, ""

}

func processWaf(req *requests.Request, url string, finger []json_new.Finger, proxy string) {
	res, error := network.ReqWafdata(req, url,proxy)
	if error == nil {
		json_new.Detect(res, finger, silent)
	}
	return

}

func main() {

	flag.Parse()

	if flag.NFlag() == 0 {
		flag.Usage()
		return
	}
	//执行的函数

	swg = sizedwaitgroup.New(threads)

	//直接扫描
	if len(url) != 0 {
		if fav {
			faviconHash := network.GetIcoHashOnce(url)
			fmt.Println(faviconHash)
		} else {
			if isOnlyWaf {
				waf, _ := json_new.Parse(waffile, silent)
				_, req, url, _ := network.Reqdata(url, proxy)
				processWaf(req, url, waf,proxy)

			} else {
				finger, _ := json_new.Parse(fingerfile, silent)
				req, url := processFinger(url, proxy, finger, nuclei,filterStatus)
				if isWaf && (req != nil) {
					waf, _ := json_new.Parse(waffile, silent)
					processWaf(req, url, waf,proxy)
				}

			}

		}
	} else if len(file) != 0 {
		domains, _ := utils.ReadFile(file)
		finger, _ := json_new.Parse(fingerfile, silent)
		waf, _ := json_new.Parse(waffile, silent)

		for i := 0; i < len(domains); i++ {
			swg.Add()
			// 开启一个并发
			go func(url string) {
				// 使用defer, 表示函数完成时将等待组值减1
				defer swg.Done()
				if isOnlyWaf {
					_, req, url, _ := network.Reqdata(url, proxy)
					processWaf(req, url, waf,proxy)

				} else {


					req, url := processFinger(url, proxy, finger, nuclei,filterStatus)

					if isWaf && (req != nil) {
						processWaf(req, url, waf,proxy)
					}

				}

			}(domains[i])

		}

	} else if fofa {
		var fofa = new(utils.FofaResp)
		fofa.Crawler(fofaFile, fofaRule, fofaSize, fofaPage, isCloud)

	}

	swg.Wait()

}

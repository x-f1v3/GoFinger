package json_new

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"scan/utils/utils"
)

type Finger struct {
	Product string `json:"product"`
	Rules   [][]struct {
		Match   string `json:"match"`
		Content string `json:"content"`
	} `json:"rules"`
}

type FetchResult struct {
	Url          string
	Content      []byte
	Headers      http.Header
	HeaderString string
	Status       string
	Certs        []byte
	FaviconHash  string
	OtherString  string
}

//返回切片数据
func Parse(filename string, silent bool) ([]Finger, error) {

	Json, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var dataArray []Finger
	err = json.Unmarshal(Json, &dataArray)

	if err != nil {

		return nil, err
	}

	if !silent {
		load_msg := fmt.Sprintf("成功加载 %s 指纹库,共加载指纹 %d", filename, len(dataArray))
		utils.Info(load_msg)
	}

	return dataArray, nil
}

func (x FetchResult) IsEmpty() bool {
	return reflect.DeepEqual(x, FetchResult{})
}

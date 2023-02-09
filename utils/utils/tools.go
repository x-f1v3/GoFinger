package utils

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
)

func StandBase64(braw []byte) string {
	bckd := base64.StdEncoding.EncodeToString(braw)
	var buffer bytes.Buffer
	for i := 0; i < len(bckd); i++ {
		ch := bckd[i]
		buffer.WriteByte(ch)
		if (i+1)%76 == 0 {
			buffer.WriteByte('\n')
		}
	}
	buffer.WriteByte('\n')
	return string(buffer.Bytes())

}

func IsContain(items []string, item string) bool {

	for _, eachItem := range items {
		if strings.Contains(item, eachItem) {
			return true
		}
	}
	return false
}

func ByteToGBK(strBuf []byte) []byte {
	if isUtf8(strBuf) {
		if GBKBuf, err := simplifiedchinese.GBK.NewEncoder().Bytes(strBuf); err == nil {
			if isUtf8(GBKBuf) == false {
				return GBKBuf
			}
		}
		if GB18030Buf, err := simplifiedchinese.GB18030.NewEncoder().Bytes(strBuf); err == nil {
			if isUtf8(GB18030Buf) == false {
				return GB18030Buf
			}
		}
		//if HZGB2312Buf, err := simplifiedchinese.HZGB2312.NewEncoder().Bytes(strBuf); err == nil {
		//	if isUtf8(HZGB2312Buf) == false {
		//		return HZGB2312Buf
		//	}
		//}
		return strBuf
	} else {
		return strBuf
	}
}

func ByteToUTF8(strBuf []byte) []byte {
	if isUtf8(strBuf) {
		return strBuf
	} else {
		if GBKBuf, err := simplifiedchinese.GBK.NewDecoder().Bytes(strBuf); err == nil {
			if isUtf8(GBKBuf) == true {
				return GBKBuf
			}
		}
		if GB18030Buf, err := simplifiedchinese.GB18030.NewDecoder().Bytes(strBuf); err == nil {
			if isUtf8(GB18030Buf) == true {
				return GB18030Buf
			}
		}
		//if HZGB2312Buf, err := simplifiedchinese.HZGB2312.NewDecoder().Bytes(strBuf); err == nil {
		//	fmt.Println("3")
		//	if isUtf8(HZGB2312Buf) == true {
		//		return HZGB2312Buf
		//	}
		//}
		return strBuf
	}
}

func isUtf8(buf []byte) bool {
	return utf8.Valid(buf)
}

func RunCommand(path, name string, arg ...string) (msg string, err error) {
	cmd := exec.Command(name, arg...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	fmt.Println(cmd.Args)
	if err != nil {
		msg = fmt.Sprint(err) + ": " + stderr.String()
		err = errors.New(msg)
		fmt.Println("err", err.Error(), "cmd", cmd.Args)
	}
	fmt.Println(out.String())
	return
}

func NucleiRun(url string, product string) {
	if product == "" {
		return
	}
	if strings.Contains(product, "Weblogic") {
		RunCommand("", "/bin/bash", "-c", "nuclei -u "+url+" -retries 3 -tags  weblogic ")

	} else if strings.Contains(product, "ThinkPHP") {
		RunCommand("", "/bin/bash", "-c", "nuclei -u "+url+" -retries 3 -tags  thinkphp ")

	}

}

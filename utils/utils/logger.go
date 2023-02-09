package utils

import (
	"fmt"
	"time"
)

var (
	conf    int    = 0
	bg      int    = 0
	timeNow string = setColor("["+time.Now().Format("15:04:05")+"] ", TextMagenta)
)

const (
	TextBlack = iota + 30
	TextRed
	TextGreen
	TextYellow
	TextBlue
	TextMagenta
	TextCyan
	TextWhite
)

func Info(msg string) {

	var infoMsg string = setColor("[INFO] ", TextGreen)
	infoMsg += msg
	fmt.Println(timeNow + infoMsg)

}

func Success(url string, title string, products string, silent bool) {

	var successMsg string = setColor("[+] ", TextBlue)
	successMsg += setColor(url+" ", TextWhite)
	successMsg += setColor(title+" ", TextYellow)
	successMsg += setColor(products, TextRed)
	if silent {

		fmt.Println(url + "-----------" + products)

	} else {
		fmt.Println(timeNow + successMsg)
	}

}

func setColor(msg string, color int) string {
	return fmt.Sprintf("%c[%d;%d;%dm%s%c[0m", 0x1B, conf, bg, color, msg, 0x1B)
}

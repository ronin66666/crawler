package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"net/http"
)

func main() {
	url := "https://www.thepaper.cn/"
	body, err := Fetch(url)
	if err != nil {
		fmt.Printf("read content failed: %v\n", err)
		return
	}

	//使用 goquery 库， 通过CSS 选择器来匹配标签中的文本
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		fmt.Println("read content failed: %v", err)
	}
	doc.Find("div.news_li h2 a[target=_blank]").Each(func(i int, selection *goquery.Selection) {
		//获取匹配标签中的文本
		title := selection.Text()
		fmt.Printf("Review %d: %s\n", i, title)
	})
}

// 其获取网页的内容，检测网页的字符编码并将文本统一转换为 UTF-8 格式。
func Fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Error status code:%v", resp.StatusCode))
	}
	bodyReader := bufio.NewReader(resp.Body)
	e := DeterminEncoding(bodyReader)

	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	return io.ReadAll(utf8Reader)

}

// 读取编码方式
func DeterminEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)

	if err != nil {
		fmt.Println("fetch error:%v", err)
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")

	return e
}

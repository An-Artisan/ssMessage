package main

import (
	"log"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"fmt"
	"io"
	"bufio"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
	"io/ioutil"
	"golang.org/x/text/encoding"
	"bytes"
	"os"
)

func main() {

	resp, err := http.Get("http://www.zhenai.com/zhenghun")


	os.Exit(0)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error:Status Code", resp.StatusCode)
		return
	}
	e:= determineEncoding(resp.Body)
	// 内容转换成utf8
	utf8Reader := transform.NewReader(resp.Body,e.NewDecoder())

	all, err := ioutil.ReadAll(utf8Reader)
	if err != nil {
		panic(err)
	}
	// 把bytes 数据转换成io.Reader  NewReader返回bytes.Reader，但是它实现了io.Reader接口，就是一个类型
	allContent := bytes.NewReader(all)
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(allContent)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("#cityList>dd>a").Each(func(i int, contentSelection *goquery.Selection) {
		url,ok := contentSelection.Attr("href")
		city := contentSelection.Text()
		if !ok {
			panic("exist href...")
			return
		}

		log.Printf("第%d个%v的url为：%v",i+1,city,url)
	})
}
//  判断传入的io内容字符集，返回字符集
func determineEncoding(r io.Reader) encoding.Encoding{
	bytes ,err := bufio.NewReader(r).Peek(1024)
	if err != nil {
		panic(err)
	}
	e,_,_ := charset.DetermineEncoding(bytes,"")
	return e
}
func getCookieHandler(w http.ResponseWriter, r *http.Request) {
	h := r.Header["Cookie"]
	fmt.Fprintln(w, h)
}
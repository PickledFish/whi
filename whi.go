package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const StopKey = 3
const ReportHtml = "<!DOCTYPE html><html lang='en'><head><meta charset='UTF-8'><title>报告</title></head><body><table border='1px'><tr><td>No.</td><td>target</td><td>title</td><td>status</td><td>error</td></tr>>replace<</table></body></html>"

func crawl(target string) map[string]string {
	fmt.Println("[Get]", target)
	resp, err := http.Get(target)
	if err != nil {
		return map[string]string{"target": target, "err": err.Error()}
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return map[string]string{"target": target, "err": err.Error(), "status": string(resp.StatusCode)}
	}
	result := make(map[string]string)
	doc.Find("title").Each(func(i int, selection *goquery.Selection) {
		result = map[string]string{"target": target, "status": strconv.FormatInt(int64(resp.StatusCode), 10), "title": selection.Text()}
	})
	return result
}

func main() {
	var (
		output string
		input  string
	)
	chTasks := make(chan string)
	chResult := make(chan interface{})
	flag.StringVar(&input, "input", "", "输入")
	flag.StringVar(&output, "output", "./result.html", "输出")
	flag.Parse()

	_, err := os.Stat(input)
	if err != nil {
		fmt.Println("Input File Err:", err.Error())
		return
	}

	data, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatal(err)
		return
	}

	split := strings.Split(string(data), "\n")
	go func() {
		for _, item := range split {
			chTasks <- strings.Trim(item, "\r")
		}
		close(chTasks)
	}()

	maxGo := 5
	if len(split) < 5 {
		maxGo = len(split)
	}

	for i := 0; i < maxGo; i++ {
		go func() {
			for {
				target, ok := <-chTasks
				if !ok {
					chResult <- StopKey
					break
				} else {
					chResult <- crawl(target)
				}
			}
		}()
	}

	func(c chan interface{}) {
		count := 0
		co := 1
		var replacePort []string
	ForLabel:
		for {
			select {
			case item := <-c:
				result, ok := item.(map[string]string)
				if ok {
					replacePort = append(replacePort, "<tr><td>"+strconv.Itoa(co)+"</td><td><a href='"+result["target"]+"' target='_blank'>"+result["target"]+"</a></td><td>"+result["title"]+"</td><td>"+result["status"]+"</td><td>"+result["err"]+"</td></tr>")
					co += 1
				}
				_, ok = item.(int)
				if ok {
					count += 1
					if count == maxGo {
						html := strings.Replace(ReportHtml, ">replace<", strings.Join(replacePort, "\n"), 1)
						err := ioutil.WriteFile(output, []byte(html), 0644)
						if err != nil {
							fmt.Println(err.Error())
						}
						break ForLabel
					}
				}
			default:
			}
		}
	}(chResult)
}
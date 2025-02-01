package reptile

import (
	"github.com/dubbogo/net/html"
	"net/http"
	"strings"
)

type CalendarInfo struct {
	Link       string
	Year       string //以学年的开始作为标记
	PDFLink    string
	ImageLinks []string
}

// 不是很成功的设计但是方便
const BASEURL = "https://jwc.ccnu.edu.cn"

type Reptile interface {
	GetCalendarLink() ([]CalendarInfo, error)
	FetchPDFOrImageLinksFromPage(url string) (string, []string, error)
}

type reptile struct{}

func NewReptile() Reptile {
	return &reptile{}
}

func (r *reptile) GetCalendarLink() ([]CalendarInfo, error) {
	//使用配置文件中的url
	url := "https://jwc.ccnu.edu.cn/index/hdxl.htm"
	//获取最初的页面
	resp, err := http.Get(url)
	if err != nil {

		return nil, err
	}
	defer resp.Body.Close()
	//解析页面
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	//创建存储获取数据的变量
	var calendarInfos []CalendarInfo
	//定义一个解析函数,有点奇妙的写法直接将数据给到calendarInfos,奇妙函数加奇妙递归
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "li" {
			for _, attr := range n.Attr {
				if attr.Key == "id" && strings.HasPrefix(attr.Val, "line_u7_") {
					var link, semester string
					for c := n.FirstChild; c != nil; c = c.NextSibling {
						if c.Type == html.ElementNode && c.Data == "a" {
							for _, aAttr := range c.Attr {
								if aAttr.Key == "href" {
									link = aAttr.Val
								}
							}
							//获取日历的开始学年
							semester = strings.TrimSpace(c.FirstChild.Data)[:4]
						}
					}
					if link != "" && semester != "" {
						//获取完成的日历链接
						fullLink := BASEURL + link[2:]
						calendarInfos = append(calendarInfos, CalendarInfo{Link: fullLink, Year: semester})
					}
				}
			}
		}
		//递归查询直到找到需要的位置为止
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)
	return calendarInfos, nil
}

// FetchPDFOrImageLinksFromPage 从给定的CalendarLinkURL中提取PDF链接和所有图片链接
func (r *reptile) FetchPDFOrImageLinksFromPage(url string) (string, []string, error) {
	// 发送HTTP GET请求以获取网页内容
	resp, err := http.Get(url)
	if err != nil {
		// 如果请求失败，返回错误
		return "", nil, err
	}
	// 确保函数结束时关闭HTTP响应体
	defer resp.Body.Close()

	// 解析HTML文档
	doc, err := html.Parse(resp.Body)
	if err != nil {
		// 如果解析失败，返回错误
		return "", nil, err
	}

	// 初始化用于存储PDF链接和图片链接的变量
	var pdfLink string
	var imageLinks []string

	// 定义一个递归遍历函数，用于遍历HTML节点树
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			// 如果当前节点是一个<a>标签，检查其是否包含PDF链接
			if n.Data == "a" {
				for _, attr := range n.Attr {
					if attr.Key == "href" && strings.HasSuffix(attr.Val, ".pdf") {
						pdfLink = BASEURL + attr.Val
					}
				}
			}
			// 找到 <div> 标签，且 class 属性为 "v_news_content"
			if n.Type == html.ElementNode && n.Data == "div" {
				for _, attr := range n.Attr {
					if attr.Key == "class" && attr.Val == "v_news_content" {
						// 在这个 <div> 块中查找所有 <img> 标签的 src 属性
						findImages(n, &imageLinks)
					}
				}
			}
		}
		// 递归遍历子节点
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	// 开始从文档的根节点遍历HTML节点树
	traverse(doc)
	// 返回找到的PDF链接和图片链接
	return pdfLink, imageLinks, nil
}

// findImages 在指定的节点树中递归查找所有 <img> 标签，并获取其 src 属性的值
func findImages(n *html.Node, imageLinks *[]string) {
	if n.Type == html.ElementNode && n.Data == "img" {
		for _, attr := range n.Attr {
			if attr.Key == "src" {
				*imageLinks = append(*imageLinks, BASEURL+attr.Val)
			}
		}
	}
	// 递归遍历子节点
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findImages(c, imageLinks)
	}
}

package apigen

import (
	"bufio"
	"strings"
	"fmt"
	"strconv"
	"io"
	"errors"
)

const (
	// 未知类型
	TYPE_UNKNOWN int = 0;
	// Api描述名称
	TYPE_NAME int = 1
	// Api描述Item
	TYPE_ITEM int = 2
)

// Aip描述
type ApiDesc struct {
	name string
	item []string
	url string
	code int
	desc string
}

func NewApiDesc(name string, item []string) *ApiDesc{
	api := ApiDesc{}
	api.name = name
	api.item = item

	ensureApiDescWithItem(&api, item)

	return &api
}

func ensureApiDesc(api *ApiDesc) *ApiDesc{
	item := api.item
	api.url = item[0]
	itemLen := len(item)
	if itemLen >= 3{
		api.code, _ = strconv.Atoi(item[1])
		api.desc = item[2]
	}else if itemLen >= 2{
		api.desc = item[1]
	}
	return api
}

func ensureApiDescWithItem(api *ApiDesc, item []string) *ApiDesc{
	api.url = item[0]
	itemLen := len(item)
	if itemLen >= 3{
		api.code, _ = strconv.Atoi(item[1])
		api.desc = item[2]
	}else if itemLen >= 2{
		api.desc = item[1]
	}
	return api
}

func (a *ApiDesc) GetName() string{
	return a.name
}

func (a *ApiDesc) GetUrl() string{
	return a.url
}

func (a *ApiDesc) GetCode() int{
	return a.code;
}

func (a *ApiDesc) GetDesc() string{
	return a.desc
}

func (a *ApiDesc) HasDesc() bool{
	return len(a.desc) > 0
}

// 是否是绝对Url
func (a *ApiDesc) IsAbsUrl() bool{
	var indexHttp, indexHttps int
	indexHttp= strings.Index(a.url, "http://")
	if indexHttp == 0{
		return true;
	}else if indexHttps = strings.Index(a.url, "http://"); indexHttps == 0{
		return true
	}
	return false;
}

// 解析
func ParserReader(rd io.Reader) []*ApiDesc{

	buffReader := bufio.NewReader(rd)

	var baseLine string
	apis := make([]*ApiDesc, 0, 10)

	var lastName string

	var lastApi *ApiDesc
	for buffLine, isPrefix, err := buffReader.ReadLine(); err == nil; buffLine, isPrefix, err = buffReader.ReadLine(){
		line := string(buffLine)
		if isPrefix{
			baseLine = baseLine + line
			continue
		}else{
			if len(baseLine) > 0{
				line = baseLine + line
			}
			baseLine = ""
		}
		element, eType, err := parserLine(line)

		if err == nil{
			switch eType{
			case TYPE_NAME:
				if lastApi != nil{
					ensureApiDesc(lastApi)
				}
				lastName = element
				lastApi = &ApiDesc{}
				lastApi.name = lastName
				apis = append(apis, lastApi)

			case TYPE_ITEM:
				if lastName != ""{
					//fmt.Println("--item--", element)
					lastApi.item = append(lastApi.item, element)
					//fmt.Println("--items--", lastApi.item)
				}else{
					fmt.Println("--ignore item--", element)
				}
			}
		}else{

		}
	}

	if lastApi != nil{
		ensureApiDesc(lastApi)
	}

	//	for _, v := range(apis){
	//		fmt.Println("apigen", v)
	//	}
	return apis
}

// 解析string
func Parser(descStr string) []*ApiDesc{
	strReader := strings.NewReader(descStr);

	apis := ParserReader(strReader)
	return apis
}

// 解析Line 得到结果和 Type
func parserLine(line string) (string, int, error){
	charLen := len(line)
	if charLen <= 0{
		e := errors.New("the line is empty")
		return "", TYPE_UNKNOWN, e;
	}
	firstNoSpace := getFirstNoSpaceIndex(line);
	//fmt.Println("first:", firstNoSpace)

	if firstNoSpace == 0{
		// 解析名称
		name, err := parserName(line)
		//fmt.Println("name is", name, err)
		return name, TYPE_NAME, err
	}else if firstNoSpace > 0 && firstNoSpace < charLen - 1{
		// 解析列表
		item, err := parserItem(line, firstNoSpace)
		return item, TYPE_ITEM, err
	}
	return "", TYPE_UNKNOWN, nil
}

// 解析Api描述的名称项
func parserName(line string) (string, error){
	end := strings.IndexByte(line, byte(":"[0]));
	if end > -1{
		return line[:end], nil
	}

	panic(": no found in line, it's not a Name")
}

// 解析Api描述的Item项
func parserItem(line string, firstNoSpace int)(string, error){
	if line[firstNoSpace] != "-"[0]{
		e := errors.New("item not start with '-'," + string(line[firstNoSpace]))
		return "", e
	}
	line = line[firstNoSpace + 1 : ]
	first := getFirstNoSpaceIndex(line)
	if first > -1{
		line = line[first:]
		line = strings.TrimSpace(line)
		return line, nil
	}else{
		e := errors.New("item is empty")
		return "", e
	}
}

// 获取第一个不是空白字符的index
func getFirstNoSpaceIndex(line string) int{
//	length := len(line)
	var first int = -1
	for i, v := range(line){
		if !IsSpace((byte)(v)) {
			first = i
			break
		}
	}
	return first
}

// 是否是空白字符
func IsSpace(c byte) bool{
	return c == '\n' || c == '\t' || c == '\f' || c == '\r' || c == ' ';
}

package api

import (
	"bufio"
	"strings"
	"fmt"
	"strconv"
	"io"
	"errors"
)

const (
	TYPE_UNKNOWN int = 0;
	TYPE_NAME int = 1
	TYPE_ITEM int = 2
)

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

func ParserReader(rd io.Reader) []*ApiDesc{

	buffReader := bufio.NewReader(rd)

	var baseLine string
	apis := make([]*ApiDesc, 0, 10)

	var lastName string

	var lastApi *ApiDesc
	for buffLine, isPerfix, err := buffReader.ReadLine(); err == nil; buffLine, isPerfix, err = buffReader.ReadLine(){
		line := string(buffLine)
		if isPerfix{
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
	//		fmt.Println("api", v)
	//	}
	return apis
}

// 解析
func Parser(yaml string) []*ApiDesc{
	strReader := strings.NewReader(yaml);

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

//		if err == nil{
//			if lastName != ""{
//				fmt.Println("----item is:", item)
//			}else{
//				fmt.Println("--ignore--item is:", item)
//			}
//		}else{
//			fmt.Println("parser item err,", line, err)
//		}

	}
	return "", TYPE_UNKNOWN, nil
}

func parserName(line string) (string, error){
	end := strings.IndexByte(line, byte(":"[0]));
	if end > -1{
		return line[:end], nil
	}

	panic(": no found in line, it's not a Name")
}

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

func getFirstNoSpaceIndex(line string) int{
//	length := len(line)
	var first int = -1
	for i, v := range(line){
		if !isSpace((byte)(v)) {
			first = i
			break
		}
	}
	return first
}

func isSpace(c byte) bool{
	return c == '\n' || c == '\t' || c == '\f' || c == '\r' || c == ' ';
}

package main

import (
	"api"
	"fmt"
	"os"
)

const (
	Test_API = `CANCEL_FOLLOW:
    - account/user/unfollow"
    - 7
    - 取消关注. /{followId}`

	NAME_HOST = "HOST"
)

func main(){
	// /** 用户关注列表. /{id}/{page}/{size}*/
	// public static final String API_GET_FOLLOWED_LIST = HOST + "api/account/user/followList";
	// public static final int CODE_GET_FOLLOWERD_LIST = 4;
	var file = "/home/holmes/Documents/Jie/API定义自动生成.txt"

	apiFile, _:= os.Open(file)
	defer apiFile.Close()

	apis := api.ParserReader(apiFile)

	var host string
	var baseCode int

	var code int

	fmt.Println("==========print=========")

	for i, v := range(apis){
		//fmt.Println("index", i)
		_ = i
		if v.GetName() == NAME_HOST{
			host = v.GetUrl()
			baseCode = v.GetCode()
			code = baseCode
		}else{
			// 输出
			printApi(v, code, true)
		}

		code ++
	}

	_ = host

//	fmt.Println("end paraser", apis)
}

func printApi(apiDesc *api.ApiDesc, code int, host bool){
	var desc string
	var apiStr string
	var codeStr string

	_ = host

	if apiDesc.HasDesc(){
		desc = fmt.Sprintf("/** %s */", apiDesc.GetDesc())
		fmt.Println(desc)
	}
	apiStr = fmt.Sprintf("public static final String API_%s = HOST + \"%s\";", apiDesc.GetName(), apiDesc.GetUrl())
	codeStr = fmt.Sprintf("public static final String CODE_%s = %d;", apiDesc.GetName(), code)

	fmt.Println(apiStr)
	fmt.Println(codeStr)
	fmt.Println("")
}

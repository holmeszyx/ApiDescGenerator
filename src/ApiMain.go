package main

import (
	"fmt"
	"os"
	apigen "./apigen"
	"text/template"
	"log"
	"path/filepath"
	"bufio"
)

const (
	Test_API = `CANCEL_FOLLOW:
    - account/user/unfollow"
    - 7
    - 取消关注. /{followId}`

	Test_API_File = "/home/holmes/IdeaProjects/ApiDescGenerator/templates/api_desc.txt"
	Test_Temp_File = "/home/holmes/IdeaProjects/ApiDescGenerator/templates/Net.java.tmpl"

	API_FILE_NAME = "api_desc.txt"
	TEMP_FILE_NAME = "Net.java.tmpl"

	DES_FILE_NAME = "Net.java"

	NAME_HOST = "HOST"
)

type ApiDescTemp struct {
	Name string
	Url string
	Code int
	Desc string
	HasDesc bool
	IsAbsUrl bool
	Host string
	IsHost bool
}

// 文件是否存在
func isFileExist(file string) bool{
	if _, err := os.Stat(file); err == nil{
		return true
	}
	return false
}

func findFiles(api string, temp string)(string, string){
	// 当前目录下的 API_FILE_NAME
	if !isFileExist(api){
		api = ""
	}
	// 当前目录下 templates/TEMP_FILE_NAME
	if !isFileExist(temp){
		temp = ""
	}
	return api, temp
}

func getCurrentDir() string{
	currentDir := filepath.Dir(os.Args[0])
	var err error
	currentDir, err = filepath.Abs(currentDir)
	if err != nil{
		log.Fatal(err)
	}
	return currentDir
}

func main(){
	// /** 用户关注列表. /{id}/{page}/{size}*/
	// public static final String API_GET_FOLLOWED_LIST = HOST + "apigen/account/user/followList";
	// public static final int CODE_GET_FOLLOWERD_LIST = 4;

	var apiFilePath, tempFilePath, descFilePath string
	currentDir := getCurrentDir()
	apiFilePath = currentDir + "/" + API_FILE_NAME
	tempFilePath = currentDir + "/templates/" + TEMP_FILE_NAME
	descFilePath = currentDir + "/" + DES_FILE_NAME

	checkApi, checkTemp := findFiles(apiFilePath, tempFilePath)
	if checkApi == ""{
		fmt.Println("Api描述文件不存在", apiFilePath)
		os.Exit(1)
	}

	if checkTemp == ""{
		fmt.Println("模板不存在", tempFilePath)
		os.Exit(2)
	}

	//apiFile = Test_API_File
	//tempFile = Test_Temp_File

	apiFile, _:= os.Open(apiFilePath)
	defer apiFile.Close()

	apis := apigen.ParserReader(apiFile)

	var apiHost *apigen.ApiDesc
	var baseCode int

	var code int

	var apiImpList = make([]*ApiDescTemp, 0, 10)

	//fmt.Println("========== print Start =========")

	var apiImpl *ApiDescTemp
	for i, v := range(apis){
		//fmt.Println("index", i)
		_ = i
		if v.GetName() == NAME_HOST{
			baseCode = v.GetCode()
			apiHost = v
			code = baseCode
			apiImpl = convertApiImpl(v, code, "")
			apiImpl.IsHost = true
		}else{
			// 输出
			//printApi(v, code, true)
			apiImpl = convertApiImpl(v, code, apiHost.GetName())
		}
		apiImpList = append(apiImpList, apiImpl)

		code ++
	}

	//fmt.Println("========== print End=========")

	apiTemp, err := template.ParseFiles(tempFilePath)
	if err != nil{
		log.Fatal(err)
	}


//	fmt.Println("end paraser", apis)
	fmt.Println("\n")
	fmt.Println("write to", descFilePath)
	if isFileExist(descFilePath) {
		os.Remove(descFilePath)
	}
	descFile, err := os.OpenFile(descFilePath, os.O_RDWR | os.O_CREATE, 0666)
	if err != nil{
		log.Fatal(err)
	}
	defer  descFile.Close()

	fileWirter := bufio.NewWriter(descFile)

	ApiDescList := apiImpList
	err = apiTemp.Execute(fileWirter, ApiDescList)
	if err != nil{
		log.Fatal(err)
	}

	fileWirter.Flush()

	fmt.Println("Finished !!")
}

// 转换Api描述成模板对像
func convertApiImpl(apiDesc *apigen.ApiDesc, code int, host string) (apiImpl *ApiDescTemp){
	apiImpl = new(ApiDescTemp)
	if code == -1{
		apiImpl.Code = apiDesc.GetCode()
	}else{
		apiImpl.Code = code
	}
	apiImpl.Name = apiDesc.GetName()
	apiImpl.Url = apiDesc.GetUrl()
	apiImpl.Desc = apiDesc.GetDesc()
	apiImpl.HasDesc = apiDesc.HasDesc()
	apiImpl.IsAbsUrl = apiDesc.IsAbsUrl()
	apiImpl.Host = host

	return
}

func printApi(apiDesc *apigen.ApiDesc, code int, host bool){
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

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
)

/**
 * gitee.com 的Webhook解析
 * 目前Content-Type只有JSON格式
 *
 */
func gitee(w http.ResponseWriter, request *http.Request) {
	contentType := request.Header.Get("Content-Type")
	if contentType == "application/json" {
		json := ParseGitEE(request)
		w.Write([]byte(json))
	} else {
		w.Write([]byte(`Hello GitEE`))
	}
}

/**
 * coding.net 的Webhook解析
 * 暂时不解析ContentType
 *
 */
func coding(w http.ResponseWriter, request *http.Request) {
	json := ParseCoding(request)
	w.Write([]byte(json))
}

func gogs(w http.ResponseWriter, request *http.Request) {
	json := ParseGogs(request)
	w.Write([]byte(json))
}

func index(w http.ResponseWriter, request *http.Request) {
	// 处理UserAgent，判断是Coding还是Gitee
	// X-Coding-Event
	// X-Gitee-Event
	// X-Gogs-Event
	json := `Hello`
	if request.Header.Get(`X-Coding-Event`) != `` {
		json = ParseCoding(request)
	} else if request.Header.Get(`X-Gitee-Event`) != `` {
		json = ParseGitEE(request)
	} else if request.Header.Get(`X-Gogs-Event`) != `` {
		json = ParseGogs(request)
	}

	w.Write([]byte(json))
}

func hook(owner string, projectName string, branch string, pwd string) string {
	// 读取文件
	filename := fmt.Sprintf("%v.%v.%v.json", owner, projectName, branch)
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Print(`无法读取文件:` + filename + `:` + err.Error())
		return "无法执行CI"
	}

	var fileJSON map[string]interface{}
	err = json.Unmarshal(b, &fileJSON)
	if err != nil {
		log.Print(`JSON解析错误:` + err.Error())
		return "CI文件解析错误"
	}

	//filePwd := fmt.Sprint(fileJSON["password"])
	workspace := fmt.Sprint(fileJSON["path"])

	// 校验密码
	//if pwd != `` {
	//	if pwd != filePwd {
	//		log.Print(`密码校验错误:` + pwd + `:正确密码:` + filePwd + `:` + err.Error())
	//		return "凭证校验异常"
	//	}
	//}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command(`./git.bat `, workspace, branch)
	} else {
		// 执行Shell 命令
		c := `./git.sh ` + workspace + ` ` + branch
		cmd = exec.Command("sh", "-c", c)
	}

	err = cmd.Start() // 该操作不阻塞
	if err != nil {
		log.Printf(`Shell执行异常: %v, err: %v`, cmd, err)
		return "任务执行异常"
	}
	return "The Job Done!"
}

/**
 *
 * 解析Coding.net的数据
 *
 */
func ParseCoding(request *http.Request) string {

	event := request.Header.Get(`X-Coding-Event`)

	if event == `ping` {
		return "这个Coding不简单"
	}

	result, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Print(`请求参数无法获取:` + err.Error())
		return "未获取到数据"
	}

	var data map[string]interface{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		log.Print(`JSON解析出错:` + err.Error())
		return "未获取到数据包"
	}

	hookID := fmt.Sprint(data["hook_id"])
	if err == nil {
		return hookID + `一切正常`
	}

	// 分支名称
	ref := fmt.Sprint(data["ref"])
	branchs := strings.Split(ref, `/`)
	branch := branchs[2]

	// 获取拥有者
	owner := getCodingOwner(data)
	projectName := getCodingProjectName(data)

	// 获取密码
	pwd := fmt.Sprint(data["token"])

	return hook(owner, projectName, branch, pwd)

}

func getCodingOwner(jsonData map[string]interface{}) string {
	repo := jsonData["repository"]
	respoM, ok := repo.(map[string]interface{})
	if !ok {
		return ""
	}

	owner := respoM["owner"]
	ownerM, ok := owner.(map[string]interface{})
	if !ok {
		return ""
	}

	return fmt.Sprint(ownerM["name"])
}

func getCodingProjectName(jsonData map[string]interface{}) string {
	repo := jsonData["repository"]
	respoM, ok := repo.(map[string]interface{})
	if !ok {
		return ""
	}

	return fmt.Sprint(respoM["name"])
}

/**
 * 解析Gitee.com
 *
 */
func ParseGitEE(request *http.Request) string {
	result, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Print(`请求参数无法获取:` + err.Error())
		return "未获取到数据"
	}

	var data map[string]interface{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		log.Print(`JSON解析出错:` + err.Error())
		return "未获取到数据包"
	}

	// 分支名称
	ref := fmt.Sprint(data["ref"])
	branchs := strings.Split(ref, `/`)
	branch := branchs[2]

	// 获取项目名称
	projName := getGiteeProjectName(data)
	projectNameArr := strings.Split(projName, `/`)
	owner := projectNameArr[0]
	projectName := projectNameArr[1]

	// 获取密码
	pwd := fmt.Sprint(data["password"])

	return hook(owner, projectName, branch, pwd)
}

func getGiteeProjectName(jsonData map[string]interface{}) string {
	repo := jsonData["repository"]
	respoM, ok := repo.(map[string]interface{})
	if !ok {
		return ""
	}

	return fmt.Sprint(respoM["path_with_namespace"])
}

func ParseGogs(request *http.Request) string {
	result, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Print(`请求参数无法获取:` + err.Error())
		return "未获取到数据"
	}

	var data map[string]interface{}
	err = json.Unmarshal(result, &data)
	if err != nil {
		log.Print(`JSON解析出错:` + err.Error())
		return "未获取到数据包"
	}

	// 分支名称
	ref := fmt.Sprint(data["ref"])
	branchs := strings.Split(ref, `/`)
	branch := branchs[2]

	// 获取项目名称
	projName := getGogsProjectName(data)
	projectNameArr := strings.Split(projName, `/`)
	owner := projectNameArr[0]
	projectName := projectNameArr[1]

	// 获取密码
	pwd := ``

	return hook(owner, projectName, branch, pwd)
}

func getGogsProjectName(jsonData map[string]interface{}) string {
	repo := jsonData["repository"]
	respoM, ok := repo.(map[string]interface{})
	if !ok {
		return ""
	}

	return fmt.Sprint(respoM["full_name"])
}

package apis

import (
	"UserInsert/models"
	"UserInsert/setting"

	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/bitly/go-simplejson"
)

func GetToken() (tokenx string) {
	// 调用接口
	// 创建各类对象
	var (
		client  *http.Client
		request *http.Request
		resp    *http.Response
		body    []byte
		err     error
	)

	scUrl := fmt.Sprintf("https://%s:10220/uums/auth/token", setting.Conf.SenseCity.SCIP)
	// 这里请注意，使用 InsecureSkipVerify: true 来跳过证书验证
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}}

	tokenUser := models.TokenUser{
		UserName:  setting.Conf.SenseCity.UserName,
		PassWord:  setting.Conf.SenseCity.PassWord,
		GrantType: setting.Conf.SenseCity.GrantType,
	}
	// Marshal
	bytesData, _ := json.Marshal(tokenUser)

	// 获取 request请求
	if request, err = http.NewRequest("POST", scUrl, bytes.NewReader(bytesData)); err != nil {
		log.Println("GetHttpSkip Request Error:", err)
		return
	}

	// 请求格式json
	request.Header.Add("content-type", "application/json")
	if resp, err = client.Do(request); err != nil {
		log.Println("GetHttpSkip Response Error:", err)
		return
	}

	// 延迟关闭
	defer request.Body.Close()

	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Println("ReadAll resp.Body Error:", err)
		return
	}

	//
	js, err := simplejson.NewJson(body)
	if err != nil {
		fmt.Println(err)
	}

	accessToken, err := js.Get("data").Get("accessToken").String()
	if err != nil {
		fmt.Println(err)
	}

	// 延迟关闭
	defer client.CloseIdleConnections()

	return accessToken
}

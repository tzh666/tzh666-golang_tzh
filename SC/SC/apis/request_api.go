package apis

import (
	"UserInsert/models"
	"UserInsert/utils"

	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

func InsertUser(urlx, tokenx string, param *models.UserInfo) (body []byte, userID int64, err error) {
	// 创建各类对象
	var (
		client    *http.Client
		request   *http.Request
		resp      *http.Response
		bytesData []byte
	)

	// 这里请注意，使用 InsecureSkipVerify: true 来跳过证书验证
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}}

	// Marshal---->参数param, 序列化数据给接口当参数传递 bytesData是json数据
	if bytesData, err = json.Marshal(param); err != nil {
		logrus.Error("Marshal param Error:", err)
		return
	}

	// 获取 request请求
	if request, err = http.NewRequest("POST", urlx, bytes.NewReader(bytesData)); err != nil {
		logrus.Error("GetHttpSkip Response Error:", err)
		return
	}

	// 加入 token; Authorization、accessToken看你接口的请求头是什么了
	request.Header.Add("accessToken", GetToken())
	request.Header.Add("content-type", "application/json")

	// client.Do
	if resp, err = client.Do(request); err != nil {
		logrus.Error("GetHttpSkip Response Error:", err)
		return
	}

	// 延迟关闭
	defer resp.Body.Close()
	// 延迟关闭
	defer client.CloseIdleConnections()

	// 返回值
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		logrus.Error("ReadAll resp.Body Error:", err)
		return
	}

	if userID, err = utils.GetUserID(string(body)); err != nil || userID == -1 {
		logrus.Error("获取用户ID失败: ", err)
		logrus.Info("userid is: ", userID)
		return
	}
	return body, userID, nil
}

// 人像库权限添加
func PortraitLib(urlx, tokenx string, param *models.PortraitLib) (code int, err error) {
	// 创建各类对象
	var (
		client    *http.Client
		request   *http.Request
		resp      *http.Response
		bytesData []byte
		body      []byte
	)

	// 这里请注意，使用 InsecureSkipVerify: true 来跳过证书验证
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}}

	// Marshal---->参数param, 序列化数据给接口当参数传递
	if bytesData, err = json.Marshal(param); err != nil {
		logrus.Error("Marshal param Error:", err)
		return
	}

	// 获取 request请求
	if request, err = http.NewRequest("PUT", urlx, bytes.NewReader(bytesData)); err != nil {
		logrus.Error("GetHttpSkip Response Error:", err)
		return
	}

	// 加入 token; Authorization、accessToken看你接口的请求头是什么了
	request.Header.Add("accessToken", GetToken())
	request.Header.Add("content-type", "application/json")

	// client.Do
	if resp, err = client.Do(request); err != nil {
		logrus.Error("GetHttpSkip Response Error:", err)
		return
	}

	// 延迟关闭
	defer resp.Body.Close()

	// 返回值
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		logrus.Error("ReadAll resp.Body Error:", err)
		return
	}
	logrus.Info("PortraitLib 返回值", string(body))

	// 延迟关闭
	defer client.CloseIdleConnections()
	return resp.StatusCode, nil
}

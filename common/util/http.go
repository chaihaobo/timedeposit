/*
 * @Author: Hugo
 * @Date: 2022-05-11 11:17:19
 * @Last Modified by: Hugo
 * @Last Modified time: 2022-05-17 09:14:35
 */
package util

import (
	"io/ioutil"

	"net/http"
	"strings"

	"gitlab.com/hugo.hu/time-deposit-eod-engine/common/constant"
	commonLog "gitlab.com/hugo.hu/time-deposit-eod-engine/common/log"
)

// Get Membu Api Calling Headers
func getMambuHeader() map[string][]string {
	return map[string][]string{
		"Content-Type": {constant.ContentType},
		"Accept":       {constant.Accept},
		"Apikey":       {constant.Apikey},
	}
}

//Post
func HttpPostData(postJsonStr, postUrl string) (string, int, error) {
	headers := getMambuHeader()

	data := strings.NewReader(postJsonStr)
	req, err := http.NewRequest("POST", postUrl, data)
	if err != nil {
		commonLog.Log.Error("Calling Api Error! Url:%v, req:%v, error:%v", postUrl, postJsonStr, err)
		return "exception!", constant.HttpStatusCodeError, err
	}
	req.Header = headers

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		commonLog.Log.Error("Calling Api Error! Url:%v, req:%v, error:%v", postUrl, postJsonStr, err)
		return "exception!", constant.HttpStatusCodeError, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "exception!", constant.HttpStatusCodeError, err
	}
	return string(body), resp.StatusCode, nil
}

//Patch
func HttpPatchData(postJsonStr, postUrl string) (string, int, error) {
	headers := getMambuHeader()

	data := strings.NewReader(postJsonStr)
	req, err := http.NewRequest("PATCH", postUrl, data)
	if err != nil {
		commonLog.Log.Error("Calling Api Error! Url:%v, req:%v, error:%v", postUrl, postJsonStr, err)
		return "exception!", constant.HttpStatusCodeError, err
	}
	req.Header = headers

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		commonLog.Log.Error("Calling Api Error! Url:%v, req:%v, error:%v", postUrl, postJsonStr, err)
		return "exception!", constant.HttpStatusCodeError, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "exception!", constant.HttpStatusCodeError, err
	}
	return string(body), resp.StatusCode, nil
}

//Get
func HttpGetData(getUrlStr string) (string, int, error) {
	headers := getMambuHeader()

	data := strings.NewReader("")
	req, err := http.NewRequest("GET", getUrlStr, data)
	if err != nil {
		commonLog.Log.Error("Calling Get Api Error! Url:%v, error:%v", getUrlStr, err)
		return "exception!", constant.HttpStatusCodeError, err
	}
	req.Header = headers

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		commonLog.Log.Error("Calling Get Api Error! Url:%v, error:%v", getUrlStr, err)
		return "exception!", constant.HttpStatusCodeError, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "exception!", constant.HttpStatusCodeError, err
	}
	return string(body), resp.StatusCode, nil
}

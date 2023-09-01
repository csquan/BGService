package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"unsafe"

	"github.com/sirupsen/logrus"
)

func Post(requestUrl string, bytesData []byte) (ret string, err error) {
	res, err := http.Post(requestUrl,
		"application/json;charset=utf-8", bytes.NewBuffer([]byte(bytesData)))
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	str := (*string)(unsafe.Pointer(&content)) //转化为string,优化内存
	return *str, nil
}

func Get(url string) (ret string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return string(body), nil
}

// CMC GET调用
func GetWithDataHeader(url string, params url.Values, cmc_api_key string) (ret string, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Info(err)
		os.Exit(1)
	}

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", cmc_api_key)
	req.URL.RawQuery = params.Encode()

	resp, err := client.Do(req)
	if err != nil {
		logrus.Info(err)
		os.Exit(1)
	}
	fmt.Println(resp.Status)
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respBody))
	return string(respBody), nil
}

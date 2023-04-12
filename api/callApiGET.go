package api

import (
	"io/ioutil"
	"net/http"
)

func CallAPIGET(url string) ([]byte, error) {
	// Tạo request với phương thức GET và đường dẫn truyền vào
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Tạo HTTP client và thực hiện request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Đọc và xử lý phản hồi từ server
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respData, nil
}

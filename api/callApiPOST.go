package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func CallAPIPOST(url string, data interface{}) ([]byte, error) {
	// Chuẩn bị dữ liệu gửi đi dưới dạng JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Tạo request với phương thức POST và đường dẫn truyền vào
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// Cấu hình request header
	req.Header.Set("Content-Type", "application/json")

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

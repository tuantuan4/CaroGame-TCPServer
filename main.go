package main

import (
	"log"
	"net"
)

func main() {
	// url := "http://localhost:8080/v1/users/login"

	// // Dữ liệu gửi đi
	// data := map[string]interface{}{
	// 	"username": "tuanuet1",
	// 	"password": "123456",
	// }

	// // Gọi hàm CallAPIPOST và xử lý phản hồi từ server
	// respData, err := api.CallAPIPOST(url, data)
	// if err != nil {
	// 	fmt.Println("Lỗi khi gọi RESTful API:", err)
	// 	return
	// }

	// // Xử lý respData, đây là dữ liệu phản hồi từ server
	// fmt.Println("Phản hồi từ server:", string(respData))
	// url := "http://localhost:8080/v1/users"

	// // Gọi hàm CallAPIGET và xử lý phản hồi từ server
	// respData, err := api.CallAPIGET(url)
	// if err != nil {
	// 	fmt.Println("Lỗi khi gọi RESTful API:", err)
	// 	return
	// }

	// // Xử lý respData, đây là dữ liệu phản hồi từ server
	// fmt.Println("Phản hồi từ server:", string(respData))

	s := newServer()
	go s.run()

	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}

	defer listener.Close()
	log.Printf("server started on :8888")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %s", err.Error())
			continue
		}

		c := s.newClient(conn)
		go c.readInput()
	}
}

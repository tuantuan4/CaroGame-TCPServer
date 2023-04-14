package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"test/api"
	"test/common"
)

type server struct {
	rooms    map[string]*room
	commands chan command
	game     Game
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listRooms(cmd.client)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client)
		case CMD_REGISTER:
			s.register(cmd.client, cmd.args)
		case CMD_LOGIN:
			s.login(cmd.client, cmd.args)
		case CMD_PLAY:
			s.play(cmd.client)
		case CMD_MOVE:
			s.move(cmd.client, cmd.args)
		case CMD_HISTORY:
			s.history(cmd.client)
		case CMD_RATE:
			s.rate(cmd.client, cmd.args)
		case CMD_TIME:
			s.time(cmd.client)
		}
	}
}

func (s *server) newClient(conn net.Conn) *client {
	log.Printf("new client has joined: %s", conn.RemoteAddr().String())

	return &client{
		conn:     conn,
		nick:     "anonymous",
		commands: s.commands,
	}
}

func (s *server) nick(c *client, args []string) {
	if len(args) < 2 {
		c.msg("nick is required. usage: /nick NAME")
		return
	}

	c.nick = args[1]
	c.msg(fmt.Sprintf("all right, I will call you %s", c.nick))
}

func (s *server) join(c *client, args []string) {
	if len(args) < 2 {
		c.msg("room name is required. usage: /join ROOM_NAME")
		return
	}

	roomName := args[1]

	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}
	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrentRoom(c)
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s joined the room", c.nick))

	c.msg(fmt.Sprintf("welcome to %s", roomName))
}

func (s *server) listRooms(c *client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.msg(fmt.Sprintf("available rooms: %s", strings.Join(rooms, ", ")))
}

func (s *server) msg(c *client, args []string) {
	if len(args) < 2 {
		c.msg("message is required, usage: /msg MSG")
		return
	}

	msg := strings.Join(args[1:], " ")
	c.room.broadcast(c, c.nick+": "+msg)
}

func (s *server) quit(c *client) {
	log.Printf("client has left the chat: %s", c.conn.RemoteAddr().String())

	s.quitCurrentRoom(c)

	c.msg("sad to see you go =(")
	c.conn.Close()
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		oldRoom := s.rooms[c.room.name]
		delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())
		oldRoom.broadcast(c, fmt.Sprintf("%s has left the room", c.nick))
	}
}

func (s *server) login(c *client, args []string) {
	if len(args) > 3 || len(args) <= 2 {
		c.msg("input invalid")
		return
	}
	c.nick = args[1]
	//args[1] la username
	//args[2] la password
	url := "http://localhost:8080/v1/users/login"

	// Dữ liệu gửi đi
	data := map[string]interface{}{
		"username": args[1],
		"password": args[2],
	}
	// Gọi hàm CallAPIPOST và xử lý phản hồi từ server
	respData, err := api.CallAPIPOST(url, data)
	if err != nil {
		fmt.Println("Lỗi khi gọi RESTful API:", err)
		return
	}
	type ResponseData struct {
		Message string `json:"message"`
		Id      int    `json:"ID"`
		Error   string `json:"error"`
	}
	var responseData ResponseData
	err = json.Unmarshal(respData, &responseData)
	if err != nil {
		fmt.Println(err)
		return
	}
	if responseData.Message == "Login success" {
		c.idUser = responseData.Id
		c.msg("Login success")
		c.msg(fmt.Sprintf("all right, I will call you %s", c.nick))
	} else if responseData.Error == "Invalid user or password" {
		c.msg("username and password is invalid")
	}
}

func (s *server) register(c *client, args []string) {
	if len(args) > 3 || len(args) <= 2 {
		c.msg("input invalid")
		return
	}
	c.nick = args[1]
	//args[1] la username
	//args[2] la password
	url := "http://localhost:8080/v1/users/register"

	// Dữ liệu gửi đi
	data := map[string]interface{}{
		"username": args[1],
		"password": args[2],
	}

	// Gọi hàm CallAPIPOST và xử lý phản hồi từ server
	respData, err := api.CallAPIPOST(url, data)
	if err != nil {
		fmt.Println("Lỗi khi gọi RESTful API:", err)
		return
	}
	type ResponseData struct {
		Message string `json:"message"`
		Error   string `json:"error"`
	}
	var responseData ResponseData
	err = json.Unmarshal(respData, &responseData)
	if err != nil {
		fmt.Println(err)
		return
	}
	if responseData.Message == "register successful" {
		c.msg("register successful")
	} else if responseData.Error == "User already exists" {
		c.msg("User already exists")
	}
}

func (s *server) play(c *client) {
	s.game.player = append(s.game.player, c.idUser)
	if len(s.game.player) == 2 {
		c.msg("ready to play")
		//goi api create game
		url := "http://localhost:8080/v1/games"

		// Dữ liệu gửi đi
		data := map[string]interface{}{
			"player_id1": s.game.player[0],
			"player_id2": s.game.player[1],
		}

		// Gọi hàm CallAPIPOST và xử lý phản hồi từ server
		respData, err := api.CallAPIPOST(url, data)
		if err != nil {
			fmt.Println("Lỗi khi gọi RESTful API:", err)
			return
		}
		type ResponseData struct {
			Message string `json:"message"`
			Data    int    `json:"data`
		}
		var responseData ResponseData
		err = json.Unmarshal(respData, &responseData)
		if err != nil {
			fmt.Println(err)
			return
		}
		if responseData.Message == "create game is success" {
			s.game.IdGame = responseData.Data
			c.room.broadcastAll(c, "ready to play")
		} else {
			c.msg("Create game is failed")
		}
	} else {
		c.msg("please wait user other")
	}
}

func (s *server) move(c *client, args []string) {
	if len(args) < 2 {
		c.msg("message is required, usage: /msg MSG")
		return
	}
	if len(args) != 3 {
		c.msg("input move invalid, please retype")
		return
	}
	if len(args[1]) != 1 || len(args[2]) != 1 {
		c.msg("move is invalid, please re-enter")
		return
	}
	x := common.StringToInt(args[1])
	y := common.StringToInt(args[2])
	if x >= 3 || x < 0 || y >= 3 || y < 0 {
		c.msg("move is invalid")
		return
	}
	url := "http://localhost:8080/v1/games/AddMove/" + common.IntToString(s.game.IdGame)

	// Dữ liệu gửi đi
	data := map[string]interface{}{
		"player_id":    c.idUser,
		"x_coordinate": x,
		"y_coordinate": y,
	}

	// Gọi hàm CallAPIPOST và xử lý phản hồi từ server
	respData, err := api.CallAPIPOST(url, data)
	if err != nil {
		fmt.Println("Lỗi khi gọi RESTful API:", err)
		return
	}
	type ResponseData struct {
		Message string `json:"message"`
	}
	var responseData ResponseData
	err = json.Unmarshal(respData, &responseData)
	if err != nil {
		fmt.Println(err)
		return
	}
	if responseData.Message == "Add move success" {
		c.msg("Add move successful")
		msg := args[1] + " " + args[2]
		c.room.broadcast(c, c.nick+": "+msg)
		checkWin(c, s)
	} else {
		c.msg("add move failed")
	}

}

func checkWin(c *client, s *server) {
	idGame := common.IntToString(s.game.IdGame)

	url := "http://localhost:8080/v1/games/CheckWin/" + idGame

	// Gọi hàm CallAPIGET và xử lý phản hồi từ server
	respData, err := api.CallAPIGET(url)
	type ResponseData struct {
		Message  string `json:"message"`
		IdWinner int    `json:"IdWinner"`
		IdLoser  int    `json:"IdLoser"`
		Data     string `json:"data"`
	}
	var responseData ResponseData
	err = json.Unmarshal(respData, &responseData)
	if err != nil {
		fmt.Println("Lỗi khi gọi RESTful API:", err)
		return
	}
	if responseData.Message == "2 player draw" {
		c.room.broadcastAll(c, "Draw")
	} else if responseData.Message == "continue play" {
		c.room.broadcastAll(c, "Continue Play")
	} else {
		playerWinner := common.IntToString(responseData.IdWinner)
		url2 := "http://localhost:8080/v1/users/" + playerWinner
		respData1, err := api.CallAPIGET(url2)
		if err != nil {
			fmt.Println("Lỗi khi gọi RESTful API:", err)
			return
		}
		err = json.Unmarshal(respData1, &responseData)
		if err != nil {
			fmt.Println("Lỗi khi gọi RESTful API:", err)
			return
		}
		c.room.broadcastAll(c, "Player "+responseData.Data+" is winner")
	}
}

func (s *server) history(c *client) {

	idPlayer := common.IntToString(c.idUser)

	url := "http://localhost:8080/v1/games/GetHistory/" + idPlayer

	// Gọi hàm CallAPIGET và xử lý phản hồi từ server
	respData, err := api.CallAPIGET(url)
	type ResponseData struct {
		Draw     int    `json:"draw"`
		Lose     int    `json:"lose"`
		Username string `json:"username"`
		Win      int    `json:"win"`
	}
	var responseData ResponseData
	err = json.Unmarshal(respData, &responseData)
	if err != nil {
		fmt.Println("Lỗi khi gọi RESTful API:", err)
		return
	}
	resultHistory := "Username: " + responseData.Username + "............." + "\n" +
		"Win: " + common.IntToString(responseData.Win) + "\n" +
		"Draw: " + common.IntToString(responseData.Draw) + "\n" +
		"Lose: " + common.IntToString(responseData.Lose)
	c.msg(resultHistory)
}

func (s *server) rate(c *client, args []string) {

	// idPlayer := common.IntToString(c.idUser)
	if len(args) < 2 {
		c.msg("message is required, usage: /msg MSG")
		return
	}
	if len(args) != 2 {
		c.msg("input invalid")
		return
	}
	url := "http://localhost:8080/v1/games/history/" + args[1]

	// Gọi hàm CallAPIGET và xử lý phản hồi từ server
	respData, err := api.CallAPIGET(url)
	type ResponseData struct {
		Draw     string `json:"draw"`
		Lose     string `json:"lose"`
		Username string `json:"username"`
		Win      string `json:"win"`
		Sum      int    `json:"sum"`
		Error    string `json:"error`
	}
	var responseData ResponseData
	err = json.Unmarshal(respData, &responseData)
	if err != nil {
		fmt.Println("Lỗi khi gọi RESTful API:", err)
		return
	}
	if responseData.Error == "Username not found" {
		c.msg("Username not found")
		return
	}

	resultHistory :=
		"Win: " + responseData.Win + "\n" +
			"Draw: " + responseData.Draw + "\n" +
			"Lose: " + responseData.Lose + "\n" +
			"Sum: " + common.IntToString(responseData.Sum)
	c.msg(resultHistory)
}

func (s *server) time(c *client) {

	idPlayer := common.IntToString(c.idUser)

	url := "http://localhost:8080/v1/games/time/" + idPlayer

	// Gọi hàm CallAPIGET và xử lý phản hồi từ server
	respData, err := api.CallAPIGET(url)
	type ResponseData struct {
		Data string `json:"data"`
	}
	var responseData ResponseData
	err = json.Unmarshal(respData, &responseData)
	if err != nil {
		fmt.Println("Lỗi khi gọi RESTful API:", err)
		return
	}
	resultTime := responseData.Data
	c.msg(resultTime)
}

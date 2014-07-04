package library

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type Player struct {
	Name   string `json:"name,omitempty"`
	Ping   string `json:"ping,omitempty"`
	Points string `json:"points,omitempty"`
}

type Config struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type Server struct {
	Players       []Player `json:"players,omitempty"`
	Configuration []Config `json:"configs,omitempty"`
}

func GetServerResponse(host string, port string) string {
	

	conn, err := net.Dial("udp", host+":"+port)
	conn.SetDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	var message string = "\377\377\377\377getstatus"
	conn.Write([]byte(message))
	var reply []byte = make([]byte, 1024)
	conn.Read(reply)
	return string(reply)
}

func ParseResponse(response string) Server {

	lines := strings.Split(response, "\n")

	if len(lines) > 1 {
		server_configs := strings.Split(lines[1], "\\")

		configs := make([]Config, 1)
		for i := 1; i < len(server_configs)-1; i = i + 2 {
			key := strings.TrimSpace(server_configs[i])
			value := strings.TrimSpace(server_configs[i+1])

			if len(key) > 0 {
				configs = append(configs, Config{key, value})
			}

		}

		players := make([]Player, 1)
		if len(lines) > 2 {
			for i := 2; i < len(lines)-1; i++ {
				player_data := strings.Split(lines[i], " ")
				points := player_data[0]
				ping := player_data[1]
				name := player_data[2][1 : len(player_data[2])-1]
				players = append(players, Player{name, ping, points})
			}

		}

		return Server{Players: players[1:], Configuration: configs[1:]}

	} else {
		//fmt.Println("Couldn't parse the response.\n")
		return Server{}

	}

}

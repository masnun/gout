package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
    "text/tabwriter"
)

type Player struct {
	name   string
	ping   string
	points string
}

type Config struct {
	key   string
	value string
}

type Server struct {
	players       []Player
	configuration []Config
}

func main() {

	port := flag.String("port", "27960", "Target port on server")
	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		host := args[0]
		
        if len(*port) > 5 {
            fmt.Println("Invalid port number\nTerminating program\n")
            os.Exit(0)
        }

		var response string = GetServerResponse(host, *port)
		var server Server = ParseResponse(response)
		if len(server.configuration) < 1 {
			fmt.Println("Nothing to display\nTerminating program\n")
			os.Exit(0)
		}

        PrintPlayerList(server)

	} else {

		fmt.Println("No host provided\nTerminating program\n")
		os.Exit(0)
	}

}

func PrintPlayerList(server Server) {
    w := new(tabwriter.Writer)
    w.Init(os.Stdout, 0, 8, 0, '\t', 0)
    fmt.Fprintln(w, "name\tping\tpoints")
    for _, player := range server.players {
        fmt.Fprintln(w, player.name + "\t" + player.ping + "\t" + player.points)
    }
    fmt.Fprintln(w)
    w.Flush()
}

func GetServerResponse(host string, port string) string {
    fmt.Println("\nUrban Terror Server Checker")
    fmt.Println("----------------------------")
    fmt.Println("Server: " + host + " // Port: " + port)
    fmt.Println("----------------------------\n\n")

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
				name := player_data[2][2 : len(player_data[2])-1]
				players = append(players, Player{name, ping, points})
			}

		}

		return Server{players: players[1:], configuration: configs[1:]}

	} else {
		//fmt.Println("Couldn't parse the response.\n")
		return Server{}

	}

}

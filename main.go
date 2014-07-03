package main

import (
	"flag"
	"fmt"
	gout "github.com/masnun/gout/library"
	"os"
	"text/tabwriter"
)

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

		var response string = gout.GetServerResponse(host, *port)
		var server gout.Server = gout.ParseResponse(response)
		if len(server.Configuration) < 1 {
			fmt.Println("Nothing to display\nTerminating program\n")
			os.Exit(0)
		}

		PrintPlayerList(server)

	} else {

		fmt.Println("No host provided\nTerminating program\n")
		os.Exit(0)
	}

}

func PrintPlayerList(server gout.Server) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, "name\tping\tpoints")
	for _, player := range server.Players {
		fmt.Fprintln(w, player.Name+"\t"+player.Ping+"\t"+player.Points)
	}
	fmt.Fprintln(w)
	w.Flush()
}

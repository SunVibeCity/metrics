package main

import (
	"./growatt"
	"flag"
	"fmt"
	"os"
)

func main() {
	username := flag.String("u", "", "Username")
	password := flag.String("p", "", "Password")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Example: %s -u john -p secret https://server.growatt.com/\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [http[s]://]hostname[:port]/[path/]\nOptions are:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() == 0 || *username == "" || *password == ""{
		flag.Usage()
		os.Exit(1)
	}

	c := growatt.New(*username, *password, flag.Arg(0))
	l := c.PlantList()
	fmt.Printf("%+v", l)
}

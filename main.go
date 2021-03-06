package main

import (
	"fmt"

	"log"
	"os"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
)

func forwardHandler(X *xgbutil.XUtil, activePos int, size int, cl []xproto.Window) {

	ewmh.ActiveWindowReq(X, cl[(int(activePos)+1)%size])

}

func backwardHandler(X *xgbutil.XUtil, activePos int, size int, cl []xproto.Window) {
	var selected int
	selected = activePos - 1
	if int(activePos)-1 < 0 {
		selected = len(cl) - 1
	}
	ewmh.ActiveWindowReq(X, cl[selected])

}

func listHandler(X *xgbutil.XUtil, clientids []xproto.Window) {

	for _, id := range clientids {
		name, err := ewmh.WmNameGet(X, id)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(name)
	}

}

func contains(s []xproto.Window, e xproto.Window) int {
	for i, a := range s {
		if a == e {
			return i
		}
	}
	return -1
}

func main() {

	if len(os.Args) == 1 {
		fmt.Println("usage: gocycle <command> ")
		fmt.Println(" +   Cycle forward")
		fmt.Println(" -  Cycle backward")
		fmt.Println(" list  List windows")
		return
	}

	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	clientList, err := ewmh.ClientListGet(X)
	if err != nil {
		log.Fatal(err)
	}

	active, err := ewmh.ActiveWindowGet(X)
	if err != nil {
		log.Fatal(err)
	}

	activePos := contains(clientList, active)

	size := len(clientList)

	switch os.Args[1] {
	case "+":
		forwardHandler(X, activePos, size, clientList)

	case "-":
		backwardHandler(X, activePos, size, clientList)

	case "list":
		listHandler(X, clientList)
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}
}

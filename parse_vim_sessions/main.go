package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"

	"github.com/k0kubun/pp"
)

type KittyTab struct {
	ID      int    `json:"id";`
	Focused bool   `json:"is_focused";`
	Title   string `json:"title";`
	File    string
}

type KittyProcess struct {
	ID       int        `json:"id";`
	Focused  bool       `json:"is_focused";`
	WindowID int        `json:"platform_window_id";`
	Tabs     []KittyTab `json:"tabs";`
}

func select_kitty_tab_by_title(title string) {
	cmd := fmt.Sprintf(`kitty @select-window -m=title:%s`, title)
	pp.Println(cmd)
}

func get_vim_file(title string) string {
	return strings.Split(title, ` `)[1]
}

func tab_is_running_vim(title string) bool {
	return (strings.HasPrefix(title, `vi `) || strings.HasPrefix(title, `vim `)) && (len(strings.Split(title, ` `)) == 2)
}

func main() {
	var kp []KittyProcess
	cmd := exec.Command(`kitty`, `@ls`)
	stdout, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	out, err := ioutil.ReadAll(stdout)

	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	uerr := json.Unmarshal(out, &kp)
	if uerr != nil {
		panic(uerr)
	}

	for _, P := range kp {
		for _, T := range P.Tabs {
			if tab_is_running_vim(T.Title) {
				T.File = get_vim_file(T.Title)

				pp.Println(T)
			}
		}
	}
	//pp.Fprintf(os.Stderr, `%s`, kp)
}

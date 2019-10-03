package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/patkruk/gophercises/clitaskmanager/cmd"
	"github.com/patkruk/gophercises/clitaskmanager/db"
)

func main() {
	home, _ := homedir.Dir()
	dbPath := filepath.Join(home, "Dropbox/Code/Go/gophercises/clitaskmanager/tasks.db")

	must(db.Init(dbPath))
	must(cmd.RootCmd.Execute())
}

func must(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

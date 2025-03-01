package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/bsuvonov/gator/internal/config"
	"github.com/bsuvonov/gator/internal/database"

	// "github.com/bsuvonov/gator/internal/database"

	_ "github.com/lib/pq"
)


type state struct {
    conf *config.Config
    db   *database.Queries
}

type command struct {
    name string
    args []string
}


type commands struct {
    handlers    map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
    c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
    handler, ok := c.handlers[cmd.name]
    if !ok {
        return errors.New("command " + cmd.name + " is not registered!")
    }
    return handler(s, cmd);
}



func main() {
	conf, err := config.ReadConfig()
    if err != nil {
        fmt.Println(err)
    }



    cmds := commands{make(map[string]func(*state, command) error)}
    cmds.register("login", handlerLogin)
    cmds.register("register", handlerRegister)

    if len(os.Args) < 2 {
        fmt.Println("error: program must be called with at least one argument.")
        os.Exit(1)
    }

    


    // Create db instance
    db, err := sql.Open("postgres", conf.DbURL)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    dbQueries := database.New(db)
    st := state{&conf, dbQueries}




    // run commands
    cmd_name := os.Args[1]
    args := []string{}
    if len(os.Args) > 2 {
        args = os.Args[2:]
    }
    cmd := command{cmd_name, args}

    err = cmds.run(&st, cmd)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
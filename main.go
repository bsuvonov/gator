package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/bsuvonov/gator/internal/config"
	"github.com/bsuvonov/gator/internal/database"

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
    descriptions map[string]string
}

func (c *commands) register(name string, f func(*state, command) error, description string) {
    c.handlers[name] = f
    c.descriptions[name] = description
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



    cmds := commands{make(map[string]func(*state, command) error), make(map[string]string)}
    cmds.register("login", handlerLogin, "login to existing user")
    cmds.register("register", handlerRegister, "register a new user")
    cmds.register("reset", handlerReset, "reset users and their feeds")
    cmds.register("users", handlerUsers, "list existing users")
    cmds.register("agg", handlerAgg, "I don't know")
    cmds.register("addfeed", handlerAddFeed, "add feed to current user")
    cmds.register("feeds", handlerFeeds, "list all the feeds of current user")

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

    if cmd_name == "help" {
        handlerHelp(cmds)
    } else {

        err = cmds.run(&st, cmd)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
    }
}
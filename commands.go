package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bsuvonov/gator/internal/database"
	"github.com/google/uuid"

	_ "github.com/lib/pq"
)



func handlerLogin(s *state, cmd command) error {
    if len(cmd.args) == 0 || len(cmd.args) > 1 {
        return errors.New("`login` handler expects a single argument - the username")
    }
	name := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("user '%s' does not exist", name)
	}

    s.conf.SetUser(name)
    return nil
}


func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 || len(cmd.args) > 1 {
		return errors.New("`register` handler expects a single argument - the username")
    }

	name := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err == nil {
		return errors.New("user already exists")
	}

	_, err = s.db.CreateUser(context.Background(), database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: name})
	if err != nil {
		return nil
	}
	fmt.Printf("The user '%s' is successfully registered\n", name)
	s.conf.SetUser(name)

	return nil
}


func handlerReset(s *state, cmd command) error {
	return s.db.ResetTable(context.Background())
}

func handlerUsers(s *state, cmd command) error {
	names, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, name := range names {

		if s.conf.CurrentUserName != nil && name == *s.conf.CurrentUserName {
			name += " (current)"
		}
		fmt.Println("* " + name)
	}
	return nil
}


func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://techcrunch.com/feed/")
	fmt.Println(feed)
	return err
}


func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return errors.New("number of arguments to command 'addfeed' must be 4")
	}

	user, err := s.db.GetUser(context.Background(), *s.conf.CurrentUserName)
	if err != nil {
		return err
	}

	feed, err := s.db.InsertFeed(context.Background(), database.InsertFeedParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.args[0], Url: cmd.args[1], UserID: user.ID})
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil
}


func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
        fmt.Printf("* %-*s %-*s %-*s\n", 40, feed.Name, 50, feed.Url, 20, feed.Name_2)
    }

	return nil
}


func handlerHelp(cmds commands) error {
	fmt.Println("Gator is a tool for managing Go source code.")
	fmt.Println("\nUsage")
	fmt.Printf("\n%-*s gator <command> [arguments]\n\n", 7, " ")
	fmt.Println("The commands are:")
	fmt.Println()

	for cmd_name := range cmds.handlers {
		fmt.Printf("%-*s %-*s %-*s\n", 7, " ", 15, cmd_name, 10, cmds.descriptions[cmd_name])
	}
	fmt.Println()
	return nil
}
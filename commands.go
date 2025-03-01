package main

import (
	"context"
	"fmt"
	// "database/sql"
	"errors"
	// "fmt"
	// "os"
	"time"

	// "github.com/bsuvonov/gator/internal/config"
	"github.com/bsuvonov/gator/internal/database"
	"github.com/google/uuid"

	// "github.com/bsuvonov/gator/internal/database"

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



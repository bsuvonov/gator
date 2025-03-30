package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bsuvonov/gator/internal/database"
	"github.com/google/uuid"
	"database/sql"

	_ "github.com/lib/pq"
)



func handlerLogin(s *state, cmd command) error {
    if len(cmd.args) == 0 || len(cmd.args) > 1 {
        return errors.New("`login` command expects a single argument - the username")
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
		return errors.New("`register` command expects a single argument - the username")
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


func scrapeFeeds(s *state) error {
	
	feedToBeFetched, err := s.db.GetNextFeedToFetch(context.Background(), *s.conf.CurrentUserName)
	if err != nil {
		return err
	}

	feedFetched, err := fetchFeed(context.Background(), feedToBeFetched.Url)
	if err != nil {
		return err
	}

	for _, item := range feedFetched.Channel.Items {
		layout := "Mon, 02 Jan 2006 15:04:05 -0700"

		parsedTime, err := time.Parse(layout, item.PubDate)
		if err != nil {
			return err
		}
		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Title: item.Title, Url: item.Link, Description: item.Description, PublishedAt: parsedTime, FeedID: feedToBeFetched.ID})
		if err != nil {
			if err.Error() == "pq: duplicate key value violates unique constraint \"posts_url_key\"" {
				continue
			}
			return err
		}
	}

	s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{ID: feedToBeFetched.ID, LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true}})

	return nil
}


func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("number of arguments to command 'addfeed' must be 1")
	}

	interval, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Println("Collecting feeds every", interval.String())

	ticker := time.NewTicker(interval)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}

	return nil
}


func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
        fmt.Printf("* %-*s %-*s %-*s\n", 20, feed.Name, 50, feed.Url, 0, feed.Name_2)
    }

	return nil
}


func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return errors.New("number of arguments to command 'addfeed' must be 2")
	}

	feed, err := s.db.InsertFeed(context.Background(), database.InsertFeedParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.args[0], Url: cmd.args[1], UserID: user.ID})
	if err != nil {
		return err
	}
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: user.ID , FeedID: feed.ID})
	if err != nil {
		return err
	}

	fmt.Println(feed)
	return nil
}


func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("number of arguments to command 'follow' must be 1")
	}
	feed_id, err := s.db.GetFeedIdByUrl(context.Background(), cmd.args[0])
	if err!=nil {
		return err
	}

	queries, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), UserID: user.ID , FeedID: feed_id})
	if err != nil {
		return err
	}
	for _, query := range queries {
		if query.UserID == user.ID && query.FeedID == feed_id {
			fmt.Printf("* %-*s %-*s\n", 20, query.FeedName, 0, query.UserName)
		}
	}

	return nil
}


func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return errors.New("number of arguments to command 'unfollow' must be 1")
	}
	err := s.db.DeleteFeedFollow(context.Background(), user.ID)
	if err != nil {
		return err
	}
	return nil
}


func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
    return func(s *state, cmd command) error {
        user, err := s.db.GetUser(context.Background(), *s.conf.CurrentUserName)
        if err != nil {
            return err
        }

        return handler(s, cmd, user)
    }
}


func handlerFollowing(s *state, cmd command) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), *s.conf.CurrentUserName)
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("* %-*s %-*s\n", 20, feed.Name, 0, *s.conf.CurrentUserName)
	}
	return nil
}


func handlerHelp(cmds commands) error {
	fmt.Println("Gator is a tool for aggregating RSS feeds.")
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
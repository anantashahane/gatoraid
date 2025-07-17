package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/anantashahane/gatoraid/internal"
	"github.com/anantashahane/gatoraid/internal/config"
	"github.com/anantashahane/gatoraid/internal/database"
	"github.com/google/uuid"
)

type state struct {
	configuration *config.Config
	dbConnection  *database.Queries
}

type command struct {
	command    string
	arguements []string
}

type commands struct {
	commandMap map[string]func(*state, command) error
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) (err error) {
		user, err := s.dbConnection.GetUser(context.Background(), s.configuration.CurrentUserName)
		if err != nil {
			return fmt.Errorf("User %s, not registered/logged in Try `users` command. Error: %s", s.configuration.CurrentUserName, err)
		}
		return handler(s, cmd, user)
	}
}

func handlerLogin(s *state, cmd command) (err error) {
	if len(cmd.arguements) != 1 {
		return fmt.Errorf("Login expects exactly 1 argument, received %v", cmd.arguements)
	}
	if cmd.arguements[0] == s.configuration.CurrentUserName {
		return fmt.Errorf("Already logged in as \"%s\"\n", s.configuration.CurrentUserName)
	}
	availableUser, err := s.dbConnection.GetUser(context.Background(), cmd.arguements[0])
	if err != nil || availableUser.Name != cmd.arguements[0] {
		return fmt.Errorf("user \"%s\" doesn't exist. Error: %s", cmd.arguements[0], err)
	}
	err = internal.SetUser(*s.configuration, availableUser.Name)
	if err != nil {
		return err
	}
	data, err := internal.Read()
	if err != nil {
		return err
	}
	s.configuration = &data
	fmt.Println("Logged in as " + s.configuration.CurrentUserName)
	return nil
}

func handlerRegister(s *state, cmd command) (err error) {
	if len(cmd.arguements) != 1 {
		return fmt.Errorf("Login expects exactly 1 argument, received %v", cmd.arguements)
	}

	availableUser, err := s.dbConnection.GetUser(context.Background(), cmd.arguements[0])
	if availableUser.Name == cmd.arguements[0] {
		return fmt.Errorf("user \"%s\" already exists.", availableUser.Name)
	}

	user, err := s.dbConnection.CreateUser(context.Background(), database.CreateUserParams{ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguements[0]})
	if err != nil {
		return err
	}
	fmt.Printf("User \"%s\" created.\n", user.Name)
	err = handlerLogin(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func presentAllUsers(s *state, cmd command) (err error) {
	users, err := s.dbConnection.GetUsers(context.Background())
	if err != nil {
		return err
	}
	currentBadge := ""

	w := tabwriter.NewWriter(os.Stdout, 2, 0, 2, ' ', 0)
	for _, user := range users {
		currentBadge = ""
		if user == s.configuration.CurrentUserName {
			currentBadge = " (current)"
		}
		fmt.Fprintf(w, "%s\t\t%s\n", user, currentBadge)
	}
	w.Flush()
	return nil
}

func resetData(s *state, cmd command) (err error) {
	err = s.dbConnection.Reset(context.Background())
	if err != nil {
		return err
	}
	err = internal.SetUser(*s.configuration, "admin")
	if err != nil {
		return err
	}
	data, err := internal.Read()
	if err != nil {
		return err
	}
	s.configuration = &data
	fmt.Println("Reset successful.")
	return nil
}

func addFeed(s *state, cmd command, user database.User) (err error) {
	if len(cmd.arguements) != 2 {
		return fmt.Errorf("Expected two arguments in \n\tname: Name of the feed.\n\turl: The url of the feed.\n Received %v", cmd.arguements)
	}

	_, err = internal.FetchFeed(context.Background(), cmd.arguements[1])
	if err != nil {
		return fmt.Errorf("Error fetching from %s. Error: %s", cmd.arguements[1], err)
	}

	feedDB, err := s.dbConnection.AddFeed(context.Background(),
		database.AddFeedParams{ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.arguements[0],
			Url:       cmd.arguements[1],
			UserID:    user.ID})

	if err != nil {
		return err
	}

	feedUserData, err := s.dbConnection.AddFeedtoUser(context.Background(),
		database.AddFeedtoUserParams{
			ID:        uuid.New(),
			CreatedAt: feedDB.CreatedAt,
			UpdatedAt: feedDB.UpdatedAt,
			UserID:    user.ID,
			FeedID:    feedDB.ID,
		})
	if err != nil {
		return err
	}

	fmt.Println("Added feed " + feedUserData.FeedName + "(" + feedDB.Url + "), for user " + feedUserData.UserName + ".")
	return nil
}

func getAllFeed(s *state, cmd command) (err error) {
	feedData, err := s.dbConnection.GetAllFeeds(context.Background())
	if err != nil {
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0)

	fmt.Fprintln(w, "Feed\tURL\tSubscriber")

	for _, feedDatum := range feedData {
		fmt.Fprintf(w, "%s\t%s\t%s\n", feedDatum.Name, feedDatum.Url, feedDatum.Name_2.String)
	}
	w.Flush()
	return nil
}

func followHandler(s *state, cmd command, user database.User) (err error) {
	if len(cmd.arguements) != 1 {
		return fmt.Errorf("Expected 1 argument in url. Received %v.", cmd.arguements)
	}

	feedDB, err := s.dbConnection.GetFeed(context.Background(), cmd.arguements[0])
	if err != nil {
		return err
	}

	followData, err := s.dbConnection.AddFeedtoUser(context.Background(), database.AddFeedtoUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feedDB.ID,
	})
	if err != nil {
		return err
	}
	fmt.Printf("%s feed added to user %s.", followData.FeedName, followData.UserName)
	return nil
}

func followingHandler(s *state, cmd command, user database.User) (err error) {
	followUser, err := s.dbConnection.GetFeedFollowesFor(context.Background(), user.ID)
	if err != nil {
		return err
	}

	fmt.Println("Content Followed by user " + user.Name + ": ")
	for index, content := range followUser {
		fmt.Printf("\t%v) %s at address %s\n", index+1, content.Feedname, content.Url)
	}
	return nil
}

func unfollowHandler(s *state, cmd command, user database.User) (err error) {
	if len(cmd.arguements) != 1 {
		return fmt.Errorf("Expected 1 arguement in url, got %v", cmd.arguements)
	}
	_, err = s.dbConnection.UnFollow(context.Background(),
		database.UnFollowParams{
			UserID: user.ID,
			Url:    cmd.arguements[0],
		})
	if err != nil {
		return fmt.Errorf("Error unfollowing %w", err)
	}
	fmt.Printf("Unfollowed %s, as %s.\n", cmd.arguements[0], user.Name)
	return nil
}

func scrapeFeedFor(s *state, feed database.Feed) (err error) {
	fmt.Printf("Aggregating from %s@%s\n", feed.Name, feed.Url)
	fetchedFeed, err := internal.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("Error fetching form address %s. Error: %s", feed.Url, err)
	}
	_, err = s.dbConnection.MarkFeedtoFetch(context.Background(),
		database.MarkFeedtoFetchParams{
			UpdatedAt: time.Now(),
			ID:        feed.ID,
		})
	if err != nil {
		return fmt.Errorf("Error updating fetch date for %s. Error: %s", feed.Url, err)
	}
	for _, fetchedSubFeed := range fetchedFeed.Channel.Item {
		publishTime, err := time.Parse(time.RFC1123, fetchedSubFeed.PubDate)
		if err != nil {
			publishTime = time.Now()
		}
		_, err = s.dbConnection.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       fetchedSubFeed.Title,
			Url:         fetchedSubFeed.Link,
			Description: sql.NullString{String: fetchedSubFeed.Description, Valid: len(fetchedSubFeed.Description) > 0},
			PublishedAt: publishTime,
			FeedID:      feed.ID,
		})
	}
	return nil

}

func aggHandler(s *state, cmd command, user database.User) (err error) {
	feeds, err := s.dbConnection.GetMyFeeds(context.Background(), user.ID)
	for _, feed := range feeds {
		err = scrapeFeedFor(s, feed)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func browseHandler(s *state, cmd command, user database.User) (err error) {
	var limit int
	var limitAddressed = false
	if len(cmd.arguements) == 1 {
		limit, err = strconv.Atoi(cmd.arguements[0])
		if err != nil {
			limit = 5
		}
		limitAddressed = true
	}

	posts, err := s.dbConnection.GetPostsForUser(context.Background(), user.ID)
	fmt.Println("Recent posts from RSS feed...")
	if limitAddressed {
		limit = min(len(posts), limit)
	} else {
		limit = len(posts)
	}
	for i := 0; i < limit; i++ {
		fmt.Println(posts[i].PublishedAt, posts[i].Title)

		fmt.Println("\t", posts[i].Url)
		fmt.Println("\t", posts[i].Description.String)
		fmt.Println("=============================================================")
		fmt.Println()
	}
	return nil
}

func (c *commands) run(s *state, cmd command) (err error) {
	runner, exist := c.commandMap[cmd.command]
	if !exist {
		return fmt.Errorf("No such command as %s", cmd.command)
	}
	err = runner(s, cmd)
	return err
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandMap[name] = f
}

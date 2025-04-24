package main

import (
	"github.com/abahnj/rssagg/internal/cli"
	"github.com/abahnj/rssagg/internal/feeds"
	"github.com/abahnj/rssagg/internal/middleware"
	"github.com/abahnj/rssagg/internal/users"
)

// registerCommands sets up all available commands
func registerCommands(commands *cli.Commands) {
	// User commands
	commands.Register("login", users.HandlerLogin)
	commands.Register("register", users.HandlerRegister)
	commands.Register("reset", users.HandlerDeleteAllUsers)
	commands.Register("users", users.HandlerListUsers)
	
	// Feed commands
	commands.Register("agg", feeds.HandlerAggregator)
	commands.Register("feeds", feeds.HandlerListFeeds)
	
	// Protected feed commands (requiring authentication)
	commands.Register("addfeed", middleware.MiddlewareLoggedIn(feeds.HandlerAddFeed))
	commands.Register("follow", middleware.MiddlewareLoggedIn(feeds.HandlerFollowFeed))
	commands.Register("following", middleware.MiddlewareLoggedIn(feeds.HandlerListFollowing))
	commands.Register("unfollow", middleware.MiddlewareLoggedIn(feeds.HandlerUnfollowFeed))
	commands.Register("browse", middleware.MiddlewareLoggedIn(feeds.HandlerBrowse))
}
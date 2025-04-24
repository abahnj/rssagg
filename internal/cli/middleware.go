package cli

import (
	"context"
	"errors"
	"fmt"

	"github.com/abahnj/rssagg/internal/database"
)

// LoggedInHandlerFunc is a handler function that requires a logged-in user
type LoggedInHandlerFunc func(*State, Command, database.User) error

// MiddlewareLoggedIn is middleware that ensures a user is logged in
// before executing a handler that requires authentication
func MiddlewareLoggedIn(handler LoggedInHandlerFunc) HandlerFunc {
	return func(s *State, cmd Command) error {
		// Check if a user is logged in
		if s.Config == nil || s.Config.CurrentUserName == "" {
			return errors.New("you must be logged in to use this command")
		}

		ctx := context.Background()
		
		// Get the current user from the database
		user, err := s.Db.GetUserByName(ctx, s.Config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("failed to get current user: %w", err)
		}

		// Call the wrapped handler with the authenticated user
		return handler(s, cmd, user)
	}
}
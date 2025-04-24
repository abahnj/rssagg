package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/abahnj/rssagg/internal/cli"
)

// HandlerLogin handles the login command
func HandlerLogin(s *cli.State, cmd cli.Command) error {
	if len(cmd.Args) < 1 {
		return ErrMissingUsername
	}

	// Check for nil config
	if s.Config == nil {
		return errors.New("config is not initialized")
	}

	username := cmd.Args[0]
	ctx := context.Background()

	service := NewService(*s.Db)
	_, err := service.Login(ctx, username)
	if err != nil {
		return err
	}

	if err := s.Config.SetUser(username); err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}

	fmt.Printf("User %s logged in\n", username)
	return nil
}

// HandlerRegister handles the register command
func HandlerRegister(s *cli.State, cmd cli.Command) error {
	if len(cmd.Args) < 1 {
		return errors.New("username is required for register")
	}

	// Check for nil config
	if s.Config == nil {
		return errors.New("config is not initialized")
	}

	username := cmd.Args[0]
	ctx := context.Background()
	
	service := NewService(*s.Db)
	_, err := service.Register(ctx, username)
	if err != nil {
		return err
	}

	return HandlerLogin(s, cmd)
}

// HandlerDeleteAllUsers handles the reset command
func HandlerDeleteAllUsers(s *cli.State, cmd cli.Command) error {
	ctx := context.Background()
	
	service := NewService(*s.Db)
	if err := service.DeleteAllUsers(ctx); err != nil {
		return err
	}
	
	fmt.Println("All users deleted")
	return nil
}

// HandlerListUsers handles the users command
func HandlerListUsers(s *cli.State, cmd cli.Command) error {
	ctx := context.Background()
	
	service := NewService(*s.Db)
	users, err := service.GetUsers(ctx)
	if err != nil {
		return err
	}
	
	currentUserName := ""
	if s.Config != nil {
		currentUserName = s.Config.CurrentUserName
	}
	
	for _, user := range users {
		if user.Name == currentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	
	return nil
}
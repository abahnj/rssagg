package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/abahnj/rssagg/internal/database"
	"github.com/google/uuid"
)

// ErrMissingUsername is returned when the login command doesn't have a username argument
var ErrMissingUsername = errors.New("username is required for login")

// Service handles user management operations
type Service struct {
	DB database.Queries
}

// NewService creates a new user service
func NewService(db database.Queries) *Service {
	return &Service{
		DB: db,
	}
}

// Login authenticates a user
func (s *Service) Login(ctx context.Context, username string) (database.User, error) {
	user, err := s.DB.GetUserByName(ctx, username)
	if err != nil {
		return database.User{}, fmt.Errorf("failed to fetch user with name %s: %w", username, err)
	}
	return user, nil
}

// Register creates a new user
func (s *Service) Register(ctx context.Context, username string) (database.User, error) {
	createUserParams := database.CreateUserParams{
		ID:   uuid.New(),
		Name: username,
	}

	user, err := s.DB.CreateUser(ctx, createUserParams)
	if err != nil {
		return database.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUsers retrieves all users
func (s *Service) GetUsers(ctx context.Context) ([]database.User, error) {
	users, err := s.DB.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	return users, nil
}

// DeleteAllUsers removes all users from the database
func (s *Service) DeleteAllUsers(ctx context.Context) error {
	if err := s.DB.DeleteAllUsers(ctx); err != nil {
		return fmt.Errorf("failed to delete all users: %w", err)
	}
	return nil
}
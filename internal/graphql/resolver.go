// Package graphql provides GraphQL schema definitions and resolvers
// for the demo-go application, implementing user management operations
// through a GraphQL API interface.
package graphql

import (
	"context"
	"strings"

	"demo-go/internal/domain"
	"demo-go/internal/logger"
)

// Resolver is the root resolver for GraphQL operations
type Resolver struct {
	userService domain.UserService
	logger      *logger.Logger
}

// NewResolver creates a new GraphQL resolver
func NewResolver(userService domain.UserService) *Resolver {
	return &Resolver{
		userService: userService,
		logger:      logger.GetGlobal().ForComponent("graphql-resolver"),
	}
}

// Query resolver
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

// Mutation resolver
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

// Subscription resolver
func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

// queryResolver implements QueryResolver interface
type queryResolver struct{ *Resolver }

// GetUser resolves the getUser query
func (r *queryResolver) GetUser(ctx context.Context, id string) (*domain.UserResponse, error) {
	log := r.logger.ForService("query", "getUser").WithField("user_id", id)

	log.Debug("Resolving getUser query")

	user, err := r.userService.GetUserByID(ctx, id)
	if err != nil {
		log.Error("Failed to get user", "error", err)
		return nil, err
	}

	log.Debug("Successfully resolved getUser query", "user_email", user.Email)
	return user, nil
}

// GetUsers resolves the getUsers query
func (r *queryResolver) GetUsers(ctx context.Context, limit *int, offset *int) ([]*domain.UserResponse, error) {
	log := r.logger.ForService("query", "getUsers")

	// Set default values if not provided
	if limit == nil {
		defaultLimit := 10
		limit = &defaultLimit
	}
	if offset == nil {
		defaultOffset := 0
		offset = &defaultOffset
	}

	log.Debug("Resolving getUsers query", "limit", *limit, "offset", *offset)

	users, _, err := r.userService.GetUsers(ctx, *limit, *offset)
	if err != nil {
		log.Error("Failed to get users", "error", err)
		return nil, err
	}

	// Apply pagination
	start := *offset
	end := start + *limit

	if start >= len(users) {
		return []*domain.UserResponse{}, nil
	}

	if end > len(users) {
		end = len(users)
	}

	paginatedUsers := users[start:end]
	log.Debug("Successfully resolved getUsers query", "total_users", len(users), "returned_users", len(paginatedUsers))

	return paginatedUsers, nil
}

// SearchUsers resolves the searchUsers query
func (r *queryResolver) SearchUsers(ctx context.Context, query string) ([]*domain.UserResponse, error) {
	log := r.logger.ForService("query", "searchUsers").WithField("search_query", query)

	log.Debug("Resolving searchUsers query")

	// Get all users and filter by name or email
	users, _, err := r.userService.GetUsers(ctx, 1000, 0) // Get up to 1000 users for search
	if err != nil {
		log.Error("Failed to get users for search", "error", err)
		return nil, err
	}

	var filteredUsers []*domain.UserResponse
	for _, user := range users {
		if containsIgnoreCase(user.Name, query) || containsIgnoreCase(user.Email, query) {
			filteredUsers = append(filteredUsers, user)
		}
	}

	log.Debug("Successfully resolved searchUsers query", "matches_found", len(filteredUsers))
	return filteredUsers, nil
}

// Me resolves the me query (returns current authenticated user)
func (r *queryResolver) Me(ctx context.Context) (*domain.UserResponse, error) {
	log := r.logger.ForService("query", "me")

	log.Debug("Resolving me query")

	// Get user ID from context (set by authentication middleware)
	userID, ok := ctx.Value("userID").(string)
	if !ok {
		log.Warn("User ID not found in context")
		return nil, domain.ErrUnauthorized
	}

	user, err := r.userService.GetUserByID(ctx, userID)
	if err != nil {
		log.Error("Failed to get current user", "user_id", userID, "error", err)
		return nil, err
	}

	log.Debug("Successfully resolved me query", "user_email", user.Email)
	return user, nil
}

// mutationResolver implements MutationResolver interface
type mutationResolver struct{ *Resolver }

// CreateUser resolves the createUser mutation
func (r *mutationResolver) CreateUser(ctx context.Context, input CreateUserInput) (*domain.UserResponse, error) {
	log := r.logger.ForService("mutation", "createUser").WithField("email", input.Email)

	log.Debug("Resolving createUser mutation")

	createReq := &domain.CreateUserRequest{
		Name:     input.Name,
		Email:    input.Email,
		Password: "default-password", // In a real app, this should be provided or generated
		Role:     "user",
	}

	user, err := r.userService.Register(ctx, createReq)
	if err != nil {
		log.Error("Failed to create user", "error", err)
		return nil, err
	}

	log.Info("Successfully created user", "user_id", user.ID, "user_email", user.Email)
	return user, nil
}

// UpdateUser resolves the updateUser mutation
func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input UpdateUserInput) (*domain.UserResponse, error) {
	log := r.logger.ForService("mutation", "updateUser").WithField("user_id", id)

	log.Debug("Resolving updateUser mutation")

	updateReq := &domain.UpdateUserRequest{}

	if input.Name != nil {
		updateReq.Name = input.Name
	}
	if input.Email != nil {
		updateReq.Email = input.Email
	}

	user, err := r.userService.UpdateProfile(ctx, id, updateReq)
	if err != nil {
		log.Error("Failed to update user", "error", err)
		return nil, err
	}

	log.Info("Successfully updated user", "user_id", user.ID, "user_email", user.Email)
	return user, nil
}

// DeleteUser resolves the deleteUser mutation
func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
	log := r.logger.ForService("mutation", "deleteUser").WithField("user_id", id)

	log.Debug("Resolving deleteUser mutation")

	err := r.userService.DeleteUser(ctx, id)
	if err != nil {
		log.Error("Failed to delete user", "error", err)
		return false, err
	}

	log.Info("Successfully deleted user", "user_id", id)
	return true, nil
}

// subscriptionResolver implements SubscriptionResolver interface
type subscriptionResolver struct{ *Resolver }

// UserCreated resolves the userCreated subscription
func (r *subscriptionResolver) UserCreated(ctx context.Context) (<-chan *domain.UserResponse, error) {
	log := r.logger.ForService("subscription", "userCreated")

	log.Debug("Setting up userCreated subscription")

	// Create a channel for user creation events
	userChan := make(chan *domain.UserResponse, 1)

	// In a real implementation, you would connect to a message broker or event system
	// For now, we'll just return an empty channel
	go func() {
		<-ctx.Done()
		close(userChan)
	}()

	return userChan, nil
}

// UserUpdated resolves the userUpdated subscription
func (r *subscriptionResolver) UserUpdated(ctx context.Context) (<-chan *domain.UserResponse, error) {
	log := r.logger.ForService("subscription", "userUpdated")

	log.Debug("Setting up userUpdated subscription")

	userChan := make(chan *domain.UserResponse, 1)

	go func() {
		<-ctx.Done()
		close(userChan)
	}()

	return userChan, nil
}

// UserDeleted resolves the userDeleted subscription
func (r *subscriptionResolver) UserDeleted(ctx context.Context) (<-chan string, error) {
	log := r.logger.ForService("subscription", "userDeleted")

	log.Debug("Setting up userDeleted subscription")

	userIDChan := make(chan string, 1)

	go func() {
		<-ctx.Done()
		close(userIDChan)
	}()

	return userIDChan, nil
}

// Helper functions

// containsIgnoreCase checks if the haystack contains the needle (case-insensitive)
func containsIgnoreCase(haystack, needle string) bool {
	return strings.Contains(strings.ToLower(haystack), strings.ToLower(needle))
}

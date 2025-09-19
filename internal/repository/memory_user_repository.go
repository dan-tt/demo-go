package repository

import (
	"context"
	"strconv"
	"sync"
	"time"

	"demo-go/internal/domain"
)

// memoryUserRepository implements domain.UserRepository using in-memory storage
type memoryUserRepository struct {
	users    map[string]*domain.User
	emails   map[string]string // email -> userID mapping for unique constraint
	mu       sync.RWMutex
	nextID   int
}

// NewMemoryUserRepository creates a new in-memory user repository
func NewMemoryUserRepository() domain.UserRepository {
	return &memoryUserRepository{
		users:  make(map[string]*domain.User),
		emails: make(map[string]string),
		nextID: 1,
	}
}

// Create creates a new user in memory
func (r *memoryUserRepository) Create(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Check if email already exists
	if _, exists := r.emails[user.Email]; exists {
		return domain.ErrUserAlreadyExists
	}
	
	// Generate ID if not provided
	if user.ID == "" {
		user.ID = strconv.Itoa(r.nextID)
		r.nextID++
	}
	
	// Set creation time
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	
	// Store user
	r.users[user.ID] = user
	r.emails[user.Email] = user.ID
	
	return nil
}

// GetByID retrieves a user by ID from memory
func (r *memoryUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	user, exists := r.users[id]
	if !exists {
		return nil, domain.ErrUserNotFound
	}
	
	// Return a copy to prevent external modifications
	userCopy := *user
	return &userCopy, nil
}

// GetByEmail retrieves a user by email from memory
func (r *memoryUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	userID, exists := r.emails[email]
	if !exists {
		return nil, domain.ErrUserNotFound
	}
	
	user := r.users[userID]
	// Return a copy to prevent external modifications
	userCopy := *user
	return &userCopy, nil
}

// Update updates a user in memory
func (r *memoryUserRepository) Update(ctx context.Context, id string, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	existingUser, exists := r.users[id]
	if !exists {
		return domain.ErrUserNotFound
	}
	
	// Check if email is being changed and if new email already exists
	if user.Email != existingUser.Email {
		if _, emailExists := r.emails[user.Email]; emailExists {
			return domain.ErrUserAlreadyExists
		}
		
		// Remove old email mapping and add new one
		delete(r.emails, existingUser.Email)
		r.emails[user.Email] = id
	}
	
	// Update user fields
	user.ID = id // Ensure ID doesn't change
	user.CreatedAt = existingUser.CreatedAt // Preserve creation time
	user.UpdatedAt = time.Now()
	
	// Store updated user
	r.users[id] = user
	
	return nil
}

// Delete deletes a user from memory
func (r *memoryUserRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	user, exists := r.users[id]
	if !exists {
		return domain.ErrUserNotFound
	}
	
	// Remove from both maps
	delete(r.users, id)
	delete(r.emails, user.Email)
	
	return nil
}

// List retrieves users with pagination from memory
func (r *memoryUserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	// Convert map to slice for sorting and pagination
	var allUsers []*domain.User
	for _, user := range r.users {
		userCopy := *user
		allUsers = append(allUsers, &userCopy)
	}
	
	// Sort by creation time (newest first)
	for i := 0; i < len(allUsers)-1; i++ {
		for j := i + 1; j < len(allUsers); j++ {
			if allUsers[i].CreatedAt.Before(allUsers[j].CreatedAt) {
				allUsers[i], allUsers[j] = allUsers[j], allUsers[i]
			}
		}
	}
	
	// Apply pagination
	start := offset
	if start > len(allUsers) {
		return []*domain.User{}, nil
	}
	
	end := start + limit
	if end > len(allUsers) {
		end = len(allUsers)
	}
	
	return allUsers[start:end], nil
}

// Count returns the total number of users in memory
func (r *memoryUserRepository) Count(ctx context.Context) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	return int64(len(r.users)), nil
}

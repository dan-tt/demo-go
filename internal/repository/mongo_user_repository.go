package repository

import (
	"context"
	"time"

	"demo-go/internal/config"
	"demo-go/internal/domain"
	"demo-go/internal/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongoUserRepository implements domain.UserRepository using MongoDB
type mongoUserRepository struct {
	collection *mongo.Collection
	timeout    time.Duration
	logger     *logger.Logger
}

// NewMongoUserRepository creates a new MongoDB user repository
func NewMongoUserRepository(client *mongo.Client, cfg *config.Config) domain.UserRepository {
	log := logger.GetGlobal().ForComponent("mongo-repository")

	collection := client.Database(cfg.Database.MongoDB.Database).Collection("users")

	// Create unique index on email
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Database.MongoDB.Timeout)
	defer cancel()

	log.Debug("Creating unique index on email field")
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Warn("Failed to create email index", "error", err)
	} else {
		log.Debug("Email index created successfully")
	}

	return &mongoUserRepository{
		collection: collection,
		timeout:    cfg.Database.MongoDB.Timeout,
		logger:     log,
	}
}

// Create creates a new user in MongoDB
func (r *mongoUserRepository) Create(ctx context.Context, user *domain.User) error {
	log := r.logger.ForRepository("user", "create").WithField("email", user.Email)

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	log.Debug("Creating user in MongoDB")

	// Set creation time
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// If ID is empty, MongoDB will generate one
	if user.ID == "" {
		user.ID = primitive.NewObjectID().Hex()
	}

	log.Debug("Inserting user document", "user_id", user.ID)

	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Warn("Duplicate email detected", "error", err)
			return domain.ErrUserAlreadyExists
		}
		log.Error("Failed to insert user", "error", err)
		return err
	}

	log.Info("User created successfully", "user_id", user.ID)
	return nil
}

// GetByID retrieves a user by ID from MongoDB
func (r *mongoUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// GetByEmail retrieves a user by email from MongoDB
func (r *mongoUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// Update updates a user in MongoDB
func (r *mongoUserRepository) Update(ctx context.Context, id string, user *domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	// Set update time
	user.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"name":       user.Name,
			"email":      user.Email,
			"role":       user.Role,
			"updated_at": user.UpdatedAt,
		},
	}

	// Only update password if it's provided
	if user.Password != "" {
		if setMap, ok := update["$set"].(bson.M); ok {
			setMap["password"] = user.Password
		}
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.ErrUserAlreadyExists
		}
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// Delete deletes a user from MongoDB
func (r *mongoUserRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// List retrieves users with pagination from MongoDB
func (r *mongoUserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	opts := options.Find().
		SetLimit(int64(limit)).
		SetSkip(int64(offset)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			// Log the error but don't fail the operation
			// since the main operation was successful
		}
	}()

	var users []*domain.User
	for cursor.Next(ctx) {
		var user domain.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Count returns the total number of users in MongoDB
func (r *mongoUserRepository) Count(ctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}

	return count, nil
}

// NewMongoClient creates a new MongoDB client
func NewMongoClient(cfg *config.Config) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Database.MongoDB.Timeout)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(cfg.Database.MongoDB.URI)

	// Safely convert int to uint64 for MaxPoolSize
	if cfg.Database.MongoDB.MaxPoolSize > 0 {
		clientOptions = clientOptions.SetMaxPoolSize(uint64(cfg.Database.MongoDB.MaxPoolSize))
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client, nil
}

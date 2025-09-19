// MongoDB initialization script
// This script runs when the MongoDB container first starts

// Switch to the demo_clean database
db = db.getSiblingDB('demo_clean');

// Create the users collection with schema validation
db.createCollection('users', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['name', 'email', 'password', 'role', 'created_at', 'updated_at'],
      properties: {
        _id: {
          bsonType: 'string',
          description: 'User ID - must be a string'
        },
        name: {
          bsonType: 'string',
          minLength: 2,
          maxLength: 100,
          description: 'User name - must be a string between 2-100 characters'
        },
        email: {
          bsonType: 'string',
          pattern: '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$',
          description: 'User email - must be a valid email format'
        },
        password: {
          bsonType: 'string',
          minLength: 6,
          description: 'User password hash - must be at least 6 characters'
        },
        role: {
          bsonType: 'string',
          enum: ['user', 'admin'],
          description: 'User role - must be either user or admin'
        },
        created_at: {
          bsonType: 'date',
          description: 'Creation timestamp - must be a date'
        },
        updated_at: {
          bsonType: 'date',
          description: 'Update timestamp - must be a date'
        }
      }
    }
  }
});

// Create unique index on email field
db.users.createIndex({ email: 1 }, { unique: true });

// Create index on created_at for efficient sorting
db.users.createIndex({ created_at: -1 });

// Create index on role for efficient role-based queries
db.users.createIndex({ role: 1 });

// Insert a default admin user for testing
db.users.insertOne({
  _id: 'admin-user-001',
  name: 'System Administrator',
  email: 'admin@demo-clean.com',
  password: '$2a$10$2Ygwc1gZnCNKGgK4KqPO3eYHGPa1ywKTaW2xlQrLFhGDCnMQN.ZJC', // bcrypt hash of 'admin123'
  role: 'admin',
  created_at: new Date(),
  updated_at: new Date()
});

// Insert a default regular user for testing
db.users.insertOne({
  _id: 'regular-user-001',
  name: 'Demo User',
  email: 'user@demo-clean.com',
  password: '$2a$10$9.qVZ.QGOvYzJKQ8LcJ8uOKE4HPVzDzH5KMX.mOlN5zG6q6YLqVWW', // bcrypt hash of 'user123'
  role: 'user',
  created_at: new Date(),
  updated_at: new Date()
});

print('MongoDB initialization completed successfully!');
print('Collections created: users');
print('Indexes created: email (unique), created_at, role');
print('Default users created:');
print('  - Admin: admin@demo-clean.com / admin123');
print('  - User: user@demo-clean.com / user123');

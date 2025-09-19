import http from 'k6/http';
import { check, group, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// Custom metrics
export let errorRate = new Rate('errors');

// Test configuration
export let options = {
  stages: [
    { duration: '2m', target: 10 }, // Ramp up to 10 users
    { duration: '5m', target: 10 }, // Stay at 10 users
    { duration: '2m', target: 20 }, // Ramp up to 20 users
    { duration: '5m', target: 20 }, // Stay at 20 users
    { duration: '2m', target: 0 },  // Ramp down to 0 users
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% of requests must complete below 500ms
    http_req_failed: ['rate<0.1'],    // Error rate must be below 10%
    errors: ['rate<0.1'],             // Custom error rate must be below 10%
  },
};

// Base URL from environment variable
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

// Test data
const users = [
  { name: 'Alice Johnson', email: 'alice@example.com', password: 'password123' },
  { name: 'Bob Smith', email: 'bob@example.com', password: 'password456' },
  { name: 'Charlie Brown', email: 'charlie@example.com', password: 'password789' },
];

let authToken = '';

export function setup() {
  // Setup phase - create a test user and get auth token
  console.log('Setting up test data...');
  
  // Health check
  let healthResponse = http.get(`${BASE_URL}/health`);
  check(healthResponse, {
    'health check status is 200': (r) => r.status === 200,
  });
  
  // Register a test user
  let registerPayload = JSON.stringify({
    name: 'Test User',
    email: 'test@example.com',
    password: 'testpassword123',
    role: 'user'
  });
  
  let registerResponse = http.post(`${BASE_URL}/auth/register`, registerPayload, {
    headers: { 'Content-Type': 'application/json' },
  });
  
  if (registerResponse.status === 201) {
    let userData = JSON.parse(registerResponse.body);
    return { authToken: userData.access_token };
  }
  
  return { authToken: '' };
}

export default function(data) {
  authToken = data.authToken;
  
  group('API Load Test', function() {
    // Test 1: Health Check
    group('Health Check', function() {
      let response = http.get(`${BASE_URL}/health`);
      check(response, {
        'health check status is 200': (r) => r.status === 200,
        'health check response time < 100ms': (r) => r.timings.duration < 100,
      }) || errorRate.add(1);
    });
    
    // Test 2: User Registration
    group('User Registration', function() {
      let user = users[Math.floor(Math.random() * users.length)];
      let payload = JSON.stringify({
        ...user,
        email: `${Math.random().toString(36).substring(7)}@example.com`, // Random email
        role: 'user'
      });
      
      let response = http.post(`${BASE_URL}/auth/register`, payload, {
        headers: { 'Content-Type': 'application/json' },
      });
      
      check(response, {
        'registration status is 201': (r) => r.status === 201,
        'registration response time < 1s': (r) => r.timings.duration < 1000,
        'registration returns access token': (r) => {
          try {
            let body = JSON.parse(r.body);
            return body.access_token && body.access_token.length > 0;
          } catch (e) {
            return false;
          }
        },
      }) || errorRate.add(1);
    });
    
    // Test 3: User Login
    group('User Login', function() {
      let payload = JSON.stringify({
        email: 'test@example.com',
        password: 'testpassword123'
      });
      
      let response = http.post(`${BASE_URL}/auth/login`, payload, {
        headers: { 'Content-Type': 'application/json' },
      });
      
      check(response, {
        'login status is 200': (r) => r.status === 200,
        'login response time < 500ms': (r) => r.timings.duration < 500,
        'login returns access token': (r) => {
          try {
            let body = JSON.parse(r.body);
            return body.access_token && body.access_token.length > 0;
          } catch (e) {
            return false;
          }
        },
      }) || errorRate.add(1);
    });
    
    // Test 4: Get User Profile (Authenticated)
    if (authToken) {
      group('Get User Profile', function() {
        let response = http.get(`${BASE_URL}/api/v1/profile`, {
          headers: { 'Authorization': `Bearer ${authToken}` },
        });
        
        check(response, {
          'profile status is 200': (r) => r.status === 200,
          'profile response time < 300ms': (r) => r.timings.duration < 300,
          'profile returns user data': (r) => {
            try {
              let body = JSON.parse(r.body);
              return body.id && body.name && body.email;
            } catch (e) {
              return false;
            }
          },
        }) || errorRate.add(1);
      });
    }
    
    // Test 5: Cache Performance Test
    group('Cache Performance', function() {
      if (authToken) {
        // First request (potential cache miss)
        let response1 = http.get(`${BASE_URL}/api/v1/profile`, {
          headers: { 'Authorization': `Bearer ${authToken}` },
        });
        
        // Second request (should be cache hit)
        let response2 = http.get(`${BASE_URL}/api/v1/profile`, {
          headers: { 'Authorization': `Bearer ${authToken}` },
        });
        
        check(response1, {
          'first request status is 200': (r) => r.status === 200,
        }) || errorRate.add(1);
        
        check(response2, {
          'second request status is 200': (r) => r.status === 200,
          'second request faster than first (cache hit)': (r) => r.timings.duration <= response1.timings.duration,
        }) || errorRate.add(1);
      }
    });
  });
  
  // Think time between requests
  sleep(Math.random() * 2 + 1); // 1-3 seconds
}

export function teardown(data) {
  // Cleanup phase
  console.log('Cleaning up test data...');
  // Add any cleanup logic here if needed
}

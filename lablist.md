# Build a Mini HTTP Framework (Like Echo/Gin)

---

## Lab 01: Framework Foundation & Project Setup

**Learning Objectives:**

- Set up Go project structure for a reusable framework
- Understand how `net/http` works under the hood
- Create the base `Engine` struct that implements `http.Handler`
- Implement basic server lifecycle (start, shutdown)

**Tasks:**

| Task | Description |
| --- | --- |
| 1.1 | Initialize Go module with proper package structure |
| 1.2 | Create `Engine` struct with `ServeHTTP` method |
| 1.3 | Implement `Run()` and `Shutdown()` methods |
| 1.4 | Add configuration options (port, timeouts) |
| 1.5 | Write unit tests for server lifecycle |

**Deliverables:**

```
goexpress/
├── goexpress.go      # Engine struct, New(), Run()
├── config.go         # Configuration options
├── goexpress_test.go # Unit tests
└── go.mod

```

**Success Criteria:**

- [ ]  Server starts on specified port
- [ ]  Graceful shutdown works correctly
- [ ]  Configuration options are applied
- [ ]  All tests pass

---

## Lab 02: Context Design & Implementation

**Learning Objectives:**

- Design the `Context` struct as request/response wrapper
- Implement request data accessors (query, header, body)
- Implement response writers (status, header, body)
- Add request-scoped storage with `Set()` and `Get()`

**Tasks:**

| Task | Description |
| --- | --- |
| 2.1 | Create `Context` struct wrapping req/res |
| 2.2 | Implement request accessors (Query, Header, Body) |
| 2.3 | Implement response methods (Status, SetHeader, Write) |
| 2.4 | Add context storage (Set, Get, MustGet) |
| 2.5 | Write comprehensive tests |

**Deliverables:**

```
goexpress/
├── context.go        # Context struct and methods
├── context_test.go   # Context unit tests
├── request.go        # Request helper methods
└── response.go       # Response helper methods

```

**Key Implementation:**

```go
type Context struct {
    Request  *http.Request
    Writer   http.ResponseWriter
    Params   map[string]string
    store    map[string]any
    index    int
    handlers []HandlerFunc
}

```

**Success Criteria:**

- [ ]  Query parameters accessible via `ctx.Query("key")`
- [ ]  Headers accessible via `ctx.GetHeader("key")`
- [ ]  Response status and headers work correctly
- [ ]  Context values stored and retrieved properly

---

## Lab 03: Router Core Implementation

**Learning Objectives:**

- Implement route registration for all HTTP methods
- Build route storage structure (method → path → handler)
- Implement basic route matching algorithm
- Handle 404 (Not Found) and 405 (Method Not Allowed)

**Tasks:**

| Task | Description |
| --- | --- |
| 3.1 | Define `HandlerFunc` type and route structure |
| 3.2 | Implement `GET`, `POST`, `PUT`, `DELETE`, `PATCH` methods |
| 3.3 | Build route matching in `ServeHTTP` |
| 3.4 | Implement 404 and 405 responses |
| 3.5 | Test all HTTP methods and error cases |

**Deliverables:**

```
goexpress/
├── router.go         # Route registration & matching
├── router_test.go    # Router unit tests
└── handlers.go       # Default handlers (404, 405)

```

**Key Implementation:**

```go
type HandlerFunc func(*Context)

type Router struct {
    routes map[string]map[string]HandlerFunc
}

func (r *Router) GET(path string, handler HandlerFunc)
func (r *Router) POST(path string, handler HandlerFunc)
func (r *Router) match(method, path string) (HandlerFunc, bool)

```

**Success Criteria:**

- [ ]  Routes register correctly for all methods
- [ ]  Exact path matching works
- [ ]  404 returned for unknown paths
- [ ]  405 returned for wrong methods

---

## Lab 04: Path Parameters Parser

**Learning Objectives:**

- Parse dynamic path segments (`:param`)
- Implement wildcard matching (`filepath`)
- Extract parameters during route matching
- Make parameters accessible via `ctx.Param()`

**Tasks:**

| Task | Description |
| --- | --- |
| 4.1 | Design path segment parser |
| 4.2 | Implement `:param` extraction |
| 4.3 | Implement `*wildcard` extraction |
| 4.4 | Integrate params into Context |
| 4.5 | Test various path patterns |

**Deliverables:**

```
goexpress/
├── tree.go           # Path parsing & matching tree
├── tree_test.go      # Tree unit tests
└── params.go         # Parameter extraction helpers

```

**Test Cases:**

```
Pattern: /users/:id          Request: /users/123      → id=123
Pattern: /users/:id/posts    Request: /users/5/posts  → id=5
Pattern: /files/*filepath    Request: /files/a/b/c    → filepath=a/b/c
Pattern: /api/:version/users Request: /api/v1/users   → version=v1

```

**Success Criteria:**

- [ ]  Single params extracted correctly
- [ ]  Multiple params in path work
- [ ]  Wildcard captures rest of path
- [ ]  `ctx.Param("name")` returns correct value

---

## Lab 05: Router Integration & Testing

**Learning Objectives:**

- Integrate all router components together
- Build comprehensive test suite
- Test edge cases and error conditions
- Benchmark route matching performance

**Tasks:**

| Task | Description |
| --- | --- |
| 5.1 | Integrate router with Engine |
| 5.2 | Write integration tests |
| 5.3 | Test edge cases (trailing slash, special chars) |
| 5.4 | Write benchmark tests |
| 5.5 | Document router API |

**Test Scenarios:**

```go
GET  /                     → 200 "Welcome"
GET  /users                → 200 "List users"
GET  /users/123            → 200 "User: 123"
GET  /users/123/posts      → 200 "Posts for user: 123"
POST /users                → 201 "User created"
GET  /files/css/style.css  → 200 "File: css/style.css"
GET  /notfound             → 404
POST /users/123            → 405 (if only GET registered)

```

**Success Criteria:**

- [ ]  All integration tests pass
- [ ]  Edge cases handled correctly
- [ ]  Benchmark shows acceptable performance
- [ ]  API documentation complete

---

## Lab 06: Middleware Pipeline Implementation

**Learning Objectives:**

- Design the middleware chain pattern
- Implement `Use()` for global middleware
- Build `Next()` function for chain execution
- Create Logger, Recovery, and CORS middleware

**Tasks:**

| Task | Description |
| --- | --- |
| 6.1 | Design middleware function signature |
| 6.2 | Implement middleware chain executor |
| 6.3 | Create Logger middleware |
| 6.4 | Create Recovery middleware |
| 6.5 | Create CORS middleware |
| 6.6 | Test middleware execution order |

**Deliverables:**

```
goexpress/
├── middleware.go
├── middleware/
│   ├── logger.go
│   ├── recovery.go
│   ├── cors.go
│   └── requestid.go
└── middleware_test.go

```

**Middleware Flow:**

```
Request → Logger.Before → Recovery.Before → Handler → Recovery.After → Logger.After → Response

```

**Success Criteria:**

- [ ]  Middleware executes in registration order
- [ ]  `Next()` passes control correctly
- [ ]  Recovery catches panics and returns 500
- [ ]  Logger outputs method, path, status, duration

---

## Lab 07: JSON Binding, Response & Validation

**Learning Objectives:**

- Implement `ctx.JSON()` for JSON responses
- Implement `ctx.Bind()` for request body binding
- Build struct tag validator (required, min, max, email)
- Return structured validation errors

**Tasks:**

| Task | Description |
| --- | --- |
| 7.1 | Implement `ctx.JSON()` and `ctx.String()` |
| 7.2 | Implement `ctx.Bind()` with content-type detection |
| 7.3 | Build tag parser for validation rules |
| 7.4 | Implement validators (required, min, max, email) |
| 7.5 | Create structured error responses |
| 7.6 | Test binding and validation |

**Deliverables:**

```
goexpress/
├── binding.go
├── render.go
├── validator.go
├── validator_test.go
└── errors.go

```

**Validation Example:**

```go
type CreateUser struct {
    Name  string `json:"name" validate:"required,min=2"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"gte=18"`
}

```

**Success Criteria:**

- [ ]  JSON response sets correct Content-Type
- [ ]  Bind correctly parses JSON body
- [ ]  Validation catches all rule violations
- [ ]  Error response includes field-level details

---

## Lab 08: Route Groups & Nested Middleware

**Learning Objectives:**

- Implement `Group()` for route prefixing
- Support nested groups with middleware inheritance
- Create modular route registration pattern
- Design clean folder structure for applications

**Tasks:**

| Task | Description |
| --- | --- |
| 8.1 | Implement `RouterGroup` struct |
| 8.2 | Add prefix handling and nesting |
| 8.3 | Implement middleware inheritance for groups |
| 8.4 | Create `Module` interface for modular routes |
| 8.5 | Test nested groups and middleware |

**Deliverables:**

```
goexpress/
├── group.go
├── group_test.go
└── module.go

```

**Usage Example:**

```go
api := app.Group("/api")
api.Use(authMiddleware)

v1 := api.Group("/v1")
v1.GET("/users", listUsers)      // → /api/v1/users

admin := api.Group("/admin")
admin.Use(adminMiddleware)
admin.GET("/stats", getStats)    // → /api/admin/stats

```

**Success Criteria:**

- [ ]  Groups prefix all child routes
- [ ]  Nested groups work correctly
- [ ]  Middleware inherited from parent groups
- [ ]  Module pattern allows clean separation

---

## Lab 09: REST API with CRUD Operations

**Learning Objectives:**

- Build complete REST API using the framework
- Implement CRUD handlers with proper patterns
- Add JWT authentication middleware
- Apply all framework features together

**Tasks:**

| Task | Description |
| --- | --- |
| 9.1 | Set up application structure |
| 9.2 | Implement User CRUD handlers |
| 9.3 | Implement Task CRUD handlers |
| 9.4 | Add JWT authentication middleware |
| 9.5 | Wire everything together |

**Application Structure:**

```
taskapi/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handler/
│   │   ├── user.go
│   │   └── task.go
│   ├── middleware/
│   │   └── auth.go
│   ├── model/
│   │   ├── user.go
│   │   └── task.go
│   └── store/
│       └── memory.go
└── go.mod

```

**API Endpoints:**

```
POST   /auth/login           # Get JWT token
GET    /api/users            # List users (protected)
POST   /api/users            # Create user
GET    /api/users/:id        # Get user
PUT    /api/users/:id        # Update user
DELETE /api/users/:id        # Delete user
GET    /api/users/:id/tasks  # List user's tasks
POST   /api/tasks            # Create task

```

**Success Criteria:**

- [ ]  All CRUD operations work correctly
- [ ]  JWT authentication protects routes
- [ ]  Validation applied to all inputs
- [ ]  Proper error responses returned

---

## Lab 10: Final Project - Complete Task Management API

**Learning Objectives:**

- Complete the full application with all features
- Add pagination, filtering, and sorting
- Implement comprehensive error handling
- Write API tests and documentation

**Tasks:**

| Task | Description |
| --- | --- |
| 10.1 | Add pagination to list endpoints |
| 10.2 | Add filtering and sorting |
| 10.3 | Implement comprehensive error handling |
| 10.4 | Write integration tests |
| 10.5 | Create API documentation |
| 10.6 | Final review and cleanup |

**Final Features:**

```
Pagination: GET /api/tasks?page=1&limit=20
Filtering:  GET /api/tasks?status=pending&priority=high
Sorting:    GET /api/tasks?sort=created_at&order=desc
Search:     GET /api/tasks?q=important

```

**Success Criteria:**

- [ ]  All 10 labs completed
- [ ]  Framework fully functional
- [ ]  REST API works end-to-end
- [ ]  Tests pass
- [ ]  Documentation complete

---

## End Project: Task Management System

Using your **GoExpress** framework, build a complete Task Management API.

### Final Project Structure

```
task-manager/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handler/
│   │   ├── auth.go
│   │   ├── user.go
│   │   ├── project.go
│   │   └── task.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── logger.go
│   │   └── recovery.go
│   ├── model/
│   │   ├── user.go
│   │   ├── project.go
│   │   └── task.go
│   ├── store/
│   │   └── memory.go
│   └── validator/
│       └── validator.go
├── pkg/
│   └── goexpress/          # Your framework
│       ├── goexpress.go
│       ├── context.go
│       ├── router.go
│       ├── group.go
│       ├── middleware.go
│       ├── binding.go
│       ├── render.go
│       └── validator.go
└── go.mod

```

### API Specification

**Authentication**

```
POST /auth/register    → Register new user
POST /auth/login       → Login and get JWT
POST /auth/refresh     → Refresh JWT token

```

**Users**

```
GET    /api/users          → List all users
GET    /api/users/:id      → Get user by ID
PUT    /api/users/:id      → Update user
DELETE /api/users/:id      → Delete user

```

**Projects**

```
GET    /api/projects           → List user's projects
POST   /api/projects           → Create project
GET    /api/projects/:id       → Get project
PUT    /api/projects/:id       → Update project
DELETE /api/projects/:id       → Delete project

```

**Tasks**

```
GET    /api/projects/:id/tasks     → List tasks in project
POST   /api/projects/:id/tasks     → Create task
GET    /api/tasks/:id              → Get task
PUT    /api/tasks/:id              → Update task
PATCH  /api/tasks/:id/status       → Update status only
DELETE /api/tasks/:id              → Delete task

```

### Required Features

**Framework Features Used:**

- Router with path parameters
- Route groups with prefix
- Middleware pipeline (Logger, Recovery, Auth)
- JSON binding and response
- Request validation
- Error handling

**Application Features:**

- JWT authentication
- CRUD for all resources
- Pagination on list endpoints
- Filtering by status/priority
- Sorting by date/priority
- Input validation
- Proper HTTP status codes
- Structured error responses

### 

---
# Create New Module

Generate a new module for the zciti-back project with all necessary files.

## Required Information
- Module name (singular, lowercase): e.g., `ticket`, `payment`, `category`
- Resource name (plural, for routes): e.g., `tickets`, `payments`, `categories`
- Fields for the schema (name, type, validations)

## Files to Generate

1. `app/database/schema/{module}.go` - GORM model
2. `app/module/{module}/index.go` - Module registration & routes
3. `app/module/{module}/controller/index.go` - Controller aggregator
4. `app/module/{module}/controller/rest_controller.go` - HTTP handlers
5. `app/module/{module}/service/index.go` - Business logic
6. `app/module/{module}/repository/index.go` - Data access
7. `app/module/{module}/request/index.go` - Request DTOs
8. `app/module/{module}/response/index.go` - Response DTOs

## After Generation

1. Register schema in `app/database/schema/index.go`
2. Import module in `cmd/main/main.go`
3. Add to fx.Options in main.go
4. Add router to `app/router/api.go`
5. Register routes in Router.Register()

## Example Usage

"Create a new module called 'ticket' with fields: title (string, required), description (text), status (enum: open/closed/pending), priority (int), userID (relation)"


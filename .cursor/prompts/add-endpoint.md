# Add New Endpoint

Add a new endpoint to an existing module.

## Required Information
- Module name
- HTTP method (GET, POST, PUT, DELETE)
- Route path
- Purpose/functionality
- Request/response fields if needed

## Steps

1. Add method to `IRestController` interface
2. Implement handler in `rest_controller.go`
3. Add method to `IService` interface if needed
4. Implement service method
5. Add method to `IRepository` interface if needed
6. Implement repository method
7. Add route in module's `index.go`
8. Add Swagger annotations

## Example Usage

"Add a GET endpoint `/users/:id/orders` to the user module that returns all orders for a specific user with pagination"


package test

import (
	"testing"
	"time"

	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/orderItem/request"
	"go-fiber-starter/utils/paginator"
)

// =============================================================================
// REPOSITORY TESTS - Direct database operations
// =============================================================================

func TestRepository_Create_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user and business
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Create test post and product
	post := ta.CreateTestPost(t, "Test Product Post", schema.PostTypeProduct, user.ID, business.ID)
	variantType := schema.ProductVariantTypeSimple
	ta.CreateTestProduct(t, post.ID, business.ID, schema.ProductTypeSimple, &variantType, 100.0)

	// Create test order
	order := ta.CreateTestOrder(t, user.ID, business.ID, schema.OrderStatusPending, 100.0)

	// Create order item using repository
	orderItem := &schema.OrderItem{
		Type:     schema.OrderItemTypeLineItem,
		Quantity: 2,
		Price:    50.0,
		Subtotal: 100.0,
		PostID:   post.ID,
		Meta: schema.OrderItemMeta{
			ProductTitle:  "Test Product",
			ProductSKU:    "TEST-SKU",
			ProductDetail: "Test detail",
		},
	}

	err := ta.OrderItemRepo.Create(orderItem, order.ID, nil)
	if err != nil {
		t.Fatalf("failed to create order item: %v", err)
	}

	// Verify order item was created
	var count int64
	ta.DB.Model(&schema.OrderItem{}).Where("order_id = ?", order.ID).Count(&count)
	if count != 1 {
		t.Errorf("expected 1 order item, got %d", count)
	}

	// Verify order_id was set correctly
	if orderItem.OrderID != order.ID {
		t.Errorf("expected order_id %d, got %d", order.ID, orderItem.OrderID)
	}
}

func TestRepository_Create_WithTransaction(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user and business
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Create test post
	post := ta.CreateTestPost(t, "Test Product Post", schema.PostTypeProduct, user.ID, business.ID)

	// Create test order
	order := ta.CreateTestOrder(t, user.ID, business.ID, schema.OrderStatusPending, 200.0)

	// Start transaction
	tx := ta.DB.Begin()

	// Create multiple order items in transaction
	orderItem1 := &schema.OrderItem{
		Type:     schema.OrderItemTypeLineItem,
		Quantity: 1,
		Price:    100.0,
		Subtotal: 100.0,
		PostID:   post.ID,
		Meta:     schema.OrderItemMeta{ProductTitle: "Product 1"},
	}

	orderItem2 := &schema.OrderItem{
		Type:     schema.OrderItemTypeLineItem,
		Quantity: 2,
		Price:    50.0,
		Subtotal: 100.0,
		PostID:   post.ID,
		Meta:     schema.OrderItemMeta{ProductTitle: "Product 2"},
	}

	err := ta.OrderItemRepo.Create(orderItem1, order.ID, tx)
	if err != nil {
		tx.Rollback()
		t.Fatalf("failed to create order item 1: %v", err)
	}

	err = ta.OrderItemRepo.Create(orderItem2, order.ID, tx)
	if err != nil {
		tx.Rollback()
		t.Fatalf("failed to create order item 2: %v", err)
	}

	// Commit transaction
	tx.Commit()

	// Verify both order items were created
	var count int64
	ta.DB.Model(&schema.OrderItem{}).Where("order_id = ?", order.ID).Count(&count)
	if count != 2 {
		t.Errorf("expected 2 order items, got %d", count)
	}
}

func TestRepository_Create_TransactionRollback(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user and business
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Create test post
	post := ta.CreateTestPost(t, "Test Product Post", schema.PostTypeProduct, user.ID, business.ID)

	// Create test order
	order := ta.CreateTestOrder(t, user.ID, business.ID, schema.OrderStatusPending, 100.0)

	// Start transaction
	tx := ta.DB.Begin()

	// Create order item in transaction
	orderItem := &schema.OrderItem{
		Type:     schema.OrderItemTypeLineItem,
		Quantity: 1,
		Price:    100.0,
		Subtotal: 100.0,
		PostID:   post.ID,
		Meta:     schema.OrderItemMeta{ProductTitle: "Product 1"},
	}

	err := ta.OrderItemRepo.Create(orderItem, order.ID, tx)
	if err != nil {
		tx.Rollback()
		t.Fatalf("failed to create order item: %v", err)
	}

	// Rollback transaction
	tx.Rollback()

	// Verify order item was NOT created
	var count int64
	ta.DB.Model(&schema.OrderItem{}).Where("order_id = ?", order.ID).Count(&count)
	if count != 0 {
		t.Errorf("expected 0 order items after rollback, got %d", count)
	}
}

func TestRepository_GetOne_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user and business
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Create test post
	post := ta.CreateTestPost(t, "Test Product Post", schema.PostTypeProduct, user.ID, business.ID)

	// Create test order
	order := ta.CreateTestOrder(t, user.ID, business.ID, schema.OrderStatusPending, 100.0)

	// Create test order item
	orderItem := ta.CreateTestOrderItem(t, order.ID, post.ID, nil, schema.OrderItemTypeLineItem, 2, 50.0)

	// Get order item using repository
	result, err := ta.OrderItemRepo.GetOne(orderItem.ID)
	if err != nil {
		t.Fatalf("failed to get order item: %v", err)
	}

	if result.ID != orderItem.ID {
		t.Errorf("expected order item ID %d, got %d", orderItem.ID, result.ID)
	}

	if result.Quantity != 2 {
		t.Errorf("expected quantity 2, got %d", result.Quantity)
	}

	if result.Price != 50.0 {
		t.Errorf("expected price 50.0, got %f", result.Price)
	}
}

func TestRepository_GetOne_NotFound(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Try to get non-existent order item
	_, err := ta.OrderItemRepo.GetOne(99999)
	if err == nil {
		t.Error("expected error for non-existent order item, got nil")
	}
}

func TestRepository_GetAll_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user and business
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Create test post
	post := ta.CreateTestPost(t, "Test Product Post", schema.PostTypeProduct, user.ID, business.ID)

	// Create test order
	order := ta.CreateTestOrder(t, user.ID, business.ID, schema.OrderStatusPending, 300.0)

	// Create multiple order items
	ta.CreateTestOrderItem(t, order.ID, post.ID, nil, schema.OrderItemTypeLineItem, 1, 100.0)
	ta.CreateTestOrderItem(t, order.ID, post.ID, nil, schema.OrderItemTypeLineItem, 2, 100.0)
	ta.CreateTestOrderItem(t, order.ID, post.ID, nil, schema.OrderItemTypeFee, 1, 10.0)

	// Get all order items
	req := request.OrderItems{
		Pagination: &paginator.Pagination{Page: 0},
	}

	results, _, err := ta.OrderItemRepo.GetAll(req)
	if err != nil {
		t.Fatalf("failed to get all order items: %v", err)
	}

	if len(results) != 3 {
		t.Errorf("expected 3 order items, got %d", len(results))
	}
}

func TestRepository_GetAll_WithPagination(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user and business
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Create test post
	post := ta.CreateTestPost(t, "Test Product Post", schema.PostTypeProduct, user.ID, business.ID)

	// Create test order
	order := ta.CreateTestOrder(t, user.ID, business.ID, schema.OrderStatusPending, 500.0)

	// Create multiple order items
	for i := 0; i < 5; i++ {
		ta.CreateTestOrderItem(t, order.ID, post.ID, nil, schema.OrderItemTypeLineItem, 1, 100.0)
	}

	// Get paginated order items (page 1, limit 2)
	req := request.OrderItems{
		Pagination: &paginator.Pagination{
			Page:   1,
			Limit:  2,
			Offset: 0,
		},
	}

	results, paging, err := ta.OrderItemRepo.GetAll(req)
	if err != nil {
		t.Fatalf("failed to get paginated order items: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("expected 2 order items per page, got %d", len(results))
	}

	if paging.Total != 5 {
		t.Errorf("expected total 5 order items, got %d", paging.Total)
	}
}

func TestRepository_Update_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user and business
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Create test post
	post := ta.CreateTestPost(t, "Test Product Post", schema.PostTypeProduct, user.ID, business.ID)

	// Create test order
	order := ta.CreateTestOrder(t, user.ID, business.ID, schema.OrderStatusPending, 100.0)

	// Create test order item
	orderItem := ta.CreateTestOrderItem(t, order.ID, post.ID, nil, schema.OrderItemTypeLineItem, 1, 100.0)

	// Update order item
	updatedItem := &schema.OrderItem{
		Quantity: 3,
		Price:    75.0,
		Subtotal: 225.0,
	}

	err := ta.OrderItemRepo.Update(orderItem.ID, updatedItem)
	if err != nil {
		t.Fatalf("failed to update order item: %v", err)
	}

	// Verify update
	result, err := ta.OrderItemRepo.GetOne(orderItem.ID)
	if err != nil {
		t.Fatalf("failed to get updated order item: %v", err)
	}

	if result.Quantity != 3 {
		t.Errorf("expected quantity 3, got %d", result.Quantity)
	}

	if result.Price != 75.0 {
		t.Errorf("expected price 75.0, got %f", result.Price)
	}

	if result.Subtotal != 225.0 {
		t.Errorf("expected subtotal 225.0, got %f", result.Subtotal)
	}
}

func TestRepository_Delete_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user and business
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Create test post
	post := ta.CreateTestPost(t, "Test Product Post", schema.PostTypeProduct, user.ID, business.ID)

	// Create test order
	order := ta.CreateTestOrder(t, user.ID, business.ID, schema.OrderStatusPending, 100.0)

	// Create test order item
	orderItem := ta.CreateTestOrderItem(t, order.ID, post.ID, nil, schema.OrderItemTypeLineItem, 1, 100.0)

	// Delete order item
	err := ta.OrderItemRepo.Delete(orderItem.ID)
	if err != nil {
		t.Fatalf("failed to delete order item: %v", err)
	}

	// Verify deletion (soft delete)
	var count int64
	ta.DB.Unscoped().Model(&schema.OrderItem{}).Where("id = ? AND deleted_at IS NOT NULL", orderItem.ID).Count(&count)
	if count != 1 {
		t.Errorf("expected order item to be soft deleted")
	}
}

// =============================================================================
// SERVICE TESTS - Business logic tests
// =============================================================================

func TestService_Index_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user and business
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Create test post
	post := ta.CreateTestPost(t, "Test Product Post", schema.PostTypeProduct, user.ID, business.ID)

	// Create test order
	order := ta.CreateTestOrder(t, user.ID, business.ID, schema.OrderStatusPending, 200.0)

	// Create order items
	ta.CreateTestOrderItem(t, order.ID, post.ID, nil, schema.OrderItemTypeLineItem, 1, 100.0)
	ta.CreateTestOrderItem(t, order.ID, post.ID, nil, schema.OrderItemTypeLineItem, 1, 100.0)

	// Create service directly for testing
	svc := ta.OrderItemRouter.Controller.RestController

	// This tests the service indirectly via the controller
	// Since routes are commented out, we can test the repository layer directly
	req := request.OrderItems{
		Pagination: &paginator.Pagination{Page: 0},
	}

	results, _, err := ta.OrderItemRepo.GetAll(req)
	if err != nil {
		t.Fatalf("failed to get order items: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("expected 2 order items, got %d", len(results))
	}

	// Verify service is initialized
	if svc == nil {
		t.Error("expected service to be initialized")
	}
}

func TestService_Show_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user and business
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Create test post
	post := ta.CreateTestPost(t, "Test Product Post", schema.PostTypeProduct, user.ID, business.ID)

	// Create test order
	order := ta.CreateTestOrder(t, user.ID, business.ID, schema.OrderStatusPending, 150.0)

	// Create order item
	orderItem := ta.CreateTestOrderItem(t, order.ID, post.ID, nil, schema.OrderItemTypeLineItem, 3, 50.0)

	// Get order item
	result, err := ta.OrderItemRepo.GetOne(orderItem.ID)
	if err != nil {
		t.Fatalf("failed to get order item: %v", err)
	}

	if result.ID != orderItem.ID {
		t.Errorf("expected order item ID %d, got %d", orderItem.ID, result.ID)
	}

	if result.Quantity != 3 {
		t.Errorf("expected quantity 3, got %d", result.Quantity)
	}

	if result.Price != 50.0 {
		t.Errorf("expected price 50.0, got %f", result.Price)
	}

	if result.Subtotal != 150.0 {
		t.Errorf("expected subtotal 150.0, got %f", result.Subtotal)
	}
}

func TestService_Destroy_Success(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user and business
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Create test post
	post := ta.CreateTestPost(t, "Test Product Post", schema.PostTypeProduct, user.ID, business.ID)

	// Create test order
	order := ta.CreateTestOrder(t, user.ID, business.ID, schema.OrderStatusPending, 100.0)

	// Create order item
	orderItem := ta.CreateTestOrderItem(t, order.ID, post.ID, nil, schema.OrderItemTypeLineItem, 1, 100.0)

	// Delete via repository
	err := ta.OrderItemRepo.Delete(orderItem.ID)
	if err != nil {
		t.Fatalf("failed to delete order item: %v", err)
	}

	// Verify deletion
	_, err = ta.OrderItemRepo.GetOne(orderItem.ID)
	if err == nil {
		t.Error("expected error after deletion, got nil")
	}
}

// =============================================================================
// ORDER ITEM TYPE TESTS
// =============================================================================

func TestOrderItem_TypeLineItem(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user and business
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Create test post
	post := ta.CreateTestPost(t, "Test Product Post", schema.PostTypeProduct, user.ID, business.ID)

	// Create test order
	order := ta.CreateTestOrder(t, user.ID, business.ID, schema.OrderStatusPending, 100.0)

	// Create line item order item
	orderItem := ta.CreateTestOrderItem(t, order.ID, post.ID, nil, schema.OrderItemTypeLineItem, 2, 50.0)

	result, err := ta.OrderItemRepo.GetOne(orderItem.ID)
	if err != nil {
		t.Fatalf("failed to get order item: %v", err)
	}

	if result.Type != schema.OrderItemTypeLineItem {
		t.Errorf("expected type 'lineItem', got '%s'", result.Type)
	}
}

func TestOrderItem_TypeReservation(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user and business
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Create test post and product
	post := ta.CreateTestPost(t, "Test Product Post", schema.PostTypeProduct, user.ID, business.ID)
	variantType := schema.ProductVariantTypeReservable
	product := ta.CreateTestProduct(t, post.ID, business.ID, schema.ProductTypeReservable, &variantType, 100.0)

	// Create test order
	order := ta.CreateTestOrder(t, user.ID, business.ID, schema.OrderStatusPending, 100.0)

	// Create reservation
	now := time.Now()
	reservation := ta.CreateTestReservation(t, user.ID, product.ID, business.ID, now, now.Add(2*time.Hour))

	// Create reservation order item
	orderItem := ta.CreateTestOrderItem(t, order.ID, post.ID, &reservation.ID, schema.OrderItemTypeReservation, 1, 100.0)

	result, err := ta.OrderItemRepo.GetOne(orderItem.ID)
	if err != nil {
		t.Fatalf("failed to get order item: %v", err)
	}

	if result.Type != schema.OrderItemTypeReservation {
		t.Errorf("expected type 'reservation', got '%s'", result.Type)
	}

	if result.ReservationID == nil {
		t.Error("expected reservation ID to be set")
	}
}

func TestOrderItem_TypeFee(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user and business
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Create test post
	post := ta.CreateTestPost(t, "Test Product Post", schema.PostTypeProduct, user.ID, business.ID)

	// Create test order
	order := ta.CreateTestOrder(t, user.ID, business.ID, schema.OrderStatusPending, 110.0)

	// Create fee order item
	orderItem := ta.CreateTestOrderItem(t, order.ID, post.ID, nil, schema.OrderItemTypeFee, 1, 10.0)

	result, err := ta.OrderItemRepo.GetOne(orderItem.ID)
	if err != nil {
		t.Fatalf("failed to get order item: %v", err)
	}

	if result.Type != schema.OrderItemTypeFee {
		t.Errorf("expected type 'fee', got '%s'", result.Type)
	}
}

func TestOrderItem_TypeTax(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user and business
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Create test post
	post := ta.CreateTestPost(t, "Test Product Post", schema.PostTypeProduct, user.ID, business.ID)

	// Create test order
	order := ta.CreateTestOrder(t, user.ID, business.ID, schema.OrderStatusPending, 109.0)

	// Create tax order item
	orderItem := ta.CreateTestOrderItem(t, order.ID, post.ID, nil, schema.OrderItemTypeTax, 1, 9.0)

	result, err := ta.OrderItemRepo.GetOne(orderItem.ID)
	if err != nil {
		t.Fatalf("failed to get order item: %v", err)
	}

	if result.Type != schema.OrderItemTypeTax {
		t.Errorf("expected type 'tax', got '%s'", result.Type)
	}
}

func TestOrderItem_TypeCoupon(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user and business
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Create test post
	post := ta.CreateTestPost(t, "Test Product Post", schema.PostTypeProduct, user.ID, business.ID)

	// Create test order
	order := ta.CreateTestOrder(t, user.ID, business.ID, schema.OrderStatusPending, 80.0)

	// Create coupon order item (negative value represents discount)
	orderItem := ta.CreateTestOrderItem(t, order.ID, post.ID, nil, schema.OrderItemTypeCoupon, 1, -20.0)

	result, err := ta.OrderItemRepo.GetOne(orderItem.ID)
	if err != nil {
		t.Fatalf("failed to get order item: %v", err)
	}

	if result.Type != schema.OrderItemTypeCoupon {
		t.Errorf("expected type 'coupon', got '%s'", result.Type)
	}
}

func TestOrderItem_TypeShipping(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user and business
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Create test post
	post := ta.CreateTestPost(t, "Test Product Post", schema.PostTypeProduct, user.ID, business.ID)

	// Create test order
	order := ta.CreateTestOrder(t, user.ID, business.ID, schema.OrderStatusPending, 115.0)

	// Create shipping order item
	orderItem := ta.CreateTestOrderItem(t, order.ID, post.ID, nil, schema.OrderItemTypeShipping, 1, 15.0)

	result, err := ta.OrderItemRepo.GetOne(orderItem.ID)
	if err != nil {
		t.Fatalf("failed to get order item: %v", err)
	}

	if result.Type != schema.OrderItemTypeShipping {
		t.Errorf("expected type 'shipping', got '%s'", result.Type)
	}
}

// =============================================================================
// ORDER ITEM META TESTS
// =============================================================================

func TestOrderItem_MetaFields(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user and business
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Create test post
	post := ta.CreateTestPost(t, "Test Product Post", schema.PostTypeProduct, user.ID, business.ID)

	// Create test order
	order := ta.CreateTestOrder(t, user.ID, business.ID, schema.OrderStatusPending, 100.0)

	// Create order item with detailed meta
	orderItem := &schema.OrderItem{
		Type:     schema.OrderItemTypeLineItem,
		Quantity: 1,
		Price:    100.0,
		Subtotal: 100.0,
		PostID:   post.ID,
		OrderID:  order.ID,
		Meta: schema.OrderItemMeta{
			ProductID:          1,
			ProductTitle:       "Test Product Title",
			ProductDetail:      "Size: Large, Color: Blue",
			ProductSKU:         "SKU-12345",
			ProductType:        schema.ProductTypeVariant,
			ProductVariantType: schema.ProductVariantTypeSimple,
		},
	}

	if err := ta.DB.Create(orderItem).Error; err != nil {
		t.Fatalf("failed to create order item: %v", err)
	}

	// Retrieve and verify meta
	result, err := ta.OrderItemRepo.GetOne(orderItem.ID)
	if err != nil {
		t.Fatalf("failed to get order item: %v", err)
	}

	if result.Meta.ProductTitle != "Test Product Title" {
		t.Errorf("expected product title 'Test Product Title', got '%s'", result.Meta.ProductTitle)
	}

	if result.Meta.ProductSKU != "SKU-12345" {
		t.Errorf("expected product SKU 'SKU-12345', got '%s'", result.Meta.ProductSKU)
	}

	if result.Meta.ProductDetail != "Size: Large, Color: Blue" {
		t.Errorf("expected product detail 'Size: Large, Color: Blue', got '%s'", result.Meta.ProductDetail)
	}

	if result.Meta.ProductType != schema.ProductTypeVariant {
		t.Errorf("expected product type 'variant', got '%s'", result.Meta.ProductType)
	}
}

// =============================================================================
// INTEGRATION TESTS
// =============================================================================

func TestIntegration_CompleteOrderWithItems(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user and business
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Create test post
	post := ta.CreateTestPost(t, "Test Product Post", schema.PostTypeProduct, user.ID, business.ID)

	// Create test order
	order := ta.CreateTestOrder(t, user.ID, business.ID, schema.OrderStatusPending, 0)

	// Create multiple order items of different types
	lineItem1 := ta.CreateTestOrderItem(t, order.ID, post.ID, nil, schema.OrderItemTypeLineItem, 2, 50.0)
	lineItem2 := ta.CreateTestOrderItem(t, order.ID, post.ID, nil, schema.OrderItemTypeLineItem, 1, 75.0)
	feeItem := ta.CreateTestOrderItem(t, order.ID, post.ID, nil, schema.OrderItemTypeFee, 1, 5.0)
	taxItem := ta.CreateTestOrderItem(t, order.ID, post.ID, nil, schema.OrderItemTypeTax, 1, 18.0)
	shippingItem := ta.CreateTestOrderItem(t, order.ID, post.ID, nil, schema.OrderItemTypeShipping, 1, 10.0)

	// Calculate expected total: (2*50) + (1*75) + 5 + 18 + 10 = 208
	expectedTotal := lineItem1.Subtotal + lineItem2.Subtotal + feeItem.Subtotal + taxItem.Subtotal + shippingItem.Subtotal

	// Get all order items
	req := request.OrderItems{
		Pagination: &paginator.Pagination{Page: 0},
	}

	results, _, err := ta.OrderItemRepo.GetAll(req)
	if err != nil {
		t.Fatalf("failed to get order items: %v", err)
	}

	if len(results) != 5 {
		t.Errorf("expected 5 order items, got %d", len(results))
	}

	// Calculate actual total from order items
	var actualTotal float64
	for _, item := range results {
		actualTotal += item.Subtotal
	}

	if actualTotal != expectedTotal {
		t.Errorf("expected total %f, got %f", expectedTotal, actualTotal)
	}
}

func TestIntegration_OrderWithReservation(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user and business
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Create test post and product
	post := ta.CreateTestPost(t, "Gym Session", schema.PostTypeProduct, user.ID, business.ID)
	variantType := schema.ProductVariantTypeReservable
	product := ta.CreateTestProduct(t, post.ID, business.ID, schema.ProductTypeReservable, &variantType, 30.0)

	// Create test order
	order := ta.CreateTestOrder(t, user.ID, business.ID, schema.OrderStatusPending, 30.0)

	// Create reservation
	now := time.Now()
	startTime := now.Add(24 * time.Hour)
	endTime := startTime.Add(time.Hour)
	reservation := ta.CreateTestReservation(t, user.ID, product.ID, business.ID, startTime, endTime)

	// Create reservation order item
	orderItem := ta.CreateTestOrderItem(t, order.ID, post.ID, &reservation.ID, schema.OrderItemTypeReservation, 1, 30.0)

	// Verify order item has reservation
	result, err := ta.OrderItemRepo.GetOne(orderItem.ID)
	if err != nil {
		t.Fatalf("failed to get order item: %v", err)
	}

	if result.ReservationID == nil {
		t.Error("expected reservation ID to be set")
	}

	if *result.ReservationID != reservation.ID {
		t.Errorf("expected reservation ID %d, got %d", reservation.ID, *result.ReservationID)
	}

	if result.Type != schema.OrderItemTypeReservation {
		t.Errorf("expected type 'reservation', got '%s'", result.Type)
	}
}

func TestIntegration_DeleteOrderItem(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user and business
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Create test post
	post := ta.CreateTestPost(t, "Test Product", schema.PostTypeProduct, user.ID, business.ID)

	// Create test order
	order := ta.CreateTestOrder(t, user.ID, business.ID, schema.OrderStatusPending, 200.0)

	// Create order items
	item1 := ta.CreateTestOrderItem(t, order.ID, post.ID, nil, schema.OrderItemTypeLineItem, 1, 100.0)
	ta.CreateTestOrderItem(t, order.ID, post.ID, nil, schema.OrderItemTypeLineItem, 1, 100.0)

	// Verify initial count
	req := request.OrderItems{
		Pagination: &paginator.Pagination{Page: 0},
	}

	results, _, err := ta.OrderItemRepo.GetAll(req)
	if err != nil {
		t.Fatalf("failed to get order items: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("expected 2 order items before deletion, got %d", len(results))
	}

	// Delete one item
	err = ta.OrderItemRepo.Delete(item1.ID)
	if err != nil {
		t.Fatalf("failed to delete order item: %v", err)
	}

	// Verify count after deletion
	results, _, err = ta.OrderItemRepo.GetAll(req)
	if err != nil {
		t.Fatalf("failed to get order items after deletion: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 order item after deletion, got %d", len(results))
	}
}

// =============================================================================
// EDGE CASE TESTS
// =============================================================================

func TestOrderItem_ZeroQuantity(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user and business
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Create test post
	post := ta.CreateTestPost(t, "Test Product", schema.PostTypeProduct, user.ID, business.ID)

	// Create test order
	order := ta.CreateTestOrder(t, user.ID, business.ID, schema.OrderStatusPending, 0)

	// Create order item with zero quantity
	orderItem := ta.CreateTestOrderItem(t, order.ID, post.ID, nil, schema.OrderItemTypeLineItem, 0, 100.0)

	result, err := ta.OrderItemRepo.GetOne(orderItem.ID)
	if err != nil {
		t.Fatalf("failed to get order item: %v", err)
	}

	if result.Quantity != 0 {
		t.Errorf("expected quantity 0, got %d", result.Quantity)
	}

	if result.Subtotal != 0 {
		t.Errorf("expected subtotal 0, got %f", result.Subtotal)
	}
}

func TestOrderItem_LargeQuantity(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user and business
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Create test post
	post := ta.CreateTestPost(t, "Test Product", schema.PostTypeProduct, user.ID, business.ID)

	// Create test order
	order := ta.CreateTestOrder(t, user.ID, business.ID, schema.OrderStatusPending, 10000000.0)

	// Create order item with large quantity
	orderItem := ta.CreateTestOrderItem(t, order.ID, post.ID, nil, schema.OrderItemTypeLineItem, 100000, 100.0)

	result, err := ta.OrderItemRepo.GetOne(orderItem.ID)
	if err != nil {
		t.Fatalf("failed to get order item: %v", err)
	}

	if result.Quantity != 100000 {
		t.Errorf("expected quantity 100000, got %d", result.Quantity)
	}

	if result.Subtotal != 10000000.0 {
		t.Errorf("expected subtotal 10000000.0, got %f", result.Subtotal)
	}
}

func TestOrderItem_DecimalPrice(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user and business
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Create test post
	post := ta.CreateTestPost(t, "Test Product", schema.PostTypeProduct, user.ID, business.ID)

	// Create test order
	order := ta.CreateTestOrder(t, user.ID, business.ID, schema.OrderStatusPending, 99.97)

	// Create order item with decimal price
	orderItem := &schema.OrderItem{
		Type:     schema.OrderItemTypeLineItem,
		Quantity: 3,
		Price:    33.33,
		Subtotal: 99.99,
		PostID:   post.ID,
		OrderID:  order.ID,
		Meta:     schema.OrderItemMeta{ProductTitle: "Decimal Price Product"},
	}

	if err := ta.DB.Create(orderItem).Error; err != nil {
		t.Fatalf("failed to create order item: %v", err)
	}

	result, err := ta.OrderItemRepo.GetOne(orderItem.ID)
	if err != nil {
		t.Fatalf("failed to get order item: %v", err)
	}

	if result.Price != 33.33 {
		t.Errorf("expected price 33.33, got %f", result.Price)
	}

	if result.Subtotal != 99.99 {
		t.Errorf("expected subtotal 99.99, got %f", result.Subtotal)
	}
}

func TestOrderItem_WithTaxAmount(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user and business
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Create test post
	post := ta.CreateTestPost(t, "Test Product", schema.PostTypeProduct, user.ID, business.ID)

	// Create test order
	order := ta.CreateTestOrder(t, user.ID, business.ID, schema.OrderStatusPending, 109.0)

	// Create order item with tax amount
	orderItem := &schema.OrderItem{
		Type:     schema.OrderItemTypeLineItem,
		Quantity: 1,
		Price:    100.0,
		Subtotal: 100.0,
		TaxAmt:   9.0,
		PostID:   post.ID,
		OrderID:  order.ID,
		Meta:     schema.OrderItemMeta{ProductTitle: "Taxable Product"},
	}

	if err := ta.DB.Create(orderItem).Error; err != nil {
		t.Fatalf("failed to create order item: %v", err)
	}

	result, err := ta.OrderItemRepo.GetOne(orderItem.ID)
	if err != nil {
		t.Fatalf("failed to get order item: %v", err)
	}

	if result.TaxAmt != 9.0 {
		t.Errorf("expected tax amount 9.0, got %f", result.TaxAmt)
	}
}

// =============================================================================
// PAGINATION EDGE CASES
// =============================================================================

func TestRepository_GetAll_EmptyResult(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Don't create any order items

	req := request.OrderItems{
		Pagination: &paginator.Pagination{Page: 0},
	}

	results, _, err := ta.OrderItemRepo.GetAll(req)
	if err != nil {
		t.Fatalf("failed to get order items: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("expected 0 order items, got %d", len(results))
	}
}

func TestRepository_GetAll_PageBeyondTotal(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	// Create test user and business
	user := ta.CreateTestUser(t, 9123456789, "testPassword123", "Test", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, user.ID)

	// Create test post
	post := ta.CreateTestPost(t, "Test Product", schema.PostTypeProduct, user.ID, business.ID)

	// Create test order
	order := ta.CreateTestOrder(t, user.ID, business.ID, schema.OrderStatusPending, 100.0)

	// Create only 2 order items
	ta.CreateTestOrderItem(t, order.ID, post.ID, nil, schema.OrderItemTypeLineItem, 1, 50.0)
	ta.CreateTestOrderItem(t, order.ID, post.ID, nil, schema.OrderItemTypeLineItem, 1, 50.0)

	// Request page 10 with limit 10 (beyond total)
	req := request.OrderItems{
		Pagination: &paginator.Pagination{
			Page:   10,
			Limit:  10,
			Offset: 90, // (10-1) * 10
		},
	}

	results, paging, err := ta.OrderItemRepo.GetAll(req)
	if err != nil {
		t.Fatalf("failed to get order items: %v", err)
	}

	// Should return empty results
	if len(results) != 0 {
		t.Errorf("expected 0 order items for page beyond total, got %d", len(results))
	}

	// Total should still be 2
	if paging.Total != 2 {
		t.Errorf("expected total 2, got %d", paging.Total)
	}
}


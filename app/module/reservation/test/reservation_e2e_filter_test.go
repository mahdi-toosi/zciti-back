package test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"go-fiber-starter/app/database/schema"
)

// =============================================================================
// TAXONOMY FILTER TESTS
// =============================================================================

func TestIndex_FilterByTaxonomies(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	setup := ta.SetupTestUser(t)

	// Create taxonomies (city -> workspace -> dormitory)
	city := ta.CreateTestTaxonomy(t, "Test City", schema.TaxonomyTypeCategory, setup.Business.ID, nil)
	workspace := ta.CreateTestTaxonomy(t, "Test Workspace", schema.TaxonomyTypeCategory, setup.Business.ID, &city.ID)
	dormitory := ta.CreateTestTaxonomy(t, "Test Dormitory", schema.TaxonomyTypeCategory, setup.Business.ID, &workspace.ID)

	// Create test posts and products
	post1 := ta.CreateTestPost(t, "Product 1", schema.PostTypeProduct, setup.Business.ID, setup.User.ID)
	post2 := ta.CreateTestPost(t, "Product 2", schema.PostTypeProduct, setup.Business.ID, setup.User.ID)
	variantType := schema.ProductVariantTypeWashingMachine
	product1 := ta.CreateTestProduct(t, post1.ID, setup.Business.ID, 10000, schema.ProductTypeVariant, &variantType)
	product2 := ta.CreateTestProduct(t, post2.ID, setup.Business.ID, 15000, schema.ProductTypeVariant, &variantType)

	// Attach taxonomies to posts
	ta.AttachTaxonomyToPost(t, post1.ID, dormitory.ID)
	ta.AttachTaxonomyToPost(t, post2.ID, workspace.ID)

	// Create reservations for different products
	now := time.Now()
	ta.CreateTestReservation(t, setup.User.ID, product1.ID, setup.Business.ID, now.Add(time.Hour), now.Add(2*time.Hour), schema.ReservationStatusReserved)
	ta.CreateTestReservation(t, setup.User.ID, product2.ID, setup.Business.ID, now.Add(3*time.Hour), now.Add(4*time.Hour), schema.ReservationStatusReserved)

	// Test filter by DormitoryID - should only get product1's reservation
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/reservations?DormitoryID=%d", setup.Business.ID, dormitory.ID), nil, setup.Token)
	AssertOK(t, resp)
	AssertDataCount(t, resp, 1)
}

func TestIndex_FilterByCityID(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	setup := ta.SetupTestUser(t)

	// Create taxonomies
	city1 := ta.CreateTestTaxonomy(t, "City 1", schema.TaxonomyTypeCategory, setup.Business.ID, nil)
	city2 := ta.CreateTestTaxonomy(t, "City 2", schema.TaxonomyTypeCategory, setup.Business.ID, nil)

	// Create test posts and products
	post1 := ta.CreateTestPost(t, "Product 1", schema.PostTypeProduct, setup.Business.ID, setup.User.ID)
	post2 := ta.CreateTestPost(t, "Product 2", schema.PostTypeProduct, setup.Business.ID, setup.User.ID)
	variantType := schema.ProductVariantTypeWashingMachine
	product1 := ta.CreateTestProduct(t, post1.ID, setup.Business.ID, 10000, schema.ProductTypeVariant, &variantType)
	product2 := ta.CreateTestProduct(t, post2.ID, setup.Business.ID, 15000, schema.ProductTypeVariant, &variantType)

	// Attach taxonomies to posts
	ta.AttachTaxonomyToPost(t, post1.ID, city1.ID)
	ta.AttachTaxonomyToPost(t, post2.ID, city2.ID)

	// Create reservations
	now := time.Now()
	ta.CreateTestReservation(t, setup.User.ID, product1.ID, setup.Business.ID, now.Add(time.Hour), now.Add(2*time.Hour), schema.ReservationStatusReserved)
	ta.CreateTestReservation(t, setup.User.ID, product2.ID, setup.Business.ID, now.Add(3*time.Hour), now.Add(4*time.Hour), schema.ReservationStatusReserved)

	// Filter by CityID
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/reservations?CityID=%d", setup.Business.ID, city1.ID), nil, setup.Token)
	AssertOK(t, resp)
	AssertDataCount(t, resp, 1)
}

// =============================================================================
// BUSINESS OBSERVER TESTS
// =============================================================================

func TestIndex_BusinessObserver_WithMeta(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	owner := ta.CreateTestUser(t, 9123456789, "testPassword123", "Owner", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, owner.ID)
	owner.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(owner)

	// Create taxonomies
	city1 := ta.CreateTestTaxonomy(t, "City 1", schema.TaxonomyTypeCategory, business.ID, nil)
	city2 := ta.CreateTestTaxonomy(t, "City 2", schema.TaxonomyTypeCategory, business.ID, nil)

	// Create test posts and products
	post1 := ta.CreateTestPost(t, "Product 1", schema.PostTypeProduct, business.ID, owner.ID)
	post2 := ta.CreateTestPost(t, "Product 2", schema.PostTypeProduct, business.ID, owner.ID)
	variantType := schema.ProductVariantTypeWashingMachine
	product1 := ta.CreateTestProduct(t, post1.ID, business.ID, 10000, schema.ProductTypeVariant, &variantType)
	product2 := ta.CreateTestProduct(t, post2.ID, business.ID, 15000, schema.ProductTypeVariant, &variantType)

	// Attach taxonomies to posts
	ta.AttachTaxonomyToPost(t, post1.ID, city1.ID)
	ta.AttachTaxonomyToPost(t, post2.ID, city2.ID)

	// Create reservations
	now := time.Now()
	ta.CreateTestReservation(t, owner.ID, product1.ID, business.ID, now.Add(time.Hour), now.Add(2*time.Hour), schema.ReservationStatusReserved)
	ta.CreateTestReservation(t, owner.ID, product2.ID, business.ID, now.Add(3*time.Hour), now.Add(4*time.Hour), schema.ReservationStatusReserved)

	// Create observer with access to city1 only
	observerMeta := &schema.UserMeta{
		TaxonomiesToObserve: schema.UserMetaTaxonomiesToObserve{
			city1.ID: {Checked: true, PartialChecked: false},
		},
	}
	observer := ta.CreateTestUserWithMeta(t, 9987654321, "testPassword456", "Observer", "User", business.ID, []schema.UserRole{schema.URBusinessObserver}, observerMeta)

	token := ta.GenerateTestToken(t, observer)

	// Observer should only see reservations for city1
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/reservations", business.ID), nil, token)
	AssertOK(t, resp)
	AssertDataCount(t, resp, 1)
}

func TestIndex_BusinessObserver_WithoutMeta(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	owner := ta.CreateTestUser(t, 9123456789, "testPassword123", "Owner", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, owner.ID)

	// Create observer without meta (should be forbidden)
	observer := ta.CreateTestUser(t, 9987654321, "testPassword456", "Observer", "User", business.ID, []schema.UserRole{schema.URBusinessObserver})

	token := ta.GenerateTestToken(t, observer)

	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/reservations", business.ID), nil, token)
	AssertForbidden(t, resp)
}

func TestIndex_BusinessObserver_FilterWithTaxonomyAccess(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	owner := ta.CreateTestUser(t, 9123456789, "testPassword123", "Owner", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, owner.ID)
	owner.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(owner)

	// Create taxonomies
	city1 := ta.CreateTestTaxonomy(t, "City 1", schema.TaxonomyTypeCategory, business.ID, nil)
	city2 := ta.CreateTestTaxonomy(t, "City 2", schema.TaxonomyTypeCategory, business.ID, nil)

	// Create test posts and products
	post1 := ta.CreateTestPost(t, "Product 1", schema.PostTypeProduct, business.ID, owner.ID)
	post2 := ta.CreateTestPost(t, "Product 2", schema.PostTypeProduct, business.ID, owner.ID)
	variantType := schema.ProductVariantTypeWashingMachine
	product1 := ta.CreateTestProduct(t, post1.ID, business.ID, 10000, schema.ProductTypeVariant, &variantType)
	product2 := ta.CreateTestProduct(t, post2.ID, business.ID, 15000, schema.ProductTypeVariant, &variantType)

	// Attach taxonomies to posts
	ta.AttachTaxonomyToPost(t, post1.ID, city1.ID)
	ta.AttachTaxonomyToPost(t, post2.ID, city2.ID)

	// Create reservations
	now := time.Now()
	ta.CreateTestReservation(t, owner.ID, product1.ID, business.ID, now.Add(time.Hour), now.Add(2*time.Hour), schema.ReservationStatusReserved)
	ta.CreateTestReservation(t, owner.ID, product2.ID, business.ID, now.Add(3*time.Hour), now.Add(4*time.Hour), schema.ReservationStatusReserved)

	// Create observer with access to city1 only
	observerMeta := &schema.UserMeta{
		TaxonomiesToObserve: schema.UserMetaTaxonomiesToObserve{
			city1.ID: {Checked: true, PartialChecked: false},
		},
	}
	observer := ta.CreateTestUserWithMeta(t, 9987654321, "testPassword456", "Observer", "User", business.ID, []schema.UserRole{schema.URBusinessObserver}, observerMeta)

	token := ta.GenerateTestToken(t, observer)

	// Observer requests with CityID filter for city1 (which they have access to)
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/reservations?CityID=%d", business.ID, city1.ID), nil, token)
	AssertOK(t, resp)
	AssertDataCount(t, resp, 1)
}

func TestIndex_BusinessObserver_FilterWithNoTaxonomyAccess(t *testing.T) {
	ta := SetupTestApp(t)
	defer ta.Cleanup()

	owner := ta.CreateTestUser(t, 9123456789, "testPassword123", "Owner", "User", 0, nil)
	business := ta.CreateTestBusiness(t, "Test Business", schema.BTypeGymManager, owner.ID)
	owner.Permissions[business.ID] = []schema.UserRole{schema.URBusinessOwner}
	ta.DB.Save(owner)

	// Create taxonomies
	city1 := ta.CreateTestTaxonomy(t, "City 1", schema.TaxonomyTypeCategory, business.ID, nil)
	city2 := ta.CreateTestTaxonomy(t, "City 2", schema.TaxonomyTypeCategory, business.ID, nil)

	// Create test posts and products
	post1 := ta.CreateTestPost(t, "Product 1", schema.PostTypeProduct, business.ID, owner.ID)
	post2 := ta.CreateTestPost(t, "Product 2", schema.PostTypeProduct, business.ID, owner.ID)
	variantType := schema.ProductVariantTypeWashingMachine
	product1 := ta.CreateTestProduct(t, post1.ID, business.ID, 10000, schema.ProductTypeVariant, &variantType)
	product2 := ta.CreateTestProduct(t, post2.ID, business.ID, 15000, schema.ProductTypeVariant, &variantType)

	// Attach taxonomies to posts
	ta.AttachTaxonomyToPost(t, post1.ID, city1.ID)
	ta.AttachTaxonomyToPost(t, post2.ID, city2.ID)

	// Create reservations
	now := time.Now()
	ta.CreateTestReservation(t, owner.ID, product1.ID, business.ID, now.Add(time.Hour), now.Add(2*time.Hour), schema.ReservationStatusReserved)
	ta.CreateTestReservation(t, owner.ID, product2.ID, business.ID, now.Add(3*time.Hour), now.Add(4*time.Hour), schema.ReservationStatusReserved)

	// Create observer with access to city1 only
	observerMeta := &schema.UserMeta{
		TaxonomiesToObserve: schema.UserMetaTaxonomiesToObserve{
			city1.ID: {Checked: true, PartialChecked: false},
		},
	}
	observer := ta.CreateTestUserWithMeta(t, 9987654321, "testPassword456", "Observer", "User", business.ID, []schema.UserRole{schema.URBusinessObserver}, observerMeta)

	token := ta.GenerateTestToken(t, observer)

	// Observer requests with CityID filter for city2 (which they DON'T have access to)
	// Should return only city1 results (fallback to observer's accessible taxonomies)
	resp := ta.MakeRequest(t, http.MethodGet, fmt.Sprintf("/v1/business/%d/reservations?CityID=%d", business.ID, city2.ID), nil, token)
	AssertOK(t, resp)
	AssertDataCount(t, resp, 1)
}

package service

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/business/response"
	"strconv"
	"strings"
)

type menuItems map[schema.BusinessType]response.MenuItem

type route struct {
	MenuItems   menuItems
	Permissions []schema.UserRole
}

var routes = []route{
	{
		Permissions: []schema.UserRole{schema.URBusinessOwner, schema.URBusinessObserver},
		MenuItems: menuItems{
			schema.BTypeWMReservation: {Title: "رزرو ها", Href: "/admin/b/:BusinessID/reservations", Icon: "pi pi-clock"},
		},
	},
	{
		Permissions: []schema.UserRole{schema.URBusinessOwner},
		MenuItems: menuItems{
			schema.BTypeWMReservation: {Title: "سفارشات", Href: "/admin/b/:BusinessID/orders", Icon: "pi pi-shopping-bag"},
		},
	},
	{
		Permissions: []schema.UserRole{schema.URBusinessOwner},
		MenuItems: menuItems{
			schema.BTypeWMReservation: {
				Title: "محصولات",
				Icon:  "pi pi-box",
				Child: []response.MenuItem{
					{Title: "همه", Icon: "pi pi-circle-fill", Href: "/admin/b/:BusinessID/products"},
					{Title: "دسته بندی ها", Icon: "pi pi-circle-fill", Href: "/admin/b/:BusinessID/products/categories"},
					//{Title: "تگ ها", Icon: "pi pi-circle-fill", Href: "/admin/b/:BusinessID/products/tags"},
					//{Title: "ویژگی ها", Icon: "pi pi-circle-fill", Href: "/admin/b/:BusinessID/products/attributes"},
				},
			},
		},
	},
	{
		Permissions: []schema.UserRole{schema.URBusinessOwner},
		MenuItems: menuItems{
			schema.BTypeWMReservation: {
				Title: "پست ها",
				Icon:  "pi pi-copy",
				Child: []response.MenuItem{
					{Title: "همه", Icon: "pi pi-circle-fill", Href: "/admin/b/:BusinessID/posts"},
					{Title: "دسته بندی ها", Icon: "pi pi-circle-fill", Href: "/admin/b/:BusinessID/posts/categories"},
					//{Title: "تگ ها", Icon: "pi pi-circle-fill", Href: "/admin/b/:BusinessID/posts/tags"},
					//{Title: "ویژگی ها", Icon: "pi pi-circle-fill", Href: "/admin/b/:BusinessID/products/attributes"},
				},
			},
		},
	},
	{
		Permissions: []schema.UserRole{schema.URBusinessOwner},
		MenuItems: menuItems{
			schema.BTypeWMReservation: {Title: "کوپن ها", Href: "/admin/b/:BusinessID/coupons", Icon: "pi pi-gift"},
		},
	},
	{
		Permissions: []schema.UserRole{schema.URBusinessOwner, schema.URBusinessObserver},
		MenuItems: menuItems{
			schema.BTypeWMReservation: {Title: "تراکنش ها", Href: "/admin/b/:BusinessID/transactions", Icon: "pi pi-arrow-right-arrow-left"},
		},
	},
	{
		Permissions: []schema.UserRole{schema.URBusinessOwner},
		MenuItems: menuItems{
			schema.BTypeWMReservation: {Title: "کاربران", Href: "/admin/b/:BusinessID/users", Icon: "pi pi-users"},
		},
	},
	{
		Permissions: []schema.UserRole{schema.URBusinessOwner},
		MenuItems: menuItems{
			schema.BTypeWMReservation: {Title: "تنظیمات", Href: "/admin/b/:BusinessID/settings", Icon: "pi pi-cog"},
		},
	},
}

// GenerateMenuItems generates menu items based on business ID, type, and user permissions.
func GenerateMenuItems(businessID uint64, businessType schema.BusinessType, user schema.User) []response.MenuItem {
	userPermissions := user.Permissions[businessID]
	var menuItems []response.MenuItem

	for _, route := range routes {
		if hasRoutePermission(userPermissions, route.Permissions) {
			menuItem := route.MenuItems[businessType]
			menuItem.Href = replaceBusinessIDPlaceholder(menuItem.Href, businessID)
			replaceBusinessIDInChildren(menuItem.Child, businessID)
			menuItems = append(menuItems, menuItem)
		}
	}
	return menuItems
}

// replaceBusinessIDPlaceholder replaces the ":BusinessID" placeholder in the URL with the actual business ID.
func replaceBusinessIDPlaceholder(href string, businessID uint64) string {
	return strings.Replace(href, ":BusinessID", strconv.FormatUint(businessID, 10), 1)
}

// replaceBusinessIDInChildren recursively replaces the ":BusinessID" placeholder in child menu items.
func replaceBusinessIDInChildren(items []response.MenuItem, businessID uint64) {
	for i := range items {
		items[i].Href = replaceBusinessIDPlaceholder(items[i].Href, businessID)
		if len(items[i].Child) > 0 {
			replaceBusinessIDInChildren(items[i].Child, businessID)
		}
	}
}

// hasRoutePermission checks if the user has any of the required permissions for the route.
func hasRoutePermission(userPermissions []schema.UserRole, routePermissions []schema.UserRole) bool {
	for _, userPermission := range userPermissions {
		for _, routePermission := range routePermissions {
			if userPermission == routePermission {
				return true
			}
		}
	}
	return false
}

package controller

import (
	"bytes"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils"
	"strconv"

	"go-fiber-starter/app/module/order/request"
	orderResponse "go-fiber-starter/app/module/order/response"
	"go-fiber-starter/app/module/order/service"
	"go-fiber-starter/utils/paginator"
	"go-fiber-starter/utils/response"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
	ptime "github.com/yaa110/go-persian-calendar"
)

type IRestController interface {
	Index(c *fiber.Ctx) error
	Show(c *fiber.Ctx) error
	Store(c *fiber.Ctx) error
	//StoreUniWash(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

func RestController(s service.IService) IRestController {
	return &controller{s}
}

type controller struct {
	service service.IService
}

// Index all Orders
// @Summary      Get all orders
// @Tags         Orders
// @Security     Bearer
// @Param        businessID path int true "Business ID"
// @Param        Status query string false "Order Status" Enums(pending, processing, onHold, completed, cancelled, refunded, failed)
// @Param        StartTime query string false "Reservation Start Time (YYYY-MM-DD)"
// @Param        EndTime query string false "Reservation End Time (YYYY-MM-DD)"
// @Param        OrderStartTime query string false "Order Creation Start Time (YYYY-MM-DD)"
// @Param        OrderEndTime query string false "Order Creation End Time (YYYY-MM-DD)"
// @Param        ProductID query int false "Product ID"
// @Param        CityID query int false "City ID"
// @Param        WorkspaceID query int false "Workspace ID"
// @Param        DormitoryID query int false "Dormitory ID"
// @Param        CouponID query int false "Coupon ID"
// @Param        FullName query string false "User Full Name"
// @Param        HasCoupon query bool false "Filter by coupon presence (true = has coupon, false = no coupon)"
// @Param        Export query string false "Export format (excel)"
// @Router       /business/:businessID/orders [get]
// @Router       /user/business/:businessID/orders [get]
func (_i *controller) Index(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	var req request.Orders
	req.Pagination = paginate
	req.BusinessID = businessID
	req.CouponID, _ = utils.GetUintInQueries(c, "CouponID")
	req.Status = c.Query("Status")
	req.FullName = c.Query("FullName")
	req.ProductID, _ = utils.GetUintInQueries(c, "ProductID")
	export := c.Query("Export")

	// Parse HasCoupon filter (optional boolean)
	if hasCouponStr := c.Query("HasCoupon"); hasCouponStr != "" {
		hasCoupon := c.QueryBool("HasCoupon")
		req.HasCoupon = &hasCoupon
	}

	req.CityID, _ = utils.GetUintInQueries(c, "CityID")
	req.WorkspaceID, _ = utils.GetUintInQueries(c, "WorkspaceID")
	req.DormitoryID, _ = utils.GetUintInQueries(c, "DormitoryID")

	// Build taxonomy filters
	if req.DormitoryID != 0 {
		req.Taxonomies = append(req.Taxonomies, req.DormitoryID)
	} else if req.WorkspaceID != 0 {
		req.Taxonomies = append(req.Taxonomies, req.WorkspaceID)
	} else if req.CityID != 0 {
		req.Taxonomies = append(req.Taxonomies, req.CityID)
	}

	req.StartTime = utils.GetDateInQueries(c, "StartTime")
	req.EndTime = utils.GetDateInQueries(c, "EndTime")
	req.OrderStartTime = utils.GetDateInQueries(c, "OrderStartTime")
	req.OrderEndTime = utils.GetDateInQueries(c, "OrderEndTime")

	orders, totalAmount, paging, err := _i.service.Index(req)
	if err != nil {
		return err
	}

	if export == "excel" {
		return _i.exportToExcel(c, orders, totalAmount)
	}

	return response.Resp(c, response.Response{
		Data: orders,
		Meta: orderResponse.Orders{
			Meta:        paging,
			TotalAmount: totalAmount,
		},
	})
}

// Show one Order
// @Summary      Get one order
// @Tags         Orders
// @Security     Bearer
// @Param        id path int true "Order ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/orders/:id [get]
func (_i *controller) Show(c *fiber.Ctx) error {
	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return err
	}
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	order, err := _i.service.Show(user.ID, id)
	if err != nil {
		return err
	}

	return c.JSON(order)
}

// Store order
// @Summary      Create order
// @Tags         Orders
// @Param 		 order body request.Order true "Order details"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/orders [post]
func (_i *controller) Store(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return err
	}

	req := new(request.Order)
	req.User = user
	req.BusinessID = businessID
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	paymentURL, orderID, err := _i.service.Store(*req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data: map[string]any{"paymentUrl": paymentURL, "orderID": orderID},
	})
}

//// StoreUniWash order
//// @Summary      Create order
//// @Tags         Orders
//// @Param 		 order body request.Order true "Order details"
//// @Param        businessID path int true "Business ID"
//// @Router       /business/:businessID/orders [post]
//func (_i *controller) StoreUniWash(c *fiber.Ctx) error {
//	businessID, err := utils.GetIntInParams(c, "businessID")
//	if err != nil {
//		return err
//	}
//
//	user, err := utils.GetAuthenticatedUser(c)
//	if err != nil {
//		return err
//	}
//
//	req := new(urequest.StoreUniWash)
//	if err := response.ParseAndValidate(c, req); err != nil {
//		return err
//	}
//
//	req.UserID = user.ID
//	req.BusinessID = businessID
//	err = _i.service.StoreUniWash(*req)
//	if err != nil {
//		return err
//	}
//
//	return c.JSON("success")
//}

// Update order
// @Summary      update order
// @Security     Bearer
// @Tags         Orders
// @Param 		 order body request.Order true "Order details"
// @Param        id path int true "Order ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/orders/:id [put]
func (_i *controller) Update(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	req := new(request.Order)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	req.BusinessID = businessID
	err = _i.service.Update(id, *req)
	if err != nil {
		return err
	}

	return c.JSON("success")
}

// Delete order
// @Summary      delete order
// @Tags         Orders
// @Security     Bearer
// @Param        id path int true "Order ID"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/orders/:id [delete]
func (_i *controller) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}

	err = _i.service.Destroy(id)
	if err != nil {
		return err
	}

	return c.JSON("success")
}

func (_i *controller) exportToExcel(c *fiber.Ctx, orders []*orderResponse.Order, totalAmount uint64) error {
	// Create a new Excel file
	f := excelize.NewFile()
	// Create a new sheet
	sheetName := "سفارشات"
	index, _ := f.NewSheet(sheetName)

	// Set RTL view
	f.SetSheetView(sheetName, 0, &excelize.ViewOptions{
		RightToLeft: utils.BoolPtr(true),
	})
	f.SetPanes(sheetName, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		XSplit:      0,
		YSplit:      1,
		TopLeftCell: "A2",
		ActivePane:  "bottomLeft",
	})

	f.SetColWidth(sheetName, "B", "I", 30)
	columnsStyles, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Family: "IRANSans",
			Size:   16,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
		},
	})

	// Headers in Persian
	f.SetCellValue(sheetName, "B1", "مجموع سفارشات تکمیل شده")
	f.SetCellInt(sheetName, "C1", int64(totalAmount))

	f.SetCellValue(sheetName, "A3", "ردیف")
	f.SetCellValue(sheetName, "B3", "نام کامل")
	f.SetCellValue(sheetName, "C3", "مبلغ (تومان)")
	f.SetCellValue(sheetName, "D3", "تاریخ سفارش")
	f.SetCellValue(sheetName, "E3", "وضعیت")
	f.SetCellValue(sheetName, "F3", "شروع رزرو")
	f.SetCellValue(sheetName, "G3", "پایان رزرو")
	f.SetCellValue(sheetName, "H3", "مکان دستگاه")
	f.SetCellValue(sheetName, "I3", "عنوان دستگاه")
	f.SetColStyle(sheetName, "A:P", columnsStyles)

	// Populate data
	for i, order := range orders {
		row := i + 4 // Start from the fourth row
		f.SetCellValue(sheetName, "A"+strconv.Itoa(row), i+1)
		f.SetCellValue(sheetName, "B"+strconv.Itoa(row), order.User.FullName)
		f.SetCellInt(sheetName, "C"+strconv.Itoa(row), int64(order.TotalAmt))
		f.SetCellValue(sheetName, "D"+strconv.Itoa(row), ptime.New(order.CreatedAt).Format("HH:mm - yyyy/MM/dd"))

		f.SetCellValue(sheetName, "E"+strconv.Itoa(row), schema.OrderStatusProxy[order.Status])

		// Add reservation and product info from first order item
		if len(order.OrderItems) > 0 {
			orderItem := order.OrderItems[0]

			// Reservation start and end time
			if orderItem.Reservation != nil {
				f.SetCellValue(sheetName, "F"+strconv.Itoa(row), ptime.New(orderItem.Reservation.StartTime).Format("HH:mm - yyyy/MM/dd"))
				f.SetCellValue(sheetName, "G"+strconv.Itoa(row), ptime.New(orderItem.Reservation.EndTime).Format("HH:mm - yyyy/MM/dd"))
			}

			// Product Title (مکان دستگاه)
			f.SetCellValue(sheetName, "H"+strconv.Itoa(row), orderItem.Meta.ProductTitle)

			// Product SKU (عنوان دستگاه)
			f.SetCellValue(sheetName, "I"+strconv.Itoa(row), orderItem.Meta.ProductSKU)
		}
		// Define background color by status
		var bgColor string
		switch order.Status {
		case schema.OrderStatusCompleted:
			bgColor = "#C6EFCE" // Light green
		case schema.OrderStatusFailed, schema.OrderStatusCancelled:
			bgColor = "#FFC7CE" // Light red
		case schema.OrderStatusPending, schema.OrderStatusProcessing, schema.OrderStatusOnHold:
			bgColor = "#FFEB9C" // Light yellow
		case schema.OrderStatusRefunded:
			bgColor = "#B4C6E7" // Light blue
		default:
			bgColor = "#D9D9D9" // Light gray
		}
		// Create style with background color
		statusStyle, _ := f.NewStyle(&excelize.Style{
			Fill: excelize.Fill{
				Pattern: 1,
				Type:    "pattern",
				Color:   []string{bgColor},
			},
			Alignment: &excelize.Alignment{
				Horizontal: "right",
				Vertical:   "center",
			},
			Font: &excelize.Font{
				Family: "IRANSans",
				Size:   16,
			},
		})
		// Apply the style to the status cell
		cell := "E" + strconv.Itoa(row)
		f.SetCellStyle(sheetName, cell, cell, statusStyle)
	}

	// Set active sheet
	f.SetActiveSheet(index)

	// Write the file to a buffer
	buf := new(bytes.Buffer)
	if err := f.Write(buf); err != nil {
		return err
	}

	// Set the content type and filename
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=orders.xlsx")

	// Send the buffer as the response
	return c.Send(buf.Bytes())
}

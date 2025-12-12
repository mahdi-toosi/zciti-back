package controller

import (
	"bytes"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/transaction/request"
	transactionResponse "go-fiber-starter/app/module/transaction/response"
	"go-fiber-starter/app/module/transaction/service"
	walletService "go-fiber-starter/app/module/wallet/service"
	"go-fiber-starter/utils"
	"go-fiber-starter/utils/paginator"
	"go-fiber-starter/utils/response"
	"slices"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
	ptime "github.com/yaa110/go-persian-calendar"
)

type IRestController interface {
	Index(c *fiber.Ctx) error
	//Show(c *fiber.Ctx) error
	//Store(c *fiber.Ctx) error
	//Update(c *fiber.Ctx) error
	//Delete(c *fiber.Ctx) error
}

func RestController(s service.IService, w walletService.IService) IRestController {
	return &controller{s, w}
}

type controller struct {
	service       service.IService
	walletService walletService.IService
}

// Index all Transactions
// @Summary      Get all transactions
// @Tags         Transactions
// @Security     Bearer
// @Param        WalletID path int true "Wallet ID"
// @Router       /wallets/:walletID/transactions [get]
func (_i *controller) Index(c *fiber.Ctx) error {
	walletID, err := utils.GetIntInParams(c, "walletID")
	if err != nil {
		return err
	}
	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	var req request.Transactions
	req.WalletID = walletID
	req.Pagination = paginate
	export := c.Query("Export")
	req.Status = c.Query("Status")
	req.ProductID, _ = utils.GetUintInQueries(c, "ProductID")

	req.CityID, _ = utils.GetUintInQueries(c, "CityID")
	req.WorkspaceID, _ = utils.GetUintInQueries(c, "WorkspaceID")
	req.DormitoryID, _ = utils.GetUintInQueries(c, "DormitoryID")
	var taxonomies []uint64
	if req.DormitoryID != 0 {
		taxonomies = append(taxonomies, req.DormitoryID)
	} else if req.WorkspaceID != 0 {
		taxonomies = append(taxonomies, req.WorkspaceID)
	} else if req.CityID != 0 {
		taxonomies = append(taxonomies, req.CityID)
	}

	req.StartTime = utils.GetDateInQueries(c, "StartTime")
	req.EndTime = utils.GetDateInQueries(c, "EndTime")
	if req.EndTime != nil && !req.EndTime.IsZero() {
		// end of the end date
		v := utils.EndOfDay(*req.EndTime)
		req.EndTime = &v
	}

	// get Wallet and check the owner
	wallet, err := _i.walletService.Show(&walletID, nil, nil)
	if err != nil {
		return err
	}

	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return err
	}

	if wallet.UserID != nil {
		if *wallet.UserID != user.ID {
			return fiber.ErrForbidden
		}
	}

	if wallet.BusinessID != nil {
		if !user.IsBusinessOwner(*wallet.BusinessID) && !user.IsBusinessObserver(*wallet.BusinessID) {
			return fiber.ErrForbidden
		}

		if user.IsBusinessObserver(*wallet.BusinessID) {
			if user.Meta == nil {
				return &fiber.Error{
					Code:    fiber.StatusForbidden,
					Message: "برای شما سطح دسترسی مشخص نشده است.",
				}
			}
			observerTaxonomies := user.Meta.GetTaxonomiesToObserve(true, false)

			if len(taxonomies) > 0 {
				for _, taxonomy := range taxonomies {
					if slices.Contains(observerTaxonomies, taxonomy) {
						req.Taxonomies = append(req.Taxonomies, taxonomy)
					}
				}
			}

			if len(req.Taxonomies) == 0 {
				req.Taxonomies = observerTaxonomies
			}
		}
	}

	transactions, totalAmount, paging, err := _i.service.Index(req)
	if err != nil {
		return err
	}

	//return response.Resp(c, response.Response{
	//	Data: transactions,
	//	Meta: paging,
	//})

	if export == "excel" {
		// Create a new Excel file
		f := excelize.NewFile()
		// Create a new sheet
		sheetName := "تراکنش ها"
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

		f.SetColWidth(sheetName, "B", "D", 30)
		columnsStyles, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{
				Family: "IRANSans", // Font name as installed in the OS
				Size:   16,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "right",
			},
		})

		// Headers in Persian
		f.SetCellValue(sheetName, "B1", "مجموع تراکنش های موفق")
		f.SetCellInt(sheetName, "C1", int64(totalAmount))

		f.SetCellValue(sheetName, "A3", "ردیف")
		f.SetCellValue(sheetName, "B3", "نام کامل")
		f.SetCellValue(sheetName, "C3", "مبلغ (تومان)")
		f.SetCellValue(sheetName, "D3", "تاریخ ")
		f.SetCellValue(sheetName, "E3", "وضعیت ")
		f.SetColStyle(sheetName, "A:P", columnsStyles)

		// Populate data
		for i, transaction := range transactions {
			row := i + 4 // Start from the second row
			f.SetCellValue(sheetName, "A"+strconv.Itoa(row), i+1)
			f.SetCellValue(sheetName, "B"+strconv.Itoa(row), transaction.User.FullName)
			f.SetCellInt(sheetName, "C"+strconv.Itoa(row), int64(transaction.Amount))
			f.SetCellValue(sheetName, "D"+strconv.Itoa(row), ptime.New(transaction.UpdatedAt).Format("HH:mm - yyyy/MM/dd"))

			f.SetCellValue(sheetName, "E"+strconv.Itoa(row), schema.TransactionStatusProxy[transaction.Status])
			// Define background color by status
			var bgColor string
			switch transaction.Status {
			case "success":
				bgColor = "#C6EFCE" // Light green
			case "failed":
				bgColor = "#FFC7CE" // Light red
			case "pending":
				bgColor = "#FFEB9C" // Light yellow
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
					Family: "IRANSans", // Font name as installed in the OS
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
		c.Set("Content-Disposition", "attachment; filename=users.xlsx")

		// Send the buffer as the response
		return c.Send(buf.Bytes())
	} else {
		return response.Resp(c, response.Response{
			Data: transactions,
			Meta: transactionResponse.Transactions{
				Meta:        paging,
				TotalAmount: totalAmount,
			},
		})
	}
}

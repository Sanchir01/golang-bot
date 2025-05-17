package excel

import (
	"fmt"
	"github.com/Sanchir01/kafka-tg/internal/entity"

	"github.com/xuri/excelize/v2"
)

func CreateFileExcel(products []entity.ProductWithQuantity, filePath, sheetName string) error {
	f := excelize.NewFile()

	f.SetSheetName(f.GetSheetName(0), sheetName)
	headers := []string{"Название", "Цена", "Количество", "Цена за эти товары", "Цена за все товары"}
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string('A'+i))
		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			return fmt.Errorf("failed to set header: %v", err)
		}
	}
	maxLength := make([]int, len(headers))
	var totalOrderPrice int
	for rowIndex, product := range products {
		row := rowIndex + 2
		err := f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), product.Candles.Slug)
		if err != nil {
			return fmt.Errorf("failed to set title: %v", err)
		}
		maxLength[0] = max(maxLength[0], len(product.Candles.Slug))
		err = f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), product.Candles.Price)
		if err != nil {
			return fmt.Errorf("failed to set price: %v", err)
		}
		maxLength[1] = max(maxLength[1], len(fmt.Sprintf("%d", product.Candles.Price)))

		err = f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), product.Quantity)
		if err != nil {
			return fmt.Errorf("failed to set quantity: %v", err)
		}
		maxLength[2] = max(maxLength[2], len(fmt.Sprintf("%d", product.Quantity)))

		totalPrice := product.Quantity * product.Candles.Price
		err = f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), totalPrice)
		if err != nil {
			return fmt.Errorf("failed to set product quantity price: %v", err)
		}
		maxLength[3] = max(maxLength[3], len(fmt.Sprintf("%d", totalPrice)))
		totalOrderPrice += totalPrice
	}
	totalRow := len(products) + 2
	err := f.SetCellValue(sheetName, fmt.Sprintf("E%d", totalRow), totalOrderPrice)
	if err != nil {
		return fmt.Errorf("failed to set total order price: %v", err)
	}
	for i, maxLength := range maxLength {
		col := string('A' + i)
		width := float64(maxLength + 10)
		if err := f.SetColWidth(sheetName, col, col, width); err != nil {
			return fmt.Errorf("failed to set column width for %s: %v", col, err)
		}
	}
	if err := f.SaveAs(filePath); err != nil {
		return fmt.Errorf("failed to save excel file: %v", err)
	}

	return nil
}

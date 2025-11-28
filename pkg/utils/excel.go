package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/web-dashboard-made-by-renz/backend/internal/models"
	"github.com/xuri/excelize/v2"
)

func ParseExcelToColoris(file *excelize.File) ([]models.Coloris, error) {
	sheets := file.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("no sheets found in Excel file")
	}

	rows, err := file.GetRows(sheets[0])
	if err != nil {
		return nil, err
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("Excel file must have at least a header row and one data row")
	}

	var colorisData []models.Coloris
	for i, row := range rows {
		if i == 0 {
			continue
		}

		if len(row) < 11 {
			continue
		}

		timestamp, err := parseTimestamp(row[0])
		if err != nil {
			timestamp = time.Now()
		}

		// Debug logging
		fmt.Printf("Row %d - Nilai PG Raw: '%s'\n", i, row[8])
		nilaiPG := parseNilaiPG(row[8])
		fmt.Printf("Row %d - Nilai PG Parsed: %.2f\n", i, nilaiPG)

		nilaiAkhir, _ := strconv.ParseFloat(row[9], 64)
		total, _ := strconv.ParseFloat(row[10], 64)

		coloris := models.Coloris{
			Timestamp:            timestamp,
			Bulan:                row[1],
			Region:               row[2],
			Cabang:               row[3],
			Materi:               row[4],
			NamaAtasanLangsung:   row[5],
			NamaToko:             row[6],
			NamaLengkapSesuaiKTP: row[7],
			NilaiPG:              nilaiPG,
			NilaiAkhir:           nilaiAkhir,
			Total:                total,
		}

		colorisData = append(colorisData, coloris)
	}

	return colorisData, nil
}

func ExportColorisToExcel(colorisData []models.Coloris) (*excelize.File, error) {
	f := excelize.NewFile()

	sheetName := "Coloris Data"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

	headers := []string{
		"Timestamp",
		"Bulan",
		"Region",
		"Cabang",
		"Materi",
		"Nama Atasan Langsung",
		"Nama Toko",
		"Nama Lengkap Sesuai KTP",
		"Nilai PG",
		"Nilai Akhir",
		"Total",
	}

	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#E0E0E0"}, Pattern: 1},
	})
	if err == nil {
		f.SetCellStyle(sheetName, "A1", fmt.Sprintf("%c1", 'A'+len(headers)-1), style)
	}

	for i, data := range colorisData {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), data.Timestamp.Format("1/2/2006 15:04:05"))
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), data.Bulan)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), data.Region)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), data.Cabang)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), data.Materi)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), data.NamaAtasanLangsung)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), data.NamaToko)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), data.NamaLengkapSesuaiKTP)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), data.NilaiPG)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), data.NilaiAkhir)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), data.Total)
	}

	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	return f, nil
}

func parseNilaiPG(nilaiStr string) float64 {
	// Trim whitespace
	nilaiStr = trimSpace(nilaiStr)

	// Handle empty string
	if nilaiStr == "" {
		return 0
	}

	// Handle format "80/100" or "80 / 100" or "90/100"
	// Check if contains "/"
	slashIndex := -1
	for i := 0; i < len(nilaiStr); i++ {
		if nilaiStr[i] == '/' {
			slashIndex = i
			break
		}
	}

	if slashIndex > 0 {
		// Extract the first number before "/"
		firstPart := nilaiStr[:slashIndex]
		firstPart = trimSpace(firstPart)

		// Try to parse as float
		if val, err := strconv.ParseFloat(firstPart, 64); err == nil {
			return val
		}
	}

	// If no "/" found or parsing failed, try to parse as regular number
	if val, err := strconv.ParseFloat(nilaiStr, 64); err == nil {
		return val
	}

	return 0
}

func trimSpace(s string) string {
	start := 0
	end := len(s)

	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}

	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}

	return s[start:end]
}

func parseFloatFromExcel(numStr string) float64 {
	// Trim whitespace
	numStr = trimSpace(numStr)

	if numStr == "" {
		return 0
	}

	// Remove common formatting:
	// - Remove "Rp" or "Rp."
	// - Remove thousand separators (dot, comma, space)
	// - Replace comma decimal separator with dot

	cleaned := ""
	hasDecimal := false
	decimalChar := byte('.')

	// First pass: detect decimal separator
	// If we find comma after many digits, it's likely decimal separator
	// If we find dot after many digits, it's likely decimal separator
	lastCommaPos := -1
	lastDotPos := -1

	for i := 0; i < len(numStr); i++ {
		if numStr[i] == ',' {
			lastCommaPos = i
		} else if numStr[i] == '.' {
			lastDotPos = i
		}
	}

	// If comma is in last 3 positions, it's decimal separator
	// Same with dot
	if lastCommaPos > 0 && lastCommaPos >= len(numStr)-3 {
		decimalChar = byte(',')
	}

	// Second pass: clean the string
	for i := 0; i < len(numStr); i++ {
		char := numStr[i]

		// Skip non-numeric characters except minus, dot, and comma
		if char >= '0' && char <= '9' {
			cleaned += string(char)
		} else if char == '-' && i == 0 {
			// Allow minus at start
			cleaned += string(char)
		} else if (char == '.' || char == ',') && !hasDecimal {
			// Check if this is decimal separator
			if (char == decimalChar && i >= len(numStr)-3) ||
				(char == '.' && decimalChar == '.' && lastDotPos == i) ||
				(char == ',' && decimalChar == ',' && lastCommaPos == i) {
				cleaned += "."
				hasDecimal = true
			}
			// Otherwise skip (it's thousand separator)
		}
		// Skip other characters (Rp, space, etc)
	}

	if cleaned == "" || cleaned == "-" {
		return 0
	}

	result, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		fmt.Printf("Warning: Failed to parse '%s' (cleaned: '%s'): %v\n", numStr, cleaned, err)
		return 0
	}

	return result
}

func parseTimestamp(dateStr string) (time.Time, error) {
	formats := []string{
		"1/2/2006 15:04:05",
		"1/2/2006 3:04:05 PM",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02",
		"01/02/2006",
		"2/1/2006 15:04:05",
		"02/01/2006 15:04:05",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

// ==================== TRAINING EXCEL UTILS ====================

func ParseExcelToTraining(file *excelize.File) ([]models.Training, error) {
	sheets := file.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("no sheets found in Excel file")
	}

	rows, err := file.GetRows(sheets[0])
	if err != nil {
		return nil, err
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("Excel file must have at least a header row and one data row")
	}

	var trainingData []models.Training
	for i, row := range rows {
		if i == 0 {
			continue
		}

		if len(row) < 11 {
			continue
		}

		timestamp, err := parseTimestamp(row[0])
		if err != nil {
			timestamp = time.Now()
		}

		// Debug logging
		fmt.Printf("Row %d - Total Nilai Raw: '%s'\n", i, row[8])
		totalNilai := parseNilaiPG(row[8]) // Reuse the same parsing function
		fmt.Printf("Row %d - Total Nilai Parsed: %.2f\n", i, totalNilai)

		nilaiEssay, _ := strconv.ParseFloat(row[9], 64)
		total, _ := strconv.ParseFloat(row[10], 64)

		training := models.Training{
			Timestamp:            timestamp,
			Bulan:                row[1],
			Region:               row[2],
			CabangArea:           row[3],
			NamaAtasanLangsung:   row[4],
			MateriPelatihan:      row[5],
			NamaLengkapSesuaiKTP: row[6],
			Jabatan:              row[7],
			TotalNilai:           totalNilai,
			NilaiEssay:           nilaiEssay,
			Total:                total,
		}

		trainingData = append(trainingData, training)
	}

	return trainingData, nil
}

func ExportTrainingToExcel(trainingData []models.Training) (*excelize.File, error) {
	f := excelize.NewFile()

	sheetName := "Training Data"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

	headers := []string{
		"Timestamp",
		"Bulan",
		"Region",
		"Cabang/Area",
		"Nama Atasan Langsung",
		"Materi Pelatihan",
		"Nama Lengkap Sesuai KTP",
		"Jabatan",
		"Total Nilai",
		"Nilai Essay",
		"Total",
	}

	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#E0E0E0"}, Pattern: 1},
	})
	if err == nil {
		f.SetCellStyle(sheetName, "A1", fmt.Sprintf("%c1", 'A'+len(headers)-1), style)
	}

	for i, data := range trainingData {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), data.Timestamp.Format("01/02/2006 15:04:05"))
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), data.Bulan)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), data.Region)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), data.CabangArea)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), data.NamaAtasanLangsung)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), data.MateriPelatihan)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), data.NamaLengkapSesuaiKTP)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), data.Jabatan)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), data.TotalNilai)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), data.NilaiEssay)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), data.Total)
	}

	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	return f, nil
}

// ==================== SELLOUT EXCEL UTILS ====================

func ParseExcelToSellout(file *excelize.File) ([]models.Sellout, error) {
	sheets := file.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("no sheets found in Excel file")
	}

	rows, err := file.GetRows(sheets[0])
	if err != nil {
		return nil, err
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("Excel file must have at least a header row and one data row")
	}

	var selloutData []models.Sellout
	for i, row := range rows {
		if i == 0 {
			continue
		}

		// minimal 18 kolom sesuai header ekspor terbaru
		if len(row) < 18 {
			continue
		}

		tahun, _ := strconv.Atoi(row[0])
		bulan, _ := strconv.Atoi(row[1])

		targetSellout := parseFloatFromExcel(row[13])
		selloutTT := parseFloatFromExcel(row[14])
		selloutRM := parseFloatFromExcel(row[15])
		primafix := parseFloatFromExcel(row[16])
		totalSellout := parseFloatFromExcel(row[17])
		masaKerja := parseFloatFromExcel(row[10])

		sellout := models.Sellout{
			Tahun:            tahun,
			Bulan:            bulan,
			Reg:              row[2],
			Cabang:           row[3],
			Outlet:           row[4],
			AreaCover:        row[5],
			MosSs:            row[6],
			NamaColorist:     row[7],
			NoReg:            row[8],
			TanggalBergabung: row[9],
			MasaKerja:        masaKerja,
			CHL:              row[11],
			Wilayah:          row[12],
			TargetSellout:    targetSellout,
			SelloutTT:        selloutTT,
			SelloutRM:        selloutRM,
			Primafix:         primafix,
			TotalSellout:     totalSellout,
		}

		selloutData = append(selloutData, sellout)
	}

	return selloutData, nil
}

func ExportSelloutToExcel(selloutData []models.Sellout) (*excelize.File, error) {
	f := excelize.NewFile()

	sheetName := "Sellout Data"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

	headers := []string{
		"Tahun",
		"Bulan",
		"Reg",
		"Cabang",
		"Outlet",
		"Area Cover",
		"MOS/SS",
		"Nama Colorist",
		"No Reg",
		"Tanggal Bergabung",
		"Masa Kerja",
		"CHL",
		"Wilayah",
		"Target Sellout",
		"Sellout TT",
		"Sellout RM",
		"Primafix",
		"Total Sellout",
	}

	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"#E0E0E0"}, Pattern: 1},
	})
	if err == nil {
		f.SetCellStyle(sheetName, "A1", fmt.Sprintf("%c1", 'A'+len(headers)-1), style)
	}

	for i, data := range selloutData {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), data.Tahun)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), data.Bulan)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), data.Reg)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), data.Cabang)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), data.Outlet)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), data.AreaCover)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), data.MosSs)
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), data.NamaColorist)
		f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), data.NoReg)
		f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), data.TanggalBergabung)
		f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), data.MasaKerja)
		f.SetCellValue(sheetName, fmt.Sprintf("L%d", row), data.CHL)
		f.SetCellValue(sheetName, fmt.Sprintf("M%d", row), data.Wilayah)
		f.SetCellValue(sheetName, fmt.Sprintf("N%d", row), data.TargetSellout)
		f.SetCellValue(sheetName, fmt.Sprintf("O%d", row), data.SelloutTT)
		f.SetCellValue(sheetName, fmt.Sprintf("P%d", row), data.SelloutRM)
		f.SetCellValue(sheetName, fmt.Sprintf("Q%d", row), data.Primafix)
		f.SetCellValue(sheetName, fmt.Sprintf("R%d", row), data.TotalSellout)
	}

	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	return f, nil
}

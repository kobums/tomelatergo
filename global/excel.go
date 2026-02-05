package global

import (
	"fmt"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"tomelater/global/log"

	"github.com/xuri/excelize/v2"

	humanize "github.com/dustin/go-humanize"
	uuid "github.com/google/uuid"
)

type Excel struct {
	File     *excelize.File
	Width    []int
	Align    []string
	Cols     int
	Rows     int
	FontSize float64

	Pos      int
	Height   float64
	Filename string
	Path     string

	Sheet string

	Index int
}

func (p *Excel) GetCell(col string, row int) string {
	pos := fmt.Sprintf("%v%v", col, row)

	str, err := p.File.GetCellValue(p.Sheet, pos)

	if err != nil {
		return ""
	}

	return str
}

func NewExcelReader(filename string) *Excel {
	var item Excel

	f, err := excelize.OpenFile(filename)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil
	}

	item.File = f

	return &item
}

func (p *Excel) SetSheet(str string) {
	p.Sheet = str
}

func OpenExcel(filename string, title string, fontSize float64, header []string, width []int, align []string) *Excel {
	var item Excel

	item.Width = width
	item.Align = align
	item.Cols = len(header)
	item.Pos = 0
	item.Rows = 0
	item.FontSize = fontSize

	f, err := excelize.OpenFile(filename)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil
	}

	item.File = f

	_, err = item.File.NewSheet("Sheet1")
	if err != nil {
		log.Error().Msg(err.Error())
		return nil
	}

	for i, value := range header {
		t := fmt.Sprintf("%v", rune('A'+i))
		err := item.File.SetColWidth("Sheet1", t, t, float64(width[i])*0.8)
		if err != nil {
			log.Error().Msg(err.Error())
			return nil
		}
		item.HeaderCell(value)
	}

	return &item
}

func New() *Excel {
	var item Excel

	item.FontSize = 10

	item.File = excelize.NewFile()

	return &item
}

func NewExcel(title string, sheet string, fontSize float64, header []string, width []int, align []string) *Excel {
	var item Excel

	item.Width = width
	item.Align = align
	item.Cols = len(header)
	item.Pos = 0
	item.Rows = 0
	item.FontSize = fontSize

	item.File = excelize.NewFile()

	if sheet == "" {
		sheet = "Sheet1"
	}
	item.Sheet = sheet

	item.Index, _ = item.File.NewSheet(sheet)

	for i, value := range header {
		t := ""

		if i > 25 {
			t = fmt.Sprintf("A%c", rune('A'+(i-26)))
		} else {
			t = fmt.Sprintf("%c", rune('A'+i))
		}
		err := item.File.SetColWidth(item.Sheet, t, t, float64(width[i])*0.8)
		if err != nil {
			log.Error().Msg(err.Error())
			return nil
		}
		item.HeaderCell(value)
	}

	if item.Sheet != "Sheet1" {
		err := item.File.DeleteSheet("Sheet1")
		if err != nil {
			log.Error().Msg(err.Error())
			return nil
		}
	}

	return &item
}

func (p *Excel) NewSheet(sheet string, header []string, width []int, align []string) {
	p.Width = width
	p.Align = align
	p.Cols = len(header)
	p.Pos = 0
	p.Rows = 0

	p.Sheet = sheet

	p.Index, _ = p.File.NewSheet(sheet)

	for i, value := range header {
		t := ""

		if i > 25 {
			t = fmt.Sprintf("A%c", rune('A'+(i-26)))
		} else {
			t = fmt.Sprintf("%c", rune('A'+i))
		}

		err := p.File.SetColWidth(p.Sheet, t, t, float64(width[i])*0.8)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		p.HeaderCell(value)
	}

	if sheet != "Sheet1" {
		err := p.File.DeleteSheet("Sheet1")
		if err != nil {
			log.Error().Msg(err.Error())
		}
	}
}

func (p *Excel) SetHeight(height float64) {
	p.Height = height
}

func (p *Excel) Save(filename string) string {
	if filename == "" {
		p.Filename = GetTempFilename() + ".xlsx"
	} else {
		p.Path = fmt.Sprintf("webdata/temp/%v", uuid.New().String())
		err := os.Mkdir(p.Path, 0755)
		if err != nil {
			return ""
		}
		p.Filename = fmt.Sprintf("%v/%v", p.Path, filename)
	}

	p.File.SetActiveSheet(p.Index)

	err := p.File.SaveAs(p.Filename)
	if err != nil {
		log.Error().Msg(err.Error())
	}

	return p.Filename
}

func (p *Excel) Remove() {
	if p.Path != "" {
		os.RemoveAll(p.Path)
	} else {
		os.Remove(p.Filename)
	}
}

func (p *Excel) HeaderCell(str string) {
	if p.Pos == 0 {
		err := p.File.SetRowHeight(p.Sheet, p.Rows+1, 30)
		if err != nil {
			log.Error().Msg(err.Error())
		}
	}

	style, _ := p.File.NewStyle(&excelize.Style{
		Border: []excelize.Border{{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Font:      &excelize.Font{Color: "000000", Size: p.FontSize},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#CCCCCC"}},
	})

	//`{"alignment":{"horizontal":"center","vertical":"center"},"border":[{"type":"left","color":"000000","style":1},{"type":"top","color":"000000","style":1},{"type":"bottom","color":"000000","style":1},{"type":"right",   "color":"000000","style":1}],"fill":{"type":"pattern","pattern":1,"color":["#CCCCCC"]},"number_format":0,"lang":"ko-kr"}`)

	t := p.GetColName()
	err := p.File.SetCellValue(p.Sheet, t, str)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
	err = p.File.SetCellStyle(p.Sheet, t, t, style)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}

	p.Pos++

	if p.Pos == p.Cols {
		p.Pos = 0
		p.Rows++
	}
}

func (p *Excel) SetHeaderStyle(col string, row int, fontSize float64) {
	t := fmt.Sprintf("%v%v", col, row)

	style, _ := p.File.NewStyle(&excelize.Style{
		Border: []excelize.Border{{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Font:      &excelize.Font{Color: "000000", Size: fontSize},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#CCCCCC"}},
	})

	err := p.File.SetCellStyle(p.Sheet, t, t, style)
	if err != nil {
		log.Error().Msg(err.Error())
	}
}

func (p *Excel) GetColName() string {
	t := ""

	if p.Pos > 25 {
		t = fmt.Sprintf("A%c%v", rune('A'+(p.Pos-26)), p.Rows+1)
	} else {
		t = fmt.Sprintf("%c%v", rune('A'+p.Pos), p.Rows+1)
	}

	return t
}

func (p *Excel) Cell(str string) string {
	if p.Pos == 0 && p.Rows > 0 {
		err := p.File.SetRowHeight(p.Sheet, p.Rows+1, p.Height)
		if err != nil {
			log.Error().Msg(err.Error())
		}
	}

	align := "center"
	if p.Align[p.Pos] == "L" {
		align = "left"
	} else if p.Align[p.Pos] == "R" {
		align = "right"
	}

	//style, _ := p.File.NewStyle(`{"alignment":{"horizontal":"` + align + `","vertical":"center"},"border":[{"type":"left","color":"000000","style":1},{"type":"top","color":"000000","style":1},{"type":"bottom","color":"000000","style":1},{"type":"right",   "color":"000000","style":1}],"fill":{"type":"pattern","pattern":1,"color":["#FFFFFF"]},"number_format":0,"lang":"ko-kr"}`)

	style, _ := p.File.NewStyle(&excelize.Style{
		Border: []excelize.Border{{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
		Font:      &excelize.Font{Color: "000000", Size: p.FontSize},
		Alignment: &excelize.Alignment{Horizontal: align, Vertical: "center", WrapText: true},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#FFFFFF"}},
	})

	t := p.GetColName()
	err := p.File.SetCellValue(p.Sheet, t, str)
	if err != nil {
		return ""
	}
	err = p.File.SetCellStyle(p.Sheet, t, t, style)
	if err != nil {
		return ""
	}

	p.Pos++

	if p.Pos == p.Cols {
		p.Pos = 0
		p.Rows++
	}

	return t
}

func (p *Excel) CellInt(value int) {
	p.Cell(fmt.Sprintf("%v", value))
}

func (p *Excel) CellInt64(value int64) {
	p.Cell(fmt.Sprintf("%v", value))
}

func (p *Excel) CellPrice(value int) {
	p.Cell(fmt.Sprintf("₩ %v", humanize.Comma(int64(value))))
}

func (p *Excel) CellImage(filename string) {
	/*
		t := p.Cell("")

		width := 100.0
		height := 100.0

		xScale := 0.2
		yScale := 0.2

		xOffset := 10
		yOffset := 10

		log.Println(xScale, yScale)

		if reader, err := os.Open(filename); err == nil {
			defer reader.Close()
			im, _, err := image.DecodeConfig(reader)
			if err != nil {
			}
			fmt.Printf("%d %d\n", im.Width, im.Height)

			xScale = width / float64(im.Width)
			yScale = height / float64(im.Height)
		} else {
			fmt.Println("Impossible to open the file:", err)
		}

			if err := p.File.AddPicture(p.Sheet, t, filename, fmt.Sprintf(`{"x_scale": %v, "y_scale": %v, "x_offset":%v, "y_offset":%v}`, xScale, yScale, xOffset, yOffset)); err != nil {
				fmt.Println(err)
			}
	*/
}

func (p *Excel) MergeCell(srcCol string, srcRow int, descCol string, descRow int) {
	src := fmt.Sprintf("%v%v", srcCol, srcRow)
	desc := fmt.Sprintf("%v%v", descCol, descRow)
	err := p.File.MergeCell(p.Sheet, src, desc)
	if err != nil {
		log.Error().Msg(err.Error())
	}
}

func (p *Excel) InsertRow(row int, n int) {
	err := p.File.InsertRows(p.Sheet, row, n)
	if err != nil {
		log.Error().Msg(err.Error())
	}
}

func (p *Excel) SetCellValue(col string, row int, value interface{}) {
	cell := fmt.Sprintf("%v%v", col, row)

	err := p.File.SetCellValue(p.Sheet, cell, value)
	if err != nil {
		log.Error().Msg(err.Error())
	}
}

func (p *Excel) SetRowHeight(row int, height float64) {
	err := p.File.SetRowHeight(p.Sheet, row, height)
	if err != nil {
		log.Error().Msg(err.Error())
	}
}

func (p *Excel) GetRows(sheet string) [][]string {
	values, _ := p.File.GetRows(sheet, excelize.Options{})
	return values
}

func (p *Excel) Close() {
	p.File.Close()
}

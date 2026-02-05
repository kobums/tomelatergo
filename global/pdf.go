package global

import (
	"fmt"
	"os"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type Pdf struct {
	Context *gofpdf.Fpdf
	Width   []int
	Align   []string
	Cols    int

	Pos      int
	Filename string
}

func NewPdf(title string, header []string, width []int, align []string, headerFont float64, bodyFont float64) *Pdf {
	t := time.Now()
	datetime := fmt.Sprintf("%02d-%02d-%04d %02d:%02d:%02d", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second())

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Times", "B", 20)
	pdf.Cell(40, 10, title)
	pdf.Ln(12)
	pdf.SetFont("Times", "", 10)
	pdf.Cell(40, 10, datetime)
	pdf.Ln(20)

	// header
	pdf.SetFont("Times", "B", headerFont)
	pdf.SetFillColor(240, 240, 240)
	for i, str := range header {
		pdf.CellFormat(float64(width[i]), 7, str, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Times", "", bodyFont)
	pdf.SetFillColor(255, 255, 255)

	var item Pdf
	item.Context = pdf

	item.Width = width
	item.Align = align
	item.Cols = len(header)
	item.Pos = 0

	return &item
}

func (p *Pdf) Save(filename string) string {
	if filename == "" {
		p.Filename = GetTempFilename()
	} else {
		p.Filename = filename
	}
	p.Context.OutputFileAndClose(p.Filename)

	return p.Filename
}

func (p *Pdf) Remove() {
	os.Remove(p.Filename)
}

func (p *Pdf) Cell(str string) {
	p.Context.CellFormat(float64(p.Width[p.Pos]), 7, str, "1", 0, p.Align[p.Pos], false, 0, "")

	p.Pos++

	if p.Pos == p.Cols {
		p.Pos = 0
		p.Context.Ln(-1)
	}
}

func (p *Pdf) CellInt(value int) {
	str := fmt.Sprintf("%v", value)

	p.Cell(str)
}

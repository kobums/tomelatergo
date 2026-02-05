package pdf

import (
	"fmt"
	"tomelater/global/log"

	"github.com/signintech/gopdf"
)

type FontWeight int

const (
	_ FontWeight = iota

	Bold
	Medium
	Regular
	Thin
)

const Left = 8        //001000
const Top = 4         //000100
const Right = 2       //000010
const Bottom = 1      //000001
const Center = 16     //010000
const Middle = 32     //100000
const AllBorders = 15 //001111

type Pdf struct {
	Context *gopdf.GoPdf

	X float64
	Y float64

	MarginX float64
	MarginY float64

	Font       string
	FontSize   float64
	FontWeight FontWeight

	TextColor Color
}

type Color struct {
	R uint8
	G uint8
	B uint8
}

func FromRGB(r int, g int, b int) Color {
	return Color{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
	}
}

func NewPdf() *Pdf {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	err := pdf.AddTTFFont("NotoSans-Bold", "./fonts/NotoSansKR-Bold.ttf")
	if err != nil {
		log.Error().Msg(err.Error())
		return nil
	}

	err = pdf.AddTTFFont("NotoSans-Medium", "./fonts/NotoSansKR-Medium.ttf")
	if err != nil {
		log.Error().Msg(err.Error())
		return nil
	}

	err = pdf.AddTTFFont("NotoSans-Regular", "./fonts/NotoSansKR-Regular.ttf")
	if err != nil {
		log.Error().Msg(err.Error())
		return nil
	}

	err = pdf.AddTTFFont("NotoSans-Thin", "./fonts/NotoSansKR-Thin.ttf")
	if err != nil {
		log.Error().Msg(err.Error())
		return nil
	}

	ret := &Pdf{
		Context:    &pdf,
		Font:       "NotoSans",
		FontSize:   19,
		FontWeight: Regular,
	}

	ret.SetFont(ret.Font)
	ret.SetFontWeight(ret.FontWeight)
	return ret
}

func (c *Pdf) GetFont() string {
	fontWeight := "Regular"

	if c.FontWeight == Bold {
		fontWeight = "Bold"
	} else if c.FontWeight == Medium {
		fontWeight = "Medium"
	} else if c.FontWeight == Thin {
		fontWeight = "Thin"
	}

	ret := fmt.Sprintf("%v-%v", c.Font, fontWeight)
	log.Println(ret)

	return ret
}

func (c *Pdf) SetFont(font string) {
	c.Font = font
	err := c.Context.SetFont(c.GetFont(), "", c.FontSize)
	if err != nil {
		log.Error().Msg(err.Error())
	}
}

func (c *Pdf) SetFontWeight(weight FontWeight) {
	c.FontWeight = weight
	err := c.Context.SetFont(c.GetFont(), "", c.FontSize)
	if err != nil {
		log.Error().Msg(err.Error())
	}
}

func (c *Pdf) SetFontSize(size float64) {
	c.FontSize = size
	err := c.Context.SetFont(c.GetFont(), "", c.FontSize)
	if err != nil {
		log.Error().Msg(err.Error())
	}
}

func (c *Pdf) AddPage() {
	c.Context.AddPage()
}

func (c *Pdf) SetLineWidth(value float64) {
	c.Context.SetLineWidth(value)
}

func (c *Pdf) SetStrokeColor(color Color) {
	c.Context.SetStrokeColor(color.R, color.G, color.B)
}

func (c *Pdf) SetFillColor(color Color) {
	c.Context.SetFillColor(color.R, color.G, color.B)
}

func (c *Pdf) SetTextColor(color Color) {
	c.TextColor = color
}

func (c *Pdf) FillRect(x float64, y float64, width float64, height float64, color Color) {
	c.Context.SetLineWidth(0.0)
	c.SetStrokeColor(color)
	c.SetFillColor(color)

	c.Context.RectFromUpperLeftWithStyle(x, y, width, height, "FD")
}

func (c *Pdf) TextOut(x float64, y float64, width float64, height float64, str string, align int) {
	c.SetStrokeColor(c.TextColor)
	c.SetFillColor(c.TextColor)

	c.Context.SetXY(x, y)

	err := c.Context.CellWithOption(&gopdf.Rect{W: width, H: height}, str, gopdf.CellOption{
		Align:  align,
		Border: 0,
		Float:  gopdf.Right,
	})
	if err != nil {
		log.Error().Msg(err.Error())
	}
}

func (c *Pdf) Save(filename string) {
	err := c.Context.WritePdf(filename)
	if err != nil {
		log.Error().Msg(err.Error())
	}
}

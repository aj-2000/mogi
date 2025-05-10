package ui

import (
	"github.com/aj-2000/mogi/color"
	"github.com/aj-2000/mogi/math"
)

type Table struct {
	Component
	Rows        []Row
	Header      []string
	HeaderColor color.RGBA
	RowColor    color.RGBA
	FontSize    float32
	FontColor   color.RGBA
}

type Row struct {
	Cells []string
}

func NewTable() *Table {
	return &Table{
		Component:   newComponentBase(TableKind),
		Rows:        []Row{},
		Header:      []string{},
		HeaderColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		RowColor:    color.RGBA{R: 255, G: 255, B: 255, A: 255},
		FontSize:    16,
		FontColor:   color.RGBA{R: 0, G: 0, B: 0, A: 255},
	}
}

func (t *Table) SetID(id string) *Table {
	t.Component.setID(id)
	return t
}

func (t *Table) SetDisplay(d Display) *Table {
	t.Component.setDisplay(d)
	return t
}

func (t *Table) SetPosition(p Position) *Table {
	t.Component.setPos(p)
	return t
}

func (t *Table) SetSize(s math.Vec2f32) *Table {
	t.Component.setSize(s)
	return t
}

func (t *Table) SetBackgroundColor(c color.RGBA) *Table {
	t.Component.setBackgroundColor(c)
	return t
}

func (t *Table) SetBorderColor(c color.RGBA) *Table {
	t.Component.setBorderColor(c)
	return t
}

func (t *Table) SetBorderWidth(w float32) *Table {
	t.Component.setBorder(math.Vec2f32{X: w, Y: w})
	return t
}

func (t *Table) SetBorderRadius(r float32) *Table {
	t.Component.setBorderRadius(r)
	return t
}

func (t *Table) SetFontSize(s float32) *Table {
	t.FontSize = s
	return t
}

func (t *Table) SetFontColor(c color.RGBA) *Table {
	t.FontColor = c
	return t
}

func (t *Table) SetHeader(header []string) *Table {
	t.Header = header
	return t
}

func (t *Table) SetHeaderColor(c color.RGBA) *Table {
	t.HeaderColor = c
	return t
}

func (t *Table) SetRowColor(c color.RGBA) *Table {
	t.RowColor = c
	return t
}

func (t *Table) AddRow(row Row) *Table {
	t.Rows = append(t.Rows, row)
	return t
}

func (t *Table) AddRows(rows []Row) *Table {
	t.Rows = append(t.Rows, rows...)
	return t
}

func (t *Table) SetZIndex(zIndex int) *Table {
	t.Component.setZIndex(zIndex)
	return t
}

func (t *Table) Kind() ComponentKind {
	return t.Component.kind
}

package main

import (
	"image/color"
)

type Label struct {
	H, W int16
	Buf  []uint16
}

func NewLabel(h, w int16) *Label {
	return &Label{
		Buf: make([]uint16, h*w),
		H:   h,
		W:   w,
	}
}

func (l *Label) Display() error {
	return nil
}

func (l *Label) Size() (x, y int16) {
	return l.W, l.H
}

func (l *Label) SetPixel(x, y int16, c color.RGBA) {
	l.Buf[y*l.W+x] = RGBATo565(c)
}

func (l *Label) FillScreen(c color.RGBA) {
	for i := range l.Buf {
		l.Buf[i] = RGBATo565(c)
	}
}

func RGBATo565(c color.RGBA) uint16 {
	r, g, b, _ := c.RGBA()
	return uint16((r & 0xF800) +
		((g & 0xFC00) >> 5) +
		((b & 0xF800) >> 11))
}

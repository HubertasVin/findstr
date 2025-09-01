package models

import "image/color"

type Style struct {
	Fg   color.RGBA
	Bg   color.RGBA
	Bold bool
}

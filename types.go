package main

import (
	"mime/multipart"
)

type CFListResponse struct {
	Result   Images   `json:"result"`
	Success  bool     `json:"success"`
	Errors   []string `json:"errors"`
	Messages []string `json:"messages"`
}

type CFPostResponse struct {
	Result   CFImage  `json:"result"`
	Success  bool     `json:"success"`
	Errors   []string `json:"errors"`
	Messages []string `json:"messages"`
}

type Images struct {
	Images            []CFImage `json:"images"`
	ContinuationToken bool      `json:"continuation_token:"`
}

type CFImage struct {
	ImageId           string   `json:"id"`
	Title             string   `json:"title"`
	Description       string   `json:"description"`
	Filename          string   `json:"filename"`
	Uploaded          string   `json:"uploaded"`
	RequireSignedURLs bool     `json:"requireSignedURLs"`
	Variants          []string `json:"variants"`
	Category          int      `json:"category"`
}

type ImageUploadDTO struct {
	Title       string                `form:"title"`
	Description string                `form:"description"`
	Category    string                `form:"category" binding:"required"`
	File        *multipart.FileHeader `form:"file" binding:"required"`
}

type Category int

const (
	Uncategorized Category = iota
	Photography
	Sculpture
	Digital
	Print
)

func (c Category) String() string {
	names := []string{"uncategorized", "photography", "sculpture", "digital", "print"}
	if c < 0 || int(c) >= len(names) {
		return "Uncategorized"
	}
	return names[c]
}

var categoryFromString = map[string]Category{
	"uncategorized": Uncategorized,
	"photography":   Photography,
	"sculpture":     Sculpture,
	"digital":       Digital,
	"print":         Print,
}

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
	Result   Image    `json:"result"`
	Success  bool     `json:"success"`
	Errors   []string `json:"errors"`
	Messages []string `json:"messages"`
}

type Images struct {
	Images            []Image `json:"images"`
	ContinuationToken bool    `json:"continuation_token:"`
}

type Image struct {
	ImageId           string   `json:"id"`
	Title             string   `json:"title"`
	Description       string   `json:"description"`
	Filename          string   `json:"filename"`
	Uploaded          string   `json:"uploaded"`
	RequireSignedURLs bool     `json:"requireSignedURLs"`
	Variants          []string `json:"variants"`
	Category          string   `json:"category"`
}

type UploadDTO struct {
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
	names := []string{"Uncategorized", "Photography", "Sculpture", "Digital", "Print"}
	if c < 0 || int(c) >= len(names) {
		return "Uncategorized"
	}
	return names[c]
}

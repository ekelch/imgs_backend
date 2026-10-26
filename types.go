package main

type ImageResponse struct {
	Result   Images   `json:"result"`
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
	Filename          string   `json:"filename"`
	Uploaded          string   `json:"uploaded"`
	RequireSignedURLs bool     `json:"requireSignedURLs"`
	Variants          []string `json:"variants"`
	Category          int      `json:"category"`
}

type Category int

const (
	Uncategorized Category = iota
	Photography
	Sculpture
	Digital
	Print
)

package main

type ImageResponse struct {
	Result Images `json:"result"`
	Success bool `json:"success"`
	Errors []string `json:"errors"`
	Messages []string `json:"messages"`
}

type Images struct {
	Images []Image `json:"images"`
	ContinuationToken bool `json:"continuation_token:"`
}

type Image struct {
	Id string `json:"id"`
	Filename string `json:"filename"`
	Uploaded string `json:"uploaded"`
	RequireSignedURLs bool `json:"requireSignedURLs"`
	Variants []string `json:"variants"`
}

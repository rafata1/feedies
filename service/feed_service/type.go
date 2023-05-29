package feed_service

type GetNewsOutput struct {
	News []New `json:"news"`
}

type New struct {
	SourceURL     string `json:"source_url"`
	SourceLogoURL string `json:"source_logo_url"`
	Title         string `json:"title"`
	ImageURL      string `json:"image_url"`
	CreatedAt     string `json:"created_at"`
}

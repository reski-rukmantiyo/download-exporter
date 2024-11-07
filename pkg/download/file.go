package download

type ImageConfig struct {
	Location       string          `yaml:"location"`
	ImageDownloads []ImageDownload `yaml:"image_downloads"`
}

type ImageDownload struct {
	Image string `yaml:"image"`
	Label string `yaml:"label"`
}

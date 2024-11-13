package download

type ImageConfig struct {
	Location       string          `yaml:"location"`
	ImageDownloads []ImageDownload `yaml:"image_downloads"`
	MinuteToPull   int             `yaml:"minute_time_to_pull"`
	ContainerType  string          `yaml:"container_type"`
}

type ImageDownload struct {
	Image string `yaml:"image"`
	Label string `yaml:"label"`
}

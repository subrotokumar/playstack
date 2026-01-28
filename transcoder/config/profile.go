package config

type QualityProfile struct {
	Name         string
	Resolution   string
	VideoBitrate string
	AudioBitrate string
	VideoCodec   string
	AudioCodec   string
	Preset       string
}

var Profiles = []QualityProfile{
	{
		Name:         "1080p",
		Resolution:   "1920x1080",
		VideoBitrate: "5000k",
		AudioBitrate: "192k",
		VideoCodec:   "libx264",
		AudioCodec:   "aac",
		Preset:       "medium",
	},
	{
		Name:         "720p",
		Resolution:   "1280x720",
		VideoBitrate: "3000k",
		AudioBitrate: "128k",
		VideoCodec:   "libx264",
		AudioCodec:   "aac",
		Preset:       "medium",
	},
	{
		Name:         "480p",
		Resolution:   "854x480",
		VideoBitrate: "1500k",
		AudioBitrate: "128k",
		VideoCodec:   "libx264",
		AudioCodec:   "aac",
		Preset:       "fast",
	},
	{
		Name:         "360p",
		Resolution:   "640x360",
		VideoBitrate: "800k",
		AudioBitrate: "96k",
		VideoCodec:   "libx264",
		AudioCodec:   "aac",
		Preset:       "fast",
	},
}

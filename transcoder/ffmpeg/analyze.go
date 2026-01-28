package ffmpeg

import (
	"encoding/json"
	"os/exec"
)

func AnalyzeVideo(inputPath string) (*VideoInfo, error) {
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		inputPath,
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var info VideoInfo
	json.Unmarshal(output, &info)

	return &info, nil
}

type VideoInfo struct {
	Format struct {
		Duration string `json:"duration"`
		Size     string `json:"size"`
		BitRate  string `json:"bit_rate"`
	} `json:"format"`
	Streams []struct {
		CodecType string `json:"codec_type"`
		Width     int    `json:"width"`
		Height    int    `json:"height"`
		FrameRate string `json:"r_frame_rate"`
	} `json:"streams"`
}

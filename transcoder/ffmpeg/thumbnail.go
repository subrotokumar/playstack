package ffmpeg

import (
	"os/exec"
)

func GenerateThumbnail(input string, output string, timestamp string) error {
	cmd := exec.Command("ffmpeg",
		"-i", input,
		"-ss", timestamp, // Timestamp (e.g., "00:00:05")
		"-vframes", "1", // Single frame
		"-vf", "scale=640:360", // Thumbnail size
		"-q:v", "2", // Quality
		output,
	)

	return cmd.Run()
}

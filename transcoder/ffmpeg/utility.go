package ffmpeg

func HLS_CMD(inputPath, outputDir string) []string {
	return []string{
		"ffmpeg",
		"-i", inputPath,

		"-filter_complex",
		"[0:v]split=3[v1][v2][v3];" +
			"[v1]scale=640:360:flags=fast_bilinear[360p];" +
			"[v2]scale=1280:720:flags=fast_bilinear[720p];" +
			"[v3]scale=1920:1080:flags=fast_bilinear[1080p]",

		// 360p video stream
		"-map", "[360p]",
		"-map", "[720p]",
		"-map", "[1080p]",

		"-c:v", "libx264",
		"-preset", "veryfast",
		"-profile:v", "high",
		"-level:v", "4.1",

		"-g", "48",
		"-keyint_min", "48",
		"-sc_threshold", "0",

		"-b:v:0", "1000k",
		"-b:v:1", "4000k",
		"-b:v:2", "8000k",

		"-f", "hls",
		"-hls_time", "6",
		"-hls_playlist_type", "vod",
		"-hls_flags", "independent_segments",
		"-hls_segment_type", "mpegts",
		"-hls_list_size", "0",

		"-master_pl_name", "master.m3u8",
		"-var_stream_map", "v:0 v:1 v:2",

		"-hls_segment_filename",
		outputDir + "/%v/segment_%03d.ts",

		outputDir + "/%v/playlist.m3u8",
	}
}

func DASH_CMD(inputPath, outputDir string) []string {
	return []string{
		"ffmpeg",
		"-i", inputPath,

		"-filter_complex",
		"[0:v]split=3[v1][v2][v3];" +
			"[v1]scale=640:360:flags=fast_bilinear[360p];" +
			"[v2]scale=1280:720:flags=fast_bilinear[720p];" +
			"[v3]scale=1920:1080:flags=fast_bilinear[1080p]",

		// 360p
		"-map", "[360p]",
		"-c:v:0", "libx264",
		"-b:v:0", "1000k",

		// 720p
		"-map", "[720p]",
		"-c:v:1", "libx264",
		"-b:v:1", "4000k",

		// 1080p
		"-map", "[1080p]",
		"-c:v:2", "libx264",
		"-b:v:2", "8000k",

		// Shared video settings
		"-preset", "veryfast",
		"-profile:v", "high",
		"-level:v", "4.1",
		"-g", "48",
		"-keyint_min", "48",
		"-sc_threshold", "0",

		// Audio
		"-map", "a:0?",
		"-c:a", "aac",
		"-b:a", "128k",

		// DASH settings
		"-use_timeline", "1",
		"-use_template", "1",
		"-window_size", "5",
		"-seg_duration", "6",
		"-adaptation_sets", "id=0,streams=v id=1,streams=a",

		"-f", "dash",
		outputDir + "/manifest.mpd",
	}
}

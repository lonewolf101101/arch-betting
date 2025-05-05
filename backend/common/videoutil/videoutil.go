package videoutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
)

type ffProbeResolution struct {
	Streams []struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"streams"`
}

func GetResolution(videoPath string) (width int, height int, error error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "json", videoPath)
	stdout, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}

	var res *ffProbeResolution
	if err := json.Unmarshal(stdout, &res); err != nil {
		return 0, 0, err
	}

	if len(res.Streams) == 0 {
		return 0, 0, errors.New("no resolution found")
	} else if len(res.Streams) > 1 {
		log.Println("Resolution FFProbe returning multiple streams. Which is unexpected!")
	}

	return res.Streams[0].Width, res.Streams[0].Height, nil
}

/*
# Generates thumbnail of specific time from video source

# Note: Replaces customer's current subscription

Gives a subscription according to calculation below:
  - balance = plan.balance * quantity
  - duration = plan.duration * quantity
*/
func GenerateThumbnail(sourcePath, exportPath string, second int) error {
	args := []string{
		"-i",
		sourcePath,
		"-ss",
		strconv.Itoa(second),
		"-vframes",
		"1",
		exportPath,
	}

	cmd := exec.Command("ffmpeg", args...)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Problem generating thumbnail", stdout, err)
		return err
	}
	return nil
}

type ExportOption struct {
	AssPath          string     // ASS format file for burning subtitle into video
	Source           string     // Source video path
	Target           string     // Target video path
	MaxRes           int        // Max res is maximum resolution of width or height based on `IsHorizontal` field
	IsHorizontal     bool       // Is the video horizontal or vertical. If its given true height will be `MaxRes`, otherwise width
	CallbackURL      string     // Can be used for receiving export progress
	Fps              int        // Frame Per Second
	WatermarkPath    string     // Places watermark image on the middle right
	HWA              bool       // Enable hardware accelerate
	HighQualityFlags bool       // Wanna get your high quality video, set this to true
	HWAOption        *HWAOption // Hardware acceleration options
}

type HWAOption struct {
	GPUID  int // GPU No: Index starts with 0
	Preset string
	CRF    int
}

// # ExportVideoFFMPEG exports a video using the FFMPEG command-line tool with the specified options.
//
// This function constructs an FFMPEG command based on the provided ExportOption struct,
// allowing for various configurations such as hardware acceleration, scaling, watermarking,
// frame rate adjustment, and high-quality encoding settings.
//
// # Behavior:
//   - If hardware acceleration (HWA) is enabled, the function uses CUDA for encoding.
//   - The input video is scaled based on the specified resolution and orientation (horizontal or vertical).
//   - If a watermark is provided, it is overlaid on the video.
//   - The function supports additional configurations such as frame rate, high-quality encoding flags,
//     and progress callback URL.
//
// Returns:
//   - error: Returns an error if the FFMPEG command fails to execute or encounters an issue.
//
// Example Usage:
//
//	err := ExportVideoFFMPEG(&ExportOption{
//	    Source:         "input.mp4",
//	    Target:         "output.mp4",
//	    MaxRes:         1080,
//	    IsHorizontal:   true,
//	    WatermarkPath:  "watermark.png",
//	    HWA:            true,
//	    HWAOption:      HWAOption{Preset: "fast", CRF: 23, GPUID: "0"},
//	    HighQualityFlags: true,
//	})
//	if err != nil {
//	    log.Fatalf("Failed to export video: %v", err)
//	}
func ExportVideoFFMPEG(option *ExportOption) error {
	filterComplex := "[0:v]"

	var args []string

	if option.HWA {
		args = append(args, "-hwaccel", "cuda")
	}

	args = append(args, "-y", "-i", option.Source)

	if option.AssPath != "" {
		filterComplex += fmt.Sprintf("ass=%v,", option.AssPath)
	}

	if option.IsHorizontal {
		filterComplex += fmt.Sprintf("scale=%v:-2", option.MaxRes)
	} else {
		filterComplex += fmt.Sprintf("scale=-2:%v", option.MaxRes)
	}

	if option.WatermarkPath != "" {
		args = append(args, []string{"-i", option.WatermarkPath}...)

		if option.IsHorizontal {
			filterComplex += fmt.Sprintf("[video];[1:v]scale=%v*0.08:-2[logo];[video][logo]overlay=W-w-30:H/4", option.MaxRes)
		} else {
			filterComplex += fmt.Sprintf("[video];[1:v]scale=-2:%v*0.08[logo];[video][logo]overlay=W-w-30:H/4", option.MaxRes)
		}
	}

	args = append(args, "-pix_fmt", "yuv420p", "-filter_complex", filterComplex)

	if option.Fps != 0 {
		args = append(args, []string{"-r", strconv.Itoa(option.Fps)}...)
	}

	if option.HWA {
		args = append(args, "-preset", option.HWAOption.Preset, "-crf", strconv.Itoa(option.HWAOption.CRF), "-vcodec", "h264_nvenc")
	}

	if option.HighQualityFlags {
		args = append(args, "-rc:v", "vbr", "-b:v", "15M", "-maxrate:v", "30M", "-profile:v", "baseline")
	}

	if option.CallbackURL != "" {
		args = append(args, "-progress", option.CallbackURL)
	}

	args = append(args, "-movflags", "+faststart", "-c:a", "aac", option.Target)

	log.Println("ffmpeg", args)

	cmd := exec.Command("ffmpeg", args...)
	if option.HWA {
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, fmt.Sprintf("CUDA_VISIBLE_DEVICES=%v", option.HWAOption.GPUID))
	}

	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Problem exporting video", string(stdout), err)
		return err
	}

	return nil
}

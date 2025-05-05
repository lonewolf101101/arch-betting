package audioutil

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func GetSizeAndDuration(audioPath string) (int, int, error) {
	wavFile, err := os.Stat(audioPath)
	if err != nil {
		return 0, 0, err
	}
	size := int(wavFile.Size())
	duration := int(size * 8 / 16000 / 16)
	return size, duration, nil
}

func SecToTimeFormat(input int) string {
	hour := int(input / 3600)
	minute := int(input/60) % 60
	second := int(input) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)
}

func TimeFormatToSec(input string) int {
	parts := strings.Split(input, ":")
	if len(parts) != 3 {
		log.Println("invalid time:", input)
		return 0
	}

	hour, _ := strconv.Atoi(parts[0])
	minute, _ := strconv.Atoi(parts[1])
	second, _ := strconv.Atoi(parts[2])

	return hour*3600 + minute*60 + second
}

func ConvertToWav(srcPath, dstPath string) error {
	args := []string{
		"-y",
		"-i",
		srcPath,
		"-acodec",
		"pcm_s16le",
		"-ar",
		"16000",
		"-ac",
		"1",
		dstPath,
	}
	cmd := exec.Command("ffmpeg", args...)
	if stdout, err := cmd.CombinedOutput(); err != nil {
		log.Println("Problem converting to wav!", string(stdout), err)
		return err
	}

	// args = []string{
	// 	"-y",
	// 	"-i",
	// 	tempFilePath,
	// 	"-map_channel",
	// 	"0.0.0",
	// 	dstPath,
	// }
	// convertMonoCmd := exec.Command("ffmpeg", args...)
	// if stdout, err := convertMonoCmd.Output(); err != nil {
	// 	s.errorLog.Printf("Problem converting to mono!, out: %+v, err: %s", stdout, err)
	// 	return err
	// }

	// if err := os.Remove(tempFilePath); err != nil {
	// 	s.errorLog.Printf("Couldn't delete unnessary converted audio, err: %s", err)
	// 	return err
	// }

	return nil
}

func MergeAudios(srcPaths []string, dstPath string) error {
	args := []string{}

	for _, srcPath := range srcPaths {
		args = append(args, "-i", srcPath)
	}

	args = append(
		args,
		"-filter_complex",
		fmt.Sprintf("amix=inputs=%v:duration=longest:dropout_transition=0:normalize=0", len(srcPaths)),
		dstPath,
	)

	log.Println("ffmpeg", args)
	cmd := exec.Command("ffmpeg", args...)
	if stdout, err := cmd.Output(); err != nil {
		log.Println("Problem converting to wav!", stdout, err)
		return err
	}

	return nil
}

func ConcatAudios(srcPaths []string, dstPath string) error {
	args := []string{}

	part := ""
	for i, srcPath := range srcPaths {
		args = append(args, "-i", srcPath)
		part += fmt.Sprintf("[%v:a]", i)
	}

	args = append(
		args,
		"-filter_complex",
		fmt.Sprintf("%vconcat=n=%v:v=0:a=1", part, len(srcPaths)),
		dstPath,
	)

	log.Println("ffmpeg", args)
	cmd := exec.Command("ffmpeg", args...)
	if stdout, err := cmd.Output(); err != nil {
		log.Println("Problem converting to wav!", stdout, err)
		return err
	}

	return nil
}

// Not working ...
// func ConcatAudios(srcPaths []string, dstPath string) error {
// 	args := []string{
// 		"-i",
// 		"concat:" + strings.Join(srcPaths, "|"),
// 		"-c",
// 		"copy",
// 		dstPath,
// 	}

// 	log.Println("ffmpeg", args)
// 	cmd := exec.Command("ffmpeg", args...)
// 	if stdout, err := cmd.Output(); err != nil {
// 		log.Println("Problem converting to wav!", stdout, err)
// 		return err
// 	}

// 	return nil
// }

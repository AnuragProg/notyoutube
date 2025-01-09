package utils

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	utilsType "github.com/anuragprog/notyoutube/preprocessor-service/types/utils"
)

// INFO: relative filenames were not producing any results, so try using full pathnames only
//       Not converting to abs because test files take their locations as base
func GetVideoResolution(ctx context.Context, filename string) (utilsType.VideoInfo, error) {

	// check for absolute path existence
	if !filepath.IsAbs(filename) {
		return utilsType.VideoInfo{}, errors.New("only absolute filenames are allowed")
	}

	// check for file's existence
	_, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err){
		return utilsType.VideoInfo{}, errors.New("file not found")
	}

	cmd := exec.CommandContext(ctx, "/bin/sh")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return utilsType.VideoInfo{}, err
	}
	defer stdin.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return utilsType.VideoInfo{}, err
	}
	defer stdout.Close()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return utilsType.VideoInfo{}, err
	}
	defer stderr.Close()

	if err := cmd.Start(); err != nil {
		return utilsType.VideoInfo{}, err
	}

	errChan := make(chan error)
	defer close(errChan)
	go func(){
		stderrReader := bufio.NewScanner(stderr)
		errMessage := ""

		for {
			hasNextToken := stderrReader.Scan()

			if hasNextToken { // if a token found then append to errMessage
				errMessage += stderrReader.Text()
			}else{ 
				if err := stderrReader.Err(); err != nil { // an error happened while reading error stream
					errChan<- err
					return
				} 

				// or EOF reached i.e. successfully retrieved err message
				if len(errMessage) > 0 { // there was no error
					errChan<- errors.New(errMessage)
				}				
				return
			}
		}
	}()

	doneChan := make(chan utilsType.VideoInfo)
	defer close(doneChan)
	go func(){ // this scope handles the setting of above videoResolution object
		defer stdin.Close()
		defer stdout.Close()
		stdoutReader := bufio.NewReader(stdout)

		videoResolution := utilsType.VideoInfo{}

		// frame - expected o/p - (integer)
		videoResolution.TotalFrames, err = processTotalFrames(stdin, stdoutReader, filename)
		if err != nil {
			errChan<- err
			return
		}

		// aspect ratio - expected o/p - (integer):(integer)
		videoResolution.AspectRatio, err = processAspectRatio(stdin, stdoutReader, filename)
		if err != nil {
			errChan<- err
			return
		}

		// bitrate - expected o/p - (integer)
		videoResolution.Bitrate, err = processBitrate(stdin, stdoutReader, filename)
		if err != nil {
			errChan<- err
			return
		}

		// resolution - expected o/p - (integer)x(integer)
		videoResolution.Width, videoResolution.Height, err = processResolution(stdin, stdoutReader, filename)
		if err != nil {
			errChan<- err
			return
		}
		doneChan<- videoResolution
	}()


	select {
	case err := <-errChan:
		return utilsType.VideoInfo{}, err
	case videoResolution := <-doneChan:
		if err := cmd.Wait(); err != nil {
			return utilsType.VideoInfo{}, err
		}
		return videoResolution, nil
	}
}


func processTotalFrames(stdin io.Writer, stdout *bufio.Reader, filename string) (uint32, error) {
	// fmt.Println("executing ffprobe command")
	frameCommand := fmt.Sprintf("ffprobe -v error -select_streams v:0 -count_frames -show_entries stream=nb_read_frames -of default=noprint_wrappers=1:nokey=1 %v", filename)
	if _, err := stdin.Write([]byte(frameCommand + "\n")); err != nil { return 0, err }
	// fmt.Println("reading from stdout")
	frameOutput, _, err := stdout.ReadLine()
	if err != nil {
		return 0, err
	}
	// fmt.Println("frameoutput=",string(frameOutput))
	frames, err := strconv.Atoi(string(frameOutput))
	if err != nil {
		return 0, err
	}
	return uint32(frames), nil
}

func processAspectRatio(stdin io.Writer, stdout *bufio.Reader, filename string) (float32, error) {
	aspectRatioCommand := fmt.Sprintf("ffprobe -v error -select_streams v:0 -show_entries stream=display_aspect_ratio -of default=noprint_wrappers=1:nokey=1 %v", filename)
	if _, err := stdin.Write([]byte(aspectRatioCommand + "\n")); err != nil { return 0, err }
	aspectRatioOutput, _, err := stdout.ReadLine()
	if err != nil {
		return 0, err
	}
	// fmt.Println("aspectratiooutput=",string(aspectRatioOutput))
	parsedAspectRatio := strings.Split(string(aspectRatioOutput), ":")
	if len(parsedAspectRatio) != 2 {
		return 0, errors.New("coudn't parse aspect ratio")
	}
	numerator, err := strconv.Atoi(parsedAspectRatio[0])
	if err != nil {
		return 0, err
	}
	denomenator, err := strconv.Atoi(parsedAspectRatio[1])
	if err != nil {
		return 0, err
	}
	if denomenator == 0 {
		return 0, errors.New(fmt.Sprintf("invalid aspect ratio: %v/%v", numerator, denomenator))
	}
	return float32(numerator)/float32(denomenator), nil
}

func processBitrate(stdin io.Writer, stdout *bufio.Reader, filename string) (uint32, error) {
	bitrateCommand := fmt.Sprintf("ffprobe -v error -select_streams v:0 -show_entries stream=bit_rate -of default=noprint_wrappers=1:nokey=1 %v", filename)
	if _, err := stdin.Write([]byte(bitrateCommand + "\n")); err != nil { return 0, err }
	bitrateOutput, _, err := stdout.ReadLine()
	if err != nil {
		return 0, err
	}
	// fmt.Println("bitrateoutput=", string(bitrateOutput))
	bitrate, err := strconv.ParseUint(string(bitrateOutput), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(bitrate), nil
}

func processResolution(stdin io.Writer, stdout *bufio.Reader, filename string)(width, height uint32, err error) {

	resolutionCommand := fmt.Sprintf("ffprobe -v error -select_streams v:0 -show_entries stream=width,height -of csv=s=x:p=0 %v", filename)
	if _, err = stdin.Write([]byte(resolutionCommand + "\n")); err != nil { 
		return
	}
	resolutionOutput, _, err := stdout.ReadLine()
	if err != nil {
		return
	}
	// fmt.Println("resolutionOutput=", string(resolutionOutput))
	parsedResolution := strings.Split(string(resolutionOutput),"x")
	if len(parsedResolution) != 2 {
		err = errors.New("couldn't parse resolution")
		return
	}
	width64, err := strconv.ParseUint(parsedResolution[0], 10, 32)
	if err != nil {
		return
	}
	height64, err := strconv.ParseUint(parsedResolution[1], 10, 32)
	if err != nil {
		return
	}
	width = uint32(width64)
	height = uint32(height64)
	return
}



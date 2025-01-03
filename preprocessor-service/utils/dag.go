package utils

import (
	"context"
	"fmt"
	"time"

	mqType "github.com/anuragprog/notyoutube/preprocessor-service/types/mq"
	"github.com/google/uuid"
)

// Contains criteria for deciding which resolution to make video for
var targetVideoResolutions = []struct {
	Encoding               mqType.VideoEncoding
	Width, Height, Bitrate uint32
	AspectRatio            float32
}{
	{
		Encoding: mqType.VideoEncoding_P144,
		Width:    256, Height: 144,
		Bitrate:     150_000,
		AspectRatio: float32(16) / float32(9),
	},
	{
		Encoding: mqType.VideoEncoding_P240,
		Width:    426, Height: 240,
		Bitrate:     240_000,
		AspectRatio: float32(16) / float32(9),
	},
	{
		Encoding: mqType.VideoEncoding_P360,
		Width:    640, Height: 360,
		Bitrate:     800_000,
		AspectRatio: float32(16) / float32(9),
	},
	{
		Encoding: mqType.VideoEncoding_P480,
		Width:    854, Height: 480,
		Bitrate:     1_200_000,
		AspectRatio: float32(16) / float32(9),
	},
	{
		Encoding: mqType.VideoEncoding_P720,
		Width:    1280, Height: 720,
		Bitrate:     2_500_000,
		AspectRatio: float32(16) / float32(9),
	},
	{
		Encoding: mqType.VideoEncoding_P1080,
		Width:    1920, Height: 1080,
		Bitrate:     5_000_000,
		AspectRatio: float32(16) / float32(9),
	},
}

var targetAsciiResolutions = []struct {
	Encoding           mqType.AsciiEncoding
	Width, Height, Fps uint32
}{
	{
		Encoding: mqType.AsciiEncoding_P208x57,
		Width:    208,
		Height:   57,
		Fps:      60,
	},
}

func CreateDAG(ctx context.Context, filename string) (*mqType.DAG, error) {

	videoExtractorWorker := &mqType.Worker{
		Id:          uuid.NewString(),
		Name:        "video-extractor-worker",
		Description: "extracts video from the video file and encode it to h.264 encoding",
		WorkerType:  mqType.WorkerType_VIDEO_EXTRACTOR,
	}
	audioExtractorWorker := &mqType.Worker{
		Id:          uuid.NewString(),
		Name:        "audio-extractor-worker",
		Description: "extracts audio from the video file and encode it to aac encoding",
		WorkerType:  mqType.WorkerType_AUDIO_EXTRACTOR,
	}
	metadataExtractorWorker := &mqType.Worker{
		Id:          uuid.NewString(),
		Name:        "metadata-extractor-worker",
		Description: "extracts metadata from the video file",
		WorkerType:  mqType.WorkerType_METADATA_EXTRACTOR,
	}

	videoEncoderWorkers := make([]*mqType.Worker, 0, len(targetVideoResolutions))
	sourceInfo, err := GetVideoResolution(ctx, filename)
	if err != nil {
		return nil, err
	}
	sourceWidth := sourceInfo.Width
	sourceHeight := sourceInfo.Height
	sourceAspectRatio := sourceInfo.AspectRatio
	sourceBitrate := sourceInfo.Bitrate
	for _, resolution := range targetVideoResolutions {
		targetEncoding := resolution.Encoding
		targetWidth := resolution.Width
		targetHeight := resolution.Height
		targetBitrate := resolution.Bitrate
		targetAspectRatio := resolution.AspectRatio

		isDimensionCompatible := sourceWidth >= targetWidth && sourceHeight >= targetHeight
		isAspectRatioCompatible := (sourceAspectRatio - targetAspectRatio) < 0.1
		isBitrateCompatible := sourceBitrate >= targetBitrate

		if isDimensionCompatible && isAspectRatioCompatible && isBitrateCompatible {
			videoEncoderWorkers = append(videoEncoderWorkers, &mqType.Worker{
				Id:          uuid.NewString(),
				Name:        fmt.Sprintf("%v-video-encoder-worker", targetEncoding.String()),
				Description: fmt.Sprintf("converts video stream to %v encoding", targetEncoding.String()),
				WorkerType:  mqType.WorkerType_VIDEO_ENCODER,
				WorkerConfig: &mqType.Worker_VideoEncoderConfig{
					VideoEncoderConfig: &mqType.VideoEncoderWorkerConfig{
						Encoding: targetEncoding,
						Width:    targetWidth,
						Height:   targetHeight,
						Bitrate:  targetBitrate,
					},
				},
			})
		}
	}

	asciiEncoderWorkers := make([]*mqType.Worker, 0, len(targetAsciiResolutions))
	for _, resolution := range targetAsciiResolutions {
		// adding all resolutions as even the pathetic video will have decent
		// quality compared to pixel mess in terminal
		asciiEncoderWorkers = append(asciiEncoderWorkers, &mqType.Worker{
			Id:          uuid.NewString(),
			Name:        fmt.Sprintf("%v-ascii-encoder-worker", resolution.Encoding.String()),
			Description: fmt.Sprintf("converts video stream to %v ascii encoding", resolution.Encoding.String()),
			WorkerType:  mqType.WorkerType_ASCII_ENCODER,
			WorkerConfig: &mqType.Worker_AsciiEncoderConfig{
				AsciiEncoderConfig: &mqType.AsciiEncoderWorkerConfig{
					Encoding: resolution.Encoding,
					Width:    resolution.Width,
					Height:   resolution.Height,
					Fps:      resolution.Fps,
				},
			},
		})
	}

	thumbnailGeneratorWorker := &mqType.Worker{
		Id:          uuid.NewString(),
		Name:        "thumbnail-generator-worker",
		Description: "generate thumbnail using video stream",
		WorkerType:  mqType.WorkerType_THUMBNAIL_GENERATOR,
	}
	assemblerWorker := &mqType.Worker{
		Id:          uuid.NewString(),
		Name:        "assembler-worker",
		Description: "assembles together video and audio stream into one single file",
		WorkerType:  mqType.WorkerType_ASSEMBLER,
	}

	// NOTE: Do not make dependencies before this
	//       for the sake of locality of code
	//       Only make nodes before this and group them together
	//       according to the logic that binds them together

	// assemble dag components together
	// 1. Create DAG
	dag := &mqType.DAG{
		Id:           uuid.NewString(),
		CreatedAt:    time.Now().Format(time.RFC3339),
		Workers:      make([]*mqType.Worker, 0),
		Dependencies: make([]*mqType.Dependency, 0),
	}

	// 2. Add Nodes/Workers
	// adding initial 3 workers video, audio and metadata extractor
	dag.Workers = append(dag.Workers, videoExtractorWorker, audioExtractorWorker, metadataExtractorWorker)
	// adding middle stage workers i.e. video encoders
	dag.Workers = append(dag.Workers, videoEncoderWorkers...)
	dag.Workers = append(dag.Workers, asciiEncoderWorkers...)
	// adding terminal nodes i.e. thumbnail generator and assembler
	dag.Workers = append(dag.Workers, thumbnailGeneratorWorker, assemblerWorker)

	// 3. Add Dependencies
	// (video-extractor [1] -- [*] video-encoder...) make all video encoder workers depend on video extractor worker
	for _, worker := range videoEncoderWorkers {
		dag.Dependencies = append(dag.Dependencies, &mqType.Dependency{
			SourceIds: []string{videoExtractorWorker.Id},
			TargetId:  worker.Id,
		})
	}
	// ((video-encoder... + audio-encoder) [set] -- [1] assembler) make assembler worker depend on all video encoder workers and audio encoder
	for _, worker := range videoEncoderWorkers {
		dag.Dependencies = append(dag.Dependencies, &mqType.Dependency{
			SourceIds: []string{audioExtractorWorker.Id, worker.Id},
			TargetId:  assemblerWorker.Id,
		})
	}
	// NOTE: currently putting audio and video together because there could be a way to run audio on terminal as well, problem for future anurag
	// (ascii-encoder... [*] -- [1] assembler) make assembler worker depend on all ascii encoder workers and audio encoder
	for _, worker := range asciiEncoderWorkers {
		dag.Dependencies = append(dag.Dependencies, &mqType.Dependency{
			SourceIds: []string{audioExtractorWorker.Id, worker.Id},
			TargetId:  assemblerWorker.Id,
		})
	}
	// (video-extractor [1] -- [1] thumbnail-extractor) make thumbnail generator depend on video extractor
	dag.Dependencies = append(dag.Dependencies, &mqType.Dependency{
		SourceIds: []string{videoExtractorWorker.Id},
		TargetId:  thumbnailGeneratorWorker.Id,
	})

	return dag, nil
}

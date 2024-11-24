package workers


import(
	mqType "github.com/anuragprog/notyoutube/preprocessor-service/types/mq"

	utilsType "github.com/anuragprog/notyoutube/preprocessor-service/types/utils"
)

var target_video_resolutions = map[string]utilsType.VideoInfo{
	"144p": {
		Width: 256, Height: 144,
		Bitrate: 150_000,   
		AspectRatio: float32(16)/float32(9),
	},
	"240p": {
		Width: 426, Height: 240,
		Bitrate: 240_000,   
		AspectRatio: float32(16)/float32(9),
	},
	"360p": {
		Width: 640, Height: 360,
		Bitrate: 800_000,   
		AspectRatio: float32(16)/float32(9),
	},
	"480p": {
		Width: 854, Height: 480,
		Bitrate: 1_200_000,   
		AspectRatio: float32(16)/float32(9),
	},
	"720p": {
		Width: 1280, Height: 720,
		Bitrate: 2_500_000,   
		AspectRatio: float32(16)/float32(9),
	},
	"1080p": {
		Width: 1920, Height: 1080,
		Bitrate: 5_000_000,   
		AspectRatio: float32(16)/float32(9),
	},
}

var target_ansi_resolutions = []interface{}{}


/*
1. Separation of video and audio
2. Video encoding to different resolutions + terminal versions as well
3. Audio encoding
4. 
*/
func DAGWorker(metadata *mqType.RawVideoMetadata) {

	/*
	Stage1:
		- video extraction
		- audio extraction
		- metadata extraction
	
	Stage2:
		- video encodings [144,240,360,480,720,1080,]
		- audio encoding
	*/

}

/*
function decideResolutionsToGenerate(source_resolution, source_aspect_ratio, source_bitrate, available_resolutions):
    # Step 1: Define the aspect ratios and bitrates for each resolution
    target_resolutions = {
        "144p": { "width": 256, "height": 144, "bitrate": 150, "aspect_ratio": 16/9 },
        "240p": { "width": 426, "height": 240, "bitrate": 400, "aspect_ratio": 16/9 },
        "360p": { "width": 640, "height": 360, "bitrate": 800, "aspect_ratio": 16/9 },
        "480p": { "width": 854, "height": 480, "bitrate": 1200, "aspect_ratio": 16/9 },
        "720p": { "width": 1280, "height": 720, "bitrate": 2500, "aspect_ratio": 16/9 },
        "1080p": { "width": 1920, "height": 1080, "bitrate": 5000, "aspect_ratio": 16/9 }
    }

    # Step 2: Check the source resolution and aspect ratio
    source_width = source_resolution["width"]
    source_height = source_resolution["height"]
    source_aspect_ratio = source_aspect_ratio # Assuming it's already calculated (width/height)
    
    # Step 3: Prepare an empty list for the resolutions to generate
    resolutions_to_generate = []
    
    # Step 4: Loop through each target resolution and decide whether to generate
    for each resolution in target_resolutions:
        target_width = target_resolutions[resolution]["width"]
        target_height = target_resolutions[resolution]["height"]
        target_bitrate = target_resolutions[resolution]["bitrate"]
        target_aspect_ratio = target_resolutions[resolution]["aspect_ratio"]

        # Step 5: Check if the source resolution is greater than or equal to the target resolution
        if (source_width >= target_width and source_height >= target_height):
            # Step 6: Check if the aspect ratios match or are close enough
            if (abs(source_aspect_ratio - target_aspect_ratio) < 0.1):
                # Step 7: Check if the source bitrate is enough for the target resolution
                if (source_bitrate >= target_bitrate):
                    # If all conditions are met, add the resolution to the list
                    resolutions_to_generate.append(resolution)
                else:
                    print("Bitrate too low for " + resolution)
            else:
                print("Aspect ratio mismatch for " + resolution)
        else:
            print("Source resolution too low for " + resolution)

    # Step 8: Return the list of resolutions to generate
    return resolutions_to_generate

# Usage Example
source_resolution = { "width": 1920, "height": 1080 }
source_aspect_ratio = 16/9
source_bitrate = 5000  # kbps
available_resolutions = ["144p", "240p", "360p", "480p", "720p", "1080p"]

decideResolutionsToGenerate(source_resolution, source_aspect_ratio, source_bitrate, available_resolutions)
*/

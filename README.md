
## Ffmpeg commands

- 60fps frames generation with 1280x720 dimension: ffmpeg -i input_video.mp4 -map 0:v -vf "fps=60,scale=1280:720" -qscale:v 2 frame_%04d.jpg


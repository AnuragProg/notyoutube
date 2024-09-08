## Useful Libraries

- Terminal

    - Playing Sounds - ```github.com/faiface/beep```


## Ffmpeg Commands

- Get total frames in a video

    - ```ffprobe -v error -select_streams v:0 -count_frames -show_entries stream=nb_read_frames -of default=noprint_wrappers=1:nokey=1 input.mp4```

    - ```Output: 739```

- Get resolution

    - ```ffprobe -v error -select_streams v:0 -show_entries stream=width,height -of csv=s=x:p=0 input.mp4```

    - ```Output: 1280x720```

- Get Aspect Ratio

    - ```ffprobe -v error -select_streams v:0 -show_entries stream=display_aspect_ratio -of default=noprint_wrappers=1:nokey=1 input.mp4```

    - ```Output: 16:9```

- Get Bitrate 

    - ```ffprobe -v error -select_streams v:0 -show_entries stream=bit_rate -of default=noprint_wrappers=1:nokey=1 input.mp4```

    - ```Output: 1032960```


## Video Encoding

- Youtube

    - [Coding With Lewis - How Instagram Stores BILLIONS of Videos](https://www.youtube.com/watch?v=HzD_Kv6IyQ0)

    - [Leo Isikdogan - How Video Compression Works](https://www.youtube.com/watch?v=QoZ8pccsYo4&t=0s)


- Articles

    - [HLS - HTTP Live Streaming (Apple Docs)](https://developer.apple.com/streaming/)

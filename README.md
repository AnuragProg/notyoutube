
## Schema

- preprocessor -> dag scheduler -> task scheduler -> [workers]*

- rendering specs:
    144p 240p 360p 480p 720p 1080p
    ansi-208x57-good enough

- actual dag processing pipeline:

                 metadata -> save to db

    raw-video -> video -> [video encodings and thumbnail generation]*
                                                                        -> merge video and audio -> final video [video + ansi(encoded) video]
                 audio -> [audio encoding]*

- schema

    

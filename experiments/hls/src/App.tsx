import { useEffect, useRef } from 'react';
import Hls from 'hls.js';
import './App.css';

const videoSrc = 'https://test-streams.mux.dev/x36xhzz/x36xhzz.m3u8';
function App() {
    const videoRef = useRef();

    useEffect(()=>{
        const hls = new Hls({
            debug: true,
        });
        if(Hls.isSupported()){
            hls.log = true;
            hls.loadSource(videoSrc);
            hls.attachMedia(videoRef.current)
            hls.on(Hls.Events.ERROR, (err) => {
                console.log(err)
            });
        }else {
            console.log('not supported');
        }
    }, []);

    return (
        <video
            ref={videoRef}
            controls
            src={videoSrc}
        />
    )
}

export default App

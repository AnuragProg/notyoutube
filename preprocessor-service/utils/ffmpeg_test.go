package utils

import (
	"context"
	"flag"
	"testing"

	"go.uber.org/goleak"
	"github.com/stretchr/testify/assert"
)


var videoFilename string

func init() {
	flag.StringVar(&videoFilename, "filename", "", "Sample Video file for testing")
}

// Currently only tests for errors returned from the GetVideoResolution method
func TestGetVideoResolution(t *testing.T) {

	defer goleak.VerifyNone(t)

	ffmpegShell := NewFFmpegShell()
	resolution, err := ffmpegShell.GetVideoResolution(context.TODO(), videoFilename)
	assert.Nil(t, err, "error happened while getting resolution: %v", err)
	t.Logf("Resolution = %+v", resolution)
}

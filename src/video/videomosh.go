package video

import (
	"fmt"

	"github.com/zergon321/reisen"
)

type VideoMosh struct {
	video *reisen.VideoStream
}

func New(path string) (*VideoMosh, error) {
	media, err := reisen.NewMedia(path)
	if err != nil {
		return nil, err
	}

	video, ok := media.Streams()[0].(*reisen.VideoStream)
	if !ok {
		return nil, fmt.Errorf("error typecasting media to VideoStream")
	}

	vm := &VideoMosh{
		video: video,
	}

	return vm, nil
}

func MoshVideo() {}

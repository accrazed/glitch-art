package video

type VideoMosh struct {
}

func New(path string) (*VideoMosh, error) {

	vm := &VideoMosh{
		// video: video,
	}

	return vm, nil
}

func MoshVideo() {}

package model

type ImageTaskStatus int

func (s ImageTaskStatus) String() string {
	switch s {
	case ImageTaskStatusDefault:
		return "default"
	case ImageTaskStatusProcessing:
		return "processing"
	case ImageTaskStatusCompleted:
		return "completed"
	case ImageTaskStatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

const (
	ImageTaskStatusDefault    ImageTaskStatus = 0
	ImageTaskStatusProcessing ImageTaskStatus = 1
	ImageTaskStatusCompleted  ImageTaskStatus = 2
	ImageTaskStatusFailed     ImageTaskStatus = 3
)

type VideoTaskStatus int

func (s VideoTaskStatus) String() string {
	switch s {
	case VideoTaskStatusDefault:
		return "default"
	case VideoTaskStatusProcessing:
		return "processing"
	case VideoTaskStatusCompleted:
		return "completed"
	case VideoTaskStatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

const (
	VideoTaskStatusDefault    VideoTaskStatus = 0
	VideoTaskStatusProcessing VideoTaskStatus = 1
	VideoTaskStatusCompleted  VideoTaskStatus = 2
	VideoTaskStatusFailed     VideoTaskStatus = 3
)

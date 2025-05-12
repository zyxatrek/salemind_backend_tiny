package model

// ImagePromptKeyword 图像生成关键词
type ImagePromptKeyword struct {
	Pose       string `json:"pose"`
	Location   string `json:"location"`
	TimeOfDay  string `json:"time_of_day"`
	HairColor  string `json:"hair_color"`
	Hairstyle  string `json:"hairstyle"`
	TopWear    string `json:"top_wear"`
	BottomWear string `json:"bottom_wear"`
	LegWear    string `json:"leg_wear"`
}

// ImageTaskRequest 创建图像任务请求
type ImageTaskRequest struct {
	Keyword ImagePromptKeyword `json:"keyword"`
}

type ImageTaskRawRequest struct {
	Prompt string `json:"prompt"`
}

// ImageTaskResponse 创建图像任务响应
type ImageTaskResponse struct {
	ImageTaskID string `json:"image_task_id"`
}

// ImageTaskStatusRequest 图像任务状态查询请求
type ImageTaskStatusRequest struct {
	TaskID string `json:"task_id"`
}

// ImageTaskStatusResponse 图像任务状态响应
type ImageTaskStatusResponse struct {
	Status   string `json:"status"`
	ImageURL string `json:"image_url,omitempty"`
}

// VideoTaskRequest 创建视频任务请求
type VideoTaskRequest struct {
	ImgURL string `json:"img_url"`
}

// VideoTaskResponse 创建视频任务响应
type VideoTaskResponse struct {
	VideoTaskID string `json:"video_task_id"`
}

// VideoTaskStatusRequest 视频任务状态查询请求
type VideoTaskStatusRequest struct {
	TaskID string `json:"task_id"`
}

// VideoTaskStatusResponse 视频任务状态响应
type VideoTaskStatusResponse struct {
	Status   string `json:"status"`
	VideoURL string `json:"video_url,omitempty"`
}

package api

import (
	"net/http"

	"salemind_backend_tiny/model"
	"salemind_backend_tiny/pkg/services"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	genSvc *services.GenerationService
}

func NewHandler(genSvc *services.GenerationService) *Handler {
	return &Handler{
		genSvc: genSvc,
	}
}

// CreateImageTaskRaw 根据提示词生成图像人物
func (h *Handler) CreateImageTaskRaw(c *gin.Context) {
	var req model.ImageTaskRawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 提交图像生成任务
	imageID, err := h.genSvc.SubmitImageTask(req.Prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交图像任务失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.ImageTaskResponse{
		ImageTaskID: imageID,
	})
}

// CreateImageTask 创建图像生成任务
func (h *Handler) CreateImageTask(c *gin.Context) {
	var req model.ImageTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 生成提示词
	prompt, err := h.genSvc.GeneratePrompt(&req.Keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成提示词失败: " + err.Error()})
		return
	}

	// 提交图像生成任务
	imageID, err := h.genSvc.SubmitImageTask(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交图像任务失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.ImageTaskResponse{
		ImageTaskID: imageID,
	})
}

// GetImageTaskStatus 获取图像生成任务状态
func (h *Handler) GetImageTaskStatus(c *gin.Context) {
	var req model.ImageTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.TaskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task_id is required"})
		return
	}

	imageURL, status, err := h.genSvc.GetImageTaskStatus(req.TaskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取图像状态失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.ImageTaskStatusResponse{
		Status:   status.String(),
		ImageURL: imageURL,
	})
}

// CreateVideoTask 创建视频生成任务
func (h *Handler) CreateVideoTask(c *gin.Context) {
	var req model.VideoTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	videoID, err := h.genSvc.CreateVideoTask(req.ImgURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建视频任务失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.VideoTaskResponse{
		VideoTaskID: videoID,
	})
}

// GetVideoTaskStatus 获取视频生成任务状态
func (h *Handler) GetVideoTaskStatus(c *gin.Context) {
	var req model.VideoTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.TaskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task_id is required"})
		return
	}

	videoURL, status, err := h.genSvc.GetVideoaskStatus(req.TaskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取视频状态失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.VideoTaskStatusResponse{
		Status:   status.String(),
		VideoURL: videoURL,
	})
}

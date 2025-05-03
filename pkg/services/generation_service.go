package services

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"salemind_backend_tiny/model"
	"salemind_backend_tiny/pkg/config"
)

type GenerationService struct {
	config *config.Config
}

func NewGenerationService(cfg *config.Config) *GenerationService {
	return &GenerationService{
		config: cfg,
	}
}

func (s *GenerationService) generateSignature(uri, secretKey string) (string, string, string) {
	timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
	nonce := fmt.Sprintf("%d", time.Now().UnixNano())
	content := fmt.Sprintf("%s&%s&%s", uri, timestamp, nonce)

	h := hmac.New(sha1.New, []byte(secretKey))
	h.Write([]byte(content))
	sign := base64.URLEncoding.EncodeToString(h.Sum(nil))
	sign = sign[:len(sign)-1] // Remove padding

	return sign, timestamp, nonce
}

func (s *GenerationService) GeneratePrompt(keyword *model.ImagePromptKeyword) (string, error) {
	originalPrompt := "This is a portrait of Ava standing at the beach at night. She has light green side parted short straight hair and wears grey crop top and black mini pants and black pantyhose."
	chineseInstruction := fmt.Sprintf(
		"请你根据以下中文修改以上描述中的对应内容（姿势、地点、时间、发色、发型、上肢服装、臀部服装、腿部服装），并只返回修改后的描述，依然返回纯英文，只是修改了对应的英文单词，不许输出任何汉字：%s、%s、%s、%s、%s、%s、%s、%s",
		keyword.Pose, keyword.Location, keyword.TimeOfDay, keyword.HairColor, keyword.Hairstyle, keyword.TopWear, keyword.BottomWear, keyword.LegWear,
	)

	// 调用Qwen生成提示词
	payload := map[string]interface{}{
		"model": "qwen-plus",
		"messages": []map[string]string{
			{"role": "system", "content": "You are a helpful assistant."},
			{"role": "user", "content": originalPrompt + " " + chineseInstruction},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("生成提示词失败: %v\n", err)
		return "", err
	}

	req, err := http.NewRequest("POST", s.config.Qwen.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("生成提示词失败: %v\n", err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.Qwen.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("生成提示词失败: %v\n", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("生成提示词失败: %v\n", err)
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("生成提示词失败: %v\n", err)
		return "", err
	}

	choices := result["choices"].([]interface{})
	message := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	customPrompt := message["content"].(string)

	return customPrompt, nil
}

func (s *GenerationService) SubmitImageTask(prompt string) (string, error) {
	uri := "/api/generate/webui/text2img"
	sign, timestamp, nonce := s.generateSignature(uri, s.config.Liblibai.SecretKey)

	params := fmt.Sprintf("?AccessKey=%s&Signature=%s&Timestamp=%s&SignatureNonce=%s",
		s.config.Liblibai.AccessKey, sign, timestamp, nonce)

	url := s.config.Liblibai.APIURL + params

	requestParams := map[string]interface{}{
		"templateUuid": "6f7c4652458d4802969f8d089cf5b91f",
		"generateParams": map[string]interface{}{
			"checkPointId": "ddb974b0717e4f96bc7789a068004831",
			"vaeId":        "",
			"prompt":       prompt,
			"clipSkip":     2,
			"steps":        20,
			"width":        1000,
			"height":       1600,
			"imgCount":     1,
			"seed":         -1,
			"restoreFaces": 0,
			"additionalNetwork": []map[string]interface{}{
				{"modelId": "4e6783ca937047bcbba4874cedd1f552", "weight": 0.4},
				{"modelId": "1aaebe53f31e489993ee5c4338e450e4", "weight": 0.5},
				{"modelId": "da6b0f7fb7004bddb7bc8d3fe698641f", "weight": 0.45},
			},
		},
	}

	jsonData, err := json.Marshal(requestParams)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if code, ok := result["code"].(float64); ok && code == 0 {
		if data, ok := result["data"].(map[string]interface{}); ok {
			if uuid, ok := data["generateUuid"].(string); ok {
				return uuid, nil
			}
		}
	}

	return "", fmt.Errorf("failed to submit image task: %s", string(body))
}

func (s *GenerationService) WaitImageResult(uUID string) (string, error) {
	uri := "/api/generate/webui/status"

	for {
		sign, timestamp, nonce := s.generateSignature(uri, s.config.Liblibai.SecretKey)
		params := fmt.Sprintf("?AccessKey=%s&Signature=%s&Timestamp=%s&SignatureNonce=%s",
			s.config.Liblibai.AccessKey, sign, timestamp, nonce)

		url := s.config.Liblibai.QueryURL + params

		reqBody := map[string]string{"generateUuid": uUID}
		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			return "", err
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return "", err
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return "", err
		}

		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			return "", err
		}

		if code, ok := result["code"].(float64); ok && code == 0 {
			if data, ok := result["data"].(map[string]interface{}); ok {
				status := data["generateStatus"].(float64)
				if images, ok := data["images"].([]interface{}); ok && len(images) > 0 {
					if image, ok := images[0].(map[string]interface{}); ok {
						if auditStatus, ok := image["auditStatus"].(float64); ok && auditStatus == 3 {
							if imageURL, ok := image["imageUrl"].(string); ok {
								return imageURL, nil
							}
						}
					}
				}
				if status == 4 || status == 5 {
					return "", fmt.Errorf("generation failed or blocked")
				}
			}
		}

		time.Sleep(5 * time.Second)
	}
}

func (s *GenerationService) CreateVideoTask(imgURL string) (string, error) {

	prompt := "美女轻轻摇摆着，表情魅惑而动人。她微微左右扭动着胯部，看起来很惬意。"

	url := s.config.Video.APIURL

	payload := map[string]interface{}{
		"model": "wanx2.1-i2v-plus",
		"input": map[string]string{
			"prompt":  prompt,
			"img_url": imgURL,
		},
		"parameters": map[string]interface{}{
			"resolution":    "720P",
			"duration":      5,
			"prompt_extend": true,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.Qwen.APIKey)
	req.Header.Set("X-DashScope-Async", "enable")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if output, ok := result["output"].(map[string]interface{}); ok {
		if taskID, ok := output["task_id"].(string); ok {
			return taskID, nil
		}
	}

	return "", fmt.Errorf("failed to create video task: %s", string(body))
}

func (s *GenerationService) PollVideo(taskID string) (string, error) {
	url := fmt.Sprintf("%s/%s", s.config.Video.TaskURL, taskID)

	for i := 0; i < 60; i++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return "", err
		}

		req.Header.Set("Authorization", "Bearer "+s.config.Qwen.APIKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return "", err
		}

		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			return "", err
		}

		if output, ok := result["output"].(map[string]interface{}); ok {
			if status, ok := output["task_status"].(string); ok {
				if status == "SUCCEEDED" {
					if videoURL, ok := output["video_url"].(string); ok {
						return videoURL, nil
					}
				} else if status == "FAILED" || status == "CANCELED" || status == "UNKNOWN" {
					return "", fmt.Errorf("video generation failed with status: %s", status)
				}
			}
		}

		time.Sleep(10 * time.Second)
	}

	return "", fmt.Errorf("video generation timeout")
}

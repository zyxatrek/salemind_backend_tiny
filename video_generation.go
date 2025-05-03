package main

import (
	"bufio"
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// 配置结构体
type Config struct {
	Qwen struct {
		APIKey  string
		BaseURL string
	}

	Liblibai struct {
		AccessKey string
		SecretKey string
		APIURL    string
		QueryURL  string
	}

	Video struct {
		APIURL  string
		TaskURL string
	}
}

// 用户输入处理函数
func withDefault(prompt, defaultValue string) string {
	fmt.Printf("%s（如：%s）: ", prompt, defaultValue)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue
	}
	return input
}

// 生成签名
func generateSignature(uri, secretKey string) (string, string, string) {
	timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
	nonce := fmt.Sprintf("%d", time.Now().UnixNano())
	content := fmt.Sprintf("%s&%s&%s", uri, timestamp, nonce)

	h := hmac.New(sha1.New, []byte(secretKey))
	h.Write([]byte(content))
	sign := base64.URLEncoding.EncodeToString(h.Sum(nil))
	sign = sign[:len(sign)-1] // 移除填充

	return sign, timestamp, nonce
}

// 提交图像生成任务
func submitImageTask(config Config, prompt string) (string, error) {
	uri := "/api/generate/webui/text2img"
	sign, timestamp, nonce := generateSignature(uri, config.Liblibai.SecretKey)

	params := fmt.Sprintf("?AccessKey=%s&Signature=%s&Timestamp=%s&SignatureNonce=%s",
		config.Liblibai.AccessKey, sign, timestamp, nonce)

	url := config.Liblibai.APIURL + params

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

// 等待图像生成结果
func waitImageResult(config Config, uUID string) (string, error) {
	uri := "/api/generate/webui/status"

	for {
		sign, timestamp, nonce := generateSignature(uri, config.Liblibai.SecretKey)
		params := fmt.Sprintf("?AccessKey=%s&Signature=%s&Timestamp=%s&SignatureNonce=%s",
			config.Liblibai.AccessKey, sign, timestamp, nonce)

		url := config.Liblibai.QueryURL + params

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
								fmt.Println("✅ 图片生成成功并通过审核！")
								return imageURL, nil
							}
						} else {
							fmt.Println("⚠️ 图片未通过审核。")
						}
					}
				}
				if status == 4 || status == 5 {
					fmt.Println("❌ 生成失败或被拦截")
					return "", fmt.Errorf("generation failed or blocked")
				} else {
					fmt.Println("⏳ 图片生成中...状态: 生成/审核中")
				}
			}
		}

		time.Sleep(5 * time.Second)
	}
}

// 创建视频生成任务
func createVideoTask(config Config, prompt, imgURL string) (string, error) {
	url := config.Video.APIURL

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
	req.Header.Set("Authorization", "Bearer "+config.Qwen.APIKey)
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

// 轮询视频生成结果
func pollVideo(config Config, taskID string) (string, error) {
	url := fmt.Sprintf("%s/%s", config.Video.TaskURL, taskID)

	for i := 0; i < 60; i++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return "", err
		}

		req.Header.Set("Authorization", "Bearer "+config.Qwen.APIKey)

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
						fmt.Println("✅ 视频生成成功！")
						return videoURL, nil
					}
				} else if status == "FAILED" || status == "CANCELED" || status == "UNKNOWN" {
					fmt.Println("❌ 视频失败，状态：", status)
					return "", fmt.Errorf("video generation failed with status: %s", status)
				} else {
					fmt.Println("⏳ 视频生成中...", status)
				}
			}
		}

		time.Sleep(10 * time.Second)
	}

	return "", fmt.Errorf("video generation timeout")
}

func main() {
	// 配置信息
	config := Config{}
	config.Qwen.APIKey = "sk-e95d2534e97b4a969ce20cb8819ddbc6"
	config.Qwen.BaseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	config.Liblibai.AccessKey = "fcULyLSOwdrOmGpFEohhZg"
	config.Liblibai.SecretKey = "Q4tlxV4CnpCN5aFSKXwOMTO1PHFRp6rS"
	config.Liblibai.APIURL = "https://openapi.liblibai.cloud/api/generate/webui/text2img"
	config.Liblibai.QueryURL = "https://openapi.liblibai.cloud/api/generate/webui/status"
	config.Video.APIURL = "https://dashscope.aliyuncs.com/api/v1/services/aigc/video-generation/video-synthesis"
	config.Video.TaskURL = "https://dashscope.aliyuncs.com/api/v1/tasks"

	// 获取用户输入
	fmt.Println("请依次输入以下字段（直接回车使用默认值）：")
	pose := withDefault("姿势", "站立")
	location := withDefault("地点", "海边")
	timeOfDay := withDefault("时间", "傍晚")
	hairColor := withDefault("发色", "浅绿色")
	hairstyle := withDefault("发型", "偏分直短发")
	topWear := withDefault("上肢服装", "白衬衫")
	bottomWear := withDefault("臀部服装", "黑色小裙子")
	legWear := withDefault("腿部服装", "黑色丝袜")

	// 生成提示词
	originalPrompt := "This is a portrait of Ava standing at the beach at night. She has light green side parted short straight hair and wears grey crop top and black mini pants and black pantyhose."
	chineseInstruction := fmt.Sprintf(
		"请你根据以下中文修改以上描述中的对应内容（姿势、地点、时间、发色、发型、上肢服装、臀部服装、腿部服装），并只返回修改后的描述，依然返回纯英文，只是修改了对应的英文单词，不许输出任何汉字：%s、%s、%s、%s、%s、%s、%s、%s",
		pose, location, timeOfDay, hairColor, hairstyle, topWear, bottomWear, legWear,
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
		return
	}

	req, err := http.NewRequest("POST", config.Qwen.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("生成提示词失败: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Qwen.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("生成提示词失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("生成提示词失败: %v\n", err)
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("生成提示词失败: %v\n", err)
		return
	}

	choices := result["choices"].([]interface{})
	message := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	customPrompt := message["content"].(string)

	fmt.Printf("\n✅ 最终英文 Prompt:\n%s\n", customPrompt)

	// 元提示词
	metaPrompt := `This is a high-resolution everyday scene image with a natural style. hsg, yiiu, one lady, cosplay. 
	Ava is a captivating character with a blend of European mixed-race heritage, exuding the charm of a K-pop idol. 
	She has beautiful hair, showcasing her unique personality. Her eyes are large and expressive, almond-shaped, with blue color that is as clear and captivating as sapphires. 
	Her skin is extremely fair, porcelain-like, a typical feature of many K-pop idols. 
	She wears bold, dramatic makeup, including smokey eyes, winged eyeliner, defined brows, and a bold red lip, enhancing her exotic beauty. 
	Her visage is exquisitely petite, with a graceful triangular contour that accentuates her delicate features. 
	Ava's fashion style leans towards edgy and avant-garde, opting for bold, statement pieces with a touch of elegance, inspired by K-pop fashion trends. 
	Her lips are thick and sexy and seductive. Her hair is natural in color. iphone photo. 
	Her eyes are large and extremely beautiful and really seductive, with bold and glamorous eye makeup, long and voluminous false lashes, intense purple and shimmery gold eyeshadow, 
	sharp and upward-angled cat-eye liner, alluring and seductive eyes, upward-slanted fox-like eye corners, exuding a bewitching charm. 
	She has sexy extremely large breasts and extremely thick thighs and long legs. The whole picture is really seductive. 
	The photo is ultra realistic and ultra detailed with bright color and sharp contrast. She is looking at the camera.`

	// 合并提示词
	combinedPrompt := metaPrompt + customPrompt
	if len(combinedPrompt) > 2000 {
		fmt.Printf("\n❌ Prompt 过长（%d 字符），请缩短描述\n", len(combinedPrompt))
		return
	}

	// 提交图像生成任务
	imageID, err := submitImageTask(config, combinedPrompt)
	if err != nil {
		fmt.Printf("提交图像任务失败: %v\n", err)
		return
	}

	fmt.Printf("📌 图像任务已提交: %s\n", imageID)

	// 等待图像生成结果
	imageURL, err := waitImageResult(config, imageID)
	if err != nil {
		fmt.Printf("等待图像结果失败: %v\n", err)
		return
	}

	fmt.Printf("图片地址: %s\n", imageURL)

	// 创建视频任务
	videoPrompt := "美女轻轻摇摆着，表情魅惑而动人。她微微左右扭动着胯部，看起来很惬意。"
	videoID, err := createVideoTask(config, videoPrompt, imageURL)
	if err != nil {
		fmt.Printf("创建视频任务失败: %v\n", err)
		return
	}

	// 等待视频生成结果
	videoURL, err := pollVideo(config, videoID)
	if err != nil {
		fmt.Printf("等待视频结果失败: %v\n", err)
		return
	}

	fmt.Printf("🎬 视频链接：%s\n", videoURL)
}

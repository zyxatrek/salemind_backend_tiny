package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"salemind_backend_tiny/model"
	"salemind_backend_tiny/pkg/config"
	"salemind_backend_tiny/pkg/services"
)

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

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		return
	}

	// 创建服务
	genSvc := services.NewGenerationService(cfg)

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

	keyword := &model.ImagePromptKeyword{
		Pose:       pose,
		Location:   location,
		TimeOfDay:  timeOfDay,
		HairColor:  hairColor,
		Hairstyle:  hairstyle,
		TopWear:    topWear,
		BottomWear: bottomWear,
		LegWear:    legWear,
	}
	fmt.Printf("\n✅ 关键词:\n%v\n", keyword)
	// 生成提示词
	customPrompt, err := genSvc.GeneratePrompt(keyword)
	if err != nil {
		fmt.Printf("生成提示词失败: %v\n", err)
		return
	}
	fmt.Printf("\n✅ 最终英文 Prompt:\n%s\n", customPrompt)

	// 提交图像生成任务
	imageID, err := genSvc.SubmitImageTask(customPrompt)
	if err != nil {
		fmt.Printf("提交图像任务失败: %v\n", err)
		return
	}

	fmt.Printf("📌 图像任务已提交: %s\n", imageID)

	// 等待图像生成结果
	imageURL, err := genSvc.WaitImageResult(imageID)
	if err != nil {
		fmt.Printf("等待图像结果失败: %v\n", err)
		return
	}
	fmt.Printf("✅ 图片生成成功并通过审核！\n图片地址: %s\n", imageURL)

	// 创建视频任务
	// videoPrompt := "美女轻轻摇摆着，表情魅惑而动人。她微微左右扭动着胯部，看起来很惬意。"
	videoID, err := genSvc.CreateVideoTask(imageURL)
	if err != nil {
		fmt.Printf("创建视频任务失败: %v\n", err)
		return
	}
	fmt.Printf("📌 视频任务已提交: %s\n", videoID)

	// 等待视频生成结果
	videoURL, err := genSvc.PollVideo(videoID)
	if err != nil {
		fmt.Printf("等待视频结果失败: %v\n", err)
		return
	}
	fmt.Printf("✅ 视频生成成功！\n🎬 视频链接：%s\n", videoURL)
}

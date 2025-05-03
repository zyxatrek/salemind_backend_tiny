package model

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

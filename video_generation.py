from openai import OpenAI
import hmac
from hashlib import sha1
import base64
import time
import uuid
import requests
import json

# =================== 用户输入 ===================
def with_default(prompt, default):
    user_input = input(f"{prompt}（如：{default}）: ").strip()
    return user_input if user_input else default

print("请依次输入以下字段（直接回车使用默认值）：")
pose = with_default("姿势", "站立")
location = with_default("地点", "海边")
time_of_day = with_default("时间", "傍晚")
hair_color = with_default("发色", "浅绿色")
hairstyle = with_default("发型", "偏分直短发")
top_wear = with_default("上肢服装", "白衬衫")
bottom_wear = with_default("臀部服装", "黑色小裙子")
leg_wear = with_default("腿部服装", "黑色丝袜")

# =================== 调用 Qwen 生成 Custom Prompt ===================
original_prompt = (
    "This is a portrait of Ava standing at the beach at night. "
    "She has light green side parted short straight hair and wears grey crop top "
    "and black mini pants and black pantyhose."
)
chinese_instruction = (
    f"请你根据以下中文修改以上描述中的对应内容（姿势、地点、时间、发色、发型、上肢服装、臀部服装、腿部服装），"
    f"并只返回修改后的描述，依然返回纯英文，只是修改了对应的英文单词，不许输出任何汉字："
    f"{pose}、{location}、{time_of_day}、{hair_color}、{hairstyle}、{top_wear}、{bottom_wear}、{leg_wear}"
)

client = OpenAI(
    api_key="sk-e95d2534e97b4a969ce20cb8819ddbc6",
    base_url="https://dashscope.aliyuncs.com/compatible-mode/v1",
)
completion = client.chat.completions.create(
    model="qwen-plus",
    messages=[
        {"role": "system", "content": "You are a helpful assistant."},
        {"role": "user", "content": original_prompt + " " + chinese_instruction}
    ]
)
CUSTOM_PROMPT = completion.model_dump()["choices"][0]["message"]["content"]
print("\n✅ 最终英文 Prompt:\n", CUSTOM_PROMPT)

# =================== 图像生成参数 ===================
ACCESS_KEY = "fcULyLSOwdrOmGpFEohhZg"
SECRET_KEY = "Q4tlxV4CnpCN5aFSKXwOMTO1PHFRp6rS"
API_URL_SUBMIT = "https://openapi.liblibai.cloud/api/generate/webui/text2img"
API_URL_QUERY = "https://openapi.liblibai.cloud/api/generate/webui/status"

META_PROMPT = (
    "This is a high-resolution everyday scene image with a natural style. hsg, yiiu, one lady, cosplay. "
    "Ava is a captivating character with a blend of European mixed-race heritage, exuding the charm of a K-pop idol. "
    "She has beautiful hair, showcasing her unique personality. Her eyes are large and expressive, almond-shaped, with blue color that is as clear and captivating as sapphires. "
    "Her skin is extremely fair, porcelain-like, a typical feature of many K-pop idols. "
    "She wears bold, dramatic makeup, including smokey eyes, winged eyeliner, defined brows, and a bold red lip, enhancing her exotic beauty. "
    "Her visage is exquisitely petite, with a graceful triangular contour that accentuates her delicate features. "
    "Ava's fashion style leans towards edgy and avant-garde, opting for bold, statement pieces with a touch of elegance, inspired by K-pop fashion trends. "
    "Her lips are thick and sexy and seductive. Her hair is natural in color. iphone photo. "
    "Her eyes are large and extremely beautiful and really seductive, with bold and glamorous eye makeup, long and voluminous false lashes, intense purple and shimmery gold eyeshadow, "
    "sharp and upward-angled cat-eye liner, alluring and seductive eyes, upward-slanted fox-like eye corners, exuding a bewitching charm. "
    "She has sexy extremely large breasts and extremely thick thighs and long legs. The whole picture is really seductive. "
    "The photo is ultra realistic and ultra detailed with bright color and sharp contrast. She is looking at the camera. "
)

combined_prompt = META_PROMPT + CUSTOM_PROMPT
if len(combined_prompt) > 2000:
    print(f"\n❌ Prompt 过长（{len(combined_prompt)} 字符），请缩短描述")
    exit(1)

REQUEST_PARAMS = {
    "templateUuid": "6f7c4652458d4802969f8d089cf5b91f",
    "generateParams": {
        "checkPointId": "ddb974b0717e4f96bc7789a068004831",
        "vaeId": "",
        "prompt": combined_prompt,
        "clipSkip": 2,
        "steps": 20,
        "width": 1000,
        "height": 1600,
        "imgCount": 1,
        "seed": -1,
        "restoreFaces": 0,
        "additionalNetwork": [
            {"modelId": "4e6783ca937047bcbba4874cedd1f552", "weight": 0.4},
            {"modelId": "1aaebe53f31e489993ee5c4338e450e4", "weight": 0.5},
            {"modelId": "da6b0f7fb7004bddb7bc8d3fe698641f", "weight": 0.45}
        ]
    }
}

def generate_signature(uri, secret_key):
    timestamp = str(int(time.time() * 1000))
    nonce = str(uuid.uuid4())
    content = f"{uri}&{timestamp}&{nonce}"
    digest = hmac.new(secret_key.encode(), content.encode(), sha1).digest()
    sign = base64.urlsafe_b64encode(digest).rstrip(b'=').decode()
    return {"signature": sign, "timestamp": timestamp, "signature_nonce": nonce}

def submit_image_task():
    uri = "/api/generate/webui/text2img"
    sign = generate_signature(uri, SECRET_KEY)
    headers = {"Content-Type": "application/json"}
    params = f"?AccessKey={ACCESS_KEY}&Signature={sign['signature']}&Timestamp={sign['timestamp']}&SignatureNonce={sign['signature_nonce']}"
    url = API_URL_SUBMIT + params
    resp = requests.post(url, headers=headers, data=json.dumps(REQUEST_PARAMS))
    if resp.status_code == 200:
        data = resp.json()
        if data.get("code") == 0:
            return data["data"]["generateUuid"]
    return None

def wait_image_result(uuid_):
    uri = "/api/generate/webui/status"
    while True:
        sign = generate_signature(uri, SECRET_KEY)
        params = f"?AccessKey={ACCESS_KEY}&Signature={sign['signature']}&Timestamp={sign['timestamp']}&SignatureNonce={sign['signature_nonce']}"
        url = API_URL_QUERY + params
        headers = {"Content-Type": "application/json"}
        resp = requests.post(url, headers=headers, json={"generateUuid": uuid_})
        if resp.status_code == 200:
            result = resp.json()
            if result.get("code") == 0:
                status = result["data"]["generateStatus"]
                images = result["data"].get("images", [])
                if images and images[0].get("auditStatus") == 3:
                    print("✅ 图片生成成功并通过审核！")
                    print("图片地址:", images[0]["imageUrl"])
                    return images[0]["imageUrl"]
                elif images:
                    print("⚠️ 图片未通过审核。")
                    return None
                elif status in [4, 5]:
                    print("❌ 生成失败或被拦截")
                    return None
                else:
                    print(f"⏳ 图片生成中...状态: 生成/审核中")
        time.sleep(5)

# =================== 图转视频 ===================
def create_video_task(prompt, img_url):
    url = "https://dashscope.aliyuncs.com/api/v1/services/aigc/video-generation/video-synthesis"
    headers = {
        "Content-Type": "application/json",
        "Authorization": "Bearer sk-e95d2534e97b4a969ce20cb8819ddbc6",
        "X-DashScope-Async": "enable"
    }
    payload = {
        "model": "wanx2.1-i2v-plus",
        "input": {"prompt": prompt, "img_url": img_url},
        "parameters": {"resolution": "720P", "duration": 5, "prompt_extend": True}
    }
    r = requests.post(url, headers=headers, json=payload)
    if r.status_code == 200:
        return r.json()["output"]["task_id"]
    else:
        print("❌ 视频任务提交失败：", r.text)
        return None

def poll_video(task_id):
    url = f"https://dashscope.aliyuncs.com/api/v1/tasks/{task_id}"
    headers = {"Authorization": "Bearer sk-e95d2534e97b4a969ce20cb8819ddbc6"}
    for _ in range(60):
        r = requests.get(url, headers=headers)
        if r.status_code == 200:
            out = r.json().get("output", {})
            status = out.get("task_status")
            if status == "SUCCEEDED":
                print("✅ 视频生成成功！")
                print("🎬 视频链接：", out["video_url"])
                return
            elif status in ["FAILED", "CANCELED", "UNKNOWN"]:
                print("❌ 视频失败，状态：", status)
                return
            else:
                print("⏳ 视频生成中...", status)
        time.sleep(10)
    print("⚠️ 视频超时未完成")

# =================== 主流程 ===================
if __name__ == "__main__":
    image_id = submit_image_task()
    if image_id:
        print("📌 图像任务已提交:", image_id)
        image_url = wait_image_result(image_id)
        if image_url:
            video_prompt = "美女轻轻摇摆着，表情魅惑而动人。她微微左右扭动着胯部，看起来很惬意。"
            video_id = create_video_task(video_prompt, image_url)
            if video_id:
                poll_video(video_id)
    else:
        print("❌ 图像任务提交失败")

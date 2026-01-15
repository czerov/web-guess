package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

// HTML 模板作为字符串常量，包含完整的页面结构和样式
const htmlTemplate = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>猜数字游戏</title>
    <style>
        body {
            font-family: 'Courier New', monospace;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            text-align: center;
            background-color: #0a0a1e;
            background-image: linear-gradient(to bottom, #0a0a1e, #1a1a3e);
            color: #00ffff;
            min-height: 100vh;
        }
        h1 {
            font-size: 48px;
            font-weight: bold;
            color: #ff00ff;
            text-shadow: 0 0 10px #ff00ff, 0 0 20px #ff00ff, 0 0 30px #ff00ff;
            margin-bottom: 30px;
        }
        .game-container {
            background-color: rgba(10, 10, 30, 0.8);
            padding: 40px;
            border-radius: 10px;
            box-shadow: 0 0 20px rgba(0, 255, 255, 0.3);
            border: 1px solid rgba(0, 255, 255, 0.5);
            backdrop-filter: blur(5px);
        }
        p {
            font-size: 18px;
            color: #00ffff;
            margin-bottom: 20px;
        }
        input[type="number"] {
            padding: 15px;
            font-size: 20px;
            width: 250px;
            margin: 20px 0;
            background-color: rgba(10, 10, 30, 0.8);
            color: #00ffff;
            border: 2px solid #00ffff;
            border-radius: 5px;
            box-shadow: 0 0 10px rgba(0, 255, 255, 0.5);
            font-family: 'Courier New', monospace;
        }
        input[type="number"]:focus {
            outline: none;
            box-shadow: 0 0 20px rgba(0, 255, 255, 0.8);
            border-color: #ff00ff;
        }
        button {
            padding: 15px 30px;
            font-size: 20px;
            background-color: rgba(255, 0, 255, 0.2);
            color: #ff00ff;
            border: 2px solid #ff00ff;
            border-radius: 5px;
            cursor: pointer;
            font-family: 'Courier New', monospace;
            font-weight: bold;
            box-shadow: 0 0 10px rgba(255, 0, 255, 0.5);
            transition: all 0.3s ease;
        }
        button:hover {
            background-color: rgba(255, 0, 255, 0.4);
            box-shadow: 0 0 20px rgba(255, 0, 255, 0.8);
            transform: translateY(-2px);
        }
        .message {
            margin: 30px 0;
            padding: 25px;
            border-radius: 8px;
            font-weight: bold;
            border: 2px solid;
            box-shadow: 0 0 15px;
        }
        .message.success {
            background-color: rgba(0, 255, 0, 0.2);
            color: #00ff00;
            border-color: #00ff00;
            box-shadow: 0 0 20px rgba(0, 255, 0, 0.7);
            font-size: 32px;
            text-shadow: 0 0 10px #00ff00;
        }
        .message.error {
            background-color: rgba(255, 255, 0, 0.2);
            color: #ffff00;
            border-color: #ffff00;
            box-shadow: 0 0 20px rgba(255, 255, 0, 0.7);
            font-size: 24px;
            text-shadow: 0 0 10px #ffff00;
        }
        .message.info {
            background-color: rgba(0, 255, 255, 0.2);
            color: #00ffff;
            border-color: #00ffff;
            box-shadow: 0 0 15px rgba(0, 255, 255, 0.5);
        }
    </style>
</head>
<body>
    <div class="game-container">
        <h1>猜数字游戏</h1>
        %s
    </div>
</body>
</html>`

// 主页面内容，包含猜测表单
const homePage = `
        <p>我已经想好了一个 1-100 之间的数字，快来猜猜看吧！</p>
        <form method="POST" action="/guess">
            <input type="number" name="number" min="1" max="100" placeholder="请输入猜测的数字" required>
            <br>
            <button type="submit">提交猜测</button>
        </form>
`

// 结果页面内容，包含反馈信息和再玩一次按钮
const resultPage = `
        <div class="message %s">%s</div>
        <form method="GET" action="/">
            <button type="submit">再玩一次</button>
        </form>
`

func main() {
	// 初始化随机数种子，确保每次运行生成不同的随机数
	rand.Seed(time.Now().UnixNano())

	// 注册 HTTP 处理函数
	// HandleFunc 函数接受两个参数：
	// 1. 路径模式（如 "/" 或 "/guess"）
	// 2. 处理函数（func(http.ResponseWriter, *http.Request)）
	// 当用户访问对应路径时，会调用相应的处理函数
	http.HandleFunc("/", homeHandler)       // 处理首页请求
	http.HandleFunc("/guess", guessHandler) // 处理猜测提交请求

	// 从环境变量读取端口，默认为 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 启动 HTTP 服务器，监听指定端口
	// ListenAndServe 函数会阻塞运行，直到服务器关闭
	fmt.Printf("服务器启动，监听端口 %s...\n", port)
	fmt.Printf("请访问 http://localhost:%s 开始游戏\n", port)
	http.ListenAndServe(":"+port, nil)
}

// homeHandler 处理首页请求 (GET /)
// 参数说明：
// - w: http.ResponseWriter 用于向客户端写入响应
// - r: *http.Request 包含客户端请求的所有信息
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// 生成 1-100 的随机数
	target := rand.Intn(100) + 1

	// 创建一个新的 Cookie，存储目标数字
	// Cookie 是存储在客户端的小型文本文件，用于在请求之间保持状态
	cookie := &http.Cookie{
		Name:     "target",             // Cookie 名称
		Value:    strconv.Itoa(target), // 将数字转换为字符串
		Path:     "/",                  // Cookie 作用路径
		MaxAge:   3600,                 // Cookie 过期时间（秒）
		HttpOnly: true,                 // 防止 JavaScript 访问，提高安全性
	}

	// 将 Cookie 添加到响应中，这样客户端会保存这个 Cookie
	http.SetCookie(w, cookie)

	// 生成完整的 HTML 页面
	// 使用 htmlTemplate 作为模板，将 homePage 插入到模板中
	html := fmt.Sprintf(htmlTemplate, homePage)

	// 设置响应的 Content-Type 为 text/html，告诉浏览器这是 HTML 内容
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// 将 HTML 内容写入响应
	fmt.Fprint(w, html)
}

// guessHandler 处理猜测提交请求 (POST /guess)
// 参数说明：
// - w: http.ResponseWriter 用于向客户端写入响应
// - r: *http.Request 包含客户端请求的所有信息
func guessHandler(w http.ResponseWriter, r *http.Request) {
	// 检查请求方法是否为 POST
	// 因为这个处理函数只应该处理 POST 请求
	if r.Method != http.MethodPost {
		// 如果不是 POST 请求，返回 405 Method Not Allowed 错误
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}

	// 解析表单数据
	// 对于 POST 请求，需要调用 ParseForm() 来获取表单数据
	if err := r.ParseForm(); err != nil {
		// 如果解析失败，返回 400 Bad Request 错误
		http.Error(w, "无法解析表单数据", http.StatusBadRequest)
		return
	}

	// 从表单中获取用户输入的数字
	// FormValue 方法会自动处理表单解析，并返回指定字段的值
	guessStr := r.FormValue("number")

	// 将字符串转换为整数
	guess, err := strconv.Atoi(guessStr)
	if err != nil {
		// 如果转换失败，返回 400 Bad Request 错误
		http.Error(w, "请输入有效的数字", http.StatusBadRequest)
		return
	}

	// 从 Cookie 中获取目标数字
	// Cookie 存储在客户端，每次请求都会发送到服务器
	cookie, err := r.Cookie("target")
	if err != nil {
		// 如果没有找到 Cookie，可能是用户直接访问了 /guess 路径
		// 重定向到首页，让用户重新开始游戏
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// 将 Cookie 中的值转换为整数
	target, err := strconv.Atoi(cookie.Value)
	if err != nil {
		// 如果转换失败，重定向到首页
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// 比较用户猜测的数字和目标数字
	var message, messageClass string
	if guess > target {
		message = "太大了！"
		messageClass = "error"
	} else if guess < target {
		message = "太小了！"
		messageClass = "error"
	} else {
		message = "恭喜猜对！"
		messageClass = "success"
	}

	// 生成结果页面
	// 使用 htmlTemplate 作为模板，将 resultPage 插入到模板中
	// resultPage 需要两个参数：消息类型和消息内容
	resultContent := fmt.Sprintf(resultPage, messageClass, message)
	html := fmt.Sprintf(htmlTemplate, resultContent)

	// 设置响应的 Content-Type 为 text/html
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// 将 HTML 内容写入响应
	fmt.Fprint(w, html)
}

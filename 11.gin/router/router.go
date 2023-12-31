package router

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"zenkie.cn/learn-gin/controllers"
)

func Router() *gin.Engine {
	router := gin.Default()

	// collectLog(router)
	setupPing(router)
	setupAsciiJSON(router)
	renderHTML(router)
	go serverPush()
	JSONP(router)
	formBinding(router)
	formHandle(router)
	pureJSON(router)
	queryAndPostform(router)
	secureJSON(router)
	renderDatas(router)
	useRouteParam(router)
	useGroup(router)
	useBasicAuth(router)
	useStatic(router)
	useCustomMiddleware(router)

	// examples
	router.GET("/user-info", controllers.GetUserInfo)
	router.GET("/user-list", controllers.GetList)
	return router
}

func setupPing(router *gin.Engine) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}

func setupAsciiJSON(router *gin.Engine) {
	router.GET("/someJSON", func(c *gin.Context) {
		data := map[string]interface{}{
			"lang": "GO 语言",
			"tag":  "<br>",
		}

		// Output {"lang":"GO \u8bed\u8a00","tag":"\u003cbr\u003e"}
		c.AsciiJSON(http.StatusOK, data)
	})
}

func renderHTML(router *gin.Engine) {
	// router.LoadHTMLFiles("templates/me.html")
	router.LoadHTMLGlob("templates/**/*")
	router.GET("/index", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "me.html", gin.H{
			"title": "Bonjour! Je suis Zenkie Bear.",
		})
	})
	router.GET("/posts/index", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "posts/index.html", gin.H{
			"title": "Posts",
		})
	})
	router.GET("/users/index", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "users/index.html", gin.H{
			"title": "Users",
		})
	})
}

var html = template.Must(template.New("https").Parse(`
<html>
<head>
  <title>Https Test</title>
	<script src="/assets/app.js"></script>
</head>
<body>
  <h1 style="color: red;">Bienvenue, Ginner!</h1>
</html>
`))

func serverPush() {
	router := gin.New()
	router.Static("/assets", "./assets")
	router.SetHTMLTemplate(html)

	router.GET("/", func(ctx *gin.Context) {
		if pusher := ctx.Writer.Pusher(); pusher != nil {
			if err := pusher.Push("/assets/app.js", nil); err != nil {
				log.Printf("Failed to push: %v", err)
			}
		}
		ctx.HTML(http.StatusOK, "https", gin.H{
			"status": "success",
		})
	})
	router.RunTLS("localhost:8081", "ssl.pem", "ssl.key")
}

func JSONP(router *gin.Engine) {
	router.GET("/jsonp", func(c *gin.Context) {
		data := map[string]interface{}{
			"foo": "bar",
		}

		// Request: /jsonp?callback=x
		// Output: x({"foo": "bar"})
		c.JSONP(http.StatusOK, data)
	})
	router.RunTLS("localhost:8081", "ssl.pem", "ssl.key")
}

type LoginForm struct {
	User     string `form:"user" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func formBinding(router *gin.Engine) {
	// curl -v --form user=john --form password=doe http://localhost:8080/login
	router.POST("/login", func(ctx *gin.Context) {
		var form LoginForm
		// Explicit binding
		// It could be XML, JSON or other formats
		// ctx.ShouldBindWith(&form, binding.Form)
		if ctx.ShouldBind(&form) == nil {
			if form.User == "john" && form.Password == "doe" {
				ctx.JSON(http.StatusOK, gin.H{"status": "You are logged in!"})
			} else {
				ctx.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			}
		}
		log.Print(form)
	})
}

func formHandle(router *gin.Engine) {
	router.POST("/form_post", func(ctx *gin.Context) {
		message := ctx.PostForm("message")
		nick := ctx.DefaultPostForm("nick", "anounymous")

		ctx.JSON(http.StatusOK, gin.H{
			"status":  "posted",
			"message": message,
			"nick":    nick,
		})
	})
	// curl -v --form message=hello http://localhost:8080/form_post
	// {"message":"hello","nick":"anounymous","status":"posted"}
}

func pureJSON(router *gin.Engine) {
	json := gin.H{
		"html": "<b>Hello, world!</b>",
	}
	router.GET("/json", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, json)
	})
	// Output: {"html":"\u003cb\u003eHello, world!\u003c/b\u003e"}

	router.GET("/purejson", func(ctx *gin.Context) {
		ctx.PureJSON(http.StatusOK, json)
	})
	// Output: {"html":"<b>Hello, world!</b>"}
}

func queryAndPostform(router *gin.Engine) {
	router.POST("/post", func(ctx *gin.Context) {
		id := ctx.Query("id")
		page := ctx.DefaultQuery("page", "0")
		name := ctx.PostForm("name")
		message := ctx.PostForm("message")

		fmt.Printf("id: %s, page: %s, name: %s, message: %s", id, page, name, message)
	})
}

func secureJSON(router *gin.Engine) {
	router.GET("/secureJSON", func(ctx *gin.Context) {
		names := []string{"lena", "austin", "foo"}

		ctx.SecureJSON(http.StatusOK, names)
	})
}

func renderDatas(router *gin.Engine) {
	result := gin.H{"message": "Hey", "status": http.StatusOK}
	router.GET("/severalJSON", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, result)
	})

	router.GET("/moreJSON", func(ctx *gin.Context) {
		// You can use a struct
		var msg struct {
			Name    string `json:"user"`
			Message string
			Number  int
		}
		msg.Name = "Milo"
		msg.Message = "Hey"
		msg.Number = 1990

		// Note: the msg.name is "user" in json
		// Output: {"user":"Milo","Message":"Hey","Number":1990}
		ctx.JSON(http.StatusOK, msg)
	})

	router.GET("/severalXML", func(ctx *gin.Context) {
		ctx.XML(http.StatusOK, result)
	})

	router.GET("/severalYAML", func(ctx *gin.Context) {
		ctx.YAML(http.StatusOK, result)
	})
}

func useRouteParam(router *gin.Engine) {
	router.GET("/user/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")
		ctx.String(http.StatusOK, "Bonjour %s！", name)
	})
}

func useGroup(router *gin.Engine) {
	v1 := router.Group("v1")
	{
		loginHandler := func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Logged in,",
				"version": "v1",
			})
		}
		registHandler := func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Registered successful.",
				"version": "v1",
			})
		}
		v1.GET("/login", loginHandler)
		v1.GET("/register", registHandler)
	}
	v2 := router.Group("v2")
	{
		loginHandler := func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Logged in,",
				"version": "v2",
			})
		}
		registHandler := func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Registered successful.",
				"version": "v2",
			})
		}
		v2.GET("/login", loginHandler)
		v2.GET("/register", registHandler)
	}
}

var secrets = gin.H{
	"foo":    gin.H{"email": "foo@bar.com", "phone": "123433"},
	"austin": gin.H{"email": "austin@example.com", "phone": "666"},
	"lena":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
}

func useBasicAuth(router *gin.Engine) {
	authorized := router.Group("/admin", gin.BasicAuth(gin.Accounts{
		"foo":    "bar",
		"austin": "1234",
		"lena":   "hello2",
		"manu":   "4321",
	}))

	authorized.GET("/secrets", func(ctx *gin.Context) {
		user := ctx.MustGet(gin.AuthUserKey).(string)
		if secret, exists := secrets[user]; exists {
			ctx.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{"user": user, "secret": "Unauthorized :("})
		}
	})
}

func collectLog(router *gin.Engine) {
	gin.DisableConsoleColor()

	// record to file
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	// also record to console
	// gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	// Define logging format
	// This will enable console logging
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("GIN %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
}

func useStatic(router *gin.Engine) {
	router.Static("/assets", "./assets")
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		c.Set("example", "12345")

		// Before request

		c.Next()

		// After request
		latency := time.Since(t)
		status := c.Writer.Status()
		log.Printf("Spend: %v\tStatus: %v", latency, status)
	}
}

func useCustomMiddleware(router *gin.Engine) {
	router.Use(Logger())

	// custom logger
	router.Use(gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		// my format
		return fmt.Sprintf("[Mon gin log] %s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			params.ClientIP,
			params.TimeStamp.Format(time.RFC1123),
			params.Method,
			params.Path,
			params.Request.Proto,
			params.StatusCode,
			params.Latency,
			params.Request.UserAgent(),
			params.ErrorMessage,
		)

	}))

	router.GET("/test", func(ctx *gin.Context) {
		log.Println(ctx.MustGet("example").(string))
	})
}

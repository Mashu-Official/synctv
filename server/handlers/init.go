package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/synctv-org/synctv/public"
	"github.com/synctv-org/synctv/server/middlewares"
	"github.com/synctv-org/synctv/utils"
)

func Init(e *gin.Engine) {
	{
		web := e.Group("/web")

		web.Use(func(ctx *gin.Context) {
			if ctx.Request.URL.Path == "/web/" {
				ctx.Header("Cache-Control", "no-store")
			} else {
				ctx.Header("Cache-Control", "public, max-age=31536000")
			}
			ctx.Next()
		})

		web.StaticFS("", http.FS(public.Public))
	}

	{
		api := e.Group("/api")

		needAuthUserApi := api.Group("")
		needAuthUserApi.Use(middlewares.AuthUserMiddleware)

		needAuthRoomApi := api.Group("")
		needAuthRoomApi.Use(middlewares.AuthRoomMiddleware)

		{
			public := api.Group("/public")

			public.GET("/settings", Settings)
		}

		{
			// TODO: admin api implement
			// admin := api.Group("/admin")
		}

		{
			room := api.Group("/room")
			needAuthRoom := needAuthRoomApi.Group("/room")
			needAuthUser := needAuthUserApi.Group("/room")

			room.GET("/ws", NewWebSocketHandler(utils.NewWebSocketServer()))

			room.GET("/check", CheckRoom)

			room.GET("/list", RoomList)

			needAuthUser.POST("/create", CreateRoom)

			needAuthUser.POST("/login", LoginRoom)

			needAuthRoom.POST("/delete", DeleteRoom)

			needAuthRoom.POST("/pwd", SetRoomPassword)

			needAuthRoom.GET("/setting", RoomSetting)
		}

		{
			movie := api.Group("/movie")
			needAuthMovie := needAuthRoomApi.Group("/movie")

			needAuthMovie.GET("/list", MovieList)

			needAuthMovie.GET("/current", CurrentMovie)

			needAuthMovie.GET("/movies", Movies)

			needAuthMovie.POST("/current", ChangeCurrentMovie)

			needAuthMovie.POST("/push", PushMovie)

			needAuthMovie.POST("/edit", EditMovie)

			needAuthMovie.POST("/swap", SwapMovie)

			needAuthMovie.POST("/delete", DelMovie)

			needAuthMovie.POST("/clear", ClearMovies)

			movie.HEAD("/proxy/:roomId/:pullKey", ProxyMovie)

			movie.GET("/proxy/:roomId/:pullKey", ProxyMovie)

			{
				live := needAuthMovie.Group("/live")

				live.POST("/publishKey", NewPublishKey)

				live.GET("/*pullKey", JoinLive)
			}
		}

		{
			user := api.Group("/user")
			needAuthUser := needAuthUserApi.Group("/user")

			user.POST("/login", LoginUser)

			user.POST("/signup", SignupUser)

			needAuthUser.POST("/logout", LogoutUser)

			needAuthUser.GET("/me", Me)

			needAuthUser.POST("/pwd", SetUserPassword)
		}
	}

	e.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/web/")
	})
}

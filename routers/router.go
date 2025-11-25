package routers

import (
	"bunker-web/configs"
	"bunker-web/middlewares"
	"net/http"

	"github.com/gin-contrib/gzip"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "bunker-web/docs"

	"github.com/gin-gonic/gin"
)

func RegisterPhoenixAPI(router *gin.Engine) {
	apiGroup := router.Group("/api")
	{
		// Session check
		apiGroup.Use(middlewares.BearerHandler())

		// APIs related to phoenix login
		phoenixLoginGroup := apiGroup.Group("/phoenix")
		{
			// Phoenix login
			phoenixLoginGroup.POST("/login", routers.API.Phoenix.Login)
		}

		// Login check
		apiGroup.Use(middlewares.LoginHandler())

		// Phoenix api for accesses from PhoenixBuilder client only
		phoenixGroup := apiGroup.Group("/phoenix")
		phoenixGroup.Use(middlewares.NormalPermissionHandler())
		{
			// Start type
			phoenixGroup.GET("/transfer_start_type", routers.API.Phoenix.TransferStartType)
			// MCP check num
			phoenixGroup.POST("/transfer_check_num", routers.API.Phoenix.TransferCheckNum)
		}
	}
}

func RegisterOpenAPI(router *gin.Engine) {
	// APIs which not require bearer token
	openApiGroup := router.Group("/openapi")
	{
		// Welcome
		openApiGroup.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "Welcome to BunkerWeb OpenAPI!")
		})

		// User
		userGroupWithoutAuth := openApiGroup.Group("/user")
		{
			// Redeem
			userGroupWithoutAuth.POST("/redeem", routers.OpenAPI.User.Redeem)
		}

		// Swagger
		openApiGroup.GET("/swagger/*any", ginSwagger.WrapHandler(
			swaggerfiles.Handler,
			ginSwagger.DefaultModelsExpandDepth(-1),
		))

		// Session check
		openApiGroup.Use(middlewares.OpenAPIKeyHandler())

		// User
		userGroupWithAuth := openApiGroup.Group("/user")
		// User
		{
			// Query info
			userGroupWithAuth.GET("/get_status", routers.OpenAPI.User.GetStatus)
		}

		// Helper
		helperGroup := openApiGroup.Group("/helper")
		{
			// Get user details
			helperGroup.POST("/get_user_details", routers.OpenAPI.Helper.GetUserDetails)
			// Get user state
			helperGroup.POST("/get_user_state", routers.OpenAPI.Helper.GetUserState)
			// Search user
			helperGroup.POST("/search_user", routers.OpenAPI.Helper.SearchUser)
			// Rental server
			rentalServerGroup := helperGroup.Group("/rental_server")
			{
				// Search rental server
				rentalServerGroup.POST("/search", routers.OpenAPI.Helper.RentalServer.Search)
				// Get rental server player list
				rentalServerGroup.POST("/get_player_list", routers.OpenAPI.Helper.RentalServer.GetPlayerList)
			}
		}

		// Owner
		ownerGroup := openApiGroup.Group("/owner")
		{
			// Rental server
			rentalServerGroup := ownerGroup.Group("/rental_server")
			{
				// Ban rental server player
				rentalServerGroup.POST("/ban_player", routers.OpenAPI.Owner.RentalServer.BanPlayer)
				// Create backup
				rentalServerGroup.POST("/create_backup", routers.OpenAPI.Owner.RentalServer.CreateBackup)
				// Kick rental server player
				rentalServerGroup.POST("/kick_player", routers.OpenAPI.Owner.RentalServer.KickPlayer)
				// Set level limitation
				rentalServerGroup.POST("/set_level_limitation", routers.OpenAPI.Owner.RentalServer.SetLevelLimitation)
				// Set password
				rentalServerGroup.POST("/set_password", routers.OpenAPI.Owner.RentalServer.SetPassword)
				// Set visibility
				rentalServerGroup.POST("/set_visibility", routers.OpenAPI.Owner.RentalServer.SetVisibility)
				// Turn off rental server
				rentalServerGroup.POST("/turn_off", routers.OpenAPI.Owner.RentalServer.TurnOff)
				// Turn on rental server
				rentalServerGroup.POST("/turn_on", routers.OpenAPI.Owner.RentalServer.TurnOn)
			}
		}
	}
}

func RegisterWebAPI(router *gin.Engine) {
	apiGroup := router.Group("/api")
	{
		// Webauthn Login APIs
		webauthnLoginGroup := apiGroup.Group("/webauthn/login")
		{
			// Get login options
			webauthnLoginGroup.GET("/options", routers.API.Webauthn.Login.Options)
			// Verify login
			webauthnLoginGroup.POST("/verification", routers.API.Webauthn.Login.Verification)
		}

		// APIs related to user
		userLoginGroup := apiGroup.Group("/user")
		{
			// Register a new user, captcha required
			userLoginGroup.POST("/register", routers.API.User.Register)
			// Normal login, token unsupported
			userLoginGroup.POST("/login", routers.API.User.Login)
			// Reset password
			userLoginGroup.POST("/reset_password", routers.API.User.ResetPassword)
			// Email group
			userLoginEmailGroup := userLoginGroup.Group("/email")
			{
				// Send email code
				userLoginEmailGroup.POST("/send_code", routers.API.User.Email.SendCode)
			}
		}

		// Session check
		apiGroup.Use(middlewares.BearerHandler())
		// Login check
		apiGroup.Use(middlewares.LoginHandler())

		// Notice APIs
		noticeGroup := apiGroup.Group("/notice")
		{
			// Get notice list
			noticeGroup.POST("/query", routers.API.Notice.Query)
			// Admin permission check
			noticeGroup.Use(middlewares.AdminPermissionHandler())
			{
				// Create a new notice
				noticeGroup.POST("/create", routers.API.Notice.Create)
				// Edit an notice
				noticeGroup.POST("/edit", routers.API.Notice.Edit)
				// Delete an notice
				noticeGroup.POST("/delete", routers.API.Notice.Delete)
			}
		}

		// Webauthn APIs
		webauthnGroup := apiGroup.Group("/webauthn")
		{
			// Remove by id
			webauthnGroup.POST("/remove", routers.API.Webauthn.Remove)
			// Register group
			webauthnRegisterGroup := webauthnGroup.Group("/register")
			{
				// Get registeration options
				webauthnRegisterGroup.GET("/options", routers.API.Webauthn.Register.Options)
				// Verify registeration
				webauthnRegisterGroup.POST("/verification", routers.API.Webauthn.Register.Verification)
			}
		}

		// APIs related to user
		userGroup := apiGroup.Group("/user")
		{
			// Change password
			userGroup.POST("/change_password", routers.API.User.ChangePassword)
			// Get user info
			userGroup.GET("/get_status", routers.API.User.GetStatus)
			// Logout
			userGroup.GET("/logout", routers.API.User.Logout)
			// Remove account
			userGroup.POST("/remove_account", routers.API.User.RemoveAccount)
			// Set response to
			userGroup.POST("/set_response_to", routers.API.User.SetResponseTo)
			// Use redeem code
			userGroup.POST("/redeem", routers.API.User.Redeem)
			// API key group
			apiKeyGroup := userGroup.Group("/api_key")
			{
				// Generate API key
				apiKeyGroup.GET("/generate", routers.API.User.APIKey.Generate)
				// Disable API key
				apiKeyGroup.GET("/disable", routers.API.User.APIKey.Disable)
			}
			// Email group
			userEmailGroup := userGroup.Group("/email")
			{
				// Email bind
				userEmailGroup.POST("/bind", routers.API.User.Email.Bind)
				// Email unbind
				userEmailGroup.POST("/unbind", routers.API.User.Email.Unbind)
			}
			// Normal permission check
			userGroup.Use(middlewares.NormalPermissionHandler())
			{
				// Bind game id
				userGroup.POST("/bind_game_id", routers.API.User.BindGameId)

				// Get fbtoken file
				userGroup.POST("/get_phoenix_token", routers.API.User.GetPhoenixToken)
			}
		}

		// APIs related to slot
		slotGroup := apiGroup.Group("/slot")
		{
			// Delete
			slotGroup.POST("/delete", routers.API.Slot.Delete)
		}
		// Normal permission check
		slotGroup.Use(middlewares.NormalPermissionHandler())
		{
			// Set game id
			slotGroup.POST("/set_game_id", routers.API.Slot.SetGameID)
			// Extend expire time
			slotGroup.POST("/extend_expire_time", routers.API.Slot.ExtendExpireTime)
		}

		// Helper is a user which used to login to server
		helperGroup := apiGroup.Group("/helper")
		{
			// Get helper info
			helperGroup.GET("/get_status", routers.API.Helper.GetStatus)
			// Delete helper info
			helperGroup.GET("/unbind", routers.API.Helper.UnBind)
			// Normal permission check
			helperGroup.Use(middlewares.NormalPermissionHandler())
			{
				// Change helper name
				helperGroup.POST("/change_name", routers.API.Helper.ChangeName)
				// Bind account group
				bindAccountGroup := helperGroup.Group("/bind_account")
				{
					// Use email account to create helper
					bindAccountGroup.POST("/email", routers.API.Helper.Email)
					// Create a helper by guest login
					bindAccountGroup.GET("/guest", routers.API.Helper.Guest)
					// Use mobile account to create helper
					bindAccountGroup.POST("/mobile", routers.API.Helper.Mobile)
					// Send smscode
					bindAccountGroup.POST("/send_sms", routers.API.Helper.SendSMS)
				}
			}
		}

		// Owner is a game account which can use to manage server
		ownerGroup := apiGroup.Group("/owner")
		{
			// Get helper info
			ownerGroup.GET("/get_status", routers.API.Owner.GetStatus)
			// Delete helper info
			ownerGroup.GET("/unbind", routers.API.Owner.UnBind)
			// Normal permission check
			ownerGroup.Use(middlewares.NormalPermissionHandler())
			{
				// Get mail reward
				ownerGroup.GET("/get_mail_reward", routers.API.Owner.GetMailReward)
				// Use gift code
				ownerGroup.POST("/use_gift_code", routers.API.Owner.UseGiftCode)
				// Bind account group
				bindAccountGroup := ownerGroup.Group("/bind_account")
				{
					// Use email account to create helper
					bindAccountGroup.POST("/email", routers.API.Owner.Email)
					// Use mobile account to create helper
					bindAccountGroup.POST("/mobile", routers.API.Owner.Mobile)
					// Send smscode
					bindAccountGroup.POST("/send_sms", routers.API.Owner.SendSMS)
				}
			}
		}

		// Admin only
		adminGroup := apiGroup.Group("/admin")
		adminGroup.Use(middlewares.AdminPermissionHandler())
		{
			// Redeem code group
			redeemCodeGroup := adminGroup.Group("/redeem_code")
			{
				// Generate redeem code
				redeemCodeGroup.POST("/generate", routers.API.Admin.RedeemCode.Generate)
			}
			// Unlimited server group
			unlimitedServerGroup := adminGroup.Group("/unlimited_server")
			{
				// Add unlimited rental server
				unlimitedServerGroup.POST("/add", routers.API.Admin.UnlimitedServer.Add)
				// Remove unlimited rental server
				unlimitedServerGroup.POST("/delete", routers.API.Admin.UnlimitedServer.Delete)
				// Get unlimited rental server list
				unlimitedServerGroup.GET("/get_list", routers.API.Admin.UnlimitedServer.GetList)
			}
			// User group
			userGroup := adminGroup.Group("/user")
			{
				// Ban user
				userGroup.POST("/ban", routers.API.Admin.User.Ban)
				// Extend user expire time
				userGroup.POST("/extend_expire_time", routers.API.Admin.User.ExtendExpireTime)
				// Extend user unlimited time
				userGroup.POST("/extend_unlimited_time", routers.API.Admin.User.ExtendUnlimitedTime)
				// Query user info by username
				userGroup.POST("/query", routers.API.Admin.User.Query)
				// Set user permission
				userGroup.POST("/set_permission", routers.API.Admin.User.SetPermission)
				// Unban user
				userGroup.POST("/unban", routers.API.Admin.User.UnBan)
			}
		}
	}
}

func InitRouter() *gin.Engine {
	gin.SetMode(configs.GIN_MODE)

	router := gin.Default()
	router.UseH2C = true

	if configs.GIN_MODE == gin.ReleaseMode {
		router.TrustedPlatform = gin.PlatformCloudflare
	}

	router.SetTrustedProxies([]string{
		"127.0.0.1",
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	})

	// Global middleware
	router.Use(gzip.Gzip(gzip.BestCompression))
	router.Use(middlewares.CORSHandler())
	router.Use(middlewares.LogHandler())
	router.Use(middlewares.GinErrorHandler())

	apiGroup := router.Group("/api")
	{
		// Welcome
		apiGroup.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "Welcome to BunkerWeb API!")
		})

		// Create a new bearer
		apiGroup.GET("/new", routers.API.New)
	}

	// Register phoenix API
	RegisterPhoenixAPI(router)

	// Register web API
	RegisterWebAPI(router)

	// Register open API
	RegisterOpenAPI(router)

	// No router
	// router.NoRoute(func(c *gin.Context) {
	// 	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
	// 		"success": false,
	// 		"message": "404 Not Found",
	// 	})
	// })

	return router
}

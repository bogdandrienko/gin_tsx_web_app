package backend

//

// IMPORTS ////////////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// IMPORTS ////////////////////////////////////////////////////////////////////////////////////////////////////////////

//

// STRUCTS ////////////////////////////////////////////////////////////////////////////////////////////////////////////

type ConfigS struct {
	// Gin
	DebugGin bool
	HostGin  string
	PortGin  int

	// Postgresql
	DriverNamePgSQL string
	HostPgSQL       string
	PortPgSQL       int
	UserPgSQL       string
	PasswordPgSQL   string
	DatabasePgSQL   string
	SslModePgSQL    string

	// Additional
	HashCost             int
	HashSalt             string
	TokenAccessLifetime  time.Duration
	TokenRefreshLifetime time.Duration
}

type RouteS struct {
	Method       string
	RelativePath string
	Function     gin.HandlerFunc
	Auth         bool
}

type User struct {
	// main
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	IsActive bool   `json:"is_active"`

	// roles
	IsModerator bool `json:"is_moderator"`
	IsAdmin     bool `json:"is_admin"`
	IsSuperuser bool `json:"is_superuser"`

	// additional
	Name              string    `json:"name"`
	Surname           string    `json:"surname"`
	Patronymic        string    `json:"patronymic"`
	Email             string    `json:"email"`
	Phone             string    `json:"phone"`
	Avatar            string    `json:"avatar"`
	DatetimeJoined    time.Time `json:"datetime_joined"`
	DatetimeLastLogin time.Time `json:"datetime_last_login"`
}

type Token struct {
	// main
	Id       int    `json:"id"`
	Username string `json:"username"`
	Access   string `json:"access"`
	Refresh  string `json:"refresh"`

	// additional
	DatetimeCreated time.Time `json:"datetime_created"`
}

type Task struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

// STRUCTS ////////////////////////////////////////////////////////////////////////////////////////////////////////////

//

// METHODS ////////////////////////////////////////////////////////////////////////////////////////////////////////////
// METHODS ////////////////////////////////////////////////////////////////////////////////////////////////////////////

//

// GLOBALS ////////////////////////////////////////////////////////////////////////////////////////////////////////////

var ConfigV = ConfigConstructor()

// GLOBALS ////////////////////////////////////////////////////////////////////////////////////////////////////////////

//

// ACTIONS ////////////////////////////////////////////////////////////////////////////////////////////////////////////

func ConfigConstructor() *ConfigS {
	config_ := &ConfigS{
		// Gin
		HostGin:  "127.0.0.1",
		PortGin:  8080,
		DebugGin: true,

		// Postgresql
		DriverNamePgSQL: "postgres",
		HostPgSQL:       "127.0.0.1",
		PortPgSQL:       5432,
		UserPgSQL:       "pgs_usr",
		PasswordPgSQL:   "12345Qwerty!",
		DatabasePgSQL:   "pgs_db",
		SslModePgSQL:    "disable",

		// Additional
		HashCost:             14,
		HashSalt:             "Qwerty!12345",
		TokenAccessLifetime:  time.Minute * 10,
		TokenRefreshLifetime: time.Minute * 60 * 24,
	}
	return config_
}

func Run() {
	err := NewApp([]RouteS{
		{"GET", "/", IndexHandler, false}, // Read     | curl -v -X GET 127.0.0.1:8080/

		// tasks(s)
		{"POST", "/api/tasks", CreateTaskHandler, false},       // Create   | curl -v -X POST 127.0.0.1:8080/tasks -d '{"title":"Amon Ra","author":"V.Pelevin"}'
		{"GET", "/api/tasks/:id", ReadTaskHandler, false},      // Read     | curl -v -X GET 127.0.0.1:8080/tasks/1
		{"GET", "/api/tasks", ReadTasksHandler, false},         // Read all | curl -v -X GET 127.0.0.1:8080/tasks
		{"PUT", "/api/tasks/:id", UpdateTaskHandler, false},    // Update   | curl -v -X PUT 127.0.0.1:8080/tasks/1 -d '{"title":"War and peace","author":"N.Tolstoy"}'
		{"DELETE", "/api/tasks/:id", DeleteTaskHandler, false}, // Delete   | curl -v -H 'x-id:1' -X DELETE 127.0.0.1:8080/tasks/1
	})
	if err != nil {
		ErrorHandler(err)
		log.Fatal(err)
		return
	}
}

func NewApp(routes []RouteS) error {
	// todo Configure engine
	if ConfigV.DebugGin {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// todo Create engine
	engine := gin.Default()

	// todo Configure static
	engine.Static("/static", "./frontend/build/static")

	// todo Configure html
	engine.LoadHTMLGlob("./frontend/build/index.html")

	engine.Use(cors.New(cors.Config{
		AllowHeaders: []string{"*"},
		AllowMethods: []string{"*"},
		AllowOrigins: []string{"*"},
		MaxAge:       12 * time.Hour,
	}))

	// todo Redirect 404
	//engine.NoRoute(RedirectEmptyHandler)

	// todo Binding routes
	for _, route := range routes {
		switch route.Method {
		case "POST": // Create
			engine.POST(route.RelativePath, Middleware(route.Function, route))
		case "GET": // Read
			engine.GET(route.RelativePath, Middleware(route.Function, route))
		case "PUT": // Update
			engine.PUT(route.RelativePath, Middleware(route.Function, route))
		case "DELETE": // Delete
			engine.DELETE(route.RelativePath, Middleware(route.Function, route))
		default:
			ErrorHandler(errors.New("unknown method"))
		}
	}

	// todo Run engine
	return engine.Run(fmt.Sprintf("%s:%d", ConfigV.HostGin, ConfigV.PortGin))
}

func Middleware(next gin.HandlerFunc, route RouteS) gin.HandlerFunc {

	time.Sleep(time.Millisecond * 1500)

	return func(context *gin.Context) {
		//log.Printf("[%s] %s\n", context.Request.Method, context.Request.RequestURI)

		if route.Auth {
			// get config.headers.Authorization
			authorizationHeader := context.GetHeader("Authorization")
			if authorizationHeader == "" {
				ErrHandlerWithContext(errors.New("authorization failed"), context)
				return
			}

			// get JWT=$Qwerty!1234567
			tokenArr := strings.Split(authorizationHeader, "=")
			if len(tokenArr) < 2 {
				ErrHandlerWithContext(errors.New("authorization failed"), context)
				return
			}

			// check token
			//status, err := CheckAccessToken(tokenArr[1])
			//if err != nil {
			//	ErrHandlerWithContext(errors.New("authorization failed"), context)
			//	return
			//}

			//if !status {
			//	ErrHandlerWithContext(errors.New("authorization failed"), context)
			//	return
			//}

			//userID := context.Request.Header.Get("x-id")
			//if userID == "" {
			//	log.Printf("[%s] %s - error: userID is not provided\n", context.Request.Method, context.Request.RequestURI)
			//	ErrHandlerWithContext(errors.New("users id is not provided"), context)
			//	return
			//}
		}

		//
		//ctx := r.Context()
		//ctx = context.WithValue(ctx, "id", userID)
		//
		//r = r.WithContext(ctx)

		next(context)
	}
}

// ACTIONS ////////////////////////////////////////////////////////////////////////////////////////////////////////////

//

// EXTRAS /////////////////////////////////////////////////////////////////////////////////////////////////////////////
// EXTRAS /////////////////////////////////////////////////////////////////////////////////////////////////////////////

//

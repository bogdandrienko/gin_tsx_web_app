package main

//

// TODO imports ////////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// TODO imports ////////////////////////////////////////////////////////////////////////////////////////////////////////

//

// TODO structs ////////////////////////////////////////////////////////////////////////////////////////////////////////

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

// TODO structs ////////////////////////////////////////////////////////////////////////////////////////////////////////

//

// TODO methods ////////////////////////////////////////////////////////////////////////////////////////////////////////

// TODO methods ////////////////////////////////////////////////////////////////////////////////////////////////////////

//

// TODO global /////////////////////////////////////////////////////////////////////////////////////////////////////////

var ConfigV = ConfigConstructor()

// TODO global /////////////////////////////////////////////////////////////////////////////////////////////////////////

//

// TODO actions ////////////////////////////////////////////////////////////////////////////////////////////////////////

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

func ErrorHandler(err error) {
	fmt.Printf("error: %s", err.Error())
}

func ErrHandlerWithContext(err error, context *gin.Context) {
	context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	ErrorHandler(err)
}

func CreateDbPgConnection() (*sql.DB, error) {
	source := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		ConfigV.HostPgSQL, ConfigV.PortPgSQL, ConfigV.UserPgSQL, ConfigV.PasswordPgSQL, ConfigV.DatabasePgSQL, ConfigV.SslModePgSQL,
	)
	dbConnection, err := sql.Open(ConfigV.DriverNamePgSQL, source)
	if err != nil {
		return nil, err
	}

	return dbConnection, nil
}

func ExecuteSelectOneDb(object []any, query string, args ...any) error {
	dbConnection, err := CreateDbPgConnection()
	if err != nil {
		return err
	}
	defer func(dbConnection *sql.DB) {
		err = dbConnection.Close()
		if err != nil {
			return
		}
	}(dbConnection)

	rows, err := dbConnection.Query(query, args...)
	if err != nil {
		return err
	}

	rows.Next()
	err = rows.Scan(object...)
	if err != nil {
		return err
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

func ExecuteSelectManyDb(objects *[]string, query string, args ...any) error {
	dbConnection, err := CreateDbPgConnection()
	if err != nil {
		return err
	}
	defer func(dbConnection *sql.DB) {
		err = dbConnection.Close()
		if err != nil {
			return
		}
	}(dbConnection)

	rows, err := dbConnection.Query(query, args...)
	if err != nil {
		return err
	}

	usernames := make([]string, 0)
	for rows.Next() {
		var obj string
		err = rows.Scan(&obj)
		if err != nil {
			return err
		}

		usernames = append(usernames, obj)
	}
	fmt.Println(usernames)
	*objects = usernames

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

func ExecuteRowsDb(query string, args ...any) (*sql.Rows, error) {
	dbConnection, err := CreateDbPgConnection()
	if err != nil {
		return nil, err
	}
	defer func(dbConnection *sql.DB) {
		err = dbConnection.Close()
		if err != nil {
			return
		}
	}(dbConnection)

	rows, err := dbConnection.Query(query, args...)
	if err != nil {
		return nil, err
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func ExecuteInsertOrDeleteDb(query string, args ...any) error {
	dbConnection, err := CreateDbPgConnection()
	if err != nil {
		return err
	}
	defer func(dbConnection *sql.DB) {
		err = dbConnection.Close()
		if err != nil {
			return
		}
	}(dbConnection)

	dbTransaction, err := dbConnection.Begin()
	if err != nil {
		return err
	}
	defer func(dbTransaction *sql.Tx) {
		_ = dbTransaction.Rollback()
	}(dbTransaction)

	_, err = dbTransaction.Exec(query, args...)
	if err != nil {
		return err
	}

	err = dbTransaction.Commit()
	if err != nil {
		return err
	}

	return nil
}

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

func IndexHandler(context *gin.Context) {
	context.HTML(http.StatusOK, "index.html", gin.H{
		"title": "SPA",
	})
}

func CreateTaskHandler(context *gin.Context) {
	// get and check param
	title := context.PostForm("title")
	if title == "" {
		ErrHandlerWithContext(errors.New("title incorrect"), context)
		return
	}

	// insert to db
	err := ExecuteInsertOrDeleteDb("insert into tasks (title) values ($1);", title)
	if err != nil {
		ErrHandlerWithContext(err, context)
		return
	}

	context.JSON(http.StatusCreated, map[string]string{"response": "successfully created"})
}

func ReadTaskHandler(context *gin.Context) {
	// get id from url
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		ErrHandlerWithContext(err, context)
		return
	}

	// create object
	task := Task{}

	// select from db
	err = ExecuteSelectOneDb([]any{&task.Id, &task.Title}, "select id, title from tasks where id = $1", id)
	if err != nil {
		ErrHandlerWithContext(err, context)
		return
	}

	context.JSON(http.StatusOK, map[string]any{"response": task})
}

func ReadTasksHandler(context *gin.Context) {
	// select from db
	rows, err := ExecuteRowsDb("select id, title from tasks order by id asc;")
	if err != nil {
		ErrHandlerWithContext(err, context)
		return
	}

	// create objects
	tasks := make([]Task, 0)

	// fulling objects
	for rows.Next() {
		task := Task{}
		err = rows.Scan(&task.Id, &task.Title)
		if err != nil {
			ErrHandlerWithContext(err, context)
			return
		}
		tasks = append(tasks, task)
	}

	context.JSON(http.StatusOK, map[string]map[string]any{"response": {"list": tasks}})
}

func UpdateTaskHandler(context *gin.Context) {
	// get id from url
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		ErrHandlerWithContext(err, context)
		return
	}

	// get and check param
	title := context.PostForm("title")
	if title == "" {
		ErrHandlerWithContext(errors.New("title incorrect"), context)
		return
	}

	// update into db
	err = ExecuteInsertOrDeleteDb("update tasks set title=$1 where id = $2;", title, id)
	if err != nil {
		ErrHandlerWithContext(err, context)
		return
	}

	context.JSON(http.StatusOK, map[string]string{"response": "successfully updated"})
}

func DeleteTaskHandler(context *gin.Context) {
	// get id from url
	id, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		ErrHandlerWithContext(err, context)
		return
	}

	// delete from db
	err = ExecuteInsertOrDeleteDb("delete from tasks where id = $1;", id)
	if err != nil {
		ErrHandlerWithContext(err, context)
		return
	}

	context.JSON(http.StatusOK, map[string]string{"response": "successfully deleted"})
}

// TODO actions ////////////////////////////////////////////////////////////////////////////////////////////////////////

//

// TODO main ///////////////////////////////////////////////////////////////////////////////////////////////////////////

func run() {
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

func main() {
	run()

	//CreateTasksDatabase()
}

// TODO main ///////////////////////////////////////////////////////////////////////////////////////////////////////////

//

// TODO extra //////////////////////////////////////////////////////////////////////////////////////////////////////////

// TODO extra //////////////////////////////////////////////////////////////////////////////////////////////////////////

//

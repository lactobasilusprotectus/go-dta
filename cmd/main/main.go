package main

import (
	authDelivery "go-dts/pkg/auth/delivery"
	authUsecase "go-dts/pkg/auth/usecase"
	commentDelivery "go-dts/pkg/comment/delivery"
	commentRepository "go-dts/pkg/comment/repository"
	commentUsecase "go-dts/pkg/comment/usecase"
	"go-dts/pkg/common/config"
	commonTime "go-dts/pkg/common/time"
	"go-dts/pkg/domain"
	photoDelivery "go-dts/pkg/photo/delivery"
	photoRepository "go-dts/pkg/photo/repository"
	photoUsecase "go-dts/pkg/photo/usecase"
	rootDelivery "go-dts/pkg/root/delivery"
	socialMediaDelivery "go-dts/pkg/socialmedia/delivery"
	socialMediaRepository "go-dts/pkg/socialmedia/repository"
	socialMediaUsecase "go-dts/pkg/socialmedia/usecase"
	userRepository "go-dts/pkg/user/repository"
	"go-dts/pkg/util/cache"
	"go-dts/pkg/util/db"
	_ "go-dts/pkg/util/http"
	httputil "go-dts/pkg/util/http"
	"go-dts/pkg/util/jwt"
	"log"
	"os"
	"os/signal"
	"reflect"
	"syscall"
)

//	@title						API Documentation for DTS project
//	@version					1.0
//	@description				This is the API documentation for DTS project
//	@termsOfService				http://swagger.io/terms/
//	@license.name				Apache 2.0
//	@license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//	@host						localhost:8080
//	@BasePath					/
//	@schemes					http
//	@securityDefinitions.apikey	JWT
//	@in							header
//	@name						Authorization
//	@Security					JWT
func main() {
	// Get env
	env := os.Getenv(config.ENV)

	if env == "" {
		env = config.LOC
	}

	// Read env file
	cfg, err := config.Read(config.GetFilePath(env))

	if err != nil {
		return
	}

	// Init utils: http server, db connection, etc.
	utils := initUtils(cfg)

	// Init repository and use case layer
	_, uc, err := initRepoAndUseCases(utils, cfg)
	if err != nil {
		log.Fatalln("initRepoAndUseCases err:", err)
	}

	// Init delivery layer (HTTP)
	httpHandler := initHttpHandler(utils, uc, env)

	// Start serving
	registerHttpHandler(utils.HttpServer, httpHandler)
	utils.HttpServer.Run(env)

	// =======================================================

	// Catching signals

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-c

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	utils.HttpServer.Stop()

	log.Println("shutting down...")
	os.Exit(0)

}

// initUtils initialises utility for the app
// this includes httputil server, db connection, and any other dependent tools
func initUtils(cfg config.Config) AppUtil {
	// database connection
	dbConn, err := db.NewDatabaseConnection(cfg.Database)

	//migration base on AppModels
	registerModels(dbConn, AppModels{})

	if err != nil {
		log.Fatalln(err)
	}

	// redis as cache
	redisClient := cache.NewRedisClient(cfg.Redis)

	// time module
	timeModule := commonTime.New()

	// JWT implementation
	jwtModule := jwt.New(timeModule)

	return AppUtil{
		HttpServer:   httputil.NewServer(cfg.Http),
		DbConnection: dbConn,
		Cache:        redisClient,
		Jwt:          jwtModule,
		Time:         timeModule,
	}
}

// initHttpHandler initialises http handler for the app
func initHttpHandler(ut AppUtil, uc AppUseCase, env string) AppHttpHandler {
	rootHandler := rootDelivery.NewRootHandler(env)

	return AppHttpHandler{
		RootHttpHandler:        rootHandler,
		AuthHttpHandler:        authDelivery.NewAuthHttpHandler(uc.AuthUseCase, uc.AuthUseCase),
		PhotoHttpHandler:       photoDelivery.NewPhotoHttpHandler(uc.PhotoUseCase, uc.AuthUseCase),
		CommentHttpHandler:     commentDelivery.NewCommentHttpHandler(uc.CommentUseCase, uc.AuthUseCase),
		SocialMediaHttpHandler: socialMediaDelivery.NewSocialMediaHttpHandler(uc.SocialMediaUseCase, uc.AuthUseCase),
	}
}

// initRepoAndUseCases initialises repo and use case layer
func initRepoAndUseCases(util AppUtil, cfg config.Config) (repo AppRepo, uc AppUseCase, err error) {
	repo.User = userRepository.NewUserRepository(util.DbConnection, util.Time)
	repo.Photo = photoRepository.NewPhotoRepository(util.DbConnection, util.Time)
	repo.Comment = commentRepository.NewCommentRepository(util.DbConnection, util.Time)
	repo.SocialMedia = socialMediaRepository.NewSocialMediaRepository(util.DbConnection, util.Time)

	//usecase
	uc.AuthUseCase = authUsecase.NewAuthUseCase(repo.User, util.Jwt, util.Cache, util.Time)
	uc.PhotoUseCase = photoUsecase.NewPhotoUseCase(repo.Photo)
	uc.CommentUseCase = commentUsecase.NewCommentUseCase(repo.Comment)
	uc.SocialMediaUseCase = socialMediaUsecase.NewSocialMediaUseCase(repo.SocialMedia)

	return repo, uc, nil
}

// registerHttpHandler registers our handlers to the http server.
// reflect docs: https://golang.org/pkg/reflect/.
// The purpose of this function is to register HTTP request handlers to an httputil.Server object based on the fields of the handlers object.
func registerHttpHandler(srv *httputil.Server, handlers AppHttpHandler) {
	h := reflect.ValueOf(handlers)

	for i := 0; i < h.NumField(); i++ {
		srv.RegisterHandler(h.Field(i).Interface().(httputil.RouterHandler))
	}
}

// registerModels registers our models to the database.
// The purpose of this function is to automatically create tables or migrate existing ones in the database based on the fields of the models object.
func registerModels(dbConn *db.DatabaseConnection, models AppModels) {
	m := reflect.ValueOf(models)

	for i := 0; i < m.NumField(); i++ {

		err := dbConn.Master.AutoMigrate(m.Field(i).Interface())

		if err != nil {
			log.Fatalf("AutoMigrate err: %v", err)
		}
	}
}

//================ TYPES =================

// AppUtil wraps utility layer with the app, includes delivery and database
type AppUtil struct {
	HttpServer   *httputil.Server
	DbConnection *db.DatabaseConnection
	Cache        cache.KVCacheInterface
	Jwt          *jwt.JwtModule
	Time         *commonTime.Time
}

// AppHttpHandler wraps HTTP handlers exposed by the app as a delivery layer
type AppHttpHandler struct {
	RootHttpHandler        *rootDelivery.RootHandler
	AuthHttpHandler        *authDelivery.AuthHttpHandler
	PhotoHttpHandler       *photoDelivery.PhotoHttpHandler
	CommentHttpHandler     *commentDelivery.CommentHttpHandler
	SocialMediaHttpHandler *socialMediaDelivery.SocialMediaHttpHandler
}

// AppUseCase wraps use case layer within the app
type AppUseCase struct {
	AuthUseCase        *authUsecase.AuthUseCase
	PhotoUseCase       *photoUsecase.PhotoUseCase
	CommentUseCase     *commentUsecase.CommentUseCase
	SocialMediaUseCase *socialMediaUsecase.SocialMediaUseCase
}

// AppRepo wraps repository layer within the app
type AppRepo struct {
	User        *userRepository.UserRepository
	Photo       *photoRepository.PhotoRepository
	Comment     *commentRepository.CommentRepository
	SocialMedia *socialMediaRepository.SocialMediaRepository
}

// AppModels wraps domain models within the app
type AppModels struct {
	User        *domain.User
	Photo       *domain.Photo
	SocialMedia *domain.SocialMedia
	Comment     *domain.Comment
}

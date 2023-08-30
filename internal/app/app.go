package app

import (
	"avito/config"
	"avito/internal/handlers"
	"avito/internal/repository"
	"avito/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	handlers   *handlers.Handlers
	useCase    *usecase.UseCase
	repository *repository.Repository
	server     *fiber.App
}

func InitApp(conf *config.Config) *App {
	app := App{}
	app.server = fiber.New()
	app.repository = repository.InitRepository(conf)
	app.useCase = usecase.InitUsecase(app.repository)
	app.handlers = handlers.InitHandlers(app.useCase)
	app.server.Get("/getById", app.handlers.GetById)
	app.server.Post("/createSegment", app.handlers.CreateSegment)
	app.server.Delete("/deleteSegment", app.handlers.DeleteSegment)
	app.server.Get("/getSegment", app.handlers.GetSegments)
	app.server.Put("/subscription", app.handlers.Subscriber)
	app.server.Put("/timeoutSubscribe", app.handlers.SubscribeWithTimeOut)
	return &app
}

func (app *App) AppStart() {
	sig := make(chan os.Signal, 1)
	sec := make(chan int, 1)
	go app.repository.TimeOutDeleter(sec)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		sec <- 1
		log.Println("Gracefully shutdown")
		if err := app.server.ShutdownWithTimeout(30 * time.Second); err != nil {
			log.Fatalln("server shutdown error: ", err)
		}
	}()
	err := app.server.Listen(":3000")
	if err != nil {
		log.Panic(err.Error())
	}

}

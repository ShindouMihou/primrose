package clients

import (
	"github.com/go-playground/validator/v10"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/middleware/cors"
)

type IIris struct {
	Client
	Module      *iris.Application
	Executables []Executable
	Logger      *golog.Logger
}

type Executable = func(app *iris.Application)

var Iris = IIris{Module: iris.New(), Executables: []Executable{}, Logger: golog.New().SetPrefix("[IRIS] ")}

func (f *IIris) Attach(executable Executable) Executable {
	f.Executables = append(f.Executables, executable)
	return executable
}

func (f *IIris) Init() {
	go func() {
		f.Logger.Info("Starting HTTP server.")
		f.Module.Use(iris.Compression)
		f.Module.Use(func(c *context.Context) {
			f.Logger.Info(c.Method(), " ", c.Path(),
				" with {\"user_agent\":\"", c.Request().UserAgent(),
				"\",\"referrer\":\"", c.Request().Referer(),
				"\",\"query\":\"", c.URLParams(), "\",\"remote_address\":\"", c.RemoteAddr(), "\"}")
			c.Next()
		})
		f.Module.UseRouter(cors.New().AllowOriginFunc(cors.AllowAnyOrigin).Handler())
		f.Module.Validator = validator.New()
		f.Module.Configure(iris.WithOptimizations, iris.WithEasyJSON, iris.WithoutBanner, iris.WithoutInterruptHandler, iris.WithProtoJSON)
		for _, executable := range f.Executables {
			executable(f.Module)
		}
		f.Logger.Info("HTTP server is now running in http://localhost:9045")
		err := f.Module.Listen(":9045")
		if err != nil {
			f.Logger.Fatal("Failed to start HTTP server: ", err)
		}
	}()
}

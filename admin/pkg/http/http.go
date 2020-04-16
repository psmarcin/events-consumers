package http

import (
	"log"
	"net/http"

	"github.com/gofiber/basicauth"
	"github.com/gofiber/fiber"
	"github.com/gofiber/helmet"
	"github.com/gofiber/logger"
	"github.com/gofiber/recover"
	"github.com/gofiber/requestid"

	"events-consumers/admin/pkg/config"
)

func Start(dependencies Dependencies) {
	cfg := recover.Config{
		Handler: ErrorHandler,
	}
	basicAuthCfg := basicauth.Config{
		Users: map[string]string{
			config.C.BasicAuthUser:  config.C.BasicAuthPassword,
		},
	}

	app := fiber.New()
	app.Settings.CaseSensitive = true
	app.Use(basicauth.New(basicAuthCfg))
	app.Use(helmet.New())
	app.Use(requestid.New())
	app.Use(logger.New())
	app.Use(recover.New(cfg))

	Router(app, dependencies)

	err := app.Listen(config.C.HttpPort)
	if err != nil {
		log.Fatal(err)
	}
}

func Router(app *fiber.App, dependencies Dependencies) {
	app.Get("/", IndexHandler(dependencies))
	app.Get("/job/:jobId/edit", JobEditFormHandler(dependencies))
	app.Post("/job/:jobId/edit", JobEditHandler(dependencies))
	app.Get("/job/create", JobCreateFormHandler())
	app.Post("/job/create", JobCreateHandler(dependencies))
	app.Post("/job/:jobId/delete", JobDeleteHandler(dependencies))
}

func IndexHandler(dependencies Dependencies) func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		jobs, err := dependencies.Job.List()
		if err != nil {
			c.Next(err)
		}

		if err := c.Render(config.C.TemplatePath + "views/index.tmpl", jobs); err != nil {
			c.Status(500).Send(err.Error())
		}
	}
}

func JobEditFormHandler(dependencies Dependencies) func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		jobId := c.Params("jobId")
		if jobId == "" {
			c.Status(http.StatusBadRequest).Send("Job identifier doesn't provided")
		}

		log.Printf("jobid %s", jobId)

		job, err := dependencies.Job.Get(jobId)
		if err != nil {
			c.Next(err)
		}

		log.Printf("job %+v", job)

		if err := c.Render(config.C.TemplatePath + "views/edit.tmpl", job); err != nil {
			c.Status(500).Send(err.Error())
		}
	}
}
func JobEditHandler(dependencies Dependencies) func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		jobId := c.Params("jobId")
		if jobId == "" {
			c.Status(http.StatusBadRequest).Send("Job identifier doesn't provided")
		}

		_, err := dependencies.Job.Update(jobId, c.Body("command"), c.Body("name"), c.Body("selector"))
		if err != nil {
			c.Next(err)
		}
		c.Redirect("/")
	}
}

func JobCreateFormHandler() func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		if err := c.Render(config.C.TemplatePath + "views/create.tmpl", nil); err != nil {
			c.Status(500).Send(err.Error())
		}
	}
}
func JobCreateHandler(dependencies Dependencies) func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		_, err := dependencies.Job.Create(c.Body("command"), c.Body("name"), c.Body("selector"))
		if err != nil {
			c.Next(err)
		}
		c.Redirect("/")
	}
}
func JobDeleteHandler(dependencies Dependencies) func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		jobId := c.Params("jobId")
		if jobId == "" {
			c.Status(http.StatusBadRequest).Send("Job identifier doesn't provided")
		}
		err := dependencies.Job.Delete(jobId)
		if err != nil {
			c.Next(err)
		}
		c.Redirect("/")
	}
}

func ErrorHandler(c *fiber.Ctx, err error) {
	c.SendString(err.Error())
	c.SendStatus(500)
}

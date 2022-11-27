package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/kardianos/service"
	"github.com/urfave/cli/v2"
)

var (
	isDev    = os.Getenv("GO_ENV") == "development"
	logger   service.Logger
	apiURL   = "https://api.yizcore.xyz"
	password string
)

// Program structures.
//
//	Define Start and Stop methods.
type program struct {
	exit chan struct{}
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		// logger.Info("Running in terminal.")

		var (
			email    string
			deviceId string
		)

		app := &cli.App{
			Name:  "file-sync",
			Usage: "Automatically sync single file.",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "login",
					Value:       "",
					Usage:       "login by email",
					Destination: &email,
					Action: func(ctx *cli.Context, s string) error {
						prompt := &survey.Password{
							Message: "Please type your password",
						}
						survey.AskOne(prompt, &password)
						hostname, err := os.Hostname()
						if err != nil {
							panic(err)
						}

						machineName := ""
						prompt2 := &survey.Input{
							Message: "Please type your device name",
							Default: hostname,
						}
						survey.AskOne(prompt2, &machineName)

						err = login(email, password, machineName)
						if err != nil {
							fmt.Println(err.Error())
						} else {
							fmt.Println("Login and register your device successfully!")
						}
						return nil
					},
				},
				&cli.StringFlag{
					Name:        "config",
					Value:       "https://file-sync.openapi.site/",
					Usage:       "HTTP API server URL",
					Destination: &email,
				},
				&cli.StringFlag{
					Name:        "remove-device",
					Aliases:     []string{"rd"},
					Value:       "current device",
					Usage:       "remove device by device id",
					Destination: &deviceId,
				},
			},
			Commands: []*cli.Command{
				{
					Name:    "add",
					Aliases: []string{"a"},
					Usage:   "Add a file to sync list",
					Action: func(cCtx *cli.Context) error {
						fmt.Println("added task: ", cCtx.Args().First())
						return nil
					},
				},
				{
					Name:    "list",
					Aliases: []string{"l"},
					Usage:   "Sync files list",
					Action: func(cCtx *cli.Context) error {
						fmt.Println("completed task: ", cCtx.Args().First())
						return nil
					},
				},
				{
					Name:    "remove",
					Aliases: []string{"r"},
					Usage:   "Remove a file config in sync list",
					Action: func(cCtx *cli.Context) error {
						fmt.Println("completed task: ", cCtx.Args().First())
						return nil
					},
				},
				{
					Name:    "service",
					Aliases: []string{"s"},
					Usage:   "Control the system service",
					Subcommands: []*cli.Command{
						{
							Name:  "enable",
							Usage: "set to boot service",
							Action: func(cCtx *cli.Context) error {
								s.Install()
								return nil
							},
						},
						{
							Name:  "disable",
							Usage: "disable the service",
							Action: func(cCtx *cli.Context) error {
								s.Uninstall()
								return nil
							},
						},
						{
							Name:  "start",
							Usage: "start the file sync service",
							Action: func(cCtx *cli.Context) error {
								s.Start()
								return nil
							},
						},
						{
							Name:  "stop",
							Usage: "stop the file sync service",
							Action: func(cCtx *cli.Context) error {
								s.Stop()
								return nil
							},
						},
						{
							Name:  "status",
							Usage: "get the service status",
							Action: func(cCtx *cli.Context) error {
								status, err := s.Status()
								if err != nil {
									fmt.Print("Info: ", err.Error())
									return nil
								}
								fmt.Print(status)
								return nil
							},
						},
					},
				},
			},
			Action: func(*cli.Context) error {
				fmt.Println("boom! I say!")
				return nil
			},
		}

		if err := app.Run(os.Args); err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	} else {
		logger.Info("Running under service manager.")
	}
	p.exit = make(chan struct{})

	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) run() error {
	logger.Infof("I'm running %v.", service.Platform())
	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case tm := <-ticker.C:
			logger.Infof("Still running at %v...", tm)
		case <-p.exit:
			ticker.Stop()
			return nil
		}
	}
}

func (p *program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	logger.Info("I'm Stopping!")
	close(p.exit)
	return nil
}

// Service setup.
//
//	Define service config.
//	Create the service.
//	Setup the logger.
//	Handle service controls (optional).
//	Run the service.
func main() {
	// svcFlag := flag.String("service", "", "Control the system service.")
	// flag.Parse()

	if isDev {
		apiURL = "http://127.0.0.1:8000"
	}

	options := make(service.KeyValue)
	options["Restart"] = "on-success"
	options["SuccessExitStatus"] = "1 2 8 SIGKILL"
	svcConfig := &service.Config{
		Name:        "file-sync",
		DisplayName: "File-sync service",
		Description: "The file-sync tool service.",
		Dependencies: []string{
			"Requires=network.target",
			"After=network-online.target syslog.target"},
		Option: options,
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	// if len(*svcFlag) != 0 {
	// 	err := service.Control(s, *svcFlag)
	// 	if err != nil {
	// 		log.Printf("Valid actions: %q\n", service.ControlAction)
	// 		log.Fatal(err)
	// 	}
	// 	return
	// }
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}

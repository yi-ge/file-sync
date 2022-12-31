package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/fatih/color"
	"github.com/kardianos/service"
	"github.com/urfave/cli/v2"
	"github.com/yi-ge/file-sync/utils"
)

var (
	isDev    = os.Getenv("GO_ENV") == "development"
	logger   service.Logger
	apiURL   = "https://api.yizcore.xyz"
	password string
	data     Data
)

// Define Start and Stop methods.
type program struct {
	exit chan struct{}
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		// logger.Info("Running in terminal.")

		app := &cli.App{
			Name:  "file-sync",
			Usage: "Automatically sync single file.",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "login",
					Value: "",
					Usage: "login by email",
					Action: func(ctx *cli.Context, email string) error {
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
							color.Red("Login failure.")
							color.Red(err.Error())
						} else {
							color.Blue("Login and register your device successfully!")
						}
						return nil
					},
				},
				&cli.StringFlag{
					Name:  "config",
					Value: "https://file-sync.openapi.site/",
					Usage: "HTTP API server URL",
					Action: func(ctx *cli.Context, url string) error {
						if url != "" && utils.CheckURL(url) != nil {
							color.Red("Invalid URL")
							return nil
						}

						err := setConfig(url)

						if err != nil {
							color.Red(err.Error())
						}

						return nil
					},
				},
				&cli.BoolFlag{
					Name:     "info",
					Required: false,
					Usage:    "Display system information",
					Action: func(ctx *cli.Context, b bool) error {
						color.Blue("HTTP API server URL: " + apiURL)

						return nil
					},
				},
				&cli.BoolFlag{
					Name:        "remove-device",
					Aliases:     []string{"rd"},
					DefaultText: "current machine",
					Required:    false,
					Usage:       "remove device by device id",
					Action: func(ctx *cli.Context, b bool) error {
						s := ""
						if ctx.NArg() > 0 {
							s = ctx.Args().Get(0)
						}
						removeMachineId := ""
						removeMachineName := ""
						data, err := getData()
						if err != nil {
							color.Red(err.Error())
						}

						if s == "" {
							removeMachineId = data.MachineId
							removeMachineName = data.MachineName
						} else {
							devices, err := listDevices(data)
							if err != nil {
								color.Red(err.Error())
								return err
							}

							for i := 0; i < devices.Size(); i++ {
								machineId := devices.Get(i, "machineId").ToString()
								if strings.Contains(machineId, s) {
									removeMachineId = machineId
									removeMachineName = devices.Get(i, "machineName").ToString()
									break
								}
							}

							if removeMachineId == "" {
								pattern := "\\d+"
								result, err := regexp.MatchString(pattern, s)
								if err != nil {
									color.Red(err.Error())
									return err
								}

								if result {
									index, err := strconv.Atoi(s)
									if err != nil {
										color.Red(err.Error())
										return err
									}
									removeMachineId = devices.Get(index-1, "machineId").ToString()
									removeMachineName = devices.Get(index-1, "machineName").ToString()
								} else {
									err = errors.New("invalid machine id")
									color.Red(err.Error())
									return err
								}
							}
						}

						if removeMachineId == "" {
							err = errors.New("invalid machine id")
							color.Red(err.Error())
							return err
						}

						del := false
						promptDel := &survey.Confirm{
							Message: "Are you sure to remove the device (" + removeMachineName + " ID:" + removeMachineId[:10] + ")?",
						}
						survey.AskOne(promptDel, &del)

						if !del {
							return nil
						}

						prompt := &survey.Password{
							Message: "Please type your password",
						}
						survey.AskOne(prompt, &password)

						machineKey, res := checkPassword(data, password)
						if res && machineKey != "" {
							res, err := removeDevice(machineKey, data, removeMachineId)
							if err != nil {
								color.Red("Remove device failure.")
								color.Red(err.Error())
							}

							if res != "" {
								color.Green("Remove device successfully!")
							} else {
								color.Red("Remove device failure. Unknown error.")
							}
						} else {
							color.Red("Password incorrect!")
						}

						return nil
					},
				},
				&cli.BoolFlag{
					Name:    "list-device",
					Aliases: []string{"ld"},
					Usage:   "list registered device",
					Action: func(cCtx *cli.Context, b bool) error {
						if b {
							data, err := getData()
							if err != nil {
								color.Red(err.Error())
							}

							devices, err := listDevices(data)
							if err != nil {
								color.Red(err.Error())
							}
							displayRowSet := mapset.NewSet("id", "machineKey")
							if devices.Size() > 0 {
								printTable(devices, displayRowSet)
							} else {
								color.Red("No registered devices.")
							}
						}

						return nil
					},
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
						data, err := getData()
						if err != nil {
							color.Red(err.Error())
						}

						configs, err := listConfigs(data)
						if err != nil {
							color.Red(err.Error())
						}
						displayRowSet := mapset.NewSet("id")
						if configs.Size() > 0 {
							printTable(configs, displayRowSet)
						} else {
							color.Red("No file config.")
						}
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
				// fmt.Println("boom! I say!")
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

	err := fsInit()
	if err != nil {
		log.Fatal(err)
	}

	if isDev {
		log.Printf("Currently in development mode!")
		apiURL = "http://localhost:8000"
	} else {
		config := getConfig()
		if config != "" {
			apiURL = config
		}
	}

	err = dataInit()
	if err != nil {
		log.Fatal(err)
	}

	err = cacheInit()
	if err != nil {
		log.Fatal(err)
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

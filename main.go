package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"
	"github.com/kardianos/service"
	"github.com/urfave/cli/v3"
	"github.com/yi-ge/file-sync/config"
	"github.com/yi-ge/file-sync/utils"
)

var (
	isDev          = false
	logger         service.Logger
	apiURL         = "https://file-sync.yizcore.xyz"
	password       string
	configInstance = config.Instance()
	watcher        *fsnotify.Watcher
)

// Define Start and Stop methods.
type program struct {
	exit chan struct{}
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		// logger.Info("Running in terminal.")

		app := &cli.App{
			Name:    "file-sync",
			Version: configInstance.GetVersion(),
			Usage:   "Automatically sync single file.",
			Suggest: true,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "login",
					Value: "",
					Usage: "login by email",
					Action: func(ctx *cli.Context, email string) error {
						machineName := ""
						if ctx.NArg() == 2 {
							password = ctx.Args().Get(0)
							machineName = ctx.Args().Get(1)
						} else {
							if ctx.NArg() == 1 {
								password = ctx.Args().Get(0)
							} else {
								prompt := &survey.Password{
									Message: "Enter your password: ",
								}
								survey.AskOne(prompt, &password)
							}

							hostname, err := os.Hostname()
							if err != nil {
								panic(err)
							}

							prompt2 := &survey.Input{
								Message: "Enter a name for this machine: ",
								Default: hostname,
							}
							survey.AskOne(prompt2, &machineName)
						}

						err := login(email, password, machineName)
						if err != nil {
							color.Red("Login failure.")
							color.Red(err.Error())
						} else {
							color.Blue("Login and register your device successfully!")
						}
						return nil
					},
				},
				&cli.BoolFlag{
					Name:        "config",
					DefaultText: "https://file-sync.openapi.site/",
					Usage:       "HTTP API server URL",
					Required:    false,
					Action: func(ctx *cli.Context, b bool) error {
						url := ""
						if ctx.NArg() > 0 {
							url = ctx.Args().Get(0)
						} else {
							reset := false
							promptReset := &survey.Confirm{
								Message: "Are you reset server URL?",
							}
							survey.AskOne(promptReset, &reset)

							if reset {
								err := delConfig()
								if err != nil {
									color.Red(err.Error())
									return nil
								}
								delCache()

								color.Blue("Reset server URL successfully!")
								return nil
							} else {
								return nil
							}
						}

						if utils.CheckURL(url) != nil {
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
					Name: "info",
					// Aliases:  []string{"version"},
					Required: false,
					Usage:    "display system information",
					Action: func(ctx *cli.Context, b bool) error {
						color.Blue("file-sync version: " + configInstance.GetVersion())
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
							if ctx.NArg() > 1 {
								password = ctx.Args().Get(1)
							}
						}
						removeMachineId := ""
						removeMachineName := ""
						data, err := getData()
						if err != nil {
							color.Red(err.Error())
							return nil
						}

						if s == "" || s == "current" {
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
									encryptedRemoveMachineName := devices.Get(i, "machineName").ToString()
									removeMachineName, err = utils.AESCTRDecryptWithBase64(encryptedRemoveMachineName, []byte(data.Verify))
									if err != nil {
										color.Red(err.Error())
										return nil
									}
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
									encryptedRemoveMachineName := devices.Get(index-1, "machineName").ToString()
									removeMachineName, err = utils.AESCTRDecryptWithBase64(encryptedRemoveMachineName, []byte(data.Verify))
									if err != nil {
										color.Red(err.Error())
										return nil
									}
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

						if ctx.NArg() < 2 {
							del := false
							promptDel := &survey.Confirm{
								Message: "Are you sure to remove the device (" + removeMachineName + " ID:" + removeMachineId[:10] + ")?",
							}
							survey.AskOne(promptDel, &del)

							if !del {
								return nil
							}

							if password == "" {
								prompt := &survey.Password{
									Message: "Enter your password: ",
								}
								survey.AskOne(prompt, &password)
							}
						}

						machineKey, res := checkPassword(data, password)
						if res && machineKey != "" {
							res, err := removeDevice(machineKey, data, removeMachineId)
							if err != nil {
								color.Red("Remove device failure.")
								color.Red(err.Error())
							}

							if res != "" {
								color.Green("Remove device successfully!")
								delCache()
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
								return nil
							}

							devices, err := listDevices(data)
							if err != nil {
								color.Red(err.Error())
								return nil
							}
							displayRowSet := mapset.NewSet("id", "machineKey")
							if devices.Size() > 0 {
								printDeviceTable(devices, displayRowSet, true, false, []byte(data.Verify))
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
					Name:            "add",
					Aliases:         []string{"a"},
					Usage:           "Add a file to sync list",
					HelpName:        "add",
					SkipFlagParsing: false,
					ArgsUsage:       "[fileId] [path]",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "name",
							DefaultText: "defined by file",
							Value:       "",
							Usage:       "name of the file display",
						},
						&cli.StringFlag{
							Name:        "machineId",
							DefaultText: "current device",
							Value:       "",
							Usage:       "target machine ID",
						},
					},
					Action: func(cCtx *cli.Context) error {
						if cCtx.NArg() == 0 {
							return cli.ShowSubcommandHelp(cCtx)
						}

						data, err := getData()
						if err != nil {
							color.Red(err.Error())
							return nil
						}

						var (
							actionMachineId = ""
							fileName        = ""
							filePath        = ""
							fileId          = ""
							isNewFile       = false
							sha256          = ""
						)

						if cCtx.Args().Get(1) != "" {
							inputArg := cCtx.Args().Get(0)
							filePath = cCtx.Args().Get(1)

							privateKeyEncrypted, err := getPrivateKey()
							if err != nil {
								color.Red(err.Error())
								return nil
							}

							privateKeyHex, err := base64.RawURLEncoding.DecodeString(string(privateKeyEncrypted))
							if err != nil {
								color.Red(err.Error())
								return nil
							}

							decrypted, privateKey, err := utils.AESMACDecryptBytes(privateKeyHex, data.RsaPrivateKeyPassword)

							if err != nil || !decrypted {
								color.Red((errors.New("secret decrypt error: " + err.Error())).Error())
								return nil
							}

							configs, err := listConfigs(data)
							if err != nil {
								color.Red(err.Error())
								return nil
							}

							for i := 0; i < configs.Size(); i++ {
								theFileId := configs.Get(i, "fileId").ToString()
								if strings.Contains(theFileId, inputArg) {
									fileId = theFileId
									fileNameEncrypted := configs.Get(i, "fileName").ToString()
									fileNameBase64, err := base64.URLEncoding.DecodeString(fileNameEncrypted)
									if err != nil {
										color.Red(err.Error())
										return nil
									}
									fileName = string(utils.RsaDecrypt([]byte(fileNameBase64), privateKey))
									break
								}
							}

							if fileId == "" {
								pattern := "\\d+"
								result, err := regexp.MatchString(pattern, inputArg)
								if err != nil {
									color.Red(err.Error())
									return err
								}

								if result {
									index, err := strconv.Atoi(inputArg)
									if err != nil {
										color.Red(err.Error())
										return err
									}
									fileId = configs.Get(index-1, "fileId").ToString()
									fileNameEncrypted := configs.Get(index-1, "fileName").ToString()
									fileNameBase64, err := base64.URLEncoding.DecodeString(fileNameEncrypted)
									if err != nil {
										color.Red(err.Error())
										return nil
									}
									fileName = string(utils.RsaDecrypt([]byte(fileNameBase64), privateKey))
								} else {
									err = errors.New("invalid file id")
									color.Red(err.Error())
									return err
								}
							}

							if fileId == "" {
								err = errors.New("invalid file id")
								color.Red(err.Error())
								return err
							}

							color.Blue("Action file (" + fileName + " ID:" + fileId[:10] + ").")
						} else {
							filePath = cCtx.Args().Get(0)
							sha256, err = utils.FileSHA256(filePath)
							if err != nil {
								color.Red(err.Error())
							}
							fileId = utils.GetSha1Str(sha256)
							isNewFile = true
						}

						if !filepath.IsAbs(filePath) {
							filePath, err = filepath.Abs(filePath)
							if err != nil {
								color.Red(err.Error())
								return nil
							}
						}

						if cCtx.String("name") != "" {
							fileName = cCtx.String("name")
						} else {
							fileName = filepath.Base(filePath)
						}

						if cCtx.String("machineId") != "" {
							s := cCtx.String("machineId")
							actionMachineId = ""
							actionMachineName := ""

							if s == "" {
								actionMachineId = data.MachineId
								actionMachineName = data.MachineName
							} else {
								devices, err := listDevices(data)
								if err != nil {
									color.Red(err.Error())
									return err
								}

								for i := 0; i < devices.Size(); i++ {
									machineId := devices.Get(i, "machineId").ToString()
									if strings.Contains(machineId, s) {
										actionMachineId = machineId
										encryptedActionMachineName := devices.Get(i, "machineName").ToString()
										actionMachineName, err = utils.AESCTRDecryptWithBase64(encryptedActionMachineName, []byte(data.Verify))
										if err != nil {
											color.Red(err.Error())
											return nil
										}
										break
									}
								}

								if actionMachineId == "" {
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
										actionMachineId = devices.Get(index-1, "machineId").ToString()
										encryptedActionMachineName := devices.Get(index-1, "machineName").ToString()
										actionMachineName, err = utils.AESCTRDecryptWithBase64(encryptedActionMachineName, []byte(data.Verify))
										if err != nil {
											color.Red(err.Error())
											return nil
										}
									} else {
										err = errors.New("invalid machine id")
										color.Red(err.Error())
										return err
									}
								}
							}

							if actionMachineId == "" {
								err = errors.New("invalid machine id")
								color.Red(err.Error())
								return err
							}

							color.Blue("Action machine: " + actionMachineName + "(" + actionMachineId + ")")
						} else {
							actionMachineId = data.MachineId
						}

						// fmt.Println("name: ", name)
						// fmt.Println("fileId: ", fileId)
						// fmt.Println("actionMachineId: ", actionMachineId)
						// fmt.Println("path: ", filePath)

						json, err := addConfig(fileId, fileName, filePath, actionMachineId, data)
						if err != nil {
							color.Red(err.Error())
							return nil
						}
						color.Blue("The file (" + fileName + " ID:" + json.Get("fileId").ToString()[:10] + ") was successfully added to the sync item.")
						delCache()

						if isNewFile {
							fileName = filepath.Base(filePath)
							f, err := os.ReadFile(filePath)
							if err != nil {
								fmt.Println("read fail", err)
							}
							fileContent := string(f)
							timestamp := time.Now().UnixNano() / 1e6
							err = fileUpload(fileId, fileName, sha256, fileContent, timestamp, data)
							if err != nil {
								color.Red(err.Error())
								return nil
							}
						} else {
							exists, _ := utils.FileExists(filePath)
							if exists {
								color.Blue("The file (" + fileName + " ID:" + json.Get("fileId").ToString()[:10] + ") already exists.")
								actionVersion := ""
								prompt := &survey.Select{
									Message: "Please select the version you want to keep:",
									Options: []string{"remote", "local"},
								}
								survey.AskOne(prompt, &actionVersion)
								if actionVersion == "local" {
									fileName = filepath.Base(filePath)
									f, err := os.ReadFile(filePath)
									if err != nil {
										fmt.Println("read fail", err)
									}
									fileContent := string(f)
									timestamp := time.Now().UnixNano() / 1e6
									err = fileUpload(fileId, fileName, sha256, fileContent, timestamp, data)
									if err != nil {
										color.Red(err.Error())
										return nil
									}
									return nil
								}
							}
							filePath, err = utils.CreateDirectoryIfNotExists(filePath, fileName)
							if err != nil {
								color.Red(err.Error())
								return nil
							}

							json, err := fileDownload(fileId, data)
							if err != nil {
								color.Red(err.Error())
							}

							content := json.Get("content").ToString()

							res, err := utils.WriteByteFile(fileName, []byte(content), 0, true)
							if err != nil {
								color.Red(err.Error())
							}

							if !res {
								color.Red("The file (" + fileName + " ID:" + json.Get("fileId").ToString()[:10] + ") write failure.")
							}
						}

						if watcher != nil {
							// add file to watch
							watcher.Add(filePath)
						}

						return nil
					},
				},
				{
					Name:    "list",
					Aliases: []string{"l"},
					Usage:   "Sync files list",
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:  "all",
							Usage: "Display all items",
						},
					},
					Action: func(cCtx *cli.Context) error {
						data, err := getData()
						if err != nil {
							color.Red(err.Error())
							return nil
						}

						privateKeyEncrypted, err := getPrivateKey()
						if err != nil {
							color.Red(err.Error())
							return nil
						}

						privateKeyHex, err := base64.RawURLEncoding.DecodeString(string(privateKeyEncrypted))
						if err != nil {
							color.Red(err.Error())
							return nil
						}

						decrypted, privateKey, err := utils.AESMACDecryptBytes(privateKeyHex, data.RsaPrivateKeyPassword)

						if err != nil || !decrypted {
							color.Red((errors.New("secret decrypt error: " + err.Error())).Error())
							return nil
						}

						isCache := false
						configs, err := listConfigs(data)
						if err != nil {
							configsCache, cacheErr := getCache()
							if cacheErr != nil {
								color.Red(err.Error())
								return err
							}
							configs = configsCache
							isCache = true
							color.Red("Request server failed, cached data is shown here: ")
						} else {
							setCache(configs)
						}
						displayRowSet := mapset.NewSet("id", "machineId", "attribute", "createdAt")
						hiddenLongPath := true
						if cCtx.Bool("all") {
							displayRowSet = mapset.NewSet("id", "attribute")
							hiddenLongPath = false
						}

						if configs != nil && configs.Size() > 0 {
							printConfigTable(configs, displayRowSet, true, hiddenLongPath, string(privateKey))
							if isCache {
								color.Red("Request server failed, above is cached data!")
							}
						} else {
							color.Red("No file config.")
						}
						return nil
					},
				},
				{
					Name:      "remove",
					Aliases:   []string{"r"},
					Usage:     "Remove a file config in sync list",
					ArgsUsage: "[fileId]",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "machineId",
							DefaultText: "current device",
							Value:       "",
							Usage:       "target machine ID",
						},
					},
					Action: func(cCtx *cli.Context) error {
						if cCtx.NArg() == 0 {
							return cli.ShowSubcommandHelp(cCtx)
						}

						data, err := getData()
						if err != nil {
							color.Red(err.Error())
							return nil
						}

						privateKeyEncrypted, err := getPrivateKey()
						if err != nil {
							color.Red(err.Error())
							return nil
						}

						privateKeyHex, err := base64.RawURLEncoding.DecodeString(string(privateKeyEncrypted))
						if err != nil {
							color.Red(err.Error())
							return nil
						}

						decrypted, privateKey, err := utils.AESMACDecryptBytes(privateKeyHex, data.RsaPrivateKeyPassword)

						if err != nil || !decrypted {
							color.Red((errors.New("secret decrypt error: " + err.Error())).Error())
							return nil
						}

						actionMachineId := data.MachineId

						if cCtx.String("machineId") != "" {
							s := cCtx.String("machineId")
							actionMachineId = ""
							actionMachineName := ""

							if s == "" {
								actionMachineId = data.MachineId
								actionMachineName = data.MachineName
							} else {
								devices, err := listDevices(data)
								if err != nil {
									color.Red(err.Error())
									return err
								}

								for i := 0; i < devices.Size(); i++ {
									machineId := devices.Get(i, "machineId").ToString()
									if strings.Contains(machineId, s) {
										actionMachineId = machineId
										encryptedActionMachineName := devices.Get(i, "machineName").ToString()
										actionMachineName, err = utils.AESCTRDecryptWithBase64(encryptedActionMachineName, []byte(data.Verify))
										if err != nil {
											color.Red(err.Error())
											return nil
										}
										break
									}
								}

								if actionMachineId == "" {
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
										actionMachineId = devices.Get(index-1, "machineId").ToString()
										actionMachineName = devices.Get(index-1, "machineName").ToString()
									} else {
										err = errors.New("invalid machine id")
										color.Red(err.Error())
										return err
									}
								}
							}

							if actionMachineId == "" {
								err = errors.New("invalid machine id")
								color.Red(err.Error())
								return err
							}

							color.Blue("Action machine: " + actionMachineName + "(" + actionMachineId + ")")
						}

						fileId := ""
						fileName := ""
						inputArg := cCtx.Args().First()

						configs, err := listConfigs(data)
						if err != nil {
							color.Red(err.Error())
						}

						for i := 0; i < configs.Size(); i++ {
							theFileId := configs.Get(i, "fileId").ToString()
							if strings.Contains(theFileId, inputArg) {
								fileId = theFileId
								fileNameEncrypted := configs.Get(i, "fileName").ToString()
								fileNameBase64, err := base64.URLEncoding.DecodeString(fileNameEncrypted)
								if err != nil {
									color.Red(err.Error())
									return nil
								}
								fileName = string(utils.RsaDecrypt([]byte(fileNameBase64), privateKey))
								break
							}
						}

						if fileId == "" {
							pattern := "\\d+"
							result, err := regexp.MatchString(pattern, inputArg)
							if err != nil {
								color.Red(err.Error())
								return err
							}

							if result {
								index, err := strconv.Atoi(inputArg)
								if err != nil {
									color.Red(err.Error())
									return err
								}
								fileId = configs.Get(index-1, "fileId").ToString()
								fileNameEncrypted := configs.Get(index-1, "fileName").ToString()
								fileNameBase64, err := base64.URLEncoding.DecodeString(fileNameEncrypted)
								if err != nil {
									color.Red(err.Error())
									return nil
								}
								fileName = string(utils.RsaDecrypt([]byte(fileNameBase64), privateKey))
							} else {
								err = errors.New("invalid file id")
								color.Red(err.Error())
								return err
							}
						}

						if fileId == "" {
							err = errors.New("invalid file id")
							color.Red(err.Error())
							return err
						}

						del := false
						promptDel := &survey.Confirm{
							Message: "Are you sure to remove the file (" + fileName + " ID:" + fileId[:10] + ") config?",
						}
						survey.AskOne(promptDel, &del)

						if !del {
							return nil
						}

						err = removeConfig(fileId, actionMachineId, data)
						if err != nil {
							color.Red("Remove config failure.")
							color.Red(err.Error())
							return nil
						}

						color.Green("Remove config successfully!")
						delCache()

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
								return s.Install()
							},
						},
						{
							Name:  "disable",
							Usage: "disable the service",
							Action: func(cCtx *cli.Context) error {
								return s.Uninstall()
							},
						},
						{
							Name:  "start",
							Usage: "start the file sync service",
							Action: func(cCtx *cli.Context) error {
								return s.Start()
							},
						},
						{
							Name:  "stop",
							Usage: "stop the file sync service",
							Action: func(cCtx *cli.Context) error {
								return s.Stop()
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
			Action: func(cCtx *cli.Context) error {
				if cCtx.NumFlags() == 0 {
					return cli.ShowAppHelp(cCtx)
				}
				return nil
			},
		}

		if err := app.Run(os.Args); err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	} else {
		logger.Info("Running under service manager.")
		// Start should not block. Do the actual work async.
		go p.run()
	}
	p.exit = make(chan struct{})

	return nil
}

func (p *program) run() error {
	data, err := getData()

	if err == nil {
		go StartSSEClient(data)
	} else {
		logger.Error(err)
	}

	logger.Infof("I'm running %v.", service.Platform())
	ticker := time.NewTicker(time.Hour)
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

	// Do't delete next line.
	isDev = true
	if isDev {
		hasEnvFile, err := utils.FileExists(".env")
		if err != nil {
			log.Fatal(err)
		}

		if !hasEnvFile {
			log.Println("No .env file found, using default values.")
		} else {
			err := godotenv.Load()
			if err != nil {
				log.Fatal("Error loading .env file")
			}
		}

	}

	err := fsInit()
	if err != nil {
		log.Fatal(err)
	}

	if isDev {
		log.Printf("Currently in development mode!")
		apiURL = os.Getenv("DEV_API_SERVER_URL")
	} else {
		conf := getConfig()
		if conf != "" {
			apiURL = conf
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

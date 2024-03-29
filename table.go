package main

import (
	"encoding/base64"
	"os"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	jsoniter "github.com/json-iterator/go"
	"github.com/yi-ge/file-sync/utils"
)

func printDeviceTable(jsonArray jsoniter.Any, displayRow mapset.Set[string], AutoMerge bool, hiddenLongPath bool, verify []byte) {
	// rowConfigAutoMerge := table.RowConfig{AutoMerge: AutoMerge}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	if jsonArray == nil {
		color.Cyan("No data found.")
		return
	}
	firstLine := jsonArray.Get(0)
	if firstLine != nil {
		keys := firstLine.Keys()
		row := make([]interface{}, len(keys)-displayRow.Cardinality())
		i := 0
		for _, v := range keys {
			if !displayRow.Contains(v) {
				row[i] = v
				i += 1
			}
		}
		t.AppendHeader(row)

		for i := 0; i < jsonArray.Size(); i++ {
			line := jsonArray.Get(i)
			col := make([]interface{}, line.Size()-displayRow.Cardinality())
			j := 0
			for _, key := range line.Keys() {
				if !displayRow.Contains(key) {
					col[j] = line.Get(key).ToString()
					if key == "machineId" || key == "fileId" {
						col[j] = line.Get(key).ToString()[:10]
					} else if key == "path" {
						value := line.Get(key).ToString()
						if hiddenLongPath && len(value) > 30 {
							col[j] = "..." + value[len(value)-30:]
						} else {
							col[j] = value
						}
					} else if key == "machineName" {
						value := line.Get(key).ToString()
						valueTmp, err := utils.AESCTRDecryptWithBase64(value, verify)
						if err != nil {
							color.Red(err.Error())
							return
						}
						col[j] = valueTmp
					}
					j += 1
				}
			}

			t.AppendRow(col)
		}

		colConfigs := []table.ColumnConfig{}

		for i := 0; i < jsonArray.Size(); i++ {
			colConfigs = append(colConfigs, table.ColumnConfig{Number: i + 1, AutoMerge: true})
		}

		t.SetColumnConfigs(colConfigs)

		t.SetAutoIndex(true)
		t.SetOutputMirror(os.Stdout)
		t.SetStyle(table.StyleLight)
		t.Style().Options.SeparateRows = true
		t.Render()
	} else {
		color.Cyan("No data found.")
	}
}

func printConfigTable(jsonArray jsoniter.Any, displayRow mapset.Set[string], AutoMerge bool, hiddenLongPath bool, privateKey string, verify []byte) {
	// rowConfigAutoMerge := table.RowConfig{AutoMerge: AutoMerge}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	if jsonArray == nil {
		color.Cyan("No data found.")
		return
	}
	firstLine := jsonArray.Get(0)
	if firstLine != nil {
		keys := firstLine.Keys()
		row := make([]interface{}, len(keys)-displayRow.Cardinality())
		i := 0
		for _, v := range keys {
			if !displayRow.Contains(v) {
				row[i] = v
				i += 1
			}
		}
		t.AppendHeader(row)

		for i := 0; i < jsonArray.Size(); i++ {
			line := jsonArray.Get(i)
			col := make([]interface{}, line.Size()-displayRow.Cardinality())
			j := 0
			for _, key := range line.Keys() {
				if !displayRow.Contains(key) {
					col[j] = line.Get(key).ToString()
					if key == "machineId" || key == "fileId" {
						col[j] = line.Get(key).ToString()[:10]
					} else if key == "path" {
						actionPathEncrypted := line.Get(key).ToString()
						actionPathBase64, err := base64.URLEncoding.DecodeString(actionPathEncrypted)
						if err != nil {
							logger.Errorf(err.Error())
							break
						}
						value := string(utils.RsaDecrypt([]byte(actionPathBase64), []byte(privateKey)))
						if hiddenLongPath && len(value) > 30 {
							col[j] = "..." + value[len(value)-30:]
						} else {
							col[j] = value
						}
					} else if key == "fileName" {
						fileNameEncrypted := line.Get(key).ToString()
						fileNameBase64, err := base64.URLEncoding.DecodeString(fileNameEncrypted)
						if err != nil {
							logger.Errorf(err.Error())
							break
						}
						value := string(utils.RsaDecrypt([]byte(fileNameBase64), []byte(privateKey)))
						if hiddenLongPath && len(value) > 30 {
							col[j] = "..." + value[len(value)-30:]
						} else {
							col[j] = value
						}
					} else if key == "machineName" {
						value := line.Get(key).ToString()
						valueTmp, err := utils.AESCTRDecryptWithBase64(value, verify)
						if err != nil {
							color.Red(err.Error())
							return
						}
						col[j] = valueTmp
					}
					j += 1
				}
			}

			t.AppendRow(col)
		}

		colConfigs := []table.ColumnConfig{}

		for i := 0; i < jsonArray.Size(); i++ {
			colConfigs = append(colConfigs, table.ColumnConfig{Number: i + 1, AutoMerge: true})
		}

		t.SetColumnConfigs(colConfigs)

		t.SetAutoIndex(true)
		t.SetOutputMirror(os.Stdout)
		t.SetStyle(table.StyleLight)
		t.Style().Options.SeparateRows = true
		t.Render()
	} else {
		color.Cyan("No data found.")
	}
}

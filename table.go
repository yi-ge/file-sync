package main

import (
	"os"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	jsoniter "github.com/json-iterator/go"
)

func printTable(jsonArray jsoniter.Any, displayRow mapset.Set[string]) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	if jsonArray == nil {
		color.Cyan("No data found.")
		return
	}
	firstLine := jsonArray.Get(0)
	if firstLine != nil {
		keys := firstLine.Keys()
		row := make([]interface{}, len(keys)-displayRow.Cardinality()+1)
		i := 1
		row[0] = "No"
		for _, v := range keys {
			if !displayRow.Contains(v) {
				row[i] = v
				i += 1
			}
		}
		t.AppendHeader(row)

		for i := 0; i < jsonArray.Size(); i++ {
			line := jsonArray.Get(i)
			col := make([]interface{}, line.Size()-displayRow.Cardinality()+1)
			j := 1
			col[0] = i + 1
			for _, key := range line.Keys() {
				if !displayRow.Contains(key) {
					col[j] = line.Get(key).ToString()
					if key == "machineId" {
						col[j] = line.Get(key).ToString()[:10]
					}
					j += 1
				}
			}

			t.AppendRow(col)
		}

		// t.AppendSeparator()
		t.SetAutoIndex(true)
		// t.SetColumnConfigs([]table.ColumnConfig{
		//     {Number: 1, AutoMerge: true},
		//     {Number: 2, AutoMerge: true},
		//     {Number: 3, AutoMerge: true},
		//     {Number: 4, AutoMerge: true},
		//     {Number: 5, Align: text.AlignCenter, AlignFooter: text.AlignCenter, AlignHeader: text.AlignCenter},
		//     {Number: 6, Align: text.AlignCenter, AlignFooter: text.AlignCenter, AlignHeader: text.AlignCenter},
		// })
		t.SetOutputMirror(os.Stdout)
		t.SetStyle(table.StyleLight)
		t.Style().Options.SeparateRows = true
		t.Render()
	} else {
		color.Cyan("No data found.")
	}
}

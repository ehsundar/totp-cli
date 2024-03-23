package cmd

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/urfave/cli/v2"

	s "github.com/yitsushi/totp-cli/internal/storage"
)

// VisualCommand is the visual live mode command
func VisualCommand() *cli.Command {
	return &cli.Command{
		Name:      "visual",
		Aliases:   []string{"vis"},
		Usage:     "View all codes in live mode",
		ArgsUsage: "",
		Flags:     []cli.Flag{},
		Action: func(ctx *cli.Context) error {
			var err error

			storage := s.NewFileStorage()
			if err = storage.Prepare(); err != nil {
				return err
			}

			app := tview.NewApplication()
			table := tview.NewTable().
				SetBorders(false)

			go func() {
				for {
					updateTableCells(storage, table)
					app.Draw()
					time.Sleep(time.Second)
				}
			}()

			table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
				if key == tcell.KeyEscape {
					app.Stop()
				}
				if key == tcell.KeyEnter {
					table.SetSelectable(true, true)
				}
			}).SetSelectedFunc(func(row int, column int) {
				table.GetCell(row, column).SetTextColor(tcell.ColorRed)
				table.SetSelectable(false, false)
			})

			if err := app.SetRoot(table, true).EnableMouse(false).Run(); err != nil {
				return err
			}

			return nil
		},
	}
}

func updateTableCells(storage *s.FileBackend, table *tview.Table) {

	table.SetCell(0, 0,
		tview.NewTableCell("Account").SetTextColor(tcell.ColorOrange),
	)

	table.SetCell(0, 1,
		tview.NewTableCell("Code").SetTextColor(tcell.ColorOrange),
	)

	table.SetCell(0, 2,
		tview.NewTableCell("Remaining").SetTextColor(tcell.ColorOrange),
	)

	rowNum := 1
	for _, namespace := range storage.ListNamespaces() {
		for _, account := range namespace.Accounts {
			code, remaining := generateCode(account)

			table.SetCell(rowNum, 0,
				tview.NewTableCell(fmt.Sprintf("%s:%s", namespace.Name, account.Name)).
					SetTextColor(tcell.ColorWhite),
			)

			table.SetCell(rowNum, 1,
				tview.NewTableCell(code).
					SetTextColor(tcell.ColorWhite),
			)

			table.SetCell(rowNum, 2,
				tview.NewTableCell(fmt.Sprintf("%ds", remaining)).
					SetTextColor(tcell.ColorRed).
					SetAlign(tview.AlignCenter),
			)

			rowNum += 1
		}
	}
}

package cmd

import (
	"fmt"
	"github.com/oustn/qtg/ui"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/oustn/qtg/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "qtg",
	Short:   "一个有趣的蜻蜓 FM 下载器",
	Version: "0.0.1",
	Args:    cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, path, err := config.ParseConfig()
		if err != nil {
			log.Fatal(err)
		}

		if cfg.Settings.EnableLogging {
			f, err := tea.LogToFile("debug.log", "debug")
			if err != nil {
				log.Fatal(err)
			}

			defer func() {
				if err = f.Close(); err != nil {
					log.Fatal(err)
				}
			}()
		}

		m := ui.NewQ(&cfg, path)
		p := tea.NewProgram(m)
		if _, err := p.Run(); err != nil {
			log.Fatal("应用打开失败", err)
		}
	},
}

// Execute runs the root command and starts the application.
func Execute() {
	rootCmd.AddCommand(updateCmd)

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// CLI app, to open the Crash Drawer from the command line.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Nimon77/cd"
	"github.com/spf13/cobra"
)

// Version is the version number, set at compile time with -ldflags "-X main.Version=1.0.0"
var Version string

func main() {
	rootCmd := &cobra.Command{
		Use:   "cd",
		Short: "cd is a BT-100U trigger for Cash Drawer",
	}

	// cmdOpen opens the cash drawer with the specified port and baud rate, or defaults
	cmdOpen := &cobra.Command{
		Use:   "open",
		Short: "Opens the cash drawer",
		Long:  `Opens the cash drawer.`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			port, err := cmd.Flags().GetString("port")
			if err != nil {
				log.Fatal(err)
			}
			baud, err := cmd.Flags().GetInt("baud")
			if err != nil {
				log.Fatal(err)
			}

			drawer, err := cd.New(port, baud)
			if err != nil {
				log.Fatal(err)
			}
			defer drawer.Close()

			err = drawer.Open(ctx)
			if err != nil {
				log.Fatal(err)
			}

			log.Println("Cash drawer opened successfully!")
		},
	}

	// cmdVersion prints the version number
	cmdVersion := &cobra.Command{
		Use:   "version",
		Short: "Print version number",
		Long:  `Print the version number of cd.`,
		Run: func(cmd *cobra.Command, args []string) {
			if Version == "" {
				fmt.Println("unknown version")
			} else {
				fmt.Println("cd " + Version)
			}
		},
	}

	rootCmd.AddCommand(cmdOpen)
	rootCmd.AddCommand(cmdVersion)
	cmdOpen.PersistentFlags().StringP("port", "p", "/dev/ttyUSB0", "TTY port")
	cmdOpen.PersistentFlags().IntP("baud", "b", 9600, "baud rate")

	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

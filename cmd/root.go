/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/anmol1vw13/grep/tool"
	"github.com/spf13/cobra"
)

var flagSet tool.FlagOptions
// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "grep",
	Short: "A CLI based search tool",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) { 
		grep := tool.GrepProps{
			Flags: flagSet,
			Args: args,
		}
		result := grep.Search()
		if result.Err != nil {
			fmt.Println(result.Err)
		} else {
			for _,line := range result.Lines {
				fmt.Println(line)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetUsageFunc(nil)
	rootCmd.SetUsageTemplate("Grep tool")
	rootCmd.Flags().StringVarP(&flagSet.OutputFile, "output", "o", "", "File to output data to")
	rootCmd.Flags().BoolVarP(&flagSet.CaseInsensitive, "ignoreCase", "i", false, "If true it performs a case insensitive search")

}



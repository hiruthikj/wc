/*
Copyright © 2024 Hiruthik J

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	countBytes bool
	countLines bool
	countWords bool
	countChars bool
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var rootCmd = &cobra.Command{
	Use:   "ccwc",
	Short: "wc - print newline, word, and byte counts for each file",
	Long: `Print newline, word, and byte counts for each FILE, and a total line if more than one FILE is specified.
A word is a non-zero-length sequence of characters delimited by white space.

With no FILE, or when FILE is -, read standard input.`,
	Run: wcRunnerFn,
}

func wcRunnerFn(cmd *cobra.Command, args []string) {
	fi, _ := os.Stdin.Stat() // get the FileInfo struct describing the standard input.
	var reader io.Reader
	var err error
	var filePath string

	if (fi.Mode() & os.ModeCharDevice) == 0 {
		fmt.Println("Data is from pipe")

		reader = os.Stdin
		// bytes, _ := io.ReadAll(os.Stdin)
		// str := string(bytes)
		// fmt.Println(str)
	} else {
		if len(args) != 1 {
			cmd.Help()
			os.Exit(1)
		}
		filePath = args[0]
		reader, err = os.Open(filePath)
		check(err)
	}


	var outputSb strings.Builder

	cmd.Flags().Visit(func(f *pflag.Flag) {
		switch f.Name {
		case "bytes":
			var count int64
			if countBytes {
				if (fi.Mode() & os.ModeCharDevice) == 0 {
					count, err = getCountBytes(reader)
				} else {
					count, err = getFileSizeBytes(filePath)
				}
				check(err)

				outputSb.WriteString(fmt.Sprint(count))
				outputSb.WriteString(" ")
			}

		case "lines":
			if countLines {
				count, err := lineCounter(reader)
				check(err)

				outputSb.WriteString(fmt.Sprint(count))
				outputSb.WriteString(" ")
			}

		case "words":
			if countWords {
				r := bufio.NewReader(reader)
				count, err := wordCounter(r)
				check(err)

				outputSb.WriteString(fmt.Sprint(count))
				outputSb.WriteString(" ")
			}

		case "chars":
			if countChars {
				r := bufio.NewReader(reader)
				count, err := runeCounter(r)
				check(err)

				outputSb.WriteString(fmt.Sprint(count))
				outputSb.WriteString(" ")
			}

		}
	})

	outputSb.WriteString(filePath)
	fmt.Println(outputSb.String())
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ccwc.yaml)")

	rootCmd.Flags().BoolVarP(&countBytes, "bytes", "c", false, "print the byte counts")
	rootCmd.Flags().BoolVarP(&countLines, "lines", "l", false, "print the newline counts")
	rootCmd.Flags().BoolVarP(&countWords, "words", "w", false, "print the word counts")
	rootCmd.Flags().BoolVarP(&countChars, "chars", "m", false, "print the character counts")

	rootCmd.Flags().SortFlags = false

	// rootCmd.MarkFlagFilename("bytes")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".ccwc" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".ccwc")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

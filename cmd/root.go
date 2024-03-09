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
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	filePath string
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
	isSet := cmd.Flags().Lookup("bytes").Changed

	var outputSb strings.Builder

	if isSet {
		count, err := getFileSizeBytes(filePath)
		check(err)
		outputSb.WriteString(fmt.Sprint(count))
		outputSb.WriteString(" ")
	}

	isSet = cmd.Flags().Lookup("lines").Changed
	if isSet {
		f, err := os.Open(filePath)
		check(err)
		count, err := lineCounter(f)
		check(err)
		outputSb.WriteString(fmt.Sprint(count))

		outputSb.WriteString(" ")
	}

	isSet = cmd.Flags().Lookup("words").Changed
	if isSet {
		f, err := os.Open(filePath)
		check(err)
		r := bufio.NewReader(f)
		count, err := wordCounter(r)
		check(err)
		outputSb.WriteString(fmt.Sprint(count))

		outputSb.WriteString(" ")
	}

	isSet = cmd.Flags().Lookup("chars").Changed
	if isSet {
		f, err := os.Open(filePath)
		check(err)
		r := bufio.NewReader(f)
		count, err := runeCounter(r)
		check(err)
		outputSb.WriteString(fmt.Sprint(count))

		outputSb.WriteString(" ")
	}

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

	rootCmd.Flags().StringVarP(&filePath, "bytes", "c", filePath, "print the byte counts")
	rootCmd.Flags().StringVarP(&filePath, "lines", "l", filePath, "print the newline counts")
	rootCmd.Flags().StringVarP(&filePath, "words", "w", filePath, "print the word counts")
	rootCmd.Flags().StringVarP(&filePath, "chars", "m", filePath, "print the character counts")

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

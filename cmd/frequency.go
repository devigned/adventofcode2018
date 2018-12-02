package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

func init() {
	freqCmd.PersistentFlags().StringVarP(&freqParams.InputPath, "input-path", "i", "", "path to input file")
	if err := freqCmd.MarkPersistentFlagRequired("input-path"); err != nil {
		log.Fatal(err)
	}
	freqCmd.AddCommand(freqSkewCmd)
	freqCmd.AddCommand(freqRepeatCmd)
	rootCmd.AddCommand(freqCmd)
}

type (
	FreqParams struct {
		InputPath string
	}
)

var (
	freqParams  FreqParams

	freqCmd = &cobra.Command{
		Use: "freq",
		Short: "Frequency related commands",
	}

	freqSkewCmd = &cobra.Command{
		Use:   "skew",
		Short: "Calculate the frequency skew from an input file",
		Run: func(cmd *cobra.Command, args []string) {
			file, err := os.Open(freqParams.InputPath)
			if err != nil {
				log.Fatal(err)
			}
			defer func() {
				if err := file.Close(); err != nil {
					log.Fatal(err)
				}
			}()

			scanner := bufio.NewScanner(file)
			total := 0
			for scanner.Scan() {
				i, err := strconv.Atoi(scanner.Text())
				if err != nil {
					log.Fatal(err)
					return
				}
				total += i
			}
			fmt.Println("freq skew: ", total)

			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
			return
		},
	}

	freqRepeatCmd = &cobra.Command{
		Use:   "repeat",
		Short: "Find first repeat frequency reached",
		Run: func(cmd *cobra.Command, args []string) {
			file, err := os.Open(freqParams.InputPath)
			if err != nil {
				log.Fatal(err)
			}
			defer func() {
				if err := file.Close(); err != nil {
					log.Fatal(err)
				}
			}()

			var items []int
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				i, err := strconv.Atoi(scanner.Text())
				if err != nil {
					log.Fatal(err)
					return
				}
				items = append(items, i)
			}

			trackedTotals := make(map[int]bool)
			total := 0
			// run until a loop is found
			for i := 0;;i++{
				idx := i % len(items)
				total += items[idx]
				if _, ok := trackedTotals[total]; ok {
					fmt.Println("first repeat: ", total)
					return
				}
				trackedTotals[total] = true
			}
		},
	}
)

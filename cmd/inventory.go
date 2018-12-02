package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	inventoryCmd.PersistentFlags().StringVarP(&inventoryParams.InputPath, "input-path", "i", "", "path to input file")
	if err := inventoryCmd.MarkPersistentFlagRequired("input-path"); err != nil {
		log.Fatal(err)
	}
	inventoryCmd.AddCommand(checksumCmd)
	inventoryCmd.AddCommand(commonCmd)
	rootCmd.AddCommand(inventoryCmd)
}

type (
	InventoryParams struct {
		InputPath string
	}
)

var (
	inventoryParams InventoryParams

	inventoryCmd = &cobra.Command{
		Use:   "inventory",
		Short: "Inventory related commands",
	}

	checksumCmd = &cobra.Command{
		Use:   "checksum",
		Short: "Find the checksum of exactly 2 and 3 letter repeats",
		Run: func(cmd *cobra.Command, args []string) {
			file, err := os.Open(inventoryParams.InputPath)
			if err != nil {
				log.Fatal(err)
			}
			defer func() {
				if err := file.Close(); err != nil {
					log.Fatal(err)
				}
			}()

			scanner := bufio.NewScanner(file)
			two := 0
			three := 0
			for scanner.Scan() {
				letterCounts := make(map[rune]int)
				for _, r := range scanner.Text() {
					if _, ok := letterCounts[r]; ok {
						letterCounts[r]++
					} else {
						letterCounts[r] = 1
					}
				}
				hasTwo, hasThree := false, false
				for _, val := range letterCounts {
					if val == 2 {
						hasTwo = true
					} else if val == 3 {
						hasThree = true
					}
				}

				if hasTwo {
					two++
				}

				if hasThree {
					three++
				}
			}
			fmt.Printf("checksum: %d * %d = %d", two, three, two*three)
		},
	}

	commonCmd = &cobra.Command{
		Use:   "common",
		Short: "Find the common runes shared between two boxes with only 1 difference",
		Run: func(cmd *cobra.Command, args []string) {
			file, err := os.Open(inventoryParams.InputPath)
			if err != nil {
				log.Fatal(err)
			}
			defer func() {
				if err := file.Close(); err != nil {
					log.Fatal(err)
				}
			}()

			var ids []string
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				ids = append(ids, scanner.Text())
			}

			for _, id := range ids {
				for _, id2 := range ids {
					differences := 0
					for idx := range id {
						if idx >= len(id2) || id[idx] != id2[idx] {
							differences++
						}
					}
					if differences == 1 {
						// winner winner, chicken dinner
						var sharedRunes []rune
						for idx, r := range id {
							if id[idx] == id2[idx] {
								sharedRunes = append(sharedRunes, r)
							}
						}
						fmt.Println("shared runes: ", string(sharedRunes))
						return
					}
				}
			}
			letterCounts := make(map[rune]int)
			for _, r := range scanner.Text() {
				if _, ok := letterCounts[r]; ok {
					letterCounts[r]++
				} else {
					letterCounts[r] = 1
				}
			}

		},
	}
)

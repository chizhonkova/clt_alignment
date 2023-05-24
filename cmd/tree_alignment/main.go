package main

import (
	"fmt"
	"os"
	alignment "tree_alignment/internal"

	"github.com/spf13/cobra"
)

var (
	flagFirstGraphPath    string
	flagSecondGraphPath   string
	flagResultPath        string
	flagTagEqualityCost   int
	flagTagUnequalityCost int
	flagDeletionCost      int
)

var rootCmd = &cobra.Command{
	Use:   "tree-alignment",
	Short: "calculates alignment for two binary trees with tags",
	Run: func(cmd *cobra.Command, args []string) {
		c := &alignment.Config{
			FirstGraphDescription:  LoadGraphDescription(flagFirstGraphPath),
			SecondGraphDescription: LoadGraphDescription(flagSecondGraphPath),
			ResultPath:             flagResultPath,
			TagEqualityCost:        flagTagEqualityCost,
			TagUnequalityCost:      flagTagUnequalityCost,
			DeletionCost:           flagDeletionCost,
		}

		if err := alignment.CalculateAlignment(c); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%+v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagFirstGraphPath, "first-graph-path", "", "path to the first graph description")
	rootCmd.PersistentFlags().StringVar(&flagSecondGraphPath, "second-graph-path", "", "path to the second graph description")
	rootCmd.PersistentFlags().StringVar(&flagResultPath, "result-path", "", "path to the resulting graph description")
	rootCmd.PersistentFlags().IntVar(&flagTagEqualityCost, "tag-equality-cost", 4, "cost for equal tags")
	rootCmd.PersistentFlags().IntVar(&flagTagUnequalityCost, "tag-unequality-cost", 3, "penalty for unequal tags")
	rootCmd.PersistentFlags().IntVar(&flagDeletionCost, "deletion-cost", 2, "penalty for deletion")
	if err := rootCmd.MarkPersistentFlagRequired("first-graph-path"); err != nil {
		panic(err)
	}
	if err := rootCmd.MarkPersistentFlagRequired("second-graph-path"); err != nil {
		panic(err)
	}
	if err := rootCmd.MarkPersistentFlagRequired("result-path"); err != nil {
		panic(err)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
}

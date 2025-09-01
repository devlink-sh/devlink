package pair

import "github.com/spf13/cobra"

var PairCmd = &cobra.Command{
	Use:   "pair",
	Short: "Now you can share your localhost with people haha",
}

func init() {
	PairCmd.AddCommand(pairShareCmd)
	PairCmd.AddCommand(pairGetCmd)
}

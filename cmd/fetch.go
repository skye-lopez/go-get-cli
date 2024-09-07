package cmd

import (
	"github.com/skye-lopez/go-get-cli/data"
	"github.com/spf13/cobra"
	"github.com/syndtr/goleveldb/leveldb"
)

var fetchCommand = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch latest index data",
	Long:  `Fetch latest index data, this can take awhile the first time you run it.`,
	Run:   Init,
}

func init() {
	rootCmd.AddCommand(fetchCommand)
}

func Init(cmd *cobra.Command, args []string) {
	db, err := leveldb.OpenFile(".go-get-cli/data", nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	data.Init(db)
}

package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/dikeert/linkman/links"
	"github.com/spf13/cobra"
)

// archiveCmd represents the archive command
var archiveCmd = &cobra.Command{
	Use:   "archive id [other ids]",
	Short: "archives a link",
	Long: `arhives link with specified ID.

Example:

linkman archive id
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		store := openStore(dataPath)

		for _, idArg := range args {
			if id, err := strconv.Atoi(idArg); err == nil {
				archiveLink(store, id)
			} else {
				fmt.Fprintf(os.Stderr, "Value %s is not an ID", idArg)
			}
		}
	},
}

func openStore(path string) links.Store {
	store, err := links.OpenStore(path)
	if err == nil {
		return store
	}

	die("Unable to open store", err)
	panic("shouldn't get there")
}

func archiveLink(store links.Store, id int) {
	if err := store.ArchiveByID(id); err != nil {
		die("Unable to archive link", err)
	}
}

func init() {
	rootCmd.AddCommand(archiveCmd)
}

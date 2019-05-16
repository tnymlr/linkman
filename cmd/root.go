package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "linkman",
	Short: "Manage links (bookmarks)",
	Long: `linkman allows you to create and list links.

A link is a bookmark of URL you want to see later.
A link consist of:
 - URL
 - Source - second (or third) level domain name
 - Title - title of the page referenced by the URL
 - List - a list the link belongs to
 - Archived status - whether to consider link to be archived

linkman is capable of maintaining multiple lists with links.
By default in adds and lists links from 'default' list

You create new links by invoking 'add' command an supplying a URL to it.
linkman will figure out source, fetch title and add the link
into 'default' list into database in $XDG_DATA_HOME/linkman/data.db.

Later you can invoke 'list' command to print links that you've saved.

Examples on how source is calculated:

| url               | source        |
| ----              | ------        |
| youtube.com       | youtube       |
| stackoverflow.com | stackoverflow |
| domain.co.uk      | domain        |
`,
	//	Run: func(cmd *cobra.Command, args []string) { },
}

var dataPath string

//Execute is the entry point into the application.
//It configures and starts the execute of root command
//which in turn passes the execution to underying commands.
func Execute(path string, args []string) {
	dataPath = path
	rootCmd.SetArgs(args)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func die(msg string, err error) {
	fmt.Fprintln(os.Stderr, fmt.Errorf("%s: %s", msg, err))
	os.Exit(1)
}

func init() {
}

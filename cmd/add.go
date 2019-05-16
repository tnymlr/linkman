package cmd

import (
	"fmt"
	"net/url"

	"github.com/dikeert/linkman/links"
	"github.com/dikeert/linkman/pages"
	"github.com/dikeert/linkman/urls"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add url [another url]",
	Short: "Creates new link for URL",
	Long: `'add' will calculate source, fetch title and
create new link for supplied URLs

You can provide list you want to use to store URL,
by default add will use 'default' list.

By default it does not allow to create links for URLs that already
had links created for them.
`,
	Args: cobra.MinimumNArgs(1),
	Run:  runAdd,
}

var skipFetchingTitle = false
var allowDuplicates = false
var targetList = "default"
var providedTitle = ""

func runAdd(cmd *cobra.Command, args []string) {
	if store, err := links.OpenStore(dataPath); err == nil {
		for _, rawurl := range args {
			saveRawURL(store, rawurl)
		}
	} else {
		die("Unable to open links store", err)
	}
}

func saveRawURL(store links.Store, rawurl string) {
	url := parseURL(rawurl)
	if allowDuplicates || !urlExists(store, url) {
		saveURL(store, url)
	} else {
		fmt.Printf("URL %s already exists, skipping\n", rawurl)
	}
}

func parseURL(rawurl string) *url.URL {
	url, err := urls.ParseURL(rawurl)
	if err == nil {
		return url
	}

	die("Unable to add URL", err)
	panic("shouldn't get here")
}

func urlExists(store links.Store, url *url.URL) bool {
	exists, err := store.LinkExists(url)
	if err == nil {
		return exists
	}

	die("Unexpected error", err)
	panic("shouldn't get there")
}

func saveURL(store links.Store, url *url.URL) {
	source := getSource(url)
	title := fetchTitle(url)
	link := store.NewLink(url, source, title, targetList)
	if err := store.SaveLink(link); err == nil {
		fmt.Println("Create link: ")
		fmt.Printf("  URL: %s\n", link.URL)
		fmt.Printf("  Title: %s\n", link.Title)
		fmt.Printf("  List: %s\n", link.List)
	} else {
		die("Unable to save link", err)
	}
}

func getSource(url *url.URL) string {
	source, err := urls.GetSource(url)
	if err == nil {
		return source
	}

	die("Mailformed URL", err)
	panic("shouldn't get here")
}

func fetchTitle(url *url.URL) string {
	if skipFetchingTitle {
		return providedTitle
	} else if title, err := pages.FetchTitle(url); err == nil {
		return title
	} else {
		die("Unable to fetch page title", err)
	}

	panic("shouldn't get there")
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().BoolVarP(&skipFetchingTitle, "skip-title-fetch", "", false,
		"Skip fetching title")
	addCmd.Flags().StringVarP(&providedTitle, "title", "t", "",
		"Use provided title instead of fetching")
	addCmd.Flags().BoolVarP(&allowDuplicates, "force", "f", false,
		"Allow duplicates")
	addCmd.Flags().StringVarP(&targetList, "list", "l", "default", "Target list")
}

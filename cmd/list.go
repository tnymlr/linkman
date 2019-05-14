package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/dikeert/linkman/links"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Prints stored links",
	Long: `'list' finds and prints links that has been added before.
It allows to specify apply certain filtering criterias:
 - by source
 - by title
 - by archived status
 - by list

By default it prints all non-archived links that belong to
'default' list.

You can specify output format for links. Available fields:

 - ID: unique identificator of the link
 - Source: source of the link
 - Title: title of the page referenced by the link
 - URL: URL of the link

Default output format:

ID:	{{.ID}}
Source:	{{.Source}}
Title:	{{.Title}}
URL:	{{.URL}}

Examples:

linkman list -s source - prints links that belong to that source
linkman list -l mylist - prints links from 'mylist'
linkman list -l '*' - prints links from all lists
linkman list -T - prints only links which has non-empty title
linkman list -t title - prints links that have 'title' in the title

linkman list -f '{{.ID}}:\t{{.Source}}' - prints links as
list of "id: source" lines
`,
	Run: runList,
}

const DEFAULT_TEMPLATE = `
ID:	{{.ID}}
Source:	{{.Source}}
Title:	{{.Title}}
URL:	{{.URL}}
`

var format string = DEFAULT_TEMPLATE
var source string = ""
var list string = ""
var title string = ""

var requireTitle bool = false
var archived bool = false
var onlyArchived bool = false

func runList(cmd *cobra.Command, args []string) {
	writer := getOutputWriter()
	template := getOutputTemplate()
	store := openLinksStore(dataPath)

	for _, link := range getLinks(store) {
		printLink(writer, template, link)
	}
	writer.Flush()
}

func openLinksStore(path string) links.Store {
	if store, err := links.OpenStore(path); err == nil {
		return store
	} else {
		die("Unable to open links store", err)
	}

	panic("Shouldn't get there")
}

func getOutputWriter() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
}

func getLinks(store links.Store) []links.Link {
	var conds []links.FilterCondition

	if source != "" {
		conds = append(conds, links.WithSource(source))
	}

	if title != "" {
		conds = append(conds, links.WithTitle(title))
	}

	if list != "" {
		conds = append(conds, links.FromList(list))
	}

	if requireTitle {
		conds = append(conds, links.TitleNotEmpty())
	}

	if onlyArchived {
		conds = append(conds, links.OnlyArchived())
	} else if archived {
		conds = append(conds, links.IncludeArchived())
	} else {
		conds = append(conds, links.NoArchived())
	}

	filter := links.NewFilter(conds...)

	if links, err := store.FindLinks(filter); err == nil {
		return links
	} else {
		die("Unable to fetch links", err)
	}

	panic("Shouldn't get there")
}

func getOutputTemplate() *template.Template {
	tpl := template.New("output template")
	format := unescapeOutputTemplate(format)

	if tpl, err := tpl.Parse(format); err == nil {
		return tpl
	} else {
		die("Unable to parse output template", err)
	}

	panic("Shouldn't get there")
}

func unescapeOutputTemplate(format string) string {
	quoted := strconv.Quote(format)
	replaced := strings.Replace(quoted, `\\`, "\\", -1)

	if result, err := strconv.Unquote(replaced); err == nil {
		return result
	} else {
		die("Unable to parse output template", err)
	}

	panic("Shouldn't get there")
}

func printLink(writer *tabwriter.Writer,
	tpl *template.Template,
	link links.Link) {

	if err := tpl.Execute(writer, link); err != nil {
		fmt.Fprintf(os.Stderr, "Error pringing link: %s\n", err)
	}
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&format,
		"format", "f",
		DEFAULT_TEMPLATE,
		"Output template. Available fields are: ID, URL, Source, Title,List")

	listCmd.Flags().StringVarP(&source,
		"source", "s", "",
		"Show only link from specified source")

	listCmd.Flags().StringVarP(&list,
		"list", "l", "default",
		"Show only links from specified list")

	listCmd.Flags().StringVarP(&title,
		"title", "t", "",
		"Show only links which title contains specified string")

	listCmd.Flags().BoolVarP(&requireTitle,
		"require-title", "T", false,
		"When specified filters out links without title")

	listCmd.Flags().BoolVarP(&archived,
		"archived", "a", false,
		"Include archived links")

	listCmd.Flags().BoolVarP(&onlyArchived,
		"only-archived", "A", false,
		"Show only archived links")
}

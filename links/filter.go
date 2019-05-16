package links

//FilterCondition is a function that modifies
//LinkFilter to filter out certain links.
type FilterCondition func(*linkFilter) *linkFilter

//LinkFilter represents a filtering configuration for links.
//Such a filtering configuration is a result of combination
//of several FilterConditions
type LinkFilter interface {
	hasSource() bool
	hasTitle() bool
	titleNotEmpty() bool

	getSource() string
	getTitle() string
	getList() string

	getArchivedFlag() archivedFlag
}

//NewFilter creates new filter with specified conditions.
func NewFilter(conditions ...FilterCondition) LinkFilter {
	filter := &linkFilter{}

	for _, condition := range conditions {
		filter = condition(filter)
	}

	return filter
}

//WithSource creates new filtering condition for Source field.
//This filtering conditions allows only links which
//Source field matches provided source string.
func WithSource(source string) FilterCondition {
	return func(filter *linkFilter) *linkFilter {
		filter.source = source
		return filter
	}
}

//WithTitle creates new filtering condition for Title field.
//This filtering condition allows only links which
//Title field contains provided title string.
func WithTitle(title string) FilterCondition {
	return func(filter *linkFilter) *linkFilter {
		filter.title = title
		return filter
	}
}

//FromList creates new filtering condition for List field.
//This filtering condition allows only links which
//List field matches provided filter string.
func FromList(list string) FilterCondition {
	return func(filter *linkFilter) *linkFilter {
		filter.list = list
		return filter
	}
}

//TitleNotEmpty creates new filtering condition for Title field.
//This filtering condition allows only links which
//Title field is not empty.
func TitleNotEmpty() FilterCondition {
	return func(filter *linkFilter) *linkFilter {
		filter.requireTitle = true
		return filter
	}
}

//IncludeArchived creates new filtering condition for Archived field.
//This filtering condition allows links which have either true or false
//value written into Archived field.
func IncludeArchived() FilterCondition {
	return archivedBuilder(includeArchived)
}

//OnlyArchived creates new filtering condition for Archived field.
//This filtering condition only allows links that have true
//value written into Archived field.
func OnlyArchived() FilterCondition {
	return archivedBuilder(onlyArchived)
}

//NoArchived creates new filtering condition for Archived field.
//This filtering condition only allows links that have false
//value written into Archived field.
func NoArchived() FilterCondition {
	return archivedBuilder(noArchived)
}

type archivedFlag int

const ( // archived flags
	includeArchived = iota + 1
	onlyArchived
	noArchived
)

type linkFilter struct {
	source       string
	title        string
	list         string
	archived     archivedFlag
	requireTitle bool
}

func (me *linkFilter) hasSource() bool {
	return me.source != ""
}

func (me *linkFilter) hasTitle() bool {
	return me.title != ""
}

func (me *linkFilter) titleNotEmpty() bool {
	return me.requireTitle
}

func (me *linkFilter) getSource() string {
	return me.source
}

func (me *linkFilter) getTitle() string {
	return me.title
}

func (me *linkFilter) getList() string {
	if me.list == "" {
		return "default"
	}

	return me.list
}

func (me *linkFilter) getArchivedFlag() archivedFlag {
	return me.archived
}

func archivedBuilder(flag archivedFlag) FilterCondition {
	return func(filter *linkFilter) *linkFilter {
		if filter.archived == 0 {
			filter.archived = flag
			return filter
		}

		panic("archived filter can be applied only once")
	}
}

package links

type FilterCondition func(*linkFilter) *linkFilter

type LinkFilter interface {
	hasSource() bool
	hasTitle() bool
	titleNotEmpty() bool

	getSource() string
	getTitle() string
	getList() string

	getArchivedFlag() archivedFlag
}

func NewFilter(conditions ...FilterCondition) LinkFilter {
	filter := &linkFilter{}

	for _, condition := range conditions {
		filter = condition(filter)
	}

	return filter
}

func WithSource(source string) FilterCondition {
	return func(filter *linkFilter) *linkFilter {
		filter.source = source
		return filter
	}
}

func WithTitle(title string) FilterCondition {
	return func(filter *linkFilter) *linkFilter {
		filter.title = title
		return filter
	}
}

func FromList(list string) FilterCondition {
	return func(filter *linkFilter) *linkFilter {
		filter.list = list
		return filter
	}
}

func TitleNotEmpty() FilterCondition {
	return func(filter *linkFilter) *linkFilter {
		filter.requireTitle = true
		return filter
	}
}

func IncludeArchived() FilterCondition {
	return archivedBuilder(include_archived)
}

func OnlyArchived() FilterCondition {
	return archivedBuilder(only_archived)
}

func NoArchived() FilterCondition {
	return archivedBuilder(no_archived)
}

type archivedFlag int

const ( // archived flags
	include_archived = iota + 1
	only_archived
	no_archived
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
	} else {
		return me.list
	}
}

func (me *linkFilter) getArchivedFlag() archivedFlag {
	return me.archived
}

func archivedBuilder(flag archivedFlag) FilterCondition {
	return func(filter *linkFilter) *linkFilter {
		if filter.archived == 0 {
			filter.archived = flag
			return filter
		} else {
			panic("archived filter can be applied only once")
		}
	}
}

package links

import (
	"fmt"
	"net/url"

	"github.com/dikeert/linkman/db"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
)

//Link holds all the data associated with stored URL in the database.
type Link struct {
	ID       int `storm:"id,increment"`
	URL      *url.URL
	Source   string `storm:"index"`
	Title    string
	List     string `storm:"index"`
	Archived bool
}

//Store provides access to storage of Links.
type Store interface {
	NewLink(url *url.URL, source string, title string, list string) *Link
	SaveLink(link *Link) error
	LinkExists(url *url.URL) (bool, error)
	FindLinks(LinkFilter) ([]Link, error)
	ArchiveByID(id int) error
}

//OpenStore creates new Store for database located
//at provided path.
func OpenStore(path string) (Store, error) {
	db, err := db.Open(path)
	if err != nil {
		return nil, err
	}

	defer db.Close()
	if err := initDatabase(db); err != nil {
		return nil, err
	}

	return &storeImpl{path: path}, nil
}

type storeImpl struct {
	path string
}

func (me *storeImpl) NewLink(url *url.URL, source string, title string, list string) *Link {
	if url == nil {
		panic("url is nil")
	}

	return &Link{
		URL:      url,
		Source:   source,
		Title:    title,
		List:     list,
		Archived: false,
	}

}

func (me *storeImpl) SaveLink(link *Link) error {
	db, err := db.Open(me.path)
	if err != nil {
		return fmt.Errorf("Unable to open database: %s", err)
	}

	defer db.Close()
	return save(db, link)
}

func (me *storeImpl) LinkExists(url *url.URL) (bool, error) {
	db, err := db.Open(me.path)
	if err != nil {
		return false, err
	}

	defer db.Close()
	links, err := findLinksByURL(db, url)
	return len(links) > 0, err
}

func (me *storeImpl) FindLinks(filter LinkFilter) ([]Link, error) {
	db, err := db.Open(me.path)
	if err != nil {
		return nil, err
	}

	return findLinks(db, filter)
}

//ArchiveByID archived the links with specified id.
func (me *storeImpl) ArchiveByID(id int) error {
	db, err := db.Open(me.path)
	if err != nil {
		return err
	}

	defer db.Close()
	return archiveByID(db, id)
}

func initDatabase(db *storm.DB) error {
	err := db.Init(&Link{})
	if err != nil {
		return err
	}

	return nil
}

func findLinks(db *storm.DB, filter LinkFilter) ([]Link, error) {
	var matchers []q.Matcher
	var result []Link

	if filter.hasSource() {
		matchers = append(matchers,
			q.Eq("Source", filter.getSource()))
	}

	if filter.hasTitle() {
		matchers = append(matchers,
			q.Re("Title", fmt.Sprintf("^.+%s.+$", filter.getTitle())))
	}

	if filter.titleNotEmpty() {
		matchers = append(matchers,
			q.Re("Title", "^.+$"))
	}

	if list := filter.getList(); list != "*" {
		matchers = append(matchers, q.Eq("List", list))
	}

	if filter.getArchivedFlag() == onlyArchived {
		matchers = append(matchers, q.Eq("Archived", true))
	} else if filter.getArchivedFlag() == noArchived {
		matchers = append(matchers, q.Eq("Archived", false))
	}

	if err := db.Select(matchers...).Find(&result); err == nil {
		return result, nil
	} else if err == storm.ErrNotFound {
		return result, nil
	} else {
		return nil, err
	}
}

func findLinksByURL(db *storm.DB, url *url.URL) ([]Link, error) {
	var links []Link
	err := db.Find("URL", url, &links)

	if err == storm.ErrNotFound {
		err = nil
	}

	return links, err
}

func archiveByID(db *storm.DB, id int) error {
	err := db.UpdateField(&Link{ID: id}, "Archived", true)
	if err != nil {
		return err
	}

	return nil
}

func save(db *storm.DB, link *Link) error {
	err := db.Save(link)
	if err != nil {
		return fmt.Errorf("Unable to save link: %s", err)
	}

	return nil
}

package links

import (
	"fmt"
	"net/url"

	"github.com/dikeert/linkman/db"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
)

type Link struct {
	ID       int `storm:"id,increment"`
	URL      *url.URL
	Source   string `storm:"index"`
	Title    string
	List     string `storm:"index"`
	Archived bool
}

type Store interface {
	NewLink(url *url.URL, source string, title string, list string) *Link
	SaveLink(link *Link) error
	LinkExists(url *url.URL) (bool, error)
	FindLinks(LinkFilter) ([]Link, error)
	ArchiveById(id int) error
}

func OpenStore(path string) (Store, error) {
	if db, err := db.Open(path); err == nil {
		defer db.Close()
		if err := initDatabase(db); err == nil {
			return &storeImpl{path: path}, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
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
	if db, err := db.Open(me.path); err == nil {
		defer db.Close()
		return save(db, link)
	} else {
		return fmt.Errorf("Unable to open database: %s", err)
	}
}

func (me *storeImpl) LinkExists(url *url.URL) (bool, error) {
	if db, err := db.Open(me.path); err == nil {
		defer db.Close()
		links, err := findLinksByUrl(db, url)
		return len(links) > 0, err
	} else {
		return false, err
	}
}

func (me *storeImpl) FindLinks(filter LinkFilter) ([]Link, error) {
	if db, err := db.Open(me.path); err == nil {
		return findLinks(db, filter)
	} else {
		return nil, err
	}
}

func (me *storeImpl) ArchiveById(id int) error {
	if db, err := db.Open(me.path); err == nil {
		defer db.Close()
		return archiveById(db, id)
	} else {
		return err
	}
}

func initDatabase(db *storm.DB) error {
	if err := db.Init(&Link{}); err != nil {
		return err
	} else {
		return nil
	}
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

	if filter.getArchivedFlag() == only_archived {
		matchers = append(matchers, q.Eq("Archived", true))
	} else if filter.getArchivedFlag() == no_archived {
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

func findLinksByUrl(db *storm.DB, url *url.URL) ([]Link, error) {
	var links []Link
	err := db.Find("URL", url, &links)

	if err == storm.ErrNotFound {
		err = nil
	}

	return links, err
}

func archiveById(db *storm.DB, id int) error {
	if err := db.UpdateField(&Link{ID: id}, "Archived", true); err != nil {
		if err == storm.ErrNotFound {
			return nil
		} else {
			return err
		}
	}

	return nil
}

func save(db *storm.DB, link *Link) error {
	if err := db.Save(link); err != nil {
		return fmt.Errorf("Unable to save link: %s", err)
	} else {
		return nil
	}
}

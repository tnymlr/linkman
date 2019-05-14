package cmd_test

import (
	"io/ioutil"
	"os"
	"reflect"
	"runtime"
	"testing"

	"github.com/dikeert/linkman/cmd"
	"github.com/dikeert/linkman/links"

	"github.com/stretchr/testify/assert"
)

type testCase func(string, links.Store, *testing.T)

var url string = "https://www.wikipedia.org/"

func TestCommands(t *testing.T) {
	tests := []testCase{
		SimpleAdd,
		AddWithoutTitle,
		AddWithCustomTitle,
		AddNoDuplcatesByDef,
		AddForceDuplicate,
	}

	for _, tc := range tests {
		tc := tc
		t.Run(getName(tc), func(t *testing.T) {
			path := getDataFile()
			defer os.Remove(path)

			if store, err := links.OpenStore(path); err == nil {
				tc(path, store, t)
			} else {
				t.Error(err)
			}
		})
	}
}

func SimpleAdd(path string, store links.Store, t *testing.T) {
	assert := assert.New(t)
	title := "Wikipedia"

	cmd.Execute(path, []string{
		"add",
		url,
	})

	if links, err := store.FindAllLinks(); err == nil {
		assert.Equal(1, len(links), "should create a link")
		assert.Equal(links[0].URL.String(), url,
			"should populate link with supplied url")
		assert.Equal(links[0].Title, title,
			"Should fetch page title")
	} else {
		t.Error(err)
	}
}

func AddWithoutTitle(path string, store links.Store, t *testing.T) {
	assert := assert.New(t)

	cmd.Execute(path, []string{
		"add",
		url,
		"--skip-title-fetch",
	})

	if links, err := store.FindAllLinks(); err == nil {
		assert.Equal(len(links), 1, "Should create a link")
		assert.Equal(links[0].URL.String(), url, "Should populare url")
		assert.Empty(links[0].Title)
	} else {
		t.Error(err)
	}
}

func AddWithCustomTitle(path string, store links.Store, t *testing.T) {
	assert := assert.New(t)
	title := "Custom Title"

	cmd.Execute(path, []string{
		"add",
		url,
		"--skip-title-fetch",
		"-t", title,
	})

	if links, err := store.FindAllLinks(); err == nil {
		assert.Equal(len(links), 1, "Should create a link")
		assert.Equal(links[0].URL.String(), url, "Should populare ULR")
		assert.Equal(links[0].Title, title, "Should use provided title")
	}
}

func AddNoDuplcatesByDef(path string, store links.Store, t *testing.T) {
	assert := assert.New(t)

	cmd.Execute(path, []string{
		"add",
		url,
		"--skip-title-fetch", //make it faster
	})

	cmd.Execute(path, []string{
		"add",
		url,
		"--skip-title-fetch", //make it faster
	})

	if links, err := store.FindAllLinks(); err == nil {
		assert.Equal(1, len(links), "Should create one link")
	} else {
		t.Error(err)
	}
}

func AddForceDuplicate(path string, store links.Store, t *testing.T) {
	assert := assert.New(t)

	cmd.Execute(path, []string{
		"add",
		url,
		"--skip-title-fetch", //make it faster
		"-f",
	})

	cmd.Execute(path, []string{
		"add",
		url,
		"--skip-title-fetch", //make it faster
		"--force",
	})

	if links, err := store.FindAllLinks(); err == nil {
		assert.Equal(2, len(links), "Should create two links")
	} else {
		t.Error(err)
	}
}

func getName(c testCase) string {
	return runtime.FuncForPC(reflect.ValueOf(c).Pointer()).Name()
}

func getDataFile() string {
	tmpfile, err := ioutil.TempFile("", "linkman.*.db")
	if err != nil {
		panic(err)
	}
	defer tmpfile.Close()
	return tmpfile.Name()
}

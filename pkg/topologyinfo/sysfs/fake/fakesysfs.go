package fake

import (
	"io/ioutil"
	"os"
	"strings"
)

type DebugLogger func(string, ...interface{})

func NullDebugLogger(fmt string, args ...interface{}) {
	return
}

var DebugLog DebugLogger = NullDebugLogger

type Tree interface {
	Add(name string, attrs map[string]string) Tree
	Name() string
	Items() []Tree
	SetAttrs() error
	Attrs() map[string]string
	Create() error
}

type Creator interface {
	Create(Tree) error
}

func NewTree(name string, attrs map[string]string) Tree {
	return &tree{
		name:  name,
		attrs: attrs,
		items: []Tree{},
	}
}

type tree struct {
	name  string
	attrs map[string]string
	items []Tree
}

func (t *tree) Add(name string, attrs map[string]string) Tree {
	n := NewTree(name, attrs)
	t.items = append(t.items, n)
	return n
}

func (t *tree) Items() []Tree {
	return t.items
}

func (t *tree) Name() string {
	return t.name
}

func (t *tree) Create() error {
	return newCreator().Create(t)
}

func (t *tree) Attrs() map[string]string {
	res := make(map[string]string)
	for key, val := range t.attrs {
		res[key] = val
	}
	return res
}

func (t *tree) SetAttrs() error {
	if t.attrs == nil {
		DebugLog("%q attrs NONE", t.name)
		return nil
	}
	var err error
	for name, content := range t.attrs {
		DebugLog("%q attrs %q", t.name, name)
		err = ioutil.WriteFile(name, []byte(content), 0644)
		if err != nil {
			break
		}
	}
	return err
}

func MakeAttrs(attrs map[string]string) map[string]string {
	resAttrs := make(map[string]string)
	for key, value := range attrs {
		if strings.HasSuffix(value, "\n") {
			resAttrs[key] = value
		} else {
			resAttrs[key] = value + "\n"
		}
	}
	return resAttrs
}

type creator struct{}

func newCreator() Creator {
	return &creator{}
}

func (c *creator) Create(t Tree) error {
	DebugLog("%q attrs: %v", t.Name(), t.Attrs())
	if err := t.SetAttrs(); err != nil {
		return err
	}
	DebugLog("%q items: %v", t.Name(), t.Items())
	return c.createItems(t.Items())
}

func (c *creator) createItems(t []Tree) error {
	var err error
	for _, st := range t {
		err = os.Mkdir(st.Name(), 0755)
		if err != nil {
			break
		}
		err = c.createItem(st)
	}
	return err
}

func (c *creator) createItem(st Tree) error {
	os.Chdir(st.Name())
	defer os.Chdir("..")
	return st.Create()
}

type FakeSysfs struct {
	base string
	root Tree
}

func NewFakeSysfs(base string) (*FakeSysfs, error) {
	return &FakeSysfs{
		base: base,
		// DO NOT USE NEITHER "." or "" HERE!!
		root: NewTree("_", nil),
	}, nil
}

func (fs *FakeSysfs) AddTree(entries ...string) Tree {
	DebugLog("%q adding: %v", fs.base, entries)
	pos := fs.root
	for _, entry := range entries {
		pos = pos.Add(entry, nil)
	}
	return pos
}

func (fs *FakeSysfs) Base() string {
	return fs.base
}

func (fs *FakeSysfs) Root() Tree {
	return fs.root
}

func (fs *FakeSysfs) Setup() error {
	oldWd, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(oldWd)
	os.Chdir(fs.base)
	DebugLog("Setup(%q)", fs.base)
	return fs.root.Create()
}

func (fs *FakeSysfs) Teardown() error {
	return os.RemoveAll(fs.base)
}

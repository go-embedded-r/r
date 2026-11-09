package eval

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

// The file system a session sees. By default it is a private, in-memory tree
// (nothing of the host is reachable, and its size is bounded by
// Options.MaxFileBytes); Options.FileRoot gives a session a directory of the
// host as its root instead, confined to it with os.Root (".." and symbolic
// links cannot leave it). Paths are slash-separated, absolute and cleaned
// before they get here.

var (
	errIsDir    = errors.New("Is a directory")
	errNotDir   = errors.New("Not a directory")
	errNotEmpty = errors.New("Directory not empty")
	errTooBig   = errors.New("File too large")
	errPerm     = errors.New("Permission denied")
	errInvalid  = errors.New("Invalid argument")
)

// errText is the strerror text R shows for a failure.
func errText(err error) string {
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return "No such file or directory"
	case errors.Is(err, fs.ErrExist):
		return "File exists"
	case errors.Is(err, errIsDir):
		return "Is a directory"
	case errors.Is(err, errNotDir):
		return "Not a directory"
	case errors.Is(err, errNotEmpty):
		return "Directory not empty"
	case errors.Is(err, errTooBig):
		return "File too large"
	case errors.Is(err, errInvalid):
		return "Invalid argument"
	}
	return "Permission denied"
}

type fsStat struct {
	dir   bool
	size  int64
	mtime time.Time
}

type fileSystem interface {
	Read(p string) ([]byte, error)
	Write(p string, data []byte, appendTo bool) error
	Stat(p string) (fsStat, error)
	Mkdir(p string) error
	Remove(p string) error
	List(p string) ([]string, error)
	Rename(from, to string) error
}

// ---- in memory ------------------------------------------------------------

type memNode struct {
	dir   bool
	data  []byte
	mtime time.Time
}

type memFS struct {
	nodes map[string]*memNode
	used  int64
	max   int64 // 0: unbounded
}

const maxMemNodes = 1 << 16

func newMemFS(max int64) *memFS {
	return &memFS{nodes: map[string]*memNode{"/": {dir: true, mtime: time.Now()}}, max: max}
}

func (m *memFS) parentOK(p string) error {
	d := path.Dir(p)
	n, ok := m.nodes[d]
	switch {
	case !ok:
		return fs.ErrNotExist
	case !n.dir:
		return errNotDir
	}
	return nil
}

func (m *memFS) Read(p string) ([]byte, error) {
	n, ok := m.nodes[p]
	switch {
	case !ok:
		return nil, fs.ErrNotExist
	case n.dir:
		return nil, errIsDir
	}
	return n.data, nil
}

func (m *memFS) Write(p string, data []byte, appendTo bool) error {
	n, ok := m.nodes[p]
	if ok && n.dir {
		return errIsDir
	}
	if !ok {
		if err := m.parentOK(p); err != nil {
			return err
		}
		if len(m.nodes) >= maxMemNodes {
			return errTooBig
		}
		n = &memNode{}
	}
	old := int64(len(n.data))
	nw := int64(len(data))
	if appendTo {
		nw += old
	}
	if m.max > 0 && m.used-old+nw > m.max {
		return errTooBig
	}
	if appendTo {
		// bytes already stored are never modified, only appended to, so a
		// slice a reader holds stays valid
		n.data = append(n.data, data...)
	} else {
		n.data = append([]byte(nil), data...)
	}
	m.used += nw - old
	n.mtime = time.Now()
	m.nodes[p] = n
	return nil
}

func (m *memFS) Stat(p string) (fsStat, error) {
	n, ok := m.nodes[p]
	if !ok {
		return fsStat{}, fs.ErrNotExist
	}
	return fsStat{dir: n.dir, size: int64(len(n.data)), mtime: n.mtime}, nil
}

func (m *memFS) Mkdir(p string) error {
	if _, ok := m.nodes[p]; ok {
		return fs.ErrExist
	}
	if err := m.parentOK(p); err != nil {
		return err
	}
	if len(m.nodes) >= maxMemNodes {
		return errTooBig
	}
	m.nodes[p] = &memNode{dir: true, mtime: time.Now()}
	return nil
}

func (m *memFS) Remove(p string) error {
	n, ok := m.nodes[p]
	if !ok {
		return fs.ErrNotExist
	}
	if p == "/" {
		return errPerm
	}
	if n.dir {
		pre := p + "/"
		for k := range m.nodes {
			if strings.HasPrefix(k, pre) {
				return errNotEmpty
			}
		}
	}
	m.used -= int64(len(n.data))
	delete(m.nodes, p)
	return nil
}

func (m *memFS) List(p string) ([]string, error) {
	n, ok := m.nodes[p]
	switch {
	case !ok:
		return nil, fs.ErrNotExist
	case !n.dir:
		return nil, errNotDir
	}
	pre := p
	if pre != "/" {
		pre += "/"
	}
	var out []string
	for k := range m.nodes {
		if k != "/" && strings.HasPrefix(k, pre) && !strings.Contains(k[len(pre):], "/") {
			out = append(out, k[len(pre):])
		}
	}
	sort.Strings(out)
	return out, nil
}

func (m *memFS) Rename(from, to string) error {
	n, ok := m.nodes[from]
	if !ok {
		return fs.ErrNotExist
	}
	if from == "/" || to == "/" {
		return errPerm
	}
	if to == from {
		return nil
	}
	if err := m.parentOK(to); err != nil {
		return err
	}
	if t, ok := m.nodes[to]; ok {
		if t.dir != n.dir {
			if t.dir {
				return errIsDir
			}
			return errNotDir
		}
		if err := m.Remove(to); err != nil {
			return err
		}
	}
	// a directory takes its subtree along
	pre := from + "/"
	moved := map[string]*memNode{}
	for k, v := range m.nodes {
		if strings.HasPrefix(k, pre) {
			moved[to+"/"+k[len(pre):]] = v
			delete(m.nodes, k)
		}
	}
	delete(m.nodes, from)
	m.nodes[to] = n
	for k, v := range moved {
		m.nodes[k] = v
	}
	return nil
}

// ---- a directory of the host, confined ------------------------------------

type rootFS struct {
	root *os.Root
	max  int64
}

func (r *rootFS) rel(p string) string {
	p = strings.TrimPrefix(p, "/")
	if p == "" {
		return "."
	}
	return p
}

func mapOSErr(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, fs.ErrNotExist):
		return fs.ErrNotExist
	case errors.Is(err, fs.ErrExist):
		return fs.ErrExist
	}
	msg := err.Error()
	switch {
	case strings.Contains(msg, "is a directory"):
		return errIsDir
	case strings.Contains(msg, "not a directory"):
		return errNotDir
	case strings.Contains(msg, "not empty"):
		return errNotEmpty
	}
	return errPerm
}

func (r *rootFS) Read(p string) ([]byte, error) {
	st, err := r.root.Stat(r.rel(p))
	if err != nil {
		return nil, mapOSErr(err)
	}
	if st.IsDir() {
		return nil, errIsDir
	}
	if r.max > 0 && st.Size() > r.max {
		return nil, errTooBig
	}
	f, err := r.root.Open(r.rel(p))
	if err != nil {
		return nil, mapOSErr(err)
	}
	defer f.Close()
	lr := io.Reader(f)
	if r.max > 0 {
		lr = io.LimitReader(f, r.max+1)
	}
	b, err := io.ReadAll(lr)
	if err != nil {
		return nil, mapOSErr(err)
	}
	if r.max > 0 && int64(len(b)) > r.max {
		return nil, errTooBig
	}
	return b, nil
}

func (r *rootFS) Write(p string, data []byte, appendTo bool) error {
	flags := os.O_WRONLY | os.O_CREATE
	if appendTo {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}
	if r.max > 0 {
		old := int64(0)
		if st, err := r.root.Stat(r.rel(p)); err == nil && appendTo {
			old = st.Size()
		}
		if old+int64(len(data)) > r.max {
			return errTooBig
		}
	}
	f, err := r.root.OpenFile(r.rel(p), flags, 0o666)
	if err != nil {
		return mapOSErr(err)
	}
	_, werr := f.Write(data)
	cerr := f.Close()
	if werr != nil {
		return mapOSErr(werr)
	}
	return mapOSErr(cerr)
}

func (r *rootFS) Stat(p string) (fsStat, error) {
	st, err := r.root.Stat(r.rel(p))
	if err != nil {
		return fsStat{}, mapOSErr(err)
	}
	return fsStat{dir: st.IsDir(), size: st.Size(), mtime: st.ModTime()}, nil
}

func (r *rootFS) Mkdir(p string) error { return mapOSErr(r.root.Mkdir(r.rel(p), 0o777)) }

func (r *rootFS) Remove(p string) error {
	if p == "/" {
		return errPerm
	}
	return mapOSErr(r.root.Remove(r.rel(p)))
}

func (r *rootFS) List(p string) ([]string, error) {
	ents, err := fs.ReadDir(r.root.FS(), r.rel(p))
	if err != nil {
		return nil, mapOSErr(err)
	}
	out := make([]string, len(ents))
	for i, e := range ents {
		out[i] = e.Name()
	}
	sort.Strings(out)
	return out, nil
}

func (r *rootFS) Rename(from, to string) error {
	if from == "/" || to == "/" {
		return errPerm
	}
	return mapOSErr(r.root.Rename(r.rel(from), r.rel(to)))
}

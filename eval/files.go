package eval

import (
	"fmt"
	"math/rand/v2"
	"path"
	"sort"
	"strings"
	"time"
)

// File-system functions over the session's file system (see fs.go), the
// session's environment variables, and Sys.sleep. Written from the help pages
// of files, file.exists, file.copy, dir.create, list.files, unlink, tempfile,
// normalizePath, getwd, Sys.getenv and Sys.setenv, and from the oracle.

func (s *Session) stat(p string) (fsStat, error) {
	if p == "/dev/null" {
		return fsStat{}, nil
	}
	return s.io.fs.Stat(p)
}

// pathVec is a character argument naming files; absent is allowed when nil.
func (s *Session) pathVec(c *callInfo, v Value, what string) []Str {
	st, ok := v.(*String)
	if !ok {
		s.errorf(c.call, "invalid '%s' argument", what)
	}
	return st.V
}

// dotPaths gathers the character arguments of a function taking `...` of
// file names, leaving out the named ones the function knows.
func (s *Session) dotPaths(c *callInfo, what string, known ...string) (paths []Str, named map[string]Value) {
	named = map[string]Value{}
next:
	for i, a := range c.args {
		for _, k := range known {
			if c.names[i] == k {
				named[k] = a
				continue next
			}
		}
		paths = append(paths, s.pathVec(c, a, what)...)
	}
	return paths, named
}

func logicals(bs []bool) *Logical {
	out := &Logical{V: make([]int32, len(bs))}
	for i, b := range bs {
		if b {
			out.V[i] = 1
		}
	}
	return out
}

func (s *Session) tempDir() string {
	if s.io.tmp != "" {
		if st, err := s.io.fs.Stat(s.io.tmp); err == nil && st.dir {
			return s.io.tmp
		}
	}
	const alnum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for try := 0; ; try++ {
		b := make([]byte, 6)
		for i := range b {
			b[i] = alnum[rand.IntN(len(alnum))]
		}
		p := "/tmp/Rtmp" + string(b)
		_ = s.io.fs.Mkdir("/tmp")
		if err := s.io.fs.Mkdir(p); err == nil {
			s.io.tmp = p
			return p
		} else if try > 20 {
			s.errorf(s.currentCall(), "cannot create 'R_TempDir'")
		}
	}
}

func recycle3(n ...int) int {
	m := 0
	for _, k := range n {
		if k == 0 {
			return 0
		}
		m = max(m, k)
	}
	return m
}

// glob matches name against a shell pattern (*, ?, [...]).
func glob(pat, name string) bool {
	ok, err := path.Match(pat, name)
	return err == nil && ok
}

func hasWild(p string) bool { return strings.ContainsAny(p, "*?[") }

// expandGlob lists the existing paths that match p, component by component.
func (s *Session) expandGlob(p string) []string {
	abs := s.abspath(p)
	cur := []string{"/"}
	for _, comp := range strings.Split(strings.TrimPrefix(abs, "/"), "/") {
		if comp == "" {
			continue
		}
		var next []string
		for _, d := range cur {
			if !hasWild(comp) {
				next = append(next, strings.TrimSuffix(d, "/")+"/"+comp)
				continue
			}
			names, err := s.io.fs.List(d)
			if err != nil {
				continue
			}
			for _, n := range names {
				if (n[0] != '.' || comp[0] == '.') && glob(comp, n) {
					next = append(next, strings.TrimSuffix(d, "/")+"/"+n)
				}
			}
		}
		cur = next
	}
	return cur
}

func (s *Session) removeAll(p string) bool {
	st, err := s.io.fs.Stat(p)
	if err != nil {
		return true
	}
	ok := true
	if st.dir {
		names, _ := s.io.fs.List(p)
		for _, n := range names {
			if !s.removeAll(strings.TrimSuffix(p, "/") + "/" + n) {
				ok = false
			}
		}
	}
	if s.io.fs.Remove(p) != nil {
		ok = false
	}
	return ok
}

// walk lists the entries below dir (relative names), depth first in sorted
// order; dirs says whether directories are listed themselves.
func (s *Session) walk(dir, prefix string, all, recursive, dirs bool, visit func(rel string, isDir bool)) {
	names, err := s.io.fs.List(dir)
	if err != nil {
		return
	}
	for _, n := range names {
		if !all && n[0] == '.' {
			continue
		}
		st, err := s.io.fs.Stat(strings.TrimSuffix(dir, "/") + "/" + n)
		if err != nil {
			continue
		}
		if st.dir && recursive {
			if dirs {
				visit(prefix+n, true)
			}
			s.walk(strings.TrimSuffix(dir, "/")+"/"+n, prefix+n+"/", all, recursive, dirs, visit)
			continue
		}
		visit(prefix+n, st.dir)
	}
}

func init() {
	regBuiltin("file.exists", func(s *Session, c *callInfo) Value {
		if len(c.args) == 0 {
			s.errorf(c.call, "invalid '%s' argument", "file")
		}
		paths, _ := s.dotPaths(c, "file")
		out := make([]bool, len(paths))
		for i, p := range paths {
			if !p.NA {
				_, err := s.stat(s.abspath(p.S))
				out[i] = err == nil
			}
		}
		return logicals(out)
	})
	regBuiltin("dir.exists", func(s *Session, c *callInfo) Value {
		st, ok := s.nargs1(c, "dir.exists").(*String)
		if !ok {
			s.errorf(c.call, "invalid filename argument")
		}
		paths := st.V
		out := make([]bool, len(paths))
		for i, p := range paths {
			if !p.NA {
				st, err := s.stat(s.abspath(p.S))
				out[i] = err == nil && st.dir
			}
		}
		return logicals(out)
	})
	regBuiltin("file.remove", func(s *Session, c *callInfo) Value {
		paths, _ := s.dotPaths(c, "file")
		out := make([]bool, len(paths))
		for i, p := range paths {
			if p.NA {
				continue
			}
			if err := s.io.fs.Remove(s.abspath(p.S)); err != nil {
				s.warningf(c.call, "cannot remove file '%s', reason '%s'", p.S, errText(err))
				continue
			}
			out[i] = true
		}
		return logicals(out)
	})
	regBuiltin("unlink", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "recursive", "force", "expand")
		if a[0] == nil {
			s.errorf(c.call, "argument \"x\" is missing, with no default")
		}
		var paths []Str
		switch x := a[0].(type) {
		case *String:
			paths = x.V
		case nullT:
		default:
			s.errorf(c.call, "invalid '%s' argument", "x")
		}
		recursive := s.boolArg(c, a[1], false, "recursive")
		expand := s.boolArg(c, a[3], true, "expand")
		status := int32(0)
		for _, p := range paths {
			if p.NA {
				continue
			}
			targets := []string{s.abspath(p.S)}
			if expand && hasWild(p.S) {
				targets = s.expandGlob(p.S)
			}
			for _, t := range targets {
				st, err := s.io.fs.Stat(t)
				switch {
				case err != nil:
				case st.dir && !recursive:
					status = 1
				case !s.removeAll(t):
					status = 1
				}
			}
		}
		return s.invisible(mkInt(status))
	})
	regBuiltin("file.create", func(s *Session, c *callInfo) Value {
		paths, named := s.dotPaths(c, "file", "showWarnings")
		show := s.boolArg(c, named["showWarnings"], true, "showWarnings")
		out := make([]bool, len(paths))
		for i, p := range paths {
			if p.NA {
				continue
			}
			if err := s.io.fs.Write(s.abspath(p.S), nil, false); err != nil {
				if show {
					s.warningf(c.call, "cannot create file '%s', reason '%s'", p.S, errText(err))
				}
				continue
			}
			out[i] = true
		}
		return logicals(out)
	})
	regBuiltin("file.rename", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "from", "to")
		from, to := s.pathVec(c, a[0], "from"), s.pathVec(c, a[1], "to")
		if len(from) != len(to) {
			s.errorf(c.call, "'from' and 'to' are of different lengths")
		}
		out := make([]bool, len(from))
		for i := range from {
			if from[i].NA || to[i].NA {
				continue
			}
			if err := s.io.fs.Rename(s.abspath(from[i].S), s.abspath(to[i].S)); err != nil {
				s.warningf(c.call, "cannot rename file '%s' to '%s', reason '%s'", from[i].S, to[i].S, errText(err))
				continue
			}
			out[i] = true
		}
		return logicals(out)
	})
	regBuiltin("file.append", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "file1", "file2")
		f1, f2 := s.pathVec(c, a[0], "file1"), s.pathVec(c, a[1], "file2")
		n := recycle3(len(f1), len(f2))
		out := make([]bool, n)
		for i := 0; i < n; i++ {
			p1, p2 := f1[i%len(f1)], f2[i%len(f2)]
			if p1.NA || p2.NA {
				continue
			}
			data, err := s.io.fs.Read(s.abspath(p2.S))
			if err != nil {
				continue
			}
			out[i] = s.io.fs.Write(s.abspath(p1.S), data, true) == nil
		}
		return logicals(out)
	})
	regBuiltin("file.copy", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "from", "to", "overwrite", "recursive", "copy.mode", "copy.date")
		from, to := s.pathVec(c, a[0], "from"), s.pathVec(c, a[1], "to")
		recursive := s.boolArg(c, a[3], false, "recursive")
		overwrite := s.boolArg(c, a[2], recursive, "overwrite")
		if recursive {
			s.errorf(c.call, "'recursive = TRUE' is not supported by this evaluator (copying directory trees)")
		}
		intoDir := false
		if len(to) == 1 && !to[0].NA {
			if st, err := s.stat(s.abspath(to[0].S)); err == nil && st.dir {
				intoDir = true
			}
		}
		if !intoDir && len(to) != len(from) {
			s.errorf(c.call, "more 'from' files than 'to' files")
		}
		out := make([]bool, len(from))
		for i, f := range from {
			if f.NA {
				continue
			}
			var dest string
			if intoDir {
				dest = s.abspath(to[0].S) + "/" + f.S[strings.LastIndexByte(f.S, '/')+1:]
			} else if to[i].NA {
				continue
			} else {
				dest = s.abspath(to[i].S)
			}
			data, err := s.io.fs.Read(s.abspath(f.S))
			if err != nil {
				continue
			}
			if _, err := s.io.fs.Stat(dest); err == nil && !overwrite {
				continue
			}
			out[i] = s.io.fs.Write(dest, data, false) == nil
		}
		return logicals(out)
	})
	regBuiltin("file.size", func(s *Session, c *callInfo) Value {
		paths, _ := s.dotPaths(c, "file")
		out := make([]float64, len(paths))
		for i, p := range paths {
			out[i] = naReal
			if p.NA {
				continue
			}
			if st, err := s.stat(s.abspath(p.S)); err == nil {
				out[i] = float64(st.size)
				if st.dir {
					out[i] = 4096
				}
			}
		}
		return &Double{V: out}
	})
	regBuiltin("dir.create", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "path", "showWarnings", "recursive", "mode")
		p := s.stringArg(c, a[0], "path")
		show := s.boolArg(c, a[1], true, "showWarnings")
		recursive := s.boolArg(c, a[2], false, "recursive")
		abs := s.abspath(p)
		if st, err := s.io.fs.Stat(abs); err == nil {
			if show {
				s.warningf(c.call, "'%s' already exists", p)
			}
			_ = st
			return s.invisible(mkLogical(false))
		}
		if recursive {
			var chain []string
			for d := abs; d != "/"; d = path.Dir(d) {
				chain = append(chain, d)
			}
			for i := len(chain) - 1; i >= 0; i-- {
				if _, err := s.io.fs.Stat(chain[i]); err != nil {
					if err := s.io.fs.Mkdir(chain[i]); err != nil {
						if show {
							s.warningf(c.call, "cannot create dir '%s', reason '%s'", p, errText(err))
						}
						return s.invisible(mkLogical(false))
					}
				}
			}
			return s.invisible(mkLogical(true))
		}
		if err := s.io.fs.Mkdir(abs); err != nil {
			if show {
				s.warningf(c.call, "cannot create dir '%s', reason '%s'", p, errText(err))
			}
			return s.invisible(mkLogical(false))
		}
		return s.invisible(mkLogical(true))
	})
	listFiles := func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "path", "pattern", "all.files", "full.names", "recursive", "ignore.case", "include.dirs", "no..", "fixed")
		paths := []Str{{S: "."}}
		if a[0] != nil {
			paths = s.pathVec(c, a[0], "path")
		}
		var m *matcher
		if a[1] != nil {
			if _, null := a[1].(nullT); !null {
				pat := s.stringArg(c, a[1], "pattern")
				m = s.compileRegex(c, pat, nil, s.boolArg(c, a[8], false, "fixed"), s.boolArg(c, a[5], false, "ignore.case"), false)
			}
		}
		all := s.boolArg(c, a[2], false, "all.files")
		full := s.boolArg(c, a[3], false, "full.names")
		rec := s.boolArg(c, a[4], false, "recursive")
		incDirs := s.boolArg(c, a[6], false, "include.dirs")
		noDots := s.boolArg(c, a[7], false, "no..")
		var out []string
		for _, p := range paths {
			if p.NA {
				continue
			}
			var found []string
			if all && !rec && !noDots {
				if st, err := s.stat(s.abspath(p.S)); err == nil && st.dir {
					found = append(found, ".", "..")
				}
			}
			s.walk(s.abspath(p.S), "", all, rec, incDirs || !rec, func(rel string, isDir bool) {
				if isDir && rec && !incDirs {
					return
				}
				found = append(found, rel)
			})
			var keep []string
			for _, f := range found {
				base := f[strings.LastIndexByte(f, '/')+1:]
				if m != nil && !m.matches(base) {
					continue
				}
				keep = append(keep, f)
			}
			sort.SliceStable(keep, func(i, j int) bool { return collate(keep[i], keep[j]) < 0 })
			for _, f := range keep {
				if full {
					f = strings.TrimSuffix(p.S, "/") + "/" + f
				}
				out = append(out, f)
			}
		}
		return strVec(out)
	}
	regBuiltin("list.files", listFiles)
	regBuiltin("dir", listFiles)
	regBuiltin("list.dirs", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "path", "full.names", "recursive")
		paths := []Str{{S: "."}}
		if a[0] != nil {
			paths = s.pathVec(c, a[0], "path")
		}
		full := s.boolArg(c, a[1], true, "full.names")
		rec := s.boolArg(c, a[2], true, "recursive")
		var out []string
		for _, p := range paths {
			if p.NA {
				continue
			}
			st, err := s.stat(s.abspath(p.S))
			if err != nil || !st.dir {
				continue
			}
			var found []string
			if rec {
				found = append(found, "")
			}
			s.walk(s.abspath(p.S), "", true, rec, true, func(rel string, isDir bool) {
				if isDir {
					found = append(found, rel)
				}
			})
			sort.SliceStable(found, func(i, j int) bool { return collate(found[i], found[j]) < 0 })
			for _, f := range found {
				switch {
				case full && f == "":
					out = append(out, p.S)
				case full:
					out = append(out, strings.TrimSuffix(p.S, "/")+"/"+f)
				default:
					out = append(out, f)
				}
			}
		}
		return strVec(out)
	})
	regBuiltin("tempdir", func(s *Session, c *callInfo) Value { return mkString(s.tempDir()) })
	regBuiltin("tempfile", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "pattern", "tmpdir", "fileext")
		get := func(v Value, def string, what string) []Str {
			if v == nil {
				return []Str{{S: def}}
			}
			return s.pathVec(c, v, what)
		}
		pat := get(a[0], "file", "pattern")
		dir := get(a[1], "", "tmpdir")
		ext := get(a[2], "", "fileext")
		if a[1] == nil {
			dir = []Str{{S: s.tempDir()}}
		}
		n := recycle3(len(pat), len(dir), len(ext))
		out := make([]string, n)
		for i := range out {
			for {
				name := fmt.Sprintf("%s/%s%x%s", dir[i%len(dir)].S, pat[i%len(pat)].S, rand.Uint64()&0xffffffffffff, ext[i%len(ext)].S)
				if _, err := s.io.fs.Stat(s.abspath(name)); err != nil {
					out[i] = name
					break
				}
			}
		}
		return strVec(out)
	})
	regBuiltin("getwd", func(s *Session, c *callInfo) Value { return mkString(s.io.cwd) })
	regBuiltin("setwd", func(s *Session, c *callInfo) Value {
		p := s.stringArg(c, s.nargs1(c, "setwd"), "dir")
		abs := s.abspath(p)
		if st, err := s.stat(abs); err != nil || !st.dir {
			s.errorf(c.call, "cannot change working directory")
		}
		old := s.io.cwd
		s.io.cwd = abs
		return s.invisible(mkString(old))
	})
	regBuiltin("normalizePath", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "path", "winslash", "mustWork")
		paths := s.pathVec(c, a[0], "path")
		mode := 0 // NA: warn
		if a[2] != nil {
			if l, ok := a[2].(*Logical); ok && len(l.V) == 1 && l.V[0] != naInt {
				mode = 1 + int(1-l.V[0]) // TRUE: 1 (error), FALSE: 2 (silent)
			}
		}
		out := make([]Str, len(paths))
		for i, p := range paths {
			out[i] = p
			if p.NA {
				continue
			}
			abs := s.abspath(p.S)
			if _, err := s.stat(abs); err == nil {
				out[i] = Str{S: abs}
				continue
			}
			msg := fmt.Sprintf("path[%d]=\"%s\": No such file or directory", i+1, p.S)
			switch mode {
			case 0:
				s.warningf(c.call, "%s", msg)
			case 1:
				s.errorf(c.call, "%s", msg)
			}
		}
		return &String{V: out}
	})
	regBuiltin("path.expand", func(s *Session, c *callInfo) Value {
		st := s.pathVec(c, s.nargs1(c, "path.expand"), "path")
		return &String{V: append([]Str(nil), st...)}
	})
	regBuiltin("Sys.getenv", func(s *Session, c *callInfo) Value {
		a := s.matchFormals(c, "x", "unset", "names")
		unset := Str{S: ""}
		if a[1] != nil {
			switch u := a[1].(type) {
			case *String:
				if len(u.V) != 1 {
					s.errorf(c.call, "invalid '%s' argument", "unset")
				}
				unset = u.V[0]
			case *Logical:
				if len(u.V) == 1 && u.V[0] == naInt {
					unset = Str{NA: true}
				} else {
					s.errorf(c.call, "invalid '%s' argument", "unset")
				}
			default:
				s.errorf(c.call, "invalid '%s' argument", "unset")
			}
		}
		var keys []string
		var vals []Str
		all := a[0] == nil
		if _, null := a[0].(nullT); null {
			all = true
		}
		if all {
			for k := range s.io.env {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				vals = append(vals, Str{S: s.io.env[k]})
			}
		} else {
			x, ok := a[0].(*String)
			if !ok {
				s.errorf(c.call, "wrong type for argument")
			}
			for _, k := range x.V {
				keys = append(keys, k.S)
				if v, ok := s.io.env[k.S]; ok && !k.NA {
					vals = append(vals, Str{S: v})
				} else {
					vals = append(vals, unset)
				}
			}
		}
		named := all || len(keys) > 1
		if a[2] != nil {
			if l, ok := a[2].(*Logical); ok && len(l.V) == 1 && l.V[0] != naInt {
				named = l.V[0] != 0
			}
		}
		out := &String{V: vals}
		if named {
			out.Attr = namesAttr(strsOf(keys))
		}
		return out
	})
	regBuiltin("Sys.setenv", func(s *Session, c *callInfo) Value {
		out := make([]bool, len(c.args))
		for i, v := range c.args {
			if c.names[i] == "" {
				s.errorf(c.call, "all arguments must be named")
			}
			st := asStrings(v)
			if len(st) != 1 {
				s.errorf(c.call, "'names' and 'val' are of different lengths")
			}
			if st[0].NA {
				st[0] = Str{S: "NA"}
			}
			s.io.env[c.names[i]] = st[0].S
			out[i] = true
		}
		return s.invisible(logicals(out))
	})
	regBuiltin("Sys.unsetenv", func(s *Session, c *callInfo) Value {
		x, ok := s.nargs1(c, "Sys.unsetenv").(*String)
		if !ok {
			s.errorf(c.call, "wrong type for argument")
		}
		out := make([]bool, len(x.V))
		for i, k := range x.V {
			delete(s.io.env, k.S)
			out[i] = true
		}
		return s.invisible(logicals(out))
	})
	regBuiltin("Sys.sleep", func(s *Session, c *callInfo) Value {
		d := asDoublesOrNil(s.nargs1(c, "Sys.sleep"))
		if len(d) != 1 || d[0] != d[0] || d[0] < 0 {
			s.errorf(c.call, "invalid '%s' value", "time")
		}
		end := time.Now().Add(time.Duration(min(d[0], 1e9) * float64(time.Second)))
		for {
			s.checkInterrupt()
			left := time.Until(end)
			if left <= 0 {
				break
			}
			time.Sleep(min(left, 5*time.Millisecond))
		}
		return s.invisible(Null)
	})
}

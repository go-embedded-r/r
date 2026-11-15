package eval

// file.info and its relatives ask the session's file system, in memory or
// confined to Options.FileRoot, for what it holds: the size, whether it is a
// directory and its modification time. The data frame is built in R
// (baselib/files16.R).

func init() {
	regBuiltin("..fileStat", func(s *Session, c *callInfo) Value {
		if len(c.args) != 1 {
			s.errorf(c.call, "%d arguments passed to '..fileStat' which requires 1", len(c.args))
		}
		paths := s.pathVec(c, c.args[0], "fn")
		size := make([]float64, len(paths))
		isDir := make([]int32, len(paths))
		mtime := make([]float64, len(paths))
		for i, p := range paths {
			size[i], isDir[i], mtime[i] = naReal, naInt, naReal
			if p.NA {
				continue
			}
			st, err := s.stat(s.abspath(p.S))
			if err != nil {
				continue
			}
			size[i] = float64(st.size)
			if st.dir {
				size[i] = 4096
				isDir[i] = 1
			} else {
				isDir[i] = 0
			}
			if !st.mtime.IsZero() {
				mtime[i] = float64(st.mtime.UnixNano()) / 1e9
			}
		}
		return listOf([]Value{&Double{V: size}, &Logical{V: isDir}, &Double{V: mtime}},
			&String{V: []Str{{S: "size"}, {S: "isdir"}, {S: "mtime"}}})
	})
}

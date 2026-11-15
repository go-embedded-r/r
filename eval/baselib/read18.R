# Base library, slice 18 (1): the readers. Written from the help pages of
# R 4.6.2 (type.convert, read.table, count.fields, scan) and from experiments
# with the oracle. scan, count.fields, read.table, read.csv and its relatives
# and the methods of type.convert are Go natives (eval/scan18.go,
# eval/readtable18.go); the generic is R.

type.convert <- function(x, ...) UseMethod("type.convert")

# the methods are not exported by R: registered, not visible by name
registerS3method("type.convert", "default", ..type.convert.default)
registerS3method("type.convert", "list", ..type.convert.list)
registerS3method("type.convert", "data.frame", ..type.convert.data.frame)

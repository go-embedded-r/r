## The raw type and the bit functions. A line holding only #--- separates cells.
#---
raw(3); raw(0); as.raw(c(1,255,16)); as.raw("12"); as.raw(TRUE); as.raw(2.9); as.raw(NULL); as.raw(c(a=1))
#---
as.raw(256)
#---
as.raw(-1)
#---
as.raw(NA)
#---
as.raw(1+0i)
#---
as.raw(list(1))
#---
typeof(as.raw(1)); class(as.raw(1)); mode(as.raw(1)); storage.mode(as.raw(1)); is.raw(as.raw(1)); is.raw(1); is.atomic(as.raw(1)); is.numeric(as.raw(1)); is.vector(as.raw(1))
#---
charToRaw("Hello"); charToRaw("é"); charToRaw(""); rawToChar(as.raw(c(72,105))); rawToChar(as.raw(c(72,105)), multiple=TRUE); rawToChar(raw(0)); rawToChar(as.raw(c(72,0)))
#---
rawToChar(as.raw(c(72,0,105)))
#---
charToRaw(c("a","b"))
#---
charToRaw(1)
#---
rawToChar("a")
#---
rawToChar(as.raw(c(0x41,0x42)), multiple=NA)
#---
as.integer(as.raw(255)); as.numeric(as.raw(7)); as.character(as.raw(c(1,171))); as.logical(as.raw(c(0,1,2))); as.complex(as.raw(3))
#---
c(as.raw(1), as.raw(2)); c(as.raw(1), 2L); c(as.raw(1), TRUE); c(as.raw(1), "a"); c(as.raw(1), list(1)); c(a=as.raw(1), b=as.raw(2)); c(as.raw(1), NULL)
#---
x <- as.raw(1:3); x[2]; x[-1]; x[[2]]; x[5]; x[c(TRUE,FALSE)]; x[0]; rev(x); rep(as.raw(1), 3); unique(as.raw(c(1,1,2))); length(raw(5)); rev(raw(0))
#---
x <- as.raw(1:3); x[2] <- as.raw(255); x
#---
x <- as.raw(1:3); x[5] <- as.raw(9); x
#---
x <- as.raw(1:3); names(x) <- c("a","b","c"); x
#---
x <- as.raw(1:3); x[2] <- 5L
#---
x <- 1:3; x[2] <- as.raw(5)
#---
x <- as.raw(1:3); x[[2]] <- as.raw(7); x
#---
x <- as.raw(1:3); x[2] <- list(1); x
#---
as.raw(1:3) == as.raw(c(1,5,3)); as.raw(1) == 1; as.raw(1:3) < as.raw(2); as.raw(1) == "01"; as.raw(10) != as.raw(10)
#---
as.raw(1:3) & as.raw(2); as.raw(1:3) | as.raw(4); !as.raw(1); xor(as.raw(12), as.raw(10)); as.raw(c(255, 0)) & as.raw(15)
#---
as.raw(1) + as.raw(1)
#---
-as.raw(1)
#---
as.raw(12) & 1L
#---
sort(as.raw(c(3,1,2)))
#---
order(as.raw(c(3,1,2)))
#---
identical(as.raw(1), as.raw(1)); identical(as.raw(1), 1L); identical(raw(2), raw(2)); identical(as.raw(1:2), as.raw(c(1,3)))
#---
rawToBits(as.raw(5)); intToBits(5L)[1:8]; typeof(intToBits(5L)); intToBits(-1L)[28:32]; length(intToBits(1:2)); intToBits(2.9)[1:4]
#---
packBits(intToBits(5L), "integer"); packBits(rawToBits(as.raw(c(1,2)))); packBits(c(TRUE,FALSE,TRUE,FALSE,FALSE,FALSE,FALSE,FALSE)); packBits(intToBits(c(1L,2L)), "integer"); packBits(rep(c(1L,0L), 4)); packBits(as.raw(1:8))
#---
packBits(1:3)
#---
packBits(c(TRUE,FALSE,TRUE), "integer")
#---
packBits(intToBits(1:2), "double")
#---
packBits("a")
#---
rawToBits(1)
#---
rawShift(as.raw(1), 3); rawShift(as.raw(255), -4); rawShift(as.raw(c(1,2,4)), 1); rawShift(raw(0), 2)
#---
rawShift(as.raw(1), 9)
#---
rawShift(1, 1)
#---
bitwAnd(12L, 10L); bitwOr(12L, 10L); bitwXor(12L, 10L); bitwNot(5L); bitwShiftL(1L, 4L); bitwShiftR(256L, 4L); bitwAnd(c(1,2,3), 3); bitwShiftL(1L, 31L); bitwShiftL(-1L, 1L); bitwShiftR(-1L, 1L); bitwAnd(1.5, 1)
#---
bitwAnd(c(1L, NA, 3L), 1L); bitwOr(1:4, c(8L, 16L)); bitwXor(NA_integer_, 1L); bitwNot(c(0L, -1L, NA)); bitwShiftL(1:3, 2L); bitwShiftL(1L, 32L); bitwShiftR(8L, c(1L, 2L, 3L)); bitwShiftL(1L, NA_integer_)
#---
bitwNot(NA)
#---
bitwAnd(5L, NA)
#---
bitwAnd("a", "b")
#---
utf8ToInt("A"); intToUtf8(c(72,105)); utf8ToInt("é"); intToUtf8(c(72,105), TRUE); intToUtf8(NULL); intToUtf8(0); utf8ToInt(NA_character_)
#---
utf8ToInt(c("a","b"))
#---
utf8ToInt(NA)
#---
is.na(as.raw(1)); is.na(raw(2)); as.vector(as.raw(1)); as.vector(as.raw(1:2), "list"); unlist(list(as.raw(1), as.raw(2))); vector("raw", 2); anyNA(raw(2)); is.finite(as.raw(1))
#---
matrix(as.raw(1:4), 2); array(as.raw(1:4), c(2,2)); t(matrix(as.raw(1:4), 2)); dim(matrix(raw(6), 2))
#---
max(as.raw(1))
#---
sum(as.raw(1))
#---
any(as.raw(1))
#---
paste(as.raw(10)); format(as.raw(c(1,255))); toString(as.raw(1:3)); cat(as.raw(c(1,255)), "\n"); sprintf("%s", as.raw(10)); as.character(as.raw(0:3)); nchar(as.raw(1))
#---
rawToChar(charToRaw("héllo")) == "héllo"; nchar(rawToChar(as.raw(c(0xc3,0xa9)))); strtoi("ff", 16L); as.raw(strtoi("ff", 16L))
#---
deparse(as.raw(1:3)); deparse(raw(0)); deparse(as.raw(255)); deparse(c(a=as.raw(1))); deparse(list(as.raw(1)))
#---
as.raw(c(1,2))[3]; as.raw(1:3)[c(TRUE, NA)]; match(as.raw(2), as.raw(1:3)); as.raw(2) %in% as.raw(1:3); duplicated(as.raw(c(1,1,2))); which(as.raw(1:3) == as.raw(2))
#---
sapply(as.raw(1:3), as.integer); lapply(as.raw(1:2), function(b) b); vapply(as.raw(1:3), function(b) as.integer(b) * 2L, 0); for (b in as.raw(1:2)) print(b)
#---
structure(as.raw(1), class="foo"); unclass(structure(as.raw(1), a=1)); attributes(structure(as.raw(1:2), names=c("a","b")))
#---
deparse(as.raw(1:40)); deparse(as.raw(c(1:25, 1:25)))

## The complex type: construction, printing, arithmetic, the math group, coercion. A line holding only #--- separates cells.
#---
1+2i; c(1.5+2i, -3+0.25i); complex(real=1e10, imaginary=1e-10); 1e-20+1i; 0i
#---
c(1+2i, NA); complex(real=NA, imaginary=1); complex(real=NaN, imaginary=Inf); complex(real=-Inf, imaginary=-Inf); NA_complex_
#---
c(123456.789+1i, 0.001+0.5i); 100000+0i; 1e5+1e-5i; 1/3+2/3i; c(0.1+0.2i, 100+200i); c(1e-10+1i, 1+1e10i); c(1.123456789+2i, 3+4.123456789i)
#---
-0.5+0i; 0-0.5i; complex(real=0, imaginary=-0); 1-0i; complex(real=1e-15, imaginary=1); 123456789+1i; 1e15+1i
#---
typeof(1i); class(1i); mode(1i); storage.mode(1i); is.numeric(1i); is.complex(1i); is.atomic(1i); is.complex(1); is.complex(NULL); is.vector(1i)
#---
Re(1+2i); Im(1+2i); Mod(3+4i); Arg(1i); Conj(1+2i); Re(3L); Im(3); Mod(-2); Arg(-2); Conj(2); Re(NA_complex_); Mod(NA_complex_)
#---
(1+2i)*(3-1i); (1+2i)^2; (1+2i)^3; 1i^2; (1+3i)-(2i); -1i; +1i
#---
1:3 + 1i; 1i * 1:3; TRUE + 1i; c(1+2i, 3) * c(1i, 2); 3L - 1i
#---
signif((1+2i)^0.5, 6); signif((1+2i)^1.5, 6); signif(2^(1+1i), 6); signif((1+1i)^(1+1i), 6); signif((-8+0i)^(1/3), 6)
#---
c(0i^2, 0i^0, 0i^-1, 0i^0.5, 1i^0, NA_complex_^0)
#---
c(0i^1i, 0i^(1+1i))
#---
complex(real=Inf, imaginary=0)^2; Inf*(0+0i)
#---
1e300*(1+1i)*(1e10+1e10i); (1e200+1e200i)*(1e200+1e200i); (0.1+0.2i)*(0.3+0.4i)
#---
1i %% 2
#---
1i %/% 2
#---
1i > 2
#---
1i < 2i
#---
1i == 1i; 1i == 2; 1i != 1i; 1i == "0+1i"; 1i == TRUE; 1+0i == 1L; NA_complex_ == 1i; NaN+0i == 1i
#---
1i & 1; 0i | 0; !1i; !0i; xor(1i, 0i); 1i && TRUE; 0i || FALSE; if (1i) "yes" else "no"; if (0i) "yes" else "no"
#---
signif(sqrt(as.complex(-1)), 6); sqrt(-1+0i); sqrt(-1-0i); sqrt(4+0i); signif(sqrt(1i), 6); signif(sqrt(-4+1i), 6)
#---
signif(exp(1i*pi), 6); signif(exp(1+1i), 6); exp(710+0i); exp(-1000+0i)
#---
signif(log(1i), 6); log(-1+0i); signif(log(1+1i), 6); signif(log(8+0i, base=2), 6)
#---
signif(sin(1i), 6); signif(cos(1+1i), 6); signif(tan(1i), 6); signif(sinh(1i), 6); signif(cosh(1i), 6); signif(tanh(1i), 6)
#---
signif(sin(1+2i), 6); signif(cos(2-1i), 6); signif(tan(1+1i), 6); signif(sinh(1+2i), 6); signif(cosh(1+2i), 6); signif(tanh(1+2i), 6)
#---
sin(800i); cos(800i); tan(800i); tan(1e3+1e3i); sinh(800+0i); cosh(800+0i); tanh(800+0i); tanh(-800+0i)
#---
asin(c(2+0i, -2+0i, 2-0i, 3+0i)); acos(c(2+0i, -2+0i, 3+0i)); atanh(c(2+0i, -2+0i, 2-0i))
#---
signif(asin(c(0.5+0i, 2+1i, 1i, 2i)), 6); signif(acos(c(0.5+0i, 2+1i, 1i, 2i)), 6); signif(atan(c(0.5i, 1+2i, 1+1i)), 6)
#---
atan(c(2i, -2i, 3i, -3i, 1i, -1i, 0i, -0+2i))
#---
signif(asinh(c(1i, 2i, 3+0i, 1+2i)), 6); signif(acosh(c(0i, 1i, -1+0i, 1+1i, 5+0i, 0.5+0i)), 6); signif(atanh(c(2i, 0.5+0i, 1+1i)), 6)
#---
atanh(c(1+0i, -1+0i))
#---
abs(3+4i); abs(-1+0i); abs(c(a=3+4i))
#---
floor(1+1i)
#---
trunc(1+1i)
#---
sum(1:3+1i); prod(c(1i,1i)); mean(c(1i,3+1i)); sum(1i, 2L, NA, na.rm=TRUE); prod(1i, 2i, 3); sum(c(1i, NaN+1i)); sum(1i, NA); mean(c(1i, NA), na.rm=TRUE); mean(c(1+1i, 3+3i, 5+0i)); mean(complex(0))
#---
cumsum(c(1i, 2, 3i)); cumprod(c(1i,1i,1i)); cumsum(c(1i, NA, 2)); cumprod(c(1i, 2, 1+1i))
#---
sum(1i, "a")
#---
max(1i)
#---
range(1i)
#---
cummax(1i)
#---
round(1.2345+6.789i, 2); round(2.567+3.123i); round(c(2.5+3.5i, -2.5-0.5i), 0); round(1234.5678+8765.4321i, -2); round(NA_complex_); round(1i, NA)
#---
signif(1/3+2/3i, 2); signif(123.456+0.00123i, 3); signif(c(123.456+0.00123i, 0.000123456+0.001i), 3); signif(1+1e-10i, 3); signif(0+0i, 3); signif(Inf+1i, 2); signif(1.5+2.5i, 1)
#---
as.complex(2); as.complex("1+2i"); as.complex(TRUE); as.complex(NA); as.complex(NULL); as.complex(c(a=1))
#---
as.complex("abc")
#---
as.complex(c("1", "i", "1+i", "1+2i", " 3-4i ", "1e2+3e-1i", "NaN", "Inf+1i", "NA", "+1i"))
#---
as.complex(list(1, 2i))
#---
as.complex(list(1:2))
#---
as.numeric(1+0i); as.integer(2+0i); as.logical(0i); as.logical(1i); as.logical(c(0i, 1i, NA_complex_)); as.integer(c(1+0i, 2.7+0i, NA_complex_))
#---
as.numeric(1+2i)
#---
as.integer(1+1i)
#---
as.numeric(c(1+0i, NA_complex_, NaN+0i))
#---
as.character(1+2i); as.character(1.5+2.25i); as.character(1/3+1i); as.character(c(1e5+1i, 1e-5+1e15i, NaN+1i, Inf-Inf*1i, NA_complex_, 100000+1i))
#---
paste(1+2i); paste0("z=", 3i); paste(c(1+2i, NA, 1/3+0i)); toString(c(1i, 2)); toString(c(1+2i, NA))
#---
format(1+2i); format(c(1+2i, 10.5-3i)); format(c(1+2i, NA, 3.5-1i)); format(1e10+1i); format(NA_complex_); format(complex(0)); format(1/3+2/3i, digits=3)
#---
format(c(1.5+2i, 3+4i), nsmall=2); format(1+2i, width=10); format(-1.5+2i, nsmall=1); format(1+2i, scientific=TRUE); format(123456+1i, big.mark=","); format(c(1+1i, 10+10i), trim=TRUE); format(c(1i, NA))
#---
cat(1+2i, 3i, "\n"); cat(1/3+1i, 1e5+1e-5i, NA_complex_, "\n"); cat(c(1+1i, 2+22i), sep=" | "); cat(complex(real=Inf, imaginary=NaN), "\n")
#---
sprintf("%s", c(1+2i, NA)); sprintf("%5s|", 1i); sprintf("%f", 1+0i)
#---
sprintf("%d", 1i)
#---
sprintf("%f", 1i)
#---
c(a=1+1i, b=2); c(1+1i, 2L, TRUE, 2.5); c(1i, "a"); c(1i, NULL, 2); c(NULL, 1i); c(a=1i, 2); c(1i, c(b=2)); c(1i, list(1))
#---
c(1i, NA)[2]; Im(c(1i, NA)[2]); Re(c(1i, NA)[2]); Im(as.complex(NA)); Re(complex(real=1, imaginary=NA) + 1); is.na(complex(real=1, imaginary=NA))
#---
x <- 1:3 + 0i; x[2]; x[[3]]; x[-1]; x[c(TRUE, FALSE)]; x[5]; x[0]; rev(x); rep(x, 2); rep(1i, 3); length(x); head(x, 2); tail(x, 1); append(1i, 2i)
#---
x <- c(a=1i, b=2i); x["a"]; x["z"]; x[["b"]]; names(x); rev(x)
#---
x <- c(1i, 2i); x[5] <- 3; x
#---
x <- 1:3; x[2] <- 1i; x
#---
x <- c(1i, 2i); x[[2]] <- 5L; x
#---
x <- c(1i, 2i); x[2] <- "a"; x
#---
x <- c(1i, 2i); x[2] <- TRUE; x
#---
x <- c("a"); x[2] <- 1i; x
#---
x <- list(1); x[[2]] <- 1i; x
#---
x <- c(1i, 2i); length(x) <- 4; x
#---
x <- matrix(1i, 2, 2); x[1, 2] <- 5; x
#---
is.na(c(1i, NA)); anyNA(c(1i, NA)); is.na(complex(real=NaN, imaginary=1)); is.nan(complex(real=NaN, imaginary=1)); is.finite(c(1i, complex(real=Inf, imaginary=0), NA)); is.infinite(complex(real=Inf, imaginary=0))
#---
identical(1i, 1i); identical(1i, 2i); identical(NaN+1i, NaN+1i); identical(NA_complex_, complex(real=NA, imaginary=NA)); identical(1+0i, as.complex(1)); identical(1i, 1)
#---
unique(c(1i, 1i, 2i)); duplicated(c(1i, 1i, 2)); match(1+0i, c(2, 1)); 1i %in% 1i; match(2i, c(1i,2i)); which(c(1i,0i)==1i); any(1i); all(c(1i, 0i))
#---
sort(c(2+1i, 1+5i, 1+2i)); order(c(2+1i, 1+5i, 1+2i)); sort(c(2i, 1i, NA)); sort(c(2i, 1i, 1+1i, 1-1i), decreasing=TRUE); order(c(1+1i, 1i, 1-1i)); rank(c(2+1i, 1+5i))
#---
matrix(c(1i, 2, 3, 4i), 2); t(matrix(c(1i, 2, 3, 4i), 2)); diag(2) * 1i; matrix(c(1+1i, NA, 3.5, 4i), 2, dimnames=list(c("a","b"), NULL)); array(1:8 * 1i, c(2,2,2))
#---
matrix(1:4 + 0i, 2) %*% matrix(1:4 + 0i, 2); crossprod(matrix(1:4+0i, 2)); outer(1:2, 1:2 * 1i); diag(c(1i, 2i)); rowSums(matrix(1:4 * 1i, 2)); colMeans(matrix(1:4 * 1i, 2)); apply(matrix(1:4 * 1i, 2), 1, sum)
#---
vapply(1:2, function(i) i*1i, 0i); sapply(c(1i, 2i), Mod); unlist(list(1i, "a")); unlist(lapply(1:2, function(i) i * 1i)); Reduce(`*`, list(1i, 1i, 1i)); do.call(sum, list(1i, 2i)); sapply(1:2, function(i) i * 1i)
#---
deparse(1+2i); deparse(c(1i, 2)); deparse(complex(0)); deparse(NA_complex_); deparse(c(a=1i)); deparse(c(1+2i, NA, 3i)); deparse(complex(real=Inf, imaginary=1)); deparse(complex(real=NaN, imaginary=1)); deparse(1e5+1i); deparse(1/3+1i); deparse(-1i)
#---
quote(1i); typeof(quote(1i)); class(quote(1i)); eval(quote(2i)); deparse(quote(1+2i)); substitute(x + 1i, list(x = 2i)); bquote(.(1i) + x)
#---
complex(3); complex(real=1:3, imaginary=0); complex(length.out=2, real=1); complex(real=1, imaginary=1:2); complex(0); vector("complex", 2); numeric(0) + 0i; integer(0) + 1i; c(1i)[0]
#---
signif(complex(modulus=2, argument=pi/2), 6); complex(modulus=1:2, argument=0); complex(modulus=2)
#---
complex(real="a")
#---
complex(real=1, imaginary=c(1,2,3), length.out=2)
#---
storage.mode(1i); x <- 1:2; storage.mode(x) <- "complex"; x; as.vector(1i, "complex"); as.vector(1:2, "complex"); as.vector(1i, "list"); as.vector(c(a=1i), "complex")
#---
x <- c(1i, 2i); names(x) <- c("a", "b"); x[["a"]] <- 0; x
#---
nchar(c(1+2i, 1/3+1i)); is.na(sprintf("%s", NA_complex_)); typeof(NA_complex_ + 1); typeof(c(NA_complex_, 1))
#---
structure(1i, class="foo"); unclass(structure(1i, a=1)); attributes(structure(1i, a=1))
#---
rowSums(matrix(1i, 2, 2)); colSums(matrix(1:4 * 1i, 2))
#---
1:3 + 1i + 1:2
#---
deparse(c(1+1i, 2+2i, 3+3i, 4+4i, 5+5i, 6+6i, 7+7i, 8+8i, 9+9i, 10+10i, 11+11i)); deparse(c(1.5+1.5i, 2.5+2.5i, 3.5+3.5i, 4.5+4.5i, 5.5+5.5i, 6.5+6.5i, 7.5+7.5i)); deparse(c(11111111111.5+1.5i, 2.5+2.5i, 3.5+3.5i, 4.5+4.5i, 5.5+5.5i))
#---
deparse(c(NA_complex_, 1i)); deparse(c(1i, NA)); deparse(c(NA_complex_, NA_complex_)); deparse(c(NA_real_, 1)); deparse(c(NA_integer_, 1L)); deparse(c(NA_character_, "a")); deparse(c(NA, 1L, NA)); deparse(c(NA_integer_, NA_integer_))
#---
deparse(c(1.11111111111111, 2.22222222222222, 3.33333333333333, 4.4444444444444)); deparse(c(1.11111111111111, 2.22222222222222, 3.33333333333333, 4.44444444444444, 5.5)); deparse(c("aaaaaaaaaaaaaaaaaaaa", "bbbbbbbbbbbbbbbbbbbbbbbb", "cccccccccccccccccccccc", "ddd"))

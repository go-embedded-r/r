format(as.Date("0005-03-10"), "%Y|%y|%C")
#---
as.POSIXct("2024-03-31 02:30:00", tz = "Europe/Paris")
#---
as.POSIXct("2024-10-27 02:30:00", tz = "Europe/Paris")
#---
format(as.POSIXct("2024-07-08 07:32:13", tz = "Asia/Tokyo"), "%s %Z")
#---
format(as.POSIXct("2024-07-15 12:00:00", tz = "Europe/Paris"), "%s")
#---
format(as.POSIXct("2024-07-15 12:00:00", tz = "Europe/Paris"), "%s %Z|%s")
#---
format(as.POSIXct("2024-01-15 12:00:00", tz = "America/New_York"), "%s %z")
#---
format(as.POSIXct("2024-07-15 12:00:00", tz = "Australia/Sydney"), "%Z %s")
#---
format(as.POSIXct("2024-07-15 12:00:00", tz = "Asia/Tokyo"), "%s")
#---
(1+2i)/(3-1i); (1+2i)^-2; 1i / 2; 2 / 1i; (1+1i)/(1-1i); 1/(1+3i); (1+3i)^-1; (3+4i)/(1+2i); (0.1+0.2i)/(0.3+0.4i); 1/3i
#---
(1+2i)/(0+0i); 0i/0i; (1+1i)/Inf; (1e308+0i)/(-Inf+2i); Inf/(1+1i); (1+2i)/Inf; 0i^-1; 0i^-2
#---
prod(c(1+1i, 1e-5-0.25i, 10-2.5i, 123456.789+2i)); prod(c(123456.789+10i, 10+1/3i)); cumprod(c(pi-2.5i, -0.25+0i, 1/3i, -0.25+7i)); cumprod(c(0.1+1e-5i, 10-2.5i))
#---
cumsum(c(1i, NA, 2)); cumsum(c(NA_real_+0.5i, 1e-5+0.1i, NaN+10i, pi+1i)); cumsum(c(1e5-2.5i, NA_real_+1e5i, 1/3+1e-5i))
#---
Re(1+2i); Im(1+2i); Mod(3+4i); Arg(1i); Conj(1+2i); Re(3L); Im(3); Mod(-2); Arg(-2); Conj(2); Re(NA_complex_); Mod(NA_complex_)
#---
c(0i^2, 0i^0, 0i^-1, 0i^0.5, 1i^0, NA_complex_^0)
#---
c(0i^1i, 0i^(1+1i))
#---
complex(real=Inf, imaginary=0)^2; Inf*(0+0i)
#---
1e300*(1+1i)*(1e10+1e10i); (1e200+1e200i)*(1e200+1e200i); (0.1+0.2i)*(0.3+0.4i)
#---
signif(sqrt(as.complex(-1)), 6); sqrt(-1+0i); sqrt(-1-0i); sqrt(4+0i); signif(sqrt(1i), 6); signif(sqrt(-4+1i), 6)
#---
signif(exp(1i*pi), 6); signif(exp(1+1i), 6); exp(710+0i); exp(-1000+0i)
#---
signif(log(1i), 6); log(-1+0i); signif(log(1+1i), 6); signif(log(8+0i, base=2), 6)
#---
sin(800i); cos(800i); tan(800i); tan(1e3+1e3i); sinh(800+0i); cosh(800+0i); tanh(800+0i); tanh(-800+0i)
#---
asin(c(2+0i, -2+0i, 2-0i, 3+0i)); acos(c(2+0i, -2+0i, 3+0i)); atanh(c(2+0i, -2+0i, 2-0i))
#---
atan(c(2i, -2i, 3i, -3i, 1i, -1i, 0i, -0+2i))
#---
signif(asinh(c(1i, 2i, 3+0i, 1+2i)), 6); signif(acosh(c(0i, 1i, -1+0i, 1+1i, 5+0i, 0.5+0i)), 6); signif(atanh(c(2i, 0.5+0i, 1+1i)), 6)
#---
atanh(c(1+0i, -1+0i))
#---
floor(1+1i)
#---
cumsum(c(1i, 2, 3i)); cumprod(c(1i,1i,1i)); cumsum(c(1i, NA, 2)); cumprod(c(1i, 2, 1+1i))
#---
signif(1/3+2/3i, 2); signif(123.456+0.00123i, 3); signif(c(123.456+0.00123i, 0.000123456+0.001i), 3); signif(1+1e-10i, 3); signif(0+0i, 3); signif(Inf+1i, 2); signif(1.5+2.5i, 1)
#---
as.complex(c("1", "i", "1+i", "1+2i", " 3-4i ", "1e2+3e-1i", "NaN", "Inf+1i", "NA", "+1i"))

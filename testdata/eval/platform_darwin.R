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

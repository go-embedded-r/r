scale <- function(x, center = TRUE, scale = TRUE) UseMethod("scale")

scale.default <- function(x, center = TRUE, scale = TRUE) {
  m <- as.matrix(x)
  if (!is.logical(center) && length(center) != ncol(m))
    stop("length of 'center' must equal the number of columns of 'x'")
  if (!is.logical(scale) && length(scale) != ncol(m))
    stop("length of 'scale' must equal the number of columns of 'x'")
  ..scale_matrix(m, center, scale)
}

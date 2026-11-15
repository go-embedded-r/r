# Base library, slice 17 (1): formulas. Written from the help pages of R 4.6.2
# (formula, terms, terms.object, terms.formula, update.formula, reformulate,
# all.vars, model.frame, model.response, get_all_vars) and from experiments
# with the oracle. The term algebra itself is the Go native ..termsFormula.

terms <- function(x, ...) UseMethod("terms")

registerS3method("terms", "default", function(x, ...) {
  found <- x$terms
  if (is.null(found)) found <- attr(x, "terms")
  if (is.null(found)) stop("no terms component nor attribute")
  found
})

registerS3method("terms", "terms", function(x, ...) x)

terms.formula <- function(x, specials = NULL, abb = NULL, data = NULL,
                          neg.out = TRUE, keep.order = FALSE, simplify = FALSE, ...,
                          allowDotAsName = FALSE) {
  if (!is.null(abb)) warning("setting 'abb' in terms.formula() is deprecated")
  if (!is.null(specials)) specials <- as.character(specials)
  ..termsFormula(sys.call(), x, specials, if (!is.null(data)) names(data),
                 isTRUE(keep.order), isTRUE(simplify), isTRUE(allowDotAsName))
}

is.empty.model <- function(x) {
  tt <- terms(x)
  length(attr(tt, "factors")) == 0L && attr(tt, "intercept") == 0L
}

# ---- formula objects ---------------------------------------------------------

formula <- function(x, ...) UseMethod("formula")

registerS3method("formula", "default", function(x = NULL, env = parent.frame(), ...) {
  if (!is.atomic(x)) {
    if (!is.null(x$formula)) return(eval(x$formula))
    if (!is.null(x$terms)) {
      out <- x$terms
      oldClass(out) <- "formula"
      return(out)
    }
    if (!is.null(x$call$formula)) return(eval(x$call$formula))
  }
  if (!is.null(attr(x, "formula"))) return(attr(x, "formula"))
  made <- switch(mode(x),
                 "NULL" = structure(list(), class = "formula"),
                 character = eval(str2lang(x)),
                 call = eval(x),
                 stop("invalid formula"))
  if (is.call(made) || is.null(attr(made, ".Environment"))) environment(made) <- env
  made
})

registerS3method("formula", "formula", function(x, ...) x)

registerS3method("formula", "terms", function(x, ...) {
  keep <- environment(x)
  attributes(x) <- list(class = "formula")
  environment(x) <- if (is.null(keep)) globalenv() else keep
  x
})

registerS3method("formula", "character", function(x, env = parent.frame(), ...) {
  made <- eval(str2lang(x))
  if (!inherits(made, "formula")) stop("invalid formula")
  environment(made) <- env
  made
})

registerS3method("formula", "data.frame", function(x, ...) {
  nms <- names(x)
  if (length(nms) < 2L) stop("cannot make a formula from fewer than two columns")
  quoted <- ifelse(make.names(nms) == nms, nms, paste0("`", nms, "`"))
  out <- eval(str2lang(paste(quoted[1L], "~", paste(quoted[-1L], collapse = " + "))))
  environment(out) <- parent.frame()
  out
})

as.formula <- function(object, env = parent.frame()) {
  if (inherits(object, "formula")) return(object)
  made <- formula(object, env = baseenv())
  if (identical(environment(made), baseenv()) || !missing(env)) environment(made) <- env
  made
}

registerS3method("print", "formula", function(x, showEnv = !identical(e, .GlobalEnv), ...) {
  e <- environment(.x <- x)
  attr(.x, ".Environment") <- NULL
  print.default(unclass(.x), ...)
  if (showEnv) print(e)
  invisible(x)
})

registerS3method("format", "formula", function(x, ...) {
  attr(x, ".Environment") <- NULL
  paste(deparse(unclass(x)), collapse = " ")
})

registerS3method("[", "formula", function(x, i) {
  out <- unclass(x)
  keep <- attr(x, ".Environment")
  attr(out, ".Environment") <- NULL
  out <- out[i]
  if (!length(out)) {
    out <- list()
    class(out) <- "formula"
    environment(out) <- keep
  } else if (is.call(out) && identical(out[[1L]], quote(`~`))) {
    class(out) <- "formula"
    environment(out) <- keep
  }
  out
})

# ---- building and changing formulas -------------------------------------------

update <- function(object, ...) UseMethod("update")

update.formula <- function(old, new, ...) {
  old <- as.formula(old)
  new <- as.formula(new)
  put <- function(part, with, wrap = TRUE) {
    if (wrap && is.call(with)) with <- call("(", with)
    do.call("substitute", list(part, list(. = with)))
  }
  oldrhs <- old[[length(old)]]
  oldlhs <- if (length(old) == 3L) old[[2L]]
  made <- new
  made[[length(new)]] <- put(new[[length(new)]], oldrhs)
  if (length(new) == 3L) {
    made[[2L]] <- if (is.null(oldlhs)) new[[2L]] else put(new[[2L]], oldlhs, FALSE)
  } else if (!is.null(oldlhs)) {
    made <- eval(call("~", oldlhs, made[[2L]]))
  }
  attr(made, ".Environment") <- NULL
  out <- formula(terms.formula(made, simplify = TRUE))
  environment(out) <- environment(old)
  out
}

reformulate <- function(termlabels, response = NULL, intercept = TRUE, env = parent.frame()) {
  if (!is.character(termlabels)) stop("'termlabels' must be a character vector")
  text <- paste(termlabels, collapse = "+")
  if (!intercept) text <- paste(text, "- 1")
  rhs <- str2lang(text)
  lhs <- NULL
  if (!is.null(response)) {
    if (is.character(response)) {
      if (length(response) != 1L) stop("'response' must be a character string")
      lhs <- tryCatch(str2lang(response), error = function(e) {
        warning(gettextf("Unparseable 'response' \"%s\"; use is deprecated.  Use as.name(.) or `..`!", response))
        as.name(response)
      })
    } else if (is.name(response) || is.call(response)) {
      lhs <- response
    } else stop("'response' must be a character string, name or call")
  }
  out <- if (is.null(lhs)) call("~", rhs) else call("~", lhs, rhs)
  out <- eval(out, envir = env)
  environment(out) <- env
  out
}

# ---- model frames --------------------------------------------------------------

na.fail <- function(object, ...) UseMethod("na.fail")
registerS3method("na.fail", "default", function(object, ...) {
  ok <- complete.cases(object)
  if (all(ok)) object else stop("missing values in object")
})

na.pass <- function(object, ...) object

na.omit <- function(object, ...) UseMethod("na.omit")

registerS3method("na.omit", "default", function(object, ...) {
  if (!is.atomic(object)) return(object)
  if (length(dim(object)) > 1L) {
    drop <- unique(((which(is.na(object)) - 1) %% nrow(object)) + 1)
  } else {
    drop <- which(is.na(object))
  }
  if (!length(drop)) return(object)
  kept <- if (length(dim(object)) > 1L) object[-drop, , drop = FALSE] else object[-drop]
  nm <- names(object)
  if (length(dim(object)) > 1L) nm <- rownames(object)
  if (!is.null(nm)) names(drop) <- nm[drop]
  attr(drop, "class") <- "omit"
  attr(kept, "na.action") <- drop
  kept
})

registerS3method("na.omit", "data.frame", function(object, ...) {
  gone <- logical(nrow(object))
  for (col in object) {
    if (!is.atomic(col)) next
    miss <- is.na(col)
    if (length(dim(miss)) == 2L) for (j in seq_len(ncol(miss))) gone <- gone | miss[, j]
    else gone <- gone | miss
  }
  drop <- which(gone)
  if (!length(drop)) return(object)
  names(drop) <- row.names(object)[drop]
  kept <- object[-drop, , drop = FALSE]
  attr(drop, "class") <- "omit"
  attr(kept, "na.action") <- drop
  kept
})

na.exclude <- function(object, ...) UseMethod("na.exclude")
registerS3method("na.exclude", "default", function(object, ...) {
  out <- na.omit(object, ...)
  act <- attr(out, "na.action")
  if (!is.null(act)) {
    class(act) <- "exclude"
    attr(out, "na.action") <- act
  }
  out
})

complete.cases <- function(...) {
  parts <- list(...)
  n <- NULL
  ok <- NULL
  add <- function(flag, len) {
    if (is.null(n)) {
      n <<- len
      ok <<- rep(TRUE, len)
    } else if (len != n) stop("not all arguments have the same length")
    ok <<- ok & flag
  }
  for (p in parts) {
    if (is.data.frame(p) || is.list(p)) {
      for (col in p) {
        if (is.matrix(col)) add(rowSums(is.na(col)) == 0, nrow(col))
        else add(!is.na(col), length(col))
      }
    } else if (is.matrix(p)) {
      add(rowSums(is.na(p)) == 0, nrow(p))
    } else add(!is.na(p), length(p))
  }
  if (is.null(n)) stop("no input has determined the number of cases")
  ok
}

model.frame <- function(formula, ...) UseMethod("model.frame")

model.frame.default <- function(formula, data = NULL, subset = NULL, na.action = na.fail,
                                drop.unused.levels = FALSE, xlev = NULL, ...)
  ..modelFrame(missing(formula), missing(data), missing(na.action),
               if (!missing(formula)) formula, if (!missing(data)) data,
               substitute(subset), na.action, drop.unused.levels, xlev,
               substitute(list(...)), parent.frame())

get_all_vars <- function(formula, data = NULL, ...)
  ..getAllVars(missing(formula), missing(data), if (!missing(formula)) formula,
               if (!missing(data)) data, substitute(list(...)), parent.frame())

labels <- function(object, ...) UseMethod("labels")

labels.default <- function(object, ...) {
  if (length(dim(object))) {
    given <- dimnames(object)
    if (is.null(given)) given <- vector("list", length(dim(object)))
    for (i in seq_along(given)) if (is.null(given[[i]])) given[[i]] <- as.character(seq_len(dim(object)[i]))
    given
  } else {
    nm <- names(object)
    if (is.null(nm)) as.character(seq_along(object)) else nm
  }
}

registerS3method("labels", "terms", function(object, ...) attr(object, "term.labels"))

registerS3method("print", "terms", function(x, ...) {
  print.default(unclass(x), ...)
  invisible(x)
})

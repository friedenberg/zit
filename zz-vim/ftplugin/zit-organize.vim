
setlocal list
let &l:t_ut = ''
let &l:listchars = "tab:  ,trail:·,nbsp:·"
let &l:equalprg = "zit format-organize %"

let &l:foldmethod = "expr"
let &l:foldexpr = "GetZitOrganizeFold(v:lnum)"

function! GetZitOrganizeFold(lnum)
  if getline(a:lnum) =~? '\v^\s*$'
    return '-1'
  endif

  let this_indent = indent(a:lnum)

  if getline(a:lnum) =~? '\v^\s*#'
    return '>' . (this_indent + 1)
  else
    return this_indent + 1
  endif
endfunction

function! Gf()
  let l:h = trim(system("zit expand-hinweis " . expand("<cfile>")))

  if !filereadable(l:h)
    echo system("zit checkout -mode both " . l:h)
  endif

  let l:f = l:h . ".md"

  execute 'tabedit' l:f
endfunction

noremap gf :call Gf()<CR>


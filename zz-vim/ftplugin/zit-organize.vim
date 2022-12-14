
setlocal list
" TODO document
let &l:t_ut = ''
let &l:listchars = "tab:  ,trail:·,nbsp:·"
let &l:equalprg = "zit format-organize -metadatei-header %"

let &l:foldmethod = "expr"
let &l:foldexpr = "GetZitOrganizeFold(v:lnum)"

" TODO implement against new organize syntax
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

function! GfOrganize()
  let l:h = trim(system("zit expand-hinweis " . expand("<cfile>")))

  if !filereadable(l:h)
    echo system("zit checkout -mode both " . l:h)
  endif

  " TODO dynamically source zettel file extension
  let l:f = l:h . ".zettel"

  execute 'tabedit' l:f
endfunction

noremap gf :call GfOrganize()<CR>


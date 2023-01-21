
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

" TODO refactor into common
function! GfOrganize()
  let l:cfile = expand("<cfile>")
  let l:expanded = trim(system("zit expand-hinweis " .. l:cfile))

  if !filereadable(l:expanded .. ".zettel")
    echom trim(system("zit checkout -mode both " .. l:expanded))
  endif

  " let l:cmd = 'tabedit ' .. fnameescape(l:expanded .. ".zettel")
  execute 'tabedit' fnameescape(l:expanded .. ".zettel")
endfunction

" <buffer> restricts this remap to the buffer in which it was defined
" https://learnvimscriptthehardway.stevelosh.com/chapters/11.html
noremap <buffer> gf :call GfOrganize()<CR>


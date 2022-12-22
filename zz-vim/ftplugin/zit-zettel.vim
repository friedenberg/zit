
let &l:equalprg = "zit format-zettel -include-cwd %:r"
let &l:comments = "fb:*,fb:-,fb:+,n:>"
let &l:commentstring = "<!--%s-->"

function! GfZettel()
  let l:f = expand("<cfile>") . ".md"

  if !filereadable(l:f)
    echo system("zit checkout -mode both " . expand("<cfile>"))
  endif

  execute 'tabedit' l:f
  " try
  "   " exec "normal! \<c-w>gf"
  " catch /E447/
  " endtry
endfunction

noremap gf :call GfZettel()<CR>

" TODO support external akte
function! ZitTypActionMenu()
  let l:items = systemlist("zit show -format action-names " . expand("%:r"))

  func! ZitTypActionMenuItemPicked(id, result) closure
    if a:result == -1
      return
    endif

    let l:val = substitute(l:items[a:result-1], '\t.*$', '', '')
    execute("!zit exec-action -action " . l:val .  " " . expand("%:r"))
  endfunc

  if len(l:items) == 0
    echom "No Zettel-Typ-specific actions available"
    return
  endif

  call popup_menu(
        \ items,
        \ #{ title: "Run a Zettel-Typ-Specific Action", 
        \ callback: 'ZitTypActionMenuItemPicked', 
        \ line: 25, col: 40,
        \ highlight: 'Question', border: [], close: 'click',  padding: [1,1,0,1]} )
endfunction

function! ZitCopyFormattersMenu()
  let l:rawItems = systemlist("zit show -format typ-formatter-uti-groups " . expand("%:r"))
  let l:processedItems = []
  let l:items = []

  for i in l:rawItems
    let l:groupName = substitute(i, '\s.*$', '', '')
    let l:group = i[len(l:groupName) +1:]
    call add(l:items, l:groupName)
    call add(l:processedItems, l:group)
  endfor

  func! ZitCopyMenuItemPicked(id, result) closure
    if a:result == -1
      return
    endif

    let l:uti_group = l:items[a:result-1]
    let l:val = substitute(l:processedItems[a:result-1], '\t.*$', '', '')
    let l:cmd_args_unprocessed_list = split(l:val)
    let l:cmd_args_list = []

    let i = 0

    while i < len(l:cmd_args_unprocessed_list)
      let l:uti = l:cmd_args_unprocessed_list[i]
      let l:formatter = l:cmd_args_unprocessed_list[i+1]
      call add(l:cmd_args_list, "-i")
      call add(l:cmd_args_list, l:uti)
      let l:cmd_sub_args = [
            \ "zit", "format-zettel", "-include-cwd", "-mode akte",
            \ "-uti-group", l:uti_group,
            \ l:uti, expand("%:r"),
            \ ]

      call add(l:cmd_args_list, "<(" . join(l:cmd_sub_args, " ") . ")")

      let i += 2
    endwhile

    execute("!tacky copy " . join(l:cmd_args_list, " "))
  endfunc

  if len(l:processedItems) == 1
    call ZitCopyMenuItemPicked("", 1)
    return
  endif

  if len(l:items) == 0
    echom "No Zettel-Typ-specific actions available"
    return
  endif

  call popup_menu(
        \ items,
        \ #{ title: "Run a Zettel-Typ-Specific Action", 
        \ callback: 'ZitCopyMenuItemPicked', 
        \ line: 25, col: 40,
        \ highlight: 'Question', border: [], close: 'click',  padding: [1,1,0,1]} )
endfunction

let maplocalleader = "-"

nnoremap <localleader>z :call ZitTypActionMenu()<cr>
nnoremap <localleader>c :call ZitCopyFormattersMenu()<cr>

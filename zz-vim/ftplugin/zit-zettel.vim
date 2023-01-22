
let &l:equalprg = "zit format-zettel -include-cwd %:r"
let &l:comments = "fb:*,fb:-,fb:+,n:>"
let &l:commentstring = "<!--%s-->"

function! GfZettel()
  let l:h = expand("<cfile>")
  let l:expanded = system("zit expand-hinweis " . l:h)
  let l:f = l:expanded . ".zettel"

  if !filereadable(l:f)
    echo system("zit checkout -mode both " . l:expanded)
  endif

  execute 'tabedit' l:f
  " try
  "   " exec "normal! \<c-w>gf"
  " catch /E447/
  " endtry
endfunction

noremap <buffer> gf :call GfZettel()<CR>

" TODO support external akte
function! ZitAction()
  let l:items = ZitGetActionNames()

  func! ZitActionItemPicked(id, result) closure
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
        \ callback: 'ZitActionItemPicked', 
        \ line: 25, col: 40,
        \ highlight: 'Question', border: [], close: 'click',  padding: [1,1,0,1]} )
endfunction

function! ZitMakeUTIGroupCommand(uti_group, cmd_args_unprocessed_list)
  let l:cmd_args_list = []

  let i = 0

  while i < len(a:cmd_args_unprocessed_list)
    let l:uti = a:cmd_args_unprocessed_list[i]
    let l:formatter = a:cmd_args_unprocessed_list[i+1]
    call add(l:cmd_args_list, "-i")
    call add(l:cmd_args_list, l:uti)
    let l:cmd_sub_args = [
          \ "zit", "format-zettel", "-include-cwd", "-mode akte",
          \ "-uti-group", a:uti_group, l:uti,
          \ expand("%:r"),
          \ "2>/dev/null",
          \ ]

    call add(l:cmd_args_list, "<(" . join(l:cmd_sub_args, " ") . ")")

    let i += 2
  endwhile

  return l:cmd_args_list
endfunction

function! ZitGetUTIGroups()
  let l:rawItems = sort(systemlist("zit show -include-cwd -format typ-formatter-uti-groups " . expand("%:r")))
  let l:processedItems = []
  let l:items = []

  for i in l:rawItems
    let l:groupName = substitute(i, '\s.*$', '', '')
    let l:group = i[len(l:groupName) +1:]
    call add(l:items, l:groupName)
    call add(l:processedItems, l:group)
  endfor

  return [l:items, l:processedItems]
endfunction

function! ZitGetActionNames()
  return sort(systemlist("zit show -format action-names " .. expand("%:r")))
endfunction

function! ZitGetFormats()
  return sort(systemlist("zit show -format formatters " .. expand("%:r")))
endfunction

function! ZitPreview()
  let l:items = ZitGetFormats()

  func! ZitPreviewMenuItemPicked(id, result) closure
    if a:result == -1
      return
    endif

    let l:format = substitute(l:items[a:result-1], '\t.*$', '', '')
    let l:hinweis = expand("%:r")

    let l:tempfile = tempname() .. "." .. l:format

    let l:cmd_args_list = [
          \ "zit format-zettel -mode akte",
          \ l:format,
          \ l:hinweis,
          \ ">",
          \ l:tempfile,
          \ ]

    call system(join(l:cmd_args_list, " "))
    echom l:tempfile
    call system("qlmanage -p "..l:tempfile..">/dev/null 2>&1 &")
  endfunc

  if len(l:items) == 1
    call ZitPreviewMenuItemPicked("", 1)
    return
  endif

  if len(l:items) == 0
    echom "No Zettel-Typ-specific actions available"
    return
  endif

  call popup_menu(
        \ items,
        \ #{ title: "Preview a Zettel-Typ format", 
        \ callback: 'ZitPreviewMenuItemPicked', 
        \ line: 25, col: 40,
        \ highlight: 'Question', border: [], close: 'click',  padding: [1,1,0,1]} )
endfunction

function! ZitCopy()
  let [l:items, l:processedItems] = ZitGetUTIGroups()

  func! ZitCopyMenuItemPicked(id, result) closure
    if a:result == -1
      return
    endif

    let l:uti_group = l:items[a:result-1]
    let l:val = substitute(l:processedItems[a:result-1], '\t.*$', '', '')
    let l:cmd_args_list = ZitMakeUTIGroupCommand(l:uti_group, split(l:val))

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

nnoremap <localleader>z :call ZitAction()<cr>
nnoremap <localleader>c :call ZitCopy()<cr>
nnoremap <localleader>p :call ZitPreview()<cr>

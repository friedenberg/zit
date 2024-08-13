
" if exists("b:current_syntax")
"   finish
" endif

if $BIN_ZIT == ""
  let BIN_ZIT = "zit"
endif

let zettel = expand("%")

let g:markdown_syntax_conceal = 0

if zettel != ""
  let cmdFormat = "$BIN_ZIT show -quiet -format typ.vim-syntax-type " . zettel
  let zettelTypSyntax = trim(system(cmdFormat))

  if v:shell_error
    echom "Error getting vim syntax type: " . zettelTypSyntax
    let zettelTypSyntax = "markdown"
  elseif zettelTypSyntax == ""
    echom "Zettel Typ has no vim syntax set"
    let zettelTypSyntax = "markdown"
  endif

  " let syntaxFile = $VIMRUNTIME . "/syntax/" . zettelTypSyntax . ".vim"
  " let ftpluginFile = $VIMRUNTIME . "/ftplugin/" . zettelTypSyntax . ".vim"

  execute "syntax include @akte" "syntax/" . zettelTypSyntax . ".vim"
  " if filereadable(syntaxFile)
  "   exec "source " . syntaxFile
  "   " TODO-P3
  "   " exec "source " . ftpluginFile
  " else
  "   let syntaxFile = $VIMRUNTIME . "/syntax/markdown.vim"
  "   let ftpluginFile = $VIMRUNTIME . "/ftplugin/markdown.vim"

  "   exec "source " . syntaxFile
  "   " TODO-P3
  "   " exec "source " . ftpluginFile
  " endif
endif

" syn case match
"
syn region zitAkte start=// end=// contains=@akte

let m = expand("<sfile>:h") . "/zit-metadatei.vim"
exec "source " . m

let b:current_syntax = 'zit-zettel'

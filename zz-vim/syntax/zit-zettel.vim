
" if exists("b:current_syntax")
"   finish
" endif

let zettel = expand("%")
echom zettel

let g:markdown_syntax_conceal = 0

if zettel != ""
  let cmdFormat = "zit show -format typ-vim-syntax-type " . zettel
  let zettelTypSyntax = trim(system(cmdFormat))
  echom cmdFormat
  echom zettelTypSyntax

  if zettelTypSyntax == ""
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

syn region zitMetadatei start=/\v%^---$/ end=/\v^---$/ 
      \ contains=zitMetadateiBezeichnungRegion,zitMetadateiEtikettRegion,zitMetadateiAkteRegion
      \ nextgroup=zitAkte

syn match zitMetadateiBezeichnung /\v[^\n]+/ contained
syn match zitMetadateiBezeichnungPrefix /\v^# / contained nextgroup=zitMetadateiBezeichnung
syn region zitMetadateiBezeichnungRegion start=/\v^# / end=/$/ oneline contained contains=zitMetadateiBezeichnungPrefix,zitMetadateiBezeichnung

syn match zitMetadateiEtikett /\v[^\n]+/ contained contains=@NoSpell
syn match zitMetadateiEtikettPrefix /\v^- / contained
syn region zitMetadateiEtikettRegion start=/\v^- / end=/$/ oneline contained contains=zitMetadateiEtikett,zitMetadateiEtikettPrefix

syn match zitMetadateiAkteBase /\v[^\n]*\.@=/ contained contains=@NoSpell nextgroup=zitMetadateiAkteDot
syn match zitMetadateiAkteDot /\v\./ contained contains=@NoSpell nextgroup=zitMetadateiAkteExt
syn match zitMetadateiAkteExt /\v\w+/ contained contains=@NoSpell
syn match zitMetadateiAktePrefix /\v^! / contained nextgroup=zitMetadateiAkteBase
syn region zitMetadateiAkteRegion start=/\v^! / end=/$/ oneline contained contains=zitMetadateiAkte,zitMetadateiAktePrefix,zitMetadateiAkteBase,zitMetadateiAkteExt

" highlight default link zitLinePrefixRoot Special
highlight default link zitMetadatei Normal
highlight default link zitMetadateiBezeichnung Title
highlight default link zitMetadateiEtikett Constant
highlight default link zitMetadateiAkteBase Underlined
highlight default link zitMetadateiAkteExt Type

let b:current_syntax = 'zit.zettel'

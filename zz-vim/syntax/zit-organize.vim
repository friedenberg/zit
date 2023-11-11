
" if exists("b:current_syntax")
"   finish
" endif

let m = expand("<sfile>:h") . "/zit-metadatei.vim"
exec "source " . m

syn match zitEtikett /\v[^#,]+/ contained contains=@NoSpell
syn match zitEtikettPrefix /\v#+/ contained
syn region zitEtikettRegion start=/\v^\s*#+ / end=/$/ oneline
      \ contains=zitEtikett,zitEtikettPrefix

syn match zitZettelHinweis /\v\w+/ contained contains=@NoSpell
syn match zitZettelSeparator /\v\// contained
syn match zitZettelPrefix /\v^\s*- / contained
" don't include the newline because this is within a region
syn match zitZettelBezeichnung /\v.*/ contained contains=@NoSpell

syn region zitZettelHinweisRegion start=/\v\[/ end=/]/ oneline contained
      \ contains=zitZettelHinweis,zitZettelHinweisSeparator
      \ nextgroup=zitZettelBezeichnung

syn region zitZettelRegion start=/\v^\s*- / end=/$/ oneline
      \ contains=zitZettelHinweisRegion ",zitZettelBezeichnung

highlight default link zitEtikett Title
highlight default link zitZettelHinweis Identifier
highlight default link zitZettelBezeichnung String

let b:current_syntax = 'zit-organize'

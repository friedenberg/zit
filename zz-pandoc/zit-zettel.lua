
function Image(el)
  if el.src == nil then
    return nil
  end

  if el.src:find("^zit://") == nil then
    return nil
  end

  local dir_zit = os.getenv("ZIT_DIR")

  if dir_zit == nil then
    error("expected ZIT_DIR env variable to be set")
  end

  local kennung = el.src:sub(7)
  local typ = pandoc.pipe("zit", {"show", "-dir-zit", dir_zit, "-format", "typ", kennung}, "")
  local data = pandoc.pipe("zit", {"show", "-dir-zit", dir_zit, "-format", "akte", kennung}, "")
  local fname = kennung .. "." .. typ
  -- TODO-P1 use mime type from typ
  pandoc.mediabag.insert(fname, "image/png", data)
  return pandoc.Image(el.caption, fname)
end

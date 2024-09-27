
function Image(el)
  if el.src == nil then
    return nil
  end

  -- TODO-P1 switch to using box format instead of url
  if el.src:find("^zit://") == nil then
    return nil
  end

  local dir_zit = os.getenv("ZIT_DIR")

  if dir_zit == nil then
    error("expected ZIT_DIR env variable to be set")
  end

  local objectID = el.src:sub(7)
  local typ = pandoc.pipe("zit", {"show", "-dir-zit", dir_zit, "-format", "type", objectID}, "")
  -- TODO-P1 load MIMEs from type and pick the best one
  local data = pandoc.pipe("zit", {"format-object", "-dir-zit", dir_zit, "-mode", "blob", "png", objectID}, "")
  local fname = objectID .. ".png"
  pandoc.mediabag.insert(fname, "image/png", data)

  return pandoc.Image(el.caption, fname)
end

-- TODO add code formatter
function CodeBlock(el)
  local classes = el.classes

  if #classes < 1 then
    return nil
  end

  local type = classes[1]

  if type:find("^!") == nil then
    return nil
  end

  local dir_zit = os.getenv("ZIT_DIR")

  if dir_zit == nil then
    error("expected ZIT_DIR env variable to be set")
  end

  local mimeGroup = "text"
  local isBinary = FORMAT:find("^markdown") == nil

  if isBinary then
    mimeGroup = "png"
  end

  local data = pandoc.pipe("zit", {"format-object", "-dir-zit", dir_zit, "-stdin", mimeGroup, type}, el.text)

  if isBinary then
    local id = pandoc.utils.sha1(el.text)
    local fname = id .. ".png"
    pandoc.mediabag.insert(fname, "image/png", data)
    return pandoc.Image("", fname)
  else
    el.text = data
    return el
  end
end

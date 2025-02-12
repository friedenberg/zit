-- TODO switch to a socket or lua module that exposes the zit API
-- DirZit = os.getenv("ZIT_DIR")

-- if DirZit == nil then
--   error("expected ZIT_DIR env variable to be set")
-- end

IsBinary = FORMAT:find("^markdown") == nil

function Hex_to_char(x)
  return string.char(tonumber(x, 16))
end

function Unescape(url)
  return url:gsub("%%(%x%x)", Hex_to_char)
end

function FormatObjectImage(imgSrc, format)
  local objectID = Unescape(imgSrc)
  return pandoc.pipe("zit", { "format-object", objectID, format }, ""), objectID
end

-- [fib/chil @b7a8 !task project-2021-zit-bugs zz-inbox] modify pandoc filters to handle non-objects in images
function ReplaceObjectImageWithTextIfNecessary(img)
  if img.src == nil then
    return nil
  end

  local format = "text"

  local data, _ = FormatObjectImage(img.src, format)

  return pandoc.RawInline("markdown", data)
end

-- [fib/chil @b7a8 !task project-2021-zit-bugs zz-inbox] modify pandoc filters to handle non-objects in images
function ReplaceObjectImageWithImageIfNecessary(img)
  if img.src == nil then
    return nil
  end

  local format = "png"

  local data, objectID = FormatObjectImage(img.src, format)

  local id = pandoc.utils.sha1(objectID)
  local fname = id .. "." .. format
  pandoc.mediabag.insert(fname, "image/png", data)
  return pandoc.Image(img.caption, fname)
end

-- [fib/chil @b7a8 !task project-2021-zit-bugs zz-inbox] modify pandoc filters to handle non-objects in images
function IsSku(str)
  return true
end

function UnescapeIfSku(table, key)
  local el = table[key]

  if not IsSku(el) then
    return
  end

  table[key] = Unescape(el)

  return
end

return {
  -- DirZit = DirZit,
  IsBinary = IsBinary,
  Hex_to_char = Hex_to_char,
  Unescape = Unescape,
  FormatObjectImage = FormatObjectImage,
  ReplaceObjectImageWithImageIfNecessary = ReplaceObjectImageWithImageIfNecessary,
  ReplaceObjectImageWithTextIfNecessary = ReplaceObjectImageWithTextIfNecessary,
  IsSku = IsSku,
  UnescapeIfSku = UnescapeIfSku,
}

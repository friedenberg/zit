
DirZit = os.getenv("ZIT_DIR")

if DirZit == nil then
  error("expected ZIT_DIR env variable to be set")
end

IsBinary = FORMAT:find("^markdown") == nil

function Hex_to_char(x)
  return string.char(tonumber(x, 16))
end

function Unescape(url)
  return url:gsub("%%(%x%x)", Hex_to_char)
end

function FormatObjectImage(imgSrc, mimeGroup)
  local objectID = Unescape(imgSrc)
  return pandoc.pipe("zit", { "format-object", "-dir-zit", DirZit, "-mode", "blob", mimeGroup, objectID }, ""), objectID
end

function ReplaceObjectImageWithTextIfNecessary(img)
  if img.src == nil then
    return nil
  end

  local mimeGroup = "text"

  local data, _ = FormatObjectImage(img.src, mimeGroup)

  return pandoc.RawInline("markdown", data)
end

function ReplaceObjectImageWithImageIfNecessary(img)
  if img.src == nil then
    return nil
  end

  local mimeGroup = "png"

  local data, objectID = FormatObjectImage(img.src, mimeGroup)

  local id = pandoc.utils.sha1(objectID)
  local fname = id .. "." .. mimeGroup
  pandoc.mediabag.insert(fname, "image/png", data)
  return pandoc.Image(img.caption, fname)
end

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
  DirZit = DirZit,
  IsBinary = IsBinary,
  Hex_to_char = Hex_to_char,
  Unescape = Unescape,
  FormatObjectImage = FormatObjectImage,
  ReplaceObjectImageWithImageIfNecessary = ReplaceObjectImageWithImageIfNecessary,
  ReplaceObjectImageWithTextIfNecessary = ReplaceObjectImageWithTextIfNecessary,
  IsSku = IsSku,
  UnescapeIfSku = UnescapeIfSku,
}

package.path = package.path .. string.format(";%s/.local/share/pandoc/filters/?.lua", os.getenv("HOME"))

local common = require("zit-common")

if common.IsBinary then
  function Image(img)
    return common.ReplaceObjectImageWithImageIfNecessary(img)
  end
else
  function Figure(el)
    local content = el.content

    if #content ~= 1 then
      return el
    end

    if content[1].t ~= "Plain" then
      return el
    end

    local plain = content[1]
    local image = plain.content[1]
    local replacement = common.ReplaceObjectImageWithTextIfNecessary(image)

    if replacement == nil then
      return el
    else
      return replacement
    end
  end
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

  local mimeGroup = "text"

  if common.IsBinary then
    mimeGroup = "png"
  end

  local data = pandoc.pipe("zit", { "format-object", "-dir-zit", common.DirZit, "-stdin", mimeGroup, type }, el.text)

  if common.IsBinary then
    local id = pandoc.utils.sha1(el.text)
    local fname = id .. ".png"
    pandoc.mediabag.insert(fname, "image/png", data)
    return pandoc.Image("", fname)
  else
    el.text = data
    return el
  end
end
package.path = package.path .. string.format(";%s/.local/share/pandoc/filters/?.lua", os.getenv("HOME"))

local common = require("zit-common")

function CodeBlock(el)
  local classes = el.classes

  if #classes < 1 then
    return nil
  end

  local type = classes[1]

  if type:find("^!") == nil then
    return nil
  end

  local data = pandoc.pipe("zit", { "format-object", "-stdin", type }, el.text)

  el.text = data

  return el
end

function Image(el)
  common.UnescapeIfSku(el, "src")
  return el
end

function Link(el)
  common.UnescapeIfSku(el, "target")
  return el
end

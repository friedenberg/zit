#! /usr/bin/env jq -f

def make_matches:
  [.[] | if (. | type) == "array" then .[] else . end] | join(" ")
  ;

def item:
  {
    title: .Zettel.Bezeichnung,
    subtitle: ([.Hinweis+":"] + .Zettel.Etiketten | join(" ")),
    arg: .Hinweis,
    match: ([.Zettel.Etiketten, .Zettel.AkteExt, .Zettel.Bezeichnung, .Hinweis] | make_matches),
    # icon: {
    #   type: "fileicon",
    #   path: "",
    # },
    type: "file:skipcheck",
    # quicklookurl: "",
    text: {
      copy: .Hinweis,
    },
  }
  ;

(. | item)

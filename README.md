# what is zit

Zit is a combination of "Zettelkasten" and "Git". Zettelkasten is a
knowledge-management method that is used by academics and professionals to graph
their thinking to make it more usable and navigable.

Zit both Zettelkasten and Git as inspiration:

-   everything has an automatically-assigned identifier
-   flat hierarchy (no directories)
-   tags
-   everything is a note or blob
-   every change is tracked and saved

Normally in Zettelkasten, timestamps are used as identifiers. Zit deviates from
this as timestamps are not very ergonomic to type, autocomplete, or
disambiguate. Zit instead uses an identifier system that combines two lists of
user-supplied identifiers to generate unique combinations. Here's an example:

- user supplies two lists
  - list a
    - red
    - green
    - blue
  - list b
    - apple
    - banana
    - orange
  - possible identifiers
    - red/apple
    - red/banana
    - red/orange
    ...
    - blue/orange
- as new zettels are created, new and unused combinations are used to generate
  new identifiers that are then assigned to the zettel

# directory structure

## top-level

I originally wrote this project in Go. However, as the project has grown larger,
I've bumped into some of Go's tradeoffs that have led me to attempt rewriting
small parts of Zit in Rust.

## per-edition

Besides the entrypoints (e.g., main.go, main.rs), each edition uses a strategy
for forcing modularity and unidirectional dependencies like so:

- zit/go/zit
  - main.go
  - src/
    - alfa
      - errors
    - bravo
      - files
    - charlie
      - file_lock

In the above selection, code in the `file_lock` module may depend on anything at
the `alfa` and `bravo` levels. Code in the `files` module may depend on
anything at the `alfa` level. And code in `alfa` may only depend on stdlib or
external modules.

For the code "levels", the nato phonetic alphabet is used. The benefits of this
approach are it encourages modularity and creating a pyramid-like structure of
dependency. For languages like Go, circular dependencies are not allowed and so
this approach prevents hard-to-untangle dependency refactors. For Rust, this
approach may help with build-times.

# how do I use it?

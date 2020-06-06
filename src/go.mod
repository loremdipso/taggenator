module taggenator

go 1.14

require github.com/tidwall/buntdb v1.1.2 // indirect

require internal/data v1.0.0

replace internal/data => ./internal/data

require internal/database v1.0.0

replace internal/database => ./internal/database

require internal/processor v1.0.0

replace internal/processor => ./internal/processor

replace internal/searcher => ./internal/searcher

require (
	github.com/cloudfoundry/bytefmt v0.0.0-20200131002437-cf55d5288a48 // indirect
	github.com/fatih/color v1.9.0
)

require (
	github.com/karrick/godirwalk v1.15.6 // indirect
	github.com/nsf/termbox-go v0.0.0-20200418040025-38ba6e5628f1 // indirect
	internal/searcher v0.0.0-00010101000000-000000000000 // indirect
)

require (
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/loremdipso/fancy_printer v0.0.0-20200605023349-e213bcb72bc1 // indirect
	github.com/loremdipso/liner v1.3.0 // indirect
	internal/go_utils v0.0.0-00010101000000-000000000000 // indirect
	internal/opener v1.0.0
)

replace internal/go_utils => ./internal/go_utils

replace internal/opener => ./internal/opener

# About

Traditional file systems are essentially unbalanced B-trees. This can be a great abstraction for organizing your files when you have few files. However, the more files you have the harder it becomes to find the exact file you want. Maybe it's that .pdf file you know you downloaded last week but forgot exactly which subdirectory it's in, or that .doc you know had the word "final" somewhere in its name. Maybe you remember nothing about the file but know it exists. How do you find it?

## Enter taggenator

Developers will know it's not a new idea to ignore a tree-based file system in favor of fuzzy-searching. If you know some of a file's name it's often much faster to find it using just a few keystrokes and letting the system filter the list of possibilities down for you.

This project attempts to leverage search and filtering to your filesystem when you need it.


# Examples

Find all files you haven't opened yet via taggenator in the current directory down, sort them by their creation date, and start opening them one by one:
```bash
taggenator open unseen -sort newest
```

You can define how you want files to be opened in taggenator_settings.json. Here I have my VLC set to open new videos in the same instance, and I use gnome-open as a sensible default:
```json
	"openerconfig": {
		"default": {
			"Open": "gnome-open \"%s\""
		},

		"video": {
			"Open": "vlc \"%s\"",
			"Update": "vlc \"%s\"",
			"Close": "killall vlc"
		}
		...
	}
```

Filter options include:
```
seen
unseen
	has taggenator opened (seen) this file before?

touched
untouched
	have you added a tag (touched) to this file before?

search
search_inclusive
	very useful. Combination of name/path/tags, though tags don't need to match exactly. Substrings are okay.
	Inclusive means any of the search terms need to match.
	Exclusive means they all do.

name | name_exclusive
name_inclusive | name_includes
	does the name of the file include the search term?

path | path_exclusive
path_includes | path_inclusive
	does the path of the file include the search term?

tags | tags_exclusive
tags_inclusive
	do any of the tags of the file match the search term?
```

Search options include:
```
newest
oldest
	sort by when files were added to taggenator

largest | biggest
smallest
	sort by filesize (isn't updated automatically)

most_recently_opened
least_recently_opened
	sort by the last time taggenator opened each file

most_frequently_opened
least_frequently_opened
	sort by how often taggenator has opened each file

reverse

random
```


# Sub-Modules Organization


The project is split into several sub-modules, some externally maintained, others specific to this project.

# Building
1. Checkout this repository
2. run `go build .`
3. Put the generated binary anywhere you want.


## Requirements
* Go runtime
* Git (obviously)


# Performance
golang was chosen for this because:
1. I care about performance enough to want something compiled with easy concurrency
2. I don't want to manually manage memory.


## Startup Time
While a daemon would be helpful in minimizing startup time, I felt that was something of a cop-out and wanted to get the startup time down just by having a quick database (thank you [buntdb](github.com/tidwall/buntdb)]) and being careful with how I navigate the filesystem (thanks [godirwalk](github.com/karrick/godirwalk)).

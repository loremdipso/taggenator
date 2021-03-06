package main 

const (
TaggenatorSettings = `{
	"extensions": [".txt"],

	"synonyms": {
		"pre": "prefix",
		"post": "postfix"
	},

	"prefixes": [
		"getmore:"
	],

	"commands": {
		"ls": "ls"
	},

	"tagger": {
		"test": "echo \"new tag1\nnew tag2\n%s\""
	},

	"openerconfig": {
		"default": {
			"Open": "gnome-open \"%s\""
		},

		"video": {
			"Open": "vlc \"%s\"",
			"Update": "vlc \"%s\"",
			"Close": "killall vlc"
		},

		"comic": {
			"Open": "mcomix \"%s\"",
			"Close": "killall mcomix"
		}
	}
}`
)

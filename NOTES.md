## NOTES

### File Walking

- I'm using the godirwalk project for faster file traversal, since go's default will also perform a STAT, which we only need for new files. Under semi-extreme testing (1M files split between 1000 folders):
	- if all the files are new (meaning we need to do a stat regardless) the github repo took 9.3s to process, whereas the standard walk took only 7.6s
	- if none of the files are new (meaning a stat is unnecessary) then the github repo only takes 1.0s to process, whereas the standard walk takes 2.26s

	So, we will use the repo for updates, and might use the standard walk for brand-new databases, if I feel that optimization is worth the added complexity.

	UPDATE: found a place to put in a coroutine that speeds up the whole initialization process by, like, 40%, so I took out the special-case option since that didn't seem to be adding much. The only hard part is deciding how big the filename cache should be. I think that'll take some fiddling, perhaps some run-time monitoring to auto-adjust/disable entirely if filesystems get to be so fast that it just isn't worth it.

	Currently, 100K files takes 600ms to initialize from nothing and 60-70ms to just iterate over. I think that's plenty fast.


## TODOs

- This does not handle renaming files !!!

- This doesn't handle multiple processes. Come up with something for that.
	- Maybe make subsequent processes read-only, or do some sort of cool daemon and support some inter-process-communication?

- Disable logging on release, maybe?
	- log.SetFlags(0)

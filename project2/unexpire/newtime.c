/*
 * Name: Christopher Naporlee
 * netID: cmn134
 * RUID: 187008361
 */
#define _GNU_SOURCE
#include <dlfcn.h>
#include <err.h>
#include <stdio.h>
#include <stdlib.h>
#include <time.h>

/* Keep date in the format month-day-year hour-min-sec */
#define DATE_TO_USE "07-31-2021 00:00:00"

/*
 * Malicious time()
 * The program we are lying to about the time expects @t and the return value
 * to be the exact same!
 */
time_t time(time_t *t)
{
	static int should_lie = 1;
	time_t ret;

	/* Only on our first run should we lie about the time */
	if (should_lie) {
		struct tm tm;

		if (!strptime(DATE_TO_USE, "%m-%d-%Y %H:%M:%S", &tm))
			err(1, "newtime: strptime() fail at line (%lu)", __LINE__);
		ret = mktime(&tm);
		if (ret == (time_t) -1)
			err(1, "Failed to mktime()");
		if (t)
			*t = ret;
		should_lie = 0;
	} else {
		time_t (*real_time)(time_t *) = dlsym(RTLD_NEXT, "time");
		if (!real_time)
			err(1, "dlsym fail(). Could not get original time()");
		ret = real_time(t);
	}
	return ret;
}

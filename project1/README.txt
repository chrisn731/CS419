Quick Start
===========

To create the authentication program use ``make``. You can then run the program
as expected.
	e.g. ./auth AddUser username password

To clean out the directory of all generated content, use ``make clean``.
To clean out just the generated tables, use ``make cleantables``.


Notes
=====

This program does not rely on any non standard Go libraries.

Tables that contain important user, domain, and type information are all
generated when they are not detected within the current directory. Therefore,
no user action is needed to set up the tables.


Tests Used
==========

Sample Test with output:

	./auth AddUser tim pass123
		Success

	./auth AddUser kayla secretpass
		Success

	./auth SetDomain kayla admin
		Success

	./auth SetDomain tim moderator
		Success

	./auth SetType voiceserver servers
		Success

	./auth SetType whitelist user_lists
		Success

	./auth AddAccess ban admin servers
		Success

	./auth AddAccess kick moderator servers
		Success

	./auth AddAccess add_user admin user_lists
		Success

	./auth CanAccess ban kayla voiceserver
		Success

	./auth CanAccess ban tim voiceserver
		Error: access denied

	./auth CanAccess add_user tim whitelist
		Error: access denied

	./auth DomainInfo admin
		kayla

	./auth TypeInfo user_lists
		whitelist


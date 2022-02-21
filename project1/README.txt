Quick Start
===========

To create the authentication program use ``make``. You can then run the program
as expected.
	e.g. ./auth AddUser username password

To clean out the directory of all generated content, use ``make clean``.
To clean out just the generated tables, use ``make cleantables``.


Notes
=====

This program does not rely on any non standard Go libraries. Only thing that
is required is the golang package to build it. (This is already installed
on iLab!)

Tables that contain important user, domain, and type information are all
generated when they are not detected within the current directory. Therefore,
no user action is needed to set up the tables.


Tests Used (tested on snow.cs.rutgers.edu)
==========================================

Automated Test (from course website):
	// Ensure we have a clean slate
	make clean
	make

	// Builds 50 users
	for i in `seq 50`;do ./auth AddUser user-$i password-$i;done
		[Returns 50 successes]

	// Authenticates they are in the system
	for i in `seq 50`;do ./auth Authenticate user-$i password-$i;done
		[Returns 50 successes]

	// Authenticates we can not use bad passwords
	for i in `seq 50`;do ./auth Authenticate user-$i bad-password-$i;done
		[Returns 50 "Error: bad password"]

	// Creates 50 unique servers within the "servers" type
	for i in `seq 50`; do ./auth SetType server-$i servers; done
		[Returns 50 "Success"]

	// Validate they are stored
	./auth TypeInfo servers
		[Returns 50 "servers" labeled server-x]

	// Set all the users we created to be of basic type "user"
	for i in `seq 50`;do ./auth SetDomain user-$i user; done
		[Returns 50 "Success"]

	// Set half of those users to be "engineers"
	for i in `seq 25`;do ./auth SetDomain user-$i engineer; done
		[Returns 25 "Success"]

	./auth DomainInfo engineer
		[Returns 25 users from user-1 to user-25]

	./auth DomainInfo user
		[Returns 50 users from user-1 to user-50]

	// Allow all users to login to the servers
	for i in `seq 50`; do ./auth AddAccess login user servers; done
		[Returns 50 "Success"]

	// Allow all engineers to make edits to the server type
	for i in `seq 25`; do ./auth AddAccess edit engineer servers; done
		[Returns 25 "Success"]

	// Ensure that all users can access the server by loggin in
	for i in `seq 50`; do ./auth CanAccess login user-$i server-1; done
		[Returns 50 "Success"]

	// Ensure that 25 users can make edits to the server
	// Ensure that the other 25 users are denied editing the server
	for i in `seq 50`; do ./auth CanAccess edit user-$i server-1; done
		[Returns 25 "Success" and 25 "Error: access denied"]



Manual Test Trace with output:

	./auth AddUser tim pass123
		Success

	./auth AddUser kayla secretpass
		Success

	./auth Authenticate tim pass123
		Success

	./auth SetDomain kayla admin
		Success

	./auth Authenticate kayla admin_password123
		Error: bad password

	./auth Authenticate root rootpass
		Error: no such user

	./auth SetDomain tim moderator
		Success

	./auth SetDomain tim ""
		Error: missing domain

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


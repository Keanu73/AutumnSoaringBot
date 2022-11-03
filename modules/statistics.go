package modules

/*
stats collection

log:
* everyone who joined a session (for review)
* number of sessions that have actually taken place
* length of sessions aka. amount of time people spend in VC. (AFTER 10 mins exercise, 5 mins meditation & journalling)
*/

/*
we'll use a central collection: db.configuration
configuration will hold:
SESSION_START_TIME (default: 6:30:00)
ACTIVITY_ROLE_ID (default: <insert ID here> - ID of "active" role)
*/

/* db.user_statistics:
1 document:
id: 115156616256552962
lastSessionJoined: 2022-10-29
sessionsJoined: 90

db.call_statistics:
date: 2022-10-30
startedAt: "06:25:00" (time when first participant joins)
finished: false
participants: [{
	id: 115156616256552962
	joinedVC: "06:30:00"
	leftVC: ""
}]
*/

/* on voice state update:
Check: do they already have a document?
If so, use existing one
If not, create new one under "statistics" collection

1) if a new person joins THE VC (BOT_VC_ID)
* Log the time they joined the VC. (e.g. 2022-10-28 06:30:00)

2) if a person leaves THE VC
* Log the time they left the VC
* If endTime - startTime < 10, don't bother logging.

*/

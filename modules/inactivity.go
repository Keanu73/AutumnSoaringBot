package modules

/*
inactivity module

let's use a cron for this. every day at midnight,
we'll check how long each of the members have been joining sessions in the discord for.
IF:
* they have not joined a single session for 6 weeks:
then:
* kick them from the server.

we'll aim to give dm notifications in advance - at certain intervals, namely
if someone hasn't attended a session for:
* 1 week (7 days)
* 2 weeks (14 days)
* 4 weeks (
* 5 weeks
* 3 days before 6th week
* 24 hours before 6th week

exempt people with "active" role

Use mongo stats DB to keep track of this
*/

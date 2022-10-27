package main

/*
inactivity module

let's use a cron for this. every day at midnight, we'll check how long each of the members have been in the discord for.
IF:
* they have not joined a single session for 6 weeks:
then:
* kick them from the server.

Use postgres stats DB to keep track of this
*/

# go-bratok #

## System for starting, control and synchronization programs on different servers without main central server. ##

### Introduction ###

No yet

## DEVELOPER RELEASE. ##

### REST API. ###


### 
There are 3 kind of servers:

- master

- regular

- dynamic


"master" and "regular" servers have the one different: if we have changed config on different servers in the same time
when "master" has an advantage over "regular" server.

"dynamic" servers cannot process request for config update.

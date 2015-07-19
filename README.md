MySQL Facility in Go
===================

* Tools
	* myproxy
		* Proxy for MySQL-Server in Go
		* [ ] Packet proxy
		* [ ] Binlog proxy
	* mymon
		* MySQL monitor by packet capture
		* [ ] Monitor packet
		* [ ] Monitor binlog
* Libraries
	* proto
		* Speak MySQL protocol in go way
		* [x] Read packet
		* [x] Write packet
	* binlog
		* MySQL binlog replication protocol in go
		* [x] Read binary log
		* [ ] Write binary log
		* [ ] Semi-Synchronization
		* [ ] Binlog client
		* [ ] Binlog server
	* [ ] server
		* MySQL server
	* [ ] index
		* Secondary index creator by repl
		* Store in redis,ledis
	* [ ] driver
		* Yet another mysql driver implement database/sql

Motivation
==========
* Learning go
* Thinking in go
* Trying to go
* Understand MySQL protocol

See Also
========
* [JSS - Java sql server]()
	* Same as myfacility, with a bad design, not complete.

USE mojo;
INSERT INTO People (FirstName,LastName,Email1) VALUES
	("Billy","Roberts","billybob@xyz.com"),
	("Lemar","Jones","lemarj@bubba.com"),
	("Ebb","Tide","ebb@tidetunes.com"),
	("Tongue","Tide","tee@tidetunes.com"),
	("Sally","Jones","sallybob@foo.com");

INSERT INTO People (FirstName,LastName,Email1,Status,OptOutDate) VALUES
	("Wilma","Whiner","wilma@gripe.com",1,"2016-02-14"),
	("Wendy","Whiner","wendy@omg.com",1,"2016-07-04");

INSERT INTO EGroup (GroupName) VALUES ("FAA"),("Isola");

INSERT INTO PGroup (PID,GID) VALUES (1,1), (2,1), (3,1), (4,1), (5,1), (6,1), (7,1), (5,2), (7,2);

INSERT INTO DataUpdate (GID,DtStart,DtStop) VALUES 
	(1,"2017-04-03 00:08:12", "2017-04-03 00:09:14"),
	(2,"2017-04-02 00:12:10", "2017-04-03 00:12:57");

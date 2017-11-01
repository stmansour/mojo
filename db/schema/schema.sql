--  FAA personnel database

DROP DATABASE IF EXISTS mojo;
CREATE DATABASE mojo;
USE mojo;
GRANT ALL PRIVILEGES ON mojo TO 'ec2-user'@'localhost';
GRANT ALL PRIVILEGES ON mojo.* TO 'ec2-user'@'localhost';
GRANT ALL PRIVILEGES ON mojo TO 'adbuser'@'localhost';
GRANT ALL PRIVILEGES ON mojo.* TO 'adbuser'@'localhost';
set GLOBAL sql_mode='ALLOW_INVALID_DATES';

-- **************************************
-- ****                              ****
-- ****        MOJO  DATABASE        ****
-- ****                              ****
-- **************************************

CREATE TABLE People (
    PID BIGINT NOT NULL AUTO_INCREMENT,                         -- person id
    FirstName VARCHAR(100) DEFAULT '',
    MiddleName VARCHAR(100) DEFAULT '',
    LastName VARCHAR(100) DEFAULT '',
    PreferredName VARCHAR(100) DEFAULT '',
    JobTitle VARCHAR(100) DEFAULT '',
    OfficePhone VARCHAR(100) DEFAULT '',
    OfficeFax VARCHAR(100) DEFAULT '',
    Email1 VARCHAR(50) DEFAULT '',
    Email2 VARCHAR(50) NOT NULL DEFAULT '',
    MailAddress VARCHAR(50) DEFAULT '',
    MailAddress2 VARCHAR(50) DEFAULT '',
    MailCity VARCHAR(100) DEFAULT '',
    MailState VARCHAR(50) DEFAULT '',
    MailPostalCode VARCHAR(50) DEFAULT '',
    MailCountry VARCHAR(50) DEFAULT '',
    RoomNumber VARCHAR(50) DEFAULT '',
    MailStop VARCHAR(100) DEFAULT '',
    Status SMALLINT DEFAULT 0,      -- 0 = ok, 1 = they've opted out, 2 = address bounced, 3 = opt out via complaint
    OptOutDate DATETIME NOT NULL DEFAULT '1970-01-01 00:00:00',     -- if State is 1, the date/time when the person opted out
    LastModTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    LastModBy BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (PID)     
);      

CREATE TABLE EGroup (
    GID BIGINT NOT NULL AUTO_INCREMENT,                         -- Group ID
    GroupName VARCHAR(50) NOT NULL DEFAULT '',                  -- Name of the group
    GroupDescription VARCHAR(1000) NOT NULL DEFAULT '',         -- Description of the group
    DtStart DATETIME NOT NULL DEFAULT '1970-01-01 00:00:00',    -- start time of last scrape
    DtStop DATETIME NOT NULL DEFAULT '1970-01-01 00:00:00',     -- stop time of last scrap
    LastModTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    LastModBy BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (GID)     
);

CREATE TABLE PGroup (
    PID BIGINT NOT NULL DEFAULT 0,                              -- person id
    GID BIGINT NOT NULL DEFAULT 0,                              -- group
    LastModTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    LastModBy BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE DataUpdate (
    DUID BIGINT NOT NULL AUTO_INCREMENT,
    GID BIGINT NOT NULL DEFAULT 0,                              -- group
    DtStart DATETIME NOT NULL DEFAULT '1970-01-01 00:00:00',
    DtStop DATETIME NOT NULL DEFAULT '1970-01-01 00:00:00',
    LastModTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    LastModBy BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY(DUID)
);

CREATE TABLE Query (
    QID BIGINT NOT NULL AUTO_INCREMENT,
    QueryName VARCHAR(50) DEFAULT '',
    QueryDescr VARCHAR(1000) DEFAULT '',
    QueryJSON VARCHAR(3000) DEFAULT '',
    LastModTime TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    LastModBy BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY(QID)
);
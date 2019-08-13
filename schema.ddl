CREATE TABLE Users (
	UserId STRING(MAX) NOT NULL,
	Email STRING(MAX) NOT NULL,
	UserName STRING(MAX) NOT NULL,
	Created TIMESTAMP NOT NULL,
	Updated TIMESTAMP NOT NULL,
	Picture STRING(MAX) NOT NULL,
) PRIMARY KEY (UserId);

CREATE TABLE Queries (
	UserId STRING(MAX) NOT NULL,
    QueryId STRING(MAX) NOT NULL,
	Created TIMESTAMP NOT NULL,
	ImageUrl STRING(MAX) NOT NULL,
	Result STRING(MAX),
) PRIMARY KEY (UserId, QueryId),
  INTERLEAVE IN PARENT Users ON DELETE CASCADE;

CREATE TABLE Sessions (
	SessionId STRING(MAX) NOT NULL,
    UserId STRING(MAX) NOT NULL,
	UserCount INT64 NOT NULL,
) PRIMARY KEY (SessionId);
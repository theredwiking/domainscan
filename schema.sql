CREATE TABLE domain(
	id INTEGER PRIMARY KEY,
	name TEXT,
	ip TEXT,
	protocol TEXT,
	server TEXT
);

CREATE UNIQUE INDEX idx_domain_name ON domain (name);

CREATE TABLE nmap(
	id INTEGER PRIMARY KEY,
	port INTEGER UNIQUE,
	protocol TEXT,
	service TEXT
);

CREATE TABLE domainNmap(
	id INTEGER PRIMARY KEY,
	domainId INTEGER,
	nmapId INTEGER,
	state TEXT,
	FOREIGN KEY (domainId) REFERENCES domain(id),
	FOREIGN KEY (nmapId) REFERENCES nmap(id)
);

CREATE INDEX idx_domain_id ON domainNmap (domainId);

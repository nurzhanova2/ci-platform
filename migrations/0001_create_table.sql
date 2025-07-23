CREATE TABLE pipeline_runs (
	id SERIAL PRIMARY KEY,
	repo TEXT NOT NULL,
	status TEXT NOT NULL,
	logs TEXT,
	started_at TIMESTAMP,
	finished_at TIMESTAMP
);

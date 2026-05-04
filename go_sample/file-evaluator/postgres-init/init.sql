CREATE SEQUENCE eval_id;

CREATE TABLE file_evaluation (
    id INTEGER DEFAULT nextval('eval_id'),
    status TEXT,
    result boolean,
    request_ts timestamp DEFAULT now(), 
    result_ts timestamp 
);
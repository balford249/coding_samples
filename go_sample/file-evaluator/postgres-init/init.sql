CREATE SEQUENCE eval_id;

CREATE TABLE file_evaluation (
    id INTEGER DEFAULT nextval('eval_id'),
    status TEXT,
        result boolean 
);
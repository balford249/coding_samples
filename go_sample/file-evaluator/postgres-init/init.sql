CREATE TABLE file_evaluation_type (
    type TEXT PRIMARY KEY
);

CREATE SEQUENCE eval_id_seq;

CREATE TABLE file_evaluation_event (
    id INTEGER DEFAULT nextval('eval_id_seq'),
    request_ts TIMESTAMP DEFAULT now(),
    PRIMARY KEY(id)
);

CREATE TABLE file_evaluation_result (
    eval_id INTEGER,
    type TEXT,
    status TEXT,
    result_ts TIMESTAMP DEFAULT now(),
    CONSTRAINT eval_type_fk
        FOREIGN KEY(type)
            REFERENCES file_evaluation_type(type),
    CONSTRAINT eval_id_fk
        FOREIGN KEY(eval_id)
            REFERENCES file_evaluation_event(id),
    CONSTRAINT eval_status CHECK(status IN ('pending', 'passed', 'failed'))
);

INSERT INTO file_evaluation_type VALUES ('FileExists');
INSERT INTO file_evaluation_type VALUES ('IsTxtFile');

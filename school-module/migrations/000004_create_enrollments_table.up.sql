-- CREATE enrollments TABLE
CREATE TABLE IF NOT EXISTS enrollments (
    student_id   UUID           NOT NULL,
    course_id    UUID           NOT NULL,
    installments INT2           NOT NULL,
    status       TEXT           NOT NULL,
    created_at   TIMESTAMP      NOT NULL DEFAULT NOW(),
    CONSTRAINT enrollments_pk PRIMARY KEY (student_id, course_id),
    CONSTRAINT enrollments_students_fk FOREIGN KEY (student_id) REFERENCES students (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT enrollments_courses_fk FOREIGN KEY (course_id) REFERENCES courses (id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- ADD INDEX TO enrollments_student_id_students_fk
CREATE INDEX IF NOT EXISTS enrollments_student_id_idx
ON enrollments
USING btree (student_id);

CREATE TABLE users
(
    id         INT AUTO_INCREMENT PRIMARY KEY,
    name       VARCHAR(255)        NOT NULL,
    email      VARCHAR(255) UNIQUE NOT NULL,
    created_at DATETIME            NOT NULL,
    updated_at DATETIME            NOT NULL
);

INSERT INTO users (name, email, created_at, updated_at) VALUES ('Alice Smith', 'alice.smith@example.com', '2023-01-10 10:00:00', '2023-01-10 10:00:00');
INSERT INTO users (name, email, created_at, updated_at) VALUES ('Bob Johnson', 'bob.j@example.com', '2023-02-15 11:30:00', '2023-02-15 11:30:00');
INSERT INTO users (name, email, created_at, updated_at) VALUES ('Charlie Brown', 'charlie.b@example.com', '2023-03-20 14:45:00', '2023-03-20 14:45:00');
INSERT INTO users (name, email, created_at, updated_at) VALUES ('Diana Prince', 'diana.p@example.com', '2023-04-01 09:15:00', '2023-04-01 09:15:00');
INSERT INTO users (name, email, created_at, updated_at) VALUES ('Eve Adams', 'eve.a@example.com', '2023-05-05 16:00:00', '2023-05-05 16:00:00');
INSERT INTO users (name, email, created_at, updated_at) VALUES ('Frank White', 'frank.w@example.com', '2023-06-10 08:00:00', '2023-06-10 08:00:00');
INSERT INTO users (name, email, created_at, updated_at) VALUES ('Grace Taylor', 'grace.t@example.com', '2023-07-22 13:00:00', '2023-07-22 13:00:00');
INSERT INTO users (name, email, created_at, updated_at) VALUES ('Henry Green', 'henry.g@example.com', '2023-08-30 17:00:00', '2023-08-30 17:00:00');
INSERT INTO users (name, email, created_at, updated_at) VALUES ('Ivy King', 'ivy.k@example.com', '2023-09-12 10:10:00', '2023-09-12 10:10:00');
INSERT INTO users (name, email, created_at, updated_at) VALUES ('Jack Lee', 'jack.l@example.com', '2023-10-25 12:20:00', '2023-10-25 12:20:00');
INSERT INTO users (name, email, created_at, updated_at) VALUES ('Karen Hall', 'karen.h@example.com', '2023-11-01 09:00:00', '2023-11-01 09:00:00');
INSERT INTO users (name, email, created_at, updated_at) VALUES ('Liam Scott', 'liam.s@example.com', '2023-12-05 15:30:00', '2023-12-05 15:30:00');
-- roles
CREATE TABLE roles (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(100) NOT NULL,
    type_role  VARCHAR(100) NOT NULL
);

-- districts
CREATE TABLE districts (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(255) NOT NULL
);


CREATE TABLE ter_otdels (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL,
    district_id UUID NOT NULL REFERENCES districts(id) ON DELETE RESTRICT
);

CREATE INDEX idx_ter_otdels_district_id ON ter_otdels(district_id);

-- users
CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name     VARCHAR(255) NOT NULL,
    login         VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role_id       UUID NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
    district_id   UUID NOT NULL REFERENCES districts(id) ON DELETE RESTRICT,
    ter_otdel_id  UUID REFERENCES ter_otdels(id) ON DELETE SET NULL
);

CREATE INDEX idx_users_role_id      ON users(role_id);
CREATE INDEX idx_users_district_id  ON users(district_id);
CREATE INDEX idx_users_ter_otdel_id ON users(ter_otdel_id);



CREATE TABLE applicants (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name     VARCHAR(255) NOT NULL,
    position      VARCHAR(255) NOT NULL,
    district_id   UUID NOT NULL REFERENCES districts(id) ON DELETE SET NULL,
    base          BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_applicants_district_id ON applicants(district_id);

-- external_users
CREATE TABLE external_users (
    id            SERIAL PRIMARY KEY,
    username      VARCHAR(50) NOT NULL,
    email         VARCHAR(100) NOT NULL,
    first_name    VARCHAR(100) NOT NULL,
    last_name     VARCHAR(100) NOT NULL,
    phone_number  VARCHAR(100),
    city          VARCHAR(100),
    date_of_birth DATE,
    bio           VARCHAR(200),
    password      VARCHAR(255) NOT NULL,
    is_accepted   BOOLEAN NOT NULL,
    patronymic    VARCHAR(100)
);

-- requests
CREATE TABLE requests (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    district_id    UUID REFERENCES districts(id) ON DELETE SET NULL,
    ter_otdel_id   UUID REFERENCES ter_otdels(id) ON DELETE SET NULL,
    applicant_id   UUID REFERENCES applicants(id) ON DELETE SET NULL,
    address        TEXT,
    dogs_count     INTEGER,
    behavior       TEXT,
    urgency        TEXT,
    contact_person TEXT,
    status         TEXT NOT NULL DEFAULT 'Новая'::text,
    number         INTEGER,
    year           INTEGER NOT NULL DEFAULT EXTRACT(year FROM now()),
    act_file       TEXT,
    created_at     TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_requests_ter_otdel_id ON requests(ter_otdel_id);
CREATE INDEX idx_requests_applicant_id ON requests(applicant_id);
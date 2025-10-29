CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name TEXT NOT NULL,
    role TEXT NOT NULL,
    login TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL
);
CREATE TABLE request_sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE applicants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL
);

CREATE TABLE requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_id UUID REFERENCES request_sources(id) ON DELETE RESTRICT,
    applicant_id UUID REFERENCES applicants(id) ON DELETE SET NULL,
    address TEXT NOT NULL,
    dogs_count INT CHECK (dogs_count >= 1),
    behavior TEXT NOT NULL,
    urgency TEXT NOT NULL,
    contact_person TEXT,
    status TEXT DEFAULT 'Новая',
    number INT NOT NULL,  -- порядковый номер в году
    year INT NOT NULL DEFAULT EXTRACT(YEAR FROM now()),
    created_at TIMESTAMP DEFAULT now(),
    UNIQUE (number, year) -- чтобы номера не повторялись внутри года
);

CREATE TABLE request_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id UUID REFERENCES requests(id) ON DELETE CASCADE,
    action TEXT NOT NULL,
    performed_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE catch_acts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id UUID REFERENCES requests(id) ON DELETE CASCADE,
    file_path TEXT NOT NULL,
    uploaded_by UUID REFERENCES users(id) ON DELETE SET NULL,
    uploaded_at TIMESTAMP DEFAULT now()
);

CREATE TABLE contractor_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE contractor_request_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contractor_request_id UUID REFERENCES contractor_requests(id) ON DELETE CASCADE,
    request_id UUID REFERENCES requests(id) ON DELETE CASCADE
);






CREATE OR REPLACE FUNCTION assign_request_number()
RETURNS TRIGGER AS $$
DECLARE
    next_number INT;
BEGIN
    -- Определяем следующий номер в текущем году
    SELECT COALESCE(MAX(number), 0) + 1
    INTO next_number
    FROM requests
    WHERE year = EXTRACT(YEAR FROM now());

    NEW.year := EXTRACT(YEAR FROM now());
    NEW.number := next_number;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Триггер на вставку
CREATE TRIGGER trg_assign_request_number
BEFORE INSERT ON requests
FOR EACH ROW
EXECUTE FUNCTION assign_request_number();

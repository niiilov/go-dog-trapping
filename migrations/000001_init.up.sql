
-- Территориальные отделы
CREATE TABLE request_sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE
);

-- Юзеры
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name TEXT NOT NULL,
    role TEXT NOT NULL,
    login TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    source_id UUID REFERENCES request_sources(id) ON DELETE RESTRICT -- NULL для админов/подрядчиков
);



-- Заявители
CREATE TABLE applicants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    is_permanent BOOLEAN DEFAULT true -- true = в списке, false = разовый
);

CREATE INDEX idx_applicants_permanent ON applicants(is_permanent) WHERE is_permanent = true;

-- Заявки
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
    number INT NOT NULL,
    year INT NOT NULL DEFAULT EXTRACT(YEAR FROM now()),
    created_at TIMESTAMP DEFAULT now(),
    UNIQUE (number, year)
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

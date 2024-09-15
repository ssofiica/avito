CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS employee (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
);

CREATE TABLE IF NOT EXISTS organization (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type organization_type,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS organization_responsible (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    user_id UUID REFERENCES employee(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tender (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name text UNIQUE NOT NULL,
    description text,
    service_type text CHECK (service_type IN ('construction', 'delivery', 'manufacture')),
	status text CHECK (status IN ('created', 'closed', 'published')),
	organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
	creator_username text,
    version integer,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS bid (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name text UNIQUE NOT NULL,
    description text,
	status text CHECK (status IN ('created', 'canceled', 'published', 'approved', 'rejected')),
	author_type text CHECK (author_type IN ('organization', 'user')),
    author_id UUID REFERENCES employee(id) ON DELETE CASCADE,
    tender_id UUID REFERENCES tender(id) ON DELETE CASCADE,
    version integer,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

insert into employee (id, username, first_name, last_name) values ('1c2bb1bd-4d36-4d1d-8b3d-e85a603c0f83', 'ssofiica', 'София', 'Валова');
insert into organization (id, name, type) values ('90c058c5-e03a-4d4e-9817-9f0d3eb7e1cd','Пиццерия', 'IE');
insert into organization_responsible (organization_id, user_id) values ('90c058c5-e03a-4d4e-9817-9f0d3eb7e1cd', '1c2bb1bd-4d36-4d1d-8b3d-e85a603c0f83');
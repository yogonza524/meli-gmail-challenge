GRANT ALL PRIVILEGES ON DATABASE meli TO meli_user;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE challenge(
    m_id UUID NOT NULL DEFAULT uuid_generate_v1(),
    fecha character varying(128),
    from_ character varying(256),
    subject character varying(256)
)
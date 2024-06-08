-- REQUIRED TABLES

create table if not exists company (
    id uuid primary key default gen_random_uuid(),
    ruc text not null,
    name text not null,
    employees smallint default 1,
    is_active boolean default true,
    created_at timestamp with time zone default now(),

    unique (ruc),
    unique (name)
);

create table if not exists role (
    id char(2) primary key,
    name varchar(255) not null,
    created_at timestamp with time zone default now()
);

create table if not exists "user" (
    id uuid primary key default gen_random_uuid(),
    email text not null,
    name text not null,
    password text not null,
    created_at timestamp with time zone default now(),

    company_id uuid not null references company (id) on delete restrict,
    role_id char(2) references role (id) on delete restrict,

    unique (email)
);

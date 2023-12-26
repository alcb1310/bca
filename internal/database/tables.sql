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

create table if not exists logged_in (
    user_id uuid primary key references "user" (id) on delete restrict,
    token text not null,
    created_at timestamp with time zone default now()
);

create table if not exists project (
    id uuid primary key default gen_random_uuid(),
    name text not null,
    is_active boolean default true,

    company_id uuid not null references company (id) on delete restrict,
    created_at timestamp with time zone default now(),

    unique (name, company_id)
);

create table if not exists supplier (
    id uuid primary key default gen_random_uuid(),
    supplier_id text not null,
    name text not null,

    contact_name text,
    contact_email text,
    contact_phone text,

    company_id uuid not null references company (id) on delete restrict,
    created_at timestamp with time zone default now(),

    unique (supplier_id, company_id),
    unique (name, company_id)
);

create table if not exists budget_item (
    id uuid primary key default gen_random_uuid(),
    code text not null,
    name text not null,
    level smallint not null default 1,
    accumulate boolean not null default true,

    parent_id uuid references budget_item (id) on delete restrict,

    company_id uuid not null references company (id) on delete restrict,
    created_at timestamp with time zone default now(),

    unique (code, company_id),
    unique (name, company_id)
);

create table if not exists budget (
    project_id uuid not null references project (id) on delete restrict,
    buget_item_id uuid not null references budget_item (id) on delete restrict,

    initial_quantity numeric,
    initial_cost numeric,
    initial_total numeric not null,

    spent_quantity numeric,
    spent_total numeric not null,

    remaining_quantity numeric,
    remaining_cost numeric,
    remaining_total numeric not null,

    updated_budget numeric not null,

    company_id uuid not null references company (id) on delete restrict,
    created_at timestamp with time zone default now(),

    unique (project_id, buget_item_id, company_id),
    primary key (project_id, buget_item_id, company_id)
);

-- VIEWS

create or replace view vw_budget_item as
select
    b.id,
    b.code,
    b.name,
    b.level,
    b.accumulate,
    p.id as parent_id,
    p.code as parent_code,
    p.name as parent_name,
    b.company_id as company_id
from budget_item as b
left join budget_item as p on b.parent_id = p.id;

-- create table if not exists budget (
--)

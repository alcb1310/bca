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
    budget_item_id uuid not null references budget_item (id) on delete restrict,

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

    unique (project_id, budget_item_id, company_id),
    primary key (project_id, budget_item_id, company_id)
);

create table if not exists invoice (
     id uuid primary key default gen_random_uuid(),
     supplier_id uuid not null references supplier (id) on delete restrict,
     project_id uuid not null references project (id) on delete restrict,
     invoice_number text not null,
     invoice_date date not null,
     invoice_total numeric not null default 0,

     company_id uuid not null references company (id) on delete restrict,
     created_at timestamp with time zone default now(),

     unique (supplier_id, project_id, invoice_number, company_id)
);

create table if not exists invoice_details (
     invoice_id uuid not null references invoice (id) on delete restrict,
     budget_item_id uuid not null references budget_item (id) on delete restrict,
     quantity numeric not null,
     cost numeric not null,
     total numeric not null,

     company_id uuid not null references company (id) on delete restrict,
     created_at timestamp with time zone default now(),

     unique (invoice_id, budget_item_id, company_id),
     primary key (invoice_id, budget_item_id, company_id)
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

create or replace view vw_budget as
select
    bi.id as budget_item_id,
    bi.code as budget_item_code,
    bi.name as budget_item_name,
    bi.level as budget_item_level,
    bi.accumulate as budget_item_accumulate,
    p.id as project_id,
    p.name as project_name,
    b.initial_quantity,
    b.initial_cost,
    b.initial_total,
    b.spent_quantity,
    b.spent_total,
    b.remaining_quantity,
    b.remaining_cost,
    b.remaining_total,
    b.updated_budget,
    b.company_id as company_id
from budget as b
join project as p on b.project_id = p.id
join budget_item as bi on b.budget_item_id = bi.id;

create or replace view vw_invoice as
select
     i.id,
     s.id as supplier_id,
     s.supplier_id as supplier_number,
     s.name as supplier_name,
     s.contact_name as supplier_contact_name,
     s.contact_email as supplier_contact_email,
     s.contact_phone as supplier_contact_phone,
     p.id as project_id,
     p.name as project_name,
     p.is_active as project_is_active,
     i.invoice_number,
     i.invoice_date,
     i.invoice_total,
     i.company_id as company_id
from invoice i
join supplier s on i.supplier_id = s.id
join project p on i.project_id = p.id;

create or replace view vw_invoice_details as 
select
vi.id as invoice_id,
vi.supplier_id as supplier_id,
vi.supplier_number as supplier_number,
vi.supplier_name as supplier_name,
vi.supplier_contact_name as supplier_contact_name,
vi.supplier_contact_email as supplier_contact_email,
vi.supplier_contact_phone as supplier_contact_phone,
vi.project_id as project_id,
vi.project_name as project_name,
vi.project_is_active as project_is_active,
vi.invoice_number,
vi.invoice_date,
vi.invoice_total,
bi.id as budget_item_id,
bi.code as budget_item_code,
bi.name as budget_item_name,
bi.level as budget_item_level,
bi.accumulate as budget_item_accumulate,
bi.parent_id as budget_item_parent_id,
id.quantity,
id.cost,
id.total,
id.company_id as company_id
from invoice_details id
join vw_invoice vi on id.invoice_id = vi.id
join budget_item bi on id.budget_item_id = bi.id;


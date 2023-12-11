-- REQUIRED TABLES

create table if not exists company (
     id uuid PRIMARY KEY  default gen_random_uuid(),
     ruc text not null,
     name text not null,
     employees smallint  default 1,
     is_active boolean  default true,
     created_at timestamp with time zone  default now(),

     unique (ruc),
     unique (name)
);

create table if not exists role (
     id char(2) PRIMARY KEY,
     name varchar(255) not null,
     created_at timestamp with time zone default now()
);

create table if not exists "user" (
     id uuid PRIMARY KEY default gen_random_uuid(),
     email text not null, 
     name text not null, 
     password text not null,
     created_at timestamp with time zone  default now(),

     company_id uuid not null references company(id) on delete restrict,
     role_id char(2) references role(id) on delete restrict,

     unique (email)
);

create table if not exists logged_in (
     user_id uuid PRIMARY KEY references "user"(id) on delete restrict,
     token text not null,
     created_at timestamp with time zone default now()
);

create table if not exists project (
     id uuid PRIMARY KEY default gen_random_uuid(),
     name text not null,
     is_active boolean  default true,

     company_id uuid not null references company(id) on delete restrict,
     created_at timestamp with time zone default now(),

     unique (name, company_id)
);

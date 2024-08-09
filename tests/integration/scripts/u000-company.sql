insert into role (id, name) values
('a', 'admin');

insert into company (id, ruc, name, employees, is_active) values
('3308a6e7-4060-4d7c-8490-f1ccddd9c411', '12345678', 'Company 1', 10, true);

insert into "user" (id, name, email, password, company_id, role_id) values
('1d0cb66e-131c-4001-9e84-8470c1de640a', 'Test User', 'test@test.com', '$2a$08$H0HYFW4KjqGpMe6YEES4p.9UMb1IHQ1WZtTdgD.65CSUbMu/Gciru','3308a6e7-4060-4d7c-8490-f1ccddd9c411', 'a');

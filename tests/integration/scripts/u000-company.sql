insert into role (id, name) values
('a', 'admin');

insert into company (id, ruc, name, employees, is_active) values
('3308a6e7-4060-4d7c-8490-f1ccddd9c411', '12345678', 'Company 1', 10, true),
('b13462ae-e693-481f-8449-f71bd0e1aa84', '87654321', 'Company 2', 20, true);

insert into "user" (id, name, email, password, company_id, role_id) values
('1d0cb66e-131c-4001-9e84-8470c1de640a', 'Test User', 'test@test.com', '$2a$08$H0HYFW4KjqGpMe6YEES4p.9UMb1IHQ1WZtTdgD.65CSUbMu/Gciru','3308a6e7-4060-4d7c-8490-f1ccddd9c411', 'a'),
('5ca76066-8722-4cd1-851e-000f86c48f44', 'Test User', 'test2@test.com', '$2a$08$H0HYFW4KjqGpMe6YEES4p.9UMb1IHQ1WZtTdgD.65CSUbMu/Gciru','b13462ae-e693-481f-8449-f71bd0e1aa84', 'a');

insert into budget_item(id, code, name, level, accumulate, parent_id, company_id) values
('439082ad-f1bd-4228-91f2-8e744894ffdc', '500', 'Costo Directo', 1, true, null, '3308a6e7-4060-4d7c-8490-f1ccddd9c411'),
('b4b2e4e4-f22d-402e-9ab5-1d59347cbfcb', '500.1', 'Obra Gruesa', 2, false, '439082ad-f1bd-4228-91f2-8e744894ffdc', '3308a6e7-4060-4d7c-8490-f1ccddd9c411'),
('420f8bb3-bc8e-4564-be99-75cd7c1a6ff8', '200', 'Costo Urbanizacion', 1, true, null, '3308a6e7-4060-4d7c-8490-f1ccddd9c411'),
('9abc2426-a92b-46ef-b074-ddbc8ee2df1a', '200.1', 'Adoquin', 2, false, '420f8bb3-bc8e-4564-be99-75cd7c1a6ff8', '3308a6e7-4060-4d7c-8490-f1ccddd9c411'),
('5282196b-b054-497c-9f66-b32ff2d7c053', '500', 'Costo Directo', 1, true, null, 'b13462ae-e693-481f-8449-f71bd0e1aa84'),
('8940a68c-f59c-46d8-b61d-82a7eb6cdb5f', '500.1', 'Obra Gruesa', 2, false, '439082ad-f1bd-4228-91f2-8e744894ffdc', 'b13462ae-e693-481f-8449-f71bd0e1aa84');

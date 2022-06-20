INSERT INTO users (user_id, name, email, roles, password_hash, date_created, date_updated) VALUES
	('5cf37266-3473-4006-984f-9325122678b7', 'Admin Gopher', 'admin@example.com', '{ADMIN,USER}', '$2a$10$1ggfMVZV6Js0ybvJufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2/a', '2019-03-24 00:00:00', '2019-03-24 00:00:00'),
	('45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'User Gopher', 'user@example.com', '{USER}', '$2a$10$9/XASPKBbJKVfCAZKDH.UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW', '2019-03-24 00:00:00', '2019-03-24 00:00:00')
	ON CONFLICT DO NOTHING;

INSERT INTO incomes (income_id, user_id, name, category, currency, amount, reoccurrence, duration, reoccurrence_type, duration_type, date_created, date_updated) VALUES
	('a2b0639f-2cc6-44b8-b97b-15d69dbb511e', '45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'Web Development', 'Job', 'EUR', 6250, 1, 0, 'Monthly', 'None', '2019-01-01 00:00:01.000001+00', '2019-01-01 00:00:01.000001+00'),
	('72f8b983-3eb4-48db-9ed0-e45cc6bd716b', '45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'Freelance', 'Part-time Job', 'EUR', 3200, 1, 0, 'Monthly', 'None', '2019-01-01 00:00:02.000001+00', '2019-01-01 00:00:02.000001+00'),
	('7424b1d4-37f2-4d10-ae36-16365ce30dbf', '5cf37266-3473-4006-984f-9325122678b7', 'Software engineer', 'Job', 'EUR', 8200, 1, 0, 'Monthly', 'None', '2019-01-01 00:00:02.000001+00', '2019-01-01 00:00:02.000001+00'),
	('775e20aa-3e57-4b10-b5f3-bad95ee50337', '45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'Grocery store', 'Assets', 'EUR', 500, 1, 0, 'Monthly', 'None', '2019-01-01 00:00:02.000001+00', '2019-01-01 00:00:02.000001+00'),
	('484e5fdb-da74-478b-a66c-f0e28388f257', '45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'Billboard', 'Assets', 'EUR', 250, 1, 0, 'Monthly', 'None', '2019-01-01 00:00:02.000001+00', '2019-01-01 00:00:02.000001+00'),
	('c74111d6-4a5a-41d1-801f-0b8dbc6d3ef9', '45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'Mobile app', 'Part-time Job', 'EUR', 300, 1, 0, 'Monthly', 'None', '2019-01-01 00:00:02.000001+00', '2019-01-01 00:00:02.000001+00'),
	('7d32570f-5206-4ed5-93b0-40f71ecc3300', '45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'Go courses', 'Assets', 'EUR', 1500, 1, 0, 'Monthly', 'None', '2019-01-01 00:00:02.000001+00', '2019-01-01 00:00:02.000001+00'),
	('42483f2e-58d6-4b28-922c-d09c9d1e1193', '45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'Appertment rental', 'Assets', 'EUR', 700, 1, 0, 'Monthly', 'None', '2019-01-01 00:00:02.000001+00', '2019-01-01 00:00:02.000001+00')
	ON CONFLICT DO NOTHING;

INSERT INTO expenses (expense_id, user_id, name, category, currency, amount, reoccurrence, duration, reoccurrence_type, duration_type, date_created, date_updated) VALUES
	('98b6d4b8-f04b-4c79-8c2e-a0aef46854b7', '45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'Taxes', 'Fees', 'EUR', 2550, 1, 0, 'Monthly', 'None', '2019-01-01 00:00:01.000001+00', '2019-01-01 00:00:01.000001+00'),
	('85f6fb09-eb05-4874-ae39-82d1a30fe0d7', '45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'Consumer Basket', 'Other', 'EUR', 800, 1, 0, 'Monthly', 'None', '2019-01-01 00:00:02.000001+00', '2019-01-01 00:00:02.000001+00'),
	('a235be9e-ab5d-44e6-a987-fa1c749264c7', '45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'Appartment Rent', 'Routine', 'EUR', 800, 1, 0, 'Monthly', 'None', '2019-01-01 00:00:02.000001+00', '2019-01-01 00:00:02.000001+00')
	ON CONFLICT DO NOTHING;
-- Fixture data for E2E tests (seed=42)
-- Two countries: Afghanistan and Albania (first two alphabetically)
-- Time range: 2026-05-11 20:00:00 UTC to 2026-05-12 00:00:00 UTC
-- 24 rows per country at 10-minute intervals
--
-- FIXTURE_FROM_ISO = '2026-05-11T20:00:00.000Z'
-- FIXTURE_TO_ISO   = '2026-05-12T00:00:00.000Z'

IF NOT EXISTS (SELECT name FROM sys.databases WHERE name = 'grafanadb')
    CREATE DATABASE grafanadb;
GO

USE grafanadb;
GO

IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'WORLD_DATA' AND schema_id = SCHEMA_ID('dbo'))
    CREATE TABLE dbo.WORLD_DATA (
        BASE_COUNTRY    NVARCHAR(MAX),
        BIRTH_RATE      INT,
        CO2             INT,
        GDP             INT,
        DATE_TIME       DATETIME2,
        TIMESTAMP_VALUE BIGINT
    );
GO

DELETE FROM dbo.WORLD_DATA
WHERE DATE_TIME >= '2026-05-11 20:00:00'
  AND DATE_TIME <  '2026-05-12 00:00:00';
GO

-- Afghanistan rows (24 rows, 10-min intervals starting 2026-05-11 20:00:00 UTC)
INSERT INTO dbo.WORLD_DATA (BASE_COUNTRY, BIRTH_RATE, CO2, GDP, DATE_TIME, TIMESTAMP_VALUE) VALUES
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 30, 185, 470000,  '2026-05-11 20:00:00', 1778529600000),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 32, 190, 480000,  '2026-05-11 20:10:00', 1778530200000),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 31, 195, 490000,  '2026-05-11 20:20:00', 1778530800000),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 33, 200, 500000,  '2026-05-11 20:30:00', 1778531400000),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 30, 195, 510000,  '2026-05-11 20:40:00', 1778532000000),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 34, 210, 520000,  '2026-05-11 20:50:00', 1778532600000),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 31, 205, 515000,  '2026-05-11 21:00:00', 1778533200000),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 32, 198, 505000,  '2026-05-11 21:10:00', 1778533800000),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 33, 192, 495000,  '2026-05-11 21:20:00', 1778534400000),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 30, 187, 485000,  '2026-05-11 21:30:00', 1778535000000),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 31, 193, 475000,  '2026-05-11 21:40:00', 1778535600000),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 32, 199, 488000,  '2026-05-11 21:50:00', 1778536200000),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 34, 207, 502000,  '2026-05-11 22:00:00', 1778536800000),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 30, 212, 515000,  '2026-05-11 22:10:00', 1778537400000),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 31, 208, 525000,  '2026-05-11 22:20:00', 1778538000000),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 33, 203, 512000,  '2026-05-11 22:30:00', 1778538600000),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 32, 197, 498000,  '2026-05-11 22:40:00', 1778539200000),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 30, 191, 487000,  '2026-05-11 22:50:00', 1778539800000),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 31, 185, 472000,  '2026-05-11 23:00:00', 1778540400000),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 34, 196, 466000,  '2026-05-11 23:10:00', 1778541000000),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 32, 202, 478000,  '2026-05-11 23:20:00', 1778541600000),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 30, 209, 493000,  '2026-05-11 23:30:00', 1778542200000),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 33, 204, 507000,  '2026-05-11 23:40:00', 1778542800000),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 31, 198, 521000,  '2026-05-11 23:50:00', 1778543400000);
GO

-- Albania rows (24 rows, 10-min intervals starting 2026-05-11 20:00:00 UTC)
INSERT INTO dbo.WORLD_DATA (BASE_COUNTRY, BIRTH_RATE, CO2, GDP, DATE_TIME, TIMESTAMP_VALUE) VALUES
  ('{"Name":"Albania","Capital":"Tirana"}', 11, 145, 1750000, '2026-05-11 20:00:00', 1778529600000),
  ('{"Name":"Albania","Capital":"Tirana"}', 12, 150, 1780000, '2026-05-11 20:10:00', 1778530200000),
  ('{"Name":"Albania","Capital":"Tirana"}', 10, 155, 1810000, '2026-05-11 20:20:00', 1778530800000),
  ('{"Name":"Albania","Capital":"Tirana"}', 13, 160, 1840000, '2026-05-11 20:30:00', 1778531400000),
  ('{"Name":"Albania","Capital":"Tirana"}', 11, 155, 1860000, '2026-05-11 20:40:00', 1778532000000),
  ('{"Name":"Albania","Capital":"Tirana"}', 12, 148, 1830000, '2026-05-11 20:50:00', 1778532600000),
  ('{"Name":"Albania","Capital":"Tirana"}', 10, 143, 1800000, '2026-05-11 21:00:00', 1778533200000),
  ('{"Name":"Albania","Capital":"Tirana"}', 11, 152, 1770000, '2026-05-11 21:10:00', 1778533800000),
  ('{"Name":"Albania","Capital":"Tirana"}', 13, 158, 1745000, '2026-05-11 21:20:00', 1778534400000),
  ('{"Name":"Albania","Capital":"Tirana"}', 12, 163, 1760000, '2026-05-11 21:30:00', 1778535000000),
  ('{"Name":"Albania","Capital":"Tirana"}', 10, 157, 1785000, '2026-05-11 21:40:00', 1778535600000),
  ('{"Name":"Albania","Capital":"Tirana"}', 11, 151, 1815000, '2026-05-11 21:50:00', 1778536200000),
  ('{"Name":"Albania","Capital":"Tirana"}', 12, 146, 1845000, '2026-05-11 22:00:00', 1778536800000),
  ('{"Name":"Albania","Capital":"Tirana"}', 13, 153, 1870000, '2026-05-11 22:10:00', 1778537400000),
  ('{"Name":"Albania","Capital":"Tirana"}', 10, 159, 1855000, '2026-05-11 22:20:00', 1778538000000),
  ('{"Name":"Albania","Capital":"Tirana"}', 11, 164, 1835000, '2026-05-11 22:30:00', 1778538600000),
  ('{"Name":"Albania","Capital":"Tirana"}', 12, 156, 1810000, '2026-05-11 22:40:00', 1778539200000),
  ('{"Name":"Albania","Capital":"Tirana"}', 10, 149, 1792000, '2026-05-11 22:50:00', 1778539800000),
  ('{"Name":"Albania","Capital":"Tirana"}', 13, 144, 1775000, '2026-05-11 23:00:00', 1778540400000),
  ('{"Name":"Albania","Capital":"Tirana"}', 11, 150, 1758000, '2026-05-11 23:10:00', 1778541000000),
  ('{"Name":"Albania","Capital":"Tirana"}', 12, 157, 1772000, '2026-05-11 23:20:00', 1778541600000),
  ('{"Name":"Albania","Capital":"Tirana"}', 10, 163, 1790000, '2026-05-11 23:30:00', 1778542200000),
  ('{"Name":"Albania","Capital":"Tirana"}', 13, 155, 1808000, '2026-05-11 23:40:00', 1778542800000),
  ('{"Name":"Albania","Capital":"Tirana"}', 11, 148, 1832000, '2026-05-11 23:50:00', 1778543400000);
GO

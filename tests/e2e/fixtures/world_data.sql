-- Fixture data for E2E tests (seed=42)
-- Two countries: Afghanistan and Albania (first two alphabetically)
-- Dynamic time range: 24 rows per country at 10-minute intervals
-- Oldest row: GETUTCDATE() - 4h; newest row: GETUTCDATE() - 10min

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

-- Afghanistan rows (24 rows, 10-min intervals, now-240min to now-10min)
INSERT INTO dbo.WORLD_DATA (BASE_COUNTRY, BIRTH_RATE, CO2, GDP, DATE_TIME, TIMESTAMP_VALUE) VALUES
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 30, 185, 470000,  DATEADD(minute, -240, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute, -240, GETUTCDATE()))),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 32, 190, 480000,  DATEADD(minute, -230, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute, -230, GETUTCDATE()))),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 31, 195, 490000,  DATEADD(minute, -220, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute, -220, GETUTCDATE()))),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 33, 200, 500000,  DATEADD(minute, -210, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute, -210, GETUTCDATE()))),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 30, 195, 510000,  DATEADD(minute, -200, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute, -200, GETUTCDATE()))),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 34, 210, 520000,  DATEADD(minute, -190, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute, -190, GETUTCDATE()))),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 31, 205, 515000,  DATEADD(minute, -180, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute, -180, GETUTCDATE()))),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 32, 198, 505000,  DATEADD(minute, -170, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute, -170, GETUTCDATE()))),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 33, 192, 495000,  DATEADD(minute, -160, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute, -160, GETUTCDATE()))),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 30, 187, 485000,  DATEADD(minute, -150, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute, -150, GETUTCDATE()))),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 31, 193, 475000,  DATEADD(minute, -140, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute, -140, GETUTCDATE()))),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 32, 199, 488000,  DATEADD(minute, -130, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute, -130, GETUTCDATE()))),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 34, 207, 502000,  DATEADD(minute, -120, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute, -120, GETUTCDATE()))),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 30, 212, 515000,  DATEADD(minute, -110, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute, -110, GETUTCDATE()))),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 31, 208, 525000,  DATEADD(minute, -100, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute, -100, GETUTCDATE()))),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 33, 203, 512000,  DATEADD(minute,  -90, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute,  -90, GETUTCDATE()))),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 32, 197, 498000,  DATEADD(minute,  -80, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute,  -80, GETUTCDATE()))),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 30, 191, 487000,  DATEADD(minute,  -70, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute,  -70, GETUTCDATE()))),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 31, 185, 472000,  DATEADD(minute,  -60, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute,  -60, GETUTCDATE()))),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 34, 196, 466000,  DATEADD(minute,  -50, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute,  -50, GETUTCDATE()))),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 32, 202, 478000,  DATEADD(minute,  -40, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute,  -40, GETUTCDATE()))),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 30, 209, 493000,  DATEADD(minute,  -30, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute,  -30, GETUTCDATE()))),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 33, 204, 507000,  DATEADD(minute,  -20, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute,  -20, GETUTCDATE()))),
  ('{"Name":"Afghanistan","Capital":"Kabul"}', 31, 198, 521000,  DATEADD(minute,  -10, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute,  -10, GETUTCDATE())));
GO

-- Albania rows (24 rows, 10-min intervals, now-240min to now-10min)
INSERT INTO dbo.WORLD_DATA (BASE_COUNTRY, BIRTH_RATE, CO2, GDP, DATE_TIME, TIMESTAMP_VALUE) VALUES
  ('{"Name":"Albania","Capital":"Tirana"}', 11, 145, 1750000, DATEADD(minute, -240, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute, -240, GETUTCDATE()))),
  ('{"Name":"Albania","Capital":"Tirana"}', 12, 150, 1780000, DATEADD(minute, -230, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute, -230, GETUTCDATE()))),
  ('{"Name":"Albania","Capital":"Tirana"}', 10, 155, 1810000, DATEADD(minute, -220, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute, -220, GETUTCDATE()))),
  ('{"Name":"Albania","Capital":"Tirana"}', 13, 160, 1840000, DATEADD(minute, -210, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute, -210, GETUTCDATE()))),
  ('{"Name":"Albania","Capital":"Tirana"}', 11, 155, 1860000, DATEADD(minute, -200, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute, -200, GETUTCDATE()))),
  ('{"Name":"Albania","Capital":"Tirana"}', 12, 148, 1830000, DATEADD(minute, -190, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute, -190, GETUTCDATE()))),
  ('{"Name":"Albania","Capital":"Tirana"}', 10, 143, 1800000, DATEADD(minute, -180, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute, -180, GETUTCDATE()))),
  ('{"Name":"Albania","Capital":"Tirana"}', 11, 152, 1770000, DATEADD(minute, -170, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute, -170, GETUTCDATE()))),
  ('{"Name":"Albania","Capital":"Tirana"}', 13, 158, 1745000, DATEADD(minute, -160, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute, -160, GETUTCDATE()))),
  ('{"Name":"Albania","Capital":"Tirana"}', 12, 163, 1760000, DATEADD(minute, -150, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute, -150, GETUTCDATE()))),
  ('{"Name":"Albania","Capital":"Tirana"}', 10, 157, 1785000, DATEADD(minute, -140, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute, -140, GETUTCDATE()))),
  ('{"Name":"Albania","Capital":"Tirana"}', 11, 151, 1815000, DATEADD(minute, -130, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute, -130, GETUTCDATE()))),
  ('{"Name":"Albania","Capital":"Tirana"}', 12, 146, 1845000, DATEADD(minute, -120, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute, -120, GETUTCDATE()))),
  ('{"Name":"Albania","Capital":"Tirana"}', 13, 153, 1870000, DATEADD(minute, -110, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute, -110, GETUTCDATE()))),
  ('{"Name":"Albania","Capital":"Tirana"}', 10, 159, 1855000, DATEADD(minute, -100, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute, -100, GETUTCDATE()))),
  ('{"Name":"Albania","Capital":"Tirana"}', 11, 164, 1835000, DATEADD(minute,  -90, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute,  -90, GETUTCDATE()))),
  ('{"Name":"Albania","Capital":"Tirana"}', 12, 156, 1810000, DATEADD(minute,  -80, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute,  -80, GETUTCDATE()))),
  ('{"Name":"Albania","Capital":"Tirana"}', 10, 149, 1792000, DATEADD(minute,  -70, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute,  -70, GETUTCDATE()))),
  ('{"Name":"Albania","Capital":"Tirana"}', 13, 144, 1775000, DATEADD(minute,  -60, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute,  -60, GETUTCDATE()))),
  ('{"Name":"Albania","Capital":"Tirana"}', 11, 150, 1758000, DATEADD(minute,  -50, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute,  -50, GETUTCDATE()))),
  ('{"Name":"Albania","Capital":"Tirana"}', 12, 157, 1772000, DATEADD(minute,  -40, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute,  -40, GETUTCDATE()))),
  ('{"Name":"Albania","Capital":"Tirana"}', 10, 163, 1790000, DATEADD(minute,  -30, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute,  -30, GETUTCDATE()))),
  ('{"Name":"Albania","Capital":"Tirana"}', 13, 155, 1808000, DATEADD(minute,  -20, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute,  -20, GETUTCDATE()))),
  ('{"Name":"Albania","Capital":"Tirana"}', 11, 148, 1832000, DATEADD(minute,  -10, GETUTCDATE()), DATEDIFF_BIG(millisecond, '1970-01-01', DATEADD(minute,  -10, GETUTCDATE())));
GO

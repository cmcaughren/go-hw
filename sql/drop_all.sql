USE HWDB;
GO

-- Drop child tables first
IF OBJECT_ID('ReturnItem', 'U') IS NOT NULL DROP TABLE ReturnItem;
GO

IF OBJECT_ID('LineItem', 'U') IS NOT NULL DROP TABLE LineItem;
GO

IF OBJECT_ID('[Order]', 'U') IS NOT NULL DROP TABLE [Order];
GO

IF OBJECT_ID('CustomerAddress', 'U') IS NOT NULL DROP TABLE CustomerAddress;
GO

-- Drop parent tables last
IF OBJECT_ID('Customer', 'U') IS NOT NULL DROP TABLE Customer;
GO

IF OBJECT_ID('Product', 'U') IS NOT NULL DROP TABLE Product;
GO

PRINT 'All tables dropped';
GO
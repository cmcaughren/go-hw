USE HWDB;
GO

CREATE TABLE Customer (
   CustomerID INT IDENTITY(1,1) PRIMARY KEY,
   FirstName NVARCHAR(50),
   LastName NVARCHAR(50),
   DOB DATETIME,
   OrderTotal DECIMAL(8,2),
   OrderTaxTotal DECIMAL(8,2)
);
GO

CREATE TABLE CustomerAddress (
   CustomerAddressID INT IDENTITY(1,1) PRIMARY KEY,
   CustomerID INT FOREIGN KEY REFERENCES Customer(CustomerID),
   Line1 NVARCHAR(100),
   Line2 NVARCHAR(100),
   City NVARCHAR(50),
   StateProvince NVARCHAR(50),
   Country NVARCHAR(50)
);
GO

-- Not listed in the table set but required for LineItem 
CREATE TABLE Product (
   ProductID INT IDENTITY(1,1) PRIMARY KEY,
   ProductName NVARCHAR(50)
);
GO

CREATE TABLE [Order] (
   OrderID INT IDENTITY(1,1) PRIMARY KEY,
   OrderNumber NVARCHAR(50) NOT NULL UNIQUE,
   CustomerID INT FOREIGN KEY REFERENCES Customer(CustomerID),
   OrderCreateDate DATETIME DEFAULT GETDATE(),
   OrderFulfilledDate DATETIME,
   OrderTotal DECIMAL(8,2),
   OrderTaxTotal DECIMAL(8,2)
);
GO

CREATE TABLE LineItem (
   LineItemID INT IDENTITY(1,1) PRIMARY KEY,
   OrderID INT FOREIGN KEY REFERENCES [Order](OrderID),
   ProductID INT FOREIGN KEY REFERENCES Product(ProductID),
   LineItemUnitPrice DECIMAL(8,2),
   Quantity INT,
   LineItemDiscount DECIMAL(8,2)
);
GO

CREATE TABLE ReturnItem (
   ReturnedItemID INT IDENTITY(1,1) PRIMARY KEY,
   OrderID INT FOREIGN KEY REFERENCES [Order](OrderID),
   LineItemID INT FOREIGN KEY REFERENCES LineItem(LineItemID),
   Quantity INT,
   AmountReturned DECIMAL(8,2)
);
GO
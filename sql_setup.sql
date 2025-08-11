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

USE HWDB;
GO


-- CREATE TEST DATA
INSERT INTO Customer (FirstName, LastName, DOB, OrderTotal, OrderTaxTotal) VALUES
   ('John', 'Smith', '1985-03-15', 524.66, 58.78),
   ('Sarah', 'Johnson', '1990-07-22', 655.06, 70.18);
GO

INSERT INTO CustomerAddress (CustomerID, Line1, Line2, City, StateProvince, Country) VALUES
   (1, '123 Main Street', 'Apt 4B', 'Vancouver', 'BC', 'Canada'),
   (2, '456 Oak Avenue', NULL, 'Victoria', 'BC', 'Canada');
GO

INSERT INTO Product (ProductName) VALUES
   ('Wireless Mouse'),
   ('USB-C Cable'),
   ('Bluetooth Headphones'),
   ('Laptop Stand'),
   ('Webcam HD'),
   ('Keyboard Mechanical'),
   ('Notebook Set'),
   ('Pen Pack'),
   ('Desk Organizer'),
   ('Stapler'),
   ('Paper Ream'),
   ('Folder Set'),
   ('Coffee Mug'),
   ('Water Bottle'),
   ('Desk Lamp'),
   ('Plant Pot'),
   ('Wall Clock'),
   ('Picture Frame'),
   ('SQL Server Guide'),
   ('Go Programming'),
   ('Clean Code'),
   ('Design Patterns'),
   ('Algorithms Book'),
   ('Database Design');
GO

-- Brackets around "Order" because it's a reserved word
INSERT INTO [Order] (OrderNumber, CustomerID, OrderCreateDate, OrderFulfilledDate, OrderTotal, OrderTaxTotal) VALUES
   ('ORD-2024-001', 1, '2024-01-15 10:30:00', '2024-01-17 14:00:00', 289.94, 34.79),
   ('ORD-2024-002', 1, '2024-02-20 09:15:00', '2024-02-22 16:30:00', 234.72, 23.99),
   ('ORD-2024-003', 2, '2024-03-10 11:45:00', '2024-03-12 13:20:00', 312.45, 31.19),
   ('ORD-2024-004', 2, '2024-04-05 14:20:00', NULL, 342.61, 38.99);
GO


-- Order 1
INSERT INTO LineItem (OrderID, ProductID, LineItemUnitPrice, Quantity, LineItemDiscount) VALUES
   (1, 1, 29.99, 2, 5.00),
   (1, 3, 89.99, 1, 0.00),
   (1, 7, 15.99, 3, 2.00),
   (1, 13, 11.99, 2, 0.00),
   (1, 19, 49.99, 1, 0.00),
   (1, 20, 44.99, 1, 0.00);

-- Order 2  
INSERT INTO LineItem (OrderID, ProductID, LineItemUnitPrice, Quantity, LineItemDiscount) VALUES
   (2, 2, 12.99, 3, 1.00),
   (2, 8, 8.99, 5, 3.00),
   (2, 14, 19.99, 2, 0.00),
   (2, 15, 34.99, 1, 0.00),
   (2, 21, 39.99, 1, 0.00),
   (2, 11, 9.99, 2, 0.00);

-- Order 3
INSERT INTO LineItem (OrderID, ProductID, LineItemUnitPrice, Quantity, LineItemDiscount) VALUES
   (3, 4, 45.99, 1, 0.00),
   (3, 5, 67.99, 1, 5.00),
   (3, 9, 24.99, 2, 2.00),
   (3, 16, 16.99, 3, 0.00),
   (3, 22, 54.99, 1, 0.00),
   (3, 10, 18.99, 1, 0.00);

-- Order 4
INSERT INTO LineItem (OrderID, ProductID, LineItemUnitPrice, Quantity, LineItemDiscount) VALUES
   (4, 6, 129.99, 1, 10.00),
   (4, 12, 12.99, 4, 2.00),
   (4, 17, 22.99, 2, 0.00),
   (4, 18, 14.99, 3, 1.50),
   (4, 23, 59.99, 1, 0.00),
   (4, 24, 47.99, 1, 0.00);
GO

INSERT INTO ReturnItem (OrderID, LineItemID, Quantity, AmountReturned) VALUES
   (1, 1, 1, 29.99),
   (1, 4, 1, 11.99),
   (2, 8, 2, 17.98),
   (3, 15, 1, 16.99);
GO
-- AI GENERATED TEST DATA! 

USE HWDB;
GO

-- ============================================
-- PRODUCTS (24 total - MUST be inserted first!)
-- ============================================
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

PRINT '24 Products inserted';
GO

-- ============================================
-- CUSTOMERS (20 total, various ages and countries)
-- Start with NULL totals - will update after orders
-- ============================================
INSERT INTO Customer (FirstName, LastName, DOB, OrderTotal, OrderTaxTotal) VALUES
-- Canadian customers
('John', 'Smith', '1985-03-15', NULL, NULL),      -- Age 39
('Sarah', 'Johnson', '1990-07-22', NULL, NULL),   -- Age 34
('Michael', 'Brown', '2007-11-03', NULL, NULL),   -- Age 16 (under 18)
('Emma', 'Davis', '1958-02-28', NULL, NULL),      -- Age 66 (senior)
('Lucas', 'Wilson', '2008-05-15', NULL, NULL),    -- Age 16 (under 18)

-- US customers
('Jennifer', 'Martinez', '1975-09-10', NULL, NULL), -- Age 49
('Robert', 'Anderson', '1955-12-25', NULL, NULL),   -- Age 68 (senior)
('Maria', 'Garcia', '1992-04-17', NULL, NULL),      -- Age 32
('William', 'Taylor', '2006-08-30', NULL, NULL),    -- Age 18

-- UK customers  
('Oliver', 'Thomas', '1988-01-20', NULL, NULL),     -- Age 36
('Charlotte', 'Jackson', '1950-06-12', NULL, NULL), -- Age 74 (senior)
('James', 'White', '1995-03-08', NULL, NULL),       -- Age 29

-- German customers
('Hans', 'Mueller', '1982-11-15', NULL, NULL),      -- Age 41
('Anna', 'Schmidt', '2009-02-01', NULL, NULL),      -- Age 15 (under 18)

-- Japanese customers
('Yuki', 'Tanaka', '1978-07-07', NULL, NULL),       -- Age 46
('Kenji', 'Yamamoto', '1962-04-22', NULL, NULL),    -- Age 62

-- Australian customers
('Liam', 'O''Connor', '1991-10-30', NULL, NULL),    -- Age 32
('Sophie', 'Chen', '1987-05-18', NULL, NULL),       -- Age 37
('Jack', 'Thompson', '2010-12-05', NULL, NULL),     -- Age 14 (under 18)
('Grace', 'Williams', '1948-09-14', NULL, NULL);    -- Age 76 (senior)
GO

PRINT '20 Customers inserted';
GO

-- ============================================
-- CUSTOMER ADDRESSES
-- ============================================
INSERT INTO CustomerAddress (CustomerID, Line1, Line2, City, StateProvince, Country) VALUES
-- Canadian addresses
(1, '123 Main Street', 'Apt 4B', 'Vancouver', 'BC', 'Canada'),
(2, '456 Oak Avenue', NULL, 'Toronto', 'ON', 'Canada'),
(3, '789 Maple Drive', NULL, 'Calgary', 'AB', 'Canada'),
(4, '321 Pine Street', 'Unit 12', 'Montreal', 'QC', 'Canada'),
(5, '654 Elm Road', NULL, 'Ottawa', 'ON', 'Canada'),

-- US addresses
(6, '100 Broadway', 'Suite 500', 'New York', 'NY', 'United States'),
(7, '200 Sunset Blvd', NULL, 'Los Angeles', 'CA', 'United States'),
(8, '300 Main Plaza', 'Apt 7', 'Chicago', 'IL', 'United States'),
(9, '400 Tech Drive', NULL, 'Austin', 'TX', 'United States'),

-- UK addresses
(10, '10 High Street', 'Flat 3', 'London', 'England', 'United Kingdom'),
(11, '20 Castle Road', NULL, 'Edinburgh', 'Scotland', 'United Kingdom'),
(12, '30 Queen Street', NULL, 'Manchester', 'England', 'United Kingdom'),

-- German addresses
(13, 'Hauptstrasse 15', NULL, 'Berlin', 'Berlin', 'Germany'),
(14, 'Bahnhofstrasse 42', 'Wohnung 5', 'Munich', 'Bavaria', 'Germany'),

-- Japanese addresses
(15, '1-2-3 Shibuya', 'Apt 801', 'Tokyo', 'Tokyo', 'Japan'),
(16, '4-5-6 Namba', NULL, 'Osaka', 'Osaka', 'Japan'),

-- Australian addresses
(17, '50 George Street', NULL, 'Sydney', 'NSW', 'Australia'),
(18, '75 Collins Street', 'Level 10', 'Melbourne', 'VIC', 'Australia'),
(19, '25 Beach Road', NULL, 'Brisbane', 'QLD', 'Australia'),
(20, '88 King William', 'Unit 3A', 'Adelaide', 'SA', 'Australia');
GO

PRINT '20 Customer Addresses inserted';
GO

-- ============================================
-- ORDERS (200 total with varying fulfillment times)
-- ============================================
DECLARE @OrderID INT = 1;
DECLARE @OrderDate DATETIME;
DECLARE @FulfillDate DATETIME;
DECLARE @CustomerID INT;
DECLARE @DaysToFulfill INT;

WHILE @OrderID <= 200
BEGIN
    -- Randomly assign to customers (weighted so some customers have more orders)
    SET @CustomerID = CASE 
        WHEN @OrderID % 7 = 0 THEN 1  -- Customer 1 gets more orders
        WHEN @OrderID % 11 = 0 THEN 2 -- Customer 2 gets more orders
        WHEN @OrderID % 13 = 0 THEN 10 -- Customer 10 gets more orders
        ELSE 1 + (@OrderID % 20)      -- Others distributed
    END;
    
    -- Generate order date (spread over last 6 months)
    SET @OrderDate = DATEADD(DAY, -180 + (@OrderID * 180 / 200), GETDATE());
    
    -- Vary fulfillment time
    SET @DaysToFulfill = CASE 
        WHEN @OrderID % 10 = 0 THEN 0    -- 10% same day
        WHEN @OrderID % 5 = 0 THEN 1     -- 20% next day
        WHEN @OrderID % 3 = 0 THEN 2     -- 33% 2 days
        WHEN @OrderID % 7 = 0 THEN 5     -- ~14% 5 days
        WHEN @OrderID % 11 = 0 THEN 7    -- ~9% week
        WHEN @OrderID > 190 THEN NULL    -- Last 10 orders not fulfilled yet
        ELSE 3                           -- Rest 3 days
    END;
    
    -- Calculate fulfillment date
    SET @FulfillDate = CASE 
        WHEN @DaysToFulfill IS NULL THEN NULL
        ELSE DATEADD(DAY, @DaysToFulfill, @OrderDate)
    END;
    
    -- Insert order with NULL totals (will calculate after line items)
    INSERT INTO [Order] (OrderNumber, CustomerID, OrderCreateDate, OrderFulfilledDate, OrderTotal, OrderTaxTotal)
    VALUES (
        'ORD-2024-' + RIGHT('000' + CAST(@OrderID AS VARCHAR(3)), 3),
        @CustomerID,
        @OrderDate,
        @FulfillDate,
        NULL,  -- Will update after line items
        NULL   -- Will update after line items
    );
    
    SET @OrderID = @OrderID + 1;
END;
GO

PRINT '200 Orders inserted';
GO

-- ============================================
-- LINE ITEMS (6-8 items per order, ~1400 total)
-- ============================================
DECLARE @LineOrderID INT = 1;
DECLARE @ItemCount INT;
DECLARE @ItemNum INT;
DECLARE @ProdID INT;
DECLARE @Price DECIMAL(8,2);

WHILE @LineOrderID <= 200
BEGIN
    -- Vary items per order (6-8 items)
    SET @ItemCount = 6 + (@LineOrderID % 3);
    SET @ItemNum = 1;
    
    WHILE @ItemNum <= @ItemCount
    BEGIN
        -- Select different products (cycling through 24 products)
        SET @ProdID = 1 + ((@LineOrderID + @ItemNum - 2) % 24);
        
        -- Set price based on product type
        SET @Price = CASE 
            WHEN @ProdID BETWEEN 1 AND 6 THEN 29.99 + (@ProdID * 10)    -- Electronics: $40-90
            WHEN @ProdID BETWEEN 7 AND 12 THEN 8.99 + (@ProdID * 2)     -- Office: $10-30
            WHEN @ProdID BETWEEN 13 AND 18 THEN 11.99 + (@ProdID * 1.5) -- Home: $15-40
            ELSE 39.99 + (@ProdID * 2)                                   -- Books: $40-80
        END;
        
        INSERT INTO LineItem (OrderID, ProductID, LineItemUnitPrice, Quantity, LineItemDiscount)
        VALUES (
            @LineOrderID,
            @ProdID,
            @Price,
            1 + (@ItemNum % 3),     -- Quantity 1-3
            CASE WHEN @ItemNum % 5 = 0 THEN (@Price * 0.10) ELSE 0 END  -- 10% discount on some items
        );
        
        SET @ItemNum = @ItemNum + 1;
    END;
    
    SET @LineOrderID = @LineOrderID + 1;
END;
GO

PRINT '~1400 Line Items inserted';
GO

-- ============================================
-- UPDATE ORDER TOTALS based on actual line items
-- ============================================
UPDATE o
SET OrderTotal = ISNULL(items.SubTotal, 0),
    OrderTaxTotal = ISNULL(items.SubTotal * 0.10, 0)  -- 10% tax rate
FROM [Order] o
LEFT JOIN (
    SELECT 
        OrderID,
        SUM(LineItemUnitPrice * Quantity - ISNULL(LineItemDiscount, 0)) as SubTotal
    FROM LineItem
    GROUP BY OrderID
) items ON o.OrderID = items.OrderID;
GO

PRINT 'Order totals updated based on actual line items';
GO

-- ============================================
-- UPDATE CUSTOMER TOTALS based on actual orders
-- ============================================
UPDATE c
SET OrderTotal = ISNULL(totals.TotalAmount, 0),
    OrderTaxTotal = ISNULL(totals.TotalTax, 0)
FROM Customer c
LEFT JOIN (
    SELECT 
        CustomerID,
        SUM(OrderTotal) as TotalAmount,
        SUM(OrderTaxTotal) as TotalTax
    FROM [Order]
    GROUP BY CustomerID
) totals ON c.CustomerID = totals.CustomerID;
GO

PRINT 'Customer totals updated based on actual orders';
GO

-- ============================================
-- RETURNS (~15% of items, ~200 returns)
-- ============================================
DECLARE @ReturnCount INT = 1;
DECLARE @ReturnOrderID INT;
DECLARE @ReturnLineItemID INT;
DECLARE @ReturnAmount DECIMAL(8,2);
DECLARE @ReturnQty INT;
DECLARE @UnitPrice DECIMAL(8,2);
DECLARE @OriginalQty INT;
DECLARE @LineDiscount DECIMAL(8,2);

WHILE @ReturnCount <= 200
BEGIN
    -- Select orders to have returns (spread across orders)
    SET @ReturnOrderID = 1 + ((@ReturnCount * 7) % 180);  -- Not from unfulfilled orders
    
    -- Calculate line item ID (approximate, assuming sequential)
    SET @ReturnLineItemID = (@ReturnOrderID - 1) * 7 + (@ReturnCount % 7) + 1;
    
    -- Skip if this LineItemID doesn't exist or already has returns
    IF NOT EXISTS (SELECT 1 FROM LineItem WHERE LineItemID = @ReturnLineItemID)
       OR EXISTS (SELECT 1 FROM ReturnItem WHERE LineItemID = @ReturnLineItemID)
    BEGIN
        SET @ReturnCount = @ReturnCount + 1;
        CONTINUE;
    END
    
    -- Get the line item details
    SELECT 
        @UnitPrice = LineItemUnitPrice,
        @OriginalQty = Quantity,
        @LineDiscount = ISNULL(LineItemDiscount, 0)
    FROM LineItem 
    WHERE LineItemID = @ReturnLineItemID;
    
    -- Determine how many units to return (1 to original quantity)
    SET @ReturnQty = CASE 
        WHEN @OriginalQty = 1 THEN 1
        WHEN @ReturnCount % 3 = 0 THEN @OriginalQty  -- Sometimes return all
        ELSE 1  -- Usually return just 1
    END;
    
    -- Calculate return amount: (unit price * returned qty) - proportional discount
    SET @ReturnAmount = (@UnitPrice * @ReturnQty) - 
                        (@LineDiscount * @ReturnQty / @OriginalQty);
    
    -- Insert return
    INSERT INTO ReturnItem (OrderID, LineItemID, Quantity, AmountReturned)
    VALUES (
        @ReturnOrderID,
        @ReturnLineItemID,
        @ReturnQty,
        ROUND(@ReturnAmount, 2)
    );
    
    SET @ReturnCount = @ReturnCount + 1;
END;
GO

PRINT '~200 Returns inserted';
GO

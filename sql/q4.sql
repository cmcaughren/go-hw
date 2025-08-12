USE HWDB;
GO

-- What is the best-selling product for a given time period, based on:
--    total sales (i.e. it has generated the most revenue),
--    and unit sales (i.e. it has sold the most units)?

-- ** Should these account for returns? 
SELECT TOP 1 ProductName, SUM(Quantity * LineItemUnitPrice) as total_sales
FROM Product 
INNER JOIN LineItem
   ON Product.ProductID = LineItem.ProductID
GROUP BY ProductName
ORDER BY total_sales DESC
GO

SELECT TOP 1 ProductName, SUM(Quantity) as unit_sales
FROM Product 
INNER JOIN LineItem
   ON Product.ProductID = LineItem.ProductID
GROUP BY ProductName
ORDER BY unit_sales DESC
GO

-- Who are the top 5 customers based on their total net sales 
-- (i.e. net sales would be the total amount of sales they have minus 
-- any returns they have made)

-- **Have I used the correct "join"?
SELECT TOP 5 FullName, SUM(OrderTotal - AmountReturned) as net_sales
FROM (
   (SELECT CONCAT(FirstName, " ", LastName) as FullName, [Order].OrderTotal, [Order].OrderID
   FROM Customer
   INNER JOIN [Order]
      ON Customer.CustomerID = [Order].CustomerID) AS Table_A
      INNER JOIN
   (SELECT AmountReturned, [Order].OrderID
   FROM ReturnItem
   INNER JOIN [Order]
      ON ReturnItem.OrderID = [Order].OrderID) AS Table_B
   ON Table_A.OrderID = Table_B.OrderID )
GROUP By FullName
ORDER BY net_sales DESC;
GO


-- If we categorize customers by these age brackets:
--    Under 18 years old,
--    18 - 29 years old,
--    30 - 45 years old,
--    45 - 65 years old,
--    65+ years old
-- How would we get the following:
--    Which age group generated the greatest sales?
--    What was the top coutnry from sales perspective from each age bracket
--    and what was their total sales?


CREATE VIEW Customer_Agegroups AS
SELECT CONCAT(FirstName, " ", LastName) as FullName, 
      (DATEDIFF(year, DOB, GETDATE())) as Age,
      CASE 
         WHEN (DATEDIFF(year, DOB, GETDATE())) < 18 THEN 'Under 18 years old'
         WHEN (DATEDIFF(year, DOB, GETDATE())) BETWEEN 18 AND 29 THEN '18-29 years old'
         WHEN (DATEDIFF(year, DOB, GETDATE())) BETWEEN 30 AND 45 THEN '30-45 years old'
         WHEN (DATEDIFF(year, DOB, GETDATE())) BETWEEN 46 AND 65 THEN '46-65 years old'
         WHEN (DATEDIFF(year, DOB, GETDATE())) > 65 THEN '65+ years old'
      END AS Age_Category,
      CustomerID,
      OrderTotal
FROM Customer;
GO

SELECT TOP 1 Age_Category, SUM(OrderTotal) as total_sales
FROM Customer_Agegroups
GROUP BY Age_Category
ORDER BY total_sales DESC;
GO

SELECT Age_Category, SUM(OrderTotal) as total_sales, Country
FROM Customer_Agegroups
INNER JOIN CustomerAddress
   ON Customer_Agegroups.CustomerID = CustomerAddress.CustomerID
GROUP BY Age_Category, Country;
GO

SELECT Table_B.Age_Category, Table_B.Country, MAX(Table_B.total_sales) as highest_sales
FROM (SELECT Age_Category, SUM(OrderTotal) as total_sales, Country
      FROM Customer_Agegroups
      INNER JOIN CustomerAddress
      ON Customer_Agegroups.CustomerID = CustomerAddress.CustomerID
      GROUP BY Age_Category, Country) as Table_B
GROUP BY Table_B.Age_Category, Table_B.Country
ORDER BY highest_sales DESC;
GO

-- Please produce the underlying data for a Histogram, that shows the number of orders
-- by the number of days it took to fulfill the order after it was placed. 
-- For example, would be looking for something like the following:
--    Less than a dat - 50 orders
--    1 - 2 days - 350 orders
--    2 - 3 days - 400 orders
--    3 - 4 days - 200 orders
--    Greater than 4 days - 75 orders

-- What percentage of items that were ordered included a discount from the full price?





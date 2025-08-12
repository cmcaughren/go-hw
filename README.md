For Question 1:
- See /crudAPI/q1.go
- It works with the database set up for Questions 4. 

For Question 2 and Question 3:
- See /series/q2q3.go

For Question 4:
- PREREQS: azure sql server database & sqlcmd
- To run SQL queries, create a .env file, enter information for azure sql server database connection.
- .env.example outlines how the .env file should look
- sql_setup file expects the database is called HWDB.
- Set up the tables:
   ./run_sql.sh sql/sql_setup.sql
- Fill the database with test data:
   ./run_sql.sh sql/load_test_data.sql
- If you need a fresh start, drop the tables:
   ./run_sql.sh sql/drop_all.sql 
- Run the Q4 queries:
   ./run_sql.sh sql/q4.sql




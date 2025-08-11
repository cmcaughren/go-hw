For Question 4:
- PREREQS: azure sql server database & sqlcmd
- To run SQL queries, create a .env file, enter information for azure sql server database connection.
- .env.example outlines how the .env file should look
- sql_setup file expects the database is called HWDB.
- Set up the tables:
   ./runsql sql_setup.sql
- Fill the database with test data:
   ./runsql load_test_data.sql
- If you need a fresh start, drop the tables:
   ./runsql drop_all.sql 
- Run the Q4 queries:
   ./runsql q4.sql




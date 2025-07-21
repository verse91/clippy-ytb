## Database Migration

To apply the database schema, run the following command (replace connection string as needed):

```
psql <connection-string> -f backend/migrations/20240721_create_downloads_table.sql
```

Or, if using a migration tool, point it to the `backend/migrations` directory.

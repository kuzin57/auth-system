# auth-system

How to use?

1. Create .env file in root directory
2. Set the following variables:
   - SMTP_HOST
   - SMTP_PORT
   - SMTP_PASS
   - SMTP_FROM
   - REGISTRATION_LINK_BASE

Example:
```
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_PASS=your_password
SMTP_FROM=your_email@gmail.com
REGISTRATION_LINK_BASE=http://localhost:8080
```

REGISTRATION_LINK_BASE is the base URL of the running server.
If you run application on the remote server, you can forward port 8080 to the remote server.

3. Run the following command:
   ```bash
   docker-compose up -d
   ```

4. Open Postman and import the collection `auth-system.postman_collection.json` or send requests via curl.

Example:
```bash
curl -X POST "http://localhost:8080/send-registration-link" -H "Content-Type: application/json" -d '{"email": "user@example.com"}'
```

Example response:
```json
{
  "status": "success"
}
```

5. Open the link in the email to register.

Example link:
```
http://localhost:8080/registration?token=your-jwt-token-here
```

6. Register the user.

7. To see all users, send the following request:
```bash
curl -X GET "http://localhost:8080/users"
```

Example response:
```json
{
  "users": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "email": "user@example.com",
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    }
  ]
}
```

Example usage (the whole flow), video:

https://drive.google.com/file/d/1Hfn3PXqAb0ry0KgLcUmPl3qeMenMrNUW/view?usp=sharing 

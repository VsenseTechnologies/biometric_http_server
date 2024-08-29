# The Go Fingerprint API Documentation
This Documentation Contains Information about all the Routes and Data related to particular Route.
## Main Routing Parameter For Every Route

```bash
# This Two Parameters Are Most Important For Every Single Route
# This Is For Admin Requests
/admin/...
# This Is For User Requests
/users/...
```
## Routes Related To Admins

### 1] Authentication URL Parameters
```bash
# Post Request

# General Authentication Route
http://localhost:8080/{id}/register
# Replace The {id} Either By users Or By admin
http://localhost:8080/admin/register
http://localhost:8080/users/register
# For Login
http://localhost:8080/admin/login
http://localhost:8080/users/login
```
### Authentication JSON Data Format To Be Sent By User As Request
```bash
# JSON Data Format For Register
{
    "username": "foo",
    "password": "bar",
}
# For Login
{
    "username": "foo",
    "password": "bar",
}
```
### Authentication JSON Data Format Sent Back From The Server
```bash
# JSON Data Format For Success
{
    "message": "Success",
}
# For Login
{
    "message": "Error Related To Specific Operatoin"
}
```

### 2] Fetch All Collages Data From Database
```bash
# Post Request

# Url Format For Fetching All Collage Data From Database
http://localhost:8080/admin/getusers
```

### Fetch All Collages Data From Database Data Format Sent Back From Server
```bash
# JSON Data Format For Fetch All Collages Success
{
    "message": "Success",
    "data": [
        {
            "username": "foo",
            "user_id": "jkyxt-jztor-hbsnh-h820g"
        },
    ],
}
# JSON Data Format For Fetch All Collages Error
{
    "messaage": "Error Related To Specific Operatoin"
}
```

### 3] Give Access To The User By Sending Their Mail
```bash
# Post Request

# Url Format For Sending Mail For Collage Access
http://localhost:8080/admin/giveaccess
```

###
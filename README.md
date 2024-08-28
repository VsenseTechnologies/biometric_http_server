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
### Routes Related To Admins
```bash
# General Authentication Route
http://localhost:8080/{id}/register
# Replace The {id} Either By users Or By admin
http://localhost:8080/admin/register
http://localhost:8080/users/register
# For Login
http://localhost:8080/admin/login
http://localhost:8080/users/login
```

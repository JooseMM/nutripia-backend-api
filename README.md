# TODO

- [   ] Require the token when sending the new password 
- [   ] Implement max attempts for changing password, resending email verifications

# Endpoints
## Authentication - Nutritionist
### GET /nutritionist/{id}

### POST /nutritionist/register
```json

{
	"firstname": "juan",
	"lastname": "perez",
	"emailAddress": "email@example.com",
	"password": "Password123!",
	"birthDate": "1998-08-01",
	"rut": "22222222-2"
}

```
### POST /nutritionist/login
```json

{
	"emailAddress": "email@example.com",
	"password": "Password123!"
}

```
### POST /nutritionist/send-reset-password
```json

{
	"emailAddress": "email@example.com"
}

```
### POST /nutritionist/verify-reset-password
```json

{
	"token": "super-secret-token"
}

```
### POST /nutritionist/complete-reset-password
```json

{
	"userId": "1230-123-123",
	"password": "Password123!"
}

```
### POST /nutritionist/confirm-email
```json

{
	"token": "super-secret-token"
}

```
### PUT /nutritionist/{id}
```json
{
	"firstname": "juan",
	"lastname": "perez",
	"emailAddress": "example@example.com",
	"birthDate": "1998-01-08"
}

```


## Clients
### POST /nutritionist/client",
```json
{
	"firstname": "juan",
	"lastname": "perez",
	"emailAddress": "example@example.com",
	"birthDate": "1998-01-08"


```
### GET /clients/
### GET /client/{id}
### DELETE /nutritionist/client/{id}
### PUT /nutritionist/client
```json
{
	"firstname": "juan",
	"lastname": "perez",
	"emailAddress": "example@example.com",
	"birthDate": "1998-01-08"
}
```


## Measurements
### POST /body-measurements
```json
{
    "id": 5,
    "mass": 5,
    "stature": 5,
    "sittingHeight": 5,
    "armSpan": 5,
    
    /* SkinFolds */
    "triceps": 5,
    "subscapular": 5,
    "biceps": 5,
    "iliacCrest": 5,
    "supraspinale": 5,
    "abdominal": 5,
    "frontThigh": 5,
    "medialCalf": 5,
    
    /* Girths */
    "head": 5,
    "neck": 5,
    "armRelaxed": 5,
    "armFlex": 5,
    "forearm": 5,
    "wrist": 5,
    "chest": 5,
    "waist": 5,
    "hip": 5,
    "thighHigh": 5,
    "thighLow": 5,
    "calf": 5,
    "ankle": 5,
    
    "clientId": 5
}
### GET /body-measurements/{id}
### DELETE /body-measurements/{id}

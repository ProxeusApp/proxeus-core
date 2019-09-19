# API Authentication

## Create an API Key
You need first to get an `API Key` from your Proxeus account.

The key can then be used in your HTTP header as following:

## Get a Session Token
You get a session token using your login and API Key using basic HTTP
authentication:

### Request
```
GET /api/session/token 

Authorization: BASIC BASE64(username:password)
```

### Response
```
{
    "token": "1e17b274-9f83-4348-8742-b28fb624cef6"
}
```

The username can be either
* your account email or
* the public Ethereum ID associated with your account

The returned token can be used to access the API as described below.


## Use the token to access the API
Use the token create using the step above to access the API by adding a 
bearer authorization header as the following example:

### Request
```
GET /api/user/workflow/list

Authorization: Bearer 1e17b274-9f83-4348-8742-b28fb624cef6
```

### Response
```
[
    {
        "owner": "af25eed6-aa6d-47d0-8b31-86ca845335cc",
        "groupAndOthers": {},
        "published": false,
        "id": "3784b001-e461-4ae6-879d-ff7b0d94af9c",
        "name": "test",
        "detail": "",
        "updated": "2019-09-10T14:28:33.09133+02:00",
        "created": "2019-09-06T17:18:06.790217+02:00",
        "price": 0,
        "data": null,
        "ownerEthAddress": "",
        "deactivated": false
    }
]
```

## Delete the Token
To delete the session associated with the token, use the following request:

### Request
```
DELETE /api/session/token 

Authorization: Bearer 1e17b274-9f83-4348-8742-b28fb624cef6
```


# List all Workflow

This method provides all the workflows you have access to and should be the initial method to call if the id of a specific workflow is unknown.

This request must be [authenticated](api_auth.md).


## Query

```
GET /api/document/list
```

## Response

### Success

#### Status
```
200
```

#### Body
```
[
   {
      "owner":"5ef29def-0a72-4bdb-ae69-607d13f00e9c",
      "groupAndOthers":{

      },
      "id":"3e6ece3d-6b5d-4e79-aea0-0c06e14935cb",
      "name":"hi",
      "detail":"123",
      "updated":"2019-08-17T15:11:55.383791641+02:00",
      "created":"2019-08-06T14:16:17.749839744+02:00",
      "price":0,
      "data":null,
      "ownerEthAddress":"",
      "deactivated":false
   }
]
```

### Error

```
401
404
```


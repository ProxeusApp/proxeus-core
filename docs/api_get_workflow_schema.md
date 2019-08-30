# Get Workflow Schema 

With the workflows ID we can call this method that delivers the workflow item 
including all the fields and rules. The fields and rules are under `root.workflow.data`.

This request must be [authenticated](api_auth.md).


## Query

```
GET /api/document/3e6ece3d-6b5d-4e79-aea0-0c06e14935cb/allAtOnce/schema
```

## Response

### Success

#### Status
```
200
```

#### Body
```
{
   "workflow":{
      "owner":"5ef29def-0a72-4bdb-ae69-607d13f00e9c",
      "groupAndOthers":{

      },
      "id":"3e6ece3d-6b5d-4e79-aea0-0c06e14935cb",
      "name":"hi",
      "detail":"123",
      "updated":"2019-08-17T15:11:55.383791641+02:00",
      "created":"2019-08-06T14:16:17.749839744+02:00",
      "price":0,
      "ownerEthAddress":"",
      "deactivated":false,
      "data":{
         "AutoSteer":{
            "required":false
         },
         "AutoSteer2":{
            "required":false
         },
         "AutoSteer4":{
            "required":false
         },
         "CHFXES":{
            "required":true
         },
         "ETD":{
            "datePattern":"dd.MM.yyyy",
            "required":true
         },
         "ETD4":{
            "datePattern":"dd.MM.yyyy",
            "required":true
         },
         "EyeColor":{
            "required":true
         },
         "EyeColor2":{
            "required":true
         },
         "EyeColor4":{
            "required":true
         },
         "Hurts":{
            "required":true
         },
         "Hurts2":{
            "required":true
         },
         "Hurts4":{
            "required":true
         },
         "Name":{
            "required":false
         },
         "Name2":{
            "required":false
         },
         "Name4":{
            "required":false
         },
         "SkinType":{
            "required":false
         },
         "SkinType2":{
            "required":false
         },
         "SkinType4":{
            "required":false
         },
         "Status":{
            "required":true
         },
         "Status2":{
            "required":true
         },
         "Status4":{
            "required":true
         },
         "Zodiac":{
            "required":false
         },
         "Zodiac2":{
            "required":false
         },
         "Zodiac4":{
            "required":false
         }
      }
   }
}
```

### Error

```
401
404
```

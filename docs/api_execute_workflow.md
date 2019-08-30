# Execute Workflow

Execute the workflow by providing all the required fields. Please note, no data will be stored on Proxeus side. 
You need to make sure you save the received data on your end.

This request must be [authenticated](api_auth.md).

## Query

```
POST /api/document/3e6ece3d-6b5d-4e79-aea0-0c06e14935cb/allAtOnce

{“field”:”value”, “field2”:123, “field3”:true}
```

## Response

### Success

#### Status
```
200
```

#### Header
```
Content-Type: application/pdf
```

#### Body

Single PDF in case the workflow contains just one template with the provided data. 

### Error

```
400: Bad Request
422: Some field values are not valid, please check schema field rules.
```

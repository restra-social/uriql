### URI to N1QL Query Generator ( URIQL )

This library helps to generate Couchbase N1QL Query from URL Query
Parameter based on Defined logic of the paraeter

## Example
URI : Patient?name:contains=Mr.

Defined Logic will be where to search for name paremeter

```
case "Patient" :
    sp := map[string]models.SearchParam{
        "name": models.SearchParam{
            Type:      "string",
            FieldType: "string",
            Path:      []string{"[]name.[]family", "[]name.[]given"},
        }
    }

    if val, ok := sp[match]; ok {
        return &val
    }
```

The JSON data looks like this

```
  "name": [
    {
      "family": [
        "Levin"
      ],
      "given": [
        "Henry"
      ]
    }
  ],
```

The output will be like this

```
select * from `default` as r where r.`resourceType` = 'Patient' and ANY n IN name satisfies
(any family in n.`family` satisfies family like %Mr.% end) and (any given in n.`given` satisfies given like %Mr.% end)  end;

```
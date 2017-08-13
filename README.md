### URI to Query Language Generator ( URIQL )

## Currently Supported Language are
 * N1QL ( Couchbase SQL like Query Language )
 * CYPHER ( Graph Database Neo4j )


This library helps to generate Query Language from URL Query
Parameter based on Defined logic of the paraeter

## Example
URI : Patient?name:contains=Mr.

If the JSON data looks like this

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

Dictionary will look like this. here `[]` means array and `.` means object. E.g : `[]array.obj` or `obj.obj` or `[]array.[]array`

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


The output will be like this

```
select * from `default` as r where r.`resourceType` = 'Patient' and ANY n IN name satisfies
(any family in n.`family` satisfies family like %Mr.% end) and (any given in n.`given` satisfies given like %Mr.% end)  end;

```


### Search Implementation Status

| Search Parameter Types | Sub Types | Status | Comment |
|:---:|:---:|:---:| :---: |
| Number | None | 100% | Done
| Date/DateTime | None | 80% | Period is not implemented , Confused ! |
| String | None | 100% | Done |
| Token | code | 100% | 80% |
|	| Coding | 100% | Done |
|	| CodableConcept | 100% | Done|
|	| token 	| 100% | Done |
|	| string 	| 100% | Done |
|	| Identifier 	| 50% | Done
| Reference | None | 80% |
| Reference | None | 100% | Done |
| Composite | None | 0% |
| Quantity | None | 50% |
| URI | None | 0% |
| _id | None | 0% |
| _lastUpdated | None | 0% |
| _tag | None | 0% |
| _profile | None | 0% |
| _security | None | 0% |
| _text | None | 0% |
| _content | None | 0% |
| _list | None | 0% |
| _has | None | 0% |
| _type | None | 0% |
| _query | None | 0% |
| _sort | None | 0% |
| _count | None | 0% |
| _include | None | 0% |
| _revinclude | None | 0% |
| _summary | None | 0% |
| _elements | None | 0% |
| _contained | None | 0% |
| _containedType | None | 0% |


Note : The search pattern implemented in this library is inspired by FHIR Search
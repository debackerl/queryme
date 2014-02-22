Structured language to embed complex queries in the query part of an URL.

For example the following piece of JavaScript

```JavaScript
var filter = QM.And(QM.Not(QM.Eq("type",[QM.String("foo"),QM.String("bar")])),QM.Fts("text","belgian chocolate"));
var sort = QM.Sort(QM.Order("rooms",false),QM.Order("price"));
window.location.search = "?f=" + filter + "&s=" + sort;
```

will generate the this query string

```
?f=and(not(eq(type,'foo','bar')),fts(text,'belgian chocolate'))&s=!rooms,price
```

Once received by the server, the following piece of go code

```go
qs := NewFromURL(url)
sql, values := ToSql(qs.Predicate("f"))
fmt.Println(sql)
fmt.Println(values)
```

will print the following

```
((NOT (`type` IN (?,?))) AND MATCH (`text`) AGAINST (?))
[]interface{}{"foo", "bar", "belgian chocolate"}
```

Formal specification:

```
predicates    = predicate *("," predicate)
predicate     = (not / and / or / eq / lt / le / gt / ge)
not           = "not" "(" predicate ")"
and           = "and" "(" predicates ")"
or            = "or" "(" predicates ")"
eq            = "eq" "(" field "," values ")"
lt            = "lt" "(" field "," value ")"
le            = "le" "(" field "," value ")"
gt            = "gt" "(" field "," value ")"
ge            = "ge" "(" field "," value ")"
fts           = "fts" "(" field "," string ")"

values        = value *("," value)
value         = (null / boolean / number / string / date)
null          = "null"
boolean       = "true" / "false"
number        = 1*(DIGIT / "." / "e" / "E" / "+" / "-")
string        = "'" *(unreserved / pct-encoded) "'"
date          = 4DIGIT "-" 2DIGIT "-" 2DIGIT *1("T" 2DIGIT ":" 2DIGIT ":" 2DIGIT *1("." 3DIGIT) "Z")

fieldorders   = *1(fieldorder *("," fieldorder))
fieldorder    = *1"!" field
field         = string

unreserved    = ALPHA / DIGIT / "-" / "." / "_" / "~"
pct-encoded   = "%" HEXDIG HEXDIG
sub-delims    = "!" / "$" / "&" / "'" / "(" / ")" / "*" / "+" / "," / ";" / "="
pchar         = unreserved / pct-encoded / sub-delims / ":" / "@"
query         = *( pchar / "/" / "?" )
```

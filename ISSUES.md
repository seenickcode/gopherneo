http://neo4j.com/docs/2.2.0-M01/rest-api-transactional.html#rest-api-begin-and-commit-a-transaction-in-one-request

# Silly Things Neo4j Does
1. ``http://localhost:7474/db/data`` invalid, must do ``http://localhost:7474/db/data/``

# Examples
 curl -X GET http://localhost:7474/db/data/ \
 -H "Accept: application/json; charset=UTF-8" \
 -H "Content-Type: application/json" \
 -H "Authorization: Basic realm="Neo4j" \
 OmNlZTM1YjM1NmE1MDBmNmJmZDY0MDE0NmI0ZjNhNzcx"

 curl -X POST http://localhost:7474/db/data/transaction/commit \
 -H "Accept: application/json; charset=UTF-8" \
 -H "Content-Type: application/json" \
 -H "Authorization: Basic realm="Neo4j" OmNlZTM1YjM1NmE1MDBmNmJmZDY0MDE0NmI0ZjNhNzcx" \
 --data '{
  "statements" : [ {
    "statement" : "CREATE (n) RETURN id(n)"
  } ]
}'
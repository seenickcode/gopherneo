# gopherneo
=========

A lean, pragmatic and battle-tested Neo4j driver for Go. GopherNeo ("Go4Neo") conforms strictly to the latest [Neo4j HTTP API](http://docs.neo4j.org/chunked/stable/rest-api.html), supporting only the latest Neo4j 2.x features. 

For comprehensive, up to date documentation and code examples, checkout the wiki. 

# Examples
==
(todo)

# Completed Features
==
* execute a cypher query
* get node with label, property
* get nodes with property, paginated
* create node with label, properties
* update a node
* delete node
* link nodes
* unlink nodes
* list linked nodes

# Feature Roadmap
==
## High Priority Features
* ability to return multiple entire nodes
* solution for passing order by option for listing nodes
* set relationship properties

## Medium Priority Features

## Low Priority Features

### Troubleshooting via Curl Examples

````
curl -X POST \
  -H "Accept: application/json; charset=UTF-8" \
  -H "Content-Type: application/json" \
  -d '{
    "query": "CREATE (t:Thing { props }) RETURN t",
    "params": {
      "props": {
        "name": "897430271489321",
        "age": 45
      }
    }
  }' \
  http://localhost:7474/db/data/cypher

curl -X POST \
  -H "Accept: application/json; charset=UTF-8" \
  -H "Content-Type: application/json" \
  -d '{
    "prop1": "val1"
  }' \
  http://localhost:7474/db/data/node

curl -X POST \
  -H "Accept: application/json; charset=UTF-8" \
  -H "Content-Type: application/json" \
  -d '{
    "statements": [{ 
        "statement": "CREATE (t:Thing { props }) RETURN id(t), t.name",
        "resultDataContents" : [ "REST" ],
        "parameters": {
          "props": {
            "name": "46372819647389216478321",
            "age": "45"
          }
        }
      },
      { 
        "statement": "CREATE (t:Things2 { props }) RETURN t.age",
        "resultDataContents" : [ "REST" ],
        "parameters": {
          "props": {
            "name": "46372819647389216478321",
            "age": "45"
          }
        }
      }]
  }' \
  http://localhost:7474/db/data/transaction
  
curl -X POST \
  -H "Accept: application/json; charset=UTF-8" \
  -H "Content-Type: application/json" \
  -d '{
    "statements": [{ 
        "statement": "CREATE (t:Thing { props }) RETURN t",
        "resultDataContents" : [ "REST" ],
        "parameters": {
          "props": {
            "name": "917293432",
            "age": "34"
          }
        }
      }]
  }' \
  http://localhost:7474/db/data/transaction/commit


 ````



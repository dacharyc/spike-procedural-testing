=== PLACEHOLDER NAMING CONSISTENCY ANALYSIS ===

**GOAL:**
Identify inconsistent placeholder naming conventions across documentation.
Even when placeholders are acceptable (domain-specific), inconsistent naming adds cognitive load.

**METHODOLOGY:**
Group similar placeholders by concept (collection, database, field, etc.)
and show all the different ways we represent the same concept.

---

## PLACEHOLDER NAMING INCONSISTENCIES

Showing different ways we represent the same concept:

### COLLECTION (14 variations, 303 total usages)

**INCONSISTENT** - 14 different ways to represent the same concept:
  - `<collection>`: 89 (29.4%)
  - `<collectionName>`: 84 (27.7%)
  - `<collection-name>`: 36 (11.9%)
  - `<collection name>`: 35 (11.6%)
  - `<MongoCollection setup code here>`: 17 (5.6%)
  - `<COLLECTION-NAME>`: 12 (4.0%)
  - `<your collection>`: 8 (2.6%)
  - `<your collection to read from>`: 6 (2.0%)
  - `<your collection to write to>`: 5 (1.7%)
  - `<collName>`: 4 (1.3%)
  - `<MongoCollection set up code here>`: 3 (1.0%)
  - `<myCollection>`: 2 (0.7%)
  - `<capped collection>`: 1 (0.3%)
  - `<targetCollectionName>`: 1 (0.3%)

### DATABASE (28 variations, 425 total usages)

**INCONSISTENT** - 28 different ways to represent the same concept:
  - `<database>`: 77 (18.1%)
  - `<databaseName>`: 75 (17.6%)
  - `<db_username>`: 64 (15.1%)
  - `<db_password>`: 64 (15.1%)
  - `<database-name>`: 35 (8.2%)
  - `<database name>`: 24 (5.6%)
  - `<Your MongoDB URI>`: 19 (4.5%)
  - `<your mongodb uri>`: 11 (2.6%)
  - `<DATABASE-NAME>`: 11 (2.6%)
  - `<mongodb-connection-string>`: 9 (2.1%)
  - `<your database to read from>`: 6 (1.4%)
  - `<your database to write to>`: 5 (1.2%)
  - `<your database>`: 4 (0.9%)
  - `<dbName>`: 3 (0.7%)
  - `<db>`: 2 (0.5%)
  - `<new_mongodb_user>`: 2 (0.5%)
  - `<MONGODB_URI>`: 2 (0.5%)
  - `<dbname>`: 2 (0.5%)
  - `<authDB>`: 1 (0.2%)
  - `<DATABASES key>`: 1 (0.2%)
  - `<mongodb connection string>`: 1 (0.2%)
  - `<Your MongoDB Connection URI>`: 1 (0.2%)
  - `<Database Name>`: 1 (0.2%)
  - `<MongoDB Search index name>`: 1 (0.2%)
  - `<sourceDatabaseName>`: 1 (0.2%)
  - `<targetDatabaseName>`: 1 (0.2%)
  - `<db-name>`: 1 (0.2%)
  - `<MongoDB connection URI>`: 1 (0.2%)

### FIELD (48 variations, 386 total usages)

**INCONSISTENT** - 48 different ways to represent the same concept:
  - `<field name>`: 100 (25.9%)
  - `<field-name>`: 84 (21.8%)
  - `<field to match>`: 27 (7.0%)
  - `<fieldName>`: 20 (5.2%)
  - `<fieldToIndex>`: 20 (5.2%)
  - `<field-type>`: 18 (4.7%)
  - `<field-mapping-definition>`: 12 (3.1%)
  - `<field>`: 9 (2.3%)
  - `<sub-field-name>`: 9 (2.3%)
  - `<string-field-definition>`: 8 (2.1%)
  - `<field_name>`: 6 (1.6%)
  - `<field to index>`: 6 (1.6%)
  - `<TEXT-FIELD-NAME>`: 5 (1.3%)
  - `<FIELD_NAME>`: 4 (1.0%)
  - `<field-to-include>`: 4 (1.0%)
  - `<field-to-exclude>`: 4 (1.0%)
  - `<FIELD-NAME>`: 4 (1.0%)
  - `<field1>`: 3 (0.8%)
  - `<field2>`: 3 (0.8%)
  - `<field name 1>`: 3 (0.8%)
  - `<field name 2>`: 3 (0.8%)
  - `<field3>`: 2 (0.5%)
  - `<new document field name>`: 2 (0.5%)
  - `<array field name>`: 2 (0.5%)
  - `<new field name>`: 2 (0.5%)
  - `<FIELD-NAME-FOR-FLOAT32-TYPE>`: 2 (0.5%)
  - `<FIELD-NAME-FOR-INT8-TYPE>`: 2 (0.5%)
  - `<FIELD-NAME-FOR-INT1-TYPE>`: 2 (0.5%)
  - `<outputfield1>`: 1 (0.3%)
  - `<array field>`: 1 (0.3%)
  - `<arrayField>`: 1 (0.3%)
  - `<GeoJSON object field>`: 1 (0.3%)
  - `<field type>`: 1 (0.3%)
  - `<updated field name>`: 1 (0.3%)
  - `<updated field type>`: 1 (0.3%)
  - `<fieldType>`: 1 (0.3%)
  - `<fieldName1>`: 1 (0.3%)
  - `<fieldName2>`: 1 (0.3%)
  - `<GeoJSON object field name>`: 1 (0.3%)
  - `<field-type-definition>`: 1 (0.3%)
  - `<field-type-configuration>`: 1 (0.3%)
  - `<field-definition>`: 1 (0.3%)
  - `<fieldDefinition>`: 1 (0.3%)
  - `<DATA-FIELD-NAME>`: 1 (0.3%)
  - `<field-to-index>`: 1 (0.3%)
  - `<field-to-search>`: 1 (0.3%)
  - `<metadata-field>`: 1 (0.3%)
  - `<indexed-field-to-search>`: 1 (0.3%)

### INDEX (25 variations, 133 total usages)

**INCONSISTENT** - 25 different ways to represent the same concept:
  - `<indexName>`: 37 (27.8%)
  - `<index name>`: 27 (20.3%)
  - `<INDEX-NAME>`: 18 (13.5%)
  - `<index-name>`: 16 (12.0%)
  - `<INDEX_NAME>`: 6 (4.5%)
  - `<firstIndexName>`: 4 (3.0%)
  - `<secondIndexName>`: 3 (2.3%)
  - `<index id>`: 2 (1.5%)
  - `<searchIndexName>`: 2 (1.5%)
  - `<index to update>`: 2 (1.5%)
  - `<first-index-name>`: 2 (1.5%)
  - `<indexId>`: 1 (0.8%)
  - `<search-index-name>`: 1 (0.8%)
  - `<your_search_index_name>`: 1 (0.8%)
  - `<Vector Search index name>`: 1 (0.8%)
  - `<search index name>`: 1 (0.8%)
  - `<first index name>`: 1 (0.8%)
  - `<second index name>`: 1 (0.8%)
  - `<second-index-name>`: 1 (0.8%)
  - `<indexDefinition>`: 1 (0.8%)
  - `<last-index-name>`: 1 (0.8%)
  - `<lastIndexName>`: 1 (0.8%)
  - `<index-definition>`: 1 (0.8%)
  - `<search_index_model>`: 1 (0.8%)
  - `<indexAnalyzer>`: 1 (0.8%)

### CONNECTIONSTRING (25 variations, 994 total usages)

**INCONSISTENT** - 25 different ways to represent the same concept:
  - `<connection-string>`: 669 (67.3%)
  - `<connection string>`: 103 (10.4%)
  - `<connectionString>`: 86 (8.7%)
  - `<connection string uri>`: 38 (3.8%)
  - `<connection string URI>`: 17 (1.7%)
  - `<cluster-url>`: 13 (1.3%)
  - `<your schema registry uri>`: 12 (1.2%)
  - `<CONNECTION-STRING>`: 12 (1.2%)
  - `<your connection string here>`: 7 (0.7%)
  - `<Your Connection String>`: 6 (0.6%)
  - `<your connection uri>`: 5 (0.5%)
  - `<atlas-connection-string>`: 4 (0.4%)
  - `<connection URI>`: 4 (0.4%)
  - `<clusterUrl>`: 3 (0.3%)
  - `<your connection string>`: 3 (0.3%)
  - `<your connection URI>`: 2 (0.2%)
  - `<your Atlas connection string>`: 2 (0.2%)
  - `<other connection string>`: 1 (0.1%)
  - `<opsManagerUri>`: 1 (0.1%)
  - `<your-connection-string>`: 1 (0.1%)
  - `<YOUR-ATLAS-URI>`: 1 (0.1%)
  - `<connection_string>`: 1 (0.1%)
  - `<Your connection URI>`: 1 (0.1%)
  - `<atlas-uri>`: 1 (0.1%)
  - `<connection uri>`: 1 (0.1%)

### USERNAME (5 variations, 38 total usages)

**INCONSISTENT** - 5 different ways to represent the same concept:
  - `<username>`: 29 (76.3%)
  - `<user>`: 4 (10.5%)
  - `<IAM User Access Key ID>`: 2 (5.3%)
  - `<IAM User Secret Access Key>`: 2 (5.3%)
  - `<userProfileCity>`: 1 (2.6%)

### PASSWORD (4 variations, 48 total usages)

**INCONSISTENT** - 4 different ways to represent the same concept:
  - `<password>`: 45 (93.8%)
  - `<pwd>`: 1 (2.1%)
  - `<generated-password>`: 1 (2.1%)
  - `<app password>`: 1 (2.1%)

### APIKEY (8 variations, 24 total usages)

**INCONSISTENT** - 8 different ways to represent the same concept:
  - `<api-key>`: 7 (29.2%)
  - `<openai-api-key>`: 4 (16.7%)
  - `<VOYAGE-API-KEY>`: 4 (16.7%)
  - `<voyage-api-key>`: 3 (12.5%)
  - `<VOYAGEAI-API-KEY>`: 3 (12.5%)
  - `<voyageai-api-key>`: 1 (4.2%)
  - `<publicApiKey>`: 1 (4.2%)
  - `<privateApiKey>`: 1 (4.2%)

### TOKEN (2 variations, 16 total usages)

**INCONSISTENT** - 2 different ways to represent the same concept:
  - `<aws session token>`: 14 (87.5%)
  - `<hf-token>`: 2 (12.5%)

### HOSTNAME (2 variations, 63 total usages)

**INCONSISTENT** - 2 different ways to represent the same concept:
  - `<hostname>`: 61 (96.8%)
  - `<host>`: 2 (3.2%)

### PORT (2 variations, 51 total usages)

**INCONSISTENT** - 2 different ways to represent the same concept:
  - `<port-number>`: 37 (72.5%)
  - `<port>`: 14 (27.5%)

### VERSION (1 variations, 29 total usages)

  - `<version>`: 29 (100.0%)

### DEPENDENCY (5 variations, 94 total usages)

**INCONSISTENT** - 5 different ways to represent the same concept:
  - `<dependency>`: 26 (27.7%)
  - `<groupId>`: 26 (27.7%)
  - `<artifactId>`: 26 (27.7%)
  - `<dependencies>`: 11 (11.7%)
  - `<dependencyManagement>`: 5 (5.3%)

## LANGUAGE-SPECIFIC PLACEHOLDER ANALYSIS

Do different programming languages use different placeholder conventions?

### PYTHON (539 placeholders, 133 unique)

**Naming conventions used:**
  - kebab-case: 212 (39.3%)
  - space separated: 197 (36.5%)
  - mixed/other: 59 (10.9%)
  - camelCase: 58 (10.8%)
  - snake_case: 9 (1.7%)
  - Title Case: 2 (0.4%)
  - UPPER_CASE: 2 (0.4%)

**Most common placeholders:**
  - `<connection-string>`: 112
  - `<field name>`: 31
  - `<value>`: 26
  - `<database-name>`: 20
  - `<collection-name>`: 20

### JAVA (501 placeholders, 103 unique)

**Naming conventions used:**
  - space separated: 208 (41.5%)
  - kebab-case: 175 (34.9%)
  - camelCase: 53 (10.6%)
  - mixed/other: 32 (6.4%)
  - PascalCase: 22 (4.4%)
  - Title Case: 11 (2.2%)

**Most common placeholders:**
  - `<connection-string>`: 111
  - `<field name>`: 26
  - `<MongoCollection setup code here>`: 17
  - `<value>`: 16
  - `<your schema registry uri>`: 12

### JAVASCRIPT (445 placeholders, 116 unique)

**Naming conventions used:**
  - space separated: 155 (34.8%)
  - kebab-case: 111 (24.9%)
  - camelCase: 73 (16.4%)
  - mixed/other: 69 (15.5%)
  - UPPER_CASE: 14 (3.1%)
  - Title Case: 12 (2.7%)
  - snake_case: 10 (2.2%)
  - PascalCase: 1 (0.2%)

**Most common placeholders:**
  - `<connection-string>`: 85
  - `<collection>`: 18
  - `<version>`: 17
  - `<connection string uri>`: 16
  - `<database>`: 13

### C (286 placeholders, 46 unique)

**Naming conventions used:**
  - space separated: 99 (34.6%)
  - camelCase: 92 (32.2%)
  - kebab-case: 84 (29.4%)
  - mixed/other: 6 (2.1%)
  - PascalCase: 2 (0.7%)
  - snake_case: 2 (0.7%)
  - Title Case: 1 (0.3%)

**Most common placeholders:**
  - `<connection-string>`: 35
  - `<connectionString>`: 32
  - `<collectionName>`: 28
  - `<field-name>`: 18
  - `<field name>`: 15

### GO (283 placeholders, 59 unique)

**Naming conventions used:**
  - kebab-case: 149 (52.7%)
  - space separated: 73 (25.8%)
  - camelCase: 44 (15.5%)
  - Title Case: 9 (3.2%)
  - mixed/other: 6 (2.1%)
  - snake_case: 2 (0.7%)

**Most common placeholders:**
  - `<connection-string>`: 111
  - `<field-name>`: 16
  - `<collection>`: 12
  - `<dependency>`: 12
  - `<groupId>`: 12

### CSHARP (282 placeholders, 57 unique)

**Naming conventions used:**
  - kebab-case: 118 (41.8%)
  - camelCase: 89 (31.6%)
  - space separated: 47 (16.7%)
  - mixed/other: 14 (5.0%)
  - Title Case: 10 (3.5%)
  - snake_case: 2 (0.7%)
  - UPPER_CASE: 2 (0.7%)

**Most common placeholders:**
  - `<connection-string>`: 97
  - `<databaseName>`: 16
  - `<collectionName>`: 15
  - `<fieldToIndex>`: 14
  - `<connectionString>`: 13

### UNDEFINED (197 placeholders, 27 unique)

**Naming conventions used:**
  - snake_case: 66 (33.5%)
  - space separated: 55 (27.9%)
  - kebab-case: 43 (21.8%)
  - camelCase: 32 (16.2%)
  - mixed/other: 1 (0.5%)

**Most common placeholders:**
  - `<db_username>`: 33
  - `<db_password>`: 33
  - `<port-number>`: 32
  - `<clusterName>`: 31
  - `<hostname>`: 31

### CPP (194 placeholders, 36 unique)

**Naming conventions used:**
  - kebab-case: 81 (41.8%)
  - space separated: 73 (37.6%)
  - camelCase: 38 (19.6%)
  - mixed/other: 1 (0.5%)
  - snake_case: 1 (0.5%)

**Most common placeholders:**
  - `<connection-string>`: 44
  - `<field-name>`: 15
  - `<connection string>`: 14
  - `<field name>`: 13
  - `<value>`: 13

### TEXT (153 placeholders, 30 unique)

**Naming conventions used:**
  - space separated: 104 (68.0%)
  - snake_case: 26 (17.0%)
  - kebab-case: 11 (7.2%)
  - camelCase: 7 (4.6%)
  - mixed/other: 3 (2.0%)
  - Title Case: 2 (1.3%)

**Most common placeholders:**
  - `<other options>`: 28
  - `<aws access key id>`: 21
  - `<aws secret access key>`: 21
  - `<aws session token>`: 14
  - `<db_username>`: 12

### SCALA (94 placeholders, 32 unique)

**Naming conventions used:**
  - space separated: 65 (69.1%)
  - mixed/other: 13 (13.8%)
  - kebab-case: 12 (12.8%)
  - snake_case: 2 (2.1%)
  - camelCase: 2 (2.1%)

**Most common placeholders:**
  - `<value>`: 11
  - `<field name>`: 10
  - `<connection string>`: 9
  - `<field to match>`: 7
  - `<value to match>`: 6

# otf-align
Web-Service to determine alignment of an assessment result to NLPs

# usage
otf-align is a web service that accepts simple requests for alignment of assessment data to the National Learning Progressions.

The otf-align API is simple:

```
> curl -v http://localhost:1324/align \
  -H 'Content-Type: application/json' \
    -d '{"alignMethod":"inferred", \
         "alignCapability":"literacy", \
         "alignToken":"answers questions confidently"}'
```

the three required parameters are:
- alignMethod: choice of mapped | inferred | prescribed (see full description below)
- alignCapability: the General Capability area of the NLPs that this measurement belongs to (required whn using the inferred method), currently must be one of literacy or numeracy.
- alignToken: the text of an observation or quesstion (for inference), the identifier of a question/module from the source system (for mapped), or the identifier of an element/sub-element/development-level/indicator from the NLPs (for prescribed).

the otf-align service will respond on success with the following data structure:
```
{
  "alignCapability": "literacy",
  "alignMethod": "inferred",
  "alignServiceID": "ygd1RcKF2k1MNFLm8eZ7Nn",
  "alignServiceName": "o5lZbn",
  "alignToken": "answers questions confidently",
  "alignments": [
    {
      "developmentLevel": "UnT3",
      "element": "Reading and viewing",
      "generalCapability": "Literacy",
      "heading": "Comprehension",
      "indicator": "infers and then describes obvious cause and effect relationships (e.g. uses information in the text to infer why a character is smiling in an image)",
      "itemID": "uri/version/d30bf6bb-4f31-4182-a696-bc20c711e09f",
      "itemText": "answers and poses mainly literal questions about the text",
      "progressionLevel": "UnT3",
      "subElement": "Understanding texts"
    }
  ]
}
```
The response echoes the input parameters for completeness, and identifies the service instance that processed the request.

The *alignments* element contains an array of GESDI blocks containing the full resolution of the aligned item.

Note the *alignments* element is an array to accommodate the fact that some alignments may produce a many-to-one relationship between the input token and the

All configuration options can be set on the command-line using flags, via envronment variables, or by using a configuration file.
Configuration can use any or all of these methods in combination.
For example options such as the address and hostname of the classifier server might best be accessed from environment variables, whilst the service name of the otf-align instance might be supplied in a json configuration file.

Configuration flags are capitalised and prefixed with OTF_ALIGN_SRVC when supplied as environment variables; so flag --niasPort on the commnad-line becomes 

```
OTF_ALIGN_SRVC_NIASPORT=1323
```

when expressed as an environment variable and

```
{ "niasPort":1323 }
```

when set in a json configuration file.

These are the configuration options:

|Option name|Type|Required|Default|Description|
|---|---|---|---|---|
|config|string|no||configuration file name|
|name|string|yes|auto-generated (hashid)|name of this instance of the service|
|id|string|yes|auto-generated (nuid)|identifier for this service instance|  
|host|string|yes|localhost|host address to run this service on|
|port|int|yes|auto-generated|port to run the service on|
|niasHost|string|yes|localhost|host of n3w service|
|niasPort|int|yes|1323|port of the n3w service|
|niasToken|string|yes|a demo|jwt token for accessing the n3w server|
|tcHost|string|yes|localhost|host address for text classification service|
|tcPort|int|yes|1576|port classifier service runs on|    

# alignment methods
otf-align is a facade service which will invoke further services in order to determine the alignment of a particular assessment result or observation.
Three styles of alignment resolution are currently supported:
- mapped 
    -  alignment is resolved by mapping tokens from the original observation/assessment data to existing data structures that themselves are linked to the NLPS.
  - for example, an assessment result may contain a test or module identifier provided by the assessment system. That identifier may already be mapped to an external structure such as the Australian Curriculum. The Australian Curriculum has a set of pre-defined relationships between its structure and the NLPs. By traversing the set of known (mapped) relationships otf-align can determine the alignment of the original assessment to the NLPs
  - the map of relationships is held as graph nodes in a nias3 datastore which is queried by the otf-align facade service when a mapped alignment is requested.
- inferred
    - alignment is resolved by passing the token to a text classification service which has been populated with the NLPs
    - the classification service will find the closest matches between the submitted token and the NLPs
    - currently the otf-align service filters the list of results from the classifier and returns the top match only.
- prescribed    
    - if the submitted token provided is a reference to the NLPs, then the service will return the full GESDI block for that token.

# pre-requisites
The otf-align service requires supporting services to be available:
- otf-classifier, provides classification engine and NLP lookup service
    + binary can be created from http://github.com/nsip/otf-calssifier
- n3w, provides the lookup graphs for mapped alignments
    + binary can be created from http://github.com/nsip/n3-web
- benthos, workflow engine installed and available
- nats-streaming-server, message broker installed and avialable

# benthos workflow
As with otf-reader the service is packaged with a benthos workflow which allows the testing of otf-align in context and interacting with the other progress data management services.

The provided script in the /benthos folder will read data that has been posted to a nats-streaming-server by instances of the otf-reader.

The script sends requests from the reader to the alignment service, and then inserts the response into an *otf* block within the original message, which will then be passed on to other services in the PDM workflow.

The current benthos script will output the results of the alignment as individual json files to the benthos/msgs folder.

To set up the otf infrastructure to run the end to end test:

start an instance of nats-streaming-server:
```
> ./nats-streaming-server
```
start the classifier wherever the binary is located:
```
> ./otf-classifier
```
start the benthos workflow:
```
otf-align/cmd/benthos> ./run_benthos.sh
```
start the otf-align server
```
otf-align/cmd/otf-align> ./otf-align --niasToken=xxxyyy --port=1324
```
(port 1324 is the port used in the benthos script)

start an instance of otf-reader to begin the ingest process (in this example using the MathsPathway configuration):
```
otf-reader/cmd/otf-reader> ./otf-reader --config=./config/mp_config.yaml
```
then copy the MathsPathway.csv file from the *cmd/otf-reader/test_data* folder to the *cmd/otf-reader/in/maths-pathway* folder.

The reader will convert the csv file to json, add all necessary meta-data and publish the message to nats.

The benthos workflow will consume the messsages, extract the tokens and method from the data, submit the request to the alignment service, and map the results back into the standard format of the message.

Messages in the complete format will be written to the *cmd/otf-reader/benthos/msgs* folder for checking.

A fully processed message example:
```
{
    "meta":
    {
        "alignMethod": "mapped",
        "capability": "numeracy",
        "inputFormat": "csv",
        "levelMethod": "prescribed",
        "providerName": "MathsPathway",
        "readerID": "ApQHSHmoEIfd08pWOCC5qz",
        "readerName": "pA9em"
    },
    "original":
    {
        "available": true,
        "completed": true,
        "diagnosed": true,
        "level": 4,
        "mastered": true,
        "module_id": "00e6a88e-f481-4984-8edb-a7f6b95e23c0",
        "student_id": "ac4f28ee-486a-4672-ade2-0fb332c10995"
    },
    "otf":
    {
        "alignmentServiceID": "lIvBYJ79X9M10yo5bBG8yZ",
        "alignmentServiceName": "RQEzxG",
        "alignments": [
        {
            "developmentLevel": "AdS7",
            "element": "Number sense and algebra",
            "generalCapability": "Numeracy",
            "heading": "Flexible strategies with two-digit numbers",
            "indicator": "represents a wide range of additive problem situations involving two-digit numbers using appropriate addition and subtraction number sentences",
            "itemID": "uri/version/ec3c0b7f-a190-4e79-84ab-4f4d8823c698",
            "itemText": "chooses from a range of known strategies to solve additive problems involving two-digit numbers (e.g. uses place value knowledge, known facts and part-part-whole number knowledge to solve problems like 24 + 8 + 13, partitioning 8 as 6 and 2 more, then combining 24 and 6 to rename it as 30 before combining it with 13 to make 43, and then combining the remaining 2 to find 45 �; adding the same to both numbers 47 � 38 = 49 - 40)",
            "progressionLevel": "AdS7",
            "subElement": "Additive strategies"
        }],
        "capability": "numeracy",
        "id":
        {
            "studentID": "ac4f28ee-486a-4672-ade2-0fb332c10995"
        },
        "method": "mapped",
        "token": "00e6a88e-f481-4984-8edb-a7f6b95e23c0"
    }
}
```












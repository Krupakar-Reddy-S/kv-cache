# **Assignment: Key-Value Cache** 

## **Overview**

In this assignment, you will build an in-memory Key-Value Cache service that implements the basic operations:

* **put(key, value):** Inserts or updates a key–value pair.  
* **get(key):** Retrieves the value for a given key, if present.

Your service will be containerized using Docker and must listen on port **7171**. Implementations will be run on identical servers (AWS t3.small \- 2 core, 2 GB RAM) where a distributed locust run will simulate a random sequence of `put()` and `get()` operations. Based on their performance, all submissions will be ranked on a leaderboard.

## **Requirements**

### **Functional Specifications**

* **API Operations:**  
  * **put(key, value):**  
    * Accepts a key and a value.  
    * If the key already exists, update the associated value.  
      * Else create a new key, value mapping.   
    * Return an acknowledgment upon success.  
  * **get(key):**  
    * Returns the value associated with the key if present.  
    * If the key does not exist, return a “not found” response.  
  * API contracts are mentioned in the Appendix.  
* **Data Constraints:**  
  * Keys and values should be treated as strings.  
  * Limit key and value lengths to a maximum of 256 ASCII characters.  
* **Data Persistence:**  
  * The cache is in-memory only. Data persistence between service restarts is not required.

### **Non-Functional Specifications**

* **Performance:**  
  * Your implementation should maintain low latency under a heavy load of random `get()` and `put()` operations.  
* **Deployment:**  
  * Package your application as a Docker image.  
  * The service must listen on port **7171**.  
  * Your Docker image should run via following command:

docker run \-p 7171:7171 \<your\_image\_name\>

* **Environment:**  
  * The service will be deployed on an AWS t3.small \- 2 core, 2 GB RAM.  
* **Documentation:**  
  * Include a README file explaining:  
    * How to build and run your Docker image.  
    * \[MUST\] Design choices or optimizations you implemented.  
* **Testing:**  
  * Ensure your implementation never leads to an Out of Memory issue  
    * call cache eviction/ use other strategies to overcome this  
  * Ensure unless the resource usage is \> 70%, cache miss for keys put till then is 0%

## **Leaderboard & Evaluation**

Submissions will be automatically deployed and tested using a load testing script. Based on the observe performance (average latency of operations), marks will be assigned as follows:

* **Top 10:** 100 marks  
* **Top 25:** 75 marks  
* **Accurate Submissions (meeting all requirements (including README) but not in the top 25):** 50 marks  
* **Else (failing to achieve accurate results/ getting memory overflow/ cache miss etc):** 0 marks

Top 3 submissions will also get SST merch (placeholder 1 will get a hoodie, and 2,3 will get a t-shirt). 

## **Submission Guidelines**

* **Fill this form to submit your assignment:** [https://forms.gle/QVYS3mjR3u7kA3Cr8](https://forms.gle/QVYS3mjR3u7kA3Cr8)   
* The Github Repository should be viewable via [https://github.com/anshumansingh](https://github.com/anshumansingh) and [https://github.com/Naman-Bhalla/vc\_domain\_research\_agent](https://github.com/Naman-Bhalla/vc_domain_research_agent)  
* The Docker Image should be publicly accessible and present on Docker Hub.  
* The Docker Image should run as per instructions shared above.  
* **Deadline:** 16 Mar 2025

## **Additional Details**

* **Language and Frameworks:** You are free to choose your programming language and any supporting frameworks, as long as the final application adheres to the API and resource requirements.  
* **Resource Management:** Ensure that your solution operates within the given memory constraints.  
* **Load Test Environment:** The load testing framework will simulate a mix of `put()` and `get()` requests with randomized keys and values. Consider potential edge cases, such as high-frequency updates, a lot of keys, etc.

---

Happy coding\!

## Appendix

### API Endpoints

#### 1\. PUT Operation

* HTTP Method: `POST`  
* Endpoint: `/put`
Assignment: Key-Value Cache 
Overview
In this assignment, you will build an in-memory Key-Value Cache service that implements the basic operations:
put(key, value): Inserts or updates a key–value pair.
get(key): Retrieves the value for a given key, if present.
Your service will be containerized using Docker and must listen on port 7171. Implementations will be run on identical servers (AWS t3.small - 2 core, 2 GB RAM) where a distributed locust run will simulate a random sequence of put() and get() operations. Based on their performance, all submissions will be ranked on a leaderboard.
Requirements
Functional Specifications
API Operations:
put(key, value):
Accepts a key and a value.
If the key already exists, update the associated value.
Else create a new key, value mapping. 
Return an acknowledgment upon success.
get(key):
Returns the value associated with the key if present.
If the key does not exist, return a “not found” response.
API contracts are mentioned in the Appendix.
Data Constraints:
Keys and values should be treated as strings.
Limit key and value lengths to a maximum of 256 ASCII characters.
Data Persistence:
The cache is in-memory only. Data persistence between service restarts is not required.
Non-Functional Specifications
Performance:
Your implementation should maintain low latency under a heavy load of random get() and put() operations.
Deployment:
Package your application as a Docker image.
The service must listen on port 7171.
Your Docker image should run via following command:
docker run -p 7171:7171 <your_image_name>
Environment:
The service will be deployed on an AWS t3.small - 2 core, 2 GB RAM.
Documentation:
Include a README file explaining:
How to build and run your Docker image.
[MUST] Design choices or optimizations you implemented.
Testing:
Ensure your implementation never leads to an Out of Memory issue
call cache eviction/ use other strategies to overcome this
Ensure unless the resource usage is > 70%, cache miss for keys put till then is 0%
Leaderboard & Evaluation
Submissions will be automatically deployed and tested using a load testing script. Based on the observe performance (average latency of operations), marks will be assigned as follows:
Top 10: 100 marks
Top 25: 75 marks
Accurate Submissions (meeting all requirements (including README) but not in the top 25): 50 marks
Else (failing to achieve accurate results/ getting memory overflow/ cache miss etc): 0 marks
Top 3 submissions will also get SST merch (placeholder 1 will get a hoodie, and 2,3 will get a t-shirt). 
Submission Guidelines
Fill this form to submit your assignment: https://forms.gle/QVYS3mjR3u7kA3Cr8 
The Github Repository should be viewable via https://github.com/anshumansingh and https://github.com/Naman-Bhalla/vc_domain_research_agent
The Docker Image should be publicly accessible and present on Docker Hub.
The Docker Image should run as per instructions shared above.
Deadline: 16 Mar 2025
Additional Details
Language and Frameworks: You are free to choose your programming language and any supporting frameworks, as long as the final application adheres to the API and resource requirements.
Resource Management: Ensure that your solution operates within the given memory constraints.
Load Test Environment: The load testing framework will simulate a mix of put() and get() requests with randomized keys and values. Consider potential edge cases, such as high-frequency updates, a lot of keys, etc.

Happy coding!

Appendix
API Endpoints
1. PUT Operation
HTTP Method: POST
Endpoint: /put
Request

Body Format:
{
  "key": "string (max 256 characters)",
  "value": "string (max 256 characters)"
}
Response

On Success (HTTP 200):
{
  "status": "OK",
  "message": "Key inserted/updated successfully."
}

On Failure:
{
  "status": "ERROR",
  "message": "Error"
}

2. GET Operation
HTTP Method: GET
Endpoint: /get
Request
Parameters: A query parameter named key
Example URL: /get?key=exampleKey
Response
On Success (HTTP 200):
{
  "status": "OK",
  "key": "exampleKey",
  "value": "the corresponding value"
}

If Key Not Found:
{
  "status": "ERROR",
  "message": "Key not found."
}

On Other Failures:
{
  "status": "ERROR",
  "message": "Error description explaining what went wrong."
}



##### Request

Body Format:  
`{`  
  `"key": "string (max 256 characters)",`  
  `"value": "string (max 256 characters)"`  
`}`

##### Response

On Success (HTTP 200):  
`{`  
  `"status": "OK",`  
  `"message": "Key inserted/updated successfully."`  
`}`

On Failure:  
`{`  
  `"status": "ERROR",`  
  `"message": "Error"`  
`}`  
---

#### 2\. GET Operation

* HTTP Method: `GET`  
* Endpoint: `/get`

##### Request

* Parameters: A query parameter named `key`  
* Example URL: `/get?key=exampleKey`

##### Response

On Success (HTTP 200):  
`{`  
  `"status": "OK",`  
  `"key": "exampleKey",`  
  `"value": "the corresponding value"`  
`}`

If Key Not Found:  
`{`  
  `"status": "ERROR",`  
  `"message": "Key not found."`  
`}`

On Other Failures:  
`{`  
  `"status": "ERROR",`  
  `"message": "Error description explaining what went wrong."`  
`}`  
---


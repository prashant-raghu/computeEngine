# ComputeEngine

A Docker based sandboxing mechanism to safely execute untrusted user code in an isolated environment and prevent remote code execution vulnerability and other nasty things one can think of.

[![N|Solid](https://raw.githubusercontent.com/prashant-raghu/computeEngine/master/assets/Arch.png)](https://github.com/prashant-raghu/computeEngine/)

# How it works

For each code execution request here is what computeEngine does:
* Roll up a new Docker Container.
* Mount a shared Directory from main server to this Docker Container.
* Write the code to be executed along with execution instructions.
* Container Executes the code present in shared Directory.
* Result is written into a new file created by container witin the shared directory.
* Main server watches for creation of this file and once done Responds to user with the result.


### Development

> prerequisites: Docker and golang.

1\) Build Docker Image from dockerfile 
```sh
$ cd sandbox
$ docker image build -t sandbox:v1 .
```
2\) Start the go server
```sh
$ go run main.go
```
3\) Verify the deployment by navigating to:
```sh
127.0.0.1:5000
```

4\) To try out the code execution part make a POST request to \execute with "code"(Js for now) in body of request:
```sh
127.0.0.1:5000/execute
```
Something like this:
[![N|Solid](https://raw.githubusercontent.com/prashant-raghu/computeEngine/master/assets/postExecute.png)](https://github.com/prashant-raghu/computeEngine/)


### Todos
 - Add Timeout docker kill to executeHandler.
 - Start writing Tests.
 - Add multi Language Support.

License
----
MIT


**Free Code, Hell Yeah!**

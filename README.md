# Code Execution Engine for Firefly

Receives a job to execute a code and returns the data to the database.

## Demo
You can send a POST request to https://run.joshjms.com/submit with the following body in JSON:
* `language`: currently supports only C, C++, and Python.
* `code`: your code.
* `mem_limit`: memory limit for execution (maxrss) in KB.
* `time_limit`: time limit for execution (cpu time) in ms.
* `stdin`: standard input for your code.
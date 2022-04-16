You need to implement a **Client-Server application** with the following requirements:
* multiple-threaded server;
* clients;
* External queue between the clients and server;

Clients:
* Should be configured from a command line or from a file (you decide);
* Can read data from a file or from a command line (you decide);
* Can request server to AddItem(), RemoveItem(), GetItem(), GetAllItems()
* Data is in the form of strings;

* Clients can be added / removed while not intefring to the server or other clients ;

Server:
* Has data structure(s) that holds the data in the memory while keeping the order of items as they added (Ordered Map for C++);
  - The data structure must keep the order of items as they added. 
    For example: If client added the following keys in the following order A, B, D, E, C. 
    The GetAllItems returns A, B, D, E, C
	If item D was removed, the GetAllItems return A, B, E, C
* Server should be able to add an item, remove an item, get a single or all item from the data structure;

External queue:
* Can be Amazon Simple Queue Service (SQS) or RabbitMQ (you decide);


Clients send requests to the external queue - while the server reads those and execute them on its data structure. You define the structure of the messages (AddItem, RemoveItem, GetItem, GetAllItems)


The flow of the project:
1. Multiple clients are sending requests to the queue (and not waiting for the response).
2. Server is reading requests from the queue and processing them, the output of the server is written to a log file
3. Server should be able to process items in parallel
4. log messages (debug, error) are written to stdout

   
Definition of success:
* Working project that can be executed on your computer (preferred OS = linux);
* Being able to explain how the project works and how to deploy the project (for the first time) on another computer;
* If you take something from the Internet or consult anyone for anything, you should be able to understand it perfectly;
* Code has no bugs, no dangling references / assets / resources, no resource leaks;
* Code is clean and readable;
* Code is reasonably efficient (server idle time will be measured).
* Working with channels when needed
* You implement the data structue(s) by yourself


You should develop the project using GOLang.

Good luck!
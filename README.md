# Gnet Scanner
A multithreaded CLI tool for scanning the network and collecting results, written in Golang.
## Assembly and start-up
All you need to do to build a programme is go build in the repository directory.
```go build```
## Flags
1) ```--Network={x.x.x.x/y}```
   Network and its mask to be scanned. Only CIDR format is supported, band format is not supported yet. Mandatory flag if --save flag is not specified (these two flags cannot be used together).
2) ```--Save={"path/to/file/save data.gsv"}```
   The path to the save file. Mandatory if --network flag is not specified (these two flags cannot be used together).
3) ```--Debug={true/false} (default false)```
   This flag enables debugging mode. More specifically, it starts a profiler on localhost:6161 and a special leak detector that saves a dump to the `heap` file when 1GiB of allocated memory is reached.
4) ```--EnableSaves={true/false} (default false)```
   Enables/disables the creation of a save file.
5) ```--PathToBD={"path/to/file"} (default results/)```
   Path to the directory where the database file will appear.
6) ```--Threads={uint} (default 50)```
   Maximum number of threads that will be used to scan a chunk.
7) ```--ChunkSize={uint} (default 100)```
   Number of addresses contained in 1 chunk.
8) ```--ConnectTimeout={uint} (default 100 ms)```
   Sets the connection waiting time.
9) ```--ReadTimeout={uint} (default 250 ms)```
   Sets the time to wait for the response to be read.
10) ```--ResponseSize={uint} (default 2^30 bytes)```
    Maximum response size.
## More info
1) Gnet-Scanner uses sqlite3 to save the scan results. You can use any DB Browser with sqlite3 support to view/edit the results.
2) Gnet-Scanner uses its own format for the save file, the format is shown below:

| Field name        | Field type | Description                                                  |
|-------------------|:----------:|--------------------------------------------------------------|
| Save Version      |   uint64   | The version of the save file used to determine compatibility |
| Scanned Addresses |   uint64   | Number of addresses already scanned                          |
| Path to DB file   |   string   | Database file path                                           |
| CIDR net          |   string   | The net and the mask                                         |

Note: the string type is a merge of a simple string and its length at the beginning (1 byte). 
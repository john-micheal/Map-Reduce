# Map-Reduce

• MapReduce is a programming model used for efficient processing in parallel over large data-sets in a distributed manner. The data is first split and then combined to produce the final result.
• The purpose of MapReduce in Hadoop is to Map each of the jobs and then it will reduce it to equivalent tasks for providing less overhead over the cluster net


#    how-to-use?

0- Open 7 different terminals on vs code 
1- Run salves from 1 to 5 each on different terminal tabs
2- Run divide_chunks_on_slaves function from master.go file 
3- Comment divide_chunks_on_slaves function and uncomment master as server function then run master.go file again

if client wants chunks uncomment first get request 
4.1- Choose desired chunk number by changing id value from the  get request at client.go file then run client.go
Note : id=0 will write all chunks after gathering them from slaves

if client wants mapReduce aka base count 
4.2- leave current get request

When connecting with other devices only change ip addresses from master file don't change port number only 192.168.1.5 part
"http://192.168.1.5:8090/fasta/baseCount" ---> "http://192.168.14.78:8090/fasta/baseCount"

This Steps for first time running after that you can ignore step 2 because you had already made the chunks and 
divided them over slaves so no need to run it again except if you deleted slaves.fasta files aka chunks

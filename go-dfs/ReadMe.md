# Distributed File Server in golang
A redis middleware based distributed file system using golang

## Steps to install
```
cd go-dfs
go install
```
## Steps to launch
```
cd bin
./go-dfs <master/slave>
```

## Client Usage
```
./go-dfs client saurabh
```

Output:

```
Mode Selected client
File System commands Available...
saurabh-cmd-prompt $:  ls
dir 	--	 test/
file 	--	 b.txt
file 	--	 a.txt

saurabh-cmd-prompt $:  cd test
Moved to dir: /test/

saurabh-cmd-prompt $:  ls
file 	--	 a.xlsx
file 	--	 sau.txt

saurabh-cmd-prompt $:  cd ..
Moved to dir: /

saurabh-cmd-prompt $:  cd ..
Already at root.

saurabh-cmd-prompt $:  exit
Program exiting
```


 #!/bin/sh

go test -covermode=count -coverprofile=count.out

go tool cover -html=count.out

rm count.out

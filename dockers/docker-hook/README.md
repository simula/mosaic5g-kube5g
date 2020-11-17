
[<img align="center" width="40" height="40" src="https://mosaic5g.io/img/m5g-kube5g.png" />](https://gitlab.eurecom.fr/mosaic5g/kube5g/-/tree/develop/dockers/docker-hook)
# docker-hook

## [Install golang](https://golang.org/doc/install):
Note that it is supposed that you create it the mosaic5g project in ```$HOME```. Please change accordingly if you created the project mosaic5g elsewhere.

Note the you can install the golang using the script ```build_m5g``` of mosaic5g, that exists in mosaic5g root directory:
- Download the code
    ```bash
    $ wget https://dl.google.com/go/go1.14.1.linux-amd64.tar.gz
    ```

- Exctract it to ```/usr/local``` (you may need to be sudor):
    ```bash
    $ (sudo) tar -C /usr/local -xzf go1.14.1.linux-amd64.tar.gz
    ```
- Add ```/usr/local/go/bin``` to the PATH environment variable:
    ```bash
    $ echo 'export PATH=$PATH:/usr/local/go/bin' >> $HOME/.profile
    ```
## Test the golang installation
- create the following directory in your home directory:
    ```bash
    $ cd
    $ mkdir -p go/src/hello
    ```
- Export ```GOPATH```
    ```bash
    export GOPATH=$HOME/go/src
    ```
- Create ```hello world``` test file
    ```bash
    cat <<EOF >go/src/hello/hello.go
        package main

        import "fmt"

        func main() {
            fmt.Printf("hello, world\n")
        }
    EOF
    ```
- run the file
    ```bash
    $ cd go/src/hello/
    $ go run hello.go 
    hello, world
    ```
- build and run the file
    ```bash
    $ cd go/src/hello/
    $ go build hello.go
    $ ls
    hello  hello.go
    $ ./hello 
    hello, world
    ```
## Start with ```docker-hook```

Since golang projects should be in ```GOPATH```, we will create sympolic link for the ```docker-hook``` project:
```bash
$ export GOPATH=$HOME/go/src
$ mkdir -p go/src/mosaic5g
$ ln -s $MOSAIC5G_DIR/kube5g/dockers/docker-hook $HOME/go/src/mosaic5g/
```
where ```MOSAIC5G_DIR``` is the root directory of ```kube5g```. For example 
```$MOSAIC5G_DIR``` is  ```$HOME/mosaic5g``` ff you get ```kube5g``` trhough [mosaic5g project](https://gitlab.eurecom.fr/mosaic5g/mosaic5g)

After that, we can work on the project like any other golang projects.

In order to generate the hook that is used to create docker images:
- go to the following directory:
    ```bash
    $ cd $HOME/go/src/mosaic5g/docker-hook/cmd/hook/
    $ go build -o hook main.go
    ```
    Where the option ```-o``` in the last command is to give the name of the binary file, which is here ```hook```.
    After that, you can copy the binary file ```hook``` to ```$HOME/mosaic5g/kube5g/dockers/docker-build/build```

At this stage, you can now start building the containers for mosaic5G snaps. To do so, please check the project that exists in ```$HOME/mosaic5g/kube5g/dockers/docker-build```, which is available [here](https://gitlab.eurecom.fr/mosaic5g/kube5g/-/tree/develop/dockers/docker-build) in [https://gitlab.eurecom.fr](https://gitlab.eurecom.fr)


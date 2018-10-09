# Estimation Test Utility

## Run Locally

### Install Golang
[Official installation guidlines](https://golang.org/dl/)

### Install Node & Yarn
Brew install of yarn will install both yarn & node

```bash
brew install yarn
```

### Get the code
With a standard go installation your go workspace is generally ~/go
```bash
cd $HOME/go
```
Make a directory that mimics the package layout.  We have to do this since we aren't able to ```go get github.com/kdevar/basket-products``` in a sane way for private repos  

```bash
mkdir src/github.com/kdevar
```

```bash
git clone https://github.com/kdevar/basket-products.git 
```

### Build and Run Code
Switch to the right directory
```bash
cd $HOME/go/src/github.com/kdevar/basket-products
```
Run make
```bash
make
```

If everything has gone well, you should find the app at [http://localhost:8080](http://localhost:8080)



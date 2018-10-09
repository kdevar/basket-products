# Estimation Test Utility
A simple utility to investigate whether search time price estimation can serve as an adequate replacement for precomputed estimated values.

The question we are trying to answer is if there are simpler way to provide synthetic prices values when no actual prices are available
## Local Development
You will need to install golang, node, and yarn to be able to build and run the project.

*You will need to VPN'd inorder to reach the elastic search instance.  This might change in the future.*
#### Install Golang
<a href="https://golang.org/dl/" target="_blank">Official installation guidelines</a>

You could probably do through brew as well

#### Install Node & Yarn
Brew install of yarn will install both yarn & node

```bash
brew install yarn
```

If you have node already installed please make sure the version is greater than 6.11.5.  Webpack will need it to be more current

[Learn more about Yarn](https://yarnpkg.com/en/docs/getting-started)

#### Get the code
With a standard go installation your go workspace is generally ~/go
```bash
cd $HOME/go
```
Make a directory that mimics the package layout.  We have to do this since we aren't able to ```go get github.com/kdevar/basket-products``` in a sane way for private repos  

```bash
mkdir src/github.com/kdevar
```

```bash
cd src/github.com/kdevar
```

```bash
git clone https://<your-username>@github.com/kdevar/basket-products.git 
```

#### Build and Run Code
Switch to the right directory
```bash
cd $HOME/go/src/github.com/kdevar/basket-products
```
Run make
```bash
make
```

If everything has gone well, you should find the app at [http://localhost:8080](http://localhost:8080) *(remember to turn on the VPN)*



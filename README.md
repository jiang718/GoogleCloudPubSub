# GoogleCloudPubSub

1: Install GoogleCloud SDK
```
wget https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-186.0.0-linux-x86_64.tar.gz
tar -xvf google-cloud-sdk-186.0.0-linux-x86_64.tar.gz
cd google-cloud-sdk/bin
./install.sh
```

2.Intall fake pubsub server
```
gcloud components install pubsub-emulator
gcloud components update
```

3.Install Go
```
wget https://dl.google.com/go/go1.9.3.linux-amd64.tar.gz
tar -xvf go1.9.3.linux-amd64.tar.gz
mkdir gomaster && mv go gomaster
cd gomaster
export GOROOT=`pwd`/go
export PATH=$PATH:$GOROOT/bin
mkdir gopath
export GOPATH=`pwd`/gopath
go get -u cloud.google.com/go/pubsub
```

4.Intall Project
```
cd gopath
git clone https://github.com/jiang718/GoogleCloudPubSub
cd GoogleCloudPubSub
```

5 Set Environment for GOOGLE CLOUD API:
```
export GOOGLE_APPLICATION_CREDENTIALS="PATH_TO_GIT_REPO/service-account.json"
```

6.Start the server emulator
6.1 if the server are in a remote machine 
```
gcloud beta emulators pubsub start --host-port=<HOST>:<PORT>
```
For Example:
```
gcloud beta emulators pubsub start --host-port=maximus.cs.umn.edu:46839
```
6.2 if the server are in a local machine
```
gcloud beta emulators pubsub start --host-port=localhost:8086
```
or without parameters
```
gcloud beta emulators pubsub start
```
7.Test the program
if the server and the client are on the same machine:
```
gcloud beta emulators pubsub env-init 
```
if the server and the client are on different machines
```
export PUBSUB_EMULATOR_HOST=<host-ip-address> (eg: maximus.cs.umn.edu:46389)
export PUBSUB_PROJECT_ID="simple-pubsub"
go build client.go 
./client.go <CLIENT_NAME>  (a random name such as "Tom" is fine)
```

Note: When run multiple clients, please make sure those clients have different names

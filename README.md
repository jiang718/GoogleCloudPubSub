# GoogleCloudPubSub

----server part--------

1: Install GoogleCloud 
wget https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-186.0.0-linux-x86_64.tar.gz
tar -xvf google-cloud-sdk-186.0.0-linux-x86_64.tar.gz
cd google-cloud-sdk/bin
./install.sh

2.Intall fake pubsub server
gcloud components install pubsub-emulator
gcloud components update

3.Start the environment of emulator

3.1 if the emulator are in a remote machine 

gcloud beta emulators pubsub start --host-port=<HOST>:<PORT>

For Example:
gcloud beta emulators pubsub start --host-port=maximus.cs.umn.edu:46839

3.1 if the emulator are in a local machine
gcloud beta emulators pubsub start --host-port=localhost:8086
or
gcloud beta emulators pubsub start


------client part------
1.Install Go
wget https://dl.google.com/go/go1.9.3.linux-amd64.tar.gztar -xvf go1.9.3.linux-amd64.tar.gz
mkdir gomaster && mv go gomaster
cd gomaster
export GOROOT=`pwd`/go
export PATH=$PATH:$GOROOT/bin
mkdir gopath
export GOPATH=`pwd`/gopath
go get -u cloud.google.com/go/pubsub

2.Intall Project
cd gopath
git clone https://github.com/jiang718/GoogleCloudPubSub
cd GoogleCloudPubSub
go build client.go 
./client.go <CLIENT_NAME>

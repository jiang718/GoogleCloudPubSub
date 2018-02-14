# GoogleCloudPubSub
This is a project simulating pubsub system using google cloud API and golang.<br />

Before follow the instruction, please make sure that you are using bash, because some CSE machines use tcsh. <br/>

1: Install GoogleCloud SDK
```
wget https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-186.0.0-linux-x86_64.tar.gz
tar -xvf google-cloud-sdk-186.0.0-linux-x86_64.tar.gz
cd google-cloud-sdk
./install.sh
```

2.Intall fake pubsub server<br />
Please first open a new terminal inorder to make the bash.rc to make changes. Otherwise "gcloud" command won't be found. <br />
Go to a designed folder that you want to run the server emulator. Make sure it's not inside the Google Cloud SDK folder.
```
gcloud components install pubsub-emulator
gcloud components update
```

3.Install Go
Go to a desired folder (or current folder) that you want to install "go". The project will be installed under this folder too.
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

5.Start the server emulator<br />
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
gcloud beta emulators pubsub start --host-port=localhost:8085
```
or without parameters
```
gcloud beta emulators pubsub start
```
<br />
6.Set the client Environment
Go to the git repo folder "GoogleCloudPubSub"<br />
First build the client using go
```
./buildclient
export GOOGLE_APPLICATION_CREDENTIALS="`pwd`/service-account.json"
```
If the server and the client are on the same machine:
```
gcloud beta emulators pubsub env-init 
```
Execute the command printting from the above command (starts with "export"), and then
```
export PUBSUB_PROJECT_ID="simple-pubsub"
```


if the server and the client are on different machines
```
export PUBSUB_EMULATOR_HOST=<host-ip-address> (eg: maximus.cs.umn.edu:46389)
export PUBSUB_PROJECT_ID="simple-pubsub"
```

Finally,
```
export PUBSUB_PROJECT_ID="simple-pubsub"
```
<br />
Note: When run multiple clients, please make sure those clients have different names

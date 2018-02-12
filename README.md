# GoogleCloudPubSub

----server part--------

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

3. Get the service Account JSON for a specific <Project-id> that you want to use:
3.1 Get the service-account.json 
Follow the instructions on
https://cloud.google.com/pubsub/docs
to be able start the server emulator

3.2 Set Environment:
```
export GOOGLE_APPLICATION_CREDENTIALS="[PATH]"
```

4.Start the server emulator

4.1 if the emulator are in a remote machine 
```
gcloud beta emulators pubsub start --host-port=<HOST>:<PORT>
```

For Example:
```
gcloud beta emulators pubsub start --host-port=maximus.cs.umn.edu:46839
```

4.2 if the emulator are in a local machine
```
gcloud beta emulators pubsub start --host-port=localhost:8086
```
or without parameters
```
gcloud beta emulators pubsub start
```


------client part------
1.Install Go
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

2.Intall Project
```
cd gopath
git clone https://github.com/jiang718/GoogleCloudPubSub
cd GoogleCloudPubSub
go build client.go 
```

3.Set Environment
if the server is on the same machine:
```
gcloud beta emulators pubsub env-init 
```
else
```
export PUBSUB_EMULATOR_HOST=<host-ip-address> (eg: maximus.cs.umn.edu:46389)
export PUBSUB_PROJECT_ID=<project-id>

4.Run Client (After running the server)
```
./client.go <CLIENT_NAME>
```

5.Run multiple clients at a time
Please make sure those clients have different names

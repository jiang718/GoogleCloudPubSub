// Sample pubsub-quickstart creates a Google Cloud Pub/Sub topic. 
package main
import (
	"fmt"
	"log"
    "strings"
    "regexp"
    "os"
    "bufio"
    //"time"
    //"sync"
	// Imports the Google Cloud Pub/Sub client package. 
	"cloud.google.com/go/pubsub"
    "golang.org/x/net/context"
    "google.golang.org/api/iterator"
)

//return true or false and client
func Join(ctx context.Context, projectId string)(*pubsub.Client, bool) {
	//ctx := context.Background() 
	client, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
        log.Fatalf("Failed to create client: %v", err)
        return client, false
	}
    return client, true
}


func CreateTopic(client *pubsub.Client, ctx context.Context, topicName string) (*pubsub.Topic, error) {
	topic, err := client.CreateTopic(ctx, topicName)
	if err != nil {
	    log.Fatalf("Failed to create topic: %v", err)
	}
    return topic, err
}
func ShowAllTopics(client *pubsub.Client, ctx context.Context) {
	var topics []*pubsub.Topic
	it := client.Topics(ctx)
	for {
			topic, err := it.Next()
			if err == iterator.Done {
					break
			}
			if err != nil {
					return
			}
			topics = append(topics, topic)
	}
    for i := 0; i < len(topics); i++ {
        fmt.Println(topics[i])
    }
}

func ExistTopic(client *pubsub.Client, ctx context.Context, topicName string)(*pubsub.Topic, bool) {
    fmt.Println("Let's check all topics")
    ShowAllTopics(client, ctx)
    fmt.Println("Let's check one topic", topicName)
    var topic *pubsub.Topic
	topic = client.Topic(topicName)
    //fmt.Println("The name of currently checking topic is :", (*topic).name)
	ok, err := topic.Exists(ctx)
	if err != nil {
		log.Fatalf("Failed to operating searching topic: %v", err)
        return topic, false
	}
	if !ok {
		log.Fatalf("Topic doesn't exist")
        return topic, false
	}
    return topic, true
}


func Subscribe(client *pubsub.Client, ctx context.Context, subName string, topicName string) (bool) {
    topic, existStatus := ExistTopic(client, ctx, topicName)
    if (existStatus == false) {
        return false
    }
    sub, err := client.CreateSubscription(ctx, subName,
       pubsub.SubscriptionConfig{Topic: topic})
    if err != nil {
	    log.Fatalf("Failed to subscribe: %v for %v", err, sub)
        return false
    }
    return true
}

func Unsubscribe(client *pubsub.Client, ctx context.Context, subName string, topicName string)(bool) {
    sub := client.Subscription(subName)
    if err := sub.Delete(ctx); err != nil {
	    log.Fatalf("Failed to unsubscribe: %v", err)
        return false
    }
    return true
}

func Publish(client *pubsub.Client, ctx context.Context, content string, topicName string) (bool) {
    topic, topicExist := ExistTopic(client, ctx, topicName)
    if (topicExist == false) {
        return false
    }
	defer topic.Stop()
	var results []*pubsub.PublishResult
	r := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(content),
	})
	results = append(results, r)
    return true
}


func PrintHelp() {
    fmt.Println("Please Input Instruction:")
    fmt.Println("-j (Run this to join the server)")
    fmt.Println("-l (Run this to leave the serve)")
    fmt.Println("-r (list all subscriptions)")
    fmt.Println("-c TopicName (Run this to create a topic)")
    fmt.Println("-s TopicName (Run this to subscribe a topic)")
    fmt.Println("-u TopicName (Run this to unsubscribe the server)")
    fmt.Println("-p Content TopicName (Run this to publisth content to a specific topic)")
    fmt.Println("Input \"END\" to terminate this client.")
}

type GlobalStatus struct {
    whetherJoin bool                //whether join to server
    sublist[10000] string           //subscription strings
    subTotal int                    //total subscription
    subExist map[string]bool       //whether this supscription is canceled
    clientName string                //the name of this client
    //chans = [10000]chan int         //the channels to deal with multiple supscription
}


func TalkToServer(client *pubsub.Client, globalStatus *GlobalStatus) {
    //argsWithProg := os.Args
    //os.Args contains the name of program

    //var input string
    var topicName string
    //var subPrefix string = "sub-for-"


	ctx := context.Background()
    in := bufio.NewReader(os.Stdin)
    //var client pubsub.Client  
    client, status := Join(ctx, "simple-pubsub")
    if status == true {
    } else {
        fmt.Println("Something wrong, stop this clients")
    }
    PrintHelp()
    for {
        input,_:= in.ReadString('\n')
        //fmt.Scanln(&input)
        if input == "END\n" {
            break
        }
        r := regexp.MustCompile("[^\\s]+")
        s := r.FindAllString(input, -1)
        if (strings.Contains(input, "-c")) {
            if (globalStatus.whetherJoin == false) {
                fmt.Println("Please first run '-j' manually!")
            } else {
                if len(s) != 2 {
                    fmt.Println("Should attach topic name after '-c'!")
                    fmt.Println("current length of input string is " , len(s))
                } else {
                    topicName = s[1]
                    topic, err := CreateTopic(client, ctx, topicName)
                    if err == nil {
                        fmt.Printf("Topic %v created.\n", topic)
                    }
                }
            }
        } else if (strings.Contains(input, "-s")) {
            if (globalStatus.whetherJoin == false) {
                fmt.Println("Please first run '-j' manually for joining the server!")
            } else {
                if len(s) != 2 {
                    fmt.Println("Should attach topic name after '-s'!")
                } else {
                    topicName = s[1]
                    subName := globalStatus.clientName + "-" + topicName
                    subStatus := Subscribe(client, ctx, subName, topicName)
                    if subStatus == true {
                        fmt.Printf("Successfully subscriped to %v in list %d.\n", topicName, globalStatus.subTotal);
                        globalStatus.sublist[globalStatus.subTotal] = subName
                        globalStatus.subExist[subName] = true
                        globalStatus.subTotal++;
                        go ReceiveSingleForever(client.Subscription(subName), ctx)
				        //PrintSubs(client)
                    }
                }
            }
        } else if (strings.Contains(input, "-u")) {
            if (globalStatus.whetherJoin == false) {
                fmt.Println("Please first run '-j' manually for joining the server!")
            } else {
                if len(s) != 2 {
                    fmt.Println("Should attach topic name after '-u'!")
                } else {
                    topicName = s[1]
                    subName := globalStatus.clientName + "-" + topicName
                    unsubStatus := Unsubscribe(client, ctx, subName, topicName) 
                    if unsubStatus == true {
                        fmt.Printf("Successfully unsubscriped from %v.\n", topicName);
                        globalStatus.subExist[subName] = false
                    }
                }
            }
        } else if (strings.Contains(input, "-p")) {
            if (globalStatus.whetherJoin == false) {
                fmt.Println("Please first run '-j' manually for joining the server!")
            } else {
                if len(s) != 3 {
                    fmt.Println("Should attach Content and TopicName after '-p'")
                } else {
                    content := s[1]
                    topicName = s[2]
                    pubStatus := Publish(client, ctx, content, topicName)
                    if pubStatus == true {
                        fmt.Printf("Successfully published %v into %v\n", content, topicName)
                    } else {
                        fmt.Printf("Published %v fail\n", content, topicName)
                    }
                }
            }
        } else if (strings.Contains(input, "-j")) {
            if (globalStatus.whetherJoin == true) {
                fmt.Println("Already join the server!")
            } else {
                fmt.Println("Successfully join the server!")
            }
            globalStatus.whetherJoin = true
        } else if (strings.Contains(input, "-l")) {
            if (globalStatus.whetherJoin == false) {
                fmt.Println("Already leave the server!")
            } else {
                fmt.Println("Successfully leave the server!")
            }
            globalStatus.whetherJoin = false
        } else if (strings.Contains(input, "-r")) {
            PrintSubs(client, globalStatus)
        }
        PrintHelp()
    }
}

func PrintSubs(client *pubsub.Client, globalStatus *GlobalStatus) {
    fmt.Println("start printing subs")
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, "simple-pubsub")
	if err != nil {
		// TODO: Handle error.
	}
	// List all subscriptions of the project.
	it := client.Subscriptions(ctx)
	for {
		sub, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// TODO: Handle error.
		}
		fmt.Println(sub)
        args := strings.Split(sub.String(), "/")
        fmt.Println("Find sub:", args[3])
        name := strings.Split(args[3], "-")
        if (name[0] == globalStatus.clientName) {
            fmt.Println("Has sub: " + args[3])
        }
	}

}

func ReceiveSingle(sub *pubsub.Subscription, cctx context.Context) {
     err := sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
             fmt.Printf("Got message: %q\n", string(msg.Data))
             msg.Ack()
     })
     if err != nil {
     }
}


func ReceiveSingleForever(sub *pubsub.Subscription, ctx context.Context) {
    for {
        cctx,_ := context.WithCancel(ctx)
        err := sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
                fmt.Printf("Got message: %q\n", string(msg.Data))
                msg.Ack()
        })
        if err != nil {
        }
    }
}


//3. an subscritpion come in, sublist add by 1, insert new routine to sublist
func AddSubscription(client *pubsub.Client, sub *pubsub.Subscription) {
}

//receive message from clients in another thread
func ReceiveMessage(projectId string, globalStatus *GlobalStatus) {
    //1. read subscriptions and build an array of subscriber based on server
    //todo

    //2. start those server routine, inside each routine, they run for ever


    //4. unsubscrip -> delete the subscriber (everytime, check whether sub exist before trying to receive anything)


    //var mu sync.Mutex
    ctx := context.Background()
	client,_ := pubsub.NewClient(ctx, projectId)
    for {
        if ((*globalStatus).whetherJoin == true) {
            //fmt.Println("Into whether join")
	        it := client.Subscriptions(ctx)
            k := 0
            for {
                //time.Sleep(7500000000)
                //fmt.Println("Im receiving")
                k = k +1
                //fmt.Println("Into for ", k)
                sub, err := it.Next()
                if err == iterator.Done {
                    break
                }
                if err != nil {
                }
                //fmt.Println("after it.Next()")
                args := strings.Split(sub.String(), "/")
                //fmt.Println("Find sub:", args[3])
                name := strings.Split(args[3], "-")
                if (name[0] == globalStatus.clientName) {
                    //fmt.Println("Start check sub: " + args[3])
                    //received := 0
                    cctx,_ := context.WithCancel(ctx)
                    go ReceiveSingle(sub, cctx);
                    //err := sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
                    //        fmt.Printf("Got message: %q\n", string(msg.Data))
                    //        msg.Ack()
                    //})
                    //if err != nil {
                    //}
                }
            }
            //if (globalStatus.subTotal == 2) {
            //    fmt.Println("sub Toal ==2")
            //}
            //if (globalStatus.subTotal == 1) {
            //    fmt.Println("sub Toal ==1")
            //}
            //for i := 0; i < globalStatus.subTotal; i++ {
            //    if (globalStatus.subExist[globalStatus.sublist[i]] == true) {
            //        fmt.Println("Start check sub: " + globalStatus.sublist[i])
            //        //received := 0
            //        sub := client.Subscription(globalStatus.sublist[i])
            //        cctx,_ := context.WithCancel(ctx)
            //        err := sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
            //                fmt.Printf("Got message: %q\n", string(msg.Data))
            //                msg.Ack()
            //        })
            //        if err != nil {
            //        }
            //    }
            //}
        }
    }
    fmt.Println("Finish receiving")
}

func main() {
    if (len(os.Args) != 2) {
        fmt.Println("Please attach the client name.")
    } else if (strings.Contains(os.Args[1], "-")) {
        fmt.Println("Please don't put - in the client name.")
    } else {
        client, status := Join(context.Background(), "simple-pubsub")
        //ctx := context.Background()
        globalStatus := GlobalStatus{whetherJoin: status, subTotal: 0}
        globalStatus.whetherJoin = false
        globalStatus.subTotal = 0
        globalStatus.clientName = os.Args[1]
        globalStatus.subExist = make(map[string]bool)
        //go ReceiveMessage("simple-pubsub", &globalStatus)
        TalkToServer(client, &globalStatus)
    }
}

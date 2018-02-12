// Sample pubsub-quickstart creates a Google Cloud Pub/Sub topic. 
package main
import (
	"fmt"
	"log"
    "strings"
    "regexp"
    "os"
    "bufio"
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
    subExist map[string]bool       //whether this subscription is canceled
    clientName string                //the name of this client
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
                    subName := globalStatus.clientName + topicName
                    subStatus := Subscribe(client, ctx, subName, topicName)
                    if subStatus == true {
                        fmt.Printf("Successfully subscriped to %v.\n", topicName);
                        globalStatus.sublist[globalStatus.subTotal] = subName
                        globalStatus.subExist[subName] = true
                        globalStatus.subTotal++;
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
                    subName := globalStatus.clientName + topicName
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
            globalStatus.whetherJoin = true
            fmt.Println("Successfully join the server!")
        } else if (strings.Contains(input, "-l")) {
            globalStatus.whetherJoin = false
            fmt.Println("Successfully leave the server!")
        }
        PrintHelp()
    }
}

//receive message from clients in another thread
func ReceiveMessage(client *pubsub.Client, globalStatus *GlobalStatus) {
    //var mu sync.Mutex
    ctx := context.Background()
    for {
        if ((*globalStatus).whetherJoin == true) {
            for i := 0; i < globalStatus.subTotal; i++ {
                if (globalStatus.subExist[globalStatus.sublist[i]] == true) {
                    fmt.Println("Start check sub: " + globalStatus.sublist[i])
                    //received := 0
                    sub := client.Subscription(globalStatus.sublist[i])
                    cctx,_ := context.WithCancel(ctx)
                    err := sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
                            fmt.Printf("Got message: %q\n", string(msg.Data))
                            msg.Ack()
                    })
                    if err != nil {
                    }
                }
            }
        }
    }
}

func main() {
    if (len(os.Args) == 1) {
        fmt.Println("Please attach the client name.")
    } else {
        client, status := Join(context.Background(), "simple-pubsub")
        //ctx := context.Background()
        globalStatus := GlobalStatus{whetherJoin: status, subTotal: 0}
        globalStatus.whetherJoin = false
        globalStatus.subTotal = 0
        globalStatus.clientName = os.Args[1]
        globalStatus.subExist = make(map[string]bool)
        go ReceiveMessage(client, &globalStatus)
        TalkToServer(client, &globalStatus)
    }
}

package pubsub

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

var(
	PUBLISH="publish"
	SUBSCRIBE="subscribe"
	UNSUBSCRIBE="unsubscribe"
)
type PubSub struct {
	Clients [] Client
	Subcription []Subcription
}
type Client struct {
	Id         string
	Connection *websocket.Conn
}
type Message struct {
	Action string `json:"action"`
	Topic string `json:"topic"`
	Message json.RawMessage `json:"message"`
}
type Subcription struct {
	Topic string
	Client *Client
}

func (ps *PubSub)AddClient(client Client)(*PubSub)  {
	ps.Clients= append(ps.Clients,client)
	//fmt.Println("add new client to the list",client.Id, len(ps.Clients))
	payload:=[]byte("Hello client ID"+client.Id)
	client.Connection.WriteMessage(1,payload)
	return ps
}

func (ps *PubSub)RemoveClient(client Client)(*PubSub)  {
	for index,sub:=range ps.Subcription{
		if sub.Client.Id==client.Id{
			ps.Subcription=append(ps.Subcription[:index],ps.Subcription[index+1:]...)
		}
	}
	for index,c :=range ps.Clients{
		if c.Id==client.Id{
			ps.Clients=append(ps.Clients[:index],ps.Clients[index+1:]...)
		}
	}
	return ps
}

func (ps *PubSub)Subscribe(client *Client,topic string)(*PubSub){
	clientSubs := ps.GetSubscriptions(topic,client)
	if len(clientSubs)>0{
		return ps
	}
		newSubscriber:=Subcription{
		Topic:  topic,
		Client: client,
	}
	ps.Subcription=append(ps.Subcription,newSubscriber)
	return ps
}

func (ps *PubSub)GetSubscriptions(topic string,client *Client)([]Subcription)  {
	var subscriptionList [] Subcription
	for _,subscription:=range ps.Subcription{
		if client!=nil {
			if subscription.Client.Id==client.Id && subscription.Topic==topic{
				subscriptionList=append(subscriptionList,subscription)
			}
		}else{
			if subscription.Topic==topic{
				subscriptionList=append(subscriptionList,subscription)
			}
		}
	}
	return subscriptionList
}

func (ps *PubSub)Publish(topic string,message []byte,excludeClient *Client)  {
	subscriptions :=ps.GetSubscriptions(topic,nil)
	for _,sub:=range subscriptions{
		fmt.Printf("sending to client %s message %s\n",sub.Client.Id,message)
		//sub.Client.Connection.WriteMessage(1,message)
		sub.Client.Send(message)
	}
}

func (client *Client)Send(message []byte)(error){
	return client.Connection.WriteMessage(1,message)
}

func (ps *PubSub)Unsubscribe(client *Client, topic string)  {
	for index,sub:=range ps.Subcription{
		if sub.Client.Id==client.Id&&sub.Topic==topic{
			ps.Subcription=append(ps.Subcription[:index],ps.Subcription[index+1:]...)
		}
	}
}

func (ps *PubSub)HandleReceiveMessage(client Client, messageType int, message []byte)(*PubSub)  {
	m:=Message{}
	err:=json.Unmarshal(message,&m)
	if(err!=nil){
		fmt.Println("this not correct message")
		return ps
	}
	fmt.Println("client message", m.Action,m.Topic,m.Message )
	switch m.Action {
	case PUBLISH:
		ps.Publish(m.Topic,m.Message,nil)
		fmt.Println("message from publisher")
		break
	case SUBSCRIBE:
		ps.Subscribe(&client,m.Topic)
		fmt.Println("new subscriber with topic ",m.Topic, len(ps.Subcription),client.Id)
		break
	case UNSUBSCRIBE:
		fmt.Println("client unsubscribe topic",m.Topic,client.Id)
		ps.Unsubscribe(&client,m.Topic)
		break
	default:
		break
	}
	return ps
}
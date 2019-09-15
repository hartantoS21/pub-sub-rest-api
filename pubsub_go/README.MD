For run code 
open terminal => go to directory file
You can run code with command `go run main.go` or you can build code with command `go build`

For test pub-sub
Open inspect => tab console 

Example test pubsub:
Open chrome => set address http://localhost:1234/
Open inspect => console
type
` ws.send('{"action":"subscribe","topic":"aaaa","message":"wawan”}’) `

Open Firefox => set address http://localhost:1234/
Open inspect => console 
type
`ws.send('{"action":"publish","topic":"aaaa","message":"wawan”}’)`

Additional you can unsubscribe with input 
`ws.send('{"action":"unsubscribe","topic":"aaaa","message":"wawan”}’)`
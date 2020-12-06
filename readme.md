![](https://github.com/newity/crawler/workflows/unit-tests/badge.svg)

Crawler docs: https://godoc.org/github.com/newity/crawler

Blocklib docs: https://godoc.org/github.com/newity/crawler/blocklib

**Crawler** is a framework for parsing Hyperledger Fabric blockchain. You can implement any storage, parser or storage adapter or use default implementations. 

This is how it looks like:

    // get instance of BadgerDB storage
    stor, err := storage.NewBadger(path)
    if err != nil {
    	logrus.Fatal(err)
    }
    
    // create Crawler instance with connection to all Fabric channels specified in connection profile (crawler.WithAutoConnect(USER, ORG)),
    // using BadgerDB storage (crawler.WithStorage(stor)) and with default storage adapter (crawler.WithStorageAdapter(adapter.NewSimpleAdapter(stor)))
    engine, err := crawler.New("connection.yaml", crawler.WithAutoConnect(USER, ORG), crawler.WithStorage(stor), crawler.WithStorageAdapter(adapter.NewSimpleAdapter(stor)))
	if err != nil {
		logrus.Error(err)
	}

    // start block events (*fab.BlockEvent) listeners
	err = engine.Listen(crawler.FromBlock(), crawler.WithBlockNum(0))
	if err != nil {
		logrus.Error(err)
	}

    // start parsing and saving info
	go engine.Run()

And that's all!

Here are the main parts of a crawler:

- **Storage** is responsible for saving data fetched from blockchain. Default is BadgerDB. 

- **Parser** is responsible for processing data from the blockchain. Simply put, this is about how exactly and into what constituent parts we will disassemble the blocks. You can find default implementation in https://github.com/newity/crawler/tree/master/parser/parser.go. Default parser just packs all txs with type ENDORSER_TRANSACTION and all events into [parser.Data](https://github.com/newity/crawler/blob/master/parser/models.go#L13) format. 

- **StorageAdapter** is used for implementation specific logic of saving parsed data into the storage. Default implementation saves gob-serialized parser.Data with block number as the key and retrieves parser.Data by block number specified.

_You can replace any of these components with your own implementation._

**Crawler**  uses [Blocklib](https://godoc.org/github.com/newity/crawler/blocklib) under the hood. You can use it too for your own parsing needs. 

Examples: 

- parse block creator (orderer), tx creator and tx endorsers identities and signatures: https://github.com/newity/crawler/tree/master/examples/getsignatures
- find actual configuration blocks for each block in the network: https://github.com/newity/crawler/tree/master/examples/readlastconfig
- run crawler with BadgerDB storage and retrieve specific block from storage: https://github.com/newity/crawler/tree/master/examples/saveblocksandread
- Google Pub/Sub usage as Crawler storage: https://github.com/newity/crawler/tree/master/examples/pubsub
- NATS Streaming usage as Crawler storage: https://github.com/newity/crawler/tree/master/examples/nats
- Blocklib basic usage example: https://github.com/newity/crawler/tree/master/examples/blocklib
- Blocklib usage for getting chaincode events: https://github.com/newity/crawler/tree/master/examples/blocklib-events
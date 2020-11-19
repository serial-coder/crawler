**Crawler with NATS Streaming storage**

Start NATS Streaming server:

    docker run --name nats -p 4223:4222 -d nats-streaming:latest --cluster_id testcluster --clustered --store file --dir datastore --cluster_bootstrap
    
Run example:

    go run main.go
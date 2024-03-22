# RaftHouse

On this episode of "how much can I mangle Clickhouse?", we add the raft consensus protocol in front of Clickhouse.

"But Dan, Clickhouse already has replication with (Zoo)Keeper, why are you doing this?" I hear you ask.

The problem with the existing replication method is that it's eventually consistent. [Their strategies for guaranteeing consistency](https://clickhouse.com/docs/knowledgebase/read_consistency#talking-to-a-random-node) (read your writes) are pretty disappointing.

There are some cases where we need the power of Clickhouse (fast columnar reads, extensive query language, materialized view framework), high availability, and the ability to have consistent reads. In this case, we can stick raft in front of the HTTP interface to manage leader election and replication.

"But you lose the ability to have shards!" I hear you yell.

"Clickhouse shards", yes, but not Clickhouse shards. Multi-group raft serves as the PERFECT tool to introduce sharding at the raft level. Because Clickhouse already makes re-sharding a PITA, we can leverage this fault with using simple hashing for choosing a raft group. I've not figured out how to do that elegantly yet, but it will work once I do! (maybe a connection pool for each hash token).  

## Known Limitations

1. Cannot do non-deterministic operations like rand and time based stuff for inserts, as results will be different on each node. Do all time and randomness in your code!
2. Must use the HTTP interface for consistent operations (read and write). Can use the native direct protocol but must NEVER WRITE.
3. Not really a limitation, but the read endpoints should use a read-only user, so you never accidentally insert to it and cause permanent inconsistencies
4. (CURRENTLY) Can only run a single shard, limited by single node size (multi-group in will solve this)
5. Should never use server side async inserts, as this can cause non-deterministic behavior
6. Should have 2 DNS records pointing to RaftHouse (configured with env vars):
   1. One that always goes to the leader (writes and consistent reads)
   2. One that goes to a random node (eventually consistent reads)

## Read and Write Endpoints

RaftHouse should be exposed over two endpoints: A read-only eventually consistent endpoint (queries local node), and a leader endpoint (inserts and consistent queries). This generally means making two DNS records, kubernetes services, etc. that will have distinct `host` HTTP headers when accessed by the Clickhouse client over HTTP (eventually consistent operations can actually use the native client, but should NEVER write).

Any time a write occurs, or a setting is changed, it must go over the leader endpoint (read-write).
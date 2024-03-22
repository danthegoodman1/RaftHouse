# RaftHouse

On this episode of "how much can I mangle ClickHouse?", we add the raft consensus protocol in front of Clickhouse.

_"But Dan, ClickHouse already has replication with (Zoo)Keeper, why are you doing this?"_ I hear you ask.

The problem with the existing replication method is that it's eventually consistent. [Their strategies for guaranteeing consistency](https://clickhouse.com/docs/knowledgebase/read_consistency#talking-to-a-random-node) (read your writes) are pretty disappointing.

There are some cases where we need the power of ClickHouse (fast columnar reads, extensive query language, materialized view framework), high availability, and the ability to have consistent reads. In this case, we can stick raft in front of the HTTP interface to manage leader election and replication.

"But you lose the ability to have shards!" I hear you yell.

"Clickhouse shards", yes, but not ClickHouse shards. Multi-group raft serves as the PERFECT tool to introduce sharding at the raft level. Because ClickHouse already makes re-sharding a PITA, we can leverage this fault with using simple hashing for choosing a raft group. I've not figured out how to do that elegantly yet, but it will work once I do! (maybe a connection pool for each hash token).

RaftHouse is run on the same node as ClickHouse, similar to how you might run Keeper on the same nodes.

_"Can't you just run queries against a single node?"_

Yes, but now you're relying on that single node being up. If it goes down, how do you select another node that has the most upto date data? How do you handle waiting for the original node to come back up? Raft handles this all for you!

_"Can you just run Raft on each node, and then run your processes as observers to find the leader?"_

Insightful question! Not only does this require you integrating raft into your code, or hitting some API to ask for raft, but there are still situations where it could change after the fact. Frankly, this relies too much on the implementation/developer and is a recipe for disaster.

_"What about the `insert_quorum` setting?"_

This requires all nodes to be online. Not acceptable because it's not required for consistency. What's the point of replicas if you can't insert when they are all alive? In the use cases we're building, we're only selecting during recovery, but inserting all the time, which is the opposite pattern this setting is optimized for.

_"What about the `select_sequential_consistency` setting?"

Astute observation! Yes this works in many cases, but puts a lot of load on Keeper, which is probably not as performant as RaftHouse.

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

RaftHouse should be exposed over two endpoints: A read-only eventually consistent endpoint (queries local node), and a leader endpoint (inserts and consistent queries). This generally means making two DNS records, kubernetes services, etc. that will have distinct `host` HTTP headers when accessed by the ClickHouse client over HTTP (eventually consistent operations can actually use the native client, but should NEVER write).

Any time a write occurs, or a setting is changed, it must go over the leader endpoint (read-write).

You don't actually have to define a read-only endpoint. Any `host` value that does not match the `LEADER_HOST` env var will query the local instance of ClickHouse.

## Optimizations

Multi-group raft for shards.

Using async writes internally, then flushing them with the `Sync` method on the state machine. This would leverage the internal batching as well.
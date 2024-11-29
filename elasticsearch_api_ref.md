# REST API Reference - Open Distro Documentation
**The Open Distro project is archived.** Open Distro development has moved to OpenSearch. The Open Distro plugins will continue to work with legacy versions of Elasticsearch OSS, but we recommend upgrading to OpenSearch to take advantage of the latest features and improvements.

Elasticsearch REST API reference
--------------------------------

This reference originates from the Elasticsearch REST API specification. We’re extremely grateful to the Elasticsearch community for their numerous contributions to open source software, including this documentation.

* * *

bulk
----

Perform multiple index, update, and/or delete operations in a single request.

```
POST {index}/_bulk
PUT {index}/_bulk

```


```
POST {index}/{type}/_bulk
PUT {index}/{type}/_bulk

```


#### HTTP request body

The operation definition and data (action-data pairs), separated by newlines.

**Required**: True

#### URL parameters



* Parameter: wait_for_active_shards
  * Type: string
  * Description: Sets the number of shard copies that must be active before proceeding with the bulk operation. Defaults to 1, meaning the primary shard only. Set to all for all shard copies. Otherwise, set to any non-negative value less than or equal to the total number of copies for the shard (number of replicas + 1).
* Parameter: refresh
  * Type: enum
  * Description: If true, refresh the affected shards to make this operation visible to search. If wait_for, wait for a refresh to make this operation visible to search. If false (the default), do nothing with refreshes.
* Parameter: routing
  * Type: string
  * Description: Specific routing value.
* Parameter: timeout
  * Type: time
  * Description: Explicit operation timeout.
* Parameter: type
  * Type: string
  * Description: Default document type for items that don’t provide one.
* Parameter: _source
  * Type: list
  * Description: True or false to return the _source field or not, or the default list of fields to return, can be overridden on each sub-request.
* Parameter: _source_excludes
  * Type: list
  * Description: Default list of fields to exclude from the returned _source field, can be overridden on each sub-request.
* Parameter: _source_includes
  * Type: list
  * Description: Default list of fields to extract and return from the _source field, can be overridden on each sub-request.
* Parameter: pipeline
  * Type: string
  * Description: The pipeline ID to preprocess incoming documents with.
* Parameter: require_alias
  * Type: boolean
  * Description: Sets require_alias for all incoming documents, defaults to false (unset).


cat.aliases
-----------

Shows information about currently configured aliases to indices including filter and routing infos.

#### URL parameters



* Parameter: format
  * Type: string
  * Description: a short version of the Accept header, e.g. json, yaml
* Parameter: local
  * Type: boolean
  * Description: Return local information, do not retrieve the state from master node (default: false)
* Parameter: h
  * Type: list
  * Description: Comma-separated list of column names to display
* Parameter: help
  * Type: boolean
  * Description: Return help information
* Parameter: s
  * Type: list
  * Description: Comma-separated list of column names or column aliases to sort by
* Parameter: v
  * Type: boolean
  * Description: Verbose mode. Display column headers
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.


cat.allocation
--------------

Provides a snapshot of how many shards are allocated to each data node and how much disk space they are using.

```
GET _cat/allocation/{node_id}

```


#### URL parameters



* Parameter: format
  * Type: string
  * Description: a short version of the Accept header, e.g. json, yaml
* Parameter: bytes
  * Type: enum
  * Description: The unit in which to display byte values
* Parameter: local
  * Type: boolean
  * Description: Return local information, do not retrieve the state from master node (default: false)
* Parameter: master_timeout
  * Type: time
  * Description: Explicit operation timeout for connection to master node
* Parameter: h
  * Type: list
  * Description: Comma-separated list of column names to display
* Parameter: help
  * Type: boolean
  * Description: Return help information
* Parameter: s
  * Type: list
  * Description: Comma-separated list of column names or column aliases to sort by
* Parameter: v
  * Type: boolean
  * Description: Verbose mode. Display column headers


cat.count
---------

Provides quick access to the document count of the entire cluster, or individual indices.

#### URL parameters


|Parameter|Type   |Description                                                      |
|---------|-------|-----------------------------------------------------------------|
|format   |string |a short version of the Accept header, e.g. json, yaml            |
|h        |list   |Comma-separated list of column names to display                  |
|help     |boolean|Return help information                                          |
|s        |list   |Comma-separated list of column names or column aliases to sort by|
|v        |boolean|Verbose mode. Display column headers                             |


cat.fielddata
-------------

Shows how much heap memory is currently being used by fielddata on every data node in the cluster.

```
GET _cat/fielddata/{fields}

```


#### URL parameters


|Parameter|Type   |Description                                                      |
|---------|-------|-----------------------------------------------------------------|
|format   |string |a short version of the Accept header, e.g. json, yaml            |
|bytes    |enum   |The unit in which to display byte values                         |
|h        |list   |Comma-separated list of column names to display                  |
|help     |boolean|Return help information                                          |
|s        |list   |Comma-separated list of column names or column aliases to sort by|
|v        |boolean|Verbose mode. Display column headers                             |
|fields   |list   |A comma-separated list of fields to return in the output         |


cat.health
----------

Returns a concise representation of the cluster health.

#### URL parameters


|Parameter|Type   |Description                                                      |
|---------|-------|-----------------------------------------------------------------|
|format   |string |a short version of the Accept header, e.g. json, yaml            |
|h        |list   |Comma-separated list of column names to display                  |
|help     |boolean|Return help information                                          |
|s        |list   |Comma-separated list of column names or column aliases to sort by|
|time     |enum   |The unit in which to display time values                         |
|ts       |boolean|Set to false to disable timestamping                             |
|v        |boolean|Verbose mode. Display column headers                             |


cat.help
--------

Returns help for the Cat APIs.

#### URL parameters


|Parameter|Type   |Description                                                      |
|---------|-------|-----------------------------------------------------------------|
|help     |boolean|Return help information                                          |
|s        |list   |Comma-separated list of column names or column aliases to sort by|


cat.indices
-----------

Returns information about indices: number of primaries and replicas, document counts, disk size, …

#### URL parameters



* Parameter: format
  * Type: string
  * Description: a short version of the Accept header, e.g. json, yaml
* Parameter: bytes
  * Type: enum
  * Description: The unit in which to display byte values
* Parameter: local
  * Type: boolean
  * Description: Return local information, do not retrieve the state from master node (default: false)
* Parameter: master_timeout
  * Type: time
  * Description: Explicit operation timeout for connection to master node
* Parameter: h
  * Type: list
  * Description: Comma-separated list of column names to display
* Parameter: health
  * Type: enum
  * Description: A health status (“green”, “yellow”, or “red” to filter only indices matching the specified health status
* Parameter: help
  * Type: boolean
  * Description: Return help information
* Parameter: pri
  * Type: boolean
  * Description: Set to true to return stats only for primary shards
* Parameter: s
  * Type: list
  * Description: Comma-separated list of column names or column aliases to sort by
* Parameter: time
  * Type: enum
  * Description: The unit in which to display time values
* Parameter: v
  * Type: boolean
  * Description: Verbose mode. Display column headers
* Parameter: include_unloaded_segments
  * Type: boolean
  * Description: If set to true segment stats will include stats for segments that are not currently loaded into memory
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.


cat.master
----------

Returns information about the master node.

#### URL parameters



* Parameter: format
  * Type: string
  * Description: a short version of the Accept header, e.g. json, yaml
* Parameter: local
  * Type: boolean
  * Description: Return local information, do not retrieve the state from master node (default: false)
* Parameter: master_timeout
  * Type: time
  * Description: Explicit operation timeout for connection to master node
* Parameter: h
  * Type: list
  * Description: Comma-separated list of column names to display
* Parameter: help
  * Type: boolean
  * Description: Return help information
* Parameter: s
  * Type: list
  * Description: Comma-separated list of column names or column aliases to sort by
* Parameter: v
  * Type: boolean
  * Description: Verbose mode. Display column headers


cat.nodeattrs
-------------

Returns information about custom node attributes.

#### URL parameters



* Parameter: format
  * Type: string
  * Description: a short version of the Accept header, e.g. json, yaml
* Parameter: local
  * Type: boolean
  * Description: Return local information, do not retrieve the state from master node (default: false)
* Parameter: master_timeout
  * Type: time
  * Description: Explicit operation timeout for connection to master node
* Parameter: h
  * Type: list
  * Description: Comma-separated list of column names to display
* Parameter: help
  * Type: boolean
  * Description: Return help information
* Parameter: s
  * Type: list
  * Description: Comma-separated list of column names or column aliases to sort by
* Parameter: v
  * Type: boolean
  * Description: Verbose mode. Display column headers


cat.nodes
---------

Returns basic statistics about performance of cluster nodes.

#### URL parameters



* Parameter: bytes
  * Type: enum
  * Description: The unit in which to display byte values
* Parameter: format
  * Type: string
  * Description: a short version of the Accept header, e.g. json, yaml
* Parameter: full_id
  * Type: boolean
  * Description: Return the full node ID instead of the shortened version (default: false)
* Parameter: local
  * Type: boolean
  * Description: Calculate the selected nodes using the local cluster state rather than the state from master node (default: false)
* Parameter: master_timeout
  * Type: time
  * Description: Explicit operation timeout for connection to master node
* Parameter: h
  * Type: list
  * Description: Comma-separated list of column names to display
* Parameter: help
  * Type: boolean
  * Description: Return help information
* Parameter: s
  * Type: list
  * Description: Comma-separated list of column names or column aliases to sort by
* Parameter: time
  * Type: enum
  * Description: The unit in which to display time values
* Parameter: v
  * Type: boolean
  * Description: Verbose mode. Display column headers


cat.pending\_tasks
------------------

Returns a concise representation of the cluster pending tasks.

#### URL parameters



* Parameter: format
  * Type: string
  * Description: a short version of the Accept header, e.g. json, yaml
* Parameter: local
  * Type: boolean
  * Description: Return local information, do not retrieve the state from master node (default: false)
* Parameter: master_timeout
  * Type: time
  * Description: Explicit operation timeout for connection to master node
* Parameter: h
  * Type: list
  * Description: Comma-separated list of column names to display
* Parameter: help
  * Type: boolean
  * Description: Return help information
* Parameter: s
  * Type: list
  * Description: Comma-separated list of column names or column aliases to sort by
* Parameter: time
  * Type: enum
  * Description: The unit in which to display time values
* Parameter: v
  * Type: boolean
  * Description: Verbose mode. Display column headers


cat.plugins
-----------

Returns information about installed plugins across nodes node.

#### URL parameters



* Parameter: format
  * Type: string
  * Description: a short version of the Accept header, e.g. json, yaml
* Parameter: local
  * Type: boolean
  * Description: Return local information, do not retrieve the state from master node (default: false)
* Parameter: master_timeout
  * Type: time
  * Description: Explicit operation timeout for connection to master node
* Parameter: h
  * Type: list
  * Description: Comma-separated list of column names to display
* Parameter: help
  * Type: boolean
  * Description: Return help information
* Parameter: s
  * Type: list
  * Description: Comma-separated list of column names or column aliases to sort by
* Parameter: v
  * Type: boolean
  * Description: Verbose mode. Display column headers


cat.recovery
------------

Returns information about index shard recoveries, both on-going completed.

```
GET _cat/recovery/{index}

```


#### URL parameters



* Parameter: format
  * Type: string
  * Description: a short version of the Accept header, e.g. json, yaml
* Parameter: active_only
  * Type: boolean
  * Description: If true, the response only includes ongoing shard recoveries
* Parameter: bytes
  * Type: enum
  * Description: The unit in which to display byte values
* Parameter: detailed
  * Type: boolean
  * Description: If true, the response includes detailed information about shard recoveries
* Parameter: h
  * Type: list
  * Description: Comma-separated list of column names to display
* Parameter: help
  * Type: boolean
  * Description: Return help information
* Parameter: index
  * Type: list
  * Description: Comma-separated list or wildcard expression of index names to limit the returned information
* Parameter: s
  * Type: list
  * Description: Comma-separated list of column names or column aliases to sort by
* Parameter: time
  * Type: enum
  * Description: The unit in which to display time values
* Parameter: v
  * Type: boolean
  * Description: Verbose mode. Display column headers


cat.repositories
----------------

Returns information about snapshot repositories registered in the cluster.

#### URL parameters


|Parameter     |Type   |Description                                                         |
|--------------|-------|--------------------------------------------------------------------|
|format        |string |a short version of the Accept header, e.g. json, yaml               |
|local         |boolean|Return local information, do not retrieve the state from master node|
|master_timeout|time   |Explicit operation timeout for connection to master node            |
|h             |list   |Comma-separated list of column names to display                     |
|help          |boolean|Return help information                                             |
|s             |list   |Comma-separated list of column names or column aliases to sort by   |
|v             |boolean|Verbose mode. Display column headers                                |


cat.segments
------------

Provides low-level information about the segments in the shards of an index.

```
GET _cat/segments/{index}

```


#### URL parameters


|Parameter|Type   |Description                                                      |
|---------|-------|-----------------------------------------------------------------|
|format   |string |a short version of the Accept header, e.g. json, yaml            |
|bytes    |enum   |The unit in which to display byte values                         |
|h        |list   |Comma-separated list of column names to display                  |
|help     |boolean|Return help information                                          |
|s        |list   |Comma-separated list of column names or column aliases to sort by|
|v        |boolean|Verbose mode. Display column headers                             |


cat.shards
----------

Provides a detailed view of shard allocation on nodes.

#### URL parameters



* Parameter: format
  * Type: string
  * Description: a short version of the Accept header, e.g. json, yaml
* Parameter: bytes
  * Type: enum
  * Description: The unit in which to display byte values
* Parameter: local
  * Type: boolean
  * Description: Return local information, do not retrieve the state from master node (default: false)
* Parameter: master_timeout
  * Type: time
  * Description: Explicit operation timeout for connection to master node
* Parameter: h
  * Type: list
  * Description: Comma-separated list of column names to display
* Parameter: help
  * Type: boolean
  * Description: Return help information
* Parameter: s
  * Type: list
  * Description: Comma-separated list of column names or column aliases to sort by
* Parameter: time
  * Type: enum
  * Description: The unit in which to display time values
* Parameter: v
  * Type: boolean
  * Description: Verbose mode. Display column headers


cat.snapshots
-------------

Returns all snapshots in a specific repository.

```
GET _cat/snapshots/{repository}

```


#### URL parameters


|Parameter         |Type   |Description                                                      |
|------------------|-------|-----------------------------------------------------------------|
|format            |string |a short version of the Accept header, e.g. json, yaml            |
|ignore_unavailable|boolean|Set to true to ignore unavailable snapshots                      |
|master_timeout    |time   |Explicit operation timeout for connection to master node         |
|h                 |list   |Comma-separated list of column names to display                  |
|help              |boolean|Return help information                                          |
|s                 |list   |Comma-separated list of column names or column aliases to sort by|
|time              |enum   |The unit in which to display time values                         |
|v                 |boolean|Verbose mode. Display column headers                             |


cat.tasks
---------

Returns information about the tasks currently executing on one or more nodes in the cluster.

#### URL parameters



* Parameter: format
  * Type: string
  * Description: a short version of the Accept header, e.g. json, yaml
* Parameter: nodes
  * Type: list
  * Description: A comma-separated list of node IDs or names to limit the returned information; use _local to return information from the node you’re connecting to, leave empty to get information from all nodes
* Parameter: actions
  * Type: list
  * Description: A comma-separated list of actions that should be returned. Leave empty to return all.
* Parameter: detailed
  * Type: boolean
  * Description: Return detailed task information (default: false)
* Parameter: parent_task_id
  * Type: string
  * Description: Return tasks with specified parent task id (node_id:task_number). Set to -1 to return all.
* Parameter: h
  * Type: list
  * Description: Comma-separated list of column names to display
* Parameter: help
  * Type: boolean
  * Description: Return help information
* Parameter: s
  * Type: list
  * Description: Comma-separated list of column names or column aliases to sort by
* Parameter: time
  * Type: enum
  * Description: The unit in which to display time values
* Parameter: v
  * Type: boolean
  * Description: Verbose mode. Display column headers


cat.templates
-------------

Returns information about existing templates.

```
GET _cat/templates/{name}

```


#### URL parameters



* Parameter: format
  * Type: string
  * Description: a short version of the Accept header, e.g. json, yaml
* Parameter: local
  * Type: boolean
  * Description: Return local information, do not retrieve the state from master node (default: false)
* Parameter: master_timeout
  * Type: time
  * Description: Explicit operation timeout for connection to master node
* Parameter: h
  * Type: list
  * Description: Comma-separated list of column names to display
* Parameter: help
  * Type: boolean
  * Description: Return help information
* Parameter: s
  * Type: list
  * Description: Comma-separated list of column names or column aliases to sort by
* Parameter: v
  * Type: boolean
  * Description: Verbose mode. Display column headers


cat.thread\_pool
----------------

Returns cluster-wide thread pool statistics per node. By default the active, queue and rejected statistics are returned for all thread pools.

```
GET _cat/thread_pool/{thread_pool_patterns}

```


#### URL parameters



* Parameter: format
  * Type: string
  * Description: a short version of the Accept header, e.g. json, yaml
* Parameter: size
  * Type: enum
  * Description: The multiplier in which to display values
* Parameter: local
  * Type: boolean
  * Description: Return local information, do not retrieve the state from master node (default: false)
* Parameter: master_timeout
  * Type: time
  * Description: Explicit operation timeout for connection to master node
* Parameter: h
  * Type: list
  * Description: Comma-separated list of column names to display
* Parameter: help
  * Type: boolean
  * Description: Return help information
* Parameter: s
  * Type: list
  * Description: Comma-separated list of column names or column aliases to sort by
* Parameter: v
  * Type: boolean
  * Description: Verbose mode. Display column headers


Explicitly clears the search context for a scroll.

```
DELETE _search/scroll/{scroll_id}

```


#### HTTP request body

A comma-separated list of scroll IDs to clear if none was specified via the scroll\_id parameter

cluster.allocation\_explain
---------------------------

Provides explanations for shard allocations in the cluster.

```
GET _cluster/allocation/explain
POST _cluster/allocation/explain

```


#### HTTP request body

The index, shard, and primary flag to explain. Empty means ‘explain the first unassigned shard’

#### URL parameters



* Parameter: include_yes_decisions
  * Type: boolean
  * Description: Return ‘YES’ decisions in explanation (default: false)
* Parameter: include_disk_info
  * Type: boolean
  * Description: Return information about disk usage and shard sizes (default: false)


cluster.delete\_component\_template
-----------------------------------

Deletes a component template

```
DELETE _component_template/{name}

```


#### URL parameters


|Parameter     |Type|Description                             |
|--------------|----|----------------------------------------|
|timeout       |time|Explicit operation timeout              |
|master_timeout|time|Specify timeout for connection to master|


cluster.delete\_voting\_config\_exclusions
------------------------------------------

Clears cluster voting config exclusions.

```
DELETE _cluster/voting_config_exclusions

```


#### URL parameters



* Parameter: wait_for_removal
  * Type: boolean
  * Description: Specifies whether to wait for all excluded nodes to be removed from the cluster before clearing the voting configuration exclusions list.


cluster.exists\_component\_template
-----------------------------------

Returns information about whether a particular component template exist

```
HEAD _component_template/{name}

```


#### URL parameters



* Parameter: master_timeout
  * Type: time
  * Description: Explicit operation timeout for connection to master node
* Parameter: local
  * Type: boolean
  * Description: Return local information, do not retrieve the state from master node (default: false)


cluster.get\_component\_template
--------------------------------

Returns one or more component templates

```
GET _component_template/{name}

```


#### URL parameters



* Parameter: master_timeout
  * Type: time
  * Description: Explicit operation timeout for connection to master node
* Parameter: local
  * Type: boolean
  * Description: Return local information, do not retrieve the state from master node (default: false)


cluster.get\_settings
---------------------

Returns cluster settings.

#### URL parameters


|Parameter       |Type   |Description                                             |
|----------------|-------|--------------------------------------------------------|
|flat_settings   |boolean|Return settings in flat format (default: false)         |
|master_timeout  |time   |Explicit operation timeout for connection to master node|
|timeout         |time   |Explicit operation timeout                              |
|include_defaults|boolean|Whether to return all default clusters setting.         |


cluster.health
--------------

Returns basic information about the health of the cluster.

```
GET _cluster/health/{index}

```


#### URL parameters



* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.
* Parameter: level
  * Type: enum
  * Description: Specify the level of detail for returned information
* Parameter: local
  * Type: boolean
  * Description: Return local information, do not retrieve the state from master node (default: false)
* Parameter: master_timeout
  * Type: time
  * Description: Explicit operation timeout for connection to master node
* Parameter: timeout
  * Type: time
  * Description: Explicit operation timeout
* Parameter: wait_for_active_shards
  * Type: string
  * Description: Wait until the specified number of shards is active
* Parameter: wait_for_nodes
  * Type: string
  * Description: Wait until the specified number of nodes is available
* Parameter: wait_for_events
  * Type: enum
  * Description: Wait until all currently queued events with the given priority are processed
* Parameter: wait_for_no_relocating_shards
  * Type: boolean
  * Description: Whether to wait until there are no relocating shards in the cluster
* Parameter: wait_for_no_initializing_shards
  * Type: boolean
  * Description: Whether to wait until there are no initializing shards in the cluster
* Parameter: wait_for_status
  * Type: enum
  * Description: Wait until cluster is in a specific state


cluster.pending\_tasks
----------------------

Returns a list of any cluster-level changes (e.g. create index, update mapping, allocate or fail shard) which have not yet been executed.

```
GET _cluster/pending_tasks

```


#### URL parameters



* Parameter: local
  * Type: boolean
  * Description: Return local information, do not retrieve the state from master node (default: false)
* Parameter: master_timeout
  * Type: time
  * Description: Specify timeout for connection to master


cluster.post\_voting\_config\_exclusions
----------------------------------------

Updates the cluster voting config exclusions by node ids or node names.

```
POST _cluster/voting_config_exclusions

```


#### URL parameters



* Parameter: node_ids
  * Type: string
  * Description: A comma-separated list of the persistent ids of the nodes to exclude from the voting configuration. If specified, you may not also specify ?node_names.
* Parameter: node_names
  * Type: string
  * Description: A comma-separated list of the names of the nodes to exclude from the voting configuration. If specified, you may not also specify ?node_ids.
* Parameter: timeout
  * Type: time
  * Description: Explicit operation timeout


cluster.put\_component\_template
--------------------------------

Creates or updates a component template

```
PUT _component_template/{name}
POST _component_template/{name}

```


#### HTTP request body

The template definition

**Required**: True

#### URL parameters



* Parameter: create
  * Type: boolean
  * Description: Whether the index template should only be added if new or can also replace an existing one
* Parameter: timeout
  * Type: time
  * Description: Explicit operation timeout
* Parameter: master_timeout
  * Type: time
  * Description: Specify timeout for connection to master


cluster.put\_settings
---------------------

Updates the cluster settings.

#### HTTP request body

The settings to be updated. Can be either `transient` or `persistent` (survives cluster restart).

**Required**: True

#### URL parameters


|Parameter     |Type   |Description                                             |
|--------------|-------|--------------------------------------------------------|
|flat_settings |boolean|Return settings in flat format (default: false)         |
|master_timeout|time   |Explicit operation timeout for connection to master node|
|timeout       |time   |Explicit operation timeout                              |


cluster.remote\_info
--------------------

Returns the information about configured remote clusters.

cluster.reroute
---------------

Allows to manually change the allocation of individual shards in the cluster.

#### HTTP request body

The definition of `commands` to perform (`move`, `cancel`, `allocate`)

#### URL parameters



* Parameter: dry_run
  * Type: boolean
  * Description: Simulate the operation only and return the resulting state
* Parameter: explain
  * Type: boolean
  * Description: Return an explanation of why the commands can or cannot be executed
* Parameter: retry_failed
  * Type: boolean
  * Description: Retries allocation of shards that are blocked due to too many subsequent allocation failures
* Parameter: metric
  * Type: list
  * Description: Limit the information returned to the specified metrics. Defaults to all but metadata
* Parameter: master_timeout
  * Type: time
  * Description: Explicit operation timeout for connection to master node
* Parameter: timeout
  * Type: time
  * Description: Explicit operation timeout


cluster.state
-------------

Returns a comprehensive information about the state of the cluster.

```
GET _cluster/state/{metric}

```


```
GET _cluster/state/{metric}/{index}

```


#### URL parameters



* Parameter: local
  * Type: boolean
  * Description: Return local information, do not retrieve the state from master node (default: false)
* Parameter: master_timeout
  * Type: time
  * Description: Specify timeout for connection to master
* Parameter: flat_settings
  * Type: boolean
  * Description: Return settings in flat format (default: false)
* Parameter: wait_for_metadata_version
  * Type: number
  * Description: Wait for the metadata version to be equal or greater than the specified metadata version
* Parameter: wait_for_timeout
  * Type: time
  * Description: The maximum time to wait for wait_for_metadata_version before timing out
* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether specified concrete indices should be ignored when unavailable (missing or closed)
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes _all string or when no indices have been specified)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.


cluster.stats
-------------

Returns high-level overview of cluster statistics.

```
GET _cluster/stats/nodes/{node_id}

```


#### URL parameters


|Parameter    |Type   |Description                                    |
|-------------|-------|-----------------------------------------------|
|flat_settings|boolean|Return settings in flat format (default: false)|
|timeout      |time   |Explicit operation timeout                     |


count
-----

Returns number of documents matching a query.

```
POST {index}/_count
GET {index}/_count

```


```
POST {index}/{type}/_count
GET {index}/{type}/_count

```


#### HTTP request body

A query to restrict the results specified with the Query DSL (optional)

#### URL parameters



* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether specified concrete indices should be ignored when unavailable (missing or closed)
* Parameter: ignore_throttled
  * Type: boolean
  * Description: Whether specified concrete, expanded or aliased indices should be ignored when throttled
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes _all string or when no indices have been specified)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.
* Parameter: min_score
  * Type: number
  * Description: Include only documents with a specific _score value in the result
* Parameter: preference
  * Type: string
  * Description: Specify the node or shard the operation should be performed on (default: random)
* Parameter: routing
  * Type: list
  * Description: A comma-separated list of specific routing values
* Parameter: q
  * Type: string
  * Description: Query in the Lucene query string syntax
* Parameter: analyzer
  * Type: string
  * Description: The analyzer to use for the query string
* Parameter: analyze_wildcard
  * Type: boolean
  * Description: Specify whether wildcard and prefix queries should be analyzed (default: false)
* Parameter: default_operator
  * Type: enum
  * Description: The default operator for query string query (AND or OR)
* Parameter: df
  * Type: string
  * Description: The field to use as default where no field prefix is given in the query string
* Parameter: lenient
  * Type: boolean
  * Description: Specify whether format-based query failures (such as providing text to a numeric field) should be ignored
* Parameter: terminate_after
  * Type: number
  * Description: The maximum count for each shard, upon reaching which the query execution will terminate early


create
------

Creates a new document in the index.

Returns a 409 response when a document with a same ID already exists in the index.

```
PUT {index}/_create/{id}
POST {index}/_create/{id}

```


```
PUT {index}/{type}/{id}/_create
POST {index}/{type}/{id}/_create

```


#### HTTP request body

The document

**Required**: True

#### URL parameters



* Parameter: wait_for_active_shards
  * Type: string
  * Description: Sets the number of shard copies that must be active before proceeding with the index operation. Defaults to 1, meaning the primary shard only. Set to all for all shard copies, otherwise set to any non-negative value less than or equal to the total number of copies for the shard (number of replicas + 1)
* Parameter: refresh
  * Type: enum
  * Description: If true then refresh the affected shards to make this operation visible to search, if wait_for then wait for a refresh to make this operation visible to search, if false (the default) then do nothing with refreshes.
* Parameter: routing
  * Type: string
  * Description: Specific routing value
* Parameter: timeout
  * Type: time
  * Description: Explicit operation timeout
* Parameter: version
  * Type: number
  * Description: Explicit version number for concurrency control
* Parameter: version_type
  * Type: enum
  * Description: Specific version type
* Parameter: pipeline
  * Type: string
  * Description: The pipeline id to preprocess incoming documents with


dangling\_indices.delete\_dangling\_index
-----------------------------------------

Deletes the specified dangling index

```
DELETE _dangling/{index_uuid}

```


#### URL parameters


|Parameter       |Type   |Description                                              |
|----------------|-------|---------------------------------------------------------|
|accept_data_loss|boolean|Must be set to true in order to delete the dangling index|
|timeout         |time   |Explicit operation timeout                               |
|master_timeout  |time   |Specify timeout for connection to master                 |


dangling\_indices.import\_dangling\_index
-----------------------------------------

Imports the specified dangling index

```
POST _dangling/{index_uuid}

```


#### URL parameters


|Parameter       |Type   |Description                                              |
|----------------|-------|---------------------------------------------------------|
|accept_data_loss|boolean|Must be set to true in order to import the dangling index|
|timeout         |time   |Explicit operation timeout                               |
|master_timeout  |time   |Specify timeout for connection to master                 |


dangling\_indices.list\_dangling\_indices
-----------------------------------------

Returns all dangling indices.

delete
------

Removes a document from the index.

```
DELETE {index}/{type}/{id}

```


#### URL parameters



* Parameter: wait_for_active_shards
  * Type: string
  * Description: Sets the number of shard copies that must be active before proceeding with the delete operation. Defaults to 1, meaning the primary shard only. Set to all for all shard copies, otherwise set to any non-negative value less than or equal to the total number of copies for the shard (number of replicas + 1)
* Parameter: refresh
  * Type: enum
  * Description: If true then refresh the affected shards to make this operation visible to search, if wait_for then wait for a refresh to make this operation visible to search, if false (the default) then do nothing with refreshes.
* Parameter: routing
  * Type: string
  * Description: Specific routing value
* Parameter: timeout
  * Type: time
  * Description: Explicit operation timeout
* Parameter: if_seq_no
  * Type: number
  * Description: only perform the delete operation if the last operation that has changed the document has the specified sequence number
* Parameter: if_primary_term
  * Type: number
  * Description: only perform the delete operation if the last operation that has changed the document has the specified primary term
* Parameter: version
  * Type: number
  * Description: Explicit version number for concurrency control
* Parameter: version_type
  * Type: enum
  * Description: Specific version type


delete\_by\_query
-----------------

Deletes documents matching the provided query.

```
POST {index}/_delete_by_query

```


```
POST {index}/{type}/_delete_by_query

```


#### HTTP request body

The search definition using the Query DSL

**Required**: True

#### URL parameters



* Parameter: analyzer
  * Type: string
  * Description: The analyzer to use for the query string
  *  :  
* Parameter: analyze_wildcard
  * Type: boolean
  * Description: Specify whether wildcard and prefix queries should be analyzed (default: false)
  *  :  
* Parameter: default_operator
  * Type: enum
  * Description: The default operator for query string query (AND or OR)
  *  :  
* Parameter: df
  * Type: string
  * Description: The field to use as default where no field prefix is given in the query string
  *  :  
* Parameter: from
  * Type: number
  * Description: Starting offset (default: 0)
  *  :  
* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether specified concrete indices should be ignored when unavailable (missing or closed)
  *  :  
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes _all string or when no indices have been specified)
  *  :  
* Parameter: conflicts
  * Type: enum
  * Description: What to do when the delete by query hits version conflicts?
  *  :  
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.
  *  :  
* Parameter: lenient
  * Type: boolean
  * Description: Specify whether format-based query failures (such as providing text to a numeric field) should be ignored
  *  :  
* Parameter: preference
  * Type: string
  * Description: Specify the node or shard the operation should be performed on (default: random)
  *  :  
* Parameter: q
  * Type: string
  * Description: Query in the Lucene query string syntax
  *  :  
* Parameter: routing
  * Type: list
  * Description: A comma-separated list of specific routing values
  *  :  
* Parameter: scroll
  * Type: time
  * Description: Specify how long a consistent view of the index should be maintained for scrolled search
  *  :  
* Parameter: search_type
  * Type: enum
  * Description: Search operation type
  *  :  
* Parameter: search_timeout
  * Type: time
  * Description: Explicit timeout for each search request. Defaults to no timeout.
  *  :  
* Parameter: size
  * Type: number
  * Description: Deprecated, please use max_docs instead
  *  :  
* Parameter: max_docs
  * Type: number
  * Description: Maximum number of documents to process (default: all documents)
  *  :  
* Parameter: sort
  * Type: list
  * Description: A comma-separated list of : pairs
  *  :  
* Parameter: _source
  * Type: list
  * Description: True or false to return the _source field or not, or a list of fields to return
  *  :  
* Parameter: _source_excludes
  * Type: list
  * Description: A list of fields to exclude from the returned _source field
  *  :  
* Parameter: _source_includes
  * Type: list
  * Description: A list of fields to extract and return from the _source field
  *  :  
* Parameter: terminate_after
  * Type: number
  * Description: The maximum number of documents to collect for each shard, upon reaching which the query execution will terminate early.
  *  :  
* Parameter: stats
  * Type: list
  * Description: Specific ‘tag’ of the request for logging and statistical purposes
  *  :  
* Parameter: version
  * Type: boolean
  * Description: Specify whether to return document version as part of a hit
  *  :  
* Parameter: request_cache
  * Type: boolean
  * Description: Specify if request cache should be used for this request or not, defaults to index level setting
  *  :  
* Parameter: refresh
  * Type: boolean
  * Description: Should the effected indexes be refreshed?
  *  :  
* Parameter: timeout
  * Type: time
  * Description: Time each individual bulk request should wait for shards that are unavailable.
  *  :  
* Parameter: wait_for_active_shards
  * Type: string
  * Description: Sets the number of shard copies that must be active before proceeding with the delete by query operation. Defaults to 1, meaning the primary shard only. Set to all for all shard copies, otherwise set to any non-negative value less than or equal to the total number of copies for the shard (number of replicas + 1)
  *  :  
* Parameter: scroll_size
  * Type: number
  * Description: Size on the scroll request powering the delete by query
  *  :  
* Parameter: wait_for_completion
  * Type: boolean
  * Description: Should the request should block until the delete by query is complete.
  *  :  
* Parameter: requests_per_second
  * Type: number
  * Description: The throttle for this request in sub-requests per second. -1 means no throttle.
  *  :  
* Parameter: slices
  * Type: number
  * Description: string
  *  : The number of slices this task should be divided into. Defaults to 1, meaning the task isn’t sliced into subtasks. Can be set to auto.


delete\_by\_query\_rethrottle
-----------------------------

Changes the number of requests per second for a particular Delete By Query operation.

```
POST _delete_by_query/{task_id}/_rethrottle

```


#### URL parameters



* Parameter: requests_per_second
  * Type: number
  * Description: The throttle to set on this request in floating sub-requests per second. -1 means set no throttle.


delete\_script
--------------

Deletes a script.

#### URL parameters


|Parameter     |Type|Description                             |
|--------------|----|----------------------------------------|
|timeout       |time|Explicit operation timeout              |
|master_timeout|time|Specify timeout for connection to master|


exists
------

Returns information about whether a document exists in an index.

#### URL parameters



* Parameter: stored_fields
  * Type: list
  * Description: A comma-separated list of stored fields to return in the response
* Parameter: preference
  * Type: string
  * Description: Specify the node or shard the operation should be performed on (default: random)
* Parameter: realtime
  * Type: boolean
  * Description: Specify whether to perform the operation in realtime or search mode
* Parameter: refresh
  * Type: boolean
  * Description: Refresh the shard containing the document before performing the operation
* Parameter: routing
  * Type: string
  * Description: Specific routing value
* Parameter: _source
  * Type: list
  * Description: True or false to return the _source field or not, or a list of fields to return
* Parameter: _source_excludes
  * Type: list
  * Description: A list of fields to exclude from the returned _source field
* Parameter: _source_includes
  * Type: list
  * Description: A list of fields to extract and return from the _source field
* Parameter: version
  * Type: number
  * Description: Explicit version number for concurrency control
* Parameter: version_type
  * Type: enum
  * Description: Specific version type


exists\_source
--------------

Returns information about whether a document source exists in an index.

```
HEAD {index}/_source/{id}

```


```
HEAD {index}/{type}/{id}/_source

```


#### URL parameters



* Parameter: preference
  * Type: string
  * Description: Specify the node or shard the operation should be performed on (default: random)
* Parameter: realtime
  * Type: boolean
  * Description: Specify whether to perform the operation in realtime or search mode
* Parameter: refresh
  * Type: boolean
  * Description: Refresh the shard containing the document before performing the operation
* Parameter: routing
  * Type: string
  * Description: Specific routing value
* Parameter: _source
  * Type: list
  * Description: True or false to return the _source field or not, or a list of fields to return
* Parameter: _source_excludes
  * Type: list
  * Description: A list of fields to exclude from the returned _source field
* Parameter: _source_includes
  * Type: list
  * Description: A list of fields to extract and return from the _source field
* Parameter: version
  * Type: number
  * Description: Explicit version number for concurrency control
* Parameter: version_type
  * Type: enum
  * Description: Specific version type


explain
-------

Returns information about why a specific matches (or doesn’t match) a query.

```
GET {index}/_explain/{id}
POST {index}/_explain/{id}

```


```
GET {index}/{type}/{id}/_explain
POST {index}/{type}/{id}/_explain

```


#### HTTP request body

The query definition using the Query DSL

#### URL parameters



* Parameter: analyze_wildcard
  * Type: boolean
  * Description: Specify whether wildcards and prefix queries in the query string query should be analyzed (default: false)
* Parameter: analyzer
  * Type: string
  * Description: The analyzer for the query string query
* Parameter: default_operator
  * Type: enum
  * Description: The default operator for query string query (AND or OR)
* Parameter: df
  * Type: string
  * Description: The default field for query string query (default: _all)
* Parameter: stored_fields
  * Type: list
  * Description: A comma-separated list of stored fields to return in the response
* Parameter: lenient
  * Type: boolean
  * Description: Specify whether format-based query failures (such as providing text to a numeric field) should be ignored
* Parameter: preference
  * Type: string
  * Description: Specify the node or shard the operation should be performed on (default: random)
* Parameter: q
  * Type: string
  * Description: Query in the Lucene query string syntax
* Parameter: routing
  * Type: string
  * Description: Specific routing value
* Parameter: _source
  * Type: list
  * Description: True or false to return the _source field or not, or a list of fields to return
* Parameter: _source_excludes
  * Type: list
  * Description: A list of fields to exclude from the returned _source field
* Parameter: _source_includes
  * Type: list
  * Description: A list of fields to extract and return from the _source field


field\_caps
-----------

Returns the information about the capabilities of fields among multiple indices.

```
GET _field_caps
POST _field_caps

```


```
GET {index}/_field_caps
POST {index}/_field_caps

```


#### HTTP request body

An index filter specified with the Query DSL

#### URL parameters



* Parameter: fields
  * Type: list
  * Description: A comma-separated list of field names
* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether specified concrete indices should be ignored when unavailable (missing or closed)
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes _all string or when no indices have been specified)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.
* Parameter: include_unmapped
  * Type: boolean
  * Description: Indicates whether unmapped fields should be included in the response.


get
---

Returns a document.

#### URL parameters



* Parameter: stored_fields
  * Type: list
  * Description: A comma-separated list of stored fields to return in the response
* Parameter: preference
  * Type: string
  * Description: Specify the node or shard the operation should be performed on (default: random)
* Parameter: realtime
  * Type: boolean
  * Description: Specify whether to perform the operation in realtime or search mode
* Parameter: refresh
  * Type: boolean
  * Description: Refresh the shard containing the document before performing the operation
* Parameter: routing
  * Type: string
  * Description: Specific routing value
* Parameter: _source
  * Type: list
  * Description: True or false to return the _source field or not, or a list of fields to return
* Parameter: _source_excludes
  * Type: list
  * Description: A list of fields to exclude from the returned _source field
* Parameter: _source_includes
  * Type: list
  * Description: A list of fields to extract and return from the _source field
* Parameter: version
  * Type: number
  * Description: Explicit version number for concurrency control
* Parameter: version_type
  * Type: enum
  * Description: Specific version type


get\_script
-----------

Returns a script.

#### URL parameters


|Parameter     |Type|Description                             |
|--------------|----|----------------------------------------|
|master_timeout|time|Specify timeout for connection to master|


get\_script\_context
--------------------

Returns all script contexts.

get\_script\_languages
----------------------

Returns available script types, languages and contexts

get\_source
-----------

Returns the source of a document.

```
GET {index}/{type}/{id}/_source

```


#### URL parameters



* Parameter: preference
  * Type: string
  * Description: Specify the node or shard the operation should be performed on (default: random)
* Parameter: realtime
  * Type: boolean
  * Description: Specify whether to perform the operation in realtime or search mode
* Parameter: refresh
  * Type: boolean
  * Description: Refresh the shard containing the document before performing the operation
* Parameter: routing
  * Type: string
  * Description: Specific routing value
* Parameter: _source
  * Type: list
  * Description: True or false to return the _source field or not, or a list of fields to return
* Parameter: _source_excludes
  * Type: list
  * Description: A list of fields to exclude from the returned _source field
* Parameter: _source_includes
  * Type: list
  * Description: A list of fields to extract and return from the _source field
* Parameter: version
  * Type: number
  * Description: Explicit version number for concurrency control
* Parameter: version_type
  * Type: enum
  * Description: Specific version type


index
-----

Creates or updates a document in an index.

```
PUT {index}/_doc/{id}
POST {index}/_doc/{id}

```


```
PUT {index}/{type}/{id}
POST {index}/{type}/{id}

```


#### HTTP request body

The document

**Required**: True

#### URL parameters



* Parameter: wait_for_active_shards
  * Type: string
  * Description: Sets the number of shard copies that must be active before proceeding with the index operation. Defaults to 1, meaning the primary shard only. Set to all for all shard copies, otherwise set to any non-negative value less than or equal to the total number of copies for the shard (number of replicas + 1)
* Parameter: op_type
  * Type: enum
  * Description: Explicit operation type. Defaults to index for requests with an explicit document ID, and to createfor requests without an explicit document ID
* Parameter: refresh
  * Type: enum
  * Description: If true then refresh the affected shards to make this operation visible to search, if wait_for then wait for a refresh to make this operation visible to search, if false (the default) then do nothing with refreshes.
* Parameter: routing
  * Type: string
  * Description: Specific routing value
* Parameter: timeout
  * Type: time
  * Description: Explicit operation timeout
* Parameter: version
  * Type: number
  * Description: Explicit version number for concurrency control
* Parameter: version_type
  * Type: enum
  * Description: Specific version type
* Parameter: if_seq_no
  * Type: number
  * Description: only perform the index operation if the last operation that has changed the document has the specified sequence number
* Parameter: if_primary_term
  * Type: number
  * Description: only perform the index operation if the last operation that has changed the document has the specified primary term
* Parameter: pipeline
  * Type: string
  * Description: The pipeline id to preprocess incoming documents with
* Parameter: require_alias
  * Type: boolean
  * Description: When true, requires destination to be an alias. Default is false


indices.add\_block
------------------

Adds a block to an index.

```
PUT {index}/_block/{block}

```


#### URL parameters



* Parameter: timeout
  * Type: time
  * Description: Explicit operation timeout
* Parameter: master_timeout
  * Type: time
  * Description: Specify timeout for connection to master
* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether specified concrete indices should be ignored when unavailable (missing or closed)
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes _all string or when no indices have been specified)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.


indices.analyze
---------------

Performs the analysis process on a text and return the tokens breakdown of the text.

```
GET _analyze
POST _analyze

```


```
GET {index}/_analyze
POST {index}/_analyze

```


#### HTTP request body

Define analyzer/tokenizer parameters and the text on which the analysis should be performed

#### URL parameters


|Parameter|Type  |Description                                 |
|---------|------|--------------------------------------------|
|index    |string|The name of the index to scope the operation|


indices.clear\_cache
--------------------

Clears all or specific caches for one or more indices.

```
POST {index}/_cache/clear

```


#### URL parameters



* Parameter: fielddata
  * Type: boolean
  * Description: Clear field data
* Parameter: fields
  * Type: list
  * Description: A comma-separated list of fields to clear when using the fielddata parameter (default: all)
* Parameter: query
  * Type: boolean
  * Description: Clear query caches
* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether specified concrete indices should be ignored when unavailable (missing or closed)
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes _all string or when no indices have been specified)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.
* Parameter: index
  * Type: list
  * Description: A comma-separated list of index name to limit the operation
* Parameter: request
  * Type: boolean
  * Description: Clear request cache


indices.clone
-------------

Clones an index

```
PUT {index}/_clone/{target}
POST {index}/_clone/{target}

```


#### HTTP request body

The configuration for the target index (`settings` and `aliases`)

#### URL parameters



* Parameter: timeout
  * Type: time
  * Description: Explicit operation timeout
* Parameter: master_timeout
  * Type: time
  * Description: Specify timeout for connection to master
* Parameter: wait_for_active_shards
  * Type: string
  * Description: Set the number of active shards to wait for on the cloned index before the operation returns.


indices.close
-------------

Closes an index.

#### URL parameters



* Parameter: timeout
  * Type: time
  * Description: Explicit operation timeout
* Parameter: master_timeout
  * Type: time
  * Description: Specify timeout for connection to master
* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether specified concrete indices should be ignored when unavailable (missing or closed)
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes _all string or when no indices have been specified)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.
* Parameter: wait_for_active_shards
  * Type: string
  * Description: Sets the number of active shards to wait for before the operation returns.


indices.create
--------------

Creates an index with optional settings and mappings.

#### HTTP request body

The configuration for the index (`settings` and `mappings`)

#### URL parameters



* Parameter: include_type_name
  * Type: boolean
  * Description: Whether a type should be expected in the body of the mappings.
* Parameter: wait_for_active_shards
  * Type: string
  * Description: Set the number of active shards to wait for before the operation returns.
* Parameter: timeout
  * Type: time
  * Description: Explicit operation timeout
* Parameter: master_timeout
  * Type: time
  * Description: Specify timeout for connection to master


indices.delete
--------------

Deletes an index.

#### URL parameters



* Parameter: timeout
  * Type: time
  * Description: Explicit operation timeout
* Parameter: master_timeout
  * Type: time
  * Description: Specify timeout for connection to master
* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Ignore unavailable indexes (default: false)
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Ignore if a wildcard expression resolves to no concrete indices (default: false)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether wildcard expressions should get expanded to open or closed indices (default: open)


indices.delete\_alias
---------------------

Deletes an alias.

```
DELETE {index}/_alias/{name}

```


```
DELETE {index}/_aliases/{name}

```


#### URL parameters


|Parameter     |Type|Description                             |
|--------------|----|----------------------------------------|
|timeout       |time|Explicit timestamp for the document     |
|master_timeout|time|Specify timeout for connection to master|


indices.delete\_index\_template
-------------------------------

Deletes an index template.

```
DELETE _index_template/{name}

```


#### URL parameters


|Parameter     |Type|Description                             |
|--------------|----|----------------------------------------|
|timeout       |time|Explicit operation timeout              |
|master_timeout|time|Specify timeout for connection to master|


indices.delete\_template
------------------------

Deletes an index template.

#### URL parameters


|Parameter     |Type|Description                             |
|--------------|----|----------------------------------------|
|timeout       |time|Explicit operation timeout              |
|master_timeout|time|Specify timeout for connection to master|


indices.exists
--------------

Returns information about whether a particular index exists.

#### URL parameters



* Parameter: local
  * Type: boolean
  * Description: Return local information, do not retrieve the state from master node (default: false)
* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Ignore unavailable indexes (default: false)
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Ignore if a wildcard expression resolves to no concrete indices (default: false)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether wildcard expressions should get expanded to open or closed indices (default: open)
* Parameter: flat_settings
  * Type: boolean
  * Description: Return settings in flat format (default: false)
* Parameter: include_defaults
  * Type: boolean
  * Description: Whether to return all default setting for each of the indices.


indices.exists\_alias
---------------------

Returns information about whether a particular alias exists.

```
HEAD {index}/_alias/{name}

```


#### URL parameters



* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether specified concrete indices should be ignored when unavailable (missing or closed)
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes _all string or when no indices have been specified)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.
* Parameter: local
  * Type: boolean
  * Description: Return local information, do not retrieve the state from master node (default: false)


indices.exists\_index\_template
-------------------------------

Returns information about whether a particular index template exists.

```
HEAD _index_template/{name}

```


#### URL parameters



* Parameter: flat_settings
  * Type: boolean
  * Description: Return settings in flat format (default: false)
* Parameter: master_timeout
  * Type: time
  * Description: Explicit operation timeout for connection to master node
* Parameter: local
  * Type: boolean
  * Description: Return local information, do not retrieve the state from master node (default: false)


indices.exists\_template
------------------------

Returns information about whether a particular index template exists.

#### URL parameters



* Parameter: flat_settings
  * Type: boolean
  * Description: Return settings in flat format (default: false)
* Parameter: master_timeout
  * Type: time
  * Description: Explicit operation timeout for connection to master node
* Parameter: local
  * Type: boolean
  * Description: Return local information, do not retrieve the state from master node (default: false)


indices.exists\_type
--------------------

Returns information about whether a particular document type exists. (DEPRECATED)

```
HEAD {index}/_mapping/{type}

```


#### URL parameters



* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether specified concrete indices should be ignored when unavailable (missing or closed)
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes _all string or when no indices have been specified)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.
* Parameter: local
  * Type: boolean
  * Description: Return local information, do not retrieve the state from master node (default: false)


indices.flush
-------------

Performs the flush operation on one or more indices.

```
POST {index}/_flush
GET {index}/_flush

```


#### URL parameters



* Parameter: force
  * Type: boolean
  * Description: Whether a flush should be forced even if it is not necessarily needed ie. if no changes will be committed to the index. This is useful if transaction log IDs should be incremented even if no uncommitted changes are present. (This setting can be considered as internal)
* Parameter: wait_if_ongoing
  * Type: boolean
  * Description: If set to true the flush operation will block until the flush can be executed if another flush operation is already executing. The default is true. If set to false the flush will be skipped iff if another flush operation is already running.
* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether specified concrete indices should be ignored when unavailable (missing or closed)
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes _all string or when no indices have been specified)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.


indices.flush\_synced
---------------------

Performs a synced flush operation on one or more indices. Synced flush is deprecated and will be removed in 8.0. Use flush instead

```
POST _flush/synced
GET _flush/synced

```


```
POST {index}/_flush/synced
GET {index}/_flush/synced

```


#### URL parameters



* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether specified concrete indices should be ignored when unavailable (missing or closed)
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes _all string or when no indices have been specified)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.


indices.forcemerge
------------------

Performs the force merge operation on one or more indices.

#### URL parameters



* Parameter: flush
  * Type: boolean
  * Description: Specify whether the index should be flushed after performing the operation (default: true)
* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether specified concrete indices should be ignored when unavailable (missing or closed)
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes _all string or when no indices have been specified)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.
* Parameter: max_num_segments
  * Type: number
  * Description: The number of segments the index should be merged into (default: dynamic)
* Parameter: only_expunge_deletes
  * Type: boolean
  * Description: Specify whether the operation should only expunge deleted documents


indices.get
-----------

Returns information about one or more indices.

#### URL parameters



* Parameter: include_type_name
  * Type: boolean
  * Description: Whether to add the type name to the response (default: false)
* Parameter: local
  * Type: boolean
  * Description: Return local information, do not retrieve the state from master node (default: false)
* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Ignore unavailable indexes (default: false)
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Ignore if a wildcard expression resolves to no concrete indices (default: false)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether wildcard expressions should get expanded to open or closed indices (default: open)
* Parameter: flat_settings
  * Type: boolean
  * Description: Return settings in flat format (default: false)
* Parameter: include_defaults
  * Type: boolean
  * Description: Whether to return all default setting for each of the indices.
* Parameter: master_timeout
  * Type: time
  * Description: Specify timeout for connection to master


indices.get\_alias
------------------

Returns an alias.

```
GET {index}/_alias/{name}

```


#### URL parameters



* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether specified concrete indices should be ignored when unavailable (missing or closed)
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes _all string or when no indices have been specified)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.
* Parameter: local
  * Type: boolean
  * Description: Return local information, do not retrieve the state from master node (default: false)


indices.get\_field\_mapping
---------------------------

Returns mapping for one or more fields.

```
GET _mapping/field/{fields}

```


```
GET {index}/_mapping/field/{fields}

```


```
GET _mapping/{type}/field/{fields}

```


```
GET {index}/_mapping/{type}/field/{fields}

```


#### URL parameters



* Parameter: include_type_name
  * Type: boolean
  * Description: Whether a type should be returned in the body of the mappings.
* Parameter: include_defaults
  * Type: boolean
  * Description: Whether the default mapping values should be returned as well
* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether specified concrete indices should be ignored when unavailable (missing or closed)
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes _all string or when no indices have been specified)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.
* Parameter: local
  * Type: boolean
  * Description: Return local information, do not retrieve the state from master node (default: false)


indices.get\_index\_template
----------------------------

Returns an index template.

```
GET _index_template/{name}

```


#### URL parameters



* Parameter: flat_settings
  * Type: boolean
  * Description: Return settings in flat format (default: false)
* Parameter: master_timeout
  * Type: time
  * Description: Explicit operation timeout for connection to master node
* Parameter: local
  * Type: boolean
  * Description: Return local information, do not retrieve the state from master node (default: false)


indices.get\_mapping
--------------------

Returns mappings for one or more indices.

```
GET {index}/_mapping/{type}

```


#### URL parameters



* Parameter: include_type_name
  * Type: boolean
  * Description: Whether to add the type name to the response (default: false)
* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether specified concrete indices should be ignored when unavailable (missing or closed)
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes _all string or when no indices have been specified)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.
* Parameter: master_timeout
  * Type: time
  * Description: Specify timeout for connection to master
* Parameter: local
  * Type: boolean
  * Description: Return local information, do not retrieve the state from master node (default: false)


indices.get\_settings
---------------------

Returns settings for one or more indices.

```
GET {index}/_settings/{name}

```


#### URL parameters



* Parameter: master_timeout
  * Type: time
  * Description: Specify timeout for connection to master
* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether specified concrete indices should be ignored when unavailable (missing or closed)
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes _all string or when no indices have been specified)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.
* Parameter: flat_settings
  * Type: boolean
  * Description: Return settings in flat format (default: false)
* Parameter: local
  * Type: boolean
  * Description: Return local information, do not retrieve the state from master node (default: false)
* Parameter: include_defaults
  * Type: boolean
  * Description: Whether to return all default setting for each of the indices.


indices.get\_template
---------------------

Returns an index template.

#### URL parameters



* Parameter: include_type_name
  * Type: boolean
  * Description: Whether a type should be returned in the body of the mappings.
* Parameter: flat_settings
  * Type: boolean
  * Description: Return settings in flat format (default: false)
* Parameter: master_timeout
  * Type: time
  * Description: Explicit operation timeout for connection to master node
* Parameter: local
  * Type: boolean
  * Description: Return local information, do not retrieve the state from master node (default: false)


indices.get\_upgrade
--------------------

The \_upgrade API is no longer useful and will be removed.

#### URL parameters



* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether specified concrete indices should be ignored when unavailable (missing or closed)
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes _all string or when no indices have been specified)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.


indices.open
------------

Opens an index.

#### URL parameters



* Parameter: timeout
  * Type: time
  * Description: Explicit operation timeout
* Parameter: master_timeout
  * Type: time
  * Description: Specify timeout for connection to master
* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether specified concrete indices should be ignored when unavailable (missing or closed)
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes _all string or when no indices have been specified)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.
* Parameter: wait_for_active_shards
  * Type: string
  * Description: Sets the number of active shards to wait for before the operation returns.


indices.put\_alias
------------------

Creates or updates an alias.

```
PUT {index}/_alias/{name}
POST {index}/_alias/{name}

```


```
PUT {index}/_aliases/{name}
POST {index}/_aliases/{name}

```


#### HTTP request body

The settings for the alias, such as `routing` or `filter`

**Required**: False

#### URL parameters


|Parameter     |Type|Description                             |
|--------------|----|----------------------------------------|
|timeout       |time|Explicit timestamp for the document     |
|master_timeout|time|Specify timeout for connection to master|


indices.put\_index\_template
----------------------------

Creates or updates an index template.

```
PUT _index_template/{name}
POST _index_template/{name}

```


#### HTTP request body

The template definition

**Required**: True

#### URL parameters



* Parameter: create
  * Type: boolean
  * Description: Whether the index template should only be added if new or can also replace an existing one
* Parameter: cause
  * Type: string
  * Description: User defined reason for creating/updating the index template
* Parameter: master_timeout
  * Type: time
  * Description: Specify timeout for connection to master


indices.put\_mapping
--------------------

Updates the index mappings.

```
PUT {index}/_mapping
POST {index}/_mapping

```


```
PUT {index}/{type}/_mapping
POST {index}/{type}/_mapping

```


```
PUT {index}/_mapping/{type}
POST {index}/_mapping/{type}

```


```
PUT {index}/{type}/_mappings
POST {index}/{type}/_mappings

```


```
PUT {index}/_mappings/{type}
POST {index}/_mappings/{type}

```


```
PUT _mappings/{type}
POST _mappings/{type}

```


```
PUT {index}/_mappings
POST {index}/_mappings

```


```
PUT _mapping/{type}
POST _mapping/{type}

```


#### HTTP request body

The mapping definition

**Required**: True

#### URL parameters



* Parameter: include_type_name
  * Type: boolean
  * Description: Whether a type should be expected in the body of the mappings.
* Parameter: timeout
  * Type: time
  * Description: Explicit operation timeout
* Parameter: master_timeout
  * Type: time
  * Description: Specify timeout for connection to master
* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether specified concrete indices should be ignored when unavailable (missing or closed)
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes _all string or when no indices have been specified)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.
* Parameter: write_index_only
  * Type: boolean
  * Description: When true, applies mappings only to the write index of an alias or data stream


indices.put\_settings
---------------------

Updates the index settings.

#### HTTP request body

The index settings to be updated

**Required**: True

#### URL parameters



* Parameter: master_timeout
  * Type: time
  * Description: Specify timeout for connection to master
* Parameter: timeout
  * Type: time
  * Description: Explicit operation timeout
* Parameter: preserve_existing
  * Type: boolean
  * Description: Whether to update existing settings. If set to true existing settings on an index remain unchanged, the default is false
* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether specified concrete indices should be ignored when unavailable (missing or closed)
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes _all string or when no indices have been specified)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.
* Parameter: flat_settings
  * Type: boolean
  * Description: Return settings in flat format (default: false)


indices.put\_template
---------------------

Creates or updates an index template.

```
PUT _template/{name}
POST _template/{name}

```


#### HTTP request body

The template definition

**Required**: True

#### URL parameters



* Parameter: include_type_name
  * Type: boolean
  * Description: Whether a type should be returned in the body of the mappings.
* Parameter: order
  * Type: number
  * Description: The order for this template when merging multiple matching ones (higher numbers are merged later, overriding the lower numbers)
* Parameter: create
  * Type: boolean
  * Description: Whether the index template should only be added if new or can also replace an existing one
* Parameter: master_timeout
  * Type: time
  * Description: Specify timeout for connection to master


indices.recovery
----------------

Returns information about ongoing index shard recoveries.

#### URL parameters


|Parameter  |Type   |Description                                                 |
|-----------|-------|------------------------------------------------------------|
|detailed   |boolean|Whether to display detailed information about shard recovery|
|active_only|boolean|Display only those recoveries that are currently on-going   |


indices.refresh
---------------

Performs the refresh operation in one or more indices.

```
POST _refresh
GET _refresh

```


```
POST {index}/_refresh
GET {index}/_refresh

```


#### URL parameters



* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether specified concrete indices should be ignored when unavailable (missing or closed)
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes _all string or when no indices have been specified)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.


indices.resolve\_index
----------------------

Returns information about any matching indices, aliases, and data streams

```
GET _resolve/index/{name}

```


#### URL parameters



* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether wildcard expressions should get expanded to open or closed indices (default: open)


indices.rollover
----------------

Updates an alias to point to a new index when the existing index is considered to be too large or too old.

```
POST {alias}/_rollover/{new_index}

```


#### HTTP request body

The conditions that needs to be met for executing rollover

#### URL parameters



* Parameter: include_type_name
  * Type: boolean
  * Description: Whether a type should be included in the body of the mappings.
* Parameter: timeout
  * Type: time
  * Description: Explicit operation timeout
* Parameter: dry_run
  * Type: boolean
  * Description: If set to true the rollover action will only be validated but not actually performed even if a condition matches. The default is false
* Parameter: master_timeout
  * Type: time
  * Description: Specify timeout for connection to master
* Parameter: wait_for_active_shards
  * Type: string
  * Description: Set the number of active shards to wait for on the newly created rollover index before the operation returns.


indices.segments
----------------

Provides low-level information about segments in a Lucene index.

#### URL parameters



* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether specified concrete indices should be ignored when unavailable (missing or closed)
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes _all string or when no indices have been specified)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.
* Parameter: verbose
  * Type: boolean
  * Description: Includes detailed memory usage by Lucene.


indices.shard\_stores
---------------------

Provides store information for shard copies of indices.

```
GET {index}/_shard_stores

```


#### URL parameters



* Parameter: status
  * Type: list
  * Description: A comma-separated list of statuses used to filter on shards to get store information for
* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether specified concrete indices should be ignored when unavailable (missing or closed)
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes _all string or when no indices have been specified)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.


indices.shrink
--------------

Allow to shrink an existing index into a new index with fewer primary shards.

```
PUT {index}/_shrink/{target}
POST {index}/_shrink/{target}

```


#### HTTP request body

The configuration for the target index (`settings` and `aliases`)

#### URL parameters



* Parameter: copy_settings
  * Type: boolean
  * Description: whether or not to copy settings from the source index (defaults to false)
* Parameter: timeout
  * Type: time
  * Description: Explicit operation timeout
* Parameter: master_timeout
  * Type: time
  * Description: Specify timeout for connection to master
* Parameter: wait_for_active_shards
  * Type: string
  * Description: Set the number of active shards to wait for on the shrunken index before the operation returns.


indices.simulate\_index\_template
---------------------------------

Simulate matching the given index name against the index templates in the system

```
POST _index_template/_simulate_index/{name}

```


#### HTTP request body

New index template definition, which will be included in the simulation, as if it already exists in the system

**Required**: False

#### URL parameters



* Parameter: create
  * Type: boolean
  * Description: Whether the index template we optionally defined in the body should only be dry-run added if new or can also replace an existing one
* Parameter: cause
  * Type: string
  * Description: User defined reason for dry-run creating the new template for simulation purposes
* Parameter: master_timeout
  * Type: time
  * Description: Specify timeout for connection to master


indices.simulate\_template
--------------------------

Simulate resolving the given template name or body

```
POST _index_template/_simulate

```


```
POST _index_template/_simulate/{name}

```


#### HTTP request body

New index template definition to be simulated, if no index template name is specified

**Required**: False

#### URL parameters



* Parameter: create
  * Type: boolean
  * Description: Whether the index template we optionally defined in the body should only be dry-run added if new or can also replace an existing one
* Parameter: cause
  * Type: string
  * Description: User defined reason for dry-run creating the new template for simulation purposes
* Parameter: master_timeout
  * Type: time
  * Description: Specify timeout for connection to master


indices.split
-------------

Allows you to split an existing index into a new index with more primary shards.

```
PUT {index}/_split/{target}
POST {index}/_split/{target}

```


#### HTTP request body

The configuration for the target index (`settings` and `aliases`)

#### URL parameters



* Parameter: copy_settings
  * Type: boolean
  * Description: whether or not to copy settings from the source index (defaults to false)
* Parameter: timeout
  * Type: time
  * Description: Explicit operation timeout
* Parameter: master_timeout
  * Type: time
  * Description: Specify timeout for connection to master
* Parameter: wait_for_active_shards
  * Type: string
  * Description: Set the number of active shards to wait for on the shrunken index before the operation returns.


indices.stats
-------------

Provides statistics on operations happening in an index.

```
GET {index}/_stats/{metric}

```


#### URL parameters



* Parameter: completion_fields
  * Type: list
  * Description: A comma-separated list of fields for fielddata and suggest index metric (supports wildcards)
* Parameter: fielddata_fields
  * Type: list
  * Description: A comma-separated list of fields for fielddata index metric (supports wildcards)
* Parameter: fields
  * Type: list
  * Description: A comma-separated list of fields for fielddata and completion index metric (supports wildcards)
* Parameter: groups
  * Type: list
  * Description: A comma-separated list of search groups for search index metric
* Parameter: level
  * Type: enum
  * Description: Return stats aggregated at cluster, index or shard level
* Parameter: types
  * Type: list
  * Description: A comma-separated list of document types for the indexing index metric
* Parameter: include_segment_file_sizes
  * Type: boolean
  * Description: Whether to report the aggregated disk usage of each one of the Lucene index files (only applies if segment stats are requested)
* Parameter: include_unloaded_segments
  * Type: boolean
  * Description: If set to true segment stats will include stats for segments that are not currently loaded into memory
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.
* Parameter: forbid_closed_indices
  * Type: boolean
  * Description: If set to false stats will also collected from closed indices if explicitly specified or if expand_wildcards expands to closed indices


indices.update\_aliases
-----------------------

Updates index aliases.

#### HTTP request body

The definition of `actions` to perform

**Required**: True

#### URL parameters


|Parameter     |Type|Description                             |
|--------------|----|----------------------------------------|
|timeout       |time|Request timeout                         |
|master_timeout|time|Specify timeout for connection to master|


indices.upgrade
---------------

The \_upgrade API is no longer useful and will be removed.

#### URL parameters



* Parameter: allow_no_indices
  * Type: boolean
  * Description: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes _all string or when no indices have been specified)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.
* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether specified concrete indices should be ignored when unavailable (missing or closed)
* Parameter: wait_for_completion
  * Type: boolean
  * Description: Specify whether the request should block until the all segments are upgraded (default: false)
* Parameter: only_ancient_segments
  * Type: boolean
  * Description: If true, only ancient (an older Lucene major release) segments will be upgraded


indices.validate\_query
-----------------------

Allows a user to validate a potentially expensive query without executing it.

```
GET _validate/query
POST _validate/query

```


```
GET {index}/_validate/query
POST {index}/_validate/query

```


```
GET {index}/{type}/_validate/query
POST {index}/{type}/_validate/query

```


#### HTTP request body

The query definition specified with the Query DSL

#### URL parameters



* Parameter: explain
  * Type: boolean
  * Description: Return detailed information about the error
* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether specified concrete indices should be ignored when unavailable (missing or closed)
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes _all string or when no indices have been specified)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.
* Parameter: q
  * Type: string
  * Description: Query in the Lucene query string syntax
* Parameter: analyzer
  * Type: string
  * Description: The analyzer to use for the query string
* Parameter: analyze_wildcard
  * Type: boolean
  * Description: Specify whether wildcard and prefix queries should be analyzed (default: false)
* Parameter: default_operator
  * Type: enum
  * Description: The default operator for query string query (AND or OR)
* Parameter: df
  * Type: string
  * Description: The field to use as default where no field prefix is given in the query string
* Parameter: lenient
  * Type: boolean
  * Description: Specify whether format-based query failures (such as providing text to a numeric field) should be ignored
* Parameter: rewrite
  * Type: boolean
  * Description: Provide a more detailed explanation showing the actual Lucene query that will be executed.
* Parameter: all_shards
  * Type: boolean
  * Description: Execute validation on all shards instead of one random shard per index


info
----

Returns basic information about the cluster.

ingest.delete\_pipeline
-----------------------

Deletes a pipeline.

```
DELETE _ingest/pipeline/{id}

```


#### URL parameters


|Parameter     |Type|Description                                             |
|--------------|----|--------------------------------------------------------|
|master_timeout|time|Explicit operation timeout for connection to master node|
|timeout       |time|Explicit operation timeout                              |


ingest.get\_pipeline
--------------------

Returns a pipeline.

```
GET _ingest/pipeline/{id}

```


#### URL parameters


|Parameter     |Type|Description                                             |
|--------------|----|--------------------------------------------------------|
|master_timeout|time|Explicit operation timeout for connection to master node|


ingest.processor\_grok
----------------------

Returns a list of the built-in patterns.

```
GET _ingest/processor/grok

```


ingest.put\_pipeline
--------------------

Creates or updates a pipeline.

```
PUT _ingest/pipeline/{id}

```


#### HTTP request body

The ingest definition

**Required**: True

#### URL parameters


|Parameter     |Type|Description                                             |
|--------------|----|--------------------------------------------------------|
|master_timeout|time|Explicit operation timeout for connection to master node|
|timeout       |time|Explicit operation timeout                              |


ingest.simulate
---------------

Allows to simulate a pipeline with example documents.

```
GET _ingest/pipeline/_simulate
POST _ingest/pipeline/_simulate

```


```
GET _ingest/pipeline/{id}/_simulate
POST _ingest/pipeline/{id}/_simulate

```


#### HTTP request body

The simulate definition

**Required**: True

#### URL parameters


|Parameter|Type   |Description                                                              |
|---------|-------|-------------------------------------------------------------------------|
|verbose  |boolean|Verbose mode. Display data output for each processor in executed pipeline|


mget
----

Allows to get multiple documents in one request.

```
GET {index}/_mget
POST {index}/_mget

```


```
GET {index}/{type}/_mget
POST {index}/{type}/_mget

```


#### HTTP request body

Document identifiers; can be either `docs` (containing full document information) or `ids` (when index and type is provided in the URL.

**Required**: True

#### URL parameters



* Parameter: stored_fields
  * Type: list
  * Description: A comma-separated list of stored fields to return in the response
* Parameter: preference
  * Type: string
  * Description: Specify the node or shard the operation should be performed on (default: random)
* Parameter: realtime
  * Type: boolean
  * Description: Specify whether to perform the operation in realtime or search mode
* Parameter: refresh
  * Type: boolean
  * Description: Refresh the shard containing the document before performing the operation
* Parameter: routing
  * Type: string
  * Description: Specific routing value
* Parameter: _source
  * Type: list
  * Description: True or false to return the _source field or not, or a list of fields to return
* Parameter: _source_excludes
  * Type: list
  * Description: A list of fields to exclude from the returned _source field
* Parameter: _source_includes
  * Type: list
  * Description: A list of fields to extract and return from the _source field


msearch
-------

Allows to execute several search operations in one request.

```
GET _msearch
POST _msearch

```


```
GET {index}/_msearch
POST {index}/_msearch

```


```
GET {index}/{type}/_msearch
POST {index}/{type}/_msearch

```


#### HTTP request body

The request definitions (metadata-search request definition pairs), separated by newlines

**Required**: True

#### URL parameters



* Parameter: search_type
  * Type: enum
  * Description: Search operation type
* Parameter: max_concurrent_searches
  * Type: number
  * Description: Controls the maximum number of concurrent searches the multi search api will execute
* Parameter: typed_keys
  * Type: boolean
  * Description: Specify whether aggregation and suggester names should be prefixed by their respective types in the response
* Parameter: pre_filter_shard_size
  * Type: number
  * Description: A threshold that enforces a pre-filter roundtrip to prefilter search shards based on query rewriting if the number of shards the search request expands to exceeds the threshold. This filter roundtrip can limit the number of shards significantly if for instance a shard can not match any documents based on its rewrite method ie. if date filters are mandatory to match but the shard bounds and the query are disjoint.
* Parameter: max_concurrent_shard_requests
  * Type: number
  * Description: The number of concurrent shard requests each sub search executes concurrently per node. This value should be used to limit the impact of the search on the cluster in order to limit the number of concurrent shard requests
* Parameter: rest_total_hits_as_int
  * Type: boolean
  * Description: Indicates whether hits.total should be rendered as an integer or an object in the rest search response
* Parameter: ccs_minimize_roundtrips
  * Type: boolean
  * Description: Indicates whether network round-trips should be minimized as part of cross-cluster search requests execution


msearch\_template
-----------------

Allows to execute several search template operations in one request.

```
GET _msearch/template
POST _msearch/template

```


```
GET {index}/_msearch/template
POST {index}/_msearch/template

```


```
GET {index}/{type}/_msearch/template
POST {index}/{type}/_msearch/template

```


#### HTTP request body

The request definitions (metadata-search request definition pairs), separated by newlines

**Required**: True

#### URL parameters



* Parameter: search_type
  * Type: enum
  * Description: Search operation type
* Parameter: typed_keys
  * Type: boolean
  * Description: Specify whether aggregation and suggester names should be prefixed by their respective types in the response
* Parameter: max_concurrent_searches
  * Type: number
  * Description: Controls the maximum number of concurrent searches the multi search api will execute
* Parameter: rest_total_hits_as_int
  * Type: boolean
  * Description: Indicates whether hits.total should be rendered as an integer or an object in the rest search response
* Parameter: ccs_minimize_roundtrips
  * Type: boolean
  * Description: Indicates whether network round-trips should be minimized as part of cross-cluster search requests execution


mtermvectors
------------

Returns multiple termvectors in one request.

```
GET _mtermvectors
POST _mtermvectors

```


```
GET {index}/_mtermvectors
POST {index}/_mtermvectors

```


```
GET {index}/{type}/_mtermvectors
POST {index}/{type}/_mtermvectors

```


#### HTTP request body

Define ids, documents, parameters or a list of parameters per document here. You must at least provide a list of document ids. See documentation.

**Required**: False

#### URL parameters



* Parameter: ids
  * Type: list
  * Description: A comma-separated list of documents ids. You must define ids as parameter or set “ids” or “docs” in the request body
* Parameter: term_statistics
  * Type: boolean
  * Description: Specifies if total term frequency and document frequency should be returned. Applies to all returned documents unless otherwise specified in body “params” or “docs”.
* Parameter: field_statistics
  * Type: boolean
  * Description: Specifies if document count, sum of document frequencies and sum of total term frequencies should be returned. Applies to all returned documents unless otherwise specified in body “params” or “docs”.
* Parameter: fields
  * Type: list
  * Description: A comma-separated list of fields to return. Applies to all returned documents unless otherwise specified in body “params” or “docs”.
* Parameter: offsets
  * Type: boolean
  * Description: Specifies if term offsets should be returned. Applies to all returned documents unless otherwise specified in body “params” or “docs”.
* Parameter: positions
  * Type: boolean
  * Description: Specifies if term positions should be returned. Applies to all returned documents unless otherwise specified in body “params” or “docs”.
* Parameter: payloads
  * Type: boolean
  * Description: Specifies if term payloads should be returned. Applies to all returned documents unless otherwise specified in body “params” or “docs”.
* Parameter: preference
  * Type: string
  * Description: Specify the node or shard the operation should be performed on (default: random) .Applies to all returned documents unless otherwise specified in body “params” or “docs”.
* Parameter: routing
  * Type: string
  * Description: Specific routing value. Applies to all returned documents unless otherwise specified in body “params” or “docs”.
* Parameter: realtime
  * Type: boolean
  * Description: Specifies if requests are real-time as opposed to near-real-time (default: true).
* Parameter: version
  * Type: number
  * Description: Explicit version number for concurrency control
* Parameter: version_type
  * Type: enum
  * Description: Specific version type


nodes.hot\_threads
------------------

Returns information about hot threads on each node in the cluster.

```
GET _nodes/{node_id}/hot_threads

```


```
GET _cluster/nodes/hotthreads

```


```
GET _cluster/nodes/{node_id}/hotthreads

```


```
GET _nodes/{node_id}/hotthreads

```


```
GET _cluster/nodes/hot_threads

```


```
GET _cluster/nodes/{node_id}/hot_threads

```


#### URL parameters



* Parameter: interval
  * Type: time
  * Description: The interval for the second sampling of threads
* Parameter: snapshots
  * Type: number
  * Description: Number of samples of thread stacktrace (default: 10)
* Parameter: threads
  * Type: number
  * Description: Specify the number of threads to provide information for (default: 3)
* Parameter: ignore_idle_threads
  * Type: boolean
  * Description: Don’t show threads that are in known-idle places, such as waiting on a socket select or pulling from an empty task queue (default: true)
* Parameter: type
  * Type: enum
  * Description: The type to sample (default: cpu)
* Parameter: timeout
  * Type: time
  * Description: Explicit operation timeout


nodes.info
----------

Returns information about nodes in the cluster.

```
GET _nodes/{node_id}/{metric}

```


#### URL parameters


|Parameter    |Type   |Description                                    |
|-------------|-------|-----------------------------------------------|
|flat_settings|boolean|Return settings in flat format (default: false)|
|timeout      |time   |Explicit operation timeout                     |


nodes.reload\_secure\_settings
------------------------------

Reloads secure settings.

```
POST _nodes/reload_secure_settings

```


```
POST _nodes/{node_id}/reload_secure_settings

```


#### HTTP request body

An object containing the password for the elasticsearch keystore

**Required**: False

#### URL parameters


|Parameter|Type|Description               |
|---------|----|--------------------------|
|timeout  |time|Explicit operation timeout|


nodes.stats
-----------

Returns statistical information about nodes in the cluster.

```
GET _nodes/{node_id}/stats

```


```
GET _nodes/stats/{metric}

```


```
GET _nodes/{node_id}/stats/{metric}

```


```
GET _nodes/stats/{metric}/{index_metric}

```


```
GET _nodes/{node_id}/stats/{metric}/{index_metric}

```


#### URL parameters



* Parameter: completion_fields
  * Type: list
  * Description: A comma-separated list of fields for fielddata and suggest index metric (supports wildcards)
* Parameter: fielddata_fields
  * Type: list
  * Description: A comma-separated list of fields for fielddata index metric (supports wildcards)
* Parameter: fields
  * Type: list
  * Description: A comma-separated list of fields for fielddata and completion index metric (supports wildcards)
* Parameter: groups
  * Type: boolean
  * Description: A comma-separated list of search groups for search index metric
* Parameter: level
  * Type: enum
  * Description: Return indices stats aggregated at index, node or shard level
* Parameter: types
  * Type: list
  * Description: A comma-separated list of document types for the indexing index metric
* Parameter: timeout
  * Type: time
  * Description: Explicit operation timeout
* Parameter: include_segment_file_sizes
  * Type: boolean
  * Description: Whether to report the aggregated disk usage of each one of the Lucene index files (only applies if segment stats are requested)


nodes.usage
-----------

Returns low-level information about REST actions usage on nodes.

```
GET _nodes/{node_id}/usage

```


```
GET _nodes/usage/{metric}

```


```
GET _nodes/{node_id}/usage/{metric}

```


#### URL parameters


|Parameter|Type|Description               |
|---------|----|--------------------------|
|timeout  |time|Explicit operation timeout|


ping
----

Returns whether the cluster is running.

put\_script
-----------

Creates or updates a script.

```
PUT _scripts/{id}
POST _scripts/{id}

```


```
PUT _scripts/{id}/{context}
POST _scripts/{id}/{context}

```


#### HTTP request body

The document

**Required**: True

#### URL parameters


|Parameter     |Type  |Description                             |
|--------------|------|----------------------------------------|
|timeout       |time  |Explicit operation timeout              |
|master_timeout|time  |Specify timeout for connection to master|
|context       |string|Context name to compile script against  |


rank\_eval
----------

Allows to evaluate the quality of ranked search results over a set of typical search queries

```
GET _rank_eval
POST _rank_eval

```


```
GET {index}/_rank_eval
POST {index}/_rank_eval

```


#### HTTP request body

The ranking evaluation search definition, including search requests, document ratings and ranking metric definition.

**Required**: True

#### URL parameters



* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether specified concrete indices should be ignored when unavailable (missing or closed)
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes _all string or when no indices have been specified)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.
* Parameter: search_type
  * Type: enum
  * Description: Search operation type


reindex
-------

Allows to copy documents from one index to another, optionally filtering the source documents by a query, changing the destination index settings, or fetching the documents from a remote cluster.

#### HTTP request body

The search definition using the Query DSL and the prototype for the index request.

**Required**: True

#### URL parameters



* Parameter: refresh
  * Type: boolean
  * Description: Should the affected indexes be refreshed?
  *  :  
* Parameter: timeout
  * Type: time
  * Description: Time each individual bulk request should wait for shards that are unavailable.
  *  :  
* Parameter: wait_for_active_shards
  * Type: string
  * Description: Sets the number of shard copies that must be active before proceeding with the reindex operation. Defaults to 1, meaning the primary shard only. Set to all for all shard copies, otherwise set to any non-negative value less than or equal to the total number of copies for the shard (number of replicas + 1)
  *  :  
* Parameter: wait_for_completion
  * Type: boolean
  * Description: Should the request should block until the reindex is complete.
  *  :  
* Parameter: requests_per_second
  * Type: number
  * Description: The throttle to set on this request in sub-requests per second. -1 means no throttle.
  *  :  
* Parameter: scroll
  * Type: time
  * Description: Control how long to keep the search context alive
  *  :  
* Parameter: slices
  * Type: number
  * Description: string
  *  : The number of slices this task should be divided into. Defaults to 1, meaning the task isn’t sliced into subtasks. Can be set to auto.
* Parameter: max_docs
  * Type: number
  * Description: Maximum number of documents to process (default: all documents)
  *  :  


reindex\_rethrottle
-------------------

Changes the number of requests per second for a particular Reindex operation.

```
POST _reindex/{task_id}/_rethrottle

```


#### URL parameters



* Parameter: requests_per_second
  * Type: number
  * Description: The throttle to set on this request in floating sub-requests per second. -1 means set no throttle.


render\_search\_template
------------------------

Allows to use the Mustache language to pre-render a search definition.

```
GET _render/template
POST _render/template

```


```
GET _render/template/{id}
POST _render/template/{id}

```


#### HTTP request body

The search definition template and its params

scripts\_painless\_execute
--------------------------

Allows an arbitrary script to be executed and a result to be returned

```
GET _scripts/painless/_execute
POST _scripts/painless/_execute

```


#### HTTP request body

The script to execute

Allows to retrieve a large numbers of results from a single search request.

```
GET _search/scroll
POST _search/scroll

```


```
GET _search/scroll/{scroll_id}
POST _search/scroll/{scroll_id}

```


#### HTTP request body

The scroll ID if not passed by URL or query parameter.

#### URL parameters



* Parameter: scroll
  * Type: time
  * Description: Specify how long a consistent view of the index should be maintained for scrolled search
* Parameter: scroll_id
  * Type: string
  * Description: The scroll ID for scrolled search
* Parameter: rest_total_hits_as_int
  * Type: boolean
  * Description: Indicates whether hits.total should be rendered as an integer or an object in the rest search response


search
------

Returns results matching a query.

```
GET {index}/_search
POST {index}/_search

```


```
GET {index}/{type}/_search
POST {index}/{type}/_search

```


#### HTTP request body

The search definition using the Query DSL

#### URL parameters



* Parameter: analyzer
  * Type: string
  * Description: The analyzer to use for the query string
* Parameter: analyze_wildcard
  * Type: boolean
  * Description: Specify whether wildcard and prefix queries should be analyzed (default: false)
* Parameter: ccs_minimize_roundtrips
  * Type: boolean
  * Description: Indicates whether network round-trips should be minimized as part of cross-cluster search requests execution
* Parameter: default_operator
  * Type: enum
  * Description: The default operator for query string query (AND or OR)
* Parameter: df
  * Type: string
  * Description: The field to use as default where no field prefix is given in the query string
* Parameter: explain
  * Type: boolean
  * Description: Specify whether to return detailed information about score computation as part of a hit
* Parameter: stored_fields
  * Type: list
  * Description: A comma-separated list of stored fields to return as part of a hit
* Parameter: docvalue_fields
  * Type: list
  * Description: A comma-separated list of fields to return as the docvalue representation of a field for each hit
* Parameter: from
  * Type: number
  * Description: Starting offset (default: 0)
* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether specified concrete indices should be ignored when unavailable (missing or closed)
* Parameter: ignore_throttled
  * Type: boolean
  * Description: Whether specified concrete, expanded or aliased indices should be ignored when throttled
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes _all string or when no indices have been specified)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.
* Parameter: lenient
  * Type: boolean
  * Description: Specify whether format-based query failures (such as providing text to a numeric field) should be ignored
* Parameter: preference
  * Type: string
  * Description: Specify the node or shard the operation should be performed on (default: random)
* Parameter: q
  * Type: string
  * Description: Query in the Lucene query string syntax
* Parameter: routing
  * Type: list
  * Description: A comma-separated list of specific routing values
* Parameter: scroll
  * Type: time
  * Description: Specify how long a consistent view of the index should be maintained for scrolled search
* Parameter: search_type
  * Type: enum
  * Description: Search operation type
* Parameter: size
  * Type: number
  * Description: Number of hits to return (default: 10)
* Parameter: sort
  * Type: list
  * Description: A comma-separated list of : pairs
* Parameter: _source
  * Type: list
  * Description: True or false to return the _source field or not, or a list of fields to return
* Parameter: _source_excludes
  * Type: list
  * Description: A list of fields to exclude from the returned _source field
* Parameter: _source_includes
  * Type: list
  * Description: A list of fields to extract and return from the _source field
* Parameter: terminate_after
  * Type: number
  * Description: The maximum number of documents to collect for each shard, upon reaching which the query execution will terminate early.
* Parameter: stats
  * Type: list
  * Description: Specific ‘tag’ of the request for logging and statistical purposes
* Parameter: suggest_field
  * Type: string
  * Description: Specify which field to use for suggestions
* Parameter: suggest_mode
  * Type: enum
  * Description: Specify suggest mode
* Parameter: suggest_size
  * Type: number
  * Description: How many suggestions to return in response
* Parameter: suggest_text
  * Type: string
  * Description: The source text for which the suggestions should be returned
* Parameter: timeout
  * Type: time
  * Description: Explicit operation timeout
* Parameter: track_scores
  * Type: boolean
  * Description: Whether to calculate and return scores even if they are not used for sorting
* Parameter: track_total_hits
  * Type: boolean
  * Description: Indicate if the number of documents that match the query should be tracked
* Parameter: allow_partial_search_results
  * Type: boolean
  * Description: Indicate if an error should be returned if there is a partial search failure or timeout
* Parameter: typed_keys
  * Type: boolean
  * Description: Specify whether aggregation and suggester names should be prefixed by their respective types in the response
* Parameter: version
  * Type: boolean
  * Description: Specify whether to return document version as part of a hit
* Parameter: seq_no_primary_term
  * Type: boolean
  * Description: Specify whether to return sequence number and primary term of the last modification of each hit
* Parameter: request_cache
  * Type: boolean
  * Description: Specify if request cache should be used for this request or not, defaults to index level setting
* Parameter: batched_reduce_size
  * Type: number
  * Description: The number of shard results that should be reduced at once on the coordinating node. This value should be used as a protection mechanism to reduce the memory overhead per search request if the potential number of shards in the request can be large.
* Parameter: max_concurrent_shard_requests
  * Type: number
  * Description: The number of concurrent shard requests per node this search executes concurrently. This value should be used to limit the impact of the search on the cluster in order to limit the number of concurrent shard requests
* Parameter: pre_filter_shard_size
  * Type: number
  * Description: A threshold that enforces a pre-filter roundtrip to prefilter search shards based on query rewriting if the number of shards the search request expands to exceeds the threshold. This filter roundtrip can limit the number of shards significantly if for instance a shard can not match any documents based on its rewrite method ie. if date filters are mandatory to match but the shard bounds and the query are disjoint.
* Parameter: rest_total_hits_as_int
  * Type: boolean
  * Description: Indicates whether hits.total should be rendered as an integer or an object in the rest search response


search\_shards
--------------

Returns information about the indices and shards that a search request would be executed against.

```
GET _search_shards
POST _search_shards

```


```
GET {index}/_search_shards
POST {index}/_search_shards

```


#### URL parameters



* Parameter: preference
  * Type: string
  * Description: Specify the node or shard the operation should be performed on (default: random)
* Parameter: routing
  * Type: string
  * Description: Specific routing value
* Parameter: local
  * Type: boolean
  * Description: Return local information, do not retrieve the state from master node (default: false)
* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether specified concrete indices should be ignored when unavailable (missing or closed)
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes _all string or when no indices have been specified)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.


search\_template
----------------

Allows to use the Mustache language to pre-render a search definition.

```
GET _search/template
POST _search/template

```


```
GET {index}/_search/template
POST {index}/_search/template

```


```
GET {index}/{type}/_search/template
POST {index}/{type}/_search/template

```


#### HTTP request body

The search definition template and its params

**Required**: True

#### URL parameters



* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether specified concrete indices should be ignored when unavailable (missing or closed)
* Parameter: ignore_throttled
  * Type: boolean
  * Description: Whether specified concrete, expanded or aliased indices should be ignored when throttled
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes _all string or when no indices have been specified)
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.
* Parameter: preference
  * Type: string
  * Description: Specify the node or shard the operation should be performed on (default: random)
* Parameter: routing
  * Type: list
  * Description: A comma-separated list of specific routing values
* Parameter: scroll
  * Type: time
  * Description: Specify how long a consistent view of the index should be maintained for scrolled search
* Parameter: search_type
  * Type: enum
  * Description: Search operation type
* Parameter: explain
  * Type: boolean
  * Description: Specify whether to return detailed information about score computation as part of a hit
* Parameter: profile
  * Type: boolean
  * Description: Specify whether to profile the query execution
* Parameter: typed_keys
  * Type: boolean
  * Description: Specify whether aggregation and suggester names should be prefixed by their respective types in the response
* Parameter: rest_total_hits_as_int
  * Type: boolean
  * Description: Indicates whether hits.total should be rendered as an integer or an object in the rest search response
* Parameter: ccs_minimize_roundtrips
  * Type: boolean
  * Description: Indicates whether network round-trips should be minimized as part of cross-cluster search requests execution


snapshot.cleanup\_repository
----------------------------

Removes stale data from repository.

```
POST _snapshot/{repository}/_cleanup

```


#### URL parameters


|Parameter     |Type|Description                                             |
|--------------|----|--------------------------------------------------------|
|master_timeout|time|Explicit operation timeout for connection to master node|
|timeout       |time|Explicit operation timeout                              |


snapshot.clone
--------------

Clones indices from one snapshot into another snapshot in the same repository.

```
PUT _snapshot/{repository}/{snapshot}/_clone/{target_snapshot}

```


#### HTTP request body

The snapshot clone definition

**Required**: True

#### URL parameters


|Parameter     |Type|Description                                             |
|--------------|----|--------------------------------------------------------|
|master_timeout|time|Explicit operation timeout for connection to master node|


snapshot.create
---------------

Creates a snapshot in a repository.

```
PUT _snapshot/{repository}/{snapshot}
POST _snapshot/{repository}/{snapshot}

```


#### HTTP request body

The snapshot definition

**Required**: False

#### URL parameters



* Parameter: master_timeout
  * Type: time
  * Description: Explicit operation timeout for connection to master node
* Parameter: wait_for_completion
  * Type: boolean
  * Description: Should this request wait until the operation has completed before returning


snapshot.create\_repository
---------------------------

Creates a repository.

```
PUT _snapshot/{repository}
POST _snapshot/{repository}

```


#### HTTP request body

The repository definition

**Required**: True

#### URL parameters


|Parameter     |Type   |Description                                             |
|--------------|-------|--------------------------------------------------------|
|master_timeout|time   |Explicit operation timeout for connection to master node|
|timeout       |time   |Explicit operation timeout                              |
|verify        |boolean|Whether to verify the repository after creation         |


snapshot.delete
---------------

Deletes a snapshot.

```
DELETE _snapshot/{repository}/{snapshot}

```


#### URL parameters


|Parameter     |Type|Description                                             |
|--------------|----|--------------------------------------------------------|
|master_timeout|time|Explicit operation timeout for connection to master node|


snapshot.delete\_repository
---------------------------

Deletes a repository.

```
DELETE _snapshot/{repository}

```


#### URL parameters


|Parameter     |Type|Description                                             |
|--------------|----|--------------------------------------------------------|
|master_timeout|time|Explicit operation timeout for connection to master node|
|timeout       |time|Explicit operation timeout                              |


snapshot.get
------------

Returns information about a snapshot.

```
GET _snapshot/{repository}/{snapshot}

```


#### URL parameters



* Parameter: master_timeout
  * Type: time
  * Description: Explicit operation timeout for connection to master node
* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether to ignore unavailable snapshots, defaults to false which means a SnapshotMissingException is thrown
* Parameter: verbose
  * Type: boolean
  * Description: Whether to show verbose snapshot info or only show the basic info found in the repository index blob


snapshot.get\_repository
------------------------

Returns information about a repository.

```
GET _snapshot/{repository}

```


#### URL parameters



* Parameter: master_timeout
  * Type: time
  * Description: Explicit operation timeout for connection to master node
* Parameter: local
  * Type: boolean
  * Description: Return local information, do not retrieve the state from master node (default: false)


snapshot.restore
----------------

Restores a snapshot.

```
POST _snapshot/{repository}/{snapshot}/_restore

```


#### HTTP request body

Details of what to restore

**Required**: False

#### URL parameters



* Parameter: master_timeout
  * Type: time
  * Description: Explicit operation timeout for connection to master node
* Parameter: wait_for_completion
  * Type: boolean
  * Description: Should this request wait until the operation has completed before returning


snapshot.status
---------------

Returns information about the status of a snapshot.

```
GET _snapshot/{repository}/_status

```


```
GET _snapshot/{repository}/{snapshot}/_status

```


#### URL parameters



* Parameter: master_timeout
  * Type: time
  * Description: Explicit operation timeout for connection to master node
* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether to ignore unavailable snapshots, defaults to false which means a SnapshotMissingException is thrown


snapshot.verify\_repository
---------------------------

Verifies a repository.

```
POST _snapshot/{repository}/_verify

```


#### URL parameters


|Parameter     |Type|Description                                             |
|--------------|----|--------------------------------------------------------|
|master_timeout|time|Explicit operation timeout for connection to master node|
|timeout       |time|Explicit operation timeout                              |


tasks.cancel
------------

Cancels a task, if it can be cancelled through an API.

```
POST _tasks/{task_id}/_cancel

```


#### URL parameters



* Parameter: nodes
  * Type: list
  * Description: A comma-separated list of node IDs or names to limit the returned information; use _local to return information from the node you’re connecting to, leave empty to get information from all nodes
* Parameter: actions
  * Type: list
  * Description: A comma-separated list of actions that should be cancelled. Leave empty to cancel all.
* Parameter: parent_task_id
  * Type: string
  * Description: Cancel tasks with specified parent task id (node_id:task_number). Set to -1 to cancel all.
* Parameter: wait_for_completion
  * Type: boolean
  * Description: Should the request block until the cancellation of the task and its descendant tasks is completed. Defaults to false


tasks.get
---------

Returns information about a task.

#### URL parameters


|Parameter          |Type   |Description                                             |
|-------------------|-------|--------------------------------------------------------|
|wait_for_completion|boolean|Wait for the matching tasks to complete (default: false)|
|timeout            |time   |Explicit operation timeout                              |


tasks.list
----------

Returns a list of tasks.

#### URL parameters



* Parameter: nodes
  * Type: list
  * Description: A comma-separated list of node IDs or names to limit the returned information; use _local to return information from the node you’re connecting to, leave empty to get information from all nodes
* Parameter: actions
  * Type: list
  * Description: A comma-separated list of actions that should be returned. Leave empty to return all.
* Parameter: detailed
  * Type: boolean
  * Description: Return detailed task information (default: false)
* Parameter: parent_task_id
  * Type: string
  * Description: Return tasks with specified parent task id (node_id:task_number). Set to -1 to return all.
* Parameter: wait_for_completion
  * Type: boolean
  * Description: Wait for the matching tasks to complete (default: false)
* Parameter: group_by
  * Type: enum
  * Description: Group tasks by nodes or parent/child relationships
* Parameter: timeout
  * Type: time
  * Description: Explicit operation timeout


termvectors
-----------

Returns information and statistics about terms in the fields of a particular document.

```
GET {index}/_termvectors/{id}
POST {index}/_termvectors/{id}

```


```
GET {index}/_termvectors
POST {index}/_termvectors

```


```
GET {index}/{type}/{id}/_termvectors
POST {index}/{type}/{id}/_termvectors

```


```
GET {index}/{type}/_termvectors
POST {index}/{type}/_termvectors

```


#### HTTP request body

Define parameters and or supply a document to get termvectors for. See documentation.

**Required**: False

#### URL parameters



* Parameter: term_statistics
  * Type: boolean
  * Description: Specifies if total term frequency and document frequency should be returned.
* Parameter: field_statistics
  * Type: boolean
  * Description: Specifies if document count, sum of document frequencies and sum of total term frequencies should be returned.
* Parameter: fields
  * Type: list
  * Description: A comma-separated list of fields to return.
* Parameter: offsets
  * Type: boolean
  * Description: Specifies if term offsets should be returned.
* Parameter: positions
  * Type: boolean
  * Description: Specifies if term positions should be returned.
* Parameter: payloads
  * Type: boolean
  * Description: Specifies if term payloads should be returned.
* Parameter: preference
  * Type: string
  * Description: Specify the node or shard the operation should be performed on (default: random).
* Parameter: routing
  * Type: string
  * Description: Specific routing value.
* Parameter: realtime
  * Type: boolean
  * Description: Specifies if request is real-time as opposed to near-real-time (default: true).
* Parameter: version
  * Type: number
  * Description: Explicit version number for concurrency control
* Parameter: version_type
  * Type: enum
  * Description: Specific version type


update
------

Updates a document with a script or partial document.

```
POST {index}/_update/{id}

```


```
POST {index}/{type}/{id}/_update

```


#### HTTP request body

The request definition requires either `script` or partial `doc`

**Required**: True

#### URL parameters



* Parameter: wait_for_active_shards
  * Type: string
  * Description: Sets the number of shard copies that must be active before proceeding with the update operation. Defaults to 1, meaning the primary shard only. Set to all for all shard copies, otherwise set to any non-negative value less than or equal to the total number of copies for the shard (number of replicas + 1)
* Parameter: _source
  * Type: list
  * Description: True or false to return the _source field or not, or a list of fields to return
* Parameter: _source_excludes
  * Type: list
  * Description: A list of fields to exclude from the returned _source field
* Parameter: _source_includes
  * Type: list
  * Description: A list of fields to extract and return from the _source field
* Parameter: lang
  * Type: string
  * Description: The script language (default: painless)
* Parameter: refresh
  * Type: enum
  * Description: If true then refresh the affected shards to make this operation visible to search, if wait_for then wait for a refresh to make this operation visible to search, if false (the default) then do nothing with refreshes.
* Parameter: retry_on_conflict
  * Type: number
  * Description: Specify how many times should the operation be retried when a conflict occurs (default: 0)
* Parameter: routing
  * Type: string
  * Description: Specific routing value
* Parameter: timeout
  * Type: time
  * Description: Explicit operation timeout
* Parameter: if_seq_no
  * Type: number
  * Description: only perform the update operation if the last operation that has changed the document has the specified sequence number
* Parameter: if_primary_term
  * Type: number
  * Description: only perform the update operation if the last operation that has changed the document has the specified primary term
* Parameter: require_alias
  * Type: boolean
  * Description: When true, requires destination is an alias. Default is false


update\_by\_query
-----------------

Performs an update on every document in the index without changing the source, for example to pick up a mapping change.

```
POST {index}/_update_by_query

```


```
POST {index}/{type}/_update_by_query

```


#### HTTP request body

The search definition using the Query DSL

#### URL parameters



* Parameter: analyzer
  * Type: string
  * Description: The analyzer to use for the query string
  *  :  
* Parameter: analyze_wildcard
  * Type: boolean
  * Description: Specify whether wildcard and prefix queries should be analyzed (default: false)
  *  :  
* Parameter: default_operator
  * Type: enum
  * Description: The default operator for query string query (AND or OR)
  *  :  
* Parameter: df
  * Type: string
  * Description: The field to use as default where no field prefix is given in the query string
  *  :  
* Parameter: from
  * Type: number
  * Description: Starting offset (default: 0)
  *  :  
* Parameter: ignore_unavailable
  * Type: boolean
  * Description: Whether specified concrete indices should be ignored when unavailable (missing or closed)
  *  :  
* Parameter: allow_no_indices
  * Type: boolean
  * Description: Whether to ignore if a wildcard indices expression resolves into no concrete indices. (This includes _all string or when no indices have been specified)
  *  :  
* Parameter: conflicts
  * Type: enum
  * Description: What to do when the update by query hits version conflicts?
  *  :  
* Parameter: expand_wildcards
  * Type: enum
  * Description: Whether to expand wildcard expression to concrete indices that are open, closed or both.
  *  :  
* Parameter: lenient
  * Type: boolean
  * Description: Specify whether format-based query failures (such as providing text to a numeric field) should be ignored
  *  :  
* Parameter: pipeline
  * Type: string
  * Description: Ingest pipeline to set on index requests made by this action. (default: none)
  *  :  
* Parameter: preference
  * Type: string
  * Description: Specify the node or shard the operation should be performed on (default: random)
  *  :  
* Parameter: q
  * Type: string
  * Description: Query in the Lucene query string syntax
  *  :  
* Parameter: routing
  * Type: list
  * Description: A comma-separated list of specific routing values
  *  :  
* Parameter: scroll
  * Type: time
  * Description: Specify how long a consistent view of the index should be maintained for scrolled search
  *  :  
* Parameter: search_type
  * Type: enum
  * Description: Search operation type
  *  :  
* Parameter: search_timeout
  * Type: time
  * Description: Explicit timeout for each search request. Defaults to no timeout.
  *  :  
* Parameter: size
  * Type: number
  * Description: Deprecated, please use max_docs instead
  *  :  
* Parameter: max_docs
  * Type: number
  * Description: Maximum number of documents to process (default: all documents)
  *  :  
* Parameter: sort
  * Type: list
  * Description: A comma-separated list of : pairs
  *  :  
* Parameter: _source
  * Type: list
  * Description: True or false to return the _source field or not, or a list of fields to return
  *  :  
* Parameter: _source_excludes
  * Type: list
  * Description: A list of fields to exclude from the returned _source field
  *  :  
* Parameter: _source_includes
  * Type: list
  * Description: A list of fields to extract and return from the _source field
  *  :  
* Parameter: terminate_after
  * Type: number
  * Description: The maximum number of documents to collect for each shard, upon reaching which the query execution will terminate early.
  *  :  
* Parameter: stats
  * Type: list
  * Description: Specific ‘tag’ of the request for logging and statistical purposes
  *  :  
* Parameter: version
  * Type: boolean
  * Description: Specify whether to return document version as part of a hit
  *  :  
* Parameter: version_type
  * Type: boolean
  * Description: Should the document increment the version number (internal) on hit or not (reindex)
  *  :  
* Parameter: request_cache
  * Type: boolean
  * Description: Specify if request cache should be used for this request or not, defaults to index level setting
  *  :  
* Parameter: refresh
  * Type: boolean
  * Description: Should the affected indexes be refreshed?
  *  :  
* Parameter: timeout
  * Type: time
  * Description: Time each individual bulk request should wait for shards that are unavailable.
  *  :  
* Parameter: wait_for_active_shards
  * Type: string
  * Description: Sets the number of shard copies that must be active before proceeding with the update by query operation. Defaults to 1, meaning the primary shard only. Set to all for all shard copies, otherwise set to any non-negative value less than or equal to the total number of copies for the shard (number of replicas + 1)
  *  :  
* Parameter: scroll_size
  * Type: number
  * Description: Size on the scroll request powering the update by query
  *  :  
* Parameter: wait_for_completion
  * Type: boolean
  * Description: Should the request should block until the update by query operation is complete.
  *  :  
* Parameter: requests_per_second
  * Type: number
  * Description: The throttle to set on this request in sub-requests per second. -1 means no throttle.
  *  :  
* Parameter: slices
  * Type: number
  * Description: string
  *  : The number of slices this task should be divided into. Defaults to 1, meaning the task isn’t sliced into subtasks. Can be set to auto.


update\_by\_query\_rethrottle
-----------------------------

Changes the number of requests per second for a particular Update By Query operation.

```
POST _update_by_query/{task_id}/_rethrottle

```


#### URL parameters



* Parameter: requests_per_second
  * Type: number
  * Description: The throttle to set on this request in floating sub-requests per second. -1 means set no throttle.

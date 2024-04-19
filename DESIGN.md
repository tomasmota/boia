# Design Decisions
## Clients do not have long-lived connections. Each command is a tcp connection
- this just simplifies implementation

## Client requests only contain the command and a serial number

## keys are strings of up to 255 bits, values are 32bit integers. This is also for simplicity and can later be extended

# Algorithm Notes
## Election Restriction

A new Leader must hold all commited log entries to be elected.

## How it is achieved

- a voter will deny a vote request if its log is more up-to-date than the candidate's log
    - to this end, RequestVote RPCs send information about the candidate's log
- since entries must be in a majority of the cluster to be commited, by requiring a vote by a majority we guarantee that if the candidate is not up-to-date the majority will deny the vote
- to compare two logs, compare index and term of last entry of each.
  - if they have different terms, later term is more up-to-date
  - if they have same term, higher index (longer log) is more up-to-date


## Client Interaction

Clients connect to a random server. if it is a follower, it will send information about the most recent leader 

- AppendEntries requests include the network address of the leader  

In order to prevent the same command being executed twice on client retry, every command should include a serial number

Read operations should never return stale data, so there are the following restrictions
- Right after being elected, the leader commits a blank no-op entry in the log, in order to figure out which of it's own log entries have been commited
    - note that even though the leader is guaranteed to have all commited entries, it might not know which ones actually have been commited
- Before responding to read requests, the leader must exchange heart-beat with majority of the cluster, to ensure it has not been deposed

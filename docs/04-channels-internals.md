# Channel Internals: hchan Under the Hood

## hchan Structure

Every channel is a `runtime.hchan` struct:

```
type hchan struct {
    qcount   uint      // current elements in buffer
    dataqsiz uint      // buffer capacity (0 for unbuffered)
    buf      *[...]T   // ring buffer (circular queue)
    elemsize uint16    // element size
    closed   uint32    // closed flag
    sendx    uint      // send index in ring buffer
    recvx    uint      // receive index in ring buffer
    recvq    waitq     // list of blocked receivers (sudog linked list)
    sendq    waitq     // list of blocked senders (sudog linked list)
    lock     mutex     // protects all fields
}
```

## Ring Buffer (Buffered Channels)

```
buf: [  A  |  B  |  C  |  _  |  _  ]
           ↑              ↑
         recvx          sendx

Capacity: 5, Count: 3
Send → writes at sendx, advances sendx
Recv → reads at recvx, advances recvx
When sendx == recvx and qcount == dataqsiz → full (sender blocks)
```

## Send/Receive Mechanics

**Send (`ch <- v`)**:
1. Lock hchan
2. If receiver waiting in recvq → copy directly to receiver's stack (no buffer)
3. Else if buffer has space → copy to buf[sendx], advance sendx
4. Else → park sender in sendq, suspend goroutine

**Receive (`v := <-ch`)**:
1. Lock hchan
2. If sender waiting in sendq → receive directly from sender
3. Else if buffer has data → copy from buf[recvx], advance recvx
4. Else → park receiver in recvq, suspend goroutine

**Key optimization**: Direct send (step 2 of send) copies value directly
to the receiver's stack memory, bypassing the buffer entirely.

## Nil and Closed Channel Behavior

| Operation | Nil Channel | Closed Channel |
|-----------|-------------|----------------|
| `ch <- v` | Blocks forever | **panic** |
| `<-ch` | Blocks forever | Returns zero value (ok=false) |
| `close(ch)` | **panic** | **panic** |
| `len(ch)` | 0 | buffered count |
| `cap(ch)` | 0 | buffer capacity |

## Select Statement Internals

1. All cases are shuffled (randomized) to prevent starvation
2. All channels are locked in address order (prevents deadlock)
3. Check each case for readiness
4. If none ready + default exists → take default
5. If none ready + no default → park goroutine in ALL channel waitqueues
6. When any channel becomes ready → wake goroutine, dequeue from all others

## Performance Characteristics

| Channel Type | Allocation | Best For |
|-------------|------------|----------|
| Unbuffered (`make(chan T)`) | Just hchan | Synchronization, handoff |
| Small buffer (`make(chan T, 1)`) | hchan + buf | Signaling, latest value |
| Large buffer (`make(chan T, N)`) | hchan + buf | Producer-consumer decoupling |

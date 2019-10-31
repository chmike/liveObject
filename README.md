Live object termination
-----------------------

This code is a reference implementation showing how we can terminate the 
goroutine associated to a live object when the live object is garbage collected.

This is achieved with following key steps :

- encapsulate the fields of the live object accessed by the goroutine in a sub-object 
  referenced by an embedded pointer
- instantiate and return the live object 
  - initialize the fields of the live object and the sub-object
  - set a finalizer on the live object to request a termination of the goroutine when
    the live object is garbage collected
  - start the goroutine referencing and accessing only the sub-object and
    terminating when requested by the finalizer

The mecanism used to request the termination of the goroutine is a simple bool in this
reference implementation. It would be a closing channel if the goroutine is looping and 
blocking over a select instruction.

## Working principle

Automatic termination of a goroutine associated to a live object can be achieved by setting 
a finalizer. The finalizer will signal to the goroutine to terminate when the live object 
is garbage collected. 

Unfortunately, the live object won’t be garbage collected if the goroutine holds a reference 
to it. For this reason, we use a sub-object that is referenced by the live object and the
goroutine. Only the user will reference the live object. Once all references to it are
removed, the finalizer will signal the goroutine to terminate and the live object is garbage
collected. Once the goroutine terminates, the last reference to the sub-object is also removed 
and it is garbage collected as well. 

This reference implementation embeds a pointer to the sub-object in the live object. This
allows to access the fields of the sub-object without indirection.
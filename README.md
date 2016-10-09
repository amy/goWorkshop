Day 1: Data oriented design + mechanical sympathy
https://github.com/ardanlabs/gotraining/blob/master/courses/README.md

Language
	How to reduce pressure on garbage collection --> increases performance
	Pressure = Values that end up on the heap & their existence is short period of time
			   Want long living values. 

	///////////////////
	//// VARIABLES ////
	///////////////////
		String:
			> two word data structure
			> word 1 --> pointer to backing array
			> word 2 --> num of bytes in backing array
		Zero Value: 
			> Every variable initialized to zero value
			> Why? Because in C, there's been too many bugs that arise because
			  you don't initialize the memory 
			> Go prefers integrity first
		Conversion:
			> Go does NOT have CASTING
			> Go has conversion
				- This means that Go needs to allocate more memory to accommodate
				  the extra bytes
		Type:
			> alignment --> compiler does this for us. Requirement from the
							hardware. Alignment boundaries. (2 vs. 4 byte alignment)
							-------- --------
						    |      | |      |   What happens if you have data across
						    -------- --------   two words? Compensate with PADDING. 

							type example struct {
								flag	bool		// needs to fall on multiple of 2
								counter int16		// 					multiple of 8
								pi 		float64
							}

							- Depending on the types, you need an increase of padding
							- ONLY MATTERS IF the struct represents pure data where
							  you need a lot of this struct 

							Fix by organizing struct largest field on top

							type example struct {
								pi 		float32
								counter int16
								flag 	bool
							}

			> anonymous struct type --> don't declare a named struct type if you
										if you don't need to. If its only used
										once, use anonymous. Ex: unmarshalling JSON. 
					- you can assign anonymous type to named type
					- working on named types, must be explicit with conversion 

	////////////////////
	///// POINTERS /////
	////////////////////

	Every go routine has its own stack. If you want to share things across program bounderies, 
	you need to put the data into the heap.

	Call main()
	Initial size of stack is 2K. 
	Every function allocates frame on stack. Go stack grows down. 
	size of function stack is set at compile time
	only values where size is known on compile time, can be added on the stack
	stack is free memory --> garbage collector doesn't care. It is self cleaning
	Go routine owns stack 
	Leverage the stack as much as possible. Because you want to reduce pressure on the heap. 

	Pass By Value:
		> pass by value, you don't share the memory
	Stacks are self cleaning:
		> stacks clean on the way down when function frames need to be allocated
		> clean every time you make a function call
	Pass By Value --> pass the address of the value
		> pointer variable hold an address
		> point of sharing memory is to read & write to memory type 
		> *int --> only accept addresses associated to integer values 

		func increment(inc *int) {	// * is not an operator. It is bound to a type star Int. 
			*int++		// * is an operator here. Access the value at the address. 
		}

	Stay on Stack
		--------
		| user |  main()
		--------
		--------
		| user |  stayOnStack()
		--------
	Escape to Heap --> compiler has ESCAPE ANALYSIS
		---------
		|   *   |  main()
		---------							**HEAP**
		---------							---------
		| &user |  escapeToHeap() -----> 	|  user |  assign value to heap
		---------							---------
	go build -gcflags -m 	--> escape analysis 

	Escaping:
		> share things up = escape

	Cannot share values across stacks / go routines. Need to put it on the heap. 
	When you run out of space on a stack, you copy the stack * function frames over to readjust.

	HEAP:
		algorithm to determine garbage collection 
		Live heap = how many values in heap
		Heap size = max # of values in heap
		Garbage collector has its own set of go routines to do the work
			> takes 25% of computing power to do the work
	GC (Garbage Collector):
		Write Barrier --> tell the GC what you're doing on the heap while its running w/ coloring
		Memory leak in Go --> reference to a value on a stack, and the reference never goes away
							  reference on the stack means, the garbage collector can't clean it up
							  After a function terminates, the memory never is cleaned up. 
							  Ex: forgetting to close things, never deleting keys in a map
	**GOOGLE HOW TO** Optimize against garbage collector

	///////////////////
	/// READABILITY ///
	///////////////////

	DO NOT START OUT THE LIFE OF A VALUE WITH A POINTER

		var u user
		err = json.Unmarshal([]byte(r), &u)
		return &u, err	// escapes to heap. Shows that it is being shared.

		*****DO NOT START OUT THE LIFE OF A VALUE WITH A POINTER*****

		var u *user
		err = json.Unmarshal([]byte(r), &u)
		return u, err	// unmarshal has a pointer to my pointer (*user)
						// Walk away from readability. Does the same thing. 

	If statements
		Put negative path in if statements. (If something goes wrong)

	Variable declaration
		The closer the variable is declared to where it's being used, the shorter
		the name of the variable. The farther you declare the variable, the longer
		because you need more context.

	////////////////
	// Benchmarks //
	////////////////
	Benchmark test measures the performance of a function / object

	Break up Main Memory into "cache lines":
		> GOAL: want to try for linear array traversal 
				hardware loves array.
				Go loves slices 
		> pull out full cache line out of main memory into L3 
		> write predictable access patterns to memory
		> Predictable Access Patterns
			1. Linear traversal through memory. Group / allocate data as continuously as possible 
			   and iterate through it continuously. Prefetchers can pick up on it. Array most important in relation to 
			2. Striding; Table look aside buffer cache. Maintains relationship OS pages & physical
			   memory addresses. Page = OS virtual granularity.
		> Worst case: column traversal (TLB & cache line miss)
					  Matrix is so big, that its not on the cache line AND its not on the same TLB
		> Linked List: sitting on the same page but not same cache line 
		> TLB miss is worse
		> Data access affects efficiency

	////////////////////////////
	// Data Access Efficiency //	DATA ORIENTED DESIGN 
	////////////////////////////

	Go loves slices 

	Object oriented patterns create Linked Lists & its not sympathetic to data access in Go specifically.

	Array:
		string --> is a 2 word data structure 
		-----
		| * |	pointer to backing array ---> |A|p|p|l|e|
		-----
		| 5 |	# of bytes
		-----

		[4]int --> size of array is PART OF THE TYPE
		[5]int
		Since array is known at compile time you can put the array onto the stack.
		Array gives you continguous block of memory. 
		Arrays are static in size. 
	
	Slice:
		Reference type --> values stay on stack and point to things on heap

		slice := make([]string, 5)
		slice is 3 word data value

						   0 ....	   4
		-----			-------------------
		| * |	---> 	|  *  | nil	| nil |		--> pointers in here point to backing arrays
		-----			-------------------
		| 	|			|  5  |  0	|  0  |
		| 	|			-------------------
		| 5	|   ---> 	length (what you have access to)
		-----
		| 5	|   ---> 	capacity (total # elements that exist) 
		-----

		*
		5	--> if you try to access slice index 5, panic. Out of bounds. 
		8	--> slice initialized to size 8. 6 - 7 not accessible. 

		// declare nil slice 		data initialized to:
		var data []string 				-------
										| nil |
										-------
										|  0  |
										-------
										|  0  |
										-------

		pass slice HEADER in, get slice HEADER out. 
		data = append(data, "THING")
		  ^
		  USE THE SAME VARIABLE

		If its different variables, there's a memory leak because theres still a reference to old backing array.

		make([]string, 0, 10000) --> set the capacity to how much data you expect to append

		Backing array lives in the heap. Append algorithm doubles backing array size until 1000. 
		Then backing array increase backing array size by about 25%. 

		// declare empty slice 		slice is: 
		data := []string{}				-------
		Use this when you 				|  *  |	--> point to nil backing data
		unmarshal things				-------
										|  0  |
										-------
										|  0  |
										-------

		// [a:b] --> a to b not including b OR a to a + length of slice2
		slice2 := slice1[2:4:4]
		// slice2 and slice1 share the same backing array.
		// Once you excede the capacity on an append, you detach and get a new backing array

		BE CAREFUL
		// Declare a slice of integers with 7 values.
		x := []int{100, 200, 300, 400, 500, 600, 700}

		// Set a pointer to the second element of the slice.
		twohundred := &x[1]

		// Append a new value to the slice.
		x = append(x, 800)

		// Change the value of the second element of the slice.
		x[1] = 250
		CHANGES OLD BACKING ARRAY

		HOW TO KEEP THE BACKING ARRAY ON THE STACK 
		var b [10]int
		slice := b[:]
		// points to b which can be on the stack
		// things on stack must be of predefined size 

Day 2
	Decouple code from change

	/////////////
	// Methods //
	/////////////

	Emphasize consistency whether you use a value or pointer receiver 

		d.displayName()	// value receiver
		d.setAge(45)	// pointer receiver

		// What's happening underneath
		data.displayName(d)
		(*data).setAge(&d, 45)

	Put initializers / factory funcitons into the same file as the struct 
	Type, factory functions, methods.

		// functions are a reference type
		// f1 is a reference type
		f1 := d.displayName 

		 f1
		-----			-----
		| * | ----> 	| * | <--- Code
		-----			-----
						| * | ---> -----  OR -----
						-----	   | d |	 | d |
								   -----	 -----
								   copy		 original
								   			 Potential to escape because now data
								   			 is shared  
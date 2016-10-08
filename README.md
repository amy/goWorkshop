Day 1: Data oriented design + mechanical sympathy

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
	Go routine owns slack 
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

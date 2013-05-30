#!/usr/bin/env python

# This will create golden files in a directory passed to it.
# A Test calls this internally to create the golden files
# So it can process them (so we don't have to checkin the files).

import msgpack, sys, os

def get_test_data_list():
    # get list with all primitive types, and a combo type
    l = [ 
        -8,
         -1616,
         -32323232,
         -6464646464646464,
         192,
         1616,
         32323232,
         6464646464646464,
         192,
         -3232.0,
         -6464646464.0,
         3232.0,
         6464646464.0,
         False,
         True,
         None,
         1328148122000002,
         "someday",
         "",
         "bytestring",
         [ 
            -8,
             -1616,
             -32323232,
             -6464646464646464,
             192,
             1616,
             32323232,
             6464646464646464,
             192,
             -3232.0,
             -6464646464.0,
             3232.0,
             6464646464.0,
             False,
             True,
             None,
             1328148122000002,
             "someday",
             "",
             "bytestring" 
             ],
         { "true": True,
           "false": False },
         { "true": "True",
           "false": False,
           "uint16(1616)": 1616 },
         { "list": [1616, 32323232, True, -3232.0, {"TRUE":True, "FALSE":False}, [True, False] ],
           "int32":32323232, "bool": True, 
           "LONG STRING": "123456789012345678901234567890123456789012345678901234567890",
           "SHORT STRING": "1234567890" },	
	 { True: "true", 8: False, "false": 0 }
         ]
    return l

def build_test_data(destdir):
    l = get_test_data_list()
    for i in range(len(l)):
        packer = msgpack.Packer()
        serialized = packer.pack(l[i])
        f = open(os.path.join(destdir, str(i) + '.golden'), 'wb')
        f.write(serialized)
        f.close()

def doMain(args):
    if len(args) == 2 and args[0] == "testdata":
        build_test_data(args[1])
    else:
        print("Usage: build.py [testdata]")
    
if __name__ == "__main__":
    doMain(sys.argv[1:])

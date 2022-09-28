syntax = "proto3";

package model;
option go_package = "/model";

// the Main Components
// only have one, include some child componetns
// all the message in this struct is the componetns
// @gogs:Components
message Components {
    // don't care the filed name, we never use it
    // but you should be careful about the filed number
    BaseWorld BaseWorld = 1;
}

// componetns, 
// all the messages are used for communication between the client and the server
message BaseWorld {
    // don't care the filed name, we never use it
    // but you should be careful about the filed number
    BindUser BindUser = 1;

    BindSuccess BindSuccess = 2;
}
// common message
// the message is used for the client and the server to communicate
// the corresponding method will be generated according to this message
message BindUser {
    string uid = 1;
}


// it is only used for messages sent from the server to the client
// the server will not receive the message and will not generate the corresponding method
// @gogs:ServerMessage
message BindSuccess {
}
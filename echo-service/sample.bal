import ballerina/http;
import ballerina/io;

service / on new http:Listener(8090) {
    resource function post .(@http:Header {name: "x-request-id"} string? xRequestId, @http:Payload string textMsg) returns string {
        if xRequestId is string {
            io:println("x-request-id: ", xRequestId);
        }
        return textMsg;
    }
}

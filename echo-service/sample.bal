import ballerina/http;
import ballerina/io;

configurable int port = 3000;

service / on new http:Listener(port) {
    resource function post .(http:Request req, @http:Header {name: "x-request-id"} string? xRequestId, @http:Payload string textMsg) returns json|error {
        if xRequestId is string {
            io:println("x-request-id: ", xRequestId);
        }

        map<string[]> headers = {};
        foreach string headerName in req.getHeaderNames() {
            headers[headerName] = check req.getHeaders(headerName);
        }

        return {
            headers: headers,
            payload: textMsg
        };
    }
}

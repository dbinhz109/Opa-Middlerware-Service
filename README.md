# Opa-Middlerware-Service

# Run golang service:
go run ./cmd/go-app/

# Run Open policy agent:
./opa run --server --set "decision_logs.console=true"

# PUT rego policy to OPA Server:

curl --location --request PUT '127.0.0.1:8181/v1/policies/example' \
--header 'Content-Type: text/plain' \
--data 'package authz

default allow = false

allow {
    input.method == "GET"
    input.path == ["lancs","test"]
}

allow {
    input.method == "POST"
    input.path == ["lancs","test"]
}

allow {
    input.method == "PUT"
    input.path == ["lancs","test"]
    input.role == "admin"
}
'
# Run test:
curl --location '127.0.0.1:9999/lancs/test' #Allow
curl --location --request POST '127.0.0.1:9999/lancs/test' #Allow
curl --location --request PUT '127.0.0.1:9999/lancs/test' --header 'X-User: admin' #Allow
curl --location --request PUT '127.0.0.1:9999/lancs/test' #Unauthorized

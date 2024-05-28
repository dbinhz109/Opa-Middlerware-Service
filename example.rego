package authz

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

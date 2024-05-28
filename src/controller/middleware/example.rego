package example

default allow = false

allow {
    input.method == "POST"
    input.path == "/lancs"
}

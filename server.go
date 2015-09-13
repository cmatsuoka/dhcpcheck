package main

import (
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"
)

func update(w http.ResponseWriter, r *http.Request) {
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")

	for i := 0; ; i++ {
		io.WriteString(w, "data: "+strconv.Itoa(i)+"\n\n")
		f.Flush()
		time.Sleep(500 * time.Millisecond)
	}
}

func status(w http.ResponseWriter, r *http.Request) {
	const page = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
		<script>
    var source = new EventSource("/update/");
    source.addEventListener("message", function(e) {
        document.getElementById("result").innerHTML = event.data;
    }, false);
		</script>

		<h1>{{.Header}}</h1>
		SSE test: <span id="result"></span>
	</body>
</html>`

	t, err := template.New("webpage").Parse(page)

	data := struct {
		Title  string
		Header string
	}{
		Title:  "DHCPCheck status",
		Header: "DHCPCheck " + Version,
	}

	err = t.Execute(w, data)
	if err != nil {
		io.WriteString(w, "Internal server error")
	}
}

func serve(port int) {
	http.HandleFunc("/", status)
	http.HandleFunc("/update/", update)
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

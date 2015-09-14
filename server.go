package main

import (
	"html/template"
	"io"
	"net/http"
	"strconv"
)

func update(w http.ResponseWriter, r *http.Request) {
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")

	for {
		s := <-repch
		io.WriteString(w, "data: "+s+"\n\n")
		f.Flush()
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
	stats=JSON.parse(event.data);
	document.getElementById("packets").innerHTML=stats.Packets;
        document.getElementById("msgtype").innerHTML = "";
	for (var key in stats.MsgType) {
        	document.getElementById("msgtype").innerHTML += key+":"+stats.MsgType[key]+"<br>";
	}
        document.getElementById("vendors").innerHTML = "";
	for (var key in stats.Vendors) {
        	document.getElementById("vendors").innerHTML += key+":"+stats.Vendors[key]+"<br>";
	}
        document.getElementById("vdclass").innerHTML = "";
	for (var key in stats.VdClass) {
        	document.getElementById("vdclass").innerHTML += key+":"+stats.VdClass[key]+"<br>";
	}
    }, false);
		</script>

		<h1>{{.Header}}</h1>
		SSE test: <span id="result"></span>
		<p>
		Packets: <span id="packets">0</span>
		<h2>DHCP message types</h2>
		<div id="msgtype">No packets received.</div>
		<h2>Packets by vendor</h2>
		<div id="vendors">No packets received.</div>
		<h2>Packets by vendor class</h2>
		<div id="vdclass">No packets received.</div>
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

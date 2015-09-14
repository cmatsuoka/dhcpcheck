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
    function showSrv(id,map) {
	var s="</td><td>"
        var t="<table><thead><tr><td>Server IP"+s+"Name"+s+"Offers"+s+"ACKs"+s+"NACKs</td></tr></thead><tbody>";
	for (var key in map) {
		var v=map[key]
        	t+="<tr><td>"+key+s+v.Name+s+v.Offer+s+v.Ack+s+v.Nack+"</td></tr>";
	}
        t+="</tbody></table>";
        document.getElementById(id).innerHTML=t;
    }
    function showMap(id,map,head) {
        var t="<table><thead><tr><td>"+head+"</td><td>Packets</td></tr></thead><tbody>";
	for (var key in map) {
        	t+="<tr><td>"+key+"</td><td>"+map[key]+"</td></tr>";
	}
        t+="</tbody></table>";
        document.getElementById(id).innerHTML=t;
    }
    var source = new EventSource("/update/");
    source.addEventListener("message", function(e) {
	//document.getElementById("result").innerHTML = event.data;
	stats=JSON.parse(event.data);
	document.getElementById("packets").innerHTML=stats.Packets;
	showSrv("servers", stats.Servers)
	showMap("msgtype", stats.MsgType, "Message type")
	showMap("vendors", stats.Vendors, "Vendor")
	showMap("vdclass", stats.VdClass, "Vendor class")
    }, false);
		</script>

		<h1>{{.Header}}</h1>
		<!-- SSE test: <span id="result"></span> -->
		<p>
		Packets: <span id="packets">0</span>
		<h2>DHCP servers</h2>
		<div id="servers">No packets received.</div>
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

// © donuts-are-good & Moritz Poldrack
// SPDX-License-Identifier: BSD-4-Clause

function openFileDialog() {
	document.getElementById("fileTree").showModal();
	updateRoots();
}

function updateRoots() {
	fetch('http://localhost:21558/files')
		.then(r =>  r.json().then(data => ({status: r.status, body: data})))
		.then(obj => {
			let sel = document.getElementById("rootSelect");
			sel.innerHTML = "";

			let index = 0;
			obj.body.forEach((elem)=>{
				let opt = document.createElement('option');
				opt.innerText = elem;
				index++;
				sel.add(opt);
			});

			listDirectory(document.getElementById("rootSelect").selectedIndex, "");
		})
		.catch((error)=>{
			console.log(error);
		})
}

function listDirectory(root, relpath) {
	document.getElementById('filelist').classList.add("hide");
	document.getElementById('listloader').classList.remove("hide");

	if(relpath!=""){
		document.getElementById("rootSelect").options[root].innerText + "/"+relpath;
	}
	fetch('http://localhost:21558/files/'+root+'?relpath='+encodeURI(relpath))
		.then(r =>  r.json().then(data => ({status: r.status, body: data})))
		.then(obj => {
			let table = document.getElementById("filelist");
			table.innerHTML = "";

			let directories = [];
			let files = [];

			obj.body.content.forEach((file)=>{
				let elem = document.createElement('a');

				elem.innerText = file.name + (file.directory?"/":"");
				if(file.directory){
					let path = relpath;
					if(path !== "") {
						path += "/";
					}
					path += file.name;
					elem.href = "javascript:listDirectory("+root+",'"+path+"');"
					elem.classList.add("fw-bold");
					directories.push(elem);
				}else{
					elem.href = "javascript:playFile("+root+",'"+relpath+"/"+file.name+"');"
					files.push(elem);
				}
			});

			let body = document.createElement("tbody")

			if(relpath != "") {
				let tr = document.createElement("tr");
				let td = document.createElement("td")
				let a = document.createElement('a');

				let path = "";

				parts = relpath.split("/");
				if(parts.length > 1){
					path = parts.slice(0, parts.length-1).join("/");
				}

				a.innerText = "../";
				a.classList.add("fw-bold");
				a.href = "javascript:listDirectory("+root+",'"+path+"');";

				td.innerHTML = a.outerHTML;
				tr.innerHTML = td.outerHTML;
				body.innerHTML += tr.outerHTML;
			}

			if(files.length == 0 && directories.length == 0){
				let tr = document.createElement("tr");
				let td = document.createElement("td")

				td.innerHTML = "<i>Wow! There's tons of nothing in here!</i>";
				tr.innerHTML = td.outerHTML;
				body.innerHTML += tr.outerHTML;
			}
			directories.forEach((elem)=>{
				let tr = document.createElement("tr");
				let td = document.createElement("td")

				td.innerHTML = elem.outerHTML;
				tr.innerHTML = td.outerHTML;
				body.innerHTML += tr.outerHTML;
			});
			files.forEach((elem)=>{
				let tr = document.createElement("tr");
				let td = document.createElement("td")

				td.innerHTML = elem.outerHTML;
				tr.innerHTML = td.outerHTML;
				body.innerHTML += tr.outerHTML;
			});

			table.innerHTML = body.outerHTML;

			document.getElementById('listloader').classList.add("hide");
			document.getElementById('filelist').classList.remove("hide");
		})
		.catch((error)=>{
			console.log(error);
		})
}

function playFile(root, path) {
	fetch('http://localhost:21558/player/start',
		{
			method: "POST",
			body: JSON.stringify({
				root: root,
				relativePath: path
			}),
			headers: {
				"Content-type": "application/json; charset=UTF-8"
			}
		}	).then(r => {
			console.log(r);
		})
		.catch((error)=>{
			console.log(error);
		})
}

function showConnected() {
	document.getElementById("connected").classList.add("show");
	setTimeout(function(){document.getElementById("connected").classList.remove("show")},3000);
}

function checkServer() {
	// Make a GET request to your API endpoint
	fetch('http://localhost:21558/status')
		.then(response => {
			if (response.status === 200) {
				if(
					document.getElementById("disconnectError").classList.contains("show")||
					document.getElementById("degradedWarn").classList.contains("show")
				){
					showConnected();
					unlock()
				}
				document.getElementById("disconnectError").classList.remove("show");
				document.getElementById("degradedWarn").classList.remove("show");
			}else{
				document.getElementById("disconnectError").classList.remove("show");
				document.getElementById("degradedWarn").classList.add("show");
				lock();
			}
		})
		.catch(()=>{
			document.getElementById("disconnectError").classList.add("show");
			document.getElementById("degradedWarn").classList.remove("show");
			lock();
		})
}

function lock() {
	document.getElementById("videoSelect").disabled = true;
	document.getElementById("rootSelect").disabled = true;
	document.getElementById("listloader").classList.add("hide");
	document.getElementById("filelist").classList.add("hide");
	document.getElementById("selectError").classList.remove("hide");
}

function unlock() {
	document.getElementById("videoSelect").disabled = false;
	document.getElementById("rootSelect").disabled = false;
	document.getElementById("selectError").classList.add("hide");
	document.getElementById("fileTree").close();
}

window.onload = function() {
	setInterval(checkServer, 1000);
	document.getElementById('rootSelect').addEventListener('select', () => {listDirectory(document.getElementById("rootSelect").selectedIndex, ".")});
}

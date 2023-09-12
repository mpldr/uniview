function createStatusCircle(color) {
	let circle = document.createElement('div');
	circle.className = 'status-circle';
	circle.style.width = '20px';
	circle.style.height = '20px';
	circle.style.borderRadius = '50%';
	circle.style.position = 'absolute';
	circle.style.right = '10px';
	circle.style.bottom = '10px';
	circle.style.backgroundColor = color;
	return circle;
}

function checkServer() {
	// Make a GET request to your API endpoint
	fetch('http://localhost:21558/status')
		.then(response => {
			if (response.status === 200) {
				// Resource is available, redirect to a new page
				window.location.href = '/room';
			}
		})
		.catch()
}

function showError() {
	let circle = createStatusCircle('#ee5253');
	document.body.appendChild(circle);
	document.getElementById("connection-dialog").showModal();
}

function unlock() {
	document.getElementById('connectButton').disabled = false;
	document.getElementById('retryConnection').disabled = false;
	document.getElementById("connectButton").children[0].classList.add("hide");
	document.getElementById("retryConnection").children[0].classList.add("hide");
}

function lock() {
	document.getElementById('connectButton').disabled = true;
	document.getElementById('retryConnection').disabled = true;
	document.getElementById("connectButton").children[0].classList.remove("hide");
	document.getElementById("retryConnection").children[0].classList.remove("hide");
}

window.onload = function() {
	let connTest = (total, delay) => {
		lock();
		// Perform resource check at regular intervals
		const interval = setInterval(checkServer, delay);

		// Stop checking after 3 seconds (3000 milliseconds)
		setTimeout(() => {
			showError();
			unlock();
		}, total);
	}
	document.getElementById('connectButton').addEventListener('click', () => {connTest(3000, 200);});

	document.getElementById('retryConnection').addEventListener('click', () => {connTest(3000, 200);});
}

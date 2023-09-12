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
    fetch('http://localhost:21588')
        .then(response => {
            let circleColor = response.status === 200 ? '#1dd1a1' : '#ee5253';
            let circle = createStatusCircle(circleColor);
            document.body.appendChild(circle);

            if (response.status !== 200) {
                let errorModal = new bootstrap.Modal(document.getElementById('errorModal'), {backdrop: false});
                errorModal.show();
            }
        })
        .catch(error => {
            let circle = createStatusCircle('#ee5253');
            document.body.appendChild(circle);
            let errorModal = new bootstrap.Modal(document.getElementById('errorModal'), {backdrop: false});
            errorModal.show();
        });
}

window.onload = function() {
    checkServer();

    document.getElementById('retryConnection').addEventListener('click', function() {
        let existingCircle = document.querySelector('.status-circle');
        if (existingCircle) existingCircle.remove();

        let errorModal = new bootstrap.Modal(document.getElementById('errorModal'), {backdrop: false});
        errorModal.hide();

        setTimeout(() => {
            checkServer();
        }, 200);
    });
}

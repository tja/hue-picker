// Send color to API
function sendColor(e) {
    // Scale to x = [-1..1] and y = [-1..1]
    const rc = e.target.getBoundingClientRect();

    const x = 2 * (e.clientX - rc.left) / (rc.right - rc.left) - 1;
    const y = 1 - 2 * (e.clientY - rc.top) / (rc.bottom - rc.top);

    // Hue is counter-clockwise angle from top
    var a = ((Math.PI / 2.0 - Math.atan2(y, x)) / Math.PI * 180.0);
    if (a < 0.0) { a += 360.0; }

    const hue = Math.round( Math.min(Math.max(a * 65536.0 / 360.0, 0), 65535) )

    // Sat is distance to origin
    const d = Math.min(1.0, Math.sqrt(x * x + y * y));

    const sat = Math.round( Math.min(Math.max(d * 256.0 , 0), 255) )

    // Send to API
    var xhr = new XMLHttpRequest();

    xhr.open("POST", "/api/light", true);
    xhr.setRequestHeader('Content-Type', 'application/json');

    xhr.send(JSON.stringify({
        on:    true,
        hue:   hue,
        sat:   sat,
    }));
}

// Set up the color wheel
function startup() {
    const el = document.getElementById('color-wheel');

    el.addEventListener('touchstart', sendColor);
    el.addEventListener('touchmove', sendColor);
    el.addEventListener('click', sendColor);
}

// Call startup when document is loaded
document.addEventListener('DOMContentLoaded', startup);

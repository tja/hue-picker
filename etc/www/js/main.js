// Global "queue"
var queue = null;

// Push color to "queue"
function pushColor(e) {
    // Don't report down-stream
    e.preventDefault();

    // Scale to x = [-1..1] and y = [-1..1]
    const rc = e.target.getBoundingClientRect();

    const cx = e.clientX || e.touches[0].clientX;
    const cy = e.clientY || e.touches[0].clientY;

    const x = 2 * (cx - rc.left) / (rc.right - rc.left) - 1;
    const y = 1 - 2 * (cy - rc.top) / (rc.bottom - rc.top);

    // Hue is counter-clockwise angle from top
    var a = ((Math.PI / 2.0 - Math.atan2(y, x)) / Math.PI * 180.0);
    if (a < 0.0) { a += 360.0; }

    const hue = Math.round( Math.min(Math.max(a * 65536.0 / 360.0, 0), 65535) )

    // Sat is distance to origin
    const d = Math.min(1.0, Math.sqrt(x * x + y * y));

    const sat = Math.round( Math.min(Math.max(d * 256.0 , 0), 255) )

    // Push on "queue"
    queue = {
        on:    true,
        hue:   hue,
        sat:   sat,
    }
}

// ...
function processQueue() {
    // Bail out if nothing to send
    if (queue == null) {
        return
    }

    // Send to API
    var xhr = new XMLHttpRequest();

    xhr.open("POST", "/api/light", true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.send(JSON.stringify(queue));

    queue = null
}

// Call startup when document is loaded
document.addEventListener('DOMContentLoaded', function() {
    // Set up color wheel
    const el = document.getElementById('color-wheel');

    el.addEventListener('touchstart', pushColor, false);
    el.addEventListener('touchmove', pushColor, false);
    el.addEventListener('click', pushColor, false);

    // Never sent more than every 500ms
    setInterval(processQueue, 500);
});

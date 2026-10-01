<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Event Explorer | Error</title>
    <link rel="stylesheet" href="../static/css/index.css">
    <link rel="stylesheet" href="../static/css/event.css">
</head>
<body>
    <nav>
        <section class="nav-left-section">
            <div class="nav-logo-box">e.</div>
            <p class="company-banner">eventexplorer</p>
        </section>
    </nav>

    <main class="main-layout">
        <section class="error-section">
            {{with index . "error"}}
            <p class="error-status">{{.Status}}</p>
            <h1 class="error-title">{{.Title}}</h1>
            <p class="error-message">{{.Message}}</p>
            {{end}}
            <a class="event-detail-button error-button" href="/">
                <span>Back to search</span>
            </a>
        </section>
    </main>
    <script src="../static/js/index.js"></script>
</body>
</html>
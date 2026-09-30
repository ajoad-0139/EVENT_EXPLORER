<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Event Explorer</title>
    <link rel="stylesheet" href="../static/css/index.css">
    <link rel="stylesheet" href="../static/css/lists.css">
</head>
<body>
    <nav>
        <section class="nav-left-section">
            <div class="nav-logo-box">e.</div>
            <p class="company-banner">eventexplorer</p>
        </section>
    </nav>

    <!-- main section starts here  -->
    <main class="main-layout">

        <!-- city header section -->
        <section class="city-header-section">
            <section class="city-header-left-section">
                <p class="city-header-eyebrow">YOUR CITY. YOUR NEXT PLAN.</p>
                <h1 class="city-header-title">
                    What is on in
                    <span class="city-header-title-green" id="city-name">Toronto.</span>
                </h1>
                <p class="city-header-subtitle">Music and sports, loaded together. Find your next reason to go out.</p>
            </section>
            <button class="city-header-button">
                <p>Change city</p>
                <span>&nearr;</span>
            </button>
        </section>

    <!-- events sections -->
        {{range .events}}
        <section class="events-section">
            <section class="events-section-header">
                <h2 class="events-section-title">{{.Name}}</h2>
                <span class="events-section-count">{{len .Events}}</span>
            </section>

            {{with .Events}}
            <!-- event card container -->
            <section class="event-card-container">
                {{range .}}
                    {{template "event-card.tpl" .}}
                {{end}}
            </section>
            {{else}}
            <p class="events-empty">No events found in this category.</p>
            {{end}}
        </section>
        {{end}}
    </main>
    <script src="../static/js/index.js"></script>
</body>
</html>
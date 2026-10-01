<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Event Explorer</title>
    <link rel="stylesheet" href="../static/css/index.css">
</head>
<body>
    <nav>
        <section class="nav-left-section">
            <div class="nav-logo-box">e.</div>
            <p class="company-banner">eventexplorer</p>
        </section>

    </nav>

    <!-- preview notice strip -->
    <section class="notice-strip">
        <p><b>Interactive preview</b> / Fictional events. No live API calls or real purchases.</p>
    </section>

    <!-- main section starts here  -->
    <main class="main-layout">

        <!-- hero section -->
        <section class="hero-section">
            <section class="hero-left-section">
                <article class="hero-eyebrow">
                    <div class="hero-eyebrow-line"></div>
                    <p>LESS SCROLLING. MORE GOING.</p>
                </article>
                <!-- hero title -->
                <h1 class="hero-title">
                    A city of possibilities.
                    <span class="hero-title-green">Find your next one.</span>
                </h1>
                <p class="hero-subtitle">Discover music and sports in one place.</p>
                <p class="hero-subtitle">Choose your city. Find something worth heading out for.</p>
                <!-- hero tags -->
                <section class="hero-tag-container">
                    <span class="hero-tag">Live music</span>
                    <span class="hero-tag">Sports &amp; matchdays</span>
                    <span class="hero-tag">One simple search</span>
                </section>
            </section>

            <!-- hero poster card -->
            <section class="hero-right-section">
                <div class="hero-poster-shadow"></div>
                <article class="hero-poster">
                    <div class="hero-poster-top">
                        <p>THE CITY IS CALLING</p>
                        <p>01 / 02</p>
                    </div>
                    <div class="hero-poster-ring hero-poster-ring-outer"></div>
                    <div class="hero-poster-ring hero-poster-ring-inner"></div>
                    <h2 class="hero-poster-title">
                        <span class="hero-poster-title-go">GO</span>
                        <span class="hero-poster-title-out">OUT.</span>
                    </h2>
                    <div class="hero-poster-bottom">
                        <p>GOOD PLANS.</p>
                        <p>GREAT MEMORIES.</p>
                    </div>
                    <svg class="hero-poster-arrow" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512"><path d="M320 96c0-17.7 14.3-32 32-32l96 0c17.7 0 32 14.3 32 32l0 96c0 17.7-14.3 32-32 32s-32-14.3-32-32l0-18.7-233.4 233.4c-12.5 12.5-32.8 12.5-45.3 0s-12.5-32.8 0-45.3L370.7 128 352 128c-17.7 0-32-14.3-32-32z"/></svg>
                </article>
            </section>
        </section>

        <!-- city search section -->
        <section class="search-section">
            <section class="search-title-container">
                <h3 class="search-title">Where are we going?</h3>
                <p class="search-subtitle">Start with a city, then explore what is on.</p>
            </section>
            <p class="search-label">Choose a city</p>
            <section class="search-action-container">
                <div class="search-input-box">
                    <svg class="search-input-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 640 640"><path d="M320 576C178.6 576 64 461.4 64 320C64 178.6 178.6 64 320 64C461.4 64 576 178.6 576 320C576 461.4 461.4 576 320 576zM320 112C205.1 112 112 205.1 112 320C112 434.9 205.1 528 320 528C434.9 528 528 434.9 528 320C528 205.1 434.9 112 320 112zM320 416C267 416 224 373 224 320C224 267 267 224 320 224C373 224 416 267 416 320C416 373 373 416 320 416z" /></svg>
                    <input class="search-input" type="text" placeholder="Search a city.." value="">
                </div>
                <button class="search-button">
                    <p>Explore events</p>
                    <span>&rarr;</span>
                </button>
                <div class="suggested-cities" id="suggested-cities"></div>
            </section>
            <section class="search-hint-container">
                <p>Type at least 3 characters and select a suggestion.</p>
                <p>Local sample cities, not Google results</p>
            </section>
            <p class="search-hint-link">Choose a city from the suggestions.</p>
        </section>
    </main>
    <script src="../static/js/index.js"></script>
</body>
</html>
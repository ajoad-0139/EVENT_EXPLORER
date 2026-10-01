<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Event Explorer</title>
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

    <!-- main section starts here  -->
    <main class="main-layout">
        {{with index . "single-event"}}
        {{$segment := ""}}
        {{with .Classifications}}{{$segment = (index . 0).Segment.Name}}{{end}}

        <!-- event detail section -->
        <section class="event-detail-section">

            <!-- event detail left section -->
            <section class="event-detail-left-section">
                <div class="event-detail-banner {{if eq $segment "Sports"}}event-detail-banner-sports{{else}}event-detail-banner-music{{end}}">
                    {{if eq $segment "Sports"}}
                    <div class="event-detail-court">
                        <div class="event-detail-court-line"></div>
                        <div class="event-detail-court-box event-detail-court-box-left"></div>
                        <div class="event-detail-court-box event-detail-court-box-right"></div>
                    </div>
                    <svg class="event-detail-ball" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 200 200">
                        <circle class="event-detail-ball-base" cx="100" cy="100" r="98"/>
                        <circle class="event-detail-ball-stroke" cx="100" cy="100" r="78"/>
                        <line class="event-detail-ball-stroke" x1="100" y1="22" x2="100" y2="178"/>
                        <line class="event-detail-ball-stroke" x1="22" y1="100" x2="178" y2="100"/>
                        <path class="event-detail-ball-stroke" d="M 52 38 Q 84 100 52 162"/>
                        <path class="event-detail-ball-stroke" d="M 148 38 Q 116 100 148 162"/>
                    </svg>
                    {{else}}
                    <div class="event-detail-circle event-detail-circle-top"></div>
                    <div class="event-detail-circle event-detail-circle-bottom"></div>
                    <div class="event-detail-wave">
                        <span class="event-detail-bar"></span>
                        <span class="event-detail-bar"></span>
                        <span class="event-detail-bar"></span>
                        <span class="event-detail-bar"></span>
                        <span class="event-detail-bar"></span>
                        <span class="event-detail-bar"></span>
                    </div>
                    {{end}}
                    <div class="event-detail-dot"></div>
                    <div class="event-detail-cross"></div>
                    <p class="event-detail-tag">{{if $segment}}{{$segment}} event{{else}}Event{{end}}</p>
                </div>

                {{with $segment}}<p class="event-detail-category">{{.}}</p>{{end}}
                <h1 class="event-detail-title">{{.Name}}</h1>

                {{if or .Info .PleaseNote}}
                <section class="event-detail-about-section">
                    <h2 class="event-detail-about-title">About this event</h2>
                    {{with .Info}}<p class="event-detail-text">{{.}}</p>{{end}}
                    {{with .PleaseNote}}<p class="event-detail-text">{{.}}</p>{{end}}
                </section>
                {{end}}
            </section>

            <!-- event detail card -->
            <aside class="event-detail-card">
                <p class="event-detail-card-eyebrow">MAKE A PLAN</p>
                <h3 class="event-detail-card-title">The details</h3>

                <div class="event-detail-row">
                    <p class="event-detail-row-label">WHEN</p>
                    <p class="event-detail-row-value">{{formatEventDate .Dates.Start.LocalDate}}</p>
                    <p class="event-detail-row-note">{{formatEventTime .Dates.Start.LocalTime}}{{with .Dates.Timezone}} ({{.}}){{end}}</p>
                </div>

                <div class="event-detail-row">
                    <p class="event-detail-row-label">WHERE</p>
                    {{with .Embedded.Venues}}
                    {{with (index . 0)}}
                    <p class="event-detail-row-value">{{.Name}}</p>
                    <p class="event-detail-row-note">{{.City.Name}}, {{.Country.CountryCode}}</p>
                    {{end}}
                    {{else}}
                    <p class="event-detail-row-value">Venue to be announced</p>
                    {{end}}
                </div>

                {{with $segment}}
                <div class="event-detail-row">
                    <p class="event-detail-row-label">CATEGORY</p>
                    <p class="event-detail-row-value">{{.}}</p>
                </div>
                {{end}}

                <a class="event-detail-button" href="/redirect/{{.ID}}" target="_blank" rel="noopener noreferrer">
                    <span>View tickets</span>
                    <span>&nearr;</span>
                </a>
                <p class="event-detail-note">Tickets are sold by Ticketmaster. The link opens in a new tab.</p>
            </aside>
        </section>
        {{else}}
        <section class="event-detail-section">
            <p class="event-detail-empty">Event not found.</p>
        </section>
        {{end}}
    </main>
    <script src="../static/js/index.js"></script>
</body>
</html>
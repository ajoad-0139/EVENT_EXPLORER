<!-- event card: receives one EventResponse -->
<article class="event-card">
    <div class="event-card-banner">
        <div class="event-card-circle event-card-circle-top"></div>
        <div class="event-card-circle event-card-circle-bottom"></div>
        <div class="event-card-ring event-card-ring-outer"></div>
        <div class="event-card-ring event-card-ring-middle"></div>
        <div class="event-card-ring event-card-ring-inner"></div>
        <div class="event-card-cross"></div>
        <div class="event-card-wave">
            <span class="event-card-bar"></span>
            <span class="event-card-bar"></span>
            <span class="event-card-bar"></span>
            <span class="event-card-bar"></span>
            <span class="event-card-bar"></span>
            <span class="event-card-bar"></span>
        </div>
        <div class="event-card-dot"></div>
        <p class="event-card-tag">{{with .Classifications}}{{(index . 0).Segment.Name}} event{{else}}Event{{end}}</p>
    </div>
    <div class="event-card-body">
        <p class="event-card-date">{{formatEventDate .Dates.Start.LocalDate}}</p>
        <h3 class="event-card-title">{{.Name}}</h3>
        <p class="event-card-venue">{{with .Embedded.Venues}}{{(index . 0).Name}}{{end}}</p>
        <div class="event-card-footer">
            <p>{{with .Embedded.Venues}}{{(index . 0).City.Name}}{{end}}</p>
            <a class="event-card-link" href="" target="_blank" rel="noopener noreferrer">
                <span>View details</span>
                <span>&nearr;</span>
            </a>
        </div>
    </div>
</article>
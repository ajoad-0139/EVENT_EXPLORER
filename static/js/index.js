const searchInput = document.querySelector(".search-input");
const searchButton = document.querySelector(".search-button");
const suggestedCities = document.querySelector("#suggested-cities");

let debounceTimer;
let sessionToken = crypto.randomUUID();

let suggestions = [];

searchInput.addEventListener("input", (event) => {
    const query = event.target.value.trim();

    clearTimeout(debounceTimer);

    debounceTimer = setTimeout(async () => {
        if (!query) {
            suggestions = [];
            return;
        }

        try {
            const response = await fetch(
                `/api/locations/autocomplete?input=${encodeURIComponent(query)}&sessionToken=${sessionToken}`
            );

            if (!response.ok) {
                throw new Error(`Autocomplete failed: ${response.status}`);
            }

            suggestions = await response.json();

            console.log("Suggestions:", suggestions);
            console.log("Session:", sessionToken);
            renderSuggestions();
        } catch (error) {
            console.error(error);
        }
    }, 500);
});



function renderSuggestions() {
    suggestedCities.innerHTML = "";

    if (suggestions.length === 0) {
        suggestedCities.style.display = "none";
        return;
    }

    suggestedCities.style.display = "block";

    suggestions.forEach((suggestion) => {
        const city = document.createElement("div");

        city.classList.add("suggested-city");

        city.textContent = suggestion.text;

        city.addEventListener("click", () => {
            searchInput.value = suggestion.text;

            selectedPlace = suggestion;

            suggestedCities.innerHTML = "";
            suggestedCities.style.display = "none";

            console.log("Selected:", selectedPlace);
        });

        suggestedCities.appendChild(city);
    });
}


searchButton.addEventListener("click", async () => {
    if (!suggestions.length) {
        return;
    }

    const selectedPlace = suggestions[0];

    try {
        const response = await fetch(
            `/api/locations/${selectedPlace.placeId}?sessionToken=${sessionToken}`
        );

        if (!response.ok) {
            throw new Error(`Location details failed: ${response.status}`);
        }

        const location = await response.json();

        console.log("Location:", location);

        // The autocomplete → details session is finished.
        sessionToken = crypto.randomUUID();

        console.log("New session:", sessionToken);
    } catch (error) {
        console.error(error);
    }
});
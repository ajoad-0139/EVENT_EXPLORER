const searchInput = document.querySelector(".search-input");
const searchButton = document.querySelector(".search-button");
const suggestedCities = document.querySelector("#suggested-cities");
const changeCityBtn = document.querySelector(".city-header-button")

let debounceTimer;
let sessionToken = crypto.randomUUID();

let suggestions = [];

let selectedPlace = null;

if (changeCityBtn) {
    changeCityBtn.addEventListener("click", () => {
        window.location.href = "/";
    });
}

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
    }, 300);
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
    if (!selectedPlace) {
        return;
    }

    try {
        // Get city + countryCode from the selected place
        const response = await fetch(
            `/api/locations/${selectedPlace.placeId}?sessionToken=${sessionToken}`
        );

        if (!response.ok) {
            throw new Error(`Location details failed: ${response.status}`);
        }

        const location = await response.json();

        console.log("Location:", location);

        
        const params = new URLSearchParams({
            city: location.city,
            countryCode: location.countryCode
        });
        sessionToken = crypto.randomUUID();
        // Redirect
        window.location.href = `/events?${params.toString()}`;

    } catch (error) {
        console.error("Failed to get location details:", error);
    }
});


